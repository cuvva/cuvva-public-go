package app

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"

	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/lib/config"
	"github.com/cuvva/cuvva-public-go/lib/ptr"
	"github.com/cuvva/cuvva-public-go/lib/slicecontains"
	"github.com/cuvva/cuvva-public-go/tools/cdep"
	"github.com/cuvva/cuvva-public-go/tools/cdep/git"
	"github.com/cuvva/cuvva-public-go/tools/cdep/parsers"
	"github.com/cuvva/cuvva-public-go/tools/cdep/paths"
	gogit "github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var imageName = regexp.MustCompile(`"docker_image_name"\s*:\s*"([a-zA-Z\d_-]+)"`)

func (a App) Update(ctx context.Context, req *parsers.Params, overruleChecks []string) error {
	if req.Environment == "prod" && req.Branch != cdep.DefaultBranch {
		return cher.New("invalid_operation", nil)
	}

	log.Info("getting latest commit hash")

	latestHash, err := git.GetLatestCommitHash(ctx, req.Branch)
	if err != nil {
		return err
	}

	repoPath, err := paths.GetConfigRepo()
	if err != nil {
		return err
	}

	log.Info("fetching config repo")

	if out, err := exec.CommandContext(ctx, "git", "-C", repoPath, "fetch", "--all").Output(); err != nil {
		fmt.Println(string(out))
		return err
	}

	log.Info("opening config repo")

	configRepo, err := gogit.PlainOpen(repoPath)
	if err != nil {
		return cher.New("git_repo_error", cher.M{
			"error": err,
		})
	}

	_, err = git.CheckRepo(configRepo)
	if err != nil {
		return err
	}

	ref, err := configRepo.Head()
	if err != nil {
		return err
	}

	defaultRef := fmt.Sprintf("refs/heads/%s", cdep.DefaultBranch)
	if ref.Name().String() != defaultRef {
		return cher.New("config_not_on_mater", nil)
	}

	log.Info("pulling config repo from remote")

	if out, err := exec.CommandContext(ctx, "git", "-C", repoPath, "pull", "origin", cdep.DefaultBranch).Output(); err != nil {
		fmt.Println(string(out))
		return err
	}

	wt, err := configRepo.Worktree()
	if err != nil {
		return err
	}

	err = git.CheckWorkingCopy(wt)
	if err != nil {
		if !slicecontains.String(overruleChecks, "working_copy_dirty") {
			return err
		}

		log.Warn("working_copy_dirty overruled")
	}

	log.Info("adding hash and branch to json files")

	updatedFiles := []string{}

	envs, err := a.LoadEnvs(repoPath, req.System, req.Environment)
	if err != nil {
		return err
	}

	for env := range envs {
		switch req.Type {
		case "service":
			for _, service := range req.Items {
				p := paths.GetPathForService(repoPath, req.System, env, service)

				if _, err := os.Stat(p); err != nil {
					log.Warn(err)
				}

				err := checkECRImage(p, latestHash, req.Branch)
				if err != nil {
					e := errors.Wrap(err, "ecr:")
					log.Warn(e)
				}

				changed, err := a.AddToConfig(p, req.Branch, latestHash)
				if err != nil {
					return err
				}

				if changed {
					shorthandPath := path.Join(req.System, env, "service", service+".json")
					updatedFiles = append(updatedFiles, shorthandPath)
				}
			}
		case "lambda":
			for _, lambda := range req.Items {
				p := paths.GetPathForLambda(repoPath, req.System, env, lambda)

				if _, err := os.Stat(p); err != nil {
					log.Warn(err)
				}

				changed, err := a.AddToConfig(p, req.Branch, latestHash)
				if err != nil {
					return err
				}

				if changed {
					shorthandPath := path.Join(req.System, env, "lambda", lambda+".json")
					updatedFiles = append(updatedFiles, shorthandPath)
				}
			}
		default:
			return cher.New("unexpected_type", cher.M{"type": req.Type})
		}
	}

	if len(updatedFiles) == 0 {
		return cher.New("nothing_changed", nil)
	}

	commitMessage := fmt.Sprintf("cdep: %s", req.String("update"))

	if err := a.PublishToSlack(ctx, req, commitMessage, updatedFiles, repoPath); err != nil {
		return err
	}

	if a.DryRun {
		log.Info("Dry run only, stopping now")
		log.Infof("commit message (%s)\n", commitMessage)
		return nil
	}

	for _, p := range updatedFiles {
		log.Infof("adding %s to commit", p)
		_, err := wt.Add(p)
		if err != nil {
			return err
		}
	}

	_, err = wt.Commit(commitMessage, &gogit.CommitOptions{})
	if err != nil {
		return err
	}

	log.Info("pushing commit to config repo")

	if out, err := exec.CommandContext(ctx, "git", "-C", repoPath, "push", "origin", cdep.DefaultBranch).Output(); err != nil {
		fmt.Println(string(out))
		return err
	}

	return nil
}

func checkECRImage(filePath, latestHash, branch string) error {
	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	matches := imageName.FindSubmatch(fileContents)
	if len(matches) != 2 {
		return cher.New("invalid_docker_image_name", nil)
	}

	dockerImageName := matches[1]

	return findBuildInECR(string(dockerImageName), latestHash, branch)
}

func findBuildInECR(dockerImageName, latestHash, branch string) error {
	cfg := config.AWS{
		Region: "eu-west-1",
	}

	awsSession, err := cfg.Session()
	if err != nil {
		return errors.Wrap(err, "aws:")
	}

	c := ecr.New(awsSession)

	branchName := "master"

	if branch != "master" {
		branchName = "branch"
	}

	images, err := c.BatchGetImage(&ecr.BatchGetImageInput{
		RegistryId:     ptr.String("005717268539"),
		RepositoryName: ptr.String(dockerImageName),
		ImageIds: []*ecr.ImageIdentifier{
			{
				ImageTag: ptr.String(fmt.Sprintf("%s-%s", branchName, latestHash)),
			},
		},
	})
	if err != nil {
		return err
	}

	if len(images.Images) != 1 {
		log.Warn("Cannot find image in ECR!")
	}

	return nil
}

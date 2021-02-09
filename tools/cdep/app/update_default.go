package app

import (
	"context"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/lib/slicecontains"
	"github.com/cuvva/cuvva-public-go/tools/cdep"
	"github.com/cuvva/cuvva-public-go/tools/cdep/git"
	"github.com/cuvva/cuvva-public-go/tools/cdep/parsers"
	"github.com/cuvva/cuvva-public-go/tools/cdep/paths"
	gogit "github.com/go-git/go-git/v5"
	log "github.com/sirupsen/logrus"
)

func (a App) UpdateDefault(ctx context.Context, req *parsers.Params, overruleChecks []string) error {
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
		return cher.New("config_not_on_default", nil)
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

	envs, err := a.LoadEnvs(repoPath, req.System, req.Environment)
	if err != nil {
		return err
	}

	loadedPaths := map[string][]string{}

	for env := range envs {
		path := path.Join(repoPath, req.System, env, req.Type)
		loadedPaths[env] = append(loadedPaths[env], path)
	}

	updatedFiles := []string{}

	log.Info("editing json files")

	for env, paths := range loadedPaths {
		for _, p := range paths {
			files, err := ioutil.ReadDir(p)
			if err != nil {
				return err
			}

			for _, file := range files {
				if file.IsDir() {
					continue
				}

				fullPath := path.Join(p, file.Name())
				var changed bool

				if strings.Contains(fullPath, "_base.json") {
					changed, err = a.AddToConfig(fullPath, req.Branch, latestHash)
					if err != nil {
						return err
					}
				} else {
					changed, err = a.RemFromConfig(fullPath)
					if err != nil {
						e := cher.Coerce(err)
						if e.Code != "frozen" {
							return err
						}

						log.Warn(fmt.Sprintf("skipping %s due to cdep freeze", file.Name()))
					}
				}

				if changed {
					shortPath := path.Join(req.System, env, req.Type, file.Name())
					updatedFiles = append(updatedFiles, shortPath)
				}
			}
		}
	}

	if len(updatedFiles) == 0 {
		return cher.New("nothing_changed", nil)
	}

	user, err := exec.CommandContext(ctx, "git", "config", "user.name").Output()
	if err != nil {
		fmt.Println(string(user))
		return err
	}

	commitMessage := fmt.Sprintf("cdep: %s", req.String("update-default"))

	if !a.DryRun && req.System == "prod" {
		textTemplate := ":wrench: *command*: `%s`\n:technologist: *user*: `%s`"
		text := fmt.Sprintf(textTemplate, req.String("update"), strings.Split(string(user), "\n")[0])
		if req.Message != "" {
			text = text + fmt.Sprintf("\n\n:email: *message*: `%s`", req.Message)
		}

		arn := "arn:aws:sns:eu-west-1:005717268539:cuvva-deployments-prod"
		subject := "A prod deployment is happening"
		_, err := a.SNS.PublishWithContext(ctx, &sns.PublishInput{
			TopicArn: &arn,
			Subject:  &subject,
			Message:  &text,
		})
		if err, ok := err.(awserr.Error); !ok || err.Code() != "EndpointDisabled" {
			fmt.Println(err)
		}
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

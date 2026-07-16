package app

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/lib/config"
	"github.com/cuvva/cuvva-public-go/lib/ptr"
	"github.com/cuvva/cuvva-public-go/lib/slicecontains"
	"github.com/cuvva/cuvva-public-go/tools/cdep"
	"github.com/cuvva/cuvva-public-go/tools/cdep/git"
	"github.com/cuvva/cuvva-public-go/tools/cdep/parsers"
	"github.com/cuvva/cuvva-public-go/tools/cdep/paths"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var imageName = regexp.MustCompile(`"?docker_image_name"?\s*:\s*"?([a-zA-Z\d_-]+)"?`)

func (a App) Update(ctx context.Context, req *parsers.Params, overruleChecks []string) error {
	if req.Environment == "prod" && req.Branch != cdep.DefaultBranch {
		return cher.New("invalid_operation", nil)
	}

	// agentcore config only exists under the _system env, so reject anything else
	// (including "all") up front rather than failing later with a confusing
	// missing-file error while iterating per-service environments.
	if req.Type == "agentcore" && req.Environment != "_system" {
		return cher.New("agentcore_requires_system_env", nil)
	}

	// agentcore lives under the _system env, so it bypasses the env == "prod"
	// guard above; enforce the default branch for prod agentcore deploys directly.
	if req.Type == "agentcore" && req.System == "prod" && req.Branch != cdep.AgentcoreDefaultBranch {
		return cher.New("invalid_operation", nil)
	}

	// A single explicit --commit cannot be applied across agentcore items, as
	// each agent has its own source repo; deploy one agent at a time when pinning.
	if req.Type == "agentcore" && req.Commit != "" && len(req.Items) > 1 {
		return cher.New("commit_requires_single_agent", nil)
	}

	// agentcore items each declare their own source repo in their config, so the
	// commit is resolved per-item below rather than once from the monorepo here.
	if req.Commit == "" && req.Type != "agentcore" {
		log.Info("getting latest commit hash")
		latestHash, err := git.GetLatestCommitHash(ctx, req.Branch)
		if err != nil {
			return fmt.Errorf("failed to get commit hash: %w", err)
		}

		req.Commit = latestHash
	}

	// Validate commit hash is a full 40-character hash, not a short hash or branch name
	if req.Commit != "" {
		if err := cdep.ValidateCommitHash(req.Commit); err != nil {
			return err
		}
	}

	repoPath, err := paths.GetConfigRepo()
	if err != nil {
		return fmt.Errorf("path get config repo: %w", err)
	}

	log.Info("fetching config repo")

	if out, err := exec.CommandContext(ctx, "git", "-C", repoPath, "fetch", "--all").CombinedOutput(); err != nil {
		fmt.Println(string(out))
		return err
	}

	log.Info("checking config repo")

	_, err = git.CheckRepo(repoPath)
	if err != nil {
		return fmt.Errorf("config git check repo: %w", err)
	}

	ref, err := exec.CommandContext(ctx, "git", "-C", repoPath, "symbolic-ref", "HEAD").Output()
	if err != nil {
		return fmt.Errorf("config git head: %w", err)
	}

	defaultRef := fmt.Sprintf("refs/heads/%s", cdep.DefaultBranch)
	if strings.TrimSpace(string(ref)) != defaultRef {
		if !slicecontains.String(overruleChecks, "config_not_on_master") {
			return cher.New("config_not_on_master", nil)
		}

		log.Warn("config_not_on_master overruled")
	}

	log.Info("pulling config repo from remote")

	if out, err := exec.CommandContext(ctx, "git", "-C", repoPath, "pull").CombinedOutput(); err != nil {
		fmt.Println(string(out))
		return fmt.Errorf("git pull: %w", err)
	}

	err = git.CheckWorkingCopy(repoPath)
	if err != nil {
		if !slicecontains.String(overruleChecks, "working_copy_dirty") {
			return fmt.Errorf("config git check working copy: %w", err)
		}

		log.Warn("working_copy_dirty overruled")
	}

	log.Info("adding hash and branch to json files")

	updatedFiles := []string{}

	envs, err := a.LoadEnvs(repoPath, req.System, req.Environment)
	if err != nil {
		return fmt.Errorf("load envs: %w", err)
	}

	var envNames []string
	for env := range envs {
		envNames = append(envNames, env)

		switch req.Type {
		case "service":
			for _, service := range req.Items {
				p := paths.GetPathForService(repoPath, req.System, env, service)

				if _, err := os.Stat(p); err != nil {
					p = paths.GetYamlPathForService(repoPath, req.System, env, service)
					_, err2 := os.Stat(p)
					if err2 != nil {
						log.Warn(err)
						log.Warn(err2)
					}
				}

				err := checkECRImage(p, req.Commit, req.Branch)
				if err != nil {
					e := errors.Wrap(err, "ecr")
					log.Warn(e)
				}

				changed, err := a.AddToConfig(p, req.Branch, req.Commit)
				if err != nil {
					return fmt.Errorf("add to config: %w", err)
				}

				if changed {
					filename := path.Base(p)
					shorthandPath := path.Join(req.System, env, "service", filename)
					updatedFiles = append(updatedFiles, shorthandPath)
				}
			}
		case "lambda":
			for _, lambda := range req.Items {
				p := paths.GetPathForLambda(repoPath, req.System, env, lambda)

				if _, err := os.Stat(p); err != nil {
					log.Warn(err)
				}

				changed, err := a.AddToConfig(p, req.Branch, req.Commit)
				if err != nil {
					return err
				}

				if changed {
					shorthandPath := path.Join(req.System, env, "lambda", lambda+".json")
					updatedFiles = append(updatedFiles, shorthandPath)
				}
			}
		case "terra": // terraform
			for _, workspace := range req.Items {
				p := paths.GetPathForTerra(repoPath, req.System, env, workspace)

				if _, err := os.Stat(p); err != nil {
					log.Warn(err)
				}

				changed, err := a.AddToConfig(p, req.Branch, req.Commit)
				if err != nil {
					return err
				}

				if changed {
					shorthandPath := path.Join(req.System, env, "terra", workspace+".json")
					updatedFiles = append(updatedFiles, shorthandPath)
				}
			}
		case "agentcore":
			for _, item := range req.Items {
				p := paths.GetPathForAgentcore(repoPath, req.System, env, item)

				// Warn and skip a missing/misnamed item rather than aborting the
				// whole run; matches how AddToConfig tolerates missing files for
				// the other types, and avoids readAgentcoreRepo hard-erroring below.
				if _, err := os.Stat(p); err != nil {
					log.Warn(err)
					continue
				}

				// Each agent declares its own source repo, so resolve the latest
				// commit from that repo unless one was pinned with --commit.
				commit := req.Commit
				if commit == "" {
					repoURL, err := readAgentcoreRepo(p)
					if err != nil {
						return fmt.Errorf("read agentcore repo: %w", err)
					}

					log.Infof("resolving latest %s commit for %s from %s", req.Branch, item, repoURL)

					commit, err = git.GetLatestCommitHashForRepo(ctx, repoURL, req.Branch)
					if err != nil {
						return fmt.Errorf("get agentcore commit: %w", err)
					}

					if err := cdep.ValidateCommitHash(commit); err != nil {
						return err
					}
				}

				changed, err := a.AddToConfig(p, req.Branch, commit)
				if err != nil {
					return err
				}

				if changed {
					shorthandPath := path.Join(req.System, env, "agentcore", item+".json")
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
		return fmt.Errorf("publish to slack: %w", err)
	}

	dashboards := chooseDashboards(req, envNames)

	printDashboards(dashboards)

	if a.DryRun {
		log.Info("Dry run only, stopping now")
		log.Infof("commit message (%s)\n", commitMessage)
		return nil
	}

	for _, p := range updatedFiles {
		log.Infof("adding %s to commit", p)
		if out, err := exec.CommandContext(ctx, "git", "-C", repoPath, "add", p).CombinedOutput(); err != nil {
			fmt.Println(string(out))
			return fmt.Errorf("config git add: %w", err)
		}
	}

	if out, err := exec.CommandContext(ctx, "git", "-C", repoPath, "commit", "-m", commitMessage).CombinedOutput(); err != nil {
		fmt.Println(string(out))
		return fmt.Errorf("config git commit: %w", err)
	}

	log.Info("pushing commit to config repo")

	if out, err := exec.CommandContext(ctx, "git", "-C", repoPath, "push", "origin", "HEAD").CombinedOutput(); err != nil {
		fmt.Println(string(out))
		return fmt.Errorf("config git push: %w", err)
	}

	return nil
}

func checkECRImage(filePath, latestHash, branch string) error {
	fileContents, err := os.ReadFile(filePath)
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
		return fmt.Errorf("batch get image: %w", err)
	}

	if len(images.Images) != 1 {
		log.Warn("Cannot find image in ECR!")
	}

	return nil
}

package app

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/lib/config"
	"github.com/cuvva/cuvva-public-go/lib/ptr"
	"github.com/cuvva/cuvva-public-go/lib/slicecontains"
	"github.com/cuvva/cuvva-public-go/tools/cdep"
	"github.com/cuvva/cuvva-public-go/tools/cdep/git"
	"github.com/cuvva/cuvva-public-go/tools/cdep/helpers"
	"github.com/cuvva/cuvva-public-go/tools/cdep/parsers"
	"github.com/cuvva/cuvva-public-go/tools/cdep/paths"
	gogit "github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (a App) Update(ctx context.Context, req *parsers.Params, overruleChecks []string, goOnly, jsOnly bool) error {
	if req.Environment == "prod" && req.Branch != cdep.DefaultBranch {
		return cher.New("invalid_operation", nil)
	}

	if req.Commit == "" {
		log.Info("getting latest commit hash")
		latestHash, err := git.GetLatestCommitHash(ctx, req.Branch)
		if err != nil {
			return fmt.Errorf("failed to get commit hash: %w", err)
		}

		req.Commit = latestHash
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

	log.Info("opening config repo")

	configRepo, err := gogit.PlainOpen(repoPath)
	if err != nil {
		return cher.New("git_repo_error", cher.M{
			"error": err,
		})
	}

	_, err = git.CheckRepo(configRepo)
	if err != nil {
		return fmt.Errorf("config git check repo: %w", err)
	}

	ref, err := configRepo.Head()
	if err != nil {
		return fmt.Errorf("config git head: %w", err)
	}

	defaultRef := fmt.Sprintf("refs/heads/%s", cdep.DefaultBranch)
	if ref.Name().String() != defaultRef {
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

	wt, err := configRepo.Worktree()
	if err != nil {
		return fmt.Errorf("config git work tree: %w", err)
	}

	err = git.CheckWorkingCopy(wt)
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
			if len(req.Items) == 0 {
				// Handle "all" services - discover all service files
				err := a.updateAllServices(repoPath, req, env, &updatedFiles, goOnly, jsOnly)
				if err != nil {
					return err
				}
			} else {
				// Handle specific services
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
		_, err := wt.Add(p)
		if err != nil {
			return fmt.Errorf("config git add: %w", err)
		}
	}

	_, err = wt.Commit(commitMessage, &gogit.CommitOptions{})
	if err != nil {
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
	dockerImageName, err := helpers.ExtractDockerImageName(filePath)
	if err != nil {
		return err
	}

	if dockerImageName == "" {
		return cher.New("invalid_docker_image_name", nil)
	}

	return findBuildInECR(dockerImageName, latestHash, branch)
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

// updateAllServices discovers and updates all service files in the given environment
func (a App) updateAllServices(repoPath string, req *parsers.Params, env string, updatedFiles *[]string, goOnly, jsOnly bool) error {
	// Get the path for services in this environment
	servicePath := path.Join(repoPath, req.System, env, "service")

	// Read all files in the service directory
	files, err := os.ReadDir(servicePath)
	if err != nil {
		// If the service directory doesn't exist, that's okay - just return
		if os.IsNotExist(err) {
			log.Debugf("No service directory found for %s/%s, skipping", req.System, env)
			return nil
		}
		return fmt.Errorf("failed to read service directory %s: %w", servicePath, err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Skip non-service files (like _base.json, _default.json)
		if strings.HasPrefix(file.Name(), "_") {
			log.Debugf("Skipping system file %s", file.Name())
			continue
		}

		fullPath := path.Join(servicePath, file.Name())

		// Apply service filtering based on flags (only if filters are set)
		if goOnly || jsOnly {
			shouldUpdate, err := helpers.ShouldUpdateService(fullPath, goOnly, jsOnly)
			if err != nil {
				return fmt.Errorf("failed to check service filter for %s: %w", file.Name(), err)
			}

			if !shouldUpdate {
				log.Debugf("Skipping service %s due to filtering", file.Name())
				continue
			}
		}

		// Check ECR image
		err := checkECRImage(fullPath, req.Commit, req.Branch)
		if err != nil {
			e := errors.Wrap(err, "ecr")
			log.Warn(e)
		}

		// Update the service configuration
		changed, err := a.AddToConfig(fullPath, req.Branch, req.Commit)
		if err != nil {
			return fmt.Errorf("add to config for %s: %w", file.Name(), err)
		}

		if changed {
			shorthandPath := path.Join(req.System, env, "service", file.Name())
			*updatedFiles = append(*updatedFiles, shorthandPath)
			log.Infof("Updated service %s", file.Name())
		}
	}

	return nil
}
package app

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/lib/slicecontains"
	"github.com/cuvva/cuvva-public-go/tools/cdep"
	"github.com/cuvva/cuvva-public-go/tools/cdep/git"
	"github.com/cuvva/cuvva-public-go/tools/cdep/helpers"
	"github.com/cuvva/cuvva-public-go/tools/cdep/parsers"
	"github.com/cuvva/cuvva-public-go/tools/cdep/paths"
	gogit "github.com/go-git/go-git/v5"
	log "github.com/sirupsen/logrus"
)



func (a App) UpdateDefault(ctx context.Context, req *parsers.Params, overruleChecks []string, goOnly, jsOnly bool) error {
	log.Info("getting latest commit hash")

	if req.Commit == "" {
		latestHash, err := git.GetLatestCommitHash(ctx, req.Branch)
		if err != nil {
			return err
		}

		req.Commit = latestHash
	}

	repoPath, err := paths.GetConfigRepo()
	if err != nil {
		return err
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

	if out, err := exec.CommandContext(ctx, "git", "-C", repoPath, "pull", "origin", cdep.DefaultBranch).CombinedOutput(); err != nil {
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
			files, err := os.ReadDir(p)
			if err != nil {
				return err
			}

			for _, file := range files {
				if file.IsDir() {
					continue
				}

				fullPath := path.Join(p, file.Name())

				// Apply service filtering based on flags (skip filtering for _base.json)
				if file.Name() != "_base.json" {
					shouldUpdate, err := helpers.ShouldUpdateService(fullPath, goOnly, jsOnly)
					if err != nil {
						return err
					}

					if !shouldUpdate {
						log.Debugf("Skipping service %s due to filtering", file.Name())
						continue
					}
				} else {
					log.Debugf("Processing base configuration file %s (skipping docker image filtering)", file.Name())
				}

				changed, err := a.AddToConfig(fullPath, req.Branch, req.Commit)
				if err != nil {
					return err
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

	commitMessage := fmt.Sprintf("cdep: %s", req.String("update-default"))

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

	if out, err := exec.CommandContext(ctx, "git", "-C", repoPath, "push", "origin", cdep.DefaultBranch).CombinedOutput(); err != nil {
		fmt.Println(string(out))
		return err
	}

	return nil
}

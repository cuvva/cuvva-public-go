package app

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/lib/slicecontains"
	"github.com/cuvva/cuvva-public-go/tools/cdep"
	"github.com/cuvva/cuvva-public-go/tools/cdep/git"
	"github.com/cuvva/cuvva-public-go/tools/cdep/parsers"
	"github.com/cuvva/cuvva-public-go/tools/cdep/paths"
	log "github.com/sirupsen/logrus"
)

func (a App) UpdateDefault(ctx context.Context, req *parsers.Params, overruleChecks []string) error {
	if err := validateServiceSelection(req); err != nil {
		return err
	}

	log.Info("getting latest commit hash")

	if req.Commit == "" {
		latestHash, err := git.GetLatestCommitHash(ctx, req.Branch)
		if err != nil {
			return err
		}

		req.Commit = latestHash
	}

	// Validate commit hash is a full 40-character hash, not a short hash or branch name
	if err := cdep.ValidateCommitHash(req.Commit); err != nil {
		return err
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

	log.Info("checking config repo")

	_, err = git.CheckRepo(repoPath)
	if err != nil {
		return err
	}

	ref, err := exec.CommandContext(ctx, "git", "-C", repoPath, "symbolic-ref", "HEAD").Output()
	if err != nil {
		return err
	}

	defaultRef := fmt.Sprintf("refs/heads/%s", cdep.DefaultBranch)
	if strings.TrimSpace(string(ref)) != defaultRef {
		return cher.New("config_not_on_default", nil)
	}

	log.Info("pulling config repo from remote")

	if out, err := exec.CommandContext(ctx, "git", "-C", repoPath, "pull", "origin", cdep.DefaultBranch).CombinedOutput(); err != nil {
		fmt.Println(string(out))
		return err
	}

	err = git.CheckWorkingCopy(repoPath)
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

	updatedFiles := []string{}
	ecr := newECRChecker(req.Commit, req.Branch)

	log.Info("editing json files")

	for env := range envs {
		dir := path.Join(repoPath, req.System, env, req.Type)

		if req.Type == "service" {
			files, err := listServiceFiles(dir, true)
			if err != nil {
				return err
			}

			selected, skipped := filterServiceFiles(env, req.Only, files)

			for _, f := range selected {
				ecr.Check(f)

				changed, err := a.AddToConfig(f.Path, req.Branch, req.Commit)
				if err != nil {
					return err
				}

				if changed {
					updatedFiles = append(updatedFiles, path.Join(req.System, env, req.Type, f.Name()))

					if f.IsEnvFile() {
						warnInheritedVersions(env, f, skipped)
					}
				}
			}

			continue
		}

		files, err := os.ReadDir(dir)
		if err != nil {
			return err
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			changed, err := a.AddToConfig(path.Join(dir, file.Name()), req.Branch, req.Commit)
			if err != nil {
				return err
			}

			if changed {
				updatedFiles = append(updatedFiles, path.Join(req.System, env, req.Type, file.Name()))
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
		logDryRun(updatedFiles)
		log.Info("Dry run only, stopping now")
		log.Infof("commit message (%s)\n", commitMessage)
		return nil
	}

	for _, p := range updatedFiles {
		log.Infof("adding %s to commit", p)
		if out, err := exec.CommandContext(ctx, "git", "-C", repoPath, "add", p).CombinedOutput(); err != nil {
			fmt.Println(string(out))
			return err
		}
	}

	if out, err := exec.CommandContext(ctx, "git", "-C", repoPath, "commit", "-m", commitMessage).CombinedOutput(); err != nil {
		fmt.Println(string(out))
		return err
	}

	log.Info("pushing commit to config repo")

	if out, err := exec.CommandContext(ctx, "git", "-C", repoPath, "push", "origin", cdep.DefaultBranch).CombinedOutput(); err != nil {
		fmt.Println(string(out))
		return err
	}

	return nil
}

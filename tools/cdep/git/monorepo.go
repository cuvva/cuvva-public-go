package git

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/tools/cdep/paths"
	log "github.com/sirupsen/logrus"
)

func GetLatestCommitHash(ctx context.Context, branchName string) (string, error) {
	monorepoPath, err := paths.GetCodeRepo()
	if err != nil {
		return "", err
	}

	remote, err := CheckRepo(monorepoPath)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	branchRef := fmt.Sprintf("refs/heads/%s", branchName)

	out, err := exec.CommandContext(ctx, "git", "-C", monorepoPath, "ls-remote", remote, branchRef).Output()
	if err != nil {
		log.Error(err)

		return "", cher.New("git_repo_error", cher.M{
			"error": gitErrString(err),
		})
	}

	// each line is "<hash>\t<ref name>"
	for line := range strings.SplitSeq(strings.TrimSpace(string(out)), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[1] == branchRef {
			return fields[0], nil
		}
	}

	return "", cher.New("remote_branch_not_found", nil)
}

// FetchBranch updates the code repo's remote-tracking ref for branchName, so
// ancestry checks see commits other people deployed since the last local fetch.
func FetchBranch(ctx context.Context, branchName string) error {
	monorepoPath, err := paths.GetCodeRepo()
	if err != nil {
		return err
	}

	remote, err := CheckRepo(monorepoPath)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	if _, err := exec.CommandContext(ctx, "git", "-C", monorepoPath, "fetch", remote, branchName).Output(); err != nil {
		return cher.New("git_repo_error", cher.M{
			"error": gitErrString(err),
		})
	}

	return nil
}

// IsAncestor reports whether ancestor is reachable from descendant in the code
// repo's history. A commit the local repo does not know is an error, not false.
func IsAncestor(ctx context.Context, ancestor, descendant string) (bool, error) {
	monorepoPath, err := paths.GetCodeRepo()
	if err != nil {
		return false, err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	err = exec.CommandContext(ctx, "git", "-C", monorepoPath, "merge-base", "--is-ancestor", ancestor, descendant).Run()
	if err == nil {
		return true, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
		return false, nil
	}

	return false, cher.New("git_repo_error", cher.M{
		"error": gitErrString(err),
	})
}

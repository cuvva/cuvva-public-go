package git

import (
	"context"
	"fmt"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/tools/cdep/paths"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func GetLatestCommitHash(ctx context.Context, branchName string) (string, error) {
	monorepoPath, err := paths.GetCodeRepo()
	if err != nil {
		return "", err
	}

	monorepo, err := gogit.PlainOpen(monorepoPath)
	if err != nil {
		return "", cher.New("git_repo_error", cher.M{
			"error": err,
		})
	}

	remote, err := CheckRepo(monorepo)
	if err != nil {
		return "", err
	}

	err = remote.FetchContext(ctx, &gogit.FetchOptions{})
	if err != nil {
		if err != gogit.NoErrAlreadyUpToDate {
			return "", err
		}
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	refs, err := remote.ListContext(ctx, &gogit.ListOptions{})
	if err != nil {
		return "", err
	}

	var latestRefForBranch *plumbing.Reference

	for _, ref := range refs {
		if !ref.Name().IsBranch() {
			continue
		}

		refName := ref.Name().String()
		if refName == fmt.Sprintf("refs/heads/%s", branchName) {
			latestRefForBranch = ref
		}
	}

	if latestRefForBranch == nil {
		return "", cher.New("remote_branch_not_found", nil)
	}

	return latestRefForBranch.Hash().String(), nil
}

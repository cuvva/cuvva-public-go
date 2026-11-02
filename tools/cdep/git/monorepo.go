package git

import (
	"context"
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

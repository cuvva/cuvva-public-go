package git

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/cher"
)

// lsRemoteTimeout bounds a single `git ls-remote` call so a network stall
// against an agent's repo cannot hang the whole deploy indefinitely.
const lsRemoteTimeout = 30 * time.Second

// GetLatestCommitHashForRepo returns the tip commit of branchName for an
// arbitrary remote repo identified by repoURL, so each agentcore item can
// resolve its own source repo dynamically. No local clone required.
//
// It shells out to `git ls-remote` (the same way cdep already runs config-repo
// fetch/pull/push) rather than using go-git in-process, so it uses the caller's
// normal git and SSH configuration — identity files, ssh agent, known_hosts —
// instead of go-git's agent-only auth.
func GetLatestCommitHashForRepo(ctx context.Context, repoURL, branchName string) (string, error) {
	target := "refs/heads/" + branchName

	ctx, cancel := context.WithTimeout(ctx, lsRemoteTimeout)
	defer cancel()

	out, err := exec.CommandContext(ctx, "git", "ls-remote", repoURL, target).Output()
	if err != nil {
		return "", cher.New("git_repo_error", cher.M{
			"repo":  repoURL,
			"error": err.Error(),
		})
	}

	// Each line is "<hash>\t<ref>"; match the branch ref exactly rather than
	// trusting ls-remote's loose pattern matching.
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[1] == target {
			return fields[0], nil
		}
	}

	return "", cher.New("remote_branch_not_found", cher.M{
		"repo":   repoURL,
		"branch": branchName,
	})
}

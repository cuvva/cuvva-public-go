package git

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/cuvva/cuvva-public-go/lib/cher"
)

// CheckRepo validates the repo at repoPath has exactly one remote configured
// over SSH, and returns the remote's name
func CheckRepo(repoPath string) (string, error) {
	out, err := exec.Command("git", "-C", repoPath, "remote").Output()
	if err != nil {
		return "", cher.New("git_repo_error", cher.M{
			"error": gitErrString(err),
		})
	}

	remotes := strings.Fields(string(out))

	if len(remotes) > 1 {
		return "", cher.New("multiple_remotes", nil)
	}

	if len(remotes) == 0 {
		return "", cher.New("no_remotes", nil)
	}

	remote := remotes[0]

	remoteURL, err := exec.Command("git", "-C", repoPath, "remote", "get-url", remote).Output()
	if err != nil {
		return "", cher.New("git_repo_error", cher.M{
			"error": gitErrString(err),
		})
	}

	if !strings.Contains(string(remoteURL), "git@") {
		return "", errors.New("cuvva repo remote origin url is not ssh")
	}

	return remote, nil
}

// gitErrString returns git's stderr where available, as the exec error alone
// only reports the exit status
func gitErrString(err error) string {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && len(exitErr.Stderr) > 0 {
		return strings.TrimSpace(string(exitErr.Stderr))
	}

	return err.Error()
}

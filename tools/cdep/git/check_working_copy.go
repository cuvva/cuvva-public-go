package git

import (
	"os/exec"
	"strings"

	"github.com/cuvva/cuvva-public-go/lib/cher"
)

// CheckWorkingCopy checks the working copy to see if its dirty
func CheckWorkingCopy(repoPath string) error {
	out, err := exec.Command("git", "-C", repoPath, "status", "--porcelain").Output()
	if err != nil {
		return err
	}

	status := strings.TrimSpace(string(out))
	if status == "" {
		return nil
	}

	return cher.New("working_copy_dirty", cher.M{
		"status": status,
	})
}

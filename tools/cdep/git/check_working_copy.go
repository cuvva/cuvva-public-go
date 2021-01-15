package git

import (
	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/go-git/go-billy/v5/osfs"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
)

// CheckWorkingCopy checks the working copy to see if its dirty
func CheckWorkingCopy(wt *gogit.Worktree) error {
	fs := osfs.New("/")

	patterns, err := gitignore.LoadGlobalPatterns(fs)
	if err != nil {
		return err
	}

	wt.Excludes = append(wt.Excludes, patterns...)

	patterns, err = gitignore.LoadSystemPatterns(fs)
	if err != nil {
		return err
	}

	wt.Excludes = append(wt.Excludes, patterns...)

	status, err := wt.Status()
	if err != nil {
		return err
	}

	if len(status) == 0 {
		return nil
	}

	return cher.New("working_copy_dirty", cher.M{
		"status": status,
	})
}

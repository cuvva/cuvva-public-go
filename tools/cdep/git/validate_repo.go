package git

import (
	"github.com/cuvva/cuvva-public-go/lib/cher"
	gogit "github.com/go-git/go-git/v5"
)

func CheckRepo(repo *gogit.Repository) (*gogit.Remote, error) {
	remotes, err := repo.Remotes()
	if err != nil {
		return nil, err
	}

	if len(remotes) > 1 {
		return nil, cher.New("multiple_remotes", nil)
	}

	return remotes[0], nil
}

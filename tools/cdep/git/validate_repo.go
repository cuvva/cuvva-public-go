package git

import (
	"errors"
	"github.com/cuvva/cuvva-public-go/lib/cher"
	gogit "github.com/go-git/go-git/v5"
	"strings"
)

func CheckRepo(repo *gogit.Repository) (*gogit.Remote, error) {
	remotes, err := repo.Remotes()
	if err != nil {
		return nil, err
	}

	if len(remotes) > 1 {
		return nil, cher.New("multiple_remotes", nil)
	}

	if len(remotes) == 0 {
		return nil, cher.New("no_remotes", nil)
	}

	remote := remotes[0]

	if !isRemoteURLSSH(remote) {
		return nil, errors.New("cuvva repo remote origin must be ssh url")
	}

	return remotes[0], nil
}

func isRemoteURLSSH(remote *gogit.Remote) bool {
	remoteURL := remote.Config().URLs[0]

	return strings.Contains(remoteURL, "git@")
}

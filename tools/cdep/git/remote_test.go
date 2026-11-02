package git

import (
	"context"
	"regexp"
	"testing"

	"github.com/cuvva/cuvva-public-go/lib/cher"
)

var hex40 = regexp.MustCompile(`^[a-f0-9]{40}$`)

// Uses a public repo over HTTPS and skips when the network is unavailable, so
// it never fails CI offline and never touches private infrastructure.
func TestGetLatestCommitHashForRepo(t *testing.T) {
	const repo = "https://github.com/go-git/go-git.git"

	ctx := context.Background()

	hash, err := GetLatestCommitHashForRepo(ctx, repo, "master")
	if err != nil {
		t.Skipf("network unavailable, skipping live ls-remote: %v", err)
	}
	if !hex40.MatchString(hash) {
		t.Fatalf("expected 40-char hex hash, got %q", hash)
	}

	_, err = GetLatestCommitHashForRepo(ctx, repo, "definitely-not-a-real-branch-xyz")
	if code := cher.Coerce(err).Code; code != "remote_branch_not_found" {
		t.Fatalf("expected remote_branch_not_found, got %q (%v)", code, err)
	}
}

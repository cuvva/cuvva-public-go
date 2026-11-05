package app

import (
	"context"

	"github.com/cuvva/cuvva-public-go/tools/cdep"
	"github.com/cuvva/cuvva-public-go/tools/cdep/git"
	log "github.com/sirupsen/logrus"
)

// rollbackCheck confirms that an older master commit may replace the one a
// service is running. A rollout over many services meets the same running
// commit again and again, so each answer is kept and reused; a different
// running commit gets its own check.
type rollbackCheck struct {
	isAncestor func(ancestor, descendant string) (bool, error)
	ask        func(label string) bool
	decisions  map[string]bool
}

func newRollbackCheck(ctx context.Context) *rollbackCheck {
	var fetched bool

	return &rollbackCheck{
		isAncestor: func(ancestor, descendant string) (bool, error) {
			if !fetched {
				fetched = true

				if err := git.FetchBranch(ctx, cdep.DefaultBranch); err != nil {
					return false, err
				}
			}

			return git.IsAncestor(ctx, ancestor, descendant)
		},
		ask:       boolPrompt,
		decisions: map[string]bool{},
	}
}

// Confirm reports whether newCommit may replace running at path. It asks only
// when newCommit is older than running, or when the two cannot be compared, and
// at most once per pair of commits. A nil check never asks.
func (r *rollbackCheck) Confirm(running, newCommit, path string) bool {
	if r == nil || running == newCommit {
		return true
	}

	key := running + ">" + newCommit

	if ok, seen := r.decisions[key]; seen {
		if !ok {
			log.Infof("Keeping newer master commit %s, path: %s", running, path)
		}

		return ok
	}

	older, err := r.isAncestor(newCommit, running)

	switch {
	case err != nil:
		log.Warnf("Cannot tell whether master commit %s is older than the running %s: %v, path: %s", newCommit, running, err, path)
	case !older:
		r.decisions[key] = true

		return true
	default:
		log.Warnf("Older master commit is going to be deployed, running: %s, new: %s, path: %s", running, newCommit, path)
	}

	ok := r.ask("Deploy older master commit (rollback)?")
	r.decisions[key] = ok

	return ok
}

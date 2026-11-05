package app

import (
	"errors"
	"testing"
)

const (
	runningCommit = "1111111111111111111111111111111111111111"
	olderCommit   = "2222222222222222222222222222222222222222"
	newerCommit   = "3333333333333333333333333333333333333333"
	otherRunning  = "4444444444444444444444444444444444444444"
)

// fakeRollbackCheck treats olderCommit as an ancestor of everything, newerCommit
// as an ancestor of nothing, and otherRunning as unknown to git.
func fakeRollbackCheck(answer bool) (*rollbackCheck, *int, *int) {
	compares, asks := 0, 0

	r := &rollbackCheck{
		isAncestor: func(ancestor, descendant string) (bool, error) {
			compares++

			if descendant == otherRunning {
				return false, errors.New("unknown commit")
			}

			return ancestor == olderCommit, nil
		},
		ask: func(string) bool {
			asks++

			return answer
		},
		decisions: map[string]bool{},
	}

	return r, &compares, &asks
}

func TestRollbackCheckConfirm(t *testing.T) {
	t.Run("same commit never compares or asks", func(t *testing.T) {
		r, compares, asks := fakeRollbackCheck(false)

		if !r.Confirm(runningCommit, runningCommit, "a.yaml") {
			t.Fatal("same commit must be allowed")
		}

		if *compares != 0 || *asks != 0 {
			t.Fatalf("compares=%d asks=%d, want 0 and 0", *compares, *asks)
		}
	})

	t.Run("newer commit is allowed without asking", func(t *testing.T) {
		r, _, asks := fakeRollbackCheck(false)

		if !r.Confirm(runningCommit, newerCommit, "a.yaml") {
			t.Fatal("newer commit must be allowed")
		}

		if *asks != 0 {
			t.Fatalf("asks=%d, want 0", *asks)
		}
	})

	t.Run("older commit asks once per running commit", func(t *testing.T) {
		r, compares, asks := fakeRollbackCheck(true)

		for _, path := range []string{"a.yaml", "b.yaml", "c.json"} {
			if !r.Confirm(runningCommit, olderCommit, path) {
				t.Fatalf("%s: confirmed rollback must be allowed", path)
			}
		}

		if *compares != 1 || *asks != 1 {
			t.Fatalf("compares=%d asks=%d, want 1 and 1", *compares, *asks)
		}
	})

	t.Run("refused rollback is remembered", func(t *testing.T) {
		r, _, asks := fakeRollbackCheck(false)

		if r.Confirm(runningCommit, olderCommit, "a.yaml") {
			t.Fatal("refused rollback must not be allowed")
		}

		if r.Confirm(runningCommit, olderCommit, "b.yaml") {
			t.Fatal("refusal must carry over to the next service")
		}

		if *asks != 1 {
			t.Fatalf("asks=%d, want 1", *asks)
		}
	})

	t.Run("a different running commit is checked on its own", func(t *testing.T) {
		r, compares, asks := fakeRollbackCheck(true)

		r.Confirm(runningCommit, olderCommit, "a.yaml")
		r.Confirm(otherRunning, olderCommit, "b.yaml")
		r.Confirm(runningCommit, olderCommit, "c.yaml")
		r.Confirm(otherRunning, olderCommit, "d.yaml")

		if *compares != 2 || *asks != 2 {
			t.Fatalf("compares=%d asks=%d, want 2 and 2", *compares, *asks)
		}
	})

	t.Run("an unknown commit asks rather than guessing", func(t *testing.T) {
		r, _, asks := fakeRollbackCheck(true)

		if !r.Confirm(otherRunning, newerCommit, "a.yaml") {
			t.Fatal("confirmed deploy must be allowed")
		}

		if *asks != 1 {
			t.Fatalf("asks=%d, want 1", *asks)
		}
	})

	t.Run("nil check allows everything", func(t *testing.T) {
		var r *rollbackCheck

		if !r.Confirm(runningCommit, olderCommit, "a.yaml") {
			t.Fatal("nil check must allow the deploy")
		}
	})
}

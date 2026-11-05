package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// AddToConfig rewrites the branch value in a service's JSON or YAML config.
// These tests cover the two branch shapes cdep sees in practice: the common
// ticket-style names (e.g. "sc-1234-fix-thing", which must keep working) and
// path-separated names (e.g. "claude/my-feature", which used to be truncated at
// the first "/" in YAML configs, leaving the wrong branch and corrupting the
// value on redeploy).
func TestAddToConfigBranchNames(t *testing.T) {
	const commit = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	tests := []struct {
		name    string
		file    string
		initial string
		branch  string
		want    string // exact branch line/token expected in the output
	}{
		{
			name:    "json ticket-style branch over master",
			file:    "service.json",
			initial: "{\n\t\"branch\": \"master\",\n\t\"commit\": \"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\"\n}\n",
			branch:  "sc-1234-fix-thing",
			want:    "\"branch\": \"sc-1234-fix-thing\"",
		},
		{
			name:    "yaml ticket-style branch over master",
			file:    "service.yaml",
			initial: "branch: master\ncommit: bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\n",
			branch:  "sc-1234-fix-thing",
			want:    "branch: sc-1234-fix-thing",
		},
		{
			name:    "yaml redeploy same ticket-style branch",
			file:    "service.yaml",
			initial: "branch: sc-1234-fix-thing\ncommit: bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\n",
			branch:  "sc-1234-fix-thing",
			want:    "branch: sc-1234-fix-thing",
		},
		{
			name:    "yaml dotted branch over master",
			file:    "service.yaml",
			initial: "branch: master\ncommit: bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\n",
			branch:  "release-1.2.3",
			want:    "branch: release-1.2.3",
		},
		{
			name:    "json slashed branch over master",
			file:    "service.json",
			initial: "{\n\t\"branch\": \"master\",\n\t\"commit\": \"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\"\n}\n",
			branch:  "claude/cdep-branch-name-slashes",
			want:    "\"branch\": \"claude/cdep-branch-name-slashes\"",
		},
		{
			name:    "yaml slashed branch over master",
			file:    "service.yaml",
			initial: "branch: master\ncommit: bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\n",
			branch:  "claude/cdep-branch-name-slashes",
			want:    "branch: claude/cdep-branch-name-slashes",
		},
		{
			name:    "yaml redeploy same slashed branch is not corrupted",
			file:    "service.yaml",
			initial: "branch: claude/cdep-branch-name-slashes\ncommit: bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\n",
			branch:  "claude/cdep-branch-name-slashes",
			want:    "branch: claude/cdep-branch-name-slashes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), tt.file)
			if err := os.WriteFile(path, []byte(tt.initial), 0o600); err != nil {
				t.Fatal(err)
			}

			if _, err := (App{}).AddToConfig(path, tt.branch, commit); err != nil {
				t.Fatalf("AddToConfig: %v", err)
			}

			out, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			got := string(out)

			if !strings.Contains(got, tt.want) {
				t.Fatalf("config missing %q\n---\n%s", tt.want, got)
			}
			// Guard against the truncation/corruption bug: the branch must not
			// gain a stray trailing segment (e.g. "claude/foo" -> "claude/foo/…").
			if strings.Contains(got, tt.branch+"/") {
				t.Fatalf("branch value %q was corrupted:\n%s", tt.branch, got)
			}
			// The old default must be gone once a feature branch is written.
			if strings.Contains(got, "master") {
				t.Fatalf("stale master branch left in config:\n%s", got)
			}
		})
	}
}

// A master-to-master deploy goes through the rollback check with the commit the
// file currently holds, and a refusal leaves the file untouched.
func TestAddToConfigMasterRollback(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		initial string
	}{
		{
			name:    "json",
			file:    "service.json",
			initial: "{\n\t\"branch\": \"master\",\n\t\"commit\": \"" + runningCommit + "\"\n}\n",
		},
		{
			name:    "yaml",
			file:    "service.yaml",
			initial: "branch: master\ncommit: " + runningCommit + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, answer := range []bool{true, false} {
				path := filepath.Join(t.TempDir(), tt.file)
				if err := os.WriteFile(path, []byte(tt.initial), 0o600); err != nil {
					t.Fatal(err)
				}

				r, _, asks := fakeRollbackCheck(answer)

				changed, err := (App{rollbacks: r}).AddToConfig(path, "master", olderCommit)
				if err != nil {
					t.Fatalf("AddToConfig: %v", err)
				}

				out, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}

				if *asks != 1 {
					t.Fatalf("answer=%v: asks=%d, want 1", answer, *asks)
				}

				if changed != answer {
					t.Fatalf("answer=%v: changed=%v", answer, changed)
				}

				if strings.Contains(string(out), olderCommit) != answer {
					t.Fatalf("answer=%v: unexpected file contents:\n%s", answer, out)
				}
			}
		})
	}
}

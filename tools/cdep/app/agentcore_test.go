package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cuvva/cuvva-public-go/lib/cher"
)

func TestReadAgentcoreRepo(t *testing.T) {
	dir := t.TempDir()

	write := func(name, body string) string {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		return p
	}

	good := write("good.json", `{
	"repo": "git@github.com:example-org/example-agent.git",
	"branch": "main",
	"commit": "0000000000000000000000000000000000000000"
}`)
	repo, err := readAgentcoreRepo(good)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo != "git@github.com:example-org/example-agent.git" {
		t.Fatalf("repo = %q", repo)
	}

	missing := write("missing.json", `{"branch":"main","commit":"0000000000000000000000000000000000000000"}`)
	if _, err := readAgentcoreRepo(missing); cher.Coerce(err).Code != "missing_agentcore_repo" {
		t.Fatalf("expected missing_agentcore_repo, got %v", err)
	}

	invalid := write("invalid.json", `{not valid json`)
	if _, err := readAgentcoreRepo(invalid); cher.Coerce(err).Code != "invalid_agentcore_config" {
		t.Fatalf("expected invalid_agentcore_config, got %v", err)
	}
}

// AddToConfig is shared with the other types; this asserts it behaves correctly
// on an agentcore-shaped file (branch "main", 40-zero commit placeholder).
func TestAddToConfigAgentcore(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "agent.json")
	os.WriteFile(p, []byte("{\n\t\"repo\": \"git@github.com:example-org/example-agent.git\",\n\t\"branch\": \"main\",\n\t\"commit\": \"0000000000000000000000000000000000000000\",\n\t\"env\": {\n\t\t\"MODE\": \"a2a\"\n\t}\n}\n"), 0o644)

	hash := "7fbcf728fe0682e62f3c5e179ebf2e846a99d3c2"
	changed, err := App{}.AddToConfig(p, "main", hash)
	if err != nil {
		t.Fatalf("AddToConfig error: %v", err)
	}
	if !changed {
		t.Fatal("expected changed=true")
	}

	out, _ := os.ReadFile(p)
	s := string(out)
	if !strings.Contains(s, `"commit": "`+hash+`"`) {
		t.Fatalf("commit not updated:\n%s", s)
	}
	if !strings.Contains(s, `"branch": "main"`) {
		t.Fatalf("branch not preserved:\n%s", s)
	}
	if !strings.Contains(s, `"MODE": "a2a"`) {
		t.Fatalf("unrelated fields clobbered:\n%s", s)
	}

	// second run with the same values is a no-op
	changed, err = App{}.AddToConfig(p, "main", hash)
	if err != nil {
		t.Fatalf("AddToConfig (rerun) error: %v", err)
	}
	if changed {
		t.Fatal("expected no change on rerun")
	}
}

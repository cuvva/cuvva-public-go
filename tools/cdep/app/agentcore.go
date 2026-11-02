package app

import (
	"encoding/json"
	"os"

	"github.com/cuvva/cuvva-public-go/lib/cher"
)

// agentcoreConfig captures the fields cdep reads from an agentcore config file.
// The file itself is rewritten by AddToConfig's regex logic (not by
// re-marshalling this struct), so any other fields are left untouched.
type agentcoreConfig struct {
	Repo string `json:"repo"`
}

// readAgentcoreRepo extracts the source repo declared in an agentcore config
// file. Each agent lives in its own repo, so the repo is read per-file to
// resolve that repo's latest commit rather than assuming the monorepo.
func readAgentcoreRepo(path string) (string, error) {
	blob, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	var cfg agentcoreConfig
	if err := json.Unmarshal(blob, &cfg); err != nil {
		return "", cher.New("invalid_agentcore_config", cher.M{
			"path":  path,
			"error": err.Error(),
		})
	}

	if cfg.Repo == "" {
		return "", cher.New("missing_agentcore_repo", cher.M{"path": path})
	}

	return cfg.Repo, nil
}

package app

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/tools/cdep"
)

var commitRegRem = regexp.MustCompile(`(?m)^\s*"commit"\s*:\s*"[a-f\d]{40}"\s*,?\n`)
var branchRegRem = regexp.MustCompile(`(?m)^\s*"branch"\s*:\s*"[a-zA-Z\d_-]+"\s*,?\n`)

func (a App) RemFromConfig(path string) (bool, error) {
	if strings.Contains(path, "_default.json") {
		return false, nil
	}

	blob, err := os.ReadFile(path)
	if err != nil {
		if v, ok := err.(*os.PathError); ok {
			if v.Op != "open" {
				return false, err
			}

			return false, nil
		}

		return false, err
	}

	if commitFreeze.Match(blob) && !commitRegAdd.Match(blob) {
		return false, cher.New("frozen_without_commit", cher.M{
			"path": path,
		})
	}

	if commitFreeze.Match(blob) {
		return false, cher.New("frozen", cher.M{"path": path})
	}

	original := make([]byte, len(blob))
	copy(original, blob)

	// if master or undefined, remove reference to both commit hash
	// and branch name from the config in favour of it coming from _default.json
	if getBranchDefinition(blob) != "custom" {
		blob = attemptRemove(blob, commitRegRem)
		blob = attemptRemove(blob, branchRegRem)
	}

	return !bytes.Equal(original, blob), os.WriteFile(path, blob, os.ModePerm)
}

func getBranchDefinition(blob []byte) string {
	// if no "branch" is defined, it's default
	if !branchRegAdd.Match(blob) {
		return "default"
	}

	// if branch is set to the default specifically
	findMe := fmt.Sprintf("\"branch\": \"%s\"", cdep.DefaultBranch)
	if branchRegAdd.Match(blob) && bytes.Contains(blob, []byte(findMe)) {
		return cdep.DefaultBranch
	}

	// branch is defined, but its not the default
	return "custom"
}

func attemptRemove(blob []byte, reg *regexp.Regexp) []byte {
	return reg.ReplaceAll(blob, []byte{})
}

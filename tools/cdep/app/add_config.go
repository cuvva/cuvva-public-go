package app

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	log "github.com/sirupsen/logrus"
)

var commitFreeze = regexp.MustCompile(`(?m)"cdep_freeze"\s*:\s*true`)

var commitRegAdd = regexp.MustCompile(`"commit"\s*:\s*"[a-f\d]{40}"`)
var branchRegAdd = regexp.MustCompile(`"branch"\s*:\s*"([a-zA-Z\d_-]+)"`)

var commitRegAddYaml = regexp.MustCompile(`commit\s*:\s*"?[a-f\d]{40}"?`)
var imageTagRegAddYaml = regexp.MustCompile(`tag\s*:\s*"?[a-z\d-]+"?`)
var branchRegAddYaml = regexp.MustCompile(`branch\s*:\s*"?([a-zA-Z\d_-]+)"?`)

func (a App) AddToConfig(path, branchName, commitHash string) (bool, error) {
	blob, err := os.ReadFile(path)
	if err != nil {
		if v, ok := err.(*os.PathError); ok {
			if v.Op != "open" {
				return false, fmt.Errorf("read file path error: %w", err)
			}

			return false, nil
		}

		return false, fmt.Errorf("read file: %w", err)
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

	if strings.HasSuffix(path, ".yaml") {
		blob = a.doYamlUpdates(path, branchName, commitHash, blob)
	} else {
		blob = a.doJsonUpdates(path, branchName, commitHash, blob)
	}

	return !bytes.Equal(original, blob), os.WriteFile(path, blob, os.ModePerm)
}

func (a App) doJsonUpdates(path string, branchName string, commitHash string, blob []byte) []byte {
	attemptWarnIfOverride(blob, branchRegAdd, branchName, path)

	blob = attemptUpdate(blob, commitRegAdd, "commit", commitHash)
	blob = attemptUpdate(blob, branchRegAdd, "branch", branchName)

	blob = attemptInsert(blob, "commit", commitHash)
	blob = attemptInsert(blob, "branch", branchName)
	return blob
}

func (a App) doYamlUpdates(path string, branchName string, commitHash string, blob []byte) []byte {
	attemptWarnIfOverride(blob, branchRegAddYaml, branchName, path)

	imagePrefix := "master"
	if branchName != "master" {
		imagePrefix = "branch"
	}

	blob = attemptUpdateYaml(blob, commitRegAddYaml, "commit", commitHash)
	blob = attemptUpdateYaml(blob, branchRegAddYaml, "branch", branchName)
	blob = attemptUpdateYaml(blob, imageTagRegAddYaml, "tag", fmt.Sprintf("%s-%s", imagePrefix, commitHash))

	blob = attemptInsertYaml(blob, "commit", commitHash)
	blob = attemptInsertYaml(blob, "branch", branchName)
	return blob
}

// attemptInsert attempts to insert a key into the struct if it doesn't exist
func attemptInsert(blob []byte, key string, value interface{}) []byte {
	strBlob := string(blob)

	// if it does not exist, add it
	if pos := strings.Index(strBlob, key); pos == -1 {
		idx := strings.Index(strBlob, "{")
		strBlob = strBlob[:idx+1] + fmt.Sprintf("\n\t\"%s\": \"%s\",", key, value) + strBlob[idx+1:]
	}

	blob = []byte(strBlob)

	return blob
}

// attemptInsertYaml attempts to insert a key into the struct if it doesn't exist
func attemptInsertYaml(blob []byte, key string, value interface{}) []byte {
	strBlob := string(blob)

	// if it does not exist, add it
	if pos := strings.Index(strBlob, key); pos == -1 {
		strBlob = fmt.Sprintf("%s: %s\n", key, value) + strBlob
	}

	blob = []byte(strBlob)

	return blob
}

// attemptUpdate attempts to change a value of a key if it already exists
func attemptUpdate(blob []byte, reg *regexp.Regexp, key, value string) []byte {
	if reg.Match(blob) {
		replacement := fmt.Sprintf("\"%s\": \"%s\"", key, value)
		blob = reg.ReplaceAll(blob, []byte(replacement))
	}

	return blob
}

// attemptUpdateYaml attempts to change a value of a key if it already exists
func attemptUpdateYaml(blob []byte, reg *regexp.Regexp, key, value string) []byte {
	if reg.Match(blob) {
		replacement := fmt.Sprintf("%s: %s", key, value)
		blob = reg.ReplaceAll(blob, []byte(replacement))
	}

	return blob
}

func attemptWarnIfOverride(blob []byte, reg *regexp.Regexp, newBranch, path string) {
	matches := reg.FindSubmatch(blob)

	if len(matches) <= 1 {
		return
	}

	branchName := string(matches[1])

	if branchName != "master" && branchName != newBranch {
		log.Warnf("Custom branch is going to be overriden, old: %s, new: %s, path: %s", branchName, newBranch, path)
	}
}

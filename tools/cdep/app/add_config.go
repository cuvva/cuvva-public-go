package app

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	log "github.com/sirupsen/logrus"
)

var commitFreeze = regexp.MustCompile(`(?m)"cdep_freeze"\s*:\s*true`)

var commitRegAdd = regexp.MustCompile(`"commit"\s*:\s*"[a-f\d]{40}"`)
var branchRegAdd = regexp.MustCompile(`"branch"\s*:\s*"([a-zA-Z\d_-]+)"`)

func (a App) AddToConfig(path, branchName, commitHash string) (bool, error) {
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

	attemptWarnIfOverride(blob, branchRegAdd, branchName)

	blob = attemptUpdate(blob, commitRegAdd, "commit", commitHash)
	blob = attemptUpdate(blob, branchRegAdd, "branch", branchName)

	blob = attemptInsert(blob, "commit", commitHash)
	blob = attemptInsert(blob, "branch", branchName)

	return !bytes.Equal(original, blob), ioutil.WriteFile(path, blob, os.ModePerm)
}

// AttemptInsert attempts to insert a key into the struct if it doesn't exist
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

// AttemptUpdate attempts to change a value of a key if it already exists
func attemptUpdate(blob []byte, reg *regexp.Regexp, key, value string) []byte {
	if reg.Match(blob) {
		replacement := fmt.Sprintf("\"%s\": \"%s\"", key, value)
		blob = reg.ReplaceAll(blob, []byte(replacement))
	}

	return blob
}

func attemptWarnIfOverride(blob []byte, reg *regexp.Regexp, newBranch string) {
	matches := reg.FindSubmatch(blob)

	if len(matches) <= 1 {
		return
	}

	branchName := string(matches[1])

	if branchName != "master" && branchName != newBranch {
		log.Warnf("Custom branch is going to be overriden, old: %s, new: %s", branchName, newBranch)
	}
}

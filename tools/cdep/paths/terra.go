package paths

import (
	"path"
)

func GetPathForTerra(repo, system, env, workspace string) string {
	return path.Join(repo, system, env, "terra", workspace+".json")
}

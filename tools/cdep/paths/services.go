package paths

import (
	"path"
)

func GetPathForService(repo, system, env, service string) string {
	return path.Join(repo, system, env, "service", service+".json")
}

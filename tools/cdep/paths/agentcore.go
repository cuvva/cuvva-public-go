package paths

import (
	"path"
)

func GetPathForAgentcore(repo, system, env, item string) string {
	return path.Join(repo, system, env, "agentcore", item+".json")
}

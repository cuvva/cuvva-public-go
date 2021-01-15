package paths

import (
	"path"
)

func GetPathForLambda(repo, system, env, lambda string) string {
	return path.Join(repo, system, env, "lambda", lambda+".json")
}

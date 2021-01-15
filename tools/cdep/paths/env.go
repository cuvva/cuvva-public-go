package paths

import (
	"os"

	"github.com/cuvva/cuvva-public-go/lib/cher"
)

const (
	ConfigRepoEnv = "CUVVA_CONFIG_REPO"
	CodeRepoEnv   = "CUVVA_CODE_REPO"
)

func GetConfigRepo() (string, error) {
	return getEnvVar(ConfigRepoEnv)
}

func GetCodeRepo() (string, error) {
	return getEnvVar(CodeRepoEnv)
}

func getEnvVar(envVar string) (string, error) {
	if v, ok := os.LookupEnv(envVar); ok {
		return v, nil
	}

	return "", cher.New("missing_env_var", cher.M{
		"env_var": envVar,
	})
}

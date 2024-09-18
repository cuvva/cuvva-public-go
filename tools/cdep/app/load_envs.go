package app

import (
	"os"
	"path"

	log "github.com/sirupsen/logrus"
)

func (a App) LoadEnvs(repoPath, system, env string) (map[string]struct{}, error) {
	envs := map[string]struct{}{}

	if env != "all" {
		log.Infof("loading single env (%s)", env)
		envs[env] = struct{}{}

		return envs, nil
	}

	log.Info("loading ALL environments")

	files, err := os.ReadDir(path.Join(repoPath, system))
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.IsDir() && f.Name() != "_system" {
			envs[f.Name()] = struct{}{}
		}
	}

	return envs, nil
}

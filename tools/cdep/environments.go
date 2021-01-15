package cdep

import (
	"sort"

	"github.com/cuvva/cuvva-public-go/lib/cher"
)

// ListEnvironments returns a set of allowed environments
func ListEnvironments(system string) []string {
	allowed := map[string]struct{}{}

	switch system {
	case "*", "all", "":
		for _, envs := range Systems {
			for env := range envs {
				allowed[env] = struct{}{}
			}
		}
	default:
		for env := range Systems[system] {
			allowed[env] = struct{}{}
		}
	}

	deduped := []string{}

	for key := range allowed {
		deduped = append(deduped, key)
	}

	sort.Strings(deduped)

	return deduped
}

// ValidateSystemEnvironment against the allowed options
func ValidateSystemEnvironment(sys, env string) error {
	if _, ok := Systems[sys]; !ok {
		return cher.New("unknown_system", cher.M{
			"sys":     sys,
			"env":     env,
			"allowed": ListEnvironments("*"),
		})
	}

	envs := Systems[sys]

	if _, ok := envs[env]; !ok && env != "_system" {
		return cher.New("unknown_environment", cher.M{
			"sys":     sys,
			"env":     env,
			"allowed": ListEnvironments(sys),
		})
	}

	return nil
}

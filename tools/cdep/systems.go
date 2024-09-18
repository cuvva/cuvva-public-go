package cdep

import (
	"sort"

	"github.com/cuvva/cuvva-public-go/lib/cher"
)

// TODO(gm): Make this more dynamic
var Systems = map[string]map[string]struct{}{
	"nonprod": {
		"avocado":   {},
		"basil":     {},
		"coconut":   {},
		"pretzel":   {},
		"ephemeral": {},
		"test":      {},
		"all":       {},
	},
	"prod": {
		"prod": {},
		"all":  {},
	},
}

// ListSystems returns a set of allowed systems
func ListSystems() []string {
	allowed := []string{}

	for sys := range Systems {
		allowed = append(allowed, sys)
	}

	sort.Strings(allowed)

	return allowed
}

func ValidateSystem(sys string) error {
	if _, ok := Systems[sys]; ok {
		return nil
	}

	return cher.New("unknown_system", cher.M{
		"sys":     sys,
		"allowed": ListSystems(),
	})
}

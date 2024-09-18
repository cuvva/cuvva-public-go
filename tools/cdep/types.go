package cdep

import (
	"sort"

	"github.com/cuvva/cuvva-public-go/lib/cher"
)

var Types = map[string]struct{}{
	"lambda":   {},
	"lambdas":  {},
	"service":  {},
	"services": {},
	"terra":    {},
}

func ParseTypeArg(in string) (string, error) {
	if _, ok := Types[in]; !ok {
		return "", cher.New("unknown_type", cher.M{
			"type":    in,
			"allowed": ListTypes(),
		})
	}

	switch in {
	case "service", "services":
		return "service", nil
	case "lambda", "lambdas":
		return "lambda", nil
	case "terra":
		return "terra", nil
	default:
		return "", cher.New("impossible", nil)
	}
}

// ListTypes returns a set of allowed types
func ListTypes() []string {
	allowed := []string{}

	for t := range Types {
		allowed = append(allowed, t)
	}

	sort.Strings(allowed)

	return allowed
}

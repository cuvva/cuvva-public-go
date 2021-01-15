package cdep

import (
	"sort"

	"github.com/cuvva/cuvva-public-go/lib/cher"
)

var WebTypes = map[string]struct{}{
	"cloudfront": {},
	"cf":         {},
}

func ParseWebTypeArg(in string) (string, error) {
	if _, ok := WebTypes[in]; !ok {
		return "", cher.New("unknown_type", cher.M{
			"type":    in,
			"allowed": ListWebTypes(),
		})
	}

	switch in {
	case "cloudfront", "cf":
		return "cloudfront", nil
	default:
		return "", cher.New("impossible", nil)
	}
}

// ListWebTypes returns a set of allowed web types
func ListWebTypes() []string {
	allowed := []string{}

	for t := range WebTypes {
		allowed = append(allowed, t)
	}

	sort.Strings(allowed)

	return allowed
}

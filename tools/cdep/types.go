package cdep

import (
	"regexp"
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

// ValidateCommitHash validates that a commit hash is exactly 40 hexadecimal characters
// This prevents short hashes and branch names from being used, which cause issues
// with ECR image tags and regex matching in config files.
func ValidateCommitHash(commit string) error {
	// Git commit hashes are exactly 40 hexadecimal characters
	commitHashRegex := regexp.MustCompile(`^[a-f0-9]{40}$`)
	if !commitHashRegex.MatchString(commit) {
		return cher.New("invalid_commit_hash", cher.M{
			"commit": commit,
		})
	}
	return nil
}

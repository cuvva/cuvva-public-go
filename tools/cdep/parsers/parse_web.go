package parsers

import (
	"errors"

	"github.com/cuvva/cuvva-public-go/tools/cdep"
)

func ParseWeb(args []string, branch string, prodSys bool) (*Params, error) {
	if len(args) < 2 {
		return nil, errors.New("missing arguments")
	}

	// default to nonprod
	system := "nonprod"

	if prodSys {
		system = "prod"
	}

	t, err := cdep.ParseWebTypeArg(args[0])
	if err != nil {
		return nil, err
	}

	err = cdep.ValidateSystem(system)
	if err != nil {
		return nil, err
	}

	err = cdep.ValidateSystemEnvironment(system, args[1])
	if err != nil {
		return nil, err
	}

	itemSet := []string{}

	if len(args) > 2 {
		itemSet = append(itemSet, args[2:]...)
	}

	return &Params{
		Type:        t,
		System:      system,
		Environment: args[1],
		Branch:      branch,
		Items:       itemSet,
	}, nil
}

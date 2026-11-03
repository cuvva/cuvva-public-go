package parsers

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/tools/cdep"
)

// AllItems can be given in place of a list of services to select every
// service in the environment.
const AllItems = "all"

// exceptions of services that start with "service-"
var exceptions = map[string]struct{}{
	"web-underwriter": {},
	"web-mid":         {},
}

func Parse(args []string, prodSys bool) (*Params, error) {
	if len(args) < 2 {
		return nil, errors.New("missing arguments")
	}

	// default to nonprod
	system := "nonprod"

	if prodSys {
		system = "prod"
	}

	t, err := cdep.ParseTypeArg(args[0])
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
	allItems := false

	if slices.Contains(args[2:], AllItems) {
		if len(args) != 3 {
			return nil, cher.New("all_must_be_alone", nil)
		}

		if t != "service" {
			return nil, cher.New("all_requires_services", cher.M{"type": t})
		}

		allItems = true
	} else if len(args) > 2 {
		for _, item := range args[2:] {
			// if we're dealing with services
			if t == "service" {
				// if this service is not in the exception list
				if _, ok := exceptions[item]; !ok {
					// check for the prefix
					if !strings.HasPrefix(item, "service-") {
						// add it if it does not exist
						item = fmt.Sprintf("%s%s", "service-", item)
					}
				}
			}

			itemSet = append(itemSet, item)
		}
	}

	return &Params{
		Type:        t,
		System:      system,
		Environment: args[1],
		Items:       itemSet,
		AllItems:    allItems,
	}, nil
}

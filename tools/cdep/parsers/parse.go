package parsers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cuvva/cuvva-public-go/tools/cdep"
)

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

	if len(args) > 2 {
		// Handle "all" keyword for services
		if len(args) == 3 && args[2] == "all" {
			// Leave itemSet empty - this signals to update all services
			// The Update method will handle discovering all services
		} else {
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
	}

	return &Params{
		Type:        t,
		System:      system,
		Environment: args[1],
		Items:       itemSet,
	}, nil
}

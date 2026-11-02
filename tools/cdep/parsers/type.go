package parsers

import (
	"fmt"
	"strings"
)

type Params struct {
	Type        string
	Environment string
	Branch      string
	Commit      string
	System      string
	Items       []string
	// AllItems is set when "all" was given instead of a list of services, and
	// selects every service in the environment.
	AllItems bool
	// Only limits services to those built by one pipeline: cdep.LanguageGo or
	// cdep.LanguageJS. Empty means every service.
	Only    string
	Message string
}

func (p Params) String(command string) (out string) {
	out = fmt.Sprintf("%s %s %s", command, p.Type, p.Environment)

	if p.AllItems {
		out = fmt.Sprintf("%s all", out)
	} else if len(p.Items) > 0 {
		items := []string{}
		for _, item := range p.Items {
			items = append(items, strings.TrimPrefix(item, "service-"))
		}

		itemsStr := strings.Join(items, " ")
		out = fmt.Sprintf("%s %s", out, itemsStr)
	}

	if p.Only != "" {
		out = fmt.Sprintf("%s --%s-only", out, p.Only)
	}

	if p.Branch != "master" {
		out = fmt.Sprintf("%s -b %s", out, p.Branch)
	}

	if p.Commit != "" {
		out = fmt.Sprintf("%s -c %s", out, p.Commit)
	}

	if p.System == "prod" {
		out = fmt.Sprintf("%s --prod", out)
	}

	return
}

package parsers

import (
	"fmt"
	"strings"
)

type Params struct {
	Type        string
	Environment string
	Branch      string
	System      string
	Items       []string
	Message     string
}

func (p Params) String(command string) (out string) {
	out = fmt.Sprintf("%s %s %s", command, p.Type, p.Environment)

	if len(p.Items) > 0 {
		items := []string{}
		for _, item := range p.Items {
			items = append(items, strings.TrimPrefix(item, "service-"))
		}

		itemsStr := strings.Join(items, " ")
		out = fmt.Sprintf("%s %s", out, itemsStr)
	}

	if p.Branch != "master" {
		out = fmt.Sprintf("%s -b %s", out, p.Branch)
	}

	if p.System == "prod" {
		out = fmt.Sprintf("%s --prod", out)
	}

	return
}

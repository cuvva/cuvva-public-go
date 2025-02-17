package app

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/cuvva/cuvva-public-go/tools/cdep/parsers"
)

type Dashboard struct {
	Name string
	URL  string
}

func chooseDashboards(req *parsers.Params, envNames []string) []Dashboard {
	var dashboards []Dashboard
	switch req.Type {
	case "service":
		dashboards = append(dashboards, Dashboard{
			Name: "Deployment Status",
			URL:  deploymentStatusDashboard(envNames, req.Items),
		})
	}

	switch req.System {
	case "prod":
		dashboards = append(dashboards, Dashboard{
			Name: "Prod Errors (inc lambda)",
			URL:  "https://app.datadoghq.eu/logs?saved-view-id=48532",
		})
	default:
		dashboards = append(dashboards, Dashboard{
			Name: "Nonprod Errors",
			URL:  "https://app.datadoghq.eu/logs?saved-view-id=42541",
		})
	}

	return dashboards
}

func printDashboards(dashboards []Dashboard) {
	if len(dashboards) == 0 {
		return
	}

	for _, dashboard := range dashboards {
		fmt.Println("")
		fmt.Println(dashboard.Name)
		fmt.Println(dashboard.URL)
	}

	fmt.Println("")
}

func deploymentStatusDashboard(envs []string, services []string) string {
	for i, svc := range services {
		services[i] = strings.TrimPrefix(svc, "service-")
	}

	endpoint := "https://app.datadoghq.eu/logs"
	query := url.Values{
		"query": []string{
			fmt.Sprintf(
				"source:service env:(%s) service:(%s)",
				strings.Join(envs, " OR "),
				strings.Join(services, " OR "),
			),
		},
		"agg_m":            []string{"count"},
		"agg_m_source":     []string{"base"},
		"agg_q":            []string{"@_commit_hash,status"},
		"agg_q_source":     []string{"base,base"},
		"agg_t":            []string{"count"},
		"analyticsOptions": []string{`["bars","dog_classic",null,null,"value"]`},
		"cols":             []string{"service,@level,@rpc_method,@message,@error,@http_duration"},
		"fromUser":         []string{"true"},
		"messageDisplay":   []string{"inline"},
		"refresh_mode":     []string{"sliding"},
		"saved-view-id":    []string{"152718"},
		"sort_m":           []string{"count,count"},
		"sort_m_source":    []string{"base,base"},
		"sort_t":           []string{"count,count"},
		"storage":          []string{"hot"},
		"stream_sort":      []string{"desc"},
		"top_n":            []string{"5,5"},
		"top_o":            []string{"top,top"},
		"viz":              []string{"timeseries"},
		"x_missing":        []string{"true,true"},
		"from_ts":          []string{"1739785059959"},
		"to_ts":            []string{"1739785959959"},
		"live":             []string{"true"},
	}

	return fmt.Sprintf(
		"%s?%s",
		endpoint,
		query.Encode(),
	)
}

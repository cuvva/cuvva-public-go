package dotnot

import (
	"fmt"
)

func To(in map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{})

	for k, v := range in {
		switch t := v.(type) {
		case map[string]interface{}:
			n := To(t)
			for nk, nv := range n {
				out[fmt.Sprintf("%s.%s", k, nk)] = nv
			}
		default:
			out[k] = v
		}
	}

	return out
}

func From(in map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

package dotnot

import (
	"fmt"
)

func To(in map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{})

	for k, v := range in {
		switch t := v.(type) {
		case map[string]interface{}:
			for nk, nv := range To(t) {
				out[fmt.Sprintf("%s.%s", k, nk)] = nv
			}
		case []interface{}:
			for sk, sv := range t {
				switch u := sv.(type) {
				case map[string]interface{}:
					for nsk, nsv := range To(u) {
						out[fmt.Sprintf("%s.%d.%s", k, sk, nsk)] = nsv
					}
				default:
					out[k] = v
				}
			}
		default:
			out[k] = v
		}
	}

	return out
}

func From(in map[string]interface{}) map[string]interface{} {
	// TODO(gm): build this

	out := make(map[string]interface{})

	return out
}

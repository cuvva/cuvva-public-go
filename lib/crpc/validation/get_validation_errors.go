package validation

import (
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

func GetValidationErrorsAsString(result *gojsonschema.Result) string {
	detail := make([]string, 0)
	for _, vErr := range result.Errors() {
		detail = append(detail, vErr.String())
	}
	return strings.Join(detail, ", ")
}

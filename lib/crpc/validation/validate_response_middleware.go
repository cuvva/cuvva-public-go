package validation

import (
	"fmt"
	"github.com/cuvva/cuvva-public-go/lib/crpc"
	"github.com/xeipuuv/gojsonschema"
	"net/http"
)

// Validate response the JSON body and applies a JSON Schema validation.
func ValidateResponseMiddleware(schema *gojsonschema.Schema) crpc.MiddlewareFunc {
	return func(next crpc.HandlerFunc) crpc.HandlerFunc {
		return func(res http.ResponseWriter, req *crpc.Request) error {
			copier := newCopyResponseWriter(res)

			err := next(copier, req)
			if err != nil {
				return err
			}

			// don't do any validation on 204 - No Content
			if copier.responseCopy.statusCode != http.StatusOK {
				return nil
			}

			result, err := schema.Validate(gojsonschema.NewBytesLoader(copier.responseCopy.body.Bytes()))
			if err != nil {
				fmt.Errorf("crpc failed to validate response: %w", err)
			}

			if !result.Valid() {
				fmt.Errorf("crpc response does not fulfil schema: %s", getValidationErrorsAsString(result))
			}

			return nil
		}
	}
}

func getValidationErrorsAsString(result *gojsonschema.Result) []string {
	detail := []string{}
	for _, vErr := range result.Errors() {
		detail = append(detail, fmt.Sprintf("%s", vErr))
	}
	return detail
}

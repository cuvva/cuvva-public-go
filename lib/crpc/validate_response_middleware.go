package crpc

import (
	"fmt"
	"net/http"

	"github.com/cuvva/cuvva-public-go/lib/clog"
	"github.com/cuvva/cuvva-public-go/lib/crpc/validation"
	"github.com/xeipuuv/gojsonschema"
)

// Validate response the JSON body and applies a JSON Schema validation.
func ValidateResponseMiddleware(schema *gojsonschema.Schema) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(res http.ResponseWriter, req *Request) error {
			copier := validation.NewCopyResponseWriter(res)

			err := next(copier, req)
			if err != nil {
				return err
			}

			// don't do any validation on 204 - No Content
			if copier.StatusCode != http.StatusOK {
				return nil
			}

			result, err := schema.Validate(gojsonschema.NewBytesLoader(copier.Body.Bytes()))
			if err != nil {
				clog.Get(req.Context()).WithError(err).Warn("crpc failed to validate response")
			}

			if !result.Valid() {
				vErr := fmt.Errorf(validation.GetValidationErrorsAsString(result))
				clog.Get(req.Context()).WithError(vErr).Warn("crpc response does not fulfil schema")
			}

			return nil
		}
	}
}

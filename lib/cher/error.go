package cher

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// errors that are expected to be common across services
const (
	BadRequest        = "bad_request"
	Unauthorized      = "unauthorized"
	AccessDenied      = "access_denied"
	NotFound          = "not_found"
	RouteNotFound     = "route_not_found"
	MethodNotAllowed  = "method_not_allowed"
	Unknown           = "unknown"
	NoLongerSupported = "no_longer_supported"
	TooManyRequests   = "too_many_requests"
	ContextCanceled   = "context_canceled"
	EOF               = "eof"
	UnexpectedEOF     = "unexpected_eof"
	RequestTimeout    = "request_timeout"

	CoercionError = "unable_to_coerce_error"
)

// E implements the official Cuvva Error structure
type E struct {
	Code    string `json:"code"`
	Meta    M      `json:"meta,omitempty"`
	Reasons []E    `json:"reasons,omitempty"`
}

// New returns a new E structure with code, meta, and optional reasons.
func New(code string, meta M, reasons ...E) E {
	return E{
		Code:    code,
		Meta:    meta,
		Reasons: reasons,
	}
}

// Errorf returns a new E structure, with a message formatted by fmt.
func Errorf(code string, meta M, format string, args ...interface{}) E {
	meta["message"] = fmt.Sprintf(format, args...)

	return E{
		Code: code,
		Meta: meta,
	}
}

// StatusCode returns the HTTP Status Code associated with the
// current error code.
// Defaults to 400 Bad Request because if something's explicitly
// handled with Cher, it is considered "by design" and not
// worthy of a 500, which will alert.
func (e E) StatusCode() int {
	switch e.Code {
	case BadRequest:
		return http.StatusBadRequest

	case Unauthorized:
		return http.StatusUnauthorized

	case AccessDenied:
		return http.StatusForbidden

	case NotFound, RouteNotFound:
		return http.StatusNotFound

	case MethodNotAllowed:
		return http.StatusMethodNotAllowed

	case NoLongerSupported:
		return http.StatusGone

	case TooManyRequests:
		return http.StatusTooManyRequests

	case Unknown, CoercionError, RequestTimeout:
		return http.StatusInternalServerError
	}

	return http.StatusBadRequest
}

// Error implements the error interface.
func (e E) Error() string {
	return e.Code
}

// Serialize returns a json representation of the Cuvva Error structure
func (e E) Serialize() string {
	output, err := json.Marshal(e)
	if err != nil {
		return ""
	}

	return string(output)
}

// M it an alias type for map[string]interface{}
type M map[string]interface{}

// Coerce attempts to coerce a Cuvva Error out of any object.
// - `E` types are just returned as-is
// - strings are taken as the Code for an E object
// - bytes are unmarshaled from JSON to an E object
// - types implementing the `error` interface to an E object with the error as a reason
func Coerce(v interface{}) E {
	switch v := v.(type) {
	case E:
		return v

	case string:
		return E{Code: v}

	case []byte:
		var e E

		err := json.Unmarshal(v, &e)
		if err != nil {
			return E{
				Code: CoercionError,
				Meta: M{
					"message": err.Error(),
				},
			}
		}

		return e

	case error:
		v = errors.Cause(v)

		return E{
			Code: Unknown,
			Meta: M{
				"message": v.Error(),
			},
		}
	}

	return E{Code: CoercionError}
}

func (e E) Value() (driver.Value, error) {
	return json.Marshal(e)
}

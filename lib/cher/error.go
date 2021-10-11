package cher

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/kr/pretty"
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
	pretty.Log("in error2")
	pretty.Log(e)
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

// Unwrap returns an error from Error (or nil if there are no errors).
// This error returned will further support Unwrap to get the next error,
// etc. The order will match the order of Errors in the multierror.Error
// at the time of calling.
//
// The resulting error supports errors.As/Is/Unwrap so you can continue
// to use the stdlib errors package to introspect further.
//
// This will perform a shallow copy of the errors slice. Any errors appended
// to this error after calling Unwrap will not be available until a new
// Unwrap is called on the multierror.Error.
func (e *E) Unwrap() error {
	pretty.Log("in unwrap")
	pretty.Log(e)
	// If we have no errors then we do nothing
	if e == nil || len(e.Reasons) == 0 {
		return nil
	}

	// If we have exactly one error, we can just return that directly.
	if len(e.Reasons) == 1 {
		return e.Reasons[0]
	}

	// Shallow copy the slice
	errs := make([]E, len(e.Reasons))
	copy(errs, e.Reasons)
	return chain(errs)
}

// chain implements the interfaces necessary for errors.Is/As/Unwrap to
// work in a deterministic way with multierror. A chain tracks a list of
// errors while accounting for the current represented error. This lets
// Is/As be meaningful.
//
// Unwrap returns the next error. In the cleanest form, Unwrap would return
// the wrapped error here but we can't do that if we want to properly
// get access to all the errors. Instead, users are recommended to use
// Is/As to get the correct error type out.
//
// Precondition: []E is non-empty (len > 0)
type chain []E

// Error implements the error interface
func (e chain) Error() string {
	pretty.Log("in error")
	pretty.Log(e)
	return e[0].Error()
}

// Unwrap implements errors.Unwrap by returning the next error in the
// chain or nil if there are no more errors.
func (e chain) Unwrap() error {
	pretty.Log("in unwrap2")
	pretty.Log(e)
	if len(e) == 1 {
		return nil
	}

	return e[1:]
}

// As implements errors.As by attempting to map to the current value.
func (e chain) As(target interface{}) bool {
	pretty.Log("in as")
	pretty.Log(e)
	return errors.As(e[0], target)
}

// Is implements errors.Is by comparing the current value directly.
func (e chain) Is(target error) bool {
	pretty.Log("in is")
	pretty.Log(e)
	return errors.Is(e[0], target)
}
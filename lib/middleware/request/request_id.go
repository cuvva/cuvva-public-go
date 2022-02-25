package request

import (
	"context"
	"net/http"

	"github.com/cuvva/cuvva-public-go/lib/ksuid"
)

// RequestIDKey is the key used for acquired/generated Request ID in a requests context map
const RequestIDKey ContextKey = "Request-ID"

// GetRequestID returns the request id embedded within a HTTP request's context,
// or an empty string if no request id has been established.
func GetRequestID(r *http.Request) string {
	return GetRequestIDContext(r.Context())
}

// GetRequestIDContext returns the request id embedded within a context,
// or an empty string if no request id has been established.
func GetRequestIDContext(ctx context.Context) string {
	if str, ok := ctx.Value(RequestIDKey).(string); ok {
		return str
	}

	return ""
}

// SetRequestID sets the request id within a HTTP requests's context,
// it will overwrite any existing request id.
func SetRequestID(r *http.Request, requestID string) *http.Request {
	return r.WithContext(SetRequestIDContext(r.Context(), requestID))
}

// SetRequestIDContext sets the request id within a context,
// it will overwrite any existing request id.
func SetRequestIDContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// GetOrSetRequestID will return any request ID found in the context, and
// if one does not exist, set and return.
func GetOrSetRequestID(ctx context.Context) (context.Context, string) {
	requestID := GetRequestIDContext(ctx)
	if requestID != "" {
		return ctx, requestID
	}

	requestID = ksuid.Generate(ctx, "req").String()

	ctx = SetRequestIDContext(ctx, requestID)
	return ctx, requestID
}

// RequestID either generates a new request id or acquires an existing request id
// from the http requests headers, and embeds into the requests context and
// response headers.
func RequestID(next http.Handler) http.Handler {
	const requestIDHeader = "Request-Id"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(requestIDHeader)
		if requestID == "" {
			// if no Request-ID is passed, generate/originate a new one
			requestID = ksuid.Generate(r.Context(), "req").String()
		}

		w.Header().Set(requestIDHeader, requestID)

		r = SetRequestID(r, requestID)

		next.ServeHTTP(w, r)
	})
}

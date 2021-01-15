package jsonclient

import (
	"context"

	"github.com/cuvva/cuvva-public-go/lib/ksuid"
)

// ContextKey is our custom definition to prevent collisions with other
// third-party libraries
type ContextKey string

// RequestIDKey is the key used for acquired/generated Request ID in a requests context map
const RequestIDKey ContextKey = "Request-ID"

// GetRequestIDContext returns the request id embedded within a context,
// or an empty string if no request id has been established.
func GetRequestIDContext(ctx context.Context) string {
	if str, ok := ctx.Value(RequestIDKey).(string); ok {
		return str
	}

	return ""
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

	requestID = ksuid.Generate("req").String()

	ctx = SetRequestIDContext(ctx, requestID)
	return ctx, requestID
}

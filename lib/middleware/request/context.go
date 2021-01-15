package request

import (
	"context"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/clog"
)

// maps are type aware, define a custom string type for context keys
// to prevent collisions with third-party context that uses the same key.
type ContextKey string

// ForkContext provides a callback function with a new context inheriting values from the request context, and will log any error returned by the callback
func ForkContext(ctx context.Context, fn func(context.Context) error) {
	requestID := GetRequestIDContext(ctx)
	logger := clog.Get(ctx)

	newCtx := context.Background()

	newCtx = clog.Set(newCtx, logger)
	newCtx = SetRequestIDContext(newCtx, requestID)

	go func() {
		err := fn(newCtx)
		if err != nil {
			clog.Get(newCtx).WithError(err).Log(clog.DetermineLevel(err), "forked context errored")
		}
	}()
}

// ForkContext provides a callback function with a new context inheriting values from the request context with a timeout, , and will log any error returned by the callback
func ForkContextWithTimeout(ctx context.Context, timeout time.Duration, fn func(context.Context) error) {
	requestID := GetRequestIDContext(ctx)
	logger := clog.Get(ctx)

	newCtx := context.Background()

	newCtx = clog.Set(newCtx, logger)
	newCtx = SetRequestIDContext(newCtx, requestID)

	newCtx, cancel := context.WithTimeout(newCtx, timeout)

	go func() {
		defer cancel()
		err := fn(newCtx)
		if err != nil {
			clog.Get(newCtx).WithError(err).Log(clog.DetermineLevel(err), "forked context errored")
		}
	}()
}

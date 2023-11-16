package request

import (
	"context"
	"fmt"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/clog"
)

// forkedContext allows cloning context values
// while discarding deadlines and cancellations
type forkedContext struct {
	ctx context.Context
}

func (forkedContext) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

func (forkedContext) Done() <-chan struct{} {
	return nil
}

func (forkedContext) Err() error {
	return nil
}

func (d forkedContext) Value(key interface{}) interface{} {
	return d.ctx.Value(key)
}

func cloneContext(ctx context.Context) context.Context {
	return forkedContext{ctx}
}

// ContextKey maps are type aware, define a custom string type for context keys
// to prevent collisions with third-party context that uses the same key.
type ContextKey string

// ForkContext provides a callback function with a new context inheriting values from the request context, and will log any error returned by the callback
func ForkContext(ctx context.Context, fn func(context.Context) error) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("recovered from panic: %v", r)
			if err != nil {
				panic(r)
			}
		}
	}()
	newCtx := cloneContext(ctx)

	go func() {
		err := fn(newCtx)
		if err != nil {
			clog.Get(newCtx).WithError(err).Log(clog.DetermineLevel(err, true), "forked context errored")
		}
	}()
}

// ForkContextWithTimeout provides a callback function with a new context inheriting values from the request context with a timeout, and will log any error returned by the callback
func ForkContextWithTimeout(ctx context.Context, timeout time.Duration, fn func(context.Context) error) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("recovered from panic: %v", r)
			if err != nil {
				panic(r)
			}
		}
	}()

	newCtx, cancel := context.WithTimeout(cloneContext(ctx), timeout)

	go func() {
		defer cancel()
		err := fn(newCtx)
		if err != nil {
			clog.Get(newCtx).WithError(err).Log(clog.DetermineLevel(err, true), "forked context errored")
		}
	}()
}

package servicecontext

import (
	"context"
)

// Info type holds useful info about the currently-running service
type Info struct {
	Name        string
	Environment string
}

var service *Info

// Set gives the singleton pointer a value
func Set(name, env string) {
	service = &Info{
		Name:        name,
		Environment: env,
	}
}

// Get returns the value of the singleton instance pointer
func Get() Info {
	if service == nil {
		panic("cannot get service context before it is set")
	}

	return *service
}

// IsSet returns true if the singleton has been initialised
func IsSet() bool {
	return service != nil
}

type contextKey string

var (
	infoContextKey = contextKey("info")
)

// SetContext wraps the context with the service info
func SetContext(ctx context.Context, name, env string) context.Context {
	return context.WithValue(ctx, infoContextKey, &Info{
		Name:        name,
		Environment: env,
	})
}

// GetContext retrieves the service info from the context
func GetContext(ctx context.Context) *Info {
	if val, ok := ctx.Value(infoContextKey).(*Info); ok {
		return val
	}

	return nil
}

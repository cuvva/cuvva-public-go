package config

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func ContextWithCancelOnSignal(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		defer cancel()
		select {
		case <-stop:
		case <-ctx.Done():
		}
	}()

	return ctx
}

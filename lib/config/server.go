package config

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

// DefaultGraceful is the graceful shutdown timeout applied when no
// configuration value is given.
const DefaultGraceful = 5

// Server configures the binding and security of an HTTP server.
type Server struct {
	Addr string `json:"addr"`

	// Graceful enables graceful shutdown and is the time in seconds to wait
	// for all outstanding requests to terminate before forceably killing the
	// server. When no value is given, DefaultGraceful is used. Graceful
	// shutdown is disabled when less than zero.
	Graceful int `json:"graceful"`
}

// ListenAndServe configures a HTTP server and begins listening for clients.
func (cfg *Server) ListenAndServe(srv *http.Server) error {
	// only set listen address if none is already configured
	if srv.Addr == "" {
		srv.Addr = cfg.Addr
	}

	if cfg.Graceful == 0 {
		cfg.Graceful = DefaultGraceful
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	errs := make(chan error, 1)

	go func() {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			errs <- err
		}
	}()

	select {
	case err := <-errs:
		return err

	case <-stop:
		if cfg.Graceful > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Graceful)*time.Second)
			defer cancel()

			return srv.Shutdown(ctx)
		}

		return nil
	}
}

func (cfg *Server) Listen() (net.Listener, error) {
	l, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, err
	}

	return l, nil
}

// Serve the HTTP requests on the specified listener, and gracefully close when the context is cancelled.
func (cfg *Server) Serve(ctx context.Context, l net.Listener, srv *http.Server) (err error) {
	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		err := srv.Serve(l)
		if err == http.ErrServerClosed {
			return nil
		}

		return err
	})

	eg.Go(func() error {
		select {
		case <-ctx.Done():
			logrus.Println("shutting down gracefully")
			ctx, cancel := context.WithTimeout(context.Background(), cfg.GracefulTimeout())
			defer cancel()
			return srv.Shutdown(ctx)
		case <-egCtx.Done():
			return nil
		}
	})

	return eg.Wait()
}

func (cfg *Server) GracefulTimeout() time.Duration {
	if cfg.Graceful == 0 {
		cfg.Graceful = DefaultGraceful
	}

	return time.Duration(cfg.Graceful) * time.Second
}

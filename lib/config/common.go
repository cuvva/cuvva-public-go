package config

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/cuvva/cuvva-public-go/lib/db/mongodb"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"golang.org/x/sync/errgroup"
)

// Redis configures a connection to a Redis database.
type Redis struct {
	URI          string        `json:"uri"`
	DialTimeout  time.Duration `json:"dial_timeout"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
}

// Options returns a configured redis.Options structure.
func (r Redis) Options() (*redis.Options, error) {
	opts, err := redis.ParseURL(r.URI)
	if err != nil {
		return nil, err
	}

	opts.DialTimeout = r.DialTimeout
	opts.ReadTimeout = r.ReadTimeout
	opts.WriteTimeout = r.WriteTimeout

	return opts, nil
}

// Connect returns a connected redis.Client instance.
func (r Redis) Connect() (*redis.Client, error) {
	opts, err := r.Options()
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)

	if err := client.Ping().Err(); err != nil {
		return client, err
	}

	return client, nil
}

// MongoDB configures a connection to a Mongo database.
type MongoDB struct {
	URI string `json:"uri"`
}

// Options returns the MongoDB client options and database name.
func (m MongoDB) Options() (opts *options.ClientOptions, dbName string, err error) {
	opts = options.Client().ApplyURI(m.URI)
	err = opts.Validate()
	if err != nil {
		return
	}

	// all Go services use majority writes, and this is unlikely to change
	// if it does change, switch to accepting as an argument
	opts.WriteConcern = writeconcern.New(writeconcern.WMajority())

	cs, err := connstring.Parse(m.URI)
	if err != nil {
		return
	}

	dbName = cs.Database
	if dbName == "" {
		err = errors.New("missing mongo database name")
	}

	return
}

// Connect returns a connected mongo.Database instance.
func (m MongoDB) Connect() (*mongodb.Database, error) {
	opts, dbName, err := m.Options()
	if err != nil {
		return nil, err
	}

	// this package can only be used for service config
	// so can only happen at init-time - no need to accept context input
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return mongodb.Connect(ctx, opts, dbName)
}

// JWT configures public (and optionally private) keys and issuer for
// JSON Web Tokens. It is intended to be used in composition rather than a key.
type JWT struct {
	Issuer  string `json:"issuer"`
	Public  string `json:"public"`
	Private string `json:"private,omitempty"`
}

// AWS configures credentials for access to Amazon Web Services.
// It is intended to be used in composition rather than a key.
type AWS struct {
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`

	Region string `json:"region,omitempty"`
}

// Credentials returns a configured set of AWS credentials.
func (a AWS) Credentials() *credentials.Credentials {
	if a.AccessKeyID != "" && a.AccessKeySecret != "" {
		return credentials.NewStaticCredentials(a.AccessKeyID, a.AccessKeySecret, "")
	}

	return nil
}

// Session returns an AWS Session configured with region and credentials.
func (a AWS) Session() (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region:      aws.String(a.Region),
		Credentials: a.Credentials(),
	})
}

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

	logrus.WithField("addr", cfg.Addr).Info("listening")
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
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Graceful)*time.Second)
			defer cancel()
			return srv.Shutdown(ctx)
		case <-egCtx.Done():
			return nil
		}
	})

	return eg.Wait()
}

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

// UnderwriterOpts represents the underwriters info/models options.
type UnderwriterOpts struct {
	IncludeUnreleased bool `json:"include_unreleased"`
}

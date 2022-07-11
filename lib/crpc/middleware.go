package crpc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/blang/semver"
	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/lib/clog"
	"github.com/cuvva/cuvva-public-go/lib/middleware/request"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
)

// Logger inherits the context logger and reports RPC request success/failure.
func Logger() MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(res http.ResponseWriter, req *Request) error {
			ctx := req.Context()

			// Add fields to the request context scoped logger, if one exists
			handleCLogError(clog.SetFields(ctx, clog.Fields{
				"rpc_version": req.Version,
				"rpc_method":  req.Method,
			}))

			t1 := time.Now()
			err := next(res, req)
			t2 := time.Now()

			handleCLogError(clog.SetFields(ctx, clog.Fields{
				"rpc_duration":    t2.Sub(t1).String(),
				"rpc_duration_us": int64(t2.Sub(t1) / time.Microsecond),
			}))

			if err == nil {
				return nil
			}

			// rewrite common errors to internal error standard
			if err == io.EOF {
				err = cher.New(cher.EOF, nil)
			} else if err == io.ErrUnexpectedEOF {
				err = cher.New(cher.UnexpectedEOF, nil)
			} else if strings.Contains(err.Error(), context.Canceled.Error()) {
				err = cher.New(cher.ContextCanceled, nil)
			}

			handleCLogError(clog.SetError(ctx, err))

			return err
		}
	}
}

func handleCLogError(err error) {
	if err != nil {
		logrus.New().WithError(err).Warn("rpc log middleware failed")
	}
}

// Validate buffers the JSON body and applies a JSON Schema validation.
func Validate(ls *gojsonschema.Schema) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(res http.ResponseWriter, req *Request) error {
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				if netErr, ok := err.(net.Error); ok {
					clog.Get(req.Context()).WithError(netErr).Warn("network error reading request body")
					return io.ErrUnexpectedEOF
				}

				return fmt.Errorf("crpc failed to read request body: %w", err)
			}

			ld := gojsonschema.NewBytesLoader(body)

			result, err := ls.Validate(ld)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return cher.New(cher.BadRequest, cher.M{"message": "invalid JSON"})
				}

				return fmt.Errorf("crpc schema validation failed: %w", err)
			}

			err = CoerceJSONSchemaError(result)
			if err != nil {
				return err
			}

			req.Body = ioutil.NopCloser(bytes.NewReader(body))
			return next(res, req)
		}
	}
}

func CoerceJSONSchemaError(result *gojsonschema.Result) error {
	if result.Valid() {
		return nil
	}

	var reasons []cher.E

	errs := result.Errors()
	for _, err := range errs {
		reasons = append(reasons, cher.E{
			Code: "schema_failure",
			Meta: cher.M{
				"field":   err.Field(),
				"type":    err.Type(),
				"message": err.Description(),
			},
		})
	}

	return cher.New(cher.BadRequest, nil, reasons...)
}

// Instrument adds prometheus compatible metrics collection to RPC functions.
// The metrics collected are:
//   - duration
//   - status code
//   - total in-flight
func Instrument(r prometheus.Registerer) func(HandlerFunc) HandlerFunc {
	reqDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "rpc_request_duration_seconds",
			Help:    "Duration of an RPC request in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "version"},
	)
	reqTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rpc_request_total",
			Help: "Total number of RPC requests",
		},
		[]string{"method", "version"},
	)
	resErrorCode := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rpc_response_error_code",
			Help: "Total number of failed RPC requests",
		},
		[]string{"method", "version", "code"},
	)

	r.MustRegister(reqDuration)
	r.MustRegister(reqTotal)
	r.MustRegister(resErrorCode)

	return func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *Request) error {
			start := time.Now()

			err := next(w, r)

			reqDuration.WithLabelValues(r.Method, r.Version).Observe(time.Since(start).Seconds())

			if cuvvaErr, ok := err.(cher.E); ok {
				resErrorCode.WithLabelValues(r.Method, r.Version, cuvvaErr.Code).Inc()
			} else if err != nil {
				resErrorCode.WithLabelValues(r.Method, r.Version, "unknown").Inc()
			}

			reqTotal.WithLabelValues(r.Method, r.Version).Inc()

			return err
		}
	}
}

// VersionRequirement is our predicate type for checking client versions
type VersionRequirement struct {
	Platform string
	Version  semver.Version
}

// NewVersionRequirement creates a predicate for evaluation
// e.g. IF platform is "ios" THEN the minimum version allowed is "3.6.8"
func NewVersionRequirement(platform string, ver semver.Version) VersionRequirement {
	return VersionRequirement{
		Platform: platform,
		Version:  ver,
	}
}

// RequireMinimumClientVersions will return an error to the client if:
// - the reqest has a valid client version header, and
// - the platform matches but the semver requirement is greater than the client semver
func RequireMinimumClientVersions(requirements ...VersionRequirement) func(HandlerFunc) HandlerFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *Request) error {
			client := request.GetClientVersionContext(r.Context())
			if client == nil {
				err := next(w, r)
				return err
			}

			for _, requirement := range requirements {
				if requirement.Platform == client.Platform && requirement.Version.GT(client.Version) {
					return cher.New(cher.NoLongerSupported, nil)
				}
			}

			err := next(w, r)
			return err
		}
	}
}

// Package clog (context logger) provides shared log configuration and helpers for building Cuvva HTTP services
package clog

import (
	"context"
	"errors"
	"os"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/lib/servicecontext"
	"github.com/cuvva/cuvva-public-go/lib/version"
	"github.com/sirupsen/logrus"
)

type contextKey string

type Fields map[string]interface{}

// LoggerKey is the key used for request-scoped loggers in a requests context map
const loggerKey contextKey = "clog"

const (
	// ServiceKey is the log entry key for the name of the crpc service
	ServiceKey = "_service"

	// HostKey is the log entry key for the hostname / container ID
	HostKey = "_hostname"

	// VersionKey is the log entry key for the current version of the codebase
	VersionKey = "_commit_hash"

	// LevelKey is the log entry key for the log level
	LevelKey = "_level"

	// MessageKey is the log entry key for a generic message
	MessageKey = "_message"

	// TimestampKey is the log entry key for the timestamp
	TimestampKey = "_timestamp"
)

// Config allows services to configure the logging format, level and storage options
// for Logrus logging.
type Config struct {
	// Format configures the output format. Possible options:
	//   - text - logrus default text output, good for local development
	//   - json - fields and message encoded as json, good for storage in e.g. cloudwatch
	Format string `json:"format"`

	// Debug enables debug level logging, otherwise INFO level
	Debug bool `json:"debug"`
}

// Configure applies Cuvva standard Logging structure options to a logrus Entry.
func (c Config) Configure(ctx context.Context) (log *logrus.Entry) {
	var serviceName string
	if svc := servicecontext.GetContext(ctx); svc != nil {
		serviceName = svc.Name
	}

	log = logrus.WithFields(logrus.Fields{
		ServiceKey: serviceName,
		VersionKey: version.Revision,
	})

	hostname, err := os.Hostname()
	if err != nil {
		log.WithError(err).Warn("logger hostname configuration failed")
		hostname = "unknown"
	}

	log = log.WithField(HostKey, hostname)

	switch c.Format {
	case "json", "logstash":
		log.Logger.Formatter = &logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyLevel: LevelKey,
				logrus.FieldKeyMsg:   MessageKey,
				logrus.FieldKeyTime:  TimestampKey,
			},
		}

	default:
		log.Logger.Formatter = &logrus.TextFormatter{}
	}

	if c.Debug {
		log.Logger.Level = logrus.DebugLevel
		log.Debug("debug logging enabled")
	} else {
		log.Logger.Level = logrus.InfoLevel
	}

	return
}

// ContextLogger wraps logrus Entry to allow field mutation, which means the
// context itself can store a pointer to a ContextLogger, so it doesn't need
// replacing each time new fields are added to the logger
type ContextLogger struct {
	entry            *logrus.Entry
	timeoutsAsErrors bool
}

// NewContextLogger creates a new (mutable) ContextLogger instance from an (immutable) logrus Entry
func NewContextLogger(log *logrus.Entry) *ContextLogger {
	return &ContextLogger{entry: log}
}

// GetLogger returns (an immutable) logrus entry from a (mutable) ContextLogger instance
func (l *ContextLogger) GetLogger() *logrus.Entry {
	return l.entry
}

// SetField updates the internal field map
func (l *ContextLogger) SetField(field string, value interface{}) *ContextLogger {
	l.entry = l.entry.WithField(field, value)
	return l
}

// SetFields updates the internal field map with multiple fields at a time
func (l *ContextLogger) SetFields(fields logrus.Fields) *ContextLogger {
	l.entry = l.entry.WithFields(fields)
	return l
}

// SetError updates the internal error
func (l *ContextLogger) SetError(err error) *ContextLogger {
	l.entry = l.entry.WithError(err)
	return l
}

// getContextLogger retrieves the ContextLogger instance from the context
func getContextLogger(ctx context.Context) *ContextLogger {
	if ctxLogger, ok := ctx.Value(loggerKey).(*ContextLogger); ok {
		return ctxLogger
	}

	return nil
}

// Set sets a persistent, mutable logger for the request context.
func Set(ctx context.Context, log *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerKey, NewContextLogger(log))
}

// Get retrieves the logrus Entry from the ContextLogger in a context
// and returns a new logrus Entry if none is found
func Get(ctx context.Context) *logrus.Entry {
	ctxLogger := getContextLogger(ctx)
	if ctxLogger != nil {
		return ctxLogger.GetLogger()
	}

	logger := logrus.NewEntry(logrus.New())

	logger.Warn("no clog exists in the context")

	return logger
}

// SetField adds or updates a field to the ContextLogger in a context
func SetField(ctx context.Context, field string, value interface{}) error {
	ctxLogger := getContextLogger(ctx)

	if ctxLogger == nil {
		return errors.New("no clog exists in the context")
	}

	ctxLogger.SetField(field, value)

	return nil
}

// SetFields adds or updates fields to the ContextLogger in a context
func SetFields(ctx context.Context, fields Fields) error {
	ctxLogger := getContextLogger(ctx)

	if ctxLogger == nil {
		return errors.New("no clog exists in the context")
	}

	ctxLogger.SetFields(logrus.Fields(fields))

	return nil
}

// SetError adds or updates an error to the ContextLogger in a context
func SetError(ctx context.Context, err error) error {
	ctxLogger := getContextLogger(ctx)

	if ctxLogger == nil {
		return errors.New("no clog exists in the context")
	}

	ctxLogger.SetError(err)

	if cherErr, ok := err.(cher.E); ok {
		if len(cherErr.Reasons) > 0 {
			ctxLogger.SetField("error_reasons", cherErr.Reasons)
		}
		if cherErr.Meta != nil {
			ctxLogger.SetField("error_meta", cherErr.Meta)
		}
	}

	return nil
}

// ConfigureTimeoutsAsErrors changes to default behaviour of logging timeouts as info, to log them as errors
func ConfigureTimeoutsAsErrors(ctx context.Context) {
	ctxLogger := getContextLogger(ctx)
	if ctxLogger == nil {
		return
	}

	ctxLogger.timeoutsAsErrors = true
}

// TimeoutsAsErrors determines whether ConfigureTimeoutsAsErrors was called on the context
func TimeoutsAsErrors(ctx context.Context) bool {
	ctxLogger := getContextLogger(ctx)
	if ctxLogger == nil {
		return false
	}

	return ctxLogger.timeoutsAsErrors
}

// DetermineLevel returns a suggested logrus Level type for a given error
func DetermineLevel(err error, timeoutsAsErrors bool) logrus.Level {
	if cherError, ok := err.(cher.E); ok {
		switch cherError.Code {

		// some cher codes have specific log levels
		case cher.BadRequest, cher.RequestTimeout:
			return logrus.WarnLevel
		case cher.ContextCanceled:
			if timeoutsAsErrors {
				return logrus.ErrorLevel
			}
			return logrus.InfoLevel
		case cher.Unknown, cher.CoercionError:
			return logrus.ErrorLevel

		// default cher errors are "handled" so warrant a warning
		default:
			return logrus.WarnLevel
		}
	}

	// non-cher errors are "unhandled" so warrant an error
	return logrus.ErrorLevel
}

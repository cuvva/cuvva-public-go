package request

import (
	"net/http"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/clog"
	"github.com/sirupsen/logrus"
)

type responseWriter struct {
	http.ResponseWriter

	Status int
	Bytes  int64
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.Status == 0 {
		rw.Status = code
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.Status == 0 {
		rw.Status = http.StatusOK
		rw.WriteHeader(http.StatusOK)
	}

	rw.Bytes += int64(len(b))

	return rw.ResponseWriter.Write(b)
}

// Logger returns a middleware handler that wraps subsequent middleware/handlers and logs
// request information AFTER the request has completed. It also injects a request-scoped
// logger on the context which can be set, read and updated using clog lib
//
// Included fields:
//   - Request ID                (request_id)
//   - HTTP Method               (http_method)
//   - HTTP Path                 (http_path)
//   - HTTP Protocol Version     (http_proto)
//   - Remote Address            (http_remote_addr)
//   - User Agent Header         (http_user_agent)
//   - Referer Header            (http_referer)
//   - Duration with unit        (http_duration)
//   - Duration in microseconds  (http_duration_us)
//   - HTTP Status Code          (http_status)
//   - Response in bytes         (http_response_bytes)
//   - Client Version header     (http_client_version)
//   - User Agent header         (http_user_agent)
func Logger(log *logrus.Entry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := GetRequestID(r)

			// create a mutable logger instance which will persist for the request
			// inject pointer to the logger into the request context
			r = r.WithContext(clog.Set(r.Context(), log))

			// panics inside handlers will be logged to standard before propagation
			defer clog.HandlePanic(r.Context(), true)

			clog.SetFields(r.Context(), clog.Fields{
				"request_id": requestID,

				"http_remote_addr":    r.RemoteAddr,
				"http_user_agent":     r.UserAgent(),
				"http_client_version": r.Header.Get("cuvva-client-version"),
				"http_path":           r.URL.Path,
				"http_method":         r.Method,
				"http_proto":          r.Proto,
				"http_referer":        r.Referer(),
			})

			// wrap given response writer with one that tracks status code/bytes written
			rw := &responseWriter{ResponseWriter: w}

			t1 := time.Now()
			next.ServeHTTP(rw, r)
			t2 := time.Now()

			clog.SetFields(r.Context(), clog.Fields{
				"http_duration":       t2.Sub(t1).String(),
				"http_duration_us":    int64(t2.Sub(t1) / time.Microsecond),
				"http_status":         rw.Status,
				"http_response_bytes": rw.Bytes,
			})

			logger := clog.Get(r.Context())

			err := getError(logger)
			logger.Log(determineLevel(err, clog.TimeoutsAsErrors(r.Context())), "request")
		})
	}
}

// getError returns the error if one is set on the log entry
func getError(l *logrus.Entry) error {
	if erri, ok := l.Data[logrus.ErrorKey]; ok {
		if err, ok := erri.(error); ok {
			return err
		}
	}

	return nil
}

// determineLevel returns a suggested logrus Level type based whether an error is present and what type
func determineLevel(err error, timeoutsAsErrors bool) logrus.Level {
	if err != nil {
		return clog.DetermineLevel(err, timeoutsAsErrors)
	}

	// no error, default to info level
	return logrus.InfoLevel
}

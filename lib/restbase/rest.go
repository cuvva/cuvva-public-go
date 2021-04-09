package restbase

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/lib/clog"
)

// Send will encode a request according to the given Accept header (or default to JSON)
// to the request writer. Send respects objects implementing the StatusCoder
// interface. If body is nil, Send will response with `204 No Content`.
func Send(ctx context.Context, w http.ResponseWriter, body interface{}) {
	if body == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err := encode(w, body)
	if err == nil {
		return
	}
	if log := clog.Get(ctx); log != nil {
		log.WithError(err).Error("rest response encoding failed")
	}
}

// ErrorHandler will encode an error type to the request writer using the Send
// method.
func ErrorHandler(ctx context.Context, w http.ResponseWriter, err error) {
	// add it to the reqest log instance
	clog.SetError(ctx, err)

	var body cher.E

	switch err := err.(type) {
	case cher.E:
		body = err

	default:
		body = cher.E{Code: cher.Unknown}
	}

	Send(ctx, w, body)
}

// Wrap handles idiomatic return form and passes it to the ResponseWriter
func Wrap(fn func(context.Context, *http.Request) (interface{}, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		res, err := fn(ctx, r)
		if err != nil {
			// rewrite common errors to internal error standard
			if err == io.EOF {
				err = cher.New(cher.EOF, nil)
			} else if err == io.ErrUnexpectedEOF {
				err = cher.New(cher.UnexpectedEOF, nil)
			} else if strings.Contains(err.Error(), context.Canceled.Error()) {
				err = cher.New(cher.ContextCanceled, nil)
			}

			ErrorHandler(ctx, w, err)
			return
		}

		Send(ctx, w, res)
	}
}

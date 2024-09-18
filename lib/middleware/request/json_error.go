package request

import (
	"encoding/json"
	"net/http"

	"github.com/cuvva/cuvva-public-go/lib/cher"
)

// JSONError will encode the given error as JSON to the client with a HTTP 401 status code.
func JSONError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)

	e := json.NewEncoder(w)

	if chErr, ok := err.(cher.E); ok {
		e.Encode(chErr)
	} else {
		e.Encode(cher.New(cher.Unauthorized, nil, cher.New(cher.Unknown, cher.M{"error": err})))
	}
}

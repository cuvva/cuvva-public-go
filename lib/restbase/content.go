package restbase

import (
	"encoding/json"
	"net/http"
)

// StatusCoder is an interface that may optionally be implemented by a response
// to change the default HTTP 200 status code that will be returned to the client.
type StatusCoder interface {
	StatusCode() int
}

func encode(w http.ResponseWriter, src interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if sc, ok := src.(StatusCoder); ok {
		w.WriteHeader(sc.StatusCode())
	} else if _, ok := src.(error); ok {
		w.WriteHeader(500)
	}

	return json.NewEncoder(w).Encode(src)
}

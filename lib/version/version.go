package version

import (
	"net/http"
)

// Revision is the global build revision from source control as used by
// the ServerHeader middleware, to be set using compiler ldflags like:
// `-ldflags="-X github.com/cuvva/cuvva-public-go/lib/version.Revision=$(git rev-parse HEAD)"`
var Revision = "dev"

// Truncated is up to 7 characters of Revision - so git hashes look as they are
// typically shown in GUIs, on GitHub, etc.
var Truncated = genTruncated()

// Header returns a middleware handler that sets the `Server` header to the
// name of the application and the (truncated) build revision.
func Header(app string) func(http.Handler) http.Handler {
	hdr := app + "/" + Truncated

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Server", hdr)

			next.ServeHTTP(w, r)
		})
	}
}

func genTruncated() string {
	if len(Revision) > 7 {
		return Revision[:7]
	}

	return Revision
}

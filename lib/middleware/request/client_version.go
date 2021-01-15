package request

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/blang/semver"
)

const (
	ClientPlatformIOS     = "ios"
	ClientPlatformAndroid = "android"
)

// keep in sync with the CDN function which does this check globally
var clientHeaderRegexp = regexp.MustCompile(`^(ios|android)-(\d+\.\d+\.\d+)-(\d+)$`)

// ClientVersion represents a parsed client version header string
type ClientVersion struct {
	Platform string
	Version  semver.Version
	Build    int
}

// ClientVersionKey is the key used for holding the parsed client version in the request context map
const ClientVersionKey ContextKey = "Client-Version"

// GetClientVersion retrieves the parsed version from a request
func GetClientVersion(r *http.Request) *ClientVersion {
	return GetClientVersionContext(r.Context())
}

// GetClientVersionContext retrieves the parsed version from a request context
func GetClientVersionContext(ctx context.Context) *ClientVersion {
	if clientVersion, ok := ctx.Value(ClientVersionKey).(*ClientVersion); ok {
		return clientVersion
	}

	return nil
}

// ParseClientVersion attempts to parse the cuvva-client-version HTTP header and add
// it as a struct to the context
func ParseClientVersion(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parsedClientVersion, err := parseVersionHeader(r.Header.Get("cuvva-client-version"))
		if err == nil {
			r = r.WithContext(context.WithValue(r.Context(), ClientVersionKey, parsedClientVersion))
		}

		next.ServeHTTP(w, r)
	})
}

func parseVersionHeader(clientVersionHeader string) (*ClientVersion, error) {
	if clientVersionHeader == "" {
		return nil, errors.New("client version header is empty")
	}

	versionParts := clientHeaderRegexp.FindStringSubmatch(clientVersionHeader)
	if len(versionParts) != 4 {
		return nil, errors.New("header did not match client version pattern")
	}

	platform := versionParts[1]
	version := versionParts[2]
	buildStr := versionParts[3]

	clientSemver, err := semver.Parse(version)
	if err != nil {
		return nil, fmt.Errorf("client version header invalid: semver failed: %w", err)
	}

	build, err := strconv.Atoi(buildStr)
	if err != nil {
		return nil, fmt.Errorf("client version header invalid: build number failed: %w", err)
	}

	return &ClientVersion{
		Platform: platform,
		Version:  clientSemver,
		Build:    build,
	}, nil
}

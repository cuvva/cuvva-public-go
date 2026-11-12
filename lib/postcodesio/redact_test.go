package postcodesio_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/lib/postcodesio"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testPostcode = "EC1R 4AB"

func TestGeocode_ErrorsDoNotContainPostcode(t *testing.T) {
	notFound := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"status":404,"error":"Postcode not found"}`))
	}))
	t.Cleanup(notFound.Close)

	slow := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	t.Cleanup(slow.Close)

	closed := httptest.NewServer(http.NotFoundHandler())
	closed.Close()

	tests := []struct {
		name     string
		baseURL  string
		wantCode string
	}{
		{"not found", notFound.URL, "not_found"},
		{"timeout", slow.URL, cher.RequestTimeout},
		{"transport error", closed.URL, "postcodes_request_failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			client := postcodesio.New(tt.baseURL)
			if tt.name == "timeout" {
				client.Client.Client.Timeout = 50 * time.Millisecond
			}

			_, err := client.Geocode(ctx, testPostcode)
			require.Error(t, err)

			var cErr cher.E
			require.True(t, errors.As(err, &cErr))
			assert.Equal(t, tt.wantCode, cErr.Code)
			assert.Equal(t, "EC1R***", cErr.Meta["postcode"])
			assertNoPostcode(t, cErr)
		})
	}
}

func TestFailoverClient_Geocode_ErrorDoesNotContainPostcode(t *testing.T) {
	closed := httptest.NewServer(http.NotFoundHandler())
	closed.Close()

	fc, err := postcodesio.NewFailoverClient(postcodesio.New(closed.URL), postcodesio.New(closed.URL))
	require.NoError(t, err)

	_, err = fc.Geocode(context.Background(), testPostcode)
	require.Error(t, err)

	var cErr cher.E
	require.True(t, errors.As(err, &cErr))
	assert.Equal(t, "postcode_error", cErr.Code)
	assert.Equal(t, "EC1R***", cErr.Meta["postcode"])
	assert.Len(t, cErr.Reasons, 2)
	assertNoPostcode(t, cErr)
}

func assertNoPostcode(t *testing.T, err cher.E) {
	t.Helper()

	serialised := err.Serialize()
	for _, form := range []string{testPostcode, "EC1R%204AB", "4AB"} {
		assert.False(t, strings.Contains(serialised, form), "error contains %q: %s", form, serialised)
	}
}

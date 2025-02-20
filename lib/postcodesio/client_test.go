package postcodesio_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/cuvva/cuvva-public-go/lib/postcodesio"
)

// TestFallbackReverseGeocode tests if given two clients, the first of which
// will error that it will retrieve a response
func TestFallbackReverseGeocode(t *testing.T) {
	std := postcodesio.New(postcodesio.DefaultBaseURL + "breakthisurl")
	fallback := postcodesio.New(postcodesio.DefaultBaseURL)

	fallbackClient, err := postcodesio.NewFailoverClient(std, fallback)
	if err != nil {
		t.Error(err)
	}

	pc, err := fallbackClient.ReverseGeocode(context.Background(), 51.532322, -0.105826)
	if err != nil {
		t.Error(err)
	}

	if pc == nil {
		// NOTE(sn): Tests if a nil is returned by accident, this does test postcodes.io a little bit
		// too much for my liking and can be removed if it causes issues.
		t.Error(errors.New("no response when one expected"))
	}
}

func TestFailoverClient_Geocode(t *testing.T) {
	std := postcodesio.New(postcodesio.DefaultBaseURL + "breakthisurl")
	fallback := postcodesio.New(postcodesio.DefaultBaseURL)

	fallbackClient, err := postcodesio.NewFailoverClient(std, fallback)
	require.NoError(t, err)

	pc, err := fallbackClient.Geocode(context.Background(), "N1 1AA")
	require.NoError(t, err)

	if pc == nil {
		t.Error(errors.New("no response when one expected"))
	}
}

package postcodesio_test

import (
	"context"
	"testing"

	"github.com/cuvva/cuvva-public-go/lib/postcodesio"
	"github.com/stretchr/testify/require"
)

// TestFallbackReverseGeocode tests if given two clients, the first of which
// will error that it will retrieve a response
func TestFallbackReverseGeocode(t *testing.T) {
	std := postcodesio.New(postcodesio.DefaultBaseURL + "breakthisurl")
	fallback := postcodesio.New(postcodesio.DefaultBaseURL)

	fallbackClient, err := postcodesio.NewFailoverClient(std, fallback)
	require.NoError(t, err)

	pc, err := fallbackClient.ReverseGeocode(context.Background(), 51.532322, -0.105826)
	require.NoError(t, err)

	require.NotNil(t, pc, "no response when one expected")

	require.Equal(t, "N1 9LQ", pc.Postcode)
	require.Equal(t, "Islington", pc.Area) // this tests the custom unmarshaler
	require.Equal(t, 51.53241, pc.Latitude)
	require.Equal(t, -0.106501, pc.Longitude)
}

func TestFailoverClient_Geocode(t *testing.T) {
	std := postcodesio.New(postcodesio.DefaultBaseURL + "breakthisurl")
	fallback := postcodesio.New(postcodesio.DefaultBaseURL)

	fallbackClient, err := postcodesio.NewFailoverClient(std, fallback)
	require.NoError(t, err)

	pc, err := fallbackClient.Geocode(context.Background(), "N1 1AA")
	require.NoError(t, err)

	require.NotNil(t, pc, "no response when one expected")

	require.Equal(t, "N1 1AA", pc.Postcode)
	require.Equal(t, "Islington", pc.Area) // this tests the custom unmarshaler
	require.Equal(t, 51.539746, pc.Latitude)
	require.Equal(t, -0.103053, pc.Longitude)
}

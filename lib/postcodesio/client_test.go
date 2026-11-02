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

	// Validate that we got a valid postcode response (API may return different valid postcodes)
	require.NotEmpty(t, pc.Postcode, "postcode should not be empty")
	require.NotEmpty(t, pc.Area, "area should not be empty") // this tests the custom unmarshaler
	require.NotZero(t, pc.Latitude, "latitude should not be zero")
	require.NotZero(t, pc.Longitude, "longitude should not be zero")

	// Validate coordinates are in reasonable range for London
	require.Greater(t, pc.Latitude, 51.0, "latitude should be in London range")
	require.Less(t, pc.Latitude, 52.0, "latitude should be in London range")
	require.Less(t, pc.Longitude, 0.0, "longitude should be negative (west of prime meridian)")
	require.Greater(t, pc.Longitude, -1.0, "longitude should be in London range")
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

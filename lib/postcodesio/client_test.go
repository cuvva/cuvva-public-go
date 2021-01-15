package postcodesio_test

import (
	"context"
	"errors"
	"testing"

	"github.com/cuvva/cuvva-public-go/lib/postcodesio"
)

// TestSvc checks if the client and fallback client conform to the same interface
func TestSvc(t *testing.T) {
	var client postcodesio.Service = &postcodesio.Client{}
	var fallbackClient postcodesio.Service = &postcodesio.Client{}

	if client == nil {
		t.Error("postcodesio client does not conform to postcodes service interface")
	}

	if fallbackClient == nil {
		t.Error("postcodesio fallback client does not conform to postcodes service interface")
	}
}

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

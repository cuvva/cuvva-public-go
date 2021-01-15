package crpc

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blang/semver"
	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/lib/middleware/request"
	"github.com/stretchr/testify/assert"
)

func TestRequireMinimumClientVersions(t *testing.T) {
	t.Run("no requirement: should pass with a version header", func(t *testing.T) {
		handler := makeClientVersionHandler()
		req := makeClientVersionRequest(&request.ClientVersion{
			Platform: request.ClientPlatformIOS,
			Version:  semver.MustParse("3.6.8"),
			Build:    1337,
		})

		err := handler(httptest.NewRecorder(), req)
		assert.Nil(t, err)
	})
	t.Run("no requirement: should pass with no version header", func(t *testing.T) {
		handler := makeClientVersionHandler()
		req := makeClientVersionRequest(nil)

		err := handler(httptest.NewRecorder(), req)
		assert.Nil(t, err)
	})
	t.Run("single requirement: should pass with same version", func(t *testing.T) {
		handler := makeClientVersionHandler(
			NewVersionRequirement("ios", semver.MustParse("3.6.8")),
		)
		req := makeClientVersionRequest(&request.ClientVersion{
			Platform: request.ClientPlatformIOS,
			Version:  semver.MustParse("3.6.8"),
			Build:    1337,
		})

		err := handler(httptest.NewRecorder(), req)
		assert.Nil(t, err)
	})
	t.Run("single requirement: should pass with greater version", func(t *testing.T) {
		handler := makeClientVersionHandler(
			NewVersionRequirement("ios", semver.MustParse("3.6.8")),
		)
		req := makeClientVersionRequest(&request.ClientVersion{
			Platform: request.ClientPlatformIOS,
			Version:  semver.MustParse("3.6.9"),
			Build:    1337,
		})

		err := handler(httptest.NewRecorder(), req)
		assert.Nil(t, err)
	})
	t.Run("single requirement: should pass with wrong platform", func(t *testing.T) {
		handler := makeClientVersionHandler(
			NewVersionRequirement("ios", semver.MustParse("3.6.8")),
		)
		req := makeClientVersionRequest(&request.ClientVersion{
			Platform: request.ClientPlatformAndroid,
			Version:  semver.MustParse("0.0.1"),
			Build:    1337,
		})

		err := handler(httptest.NewRecorder(), req)
		assert.Nil(t, err)
	})
	t.Run("single requirement: should pass with no header", func(t *testing.T) {
		handler := makeClientVersionHandler(
			NewVersionRequirement("ios", semver.MustParse("3.6.8")),
		)
		req := makeClientVersionRequest(nil)

		err := handler(httptest.NewRecorder(), req)
		assert.Nil(t, err)
	})
	t.Run("single requirement: should fail with correct error if platform matches but version is less", func(t *testing.T) {
		handler := makeClientVersionHandler(
			NewVersionRequirement("ios", semver.MustParse("3.6.8")),
		)
		req := makeClientVersionRequest(&request.ClientVersion{
			Platform: request.ClientPlatformIOS,
			Version:  semver.MustParse("3.6.7"),
			Build:    1337,
		})

		err := handler(httptest.NewRecorder(), req)
		assert.NotNil(t, err)
		assert.NotNil(t, err.(cher.E))
		assert.EqualError(t, err, cher.NoLongerSupported)
	})
	t.Run("multiple requirement: should fail with one failing requirement", func(t *testing.T) {
		handler := makeClientVersionHandler(
			NewVersionRequirement("ios", semver.MustParse("3.6.8")),
			NewVersionRequirement("android", semver.MustParse("3.0.1")),
		)
		req := makeClientVersionRequest(&request.ClientVersion{
			Platform: request.ClientPlatformAndroid,
			Version:  semver.MustParse("2.9.6"),
			Build:    1337,
		})

		err := handler(httptest.NewRecorder(), req)
		assert.NotNil(t, err)
		assert.NotNil(t, err.(cher.E))
		assert.EqualError(t, err, cher.NoLongerSupported)
	})
	t.Run("multiple requirement: should fail with one failing requirement", func(t *testing.T) {
		handler := makeClientVersionHandler(
			NewVersionRequirement("ios", semver.MustParse("3.6.8")),
			NewVersionRequirement("ios", semver.MustParse("3.6.5")),
		)
		req := makeClientVersionRequest(&request.ClientVersion{
			Platform: request.ClientPlatformIOS,
			Version:  semver.MustParse("3.6.6"),
			Build:    1337,
		})

		err := handler(httptest.NewRecorder(), req)
		assert.NotNil(t, err)
		assert.NotNil(t, err.(cher.E))
		assert.EqualError(t, err, cher.NoLongerSupported)
	})
}

func makeClientVersionHandler(requirements ...VersionRequirement) HandlerFunc {
	middleware := RequireMinimumClientVersions(requirements...)
	handlerFunc := func(w http.ResponseWriter, r *Request) error {
		return nil
	}

	return middleware(handlerFunc)
}

func makeClientVersionRequest(ver *request.ClientVersion) *Request {
	if ver == nil {
		return &Request{ctx: context.Background()}
	}

	return &Request{
		ctx: context.WithValue(context.Background(), request.ClientVersionKey, ver),
	}
}

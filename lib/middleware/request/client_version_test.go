package request

import (
	"testing"

	"github.com/blang/semver"
	"github.com/stretchr/testify/assert"
)

func TestParseVersionHeader(t *testing.T) {
	t.Run("empty should error", func(t *testing.T) {
		parsed, err := parseVersionHeader("")
		assert.NotNil(t, err)
		assert.Nil(t, parsed)
	})

	t.Run("invalid should error", func(t *testing.T) {
		parsed, err := parseVersionHeader("a-1-")
		assert.NotNil(t, err)
		assert.Nil(t, parsed)
	})

	t.Run("ios should work", func(t *testing.T) {
		parsed, err := parseVersionHeader("ios-3.6.8-1337")
		assert.Nil(t, err)
		assert.Equal(t, parsed.Platform, ClientPlatformIOS)
		assert.Equal(t, parsed.Version, semver.MustParse("3.6.8"))
		assert.Equal(t, parsed.Build, 1337)
	})

	t.Run("android should work", func(t *testing.T) {
		parsed, err := parseVersionHeader("android-0.0.1-1")
		assert.Nil(t, err)
		assert.Equal(t, parsed.Platform, ClientPlatformAndroid)
		assert.Equal(t, parsed.Version, semver.MustParse("0.0.1"))
		assert.Equal(t, parsed.Build, 1)
	})
}

package jsonschema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKSUIDFormatChecker(t *testing.T) {
	assert.Equal(t, false, ksuidFormatChecker{}.IsFormat("foo"))
	assert.Equal(t, true, ksuidFormatChecker{}.IsFormat("user_000000CBEvdtGRrnrcQKCsSDNNKmR"))
	assert.Equal(t, true, ksuidFormatChecker{}.IsFormat("test_user_000000CBEvefIcrYXsZObXFKZBQrh"))
}

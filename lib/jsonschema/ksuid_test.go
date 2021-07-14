package jsonschema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKSUIDFormatChecker(t *testing.T) {
	c := ksuidFormatChecker{}

	assert.Equal(t, false, c.IsFormat("foo"))
	assert.Equal(t, true, c.IsFormat("user_000000CBEvdtGRrnrcQKCsSDNNKmR"))
	assert.Equal(t, true, c.IsFormat("test_user_000000CBEvefIcrYXsZObXFKZBQrh"))
}

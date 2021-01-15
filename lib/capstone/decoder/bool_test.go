package decoder_test

import (
	"testing"

	"github.com/cuvva/cuvva-public-go/lib/capstone/decoder"
	"github.com/stretchr/testify/assert"
)

func TestBool(t *testing.T) {
	d := decoder.Bool{
		UnavailableValue: "Z",
	}

	tb, err := d.Decode("1")
	if err != nil {
		t.Error(err)
	}

	fb, err := d.Decode("0")
	if err != nil {
		t.Error(err)
	}

	nb, err := d.Decode("Z")
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, tb)
	assert.True(t, true, *tb.(*bool))

	assert.NotNil(t, fb)
	assert.False(t, *fb.(*bool))

	assert.Nil(t, nb)
}

package decoder_test

import (
	"testing"

	"github.com/cuvva/cuvva-public-go/lib/capstone/decoder"
	"github.com/stretchr/testify/assert"
)

func TestFloatDefault(t *testing.T) {
	d := decoder.Float{
		UnavailableValue: "Z",
	}

	v, err := d.Decode("Z")
	if err != nil {
		t.Error(err)
	}

	assert.Nil(t, v)
}

func TestFloatTransform(t *testing.T) {
	d := decoder.Float{
		UnavailableValue: "Z",
		Transform: func(i int) float64 {
			return float64(i) * 100.5
		},
	}

	v, err := d.Decode("0101")
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, v)
	assert.Equal(t, 10150.5, *v.(*float64))
}

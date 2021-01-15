package decoder_test

import (
	"testing"

	"github.com/cuvva/cuvva-public-go/lib/capstone/decoder"
	"github.com/stretchr/testify/assert"
)

func TestIntDefault(t *testing.T) {
	d := decoder.Int{
		UnavailableValue: "X",
	}

	v, err := d.Decode("X")
	if err != nil {
		t.Fatal(err)
	}

	assert.Nil(t, v)
}

func TestLeadingZeroInt(t *testing.T) {
	d := decoder.Int{}

	v, err := d.Decode("01")
	if err != nil {
		t.Error(err)
	}

	i, ok := v.(*int)
	if !ok {
		t.Error("expected i to be an int")
	}

	assert.NotNil(t, i)
	assert.Equal(t, 1, *i)
}

func TestValuesInt(t *testing.T) {
	d := decoder.Int{
		Values: map[string]int{
			"X": -10,
		},
	}

	v, err := d.Decode("X")
	if err != nil {
		t.Error(err)
	}

	i, ok := v.(*int)
	if !ok {
		t.Error("expected i to be an int")
	}

	assert.NotNil(t, i)
	assert.Equal(t, -10, *i)
}

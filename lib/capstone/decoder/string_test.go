package decoder_test

import (
	"errors"
	"testing"

	"github.com/cuvva/cuvva-public-go/lib/capstone/decoder"
	"github.com/stretchr/testify/assert"
)

func TestStringDefault(t *testing.T) {
	d := decoder.String{
		UnavailableValue: "X",
		Values:           map[string]string{},
	}

	v, err := d.Decode("X")
	if err != nil {
		t.Error(err)
	}

	assert.Nil(t, v)
}

func TestStringNotFound(t *testing.T) {
	d := decoder.String{
		UnavailableValue: "ZZ",
		Values:           map[string]string{},
	}

	_, err := d.Decode("EX")
	if !errors.Is(err, decoder.ErrStringNotFound) {
		t.Error("expected ErrStringNotFound")
	}
}

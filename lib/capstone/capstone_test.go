package capstone_test

import (
	"testing"

	"github.com/cuvva/cuvva-public-go/lib/capstone"
	"github.com/cuvva/cuvva-public-go/lib/ptr"
	"github.com/stretchr/testify/assert"
)

func TestExample(t *testing.T) {
	value := "AR081091"
	var exCap capstone.Example

	err := capstone.Parse(value, &exCap)
	if err != nil {
		t.Error(err)
	}

	expected := capstone.Example{
		PolicyStatus:         ptr.String("active"),
		VehicleMatch:         ptr.String("vrn"),
		YearOfNCD:            ptr.Int(8),
		MaleUnemploymentRate: ptr.Float64(0.1),
		EverOnElectoralRoll:  ptr.Bool(true),
	}

	assert.Equal(t, expected, exCap)
}

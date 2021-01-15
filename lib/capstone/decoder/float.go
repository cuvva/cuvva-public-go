package decoder

import (
	"strconv"
)

// Float parses a given string attribute into an int, before passing it
// to a transform function to apply the correct calculation
type Float struct {
	UnavailableValue string
	Transform        func(int) float64
}

// Decode implements capstone.Decoder. See decoder.Float for more details
func (fd Float) Decode(in string) (interface{}, error) {
	if in == fd.UnavailableValue {
		return nil, nil
	}

	v, err := strconv.Atoi(in)
	if err != nil {
		return nil, err
	}

	o := fd.Transform(v)

	return &o, nil
}

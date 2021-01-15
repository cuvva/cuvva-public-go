package decoder

import (
	"errors"
	"fmt"
)

// ErrStringNotFound is returned if the input string did not match any
// field in the Values map
var ErrStringNotFound = errors.New("no string mapping found for value")

// String maps a string to a string held in `values`, mapping up
// any unavailable values to be a nil pointer
type String struct {
	UnavailableValue string
	Values           map[string]string
}

// Decode implements capstone.Decoder. See decoder.String for more details.
func (sd String) Decode(in string) (interface{}, error) {
	if in == sd.UnavailableValue {
		return nil, nil
	}

	if sd.Values == nil {
		return &in, nil
	}

	v, ok := sd.Values[in]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrStringNotFound, in)
	}

	return &v, nil
}

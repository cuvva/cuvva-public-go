package decoder

import (
	"strconv"
)

// Int decodes a leading zero string into an interger. It handles cases
// where the integer is not available through UnavailableValue.
type Int struct {
	UnavailableValue string
	Values           map[string]int
}

// Decode implements capstone.Decoder. See decoder.Int for more details.
func (id Int) Decode(in string) (interface{}, error) {
	if in == id.UnavailableValue {
		return nil, nil
	}

	v, ok := id.Values[in]
	if ok {
		return &v, nil
	}

	v, err := strconv.Atoi(in)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

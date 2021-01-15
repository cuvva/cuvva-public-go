package decoder

// Bool decodes a "1", "0" or UnavailableValue into a *bool
type Bool struct {
	UnavailableValue string
}

// Decode implements capstone.Decoder
func (bd Bool) Decode(in string) (interface{}, error) {
	if in == bd.UnavailableValue {
		return nil, nil
	}

	v := in == "1"

	return &v, nil
}

package icache

import (
	"encoding/xml"
	"fmt"
)

type Bool bool

func (b Bool) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if b {
		return e.EncodeElement("Y", start)
	}

	return e.EncodeElement("N", start)
}

func (b *Bool) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	switch s {
	default:
		return fmt.Errorf("invalid bool - %s", s)
	case "Y":
		*b = true
	case "N":
		*b = false
	}

	return nil
}

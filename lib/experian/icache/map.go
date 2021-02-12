package icache

import (
	"encoding/xml"
	"fmt"
)

type mapContainer struct {
	Content []mapElement `xml:",any"`
}

type mapElement struct {
	XMLName xml.Name
	Content string `xml:",chardata"`
}

type Map map[string]string

func (m Map) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	c := mapContainer{
		Content: make([]mapElement, 0, len(m)),
	}

	for k, v := range m {
		c.Content = append(c.Content, mapElement{
			XMLName: xml.Name{Local: k},
			Content: v,
		})
	}

	e.EncodeElement(c, start)
	return nil
}

func (m *Map) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var c mapContainer
	if err := d.DecodeElement(&c, &start); err != nil {
		return err
	}

	res := make(map[string]string, len(c.Content))

	for _, el := range c.Content {
		if _, ok := res[el.XMLName.Local]; ok {
			return fmt.Errorf("field already defined - %s", el.XMLName.Local)
		}

		res[el.XMLName.Local] = el.Content
	}

	*m = res
	return nil
}

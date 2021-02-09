package icache

import (
	"encoding/xml"
)

// pray for Go generics 🙏

type soapEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`

	Body soapBody
}

type soapBody struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`

	Content InteractiveResponse
}

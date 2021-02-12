package icache

import (
	"encoding/xml"

	"github.com/cuvva/cuvva-public-go/lib/soap"
)

// pray for Go generics 🙏

type soapEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`

	Body soapBody
}

type soapBody struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`

	Fault *soap.Fault

	Content *InteractiveResponse
}

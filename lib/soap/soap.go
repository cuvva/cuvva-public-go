package soap

import (
	"encoding/xml"
)

type Envelope struct {
	XMLName xml.Name `json:"-" xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`

	Header Header
	Body   Body
}

type Header struct {
	XMLName xml.Name `json:"-" xml:"http://schemas.xmlsoap.org/soap/envelope/ Header"`

	Content interface{} `xml:",any"`
}

type Body struct {
	XMLName xml.Name `json:"-" xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`

	Fault *Fault

	Content interface{} `xml:",any"`
}

type Fault struct {
	XMLName xml.Name `json:"-" xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault"`

	Code   string `xml:"faultcode"`
	String string `xml:"faultstring"`
	Actor  string `xml:"faultactor"`
}

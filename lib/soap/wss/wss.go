package wss

import (
	"encoding/xml"
)

type Security struct {
	XMLName xml.Name `xml:"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd Security"`

	Token interface{}
}

type BinarySecurityToken struct {
	XMLName xml.Name `xml:"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd BinarySecurityToken"`

	ValueType string `xml:",attr"`

	Token string `xml:",chardata"`
}

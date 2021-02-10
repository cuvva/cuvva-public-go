package icache

import (
	"encoding/xml"
)

type InteractiveRequest struct {
	XMLName xml.Name `xml:"http://www.uk.experian.com/experian/wbsv/peinteractive/v100 Interactive"`

	Root InputRoot
}

type InputRoot struct {
	XMLName xml.Name `xml:"http://schemas.microsoft.com/BizTalk/2003/Any Root"`

	Input Input
}

type Input struct {
	XMLName xml.Name `xml:"http://schema.uk.experian.com/experian/cems/msgs/v1.1/ConsumerData Input"`

	Control        Control
	Application    Application
	ThirdPartyData ThirdPartyData

	Applicant       Applicant
	LocationDetails LocationDetails
	Residency       Residency
}

type Application struct {
	XMLNS string `xml:"xmlns,attr"`

	ApplicationType string
}

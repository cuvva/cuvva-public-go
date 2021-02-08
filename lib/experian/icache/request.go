package icache

import (
	"encoding/xml"
)

type Interactive struct {
	XMLName xml.Name `xml:"http://www.uk.experian.com/experian/wbsv/peinteractive/v100 Interactive"`

	Root Root
}

type Root struct {
	XMLName xml.Name `xml:"http://schemas.microsoft.com/BizTalk/2003/Any Root"`

	Input Input
}

type Input struct {
	XMLName xml.Name `xml:"http://schema.uk.experian.com/experian/cems/msgs/v1.1/ConsumerData Input"`

	Control        Control
	Application    Application
	ThirdPartyData ThirdPartyData

	Applicants  []Applicant       `xml:"Applicant"`
	Locations   []LocationDetails `xml:"LocationDetails"`
	Residencies []Residency       `xml:"Residency"`
}

type Control struct {
	XMLNS string `xml:"xmlns,attr"`

	ClientAccountNumber string
	ClientBranchNumber  string
	UserIdentity        string
}

type Application struct {
	XMLNS string `xml:"xmlns,attr"`

	ApplicationType string
}

type ThirdPartyData struct {
	XMLNS string `xml:"xmlns,attr"`

	OptOut          Bool
	TransientAssocs Bool
	HHOAllowed      Bool
}

type Applicant struct {
	XMLNS string `xml:"xmlns,attr"`

	ApplicantIdentifier int

	Name ApplicantName

	DateOfBirth Date
}

type ApplicantName struct {
	Forename   string
	MiddleName string
	Surname    string
}

type LocationDetails struct {
	XMLNS string `xml:"xmlns,attr"`

	LocationIdentifier int

	UKLocation LocationDetailsUKLocation
}

type LocationDetailsUKLocation struct {
	Flat            string
	HouseName       string
	HouseNumber     string
	Street          string
	Street2         string
	District        string
	District2       string
	PostTown        string
	County          string
	Postcode        string
	POBox           string
	Country         string
	SharedLetterbox string
}

type Residency struct {
	XMLNS string `xml:"xmlns,attr"`

	ApplicantIdentifier int
	LocationIdentifier  int

	LocationCode string

	ResidencyDateFrom Date
	ResidencyDateTo   Date
}

package icache

type Control struct {
	XMLNS string `json:"-" xml:"xmlns,attr"`

	ExperianReference   *string
	ClientAccountNumber *string
	ClientBranchNumber  *string
	UserIdentity        *string
}

type ThirdPartyData struct {
	XMLNS string `json:"-" xml:"xmlns,attr"`

	OptOut          Bool
	TransientAssocs Bool
	HHOAllowed      Bool
}

type Applicant struct {
	XMLNS string `json:"-" xml:"xmlns,attr"`

	ApplicantIdentifier int

	Name ApplicantName

	DateOfBirth Date
}

type ApplicantName struct {
	Forename   string
	MiddleName *string
	Surname    string
}

type LocationDetails struct {
	XMLNS string `json:"-" xml:"xmlns,attr"`

	LocationIdentifier int

	UKLocation LocationDetailsUKLocation
}

type LocationDetailsUKLocation struct {
	Flat            *string
	HouseName       *string
	HouseNumber     *string
	Street          *string
	Street2         *string
	District        *string
	District2       *string
	PostTown        *string
	County          *string
	Postcode        *string
	POBox           *string
	Country         *string
	SharedLetterbox *string
}

type Residency struct {
	XMLNS string `json:"-" xml:"xmlns,attr"`

	ApplicantIdentifier int
	LocationIdentifier  int

	LocationCode string

	ResidencyDateFrom Date
	ResidencyDateTo   Date
}

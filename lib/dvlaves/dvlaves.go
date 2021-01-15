package dvlaves

const (
	ProductionURI = "https://driver-vehicle-licensing.api.gov.uk/"
	UATURI        = "https://uat.driver-vehicle-licensing.api.gov.uk/"

	TaxStatusTaxed              = "Taxed"
	TaxStatusUntaxed            = "Untaxed"
	TaxStatusNotTaxedForRoadUse = "Not Taxed for on Road Use"
	TaxStatusSORN               = "SORN"

	MOTStatusNotFound  = "No details held by DVLA"
	MOTStatusNoResults = "No results returned"
	MOTStatusNotValid  = "Not valid"
	MOTStatusValid     = "Valid"
)

type VESVRMRequest struct {
	RegistrationNumber string `json:"registrationNumber"`
}

type Vehicle struct {
	RegistrationNumber string `json:"registrationNumber"`

	FirstRegistrationMonth     string  `json:"monthOfFirstRegistration"`
	FirstDVLARegistrationMonth *string `json:"monthOfFirstDvlaRegistration"`
	DateOfLastV5CIssued        string  `json:"dateOfLastV5CIssued"`

	MarkedForExport bool    `json:"markedForExport"`
	TypeApproval    *string `json:"typeApproval"`
	RevenueWeight   *int    `json:"revenueWeight"`
	EuroStatus      *string `json:"euroStatus"`
	ArtEndDate      *string `json:"artEndDate"`

	TaxStatus  string  `json:"taxStatus"`
	TaxDueDate *string `json:"taxDueDate"`

	MOTStatus     string  `json:"motStatus"`
	MOTExpiryDate *string `json:"motExpiryDate"`

	Make                 string `json:"make"`
	Colour               string `json:"colour"`
	EngineCapacity       int    `json:"engineCapacity"` // cc
	CO2Emissions         *int   `json:"co2Emissions"`   // g/km
	RealDrivingEmissions string `json:"realDrivingEmissions"`
	FuelType             string `json:"fuelType"`
	Wheelplan            string `json:"wheelplan"`
	YearOfManufacture    int    `json:"yearOfManufacture"`
}

type Errors struct {
	Errors []Error `json:"errors"`
}

type Error struct {
	Status string `json:"status"`
	Code   string `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

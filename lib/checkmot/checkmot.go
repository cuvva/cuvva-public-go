package checkmot

const ErrNoResults = "no results"
const ErrMultipleVehicles = "multiple vehicles returned"

type Vehicle struct {
	Registration  string `json:"registration"`
	Make          string `json:"make"`
	Model         string `json:"model"`
	FuelType      string `json:"fuelType"`
	PrimaryColour string `json:"primaryColour"`

	// these fields only seem to be present for newer vehicles which haven't been tested yet

	ManufactureYear *string `json:"manufactureYear"`
	DVLAID          *string `json:"dvlaId"`
	MOTTestDueDate  *string `json:"motTestDueDate"`

	// these fields only seem to be present for older vehicles which have been tested

	FirstUsedDate *Date      `json:"firstUsedDate"`
	MOTTests      []*MOTTest `json:"motTests"`
}

type MOTTest struct {
	CompletedDate      Time   `json:"completedDate"`
	TestResult         string `json:"testResult"`
	ExpiryDate         *Date  `json:"expiryDate"`
	OdometerValue      string `json:"odometerValue"`
	OdometerUnit       string `json:"odometerUnit"`
	OdometerResultType string `json:"odometerResultType"`
	MOTTestNumber      string `json:"motTestNumber"`

	// RfRAndComments provide any "reasons for rejection" and comments
	RfRAndComments []*RfROrComment `json:"rfrAndComments"`
}

type RfROrComment struct {
	Text      string `json:"text"`
	Type      string `json:"type"`
	Dangerous bool   `json:"dangerous"`
}

type MOTTestsByCompletedDate []*MOTTest

func (a MOTTestsByCompletedDate) Len() int      { return len(a) }
func (a MOTTestsByCompletedDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a MOTTestsByCompletedDate) Less(i, j int) bool {
	return a[i].CompletedDate.Before(a[j].CompletedDate.Time)
}

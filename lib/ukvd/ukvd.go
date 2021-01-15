package ukvd

type FuelPriceData struct {
	BillingAccount BillingAccount `json:"billingAccount"`
	Request        Request        `json:"request"`
	Response       Response       `json:"response"`
}

type ExtraInformation struct {
}

type BillingAccount struct {
	AccountType      string           `json:"accountType"`
	AccountBalance   float64          `json:"accountBalance"`
	TransactionCost  float64          `json:"transactionCost"`
	ExtraInformation ExtraInformation `json:"extraInformation"`
}

type DataKeys struct {
	Postcode string `json:"postcode"`
}

type Request struct {
	RequestGUID     string   `json:"requestGuid"`
	PackageID       string   `json:"packageId"`
	PackageVersion  int      `json:"packageVersion"`
	ResponseVersion int      `json:"responseVersion"`
	DataKeys        DataKeys `json:"dataKeys"`
}

type Lookup struct {
	StatusCode    string `json:"statusCode"`
	StatusMessage string `json:"statusMessage"`
}

type StatusInformation struct {
	Lookup Lookup `json:"lookup"`
}

type Fuel struct {
	HasUnleaded      bool `json:"hasUnleaded"`
	HasSuperUnleaded bool `json:"hasSuperUnleaded"`
	HasDiesel        bool `json:"hasDiesel"`
	HasPremiumDiesel bool `json:"hasPremiumDiesel"`
	HasLpg           bool `json:"hasLpg"`
	HasEvCharging    bool `json:"hasEvCharging"`
}

type Services struct {
	HasCarWash   bool `json:"hasCarWash"`
	HasTyrePump  bool `json:"hasTyrePump"`
	HasWater     bool `json:"hasWater"`
	HasCashPoint bool `json:"hasCashPoint"`
	HasCarVacuum bool `json:"hasCarVacuum"`
}

type Features struct {
	Fuel     Fuel     `json:"fuel"`
	Services Services `json:"services"`
}

type LatestRecordedPrice struct {
	InPence      float64 `json:"inPence"`
	InGbp        float64 `json:"inGbp"`
	TimeRecorded string  `json:"timeRecorded"`
}

type FuelPriceList struct {
	FuelType            string              `json:"fuelType"`
	LatestRecordedPrice LatestRecordedPrice `json:"latestRecordedPrice"`
}

type FuelStation struct {
	DistanceFromSearchPostcode float64         `json:"distanceFromSearchPostcode"`
	Brand                      string          `json:"brand"`
	Name                       string          `json:"name"`
	Street                     string          `json:"street"`
	Suburb                     string          `json:"suburb"`
	Town                       string          `json:"town"`
	County                     string          `json:"county"`
	Postcode                   string          `json:"postcode"`
	Features                   Features        `json:"features"`
	FuelPriceCount             int             `json:"fuelPriceCount"`
	FuelPriceList              []FuelPriceList `json:"fuelPriceList"`
}

type FuelStationDetails struct {
	FuelStationCount int           `json:"fuelStationCount"`
	SearchRadiusUsed int           `json:"searchRadiusUsed"`
	FuelStationList  []FuelStation `json:"fuelStationList"`
}

type DataItems struct {
	FuelStationDetails FuelStationDetails `json:"fuelStationDetails"`
}

type Response struct {
	StatusCode        string            `json:"statusCode"`
	StatusMessage     string            `json:"statusMessage"`
	StatusInformation StatusInformation `json:"statusInformation"`
	DataItems         DataItems         `json:"dataItems"`
}

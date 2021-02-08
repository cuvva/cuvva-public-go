package givefood

type FoodBank struct {
	Name                      string     `json:"name"`
	Slug                      string     `json:"slug"`
	Address                   string     `json:"address"`
	Postcode                  string     `json:"postcode"`
	Country                   string     `json:"country"`
	LattLong                  string     `json:"latt_long"`
	Closed                    bool       `json:"closed"`
	Phone                     string     `json:"phone"`
	Email                     string     `json:"email"`
	URL                       string     `json:"url"`
	ShoppingListURL           string     `json:"shopping_list_url"`
	CharityNumber             string     `json:"charity_number"`
	CharityRegisteredURL      string     `json:"charity_register_url"`
	Network                   string     `json:"network"`
	ParliamentaryConstituency string     `json:"parliamentary_constituency"`
	MPParty                   string     `json:"mp_party"`
	MP                        string     `json:"mp"`
	District                  string     `json:"district"`
	Ward                      string     `json:"ward"`
	DistanceMi                float64    `json:"distance_mi"`
	NumberNeeds               int        `json:"number_needs"`
	Needs                     string     `json:"needs"`
	NeedID                    string     `json:"need_id"`
	Updated                   string     `json:"updated"`
	UpdatedText               string     `json:"updated_text"`
	Locations                 []Location `json:"locations"`
}

type Location struct {
	Phone                     string `json:"phone"`
	LattLong                  string `json:"latt_long"`
	ParliamentaryConstituency string `json:"parliamentary_constituency"`
	Name                      string `json:"name"`
	MPParty                   string `json:"mp_party"`
	Address                   string `json:"address"`
	District                  string `json:"district"`
	Ward                      string `json:"ward"`
	MP                        string `json:"mp"`
	Postcode                  string `json:"postcode"`
}

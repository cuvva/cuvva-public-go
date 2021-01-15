package capstone

import (
	"github.com/cuvva/cuvva-public-go/lib/capstone/decoder"
)

type Example struct {
	PolicyStatus         *string  `cap:"policy_status" json:"policy_status"`
	VehicleMatch         *string  `cap:"vehicle_match" json:"vehicle_match"`
	YearOfNCD            *int     `cap:"year_of_ncd" json:"year_of_ncd"`
	MaleUnemploymentRate *float64 `cap:"male_unemployment_rate" json:"male_unemployment_rate"`
	EverOnElectoralRoll  *bool    `cap:"ever_on_electoral_roll" json:"ever_on_electoral_roll"`
}

// Format example format is [A|P|C|E|Z][R|I|N|Z][00][00][-][bool]
func (e Example) Format() Format {
	return Format{
		{
			Offset:   0,
			Length:   1,
			TagValue: "policy_status",
			Decoder: decoder.String{
				UnavailableValue: "Z",
				Values: map[string]string{
					"A": "active",
					"P": "pending",
					"C": "cancelled",
					"E": "expired",
				},
			},
		},
		{
			Offset:   1,
			Length:   1,
			TagValue: "vehicle_match",
			Decoder: decoder.String{
				UnavailableValue: "Z",
				Values: map[string]string{
					"R": "vrn",
					"I": "vim",
					"N": "no_match",
				},
			},
		},
		{
			Offset:   2,
			Length:   2,
			TagValue: "year_of_ncd",
			Decoder: decoder.Int{
				UnavailableValue: "Z9",
			},
		},
		{
			Offset:   4,
			Length:   2,
			TagValue: "male_unemployment_rate",
			Decoder: decoder.Float{
				UnavailableValue: "ZZ",
				Transform: func(v int) float64 {
					return float64(v) / 100
				},
			},
		},
		// Offset 6 blank (9)
		{
			Offset:   7,
			Length:   1,
			TagValue: "ever_on_electoral_roll",
			Decoder: decoder.Bool{
				UnavailableValue: "Z",
			},
		},
	}
}

package sanctionsio

import (
	"net/url"
	"regexp"
	"strings"
)

const SourceHMTreasury = "HM TREASURY"

var SearchRequestDOBFormat = regexp.MustCompile(`^\d{4}(?:-\d{2}-\d{2}$)?`)

type SearchRequest struct {
	Name        string
	Sources     []string
	DateOfBirth *string
}

func (r SearchRequest) URLValues() url.Values {
	q := url.Values{
		"name":    {r.Name},
		"sources": {strings.Join(r.Sources, ",")},
	}

	if r.DateOfBirth != nil {
		q.Set("date_of_birth", *r.DateOfBirth)
	}

	return q
}

type SearchResponse struct {
	Count    int                     `json:"count"`
	Next     *string                 `json:"next"`
	Previous *string                 `json:"previous"`
	Results  []*SearchResponseResult `json:"results"`
}

func (r *SearchResponse) HasMatches() bool {
	return len(r.Results) > 0
}

type SearchResponseResult struct {
	Name         string `json:"name" bson:"name"`
	Source       string `json:"source" bson:"source"`
	EntityNumber int    `json:"entity_number" bson:"entity_number"`
	Type         string `json:"type" bson:"type"`
	StartDate    string `json:"start_date" bson:"start_date"`

	Addresses     []string `json:"addresses,omitempty" bson:"addresses,omitempty"`
	Remarks       *string  `json:"remarks,omitempty" bson:"remarks,omitempty"`
	Nationalities []string `json:"nationalities,omitempty" bson:"nationalities,omitempty"`
	DatesOfBirth  []string `json:"dates_of_birth,omitempty" bson:"dates_of_birth,omitempty"`
	PlacesOfBirth []string `json:"places_of_birth,omitempty" bson:"places_of_birth,omitempty"`
	Regime        *string  `json:"regime,omitempty" bson:"regime,omitempty"`
	NINumbers     []string `json:"ni_numbers,omitempty" bson:"ni_numbers,omitempty"`
}

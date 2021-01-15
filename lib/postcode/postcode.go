package postcode

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var baseRegex = regexp.MustCompile(`(?i)^([a-z]{1,2}\d[a-z\d]?)\s*(\d[a-z]{2})$`)
var outcodeRegex = regexp.MustCompile(`^(([A-Z]{1,2})\d{1,2})([A-Z])?$`)

type Parsed struct {
	// FullNormalized is in the form "EC2A 4DP", "E1 4TT"
	FullNormalized string

	// FullCompact is in the form "EC2A4DP", "E14TT"
	FullCompact string

	// Area is in the form "EC", "E"
	Area string

	// District is in the form "EC2", "E1"
	District string

	// SubDistrict is in the form "EC2A", "" - empty string if no subdistrict
	// only postcodes where the outcode ends with a letter have subdistricts
	// e.g. "EC2A" and "W1C" do, but "E1" and "WA14" don't
	SubDistrict string

	// Outcode is in the form "EC2A", "E1"
	Outcode string

	// Sector is in the form "EC2A 4", "E1 4"
	Sector string

	// Incode is in the form "4DP", "4TT"
	Incode string

	// Unit is in the form "DP", "TT"
	Unit string
}

var ErrInvalidPostcode = errors.New("invalid postcode")

func Parse(postcode string) (*Parsed, error) {
	matches := baseRegex.FindStringSubmatch(postcode)
	if matches == nil {
		return nil, ErrInvalidPostcode
	}

	outcode := strings.ToUpper(matches[1])
	incode := strings.ToUpper(matches[2])

	outcodeMatches := outcodeRegex.FindStringSubmatch(outcode)
	if outcodeMatches == nil {
		panic("impossible - outcodeRegex didn't match")
	}

	district := outcodeMatches[1]
	area := outcodeMatches[2]
	subdistrict := outcodeMatches[3]

	if subdistrict != "" {
		subdistrict = outcode
	}

	incodeNumber := incode[0:1]
	unit := incode[1:3]

	return &Parsed{
		FullNormalized: fmt.Sprintf("%s %s", outcode, incode),
		FullCompact:    fmt.Sprintf("%s%s", outcode, incode),
		Area:           area,
		District:       district,
		SubDistrict:    subdistrict,
		Outcode:        outcode,
		Sector:         fmt.Sprintf("%s %s", outcode, incodeNumber),
		Incode:         incode,
		Unit:           unit,
	}, nil
}

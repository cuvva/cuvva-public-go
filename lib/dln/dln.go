package dln

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// ErrInvalidInput is a set error message for when the input to a function is invalid.
var ErrInvalidInput = errors.New("invalid input")

var dateRegex = regexp.MustCompile(`^\d{2}(\d{1})(\d{1})-(\d{2})-(\d{2})$`)
var dlnRegex = regexp.MustCompile(`^[A-Z]{1,5}9{0,4}\d(?:[05][1-9]|[16][0-2])(?:0[1-9]|[12]\d|3[01])\d(?:99|[A-Z][A-Z9])[A-HJ-NPR-X2-9][A-Z]{2}$`)

// UserDetails contains the set of details about a user that can be converted to and from a driving license number.
type UserDetails struct {
	PersonalName string
	FamilyName   string
	Sex          string
	BirthDate    string
}

// Generate a DLN from a set of user details.
func Generate(userDetails UserDetails, includeMiddleName bool) (string, error) {
	sectA, err := generateSectionA(userDetails.FamilyName)
	if err != nil {
		return "", err
	}

	sectB, err := generateSectionB(userDetails.BirthDate, userDetails.Sex)
	if err != nil {
		return "", err
	}

	sectC := generateSectionC(userDetails.PersonalName, includeMiddleName)

	return fmt.Sprintf("%s%s%s", sectA, sectB, sectC), nil
}

// Validate a DLN and set of user details match.
func Validate(dln string, userDetails UserDetails, checkMiddleName bool) (bool, error) {
	subDLN, err := Generate(userDetails, checkMiddleName)

	if err != nil || len(dln) < len(subDLN) {
		return false, err
	}

	return subDLN == dln[:len(subDLN)], nil
}

// Parse a DLN string into a userDetails struct, it will loose data though.
func Parse(dln string, includeMiddleName bool) (*UserDetails, error) {
	if dln == "" || len(dln) != 16 || !dlnRegex.MatchString(dln) {
		return nil, ErrInvalidInput
	}

	familyName := parseSectionA(dln[:5])
	sex, birthDate, err := parseSectionB(dln[5:11])
	initials := parseSectionC(dln[11:13], includeMiddleName)

	if err != nil {
		return nil, err
	}

	return &UserDetails{
		PersonalName: initials,
		FamilyName:   familyName,
		Sex:          sex,
		BirthDate:    birthDate,
	}, nil
}

func generateSectionA(familyName string) (string, error) {
	cleaned := clean(familyName, isAlpha)

	if cleaned == "" {
		return "", ErrInvalidInput
	}

	if strings.HasPrefix(cleaned, "MAC") {
		cleaned = strings.Replace(cleaned, "MAC", "MC", 1)
	}

	return fmt.Sprintf("%s9999", cleaned)[:5], nil
}

func generateSectionB(birthDate, sex string) (string, error) {
	if birthDate == "" || (sex != "F" && sex != "M") {
		return "", ErrInvalidInput
	}

	_, err := time.Parse("2006-01-02", birthDate)
	if err != nil {
		return "", ErrInvalidInput
	}

	out := []byte(dateRegex.ReplaceAllString(birthDate, "$1$3$4$2"))

	if sex == "F" {
		out[1] = out[1] + 5
	}

	return string(out), nil
}

func generateSectionC(personalName string, includeMiddleName bool) string {
	personalName = clean(personalName, isAlphaOrSpace)
	names := strings.Fields(personalName)
	var base []byte

	if includeMiddleName {
		base = []byte{'9', '9'}
	} else {
		base = []byte{'9'}
	}

	for k, v := range names {
		if k >= len(base) {
			break
		}

		base[k] = v[0]
	}

	return string(base)
}

func parseSectionA(a string) string {
	familyName := strings.Title(strings.ToLower(a))

	return strings.Replace(familyName, "9", "", 4)
}

func parseSectionB(in string) (string, string, error) {
	b := []byte(in)

	var sex string
	if b[1] == '5' || b[1] == '6' {
		sex = "F"
		b[1] -= 5
	} else {
		sex = "M"
	}

	birthDate := fmt.Sprintf("%c%c%s", b[0], b[5], b[1:5])
	parsed, err := time.Parse("060102", birthDate)
	if err != nil {
		return "", "", err
	}

	return sex, parsed.Format("2006-01-02"), nil
}

func parseSectionC(c string, includeMiddleName bool) string {
	var initials string

	if c != "99" {
		if !includeMiddleName || c[1] == '9' {
			initials = c[:1]
		} else {
			initials = fmt.Sprintf("%c %c", c[0], c[1])
		}
	}

	return initials
}

func clean(in string, fn func(rune) bool) string {
	upper := strings.ToUpper(in)

	return strings.Map(func(r rune) rune {
		if fn(r) {
			return r
		}

		return -1
	}, upper)
}

func isAlpha(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

func isAlphaOrSpace(r rune) bool {
	return isAlpha(r) || r == ' '
}

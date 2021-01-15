package dln

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		Name              string
		DLN               string
		IncludeMiddleName bool
		UserDetails       UserDetails
		Error             error
	}{
		{
			"basic generation with middle name",
			"MILLE903083GE",
			true,
			UserDetails{
				PersonalName: "George Edward",
				FamilyName:   "Miller",
				Sex:          "M",
				BirthDate:    "1993-03-08",
			},
			nil,
		},
		{
			"basic generation without middle name",
			"MILLE903083G",
			false,
			UserDetails{
				PersonalName: "George",
				FamilyName:   "Miller",
				Sex:          "M",
				BirthDate:    "1993-03-08",
			},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			dln, err := Generate(test.UserDetails, test.IncludeMiddleName)
			if test.Error == nil {
				if assert.NoError(t, err) {
					assert.Equal(t, test.DLN, dln)
				}
			} else {
				assert.Equal(t, test.Error, err)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		Name            string
		DLN             string
		ExpectedOutput  bool
		UserDetails     UserDetails
		CheckMiddleName bool
		Error           error
	}{
		{
			"basic validation with middle name",
			"MILLE903083GE9HT",
			true,
			UserDetails{
				PersonalName: "George Edward",
				FamilyName:   "Miller",
				Sex:          "M",
				BirthDate:    "1993-03-08",
			},
			true,
			nil,
		},
		{
			"basic validation without middle name",
			"MILLE903083GE9HT",
			true,
			UserDetails{
				PersonalName: "George",
				FamilyName:   "Miller",
				Sex:          "M",
				BirthDate:    "1993-03-08",
			},
			false,
			nil,
		},
		{
			"basic validation with middle name 2",
			"BILLI906275JE",
			true,
			UserDetails{
				PersonalName: "James Edward",
				FamilyName:   "Billingham",
				Sex:          "M",
				BirthDate:    "1995-06-27",
			},
			true,
			nil,
		},
		{
			"mismatched dln",
			"WINDSO29332109WS8HT",
			false,
			UserDetails{
				PersonalName: "William Arthur Philip Louis",
				FamilyName:   "Windsor",
				Sex:          "M",
				BirthDate:    "2002-06-21",
			},
			true,
			nil,
		},
		{
			"invalid dln",
			"",
			false,
			UserDetails{
				PersonalName: "James Tiberius",
				FamilyName:   "Kirk",
				Sex:          "M",
				BirthDate:    "2228-03-22",
			},
			false,
			nil,
		},
		{
			"family name missing from user details",
			"BBAEE707015C9",
			false,
			UserDetails{
				PersonalName: "Charles",
				FamilyName:   "",
				Sex:          "M",
				BirthDate:    "1975-07-01",
			},
			false,
			ErrInvalidInput,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			actualOutput, err := Validate(test.DLN, test.UserDetails, test.CheckMiddleName)
			if test.Error == nil {
				if assert.NoError(t, err) {
					assert.Equal(t, test.ExpectedOutput, actualOutput)
				}
			} else {
				assert.Equal(t, test.Error, err)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		Name              string
		DLN               string
		IncludeMiddleName bool
		UserDetails       UserDetails
		Error             error
	}{
		{
			"parse male with middle name",
			"MILLE903083GE9HT",
			true,
			UserDetails{
				PersonalName: "G E",
				FamilyName:   "Mille",
				Sex:          "M",
				BirthDate:    "1993-03-08",
			},
			nil,
		},
		{
			"parse male without middle name",
			"MILLE903083GE9HT",
			false,
			UserDetails{
				PersonalName: "G",
				FamilyName:   "Mille",
				Sex:          "M",
				BirthDate:    "1993-03-08",
			},
			nil,
		},
		{
			"parse male without middle name check",
			"MILLE903083G99HT",
			true,
			UserDetails{
				PersonalName: "G",
				FamilyName:   "Mille",
				Sex:          "M",
				BirthDate:    "1993-03-08",
			},
			nil,
		},
		{
			"parse male with no middle name",
			"MILLE903083999HT",
			true,
			UserDetails{
				PersonalName: "",
				FamilyName:   "Mille",
				Sex:          "M",
				BirthDate:    "1993-03-08",
			},
			nil,
		},
		{
			"parse short family name",
			"CHAN9953228S9XAS",
			true,
			UserDetails{
				PersonalName: "S",
				FamilyName:   "Chan",
				Sex:          "F",
				BirthDate:    "1998-03-22",
			},
			nil,
		},
		{
			"parse mc",
			"MCMIL903083GE9HT",
			true,
			UserDetails{
				PersonalName: "G E",
				FamilyName:   "Mcmil",
				Sex:          "M",
				BirthDate:    "1993-03-08",
			},
			nil,
		},
		{
			"parse female",
			"MILLE953083999HT",
			true,
			UserDetails{
				PersonalName: "",
				FamilyName:   "Mille",
				Sex:          "F",
				BirthDate:    "1993-03-08",
			},
			nil,
		},
		{
			"parse female with middle name",
			"MILLE955107JL9HT",
			true,
			UserDetails{
				PersonalName: "J L",
				FamilyName:   "Mille",
				Sex:          "F",
				BirthDate:    "1997-05-10",
			},
			nil,
		},
		{
			"parse edge birthday upper",
			"KIRK9601019JT9IL",
			true,
			UserDetails{
				PersonalName: "J T",
				FamilyName:   "Kirk",
				Sex:          "M",
				BirthDate:    "1969-01-01",
			},
			nil,
		},
		{
			"parse edge birthday lower",
			"KIRK9612318JT9IL",
			true,
			UserDetails{
				PersonalName: "J T",
				FamilyName:   "Kirk",
				Sex:          "M",
				BirthDate:    "2068-12-31",
			},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			out, err := Parse(test.DLN, test.IncludeMiddleName)
			if test.Error == nil {
				if assert.NoError(t, err) {
					assert.Equal(t, test.UserDetails.PersonalName, out.PersonalName)
					assert.Equal(t, test.UserDetails.FamilyName, out.FamilyName)
					assert.Equal(t, test.UserDetails.Sex, out.Sex)
					assert.Equal(t, test.UserDetails.BirthDate, out.BirthDate)
				}
			} else {
				assert.Equal(t, test.Error, err)
			}
		})
	}
}

func TestGenerateA(t *testing.T) {
	tests := []struct {
		Name           string
		FamilyName     string
		ExpectedOutput string
		Error          error
	}{
		{
			"basic generate a",
			"Miller",
			"MILLE",
			nil,
		},
		{
			"handle short name",
			"Yu",
			"YU999",
			nil,
		},
		{
			"handle no name",
			"",
			"99999",
			ErrInvalidInput,
		},
		{
			"handle mac",
			"Macnamara",
			"MCNAM",
			nil,
		},
		{
			"handle weird characters 1",
			"T'challa",
			"TCHAL",
			nil,
		},
		{
			"handle weird characters 2",
			"--ATx",
			"ATX99",
			nil,
		},
		{
			"handle hypenated family name",
			"Hampshire-Gill",
			"HAMPS",
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			actualOutput, err := generateSectionA(test.FamilyName)
			if test.Error == nil {
				if assert.NoError(t, err) {
					assert.Equal(t, test.ExpectedOutput, actualOutput)
				}
			} else {
				assert.Equal(t, test.Error, err)
			}
		})
	}
}

func TestGenerateB(t *testing.T) {
	tests := []struct {
		Name           string
		BirthDate      string
		Sex            string
		ExpectedOutput string
		Error          error
	}{
		{
			"validate section b male",
			"1993-03-08",
			"M",
			"903083",
			nil,
		},
		{
			"validate section b female",
			"1997-05-10",
			"F",
			"955107",
			nil,
		},
		{
			"invalid date",
			"1993-02-30",
			"F",
			"",
			ErrInvalidInput,
		},
		{
			"invalid date format",
			"2022-00-00",
			"F",
			"",
			ErrInvalidInput,
		},
		{
			"no date given",
			"",
			"F",
			"",
			ErrInvalidInput,
		},
		{
			"no sex given",
			"1997-05-10",
			"",
			"",
			ErrInvalidInput,
		},
		{
			"sex is invalid",
			"1997-05-10",
			"D",
			"",
			ErrInvalidInput,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			actualOutput, err := generateSectionB(test.BirthDate, test.Sex)
			if test.Error == nil {
				if assert.NoError(t, err) {
					assert.Equal(t, test.ExpectedOutput, actualOutput)
				}
			} else {
				assert.Equal(t, test.Error, err)
			}
		})
	}
}

func TestGenerateC(t *testing.T) {
	tests := []struct {
		Name              string
		PersonalName      string
		ExpectedOutput    string
		IncludeMiddleName bool
	}{
		{
			"generate with one name",
			"George",
			"G9",
			true,
		},
		{
			"generate with two names",
			"George Xavier",
			"GX",
			true,
		},
		{
			"ignores third name",
			"George Xander Cage",
			"GX",
			true,
		},
		{
			"does not generate middle name if checked",
			"George Edward",
			"G",
			false,
		},
		{
			"handles odd characters",
			"-George Echo",
			"GE",
			true,
		},
		{
			"handles no names",
			"",
			"99",
			true,
		},
		{
			"handles spaces",
			"    George    Cage  ",
			"GC",
			true,
		},
		{
			"hypenated personal name",
			"Emma-Jane",
			"E9",
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			actualOutput := generateSectionC(test.PersonalName, test.IncludeMiddleName)
			assert.Equal(t, test.ExpectedOutput, actualOutput)
		})
	}
}

func TestClean(t *testing.T) {
	tests := []struct {
		Name           string
		InputString    string
		Func           func(rune) bool
		ExpectedOutput string
	}{
		{
			"basic isAlpha cleaning",
			"James Billingham",
			isAlpha,
			"JAMESBILLINGHAM",
		},
		{
			"isAlphaOrSpace allows spaces",
			"James Billingham 😄",
			isAlphaOrSpace,
			"JAMES BILLINGHAM ",
		},
		{
			"emoji is removed",
			"James 🧟‍♂️",
			isAlpha,
			"JAMES",
		},
		{
			"other characters are removed",
			"James 🧟‍♂️ -- Billing🥉 èas",
			isAlpha,
			"JAMESBILLINGAS",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			actualOutput := clean(test.InputString, test.Func)
			assert.Equal(t, test.ExpectedOutput, actualOutput)
		})
	}
}

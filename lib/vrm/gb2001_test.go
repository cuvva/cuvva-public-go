package vrm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGB2001(t *testing.T) {
	tests := []struct {
		Name   string
		VRM    string
		Result VRM
	}{
		{"Empty", "", nil},
		{"Short", "LB07SE", nil},
		{"BadArea", "0007SEO", nil},
		{"BadAge", "LBAASEO", nil},
		{"InvalidAgeIdentifier01", "LB01SEO", nil},
		{"BadSerial", "LB07000", nil},
		{"Valid", "LB07SEO", &GB2001{Area: "LB", FirstHalf: true, Year: 2007, Serial: "SEO"}},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			v := ParseGB2001(test.VRM)
			assert.Equal(t, test.Result, v)
		})
	}
}

func BenchmarkParseGB2001(b *testing.B) {
	b.SetBytes(7)

	for n := 0; n < b.N; n++ {
		ParseGB2001("LB07SEO")
	}
}

func TestGB2001(t *testing.T) {
	tests := []struct {
		Name         string
		VRM          *GB2001
		String       string
		PrettyString string
	}{
		{"Test", &GB2001{"LB", true, 2007, "SEO"}, "LB07SEO", "LB07 SEO"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, "gb_2001", test.VRM.Format())

			assert.Equal(t, test.String, test.VRM.String())
			assert.Equal(t, test.PrettyString, test.VRM.PrettyString())
		})
	}
}

func TestCalcYearAge(t *testing.T) {
	tests := []struct {
		Name      string
		Age       int
		FirstHalf bool
		Year      int
	}{
		{"Sep2001", 51, false, 2001},
		{"Mar2002", 2, true, 2002},
		{"Sep2002", 52, false, 2002},
		{"Mar2010", 10, true, 2010},
		{"Sep2010", 60, false, 2010},
		{"Mar2050", 50, true, 2050},
		{"Sep2050", 0, false, 2050},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			g := &GB2001{}

			firstHalf, year := g.calcYear(test.Age)
			assert.Equal(t, test.FirstHalf, firstHalf, "first half incorrect")
			assert.Equal(t, test.Year, year, "year incorrect")

			age := g.calcAge(test.FirstHalf, test.Year)
			assert.Equal(t, test.Age, age, "age incorrect")
		})
	}
}

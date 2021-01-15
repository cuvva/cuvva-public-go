package vrm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoerce(t *testing.T) {
	tests := []struct {
		Name    string
		VRM     string
		Formats []Parser
		Results []VRM
	}{
		{"BadVRM", "!!!", nil, nil},
		{"Short", "R8", nil, []VRM{&GB1903{false, "R", "8"}}},
		{"GB2001", "LB07 SEO", nil, []VRM{&GB2001{"LB", true, 2007, "SEO"}}},
		{"NI1966", "ABI 1234", nil, []VRM{&NI1966{"A", "BI", "1234"}}},
		{"Military", "11 AA 11", nil, []VRM{&Military{[]string{"11", "AA", "11"}}}},
		{"Diplomatic", "123 D 456", nil, []VRM{&Diplomatic{'D', 123, 456}}},
		{"GB2001Only", "ABI 1234", []Parser{ParseGB2001}, nil},
		{"GB2001Variant", "LBO7 SE0", []Parser{ParseGB2001}, []VRM{&GB2001{"LB", true, 2007, "SEO"}}},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			results := Coerce(test.VRM, test.Formats...)
			assert.Equal(t, test.Results, results)
		})
	}
}

func TestInfo(t *testing.T) {
	tests := []struct {
		Name      string
		VRM       string
		PrettyVRM string
		Formats   []Parser
		Result    VRM
	}{
		{"NotNormalised", "LB07 SEO", "", nil, nil},
		{"Short", "R8", "R 8", nil, &GB1903{false, "R", "8"}},
		{"GB2001", "LB07SEO", "LB07 SEO", nil, &GB2001{"LB", true, 2007, "SEO"}},
		{"NI1966", "ABI1234", "ABI 1234", nil, &NI1966{"A", "BI", "1234"}},
		{"Military", "11AA11", "11 AA 11", nil, &Military{[]string{"11", "AA", "11"}}},
		{"Diplomatic", "123D456", "123 D 456", nil, &Diplomatic{'D', 123, 456}},
		{"GB2001Only", "ABI1234", "ABI 1234", []Parser{ParseGB2001}, nil},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result := Info(test.VRM, test.Formats...)
			assert.Equal(t, test.Result, result)

			if result != nil {
				assert.Equal(t, test.VRM, result.String())
				assert.Equal(t, test.PrettyVRM, result.PrettyString())
			}
		})
	}
}

func TestValidVRM(t *testing.T) {
	tests := []struct {
		Name  string
		VRM   string
		Valid bool
	}{
		{"Empty", "", false},
		{"Short", "A", false},
		{"Invalid", "AAA!", false},
		{"Valid", "AA1234", true},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, test.Valid, ValidVRM(test.VRM))
		})
	}
}

func BenchmarkValidVRM(b *testing.B) {
	b.SetBytes(6)

	for n := 0; n < b.N; n++ {
		ValidVRM("AA1234")
	}
}

func TestNormaliseVRM(t *testing.T) {
	tests := []struct {
		Name          string
		VRM           string
		NormalisedVRM string
	}{
		{"Empty", "", ""},
		{"Alpha", "AAAA", "AAAA"},
		{"Numeric", "1111", "1111"},
		{"Alphanumeric", "AA11", "AA11"},
		{"Space", "AA 11", "AA11"},
		{"Case", "aa11aa", "AA11AA"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, test.NormalisedVRM, NormaliseVRM(test.VRM))
		})
	}
}

func BenchmarkNormaliseVRM(b *testing.B) {
	b.SetBytes(8)

	for n := 0; n < b.N; n++ {
		NormaliseVRM("LB07 SEO")
	}
}

func TestCombinations(t *testing.T) {
	tests := []struct {
		Name   string
		Prefix string
		Input  string
		Sub    func(rune) rune

		Expected []string
	}{
		{"Empty", "", "", substitutions, []string(nil)},
		{"NoPrefix", "", "LBO7SE0", substitutions, []string{"LB07SE0", "LB07SEO", "LBO7SEO"}},
		{"Prefix", "Foo", "LBO7SE0", substitutions, []string{"FooLB07SE0", "FooLB07SEO", "FooLBO7SEO"}},

		// real world failures
		{"EO16UHF", "", "E016UHF", substitutions, []string{"EO16UHF", "EOI6UHF", "E0I6UHF"}},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			res := combinations(test.Prefix, test.Input, test.Sub)
			assert.Equal(t, test.Expected, res)
		})
	}
}

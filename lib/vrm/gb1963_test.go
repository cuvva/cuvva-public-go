package vrm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGB1963(t *testing.T) {
	tests := []struct {
		Name   string
		VRM    string
		Result VRM
	}{
		{"Empty", "", nil},
		{"Short", "ABC1", nil},
		{"BadSerial", "0BC123D", nil},
		{"BadArea", "A00123D", nil},
		{"BadSequence", "ABCAAAD", nil},
		{"BadYear", "ABC1231", nil},
		{"Prohibited", "IBC123D", nil},
		{"Valid", "ABC123D", &GB1963{Serial: "A", Area: "BC", Sequence: "123", AgeID: "D"}},
		{"ValidShort", "ABC1D", &GB1963{Serial: "A", Area: "BC", Sequence: "1", AgeID: "D"}},
		{"ValidSkipYearI", "ABC123J", &GB1963{Serial: "A", Area: "BC", Sequence: "123", AgeID: "J"}},
		{"ValidSkipYearI", "ABC123P", &GB1963{Serial: "A", Area: "BC", Sequence: "123", AgeID: "P"}},
		{"ValidSkipYearQ", "ABC123R", &GB1963{Serial: "A", Area: "BC", Sequence: "123", AgeID: "R"}},
		{"ValidSkipYearU", "ABC123V", &GB1963{Serial: "A", Area: "BC", Sequence: "123", AgeID: "V"}},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			v := ParseGB1963(test.VRM)
			assert.Equal(t, test.Result, v)
		})
	}
}

func BenchmarkParseGB1963(b *testing.B) {
	b.SetBytes(7)

	for n := 0; n < b.N; n++ {
		ParseGB1963("ABC123D")
	}
}

func TestGB1963(t *testing.T) {
	tests := []struct {
		Name         string
		VRM          *GB1963
		String       string
		PrettyString string
	}{
		{"Long", &GB1963{"A", "BC", "123", "D"}, "ABC123D", "ABC 123D"},
		{"Short", &GB1963{"A", "BC", "1", "D"}, "ABC1D", "ABC 1D"},
		{"SkipYearI", &GB1963{"A", "BC", "123", "J"}, "ABC123J", "ABC 123J"},
		{"SkipYearO", &GB1963{"A", "BC", "123", "P"}, "ABC123P", "ABC 123P"},
		{"SkipYearQ", &GB1963{"A", "BC", "123", "R"}, "ABC123R", "ABC 123R"},
		{"SkipYearU", &GB1963{"A", "BC", "123", "V"}, "ABC123V", "ABC 123V"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, "gb_1963", test.VRM.Format())

			assert.Equal(t, test.String, test.VRM.String())
			assert.Equal(t, test.PrettyString, test.VRM.PrettyString())
		})
	}
}

package vrm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGB1983(t *testing.T) {
	tests := []struct {
		Name   string
		VRM    string
		Result VRM
	}{
		{"Empty", "", nil},
		{"Short", "A1BC", nil},
		{"BadYear", "0123BCD", nil},
		{"BadSeq", "AAAABCD", nil},
		{"BadSeqShort", "AABCD", nil},
		{"BadArea", "A123000", nil},
		{"Prohibited", "A123III", nil},
		{"Valid", "A123BCD", &GB1983{AgeID: "A", Sequence: "123", Serial: "B", Area: "CD"}},
		{"ValidShort", "A1BCD", &GB1983{AgeID: "A", Sequence: "1", Serial: "B", Area: "CD"}},
		{"ValidSkipYearI", "J123BCD", &GB1983{AgeID: "J", Sequence: "123", Serial: "B", Area: "CD"}},
		{"ValidSkipYearO", "P123BCD", &GB1983{AgeID: "P", Sequence: "123", Serial: "B", Area: "CD"}},
		{"ValidSkipYearU", "V123BCD", &GB1983{AgeID: "V", Sequence: "123", Serial: "B", Area: "CD"}},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			v := ParseGB1983(test.VRM)
			assert.Equal(t, test.Result, v)
		})
	}
}

func BenchmarkParseGB1983(b *testing.B) {
	b.SetBytes(7)

	for n := 0; n < b.N; n++ {
		ParseGB1983("A123BCD")
	}
}

func TestGB1983(t *testing.T) {
	tests := []struct {
		Name         string
		VRM          *GB1983
		String       string
		PrettyString string
	}{
		{"Long", &GB1983{"A", "123", "B", "CD"}, "A123BCD", "A123 BCD"},
		{"Short", &GB1983{"A", "1", "B", "CD"}, "A1BCD", "A1 BCD"},
		{"SkipYearI", &GB1983{"J", "123", "B", "CD"}, "J123BCD", "J123 BCD"},
		{"SkipYearO", &GB1983{"P", "123", "B", "CD"}, "P123BCD", "P123 BCD"},
		{"SkipYearU", &GB1983{"V", "123", "B", "CD"}, "V123BCD", "V123 BCD"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, "gb_1983", test.VRM.Format())

			assert.Equal(t, test.String, test.VRM.String())
			assert.Equal(t, test.PrettyString, test.VRM.PrettyString())
		})
	}
}

package vrm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseNI1903(t *testing.T) {
	tests := []struct {
		Name   string
		VRM    string
		Result VRM
	}{
		{"Empty", "", nil},
		{"Short", "AI", nil},
		{"BadSequence", "AI!!!!", nil},
		{"BadSequenceReversed", "1!!!AI", nil},
		{"BadArea", "1234A!", nil},
		{"UnacceptableArea", "AA1234", nil},
		{"Invalid", "!!!!", nil},
		{"Valid", "AI1234", &NI1903{Reversed: false, Area: "AI", Sequence: "1234"}},
		{"ValidShort", "AI1", &NI1903{Reversed: false, Area: "AI", Sequence: "1"}},
		{"ValidReversed", "1234AI", &NI1903{Reversed: true, Area: "AI", Sequence: "1234"}},
		{"ValidShortReversed", "1AI", &NI1903{Reversed: true, Area: "AI", Sequence: "1"}},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			v := ParseNI1903(test.VRM)
			assert.Equal(t, test.Result, v)
		})
	}
}

func BenchmarkParseNI1903(b *testing.B) {
	b.SetBytes(6)

	for n := 0; n < b.N; n++ {
		ParseNI1903("AI1234")
	}
}

func TestNI1903(t *testing.T) {
	tests := []struct {
		Name         string
		VRM          *NI1903
		String       string
		PrettyString string
	}{
		{"Long", &NI1903{false, "AI", "1234"}, "AI1234", "AI 1234"},
		{"Short", &NI1903{false, "AI", "1"}, "AI1", "AI 1"},
		{"LongReversed", &NI1903{true, "AI", "1234"}, "1234AI", "1234 AI"},
		{"ShortReversed", &NI1903{true, "AI", "1"}, "1AI", "1 AI"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, "ni_1903", test.VRM.Format())

			assert.Equal(t, test.String, test.VRM.String())
			assert.Equal(t, test.PrettyString, test.VRM.PrettyString())
		})
	}
}

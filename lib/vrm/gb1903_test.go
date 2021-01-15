package vrm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGB1903(t *testing.T) {
	tests := []struct {
		Name   string
		VRM    string
		Result VRM
	}{
		{"Empty", "", nil},
		{"Short", "A", nil},
		{"Invalid", "!!", nil},
		{"BadSequence", "AA!!!!", nil},
		{"BadArea", "1234!!", nil},
		{"Prohibited", "II1234", nil},
		{"Valid", "AB1234", &GB1903{Reversed: false, Area: "AB", Sequence: "1234"}},
		{"ValidShort", "A1", &GB1903{Reversed: false, Area: "A", Sequence: "1"}},
		{"ValidReversed", "1234AB", &GB1903{Reversed: true, Area: "AB", Sequence: "1234"}},
		{"ValidShortReversed", "1A", &GB1903{Reversed: true, Area: "A", Sequence: "1"}},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			v := ParseGB1903(test.VRM)
			assert.Equal(t, test.Result, v)
		})
	}
}

func BenchmarkParseGB1903(b *testing.B) {
	b.SetBytes(5)

	for n := 0; n < b.N; n++ {
		ParseGB1903("AB123")
	}
}

func TestGB1903(t *testing.T) {
	tests := []struct {
		Name         string
		VRM          *GB1903
		String       string
		PrettyString string
	}{
		{"Long", &GB1903{false, "AB", "1234"}, "AB1234", "AB 1234"},
		{"Short", &GB1903{false, "A", "1"}, "A1", "A 1"},
		{"LongReversed", &GB1903{true, "AB", "1234"}, "1234AB", "1234 AB"},
		{"ShortReversed", &GB1903{true, "A", "1"}, "1A", "1 A"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, "gb_1903", test.VRM.Format())

			assert.Equal(t, test.String, test.VRM.String())
			assert.Equal(t, test.PrettyString, test.VRM.PrettyString())
		})
	}
}

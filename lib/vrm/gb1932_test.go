package vrm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGB1932(t *testing.T) {
	tests := []struct {
		Name   string
		VRM    string
		Result VRM
	}{
		{"Empty", "", nil},
		{"Short", "ABC", nil},
		{"Invalid", "!!!!", nil},
		{"BadAreaReversed", "000A00", nil},
		{"BadSequence", "ABCAAA", nil},
		{"BadSequenceReversed", "0AAABC", nil},
		{"Prohibited", "III123", nil},
		{"Valid", "ABC123", &GB1932{Reversed: false, Serial: "A", Area: "BC", Sequence: "123"}},
		{"ValidShort", "ABC1", &GB1932{Reversed: false, Serial: "A", Area: "BC", Sequence: "1"}},
		{"ValidReversed", "123ABC", &GB1932{Reversed: true, Serial: "A", Area: "BC", Sequence: "123"}},
		{"ValidShortReversed", "1ABC", &GB1932{Reversed: true, Serial: "A", Area: "BC", Sequence: "1"}},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			v := ParseGB1932(test.VRM)
			assert.Equal(t, test.Result, v)
		})
	}
}

func BenchmarkGB1932(b *testing.B) {
	b.SetBytes(6)

	for n := 0; n < b.N; n++ {
		ParseGB1932("ABC123")
	}
}

func TestGB1932(t *testing.T) {
	tests := []struct {
		Name         string
		VRM          *GB1932
		String       string
		PrettyString string
	}{
		{"Long", &GB1932{false, "A", "BC", "123"}, "ABC123", "ABC 123"},
		{"Short", &GB1932{false, "A", "BC", "1"}, "ABC1", "ABC 1"},
		{"LongReversed", &GB1932{true, "A", "BC", "123"}, "123ABC", "123 ABC"},
		{"ShortReversed", &GB1932{true, "A", "BC", "1"}, "1ABC", "1 ABC"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, "gb_1932", test.VRM.Format())

			assert.Equal(t, test.String, test.VRM.String())
			assert.Equal(t, test.PrettyString, test.VRM.PrettyString())
		})
	}
}

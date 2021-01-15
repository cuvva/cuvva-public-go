package vrm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseNI1966(t *testing.T) {
	tests := []struct {
		Name   string
		VRM    string
		Result VRM
	}{
		{"Empty", "", nil},
		{"Short", "ABI", nil},
		{"BadSerial", "!BI1234", nil},
		{"BadArea", "AAA1234", nil},
		{"BadSequence", "ABI!!!!", nil},
		{"Valid", "ABI1234", &NI1966{Serial: "A", Area: "BI", Sequence: "1234"}},
		{"ValidShort", "ABI1", &NI1966{Serial: "A", Area: "BI", Sequence: "1"}},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			v := ParseNI1966(test.VRM)
			assert.Equal(t, test.Result, v)
		})
	}
}

func BenchmarkParseNI1966(b *testing.B) {
	b.SetBytes(7)

	for n := 0; n < b.N; n++ {
		ParseNI1966("ABI1234")
	}
}

func TestNI1966(t *testing.T) {
	tests := []struct {
		Name         string
		VRM          *NI1966
		String       string
		PrettyString string
	}{
		{"Long", &NI1966{"A", "BI", "1234"}, "ABI1234", "ABI 1234"},
		{"Short", &NI1966{"A", "BI", "1"}, "ABI1", "ABI 1"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, "ni_1966", test.VRM.Format())

			assert.Equal(t, test.String, test.VRM.String())
			assert.Equal(t, test.PrettyString, test.VRM.PrettyString())
		})
	}
}

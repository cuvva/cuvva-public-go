package vrm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDiplomatic(t *testing.T) {
	tests := []struct {
		Name   string
		VRM    string
		Result VRM
	}{
		{"Empty", "", nil},
		{"Short", "123X45", nil},
		{"Invalid", "!!!!!!!", nil},
		{"BadEntity", "!!!X456", nil},
		{"BadSerial", "123D!!!", nil},
		{"Valid", "123X456", &Diplomatic{Type: 'X', Entity: 123, Serial: 456}},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			v := ParseDiplomatic(test.VRM)
			assert.Equal(t, test.Result, v)
		})
	}
}

func BenchmarkParseDiplomatic(b *testing.B) {
	b.SetBytes(7)

	for n := 0; n < b.N; n++ {
		ParseDiplomatic("123X456")
	}
}

func TestDiplomatic(t *testing.T) {
	tests := []struct {
		Name         string
		VRM          *Diplomatic
		String       string
		PrettyString string
	}{
		{"X", &Diplomatic{'X', 123, 456}, "123X456", "123 X 456"},
		{"D", &Diplomatic{'D', 123, 456}, "123D456", "123 D 456"},
		{"D", &Diplomatic{'D', 1, 5}, "001D005", "001 D 005"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, "diplomatic", test.VRM.Format())

			assert.Equal(t, test.String, test.VRM.String())
			assert.Equal(t, test.PrettyString, test.VRM.PrettyString())
		})
	}
}

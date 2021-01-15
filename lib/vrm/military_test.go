package vrm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMilitary(t *testing.T) {
	tests := []struct {
		Name   string
		VRM    string
		Result VRM
	}{
		{"Empty", "", nil},
		{"Short", "AA11A", nil},
		{"Invalid", "!!!!!!", nil},
		{"Valid", "11AA11", &Military{[]string{"11", "AA", "11"}}},
		{"ValidReversed", "AA11AA", &Military{[]string{"AA", "11", "AA"}}},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			v := ParseMilitary(test.VRM)
			assert.Equal(t, test.Result, v)
		})
	}
}

func BenchmarkParseMilitary(b *testing.B) {
	b.SetBytes(6)

	for n := 0; n < b.N; n++ {
		ParseMilitary("11AA11")
	}
}

func TestMiliary(t *testing.T) {
	tests := []struct {
		Name         string
		VRM          *Military
		String       string
		PrettyString string
	}{
		{"Valid", &Military{[]string{"11", "AA", "11"}}, "11AA11", "11 AA 11"},
		{"ValidReversed", &Military{[]string{"AA", "11", "AA"}}, "AA11AA", "AA 11 AA"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, "military", test.VRM.Format())

			assert.Equal(t, test.String, test.VRM.String())
			assert.Equal(t, test.PrettyString, test.VRM.PrettyString())
		})
	}
}

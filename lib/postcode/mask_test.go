package postcode

import "testing"

func TestMask(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"normalised", "EC1R 4AB", "EC1R***"},
		{"compact", "EC1R4AB", "EC1R***"},
		{"lowercase", "ec1r 4ab", "EC1R***"},
		{"short outcode", "E1 4TT", "E1***"},
		{"outcode only", "EC1R", "***"},
		{"invalid", "NOT A POSTCODE", "***"},
		{"empty", "", "***"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Mask(tt.input); got != tt.want {
				t.Errorf("Mask(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

package cdep

import (
	"testing"
)

func TestParseTypeArg(t *testing.T) {
	cases := map[string]string{
		"service":   "service",
		"services":  "service",
		"lambda":    "lambda",
		"lambdas":   "lambda",
		"terra":     "terra",
		"agentcore": "agentcore",
	}

	for in, want := range cases {
		t.Run(in, func(t *testing.T) {
			got, err := ParseTypeArg(in)
			if err != nil {
				t.Fatalf("ParseTypeArg(%q) error: %v", in, err)
			}
			if got != want {
				t.Fatalf("ParseTypeArg(%q) = %q, want %q", in, got, want)
			}
		})
	}

	if _, err := ParseTypeArg("nope"); err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestValidateCommitHash(t *testing.T) {
	tests := []struct {
		name    string
		commit  string
		wantErr bool
	}{
		{
			name:    "valid 40-character hash",
			commit:  "7fbcf728fe0682e62f3c5e179ebf2e846a99d3c2",
			wantErr: false,
		},
		{
			name:    "short hash",
			commit:  "7fbcf728fe",
			wantErr: true,
		},
		{
			name:    "branch name master",
			commit:  "master",
			wantErr: true,
		},
		{
			name:    "branch name main",
			commit:  "main",
			wantErr: true,
		},
		{
			name:    "empty string",
			commit:  "",
			wantErr: true,
		},
		{
			name:    "invalid characters",
			commit:  "7fbcf728fe0682e62f3c5e179ebf2e846a99d3cg",
			wantErr: true,
		},
		{
			name:    "uppercase letters",
			commit:  "7FBCF728FE0682E62F3C5E179EBF2E846A99D3C2",
			wantErr: true,
		},
		{
			name:    "too long",
			commit:  "7fbcf728fe0682e62f3c5e179ebf2e846a99d3c2a",
			wantErr: true,
		},
		{
			name:    "too short",
			commit:  "7fbcf728f",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCommitHash(tt.commit)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCommitHash(%q) error = %v, wantErr %v", tt.commit, err, tt.wantErr)
			}
		})
	}
}

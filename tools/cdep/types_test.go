package cdep

import (
	"testing"
)

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

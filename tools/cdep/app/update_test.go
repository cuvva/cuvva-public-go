package app

import (
	"context"
	"testing"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/tools/cdep/parsers"
)

// The agentcore guards sit at the top of Update and return before any
// git/config/AWS work, so they can be exercised in isolation.
func TestUpdateAgentcoreGuards(t *testing.T) {
	tests := []struct {
		name string
		req  *parsers.Params
		want string
	}{
		{
			name: "non-_system env rejected (incl. all)",
			req:  &parsers.Params{Type: "agentcore", System: "nonprod", Environment: "all", Branch: "main", Items: []string{"a"}},
			want: "agentcore_requires_system_env",
		},
		{
			name: "pinned commit cannot span multiple agents",
			req:  &parsers.Params{Type: "agentcore", System: "nonprod", Environment: "_system", Branch: "main", Commit: "7fbcf728fe0682e62f3c5e179ebf2e846a99d3c2", Items: []string{"a", "b"}},
			want: "commit_requires_single_agent",
		},
		{
			name: "prod agentcore must be on default branch",
			req:  &parsers.Params{Type: "agentcore", System: "prod", Environment: "_system", Branch: "some-feature", Items: []string{"a"}},
			want: "invalid_operation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := App{}.Update(context.Background(), tt.req, nil)
			if code := cher.Coerce(err).Code; code != tt.want {
				t.Fatalf("got %q, want %q (err: %v)", code, tt.want, err)
			}
		})
	}
}

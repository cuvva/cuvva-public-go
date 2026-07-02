package paths

import "testing"

func TestGetPathForAgentcore(t *testing.T) {
	got := GetPathForAgentcore("/repo", "nonprod", "_system", "my-agent")
	want := "/repo/nonprod/_system/agentcore/my-agent.json"

	if got != want {
		t.Fatalf("GetPathForAgentcore = %q, want %q", got, want)
	}
}

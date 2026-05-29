package entirecli

import "testing"

func TestParseAgentList(t *testing.T) {
	agents, err := ParseAgentList(`Agents:
  ✓ claude-code
    codex
  ✓ opencode

Run 'entire agent add <name>' to enable another agent.
`)
	if err != nil {
		t.Fatalf("ParseAgentList: %v", err)
	}
	if len(agents) != 3 {
		t.Fatalf("expected 3 agents, got %#v", agents)
	}
	if agents[0] != (Agent{Name: "claude-code", Installed: true}) {
		t.Fatalf("first agent = %#v", agents[0])
	}
	if agents[1] != (Agent{Name: "codex", Installed: false}) {
		t.Fatalf("second agent = %#v", agents[1])
	}
	if agents[2] != (Agent{Name: "opencode", Installed: true}) {
		t.Fatalf("third agent = %#v", agents[2])
	}
}

func TestSpecsForConfiguredAgents(t *testing.T) {
	specs := SpecsForConfiguredAgents([]Agent{
		{Name: "claude-code", Installed: true},
		{Name: "codex", Installed: false},
		{Name: "custom-agent", Installed: true},
	})
	if len(specs) != 2 {
		t.Fatalf("expected 2 specs, got %#v", specs)
	}
	if specs[0].Binary != "claude" {
		t.Fatalf("claude binary = %q", specs[0].Binary)
	}
	if specs[1].Binary != "custom-agent" {
		t.Fatalf("custom binary = %q", specs[1].Binary)
	}
}

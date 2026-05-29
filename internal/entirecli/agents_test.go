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

func TestBuiltInYoloArgs(t *testing.T) {
	tests := map[string][]string{
		"claude-code":     {"--dangerously-skip-permissions"},
		"codex":           {"--dangerously-bypass-approvals-and-sandbox"},
		"gemini":          {"--yolo"},
		"copilot-cli":     {"--yolo"},
		"cursor":          {"--yolo"},
		"factoryai-droid": {"exec", "--skip-permissions-unsafe"},
		"opencode":        {"run", "--dangerously-skip-permissions"},
	}
	for name, want := range tests {
		t.Run(name, func(t *testing.T) {
			spec, ok := SpecFor(name)
			if !ok {
				t.Fatalf("missing spec")
			}
			if len(spec.YoloArgs) != len(want) {
				t.Fatalf("YoloArgs = %#v, want %#v", spec.YoloArgs, want)
			}
			for i := range want {
				if spec.YoloArgs[i] != want[i] {
					t.Fatalf("YoloArgs = %#v, want %#v", spec.YoloArgs, want)
				}
			}
		})
	}
}

func TestPiYoloWarning(t *testing.T) {
	spec, ok := SpecFor("pi")
	if !ok {
		t.Fatalf("missing pi spec")
	}
	if len(spec.YoloArgs) != 0 {
		t.Fatalf("Pi YoloArgs = %#v, want none", spec.YoloArgs)
	}
	if spec.YoloWarning == "" {
		t.Fatal("Pi should warn that --yolo is unsupported")
	}
}

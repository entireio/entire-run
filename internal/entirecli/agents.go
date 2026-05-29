// Package entirecli discovers Entire CLI agent integrations from the parent
// `entire` command. External plugins cannot import the parent CLI internals,
// so discovery shells out to `entire agent list`.
package entirecli

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// Agent represents one row from `entire agent list`.
type Agent struct {
	Name      string
	Installed bool
}

// AgentSpec maps an Entire agent registry key to the host CLI command this
// launcher should execute.
type AgentSpec struct {
	EntireName  string
	Display     string
	Binary      string
	DefaultArgs []string
	YoloArgs    []string
	YoloWarning string
}

var AgentSpecs = []AgentSpec{
	{EntireName: "claude-code", Display: "Claude Code", Binary: "claude", YoloArgs: []string{"--dangerously-skip-permissions"}},
	{EntireName: "codex", Display: "Codex", Binary: "codex", YoloArgs: []string{"--dangerously-bypass-approvals-and-sandbox"}},
	{EntireName: "gemini", Display: "Gemini CLI", Binary: "gemini", YoloArgs: []string{"--yolo"}},
	{EntireName: "copilot-cli", Display: "Copilot CLI", Binary: "copilot", YoloArgs: []string{"--yolo"}},
	{EntireName: "cursor", Display: "Cursor", Binary: "agent", YoloArgs: []string{"--yolo"}},
	{EntireName: "factoryai-droid", Display: "Factory AI Droid", Binary: "droid", YoloArgs: []string{"exec", "--skip-permissions-unsafe"}},
	{EntireName: "opencode", Display: "OpenCode", Binary: "opencode", YoloArgs: []string{"run", "--dangerously-skip-permissions"}},
	{EntireName: "pi", Display: "Pi", Binary: "pi", YoloWarning: "warning: Pi does not support --yolo; launching without YOLO mode\n"},
}

func SpecFor(name string) (AgentSpec, bool) {
	for _, spec := range AgentSpecs {
		if spec.EntireName == name {
			return spec, true
		}
	}
	return AgentSpec{}, false
}

func ListConfigured(ctx context.Context) ([]Agent, error) {
	binPath, err := exec.LookPath("entire")
	if err != nil {
		return nil, fmt.Errorf("entire CLI not on PATH: %w", err)
	}
	cmd := exec.CommandContext(ctx, binPath, "agent", "list")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("entire agent list: %w", err)
	}
	return ParseAgentList(string(out))
}

func ConfiguredSpecsWithError(ctx context.Context) ([]AgentSpec, error) {
	agents, err := ListConfigured(ctx)
	if err != nil {
		return nil, err
	}
	return SpecsForConfiguredAgents(agents), nil
}

func SpecsForConfiguredAgents(agents []Agent) []AgentSpec {
	specs := make([]AgentSpec, 0, len(agents))
	for _, agent := range agents {
		if !agent.Installed {
			continue
		}
		spec, ok := SpecFor(agent.Name)
		if !ok {
			spec = AgentSpec{
				EntireName: agent.Name,
				Display:    agent.Name,
				Binary:     agent.Name,
			}
		}
		specs = append(specs, spec)
	}
	return specs
}

// ParseAgentList consumes the marker-prefixed `entire agent list` output:
//
//	Agents:
//	  ✓ claude-code
//	    codex
//
// The check mark means hooks are installed for this repository.
func ParseAgentList(s string) ([]Agent, error) {
	var agents []Agent
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		raw := scanner.Text()
		if !strings.HasPrefix(raw, "  ") {
			continue
		}
		line := strings.TrimLeft(raw, " ")
		installed := false
		if strings.HasPrefix(line, "✓ ") {
			installed = true
			line = strings.TrimPrefix(line, "✓ ")
		}
		name := strings.TrimSpace(line)
		if name == "" || strings.ContainsAny(name, " .") {
			continue
		}
		agents = append(agents, Agent{Name: name, Installed: installed})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan agent list: %w", err)
	}
	return agents, nil
}

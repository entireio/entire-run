package cli

import (
	"bytes"
	"context"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/entireio/entire-run/internal/entirecli"
)

func TestRunLauncherLaunchesSelectedAgent(t *testing.T) {
	restore := stubLauncher(t, []entirecli.AgentSpec{
		{EntireName: "claude-code", Display: "Claude Code", Binary: "claude"},
		{EntireName: "codex", Display: "Codex", Binary: "codex"},
	})
	defer restore()

	var launched entirecli.AgentSpec
	var launchedArgs []string
	launchAgentFn = func(_ context.Context, spec entirecli.AgentSpec, args []string) error {
		launched = spec
		launchedArgs = args
		return nil
	}

	if err := runLauncher(context.Background(), ioDiscard{}, []string{"codex", "--ask-for-approval"}); err != nil {
		t.Fatalf("runLauncher: %v", err)
	}
	if launched.EntireName != "codex" {
		t.Fatalf("launched %q", launched.EntireName)
	}
	if !reflect.DeepEqual(launchedArgs, []string{"--ask-for-approval"}) {
		t.Fatalf("args = %#v", launchedArgs)
	}
}

func TestRunLauncherRejectsDisabledAgent(t *testing.T) {
	restore := stubLauncher(t, []entirecli.AgentSpec{
		{EntireName: "codex", Display: "Codex", Binary: "codex"},
	})
	defer restore()

	err := runLauncher(context.Background(), ioDiscard{}, []string{"gemini"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "not enabled") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPickAgentNonInteractiveDefaultsToFirst(t *testing.T) {
	old := stdinIsTerminalFn
	stdinIsTerminalFn = func(*os.File) bool { return false }
	defer func() { stdinIsTerminalFn = old }()

	var out bytes.Buffer
	spec, err := pickAgent(&out, os.Stdin, []entirecli.AgentSpec{
		{EntireName: "codex", Display: "Codex", Binary: "codex"},
		{EntireName: "gemini", Display: "Gemini CLI", Binary: "gemini"},
	})
	if err != nil {
		t.Fatalf("pickAgent: %v", err)
	}
	if spec.EntireName != "codex" {
		t.Fatalf("picked %q", spec.EntireName)
	}
	if !strings.Contains(out.String(), "Non-interactive stdin") {
		t.Fatalf("missing fallback message: %s", out.String())
	}
}

func stubLauncher(t *testing.T, specs []entirecli.AgentSpec) func() {
	t.Helper()
	oldSpecs := configuredSpecsFn
	oldLaunch := launchAgentFn
	oldSelect := selectAgentFn
	configuredSpecsFn = func(context.Context) ([]entirecli.AgentSpec, error) {
		return specs, nil
	}
	launchAgentFn = func(context.Context, entirecli.AgentSpec, []string) error {
		return nil
	}
	selectAgentFn = func(specs []entirecli.AgentSpec) (entirecli.AgentSpec, error) {
		return specs[0], nil
	}
	return func() {
		configuredSpecsFn = oldSpecs
		launchAgentFn = oldLaunch
		selectAgentFn = oldSelect
	}
}

type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) {
	return len(p), nil
}

package cli

import (
	"bytes"
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/entireio/entire-run/internal/entirecli"
	"github.com/spf13/cobra"
)

func execute(t *testing.T, cmd *cobra.Command, args ...string) (string, error) {
	t.Helper()

	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stdout)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return stdout.String(), err
}

func TestDoctorRequiresPluginDataDir(t *testing.T) {
	cmd := NewRootCommand(Options{Version: "test-version"})
	out, err := execute(t, cmd, "doctor")
	if err == nil {
		t.Fatal("doctor returned nil error without plugin data dir")
	}
	if !strings.Contains(out, "ENTIRE_PLUGIN_DATA_DIR=<unset>") {
		t.Fatalf("doctor output missing environment dump:\n%s", out)
	}
}

func TestDoctorCreatesWritablePluginDataDir(t *testing.T) {
	dataDir := t.TempDir() + "/plugin-data"
	cmd := NewRootCommand(Options{
		Version: "test-version",
		Env:     EntireEnv{PluginDataDir: dataDir},
	})

	out, err := execute(t, cmd, "doctor")
	if err != nil {
		t.Fatalf("doctor: %v", err)
	}
	if !strings.Contains(out, "plugin data dir: writable") {
		t.Fatalf("doctor output missing writable status:\n%s", out)
	}
}

func TestRootYoloFlagPassesThroughAgentFlags(t *testing.T) {
	restore := stubLauncher(t, []entirecli.AgentSpec{
		{EntireName: "codex", Display: "Codex", Binary: "codex", YoloArgs: []string{"--dangerously-bypass-approvals-and-sandbox"}},
	})
	defer restore()

	var launchedArgs []string
	launchAgentFn = func(_ context.Context, _ entirecli.AgentSpec, args []string) error {
		launchedArgs = args
		return nil
	}

	cmd := NewRootCommand(Options{Version: "test-version"})
	_, err := execute(t, cmd, "--yolo", "codex", "--ask-for-approval", "never")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	want := []string{"--dangerously-bypass-approvals-and-sandbox", "--ask-for-approval", "never"}
	if !reflect.DeepEqual(launchedArgs, want) {
		t.Fatalf("args = %#v, want %#v", launchedArgs, want)
	}
}

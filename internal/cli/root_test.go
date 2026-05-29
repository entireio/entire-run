package cli

import (
	"bytes"
	"strings"
	"testing"

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

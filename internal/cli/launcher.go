package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/entireio/entire-run/internal/entirecli"
	"github.com/entireio/entire-run/internal/ui"
)

var (
	configuredSpecsFn = entirecli.ConfiguredSpecsWithError
	launchAgentFn     = launchAgent
	stdinIsTerminalFn = isTerminal
	selectAgentFn     = selectAgentInteractive
)

func runLauncher(ctx context.Context, out io.Writer, args []string) error {
	specs, err := configuredSpecsFn(ctx)
	if err != nil {
		return err
	}
	if len(specs) == 0 {
		return fmt.Errorf("no Entire agents are enabled in this repository; run `entire agent add <name>` first")
	}

	var spec entirecli.AgentSpec
	var agentArgs []string
	if len(args) > 0 {
		var ok bool
		spec, ok = findSpec(specs, args[0])
		if !ok {
			return fmt.Errorf("agent %q is not enabled for Entire in this repository", args[0])
		}
		agentArgs = args[1:]
	} else {
		var pickErr error
		spec, pickErr = pickAgent(out, os.Stdin, specs)
		if pickErr != nil {
			return pickErr
		}
	}

	return launchAgentFn(ctx, spec, agentArgs)
}

func findSpec(specs []entirecli.AgentSpec, name string) (entirecli.AgentSpec, bool) {
	for _, spec := range specs {
		if strings.EqualFold(name, spec.EntireName) || strings.EqualFold(name, spec.Display) || strings.EqualFold(name, spec.Binary) {
			return spec, true
		}
	}
	return entirecli.AgentSpec{}, false
}

func pickAgent(out io.Writer, in *os.File, specs []entirecli.AgentSpec) (entirecli.AgentSpec, error) {
	if len(specs) == 0 {
		return entirecli.AgentSpec{}, fmt.Errorf("no agents available")
	}
	if len(specs) == 1 {
		fmt.Fprintf(out, "Launching the only enabled agent: %s\n", specs[0].Display)
		return specs[0], nil
	}
	if !stdinIsTerminalFn(in) {
		fmt.Fprintf(out, "Non-interactive stdin; defaulting to %s\n", specs[0].Display)
		return specs[0], nil
	}

	fmt.Fprintln(out, "=== entire-run ===")
	fmt.Fprintln(out)

	spec, err := selectAgentFn(specs)
	if err != nil {
		if errors.Is(err, ui.ErrSelectCanceled) {
			return entirecli.AgentSpec{}, fmt.Errorf("launch canceled")
		}
		return entirecli.AgentSpec{}, err
	}
	return spec, nil
}

func selectAgentInteractive(specs []entirecli.AgentSpec) (entirecli.AgentSpec, error) {
	labels := make([]string, 0, len(specs))
	for _, spec := range specs {
		labels = append(labels, spec.Display)
	}
	idx, err := ui.Select("Launch agent", "Choose an option, press enter to run.", labels)
	if err != nil {
		return entirecli.AgentSpec{}, err
	}
	if idx < 0 || idx >= len(specs) {
		return entirecli.AgentSpec{}, fmt.Errorf("selected agent index %d out of range", idx)
	}
	return specs[idx], nil
}

func launchAgent(ctx context.Context, spec entirecli.AgentSpec, extraArgs []string) error {
	binPath, err := exec.LookPath(spec.Binary)
	if err != nil {
		return fmt.Errorf("%s binary %q not on PATH: %w", spec.Display, spec.Binary, err)
	}

	args := append([]string{}, spec.DefaultArgs...)
	args = append(args, extraArgs...)
	cmd := exec.CommandContext(ctx, binPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	return cmd.Run()
}

func isTerminal(f *os.File) bool {
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

# entire-run

An Entire CLI plugin that launches one of the agents enabled for Entire in the
current repository.

Entire CLI plugins are plain executables named `entire-<name>` on `PATH`.
When a user runs `entire <name>`, the parent CLI dispatches to that binary and
passes the remaining arguments through unchanged.

This plugin builds a binary named `entire-run`, invoked as:

```sh
entire run
```

With no arguments, it shells out to `entire agent list`, shows the agents with
Entire hooks installed for the current repository, and launches the selected
agent in the foreground. If only one agent is enabled, it launches that agent
immediately.

## Quick Start

### Build the Plugin

```sh
mise install
mise run test
mise run build
```

### Install with the CLI

```sh
entire plugin install ./entire-run
```

### Usage

```sh
# Pick an enabled agent from the menu
entire run

# Launch a specific enabled agent
entire run codex
entire run claude-code

# Pass extra arguments to the agent after the agent name
entire run codex -- --model gpt-5
```

### Local Execution

For local development without installing the binary:

```sh
go run ./cmd/entire-run
```

### Subcommands

The `doctor` command expects to run through the Entire CLI so
`ENTIRE_PLUGIN_DATA_DIR` is present. For standalone testing, set it:

```sh
ENTIRE_PLUGIN_DATA_DIR="$(mktemp -d)" go run ./cmd/entire-run doctor
```

## Entire Plugin Contract

The parent CLI supplies these variables when it dispatches a plugin:

| Variable | Meaning |
|---|---|
| `ENTIRE_CLI_VERSION` | Parent CLI version, such as `0.42.0` or `dev`. |
| `ENTIRE_REPO_ROOT` | Absolute git worktree root when invoked inside one. |
| `ENTIRE_PLUGIN_DATA_DIR` | Per-plugin durable storage directory. The plugin should create it before writing. |

The plugin runs in the caller's current working directory. The parent CLI
filters the environment before launching third-party plugins; users can opt
additional variables in with `ENTIRE_PLUGIN_ENV`, for example:

```sh
ENTIRE_PLUGIN_ENV='AWS_*,EDITOR' entire run
```

## Useful Commands

```sh
mise run fmt        # gofmt -s -w .
mise run lint       # go vet, gofmt check, go mod tidy check, shellcheck
mise run test       # go test ./...
mise run test:ci    # go test -race ./...
mise run build      # build ./entire-run
mise run build-all  # cross-build common Entire targets
```

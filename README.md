# Entire Run

An [Entire CLI](https://github.com/entireio/cli) plugin that launches one of the agents enabled for Entire in the
current repository. This way, your sessions are always logged.

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

### Clone the Repo

```sh
git clone https://github.com/entireio/entire-upgrade.git
```

### Build the Plugin

```sh
mise trust
mise install
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

# Launch an agent in its no-approval mode when supported
entire run --yolo codex

# Pass extra arguments to the agent after the agent name
entire run codex --model gpt-5
```

## Development

### Local Execution

For local development without installing the binary:

```sh
go run ./cmd/entire-run
```

### Entire Plugin Contract

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

### Useful Mise Commands

This project uses [mise](https://mise.jdx.dev/) for task automation and dependency management.

```sh
mise run fmt        # gofmt -s -w .
mise run lint       # go vet, gofmt check, go mod tidy check, shellcheck
mise run test       # go test ./...
mise run test:ci    # go test -race ./...
mise run build      # build ./entire-run
mise run build-all  # cross-build common Entire targets
```

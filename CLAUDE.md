# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Analysis-first workflow

Before investigating any non-trivial task (understanding a bug, planning a feature, tracing a data flow), check `.claude/analysis/` for an existing analysis file on that topic. If one exists, read it before exploring the code — do not re-derive what is already recorded.

When analysis is performed, save it to `.claude/analysis/<topic>.md`. The file should capture findings, decisions, and any open questions so the next task can pick up where this one left off without re-reading the same files.

Index every new analysis file in `.claude/analysis/INDEX.md` with a one-line summary so it is discoverable at a glance.

## Commands

```bash
make build          # compile binary (outputs ./qwe)
make install        # build and install to /usr/local/bin/qwe
make dev            # run with DEV_MODE=1 (HTML served from disk, not embedded)
make tidy           # go mod tidy
make all            # bump version + tidy + install
make bumpversion    # increment patch version in VERSION file
make prepare-embed-version  # copy VERSION → internal/web/VERSION (required before build)
```

make test           # go test ./...

# run a single test
go test -run TestExtractJiraTicket ./commit/

To run the tool locally without installing:

```bash
go run . [command]
```

## Architecture

The binary is named `qwe` (not `gocommit`). The entry point `main.go` routes CLI args to subcommands.

**Main use case — interactive commit wizard:**
1. Ensures `gocommit.conf.json` exists in the working directory
2. Runs `git add .`
3. Prompts for commit type (feat/fix/chore/refactor) via `promptui`
4. Extracts a JIRA ticket from the branch name (e.g. `feature/PROJ-123-desc` → `PROJ-123`)
5. Asks for scope and message
6. Executes pre-commit hook (`onBeforeCommit` shell string, if set)
7. Runs `git commit -m "type(scope): message"`
8. Executes post-commit hook (`onAfterCommit`)
9. Optionally pushes if `pushAfterCommit: true` in config and a remote exists

**Packages:**

| Package | Role |
|---|---|
| `commit/` | All direct git operations, config loading, hook execution, JIRA parsing, colored output |
| `internal/wizard/` | Orchestrates the interactive commit flow |
| `internal/branch/` | Interactive branch-switching UI |
| `internal/check/` | Validates/creates `gocommit.conf.json` with defaults |
| `internal/web/` | Background HTTP server + web UI for PRs, branches, repos |
| `srcutils/` | Shell command execution helpers and git state checks |

**Web server lifecycle:** Server info (PID + port) persists in `~/.sensorario-qwe/server.json`. Port auto-detected from 8080 up, or via `QWE_PORT` env var. `web-stop` kills the stored PID.

**Version embedding:** `VERSION` file is the source of truth. `make prepare-embed-version` copies it to `internal/web/VERSION`. Both are embedded at compile time via `//go:embed VERSION` and passed via `-ldflags "-X main.Version=..."`.

**Dev mode:** Set `DEV_MODE=1` (via `make dev`) to have the web server load HTML templates from disk instead of the embedded filesystem — useful when iterating on the web UI.

**Config file** (`gocommit.conf.json`, git-ignored):
```json
{
  "onBeforeCommit": "",
  "onAfterCommit": "",
  "pushAfterCommit": false
}
```
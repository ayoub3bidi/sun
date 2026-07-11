# AGENTS.md

Guidance for coding agents working in this repository.

## Project Overview

`sun` is the Solar Stack CLI. It is a small Go command-line application that launches an interactive Bubble Tea TUI and delegates actual project scaffolding to `github.com/ayoub3bidi/solar-commons`.

Keep this repository focused on:

- CLI commands and help text.
- Interactive TUI flow, validation, defaults, and user-facing errors.
- Template registry metadata.
- Install/run scripts and release packaging.
- Tests that verify `sun` wires the correct inputs into the scaffold engine.

Do not move clone, rewrite, feature pruning, manifest parsing, or git initialization logic into this repo. That belongs in `solar-commons` or in template `solar.manifest.yaml` files.

## Architecture

- `main.go` is only the entry point and calls `cmd.Execute()`.
- `cmd/root.go` defines the Cobra root command, `sun version`, and registers all subcommands.
- `cmd/create.go` defines `sun create` and launches `internal/tui`.
- `cmd/update.go` defines `sun update` and handles self-updating via `go-selfupdate`.
- `cmd/version_check.go` contains `checkForUpdate()` and `runVersionCheck()` for `sun version --check`.
- `cmd/update_check.go` contains background update notification logic (opt-in via `SUN_CHECK_UPDATES=1`).
- `internal/tui/model.go` contains the Bubble Tea model, step transitions, confirmation view, and final scaffold command.
- `internal/tui/validate.go` contains local input validation and friendly error translation.
- `internal/tui/styles.go` contains Lipgloss styles only.
- `registry/registry.go` is the source of truth for available templates, labels, repository URLs, refs, variables, features, and template kind.
- `scripts/run.sh` is the ephemeral npx-style runner.
- `scripts/install.sh` is the persistent installer.
- `scripts/lib/` contains shared shell helper code, but note that `run.sh` and `install.sh` currently inline similar logic.
- `e2e/scaffold_test.go` is behind the `integration` build tag and calls `solar-commons.Scaffold` using registry data.
- `.goreleaser.yaml` builds Linux/macOS amd64/arm64 release archives and injects `cmd.version`, `cmd.commit`, and `cmd.date` via ldflags.

## Local Setup

Use the Go version declared in `go.mod`.

```bash
go build -o sun .
./sun create
```

Current `go.mod` depends on the published `github.com/ayoub3bidi/solar-commons v0.2.0` module. Some docs and CI comments still mention a sibling `../solar-commons` checkout or a `replace` directive; verify the current `go.mod` before following or updating those instructions.

## Validation Commands

Run these before handing off changes:

```bash
go vet ./...
go test ./...
go build ./...
```

Run integration tests only when network access and template repository availability are acceptable:

```bash
go test -tags integration -timeout 180s ./e2e/...
```

The integration test scaffolds from the pinned registry ref and may fail because of remote repository, tag, git, or network issues rather than local unit behavior.

## Implementation Rules

- Preserve the `ScaffoldConfig` -> `solar-commons.Scaffold` path as the single scaffold path.
- Add or update templates in `registry/registry.go`; avoid hard-coding template-specific behavior into TUI step logic.
- Treat `registry/registry.go` as the source of truth for current template refs. Do not copy version numbers into docs or tests unless the test specifically needs a fixed value.
- Keep TUI state transitions explicit and easy to follow. If a new prompt is added, update the step constants, view rendering, update handler, confirm summary, and scaffold config mapping together.
- Keep project-name validation strict: lowercase kebab-case, no path separators, and no existing destination directory.
- Prefer friendly scaffold errors in `internal/tui/validate.go` when translating known `solar-commons` errors.
- Do not introduce non-interactive scaffolding flags unless the CLI contract is intentionally changing; this tool is currently centered on `sun create` as an interactive TUI.
- Keep shell scripts POSIX `sh` compatible.
- If changing release asset names, update `.goreleaser.yaml`, `scripts/run.sh`, `scripts/install.sh`, and `scripts/lib/download.sh` together.
- If changing user-visible workflow, update `README.md`, `CONTRIBUTING.md`, and `e2e/MANUAL.md` where relevant.

## Testing Notes

- Unit tests currently cover registry lookup, local TUI validation helpers, version checking, and update path detection.
- The TUI itself is not keypress-driven in tests; behavior is mostly protected by small helper tests and integration coverage of the scaffold config path.
- `go test ./...` does not run `e2e/scaffold_test.go` because it has `//go:build integration`.
- CI currently checks `go vet ./...` and `go test ./...`; the E2E workflow additionally runs the integration test.

## Known Drift To Watch

- `scripts/lib/download.sh` includes checksum support, but the top-level `run.sh` and `install.sh` currently inline download logic without sourcing that helper.
- Template refs in `registry/registry.go` can move ahead of examples in docs. Check the registry first when updating release notes, docs, or tests.

## Agent Workflow

1. Inspect the relevant files before editing; this repo is small and changes often cross docs, registry, scripts, and tests.
2. Keep edits scoped. Avoid refactoring scaffold engine concerns into `sun`.
3. Run `gofmt` on edited Go files.
4. Run `go vet ./...` and `go test ./...`.
5. For script or release changes, also inspect `.goreleaser.yaml` and manually reason through asset names, version prefixes, OS/arch detection, and cache/install paths.

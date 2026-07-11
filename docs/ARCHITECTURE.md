# Architecture

How `sun` is structured and how data flows through it.

## Project layout

```
sun/
├── main.go                         Entry point — calls cmd.Execute()
├── cmd/
│   ├── root.go                     Cobra root command, sun version, subcommand registration
│   ├── create.go                   sun create — launches the TUI
│   ├── update.go                   sun update — self-update via go-selfupdate
│   ├── update_check.go             Background update notification (opt-in, once per day)
│   ├── version_check.go            sun version --check, GitHub API version comparison
│   └── *_test.go                   Unit tests for version/update logic
├── internal/tui/
│   ├── model.go                    Bubble Tea model, step transitions, confirmation, scaffold orchestration
│   ├── validate.go                 Local input validation and friendly error translation
│   ├── validate_test.go            Validation helper tests
│   └── styles.go                   Lipgloss styles
├── registry/
│   └── registry.go                 Template registry — labels, repos, refs, variables, features
├── e2e/
│   ├── scaffold_test.go            Integration test (//go:build integration)
│   └── MANUAL.md                   Manual E2E test matrix
├── scripts/
│   ├── run.sh                      Ephemeral npx-style runner (primary install path)
│   ├── install.sh                  Persistent installer
│   └── lib/                        Shared shell helpers
├── docs/
│   ├── DEVELOPMENT.md              This file's companion — contributor setup
│   ├── ARCHITECTURE.md             This file
│   └── SOLAR_STACK.md             Ecosystem overview
├── .github/workflows/
│   ├── ci.yml                      Unit tests
│   ├── e2e.yml                     Integration tests
│   └── release.yml                 GoReleaser builds + GitHub release
├── assets/                         Screenshot, diagram
├── AGENTS.md                       Guidance for AI coding agents
├── CONTRIBUTING.md                 Contribution rules and PR checklist
├── CHANGELOG.md                    Release history
└── .goreleaser.yaml                Release build configuration
```

## Data flow

```
User runs "sun create"
        │
        ▼
   ┌─────────┐
   │ Cobra   │  cmd/create.go parses args, calls tui.Run()
   └────┬────┘
        │
        ▼
   ┌─────────┐
   │  TUI    │  internal/tui/model.go — Bubble Tea multi-step form
   │         │  Steps: project type → name/title/version → features → confirm
   └────┬────┘
        │
        ▼
   ┌───────────────┐
   │ ScaffoldConfig│  Collected user inputs mapped to a config struct
   └────┬──────────┘
        │
        ▼
   ┌──────────────┐
   │solar-commons │  github.com/ayoub3bidi/solar-commons
   │              │  1. Shallow clone template @ pinned ref
   │              │  2. Strip .git, parse solar.manifest.yaml
   │              │  3. Apply variable rewrites
   │              │  4. Prune disabled features
   │              │  5. git init + single initial commit
   └──────────────┘
        │
        ▼
   ┌──────────────┐
   │  ./project/  │  Ready-to-use scaffolded project
   └──────────────┘
```

## Design boundaries

| Concern | Lives in |
|---------|----------|
| TUI flow, validation, user-facing errors | `sun` (internal/tui) |
| Template metadata (repos, refs, labels) | `sun` (registry) |
| CLI commands, self-update, version check | `sun` (cmd) |
| Clone, rewrite, prune, manifest parsing | `solar-commons` |
| Template-specific config (variables, features) | Each template's `solar.manifest.yaml` |
| Release builds, archives, changelogs | `.goreleaser.yaml` + `release.yml` |

The CLI stays thin. Scaffold engine logic does not belong in `sun`.

## Update system

`sun update` uses `github.com/creativeprojects/go-selfupdate` to download the latest GitHub release and replace the binary in-place. It detects Snap/Homebrew/Nix installs and redirects users to their package manager instead.

`sun version --check` and the background `SUN_CHECK_UPDATES=1` notification both call the GitHub API (`/repos/ayoub3bidi/sun/releases/latest`) and compare versions using semver.

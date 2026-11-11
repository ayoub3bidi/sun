# Changelog

All notable changes to `sun` will be documented in this file.

## [Unreleased]

## [1.2.0]

### Added

- Non-interactive `sun create` flags for scripting and CI, with JSON success output and template feature controls.
- End-to-end coverage for the direct CLI path across Mars, Mercury, and Venus.

## [1.1.2]

### Changed

- Updated the Mars and Mercury template pins to releases that include the asynchronous agent runtime, run APIs, and their database migrations.
- Expanded scaffold integration coverage to exercise Mars, Mercury, and Venus through the same registry path used by `sun create`.

## [1.1.1]

### Fixed

- `run.sh` and `install.sh` fail when piped via `curl | sh`. Scripts now detect the pipe-to-sh case and re-exec from a downloaded copy so `lib/platform.sh` is available.

### Added

- `sun update` — self-update command that downloads the latest release from GitHub and replaces the binary in-place. Detects Snap, Homebrew, and Nix installs and redirects to the appropriate package manager.
- `sun version --check` — check if a newer version is available without installing it.
- Background update notification via `SUN_CHECK_UPDATES=1` environment variable. Checks once per day and prints a notice to stderr.
- Unit tests for version comparison, update path detection, and TUI validation helpers.
- `AGENTS.md` — guidance for AI coding agents working in this repository.

### Changed

- TUI rewrite: new styles, interactive elements, boolean step handler fixes.
- Updated template refs: Mars v0.1.4, Mercury v0.4.1, Venus v1.1.2.
- Improved `install.sh` and `run.sh` scripts.
- README overhauled with features, commands, and ecosystem docs.

### Fixed

- Boolean step handlers not updating model state in the TUI.
- CodeRabbitAI review comments addressed.

## [1.0.0]

### Added

- Initial release of `sun` CLI.
- Interactive Bubble Tea TUI for scaffolding project templates.
- Mars template (FastAPI + Alembic + PostgreSQL + Redis + JWT auth).
- Mercury template (FastAPI + Flyway + PostgreSQL + Redis + JWT auth).
- Venus template (Docusaurus + Tailwind CSS + i18n + search).
- Ephemeral runner (`run.sh`) — scaffold without installing.
- Persistent installer (`install.sh`).
- `sun version` command.
- GoReleaser-based release builds for Linux/macOS (amd64 + arm64).
- Integration tests against live template repos.

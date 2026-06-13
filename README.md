# Sun CLI

Scaffold production-ready project templates from a single binary.

```bash
curl -fsSL https://raw.githubusercontent.com/ayoub3bidi/sun/main/install.sh | sh
sun create
```

## Overview

`sun` generates customized projects from Solar's template family — Mars, Mercury, and Venus — without requiring Node, Python, or any runtime on your machine. It clones a pinned template, rewrites project names and configuration, prunes disabled features, and initialises a clean git repository — all from one Go binary.

## Quick start

```bash
sun create
```

An interactive TUI will guide you through:
- **Project name** — kebab-case, e.g. `my-api`
- **API version prefix** — e.g. `v1`
- **App title** — used in page titles and documentation
- **OIDC** — optionally include Google OIDC routes and configuration
- **Confirm & scaffold** — one clean git commit in `./<project-name>/`

### Example

```
sun create
  → Project name:       my-api
  → API version prefix:  v1
  → App title:           My API
  → OIDC:                No
  → Output:              ./my-api/

cd my-api
cp .env.dist .env
docker compose up
```

## Templates

| TUI label | Template | Stack | Manifest |
|-----------|----------|-------|----------|
| API + Alembic (recommended) | **Mars** | FastAPI, Alembic, PostgreSQL, Redis, JWT | `solar.manifest.yaml` |
| API + Flyway | **Mercury** | FastAPI, Flyway, PostgreSQL, Redis, JWT | — |
| Docs site | **Venus** | Docusaurus, Tailwind CSS, i18n | — |

Mars is the default and only template wired in the initial release. Mercury and Venus will follow once their manifests are added.

## Commands

| Command | Description |
|---------|-------------|
| `sun` | Print help |
| `sun version` | Print the version number |
| `sun create` | Launch the interactive scaffold TUI |

## Architecture

```
sun CLI (Bubble Tea TUI)
  ↓  ScaffoldConfig
solar-commons (Go library)
  ├─ clone template @ pinned ref
  ├─ strip .git
  ├─ parse solar.manifest.yaml
  ├─ apply variable rewrites
  ├─ prune disabled features
  └─ git init (single commit)
  ↓
./my-project/  (clean, ready to use)
```

The CLI is intentionally thin — all scaffold logic lives in [`solar-commons`](https://github.com/ayoub3bidi/solar-commons). Templates are independent public repositories with their own `solar.manifest.yaml`.

## Development

### Prerequisites

- Go 1.25.3+
- A local checkout of `solar-commons` at `../solar-commons` (relative to `sun/`)

### Setup

```bash
git clone https://github.com/ayoub3bidi/sun
git clone https://github.com/ayoub3bidi/solar-commons
# sibling directories: ./sun/  ./solar-commons/

cd sun
go build -o sun .
./sun create
```

The `go.mod` uses a `replace` directive to point to the local `solar-commons` checkout during development.

### Testing

```bash
go vet ./...
go test ./...
```

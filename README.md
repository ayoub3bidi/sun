<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="./assets/sun.png">
    <img src="./assets/sun.png" width="160" alt="Sun CLI">
  </picture>
  <h3 align="center">Sun CLI</h3>
  <p align="center">Scaffold production-ready backends and documentation sites — from a single binary.</p>
</p>

<p align="center">
  <a href="https://github.com/ayoub3bidi/sun/releases">
    <img src="https://img.shields.io/github/v/release/ayoub3bidi/sun?style=flat&label=Release">
  </a>
  <a href="https://github.com/ayoub3bidi/sun/actions">
    <img src="https://img.shields.io/github/actions/workflow/status/ayoub3bidi/sun/.github/workflows/ci.yml?branch=main&style=flat&label=CI">
  </a>
  <a href="https://pkg.go.dev/github.com/ayoub3bidi/sun">
    <img src="https://img.shields.io/badge/Go-1.25.3-00ADD8?style=flat&logo=go">
  </a>
  <a href="./LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-green?style=flat">
  </a>
</p>

---

## Quick start

```bash
curl -fsSL https://raw.githubusercontent.com/ayoub3bidi/sun/main/install.sh | sh
sun create
```

No Node.js, Python, or Docker required. Just a single binary.

---

## Features

- **Single binary** — no language runtime required on the user's machine
- **Interactive TUI** — Bubble Tea multi-step form with real-time validation
- **Pinned templates** — every scaffold uses a stable semver tag, never `main`
- **Optional features** — toggle OIDC, i18n, search — disabled code is removed entirely
- **Clean git history** — output is a fresh repo with one initial commit, no upstream baggage
- **Manifest-driven** — templates declare their own variables, rewrites, and features in `solar.manifest.yaml`

## How it works

```
┌─────────────┐     ┌──────────────────────────────────────────────────┐
│  sun create  │ ──→ │  Bubble Tea TUI: name, version, title, OIDC…    │
└─────────────┘     └──────────────┬───────────────────────────────────┘
                                   ↓
                            ScaffoldConfig
                                   ↓
┌──────────────────────────────────────────────────────────────────────┐
│  solar-commons (Go library)                                           │
│                                                                      │
│  1. Shallow clone template @ pinned ref (v0.1.3, v1.1.1, …)         │
│  2. Strip .git history                                               │
│  3. Parse solar.manifest.yaml                                        │
│  4. Apply variable rewrites (project name, title, version)           │
│  5. Prune disabled features (OIDC paths, i18n, search)               │
│  6. git init + single initial commit                                 │
└──────────────────────────────────────────────────────────────────────┘
                                   ↓
                    ┌──────────────────────────────┐
                    │  ./my-project/               │
                    │  ├── .env.dist               │
                    │  ├── docker-compose.yml      │
                    │  ├── src/                    │
                    │  └── .git/ (clean, 1 commit) │
                    └──────────────────────────────┘
```

The CLI stays thin. All scaffold logic lives in [`solar-commons`](https://github.com/ayoub3bidi/solar-commons). Templates are independent repositories with their own `solar.manifest.yaml` — no template logic is hard-coded in `sun`.

## Templates

| TUI label | Template | Stack | Manifest |
|---|---|---|---|
| API + Alembic (recommended) | **Mars** | FastAPI, Alembic, PostgreSQL, Redis, JWT auth | ✅ `v0.1.3` |
| API + Flyway | **Mercury** | FastAPI, Flyway, PostgreSQL, Redis, JWT auth | ❌ Coming soon |
| Docs site | **Venus** | Docusaurus 3, Tailwind CSS, i18n | ❌ Coming soon |

Mars is the only template currently wired in the TUI. Mercury and Venus will be added once their manifests are ready.

### Template lifecycle

Each template is a standard public GitHub repository. A `solar.manifest.yaml` at the root declares:

- **`variables`** — what the CLI asks the user (project name, API version, etc.)
- **`rewrites`** — literal string replacements across env files, compose files, source code
- **`features`** — optional components (OIDC, i18n, search) with their removal rules when disabled
- **`default_ref`** — a pinned semver tag for reproducible scaffolds

This means templates evolve independently. A new Mars release only requires bumping `default_ref` — no `sun` release needed.

## Commands

| Command | Description |
|---|---|
| `sun` | Print help and available commands |
| `sun version` | Print the current version |
| `sun create` | Launch the interactive scaffold TUI |

### `sun create` flow

```
Step 1 — Project type
├── API project (Alembic recommended, or Flyway)
└── Docs project            ← not yet wired

Step 2a — If API
│   ├── Project name, API version prefix, app title
│   └── Include Google OIDC? (default: no)
│   → Output: ./<name>/ (customized Mars or Mercury clone)

Step 2b — If Docs
    ├── Project name, site title
    ├── Include i18n? (default: yes)
    └── Include local search? (default: yes)
    → Output: ./<name>/ (customized Venus clone)
```

### Example

```
sun create
  → Project name:       my-api
  → API version prefix: v1
  → App title:          My API
  → OIDC:               No
  → Output:             ./my-api/

cd my-api
cp .env.dist .env
docker compose up
```

## Output structure

A scaffolded API project looks like this:

```
my-api/
├── .env.dist              # renamed, OIDC keys removed
├── docker-compose.yml     # service names renamed
├── Dockerfile
├── alembic/
├── src/
│   ├── main.py            # app title rewritten
│   ├── constants/settings.py
│   ├── routes/
│   ├── controllers/
│   └── utils/
├── .git/                  # single clean commit
├── requirements.txt
└── README.md
```

## Development

### Prerequisites

- Go 1.25.3+
- A local checkout of `solar-commons` as a sibling directory

### Setup

```bash
git clone https://github.com/ayoub3bidi/sun
git clone https://github.com/ayoub3bidi/solar-commons

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

### Project layout

```
sun/
├── main.go                # Entry point
├── cmd/
│   ├── root.go            # Cobra root — sun, sun version
│   └── create.go          # sun create — launches TUI
├── registry/
│   └── registry.go        # Template registry (mars currently)
├── internal/tui/
│   ├── model.go           # Bubble Tea model, steps, scaffold orchestration
│   └── styles.go          # Lipgloss styling
├── .github/
│   ├── workflows/ci.yml
│   └── ISSUE_TEMPLATE/
└── assets/
```

## Roadmap

- **v0.2** — Mercury and Venus manifests, non-interactive flags (`--ci`), `sun doctor`
- **v0.3** — Cache directory (`~/.sun/cache/`), offline mode, template registry JSON
- **v0.4** — `sun upgrade` for template updates, `sun cache clear`
- **v1.0** — GoReleaser, `install.sh`, `brew install`, stable API

See the [full plan](./PLAN.md) for details.

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) for development workflow and guidelines.

## License

[MIT](./LICENSE) — Copyright (c) 2026 Ayoub Bidi

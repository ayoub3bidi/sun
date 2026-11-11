# Development

How to build, test, and contribute to `sun`.

## Prerequisites

- Go 1.25.12+ (see `go.mod` for the exact version)

That's it. `go.mod` depends on the published `github.com/ayoub3bidi/solar-commons` module — no local checkout or `replace` directive is required.

## Local setup

```bash
git clone https://github.com/ayoub3bidi/sun
cd sun
go build -o sun .
./sun create
```

## Running tests

Unit tests:

```bash
go vet ./...
go test ./...
```

Integration tests (requires network access to clone template repos):

```bash
go test -tags integration -timeout 300s ./e2e/...
```

The integration test scaffolds from the pinned registry refs and may fail due to remote repository or network issues rather than local code problems.

## CI workflows

| Workflow | Trigger | What it does |
|----------|---------|-------------|
| `ci.yml` | Push/PR to `main` | `go vet` + `go test` |
| `e2e.yml` | Push/PR to `main` | Unit tests + Mars, Mercury, and Venus scaffold integration tests |
| `release.yml` | Push tag `v*` | Tests + GoReleaser builds Linux/macOS archives + GitHub release |

## Release process

1. Update `CHANGELOG.md` with the new version.
2. Tag and push: `git tag v1.1.0 && git push origin v1.1.0`.
3. The `release.yml` workflow builds binaries and creates the GitHub release.

GoReleaser injects `cmd.version`, `cmd.commit`, and `cmd.date` via ldflags so the binary knows its own version at runtime.

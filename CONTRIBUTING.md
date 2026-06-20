# Contributing

Thanks for helping improve `sun`.

`sun` is the Solar CLI — a Bubble Tea TUI that scaffolds project templates using `solar-commons`. See [docs/SOLAR_STACK.md](./docs/SOLAR_STACK.md) for how this repo fits the ecosystem.

## Local Setup

```bash
git clone https://github.com/ayoub3bidi/sun
git clone https://github.com/ayoub3bidi/solar-commons
# sibling directories: ./sun/  ./solar-commons/

cd sun
go build -o sun .
```

Use the Go version declared in `go.mod`.

```bash
go test ./...
go vet ./...
```

## Contribution Rules

- Keep `sun` focused on user experience: prompts, templates, flags, releases.
- Keep scaffold engine logic in `solar-commons`. Do not duplicate clone, rewrite, prune, or git init here.
- Template-specific branching logic belongs in `solar.manifest.yaml`, not in `sun` code.
- Use the template registry (`registry/`) to map user-facing labels to template IDs.
- Preserve the `ScaffoldConfig` → `solar-commons.Scaffold` pipeline as the single code path.
- New templates are added to the registry, not to the TUI's step logic.

## Pull Requests

Before opening a pull request:

1. Run `go build ./...`.
2. Run `go vet ./...`.
3. Run `go test ./...`.
4. Update `README.md` if commands, templates, or workflows change.

Small, focused pull requests are easier to review.

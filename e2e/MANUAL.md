# Manual E2E matrix

Run on a **fresh Linux** and **fresh macOS** VM (no Go, no prior `sun` install).

| # | Scenario | Pass criteria |
|---|----------|---------------|
| 1 | `curl -fsSL .../scripts/run.sh \| sh` → API → Alembic, OIDC off | `./my-api/` valid; `which sun` fails |
| 2 | Ephemeral → API → Flyway, OIDC on | OIDC routes present; no `mercury` in `.env.dist` |
| 3 | Ephemeral → Docs, i18n on, search on | `pnpm install && pnpm build` succeeds |
| 4 | Ephemeral → Docs, i18n off, search off | No `i18n/fr/`; no search plugin |
| 5 | `install.sh` then `sun create` | Second scaffold from PATH works |
| 6 | Existing destination dir | Friendly error, no partial directory |
| 7 | Unsupported platform | Clear failure message |
| 8 | `sun create --template mars --name my-api --format json` | No prompts; valid Mars scaffold and parseable JSON result |
| 9 | `sun create --template mercury --name my-api --with-oidc` | OIDC files remain in the scaffold |
| 10 | `sun create --template venus --name my-docs --without-i18n --without-search` | i18n and search files are removed |
| 11 | Missing name, unknown template, or template-incompatible feature flag | Non-zero exit, actionable stderr, no scaffold directory |

Pin version when testing a release candidate:

```bash
SUN_VERSION=v1.0.0 curl -fsSL https://raw.githubusercontent.com/ayoub3bidi/sun/main/scripts/run.sh | sh
```

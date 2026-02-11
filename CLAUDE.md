# Bunny Database Shell

## What this is

An NPM package that builds and runs a Go binary to connect to a Bunny Database shell. Bunny Database uses libsql/libsql-server under the hood, so we use [libsql-shell-go](https://github.com/tursodatabase/libsql-shell-go) as the Go dependency for the interactive SQL shell.

## Architecture

- **Go binary** (`cmd/bunny-database-shell/main.go`): Uses Cobra for CLI flag parsing and `libsql-shell-go/pkg/shell` for the interactive shell. Accepts `--url` and `--auth-token` flags. Values are resolved in order: CLI flags -> `.env` file (`BUNNY_DB_URL`, `BUNNY_DB_TOKEN`) -> interactive prompt.
- **NPM package** (`package.json`): Publishes to NPM as `bunny-database-shell`. The `bin` field points to `./bin/cli.js` (a Node wrapper that spawns the Go binary). The `postinstall` script (`scripts/install.js`) downloads the correct prebuilt binary from GitHub Releases.
- **Prebuilt binaries**: CI cross-compiles for linux (x64/arm64), darwin (x64/arm64), and windows (x64). Binaries are attached to GitHub Releases and downloaded at install time. Users do NOT need Go installed.

## Key files

- `cmd/bunny-database-shell/main.go` — CLI entry point (Go source)
- `bin/cli.js` — Node wrapper that npm links to, spawns the Go binary
- `scripts/install.js` — postinstall script, downloads prebuilt binary from GitHub Releases
- `.github/workflows/release.yml` — CI workflow, cross-compiles and creates GitHub Releases
- `package.json` — NPM package config
- `go.mod` / `go.sum` — Go module dependencies

## Dependencies

- Go: `github.com/libsql/libsql-shell-go` (provides `pkg/shell.RunShell`)
- Go: `github.com/spf13/cobra` (CLI framework)
- Go: `github.com/joho/godotenv` (.env file loading)
- The `internal/` packages of libsql-shell-go are NOT importable. We use the public API at `pkg/shell`.

## Bunny Database vs Turso

Bunny Database URLs follow the same format as Turso (`libsql://`, `wss://`) but point to Bunny Database instances. Where libsql-shell-go references Turso, we are connecting to Bunny Database.

## Build & test

```sh
go build -o bin/bunny-database-shell ./cmd/bunny-database-shell/   # local dev build
node bin/cli.js --help                                             # test wrapper
npx bunny-database-shell --url <URL> --auth-token <TOKEN>          # run
```

## Release flow

1. Bump `version` in `package.json`
2. `git commit -m "v{version}" && git tag v{version} && git push --tags`
3. CI builds binaries via GoReleaser and creates GitHub Release, then publishes to NPM

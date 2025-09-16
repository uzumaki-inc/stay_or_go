# Repository Guidelines

## Project Structure & Module Organization
- `cmd/`: Cobra CLI entry (`root.go`), flag parsing and command execution.
- `parser/`: Reads `go.mod`/`Gemfile`, resolves repo URLs.
- `analyzer/`: Fetches GitHub repo stats and computes scores.
- `presenter/`: Renders results (`markdown`, `csv`, `tsv`).
- `utils/`: Shared helpers (logging/verbosity).
- Top level: `main.go` (loads `.env`, runs CLI), `.golangci.yml`, `.air.toml`, `Makefile`, `.sample_files/`.

## Build, Test, and Development Commands
- `go run . go -g $GITHUB_TOKEN`: Analyze Go deps in current repo.
- `go run . ruby -i ./Gemfile -g $GITHUB_TOKEN`: Analyze Ruby deps.
- `GITHUB_TOKEN=... go run . go`: Use env var instead of `-g`.
- `go test ./...`: Run unit tests across packages.
- `make lint` / `make lintFix`: Run `golangci-lint` (and autofix where possible).
- `make air`: Live dev with Air; runs lint + tests before rebuild.
- `make cover`: Run tests with coverage and open `coverage.html`.

## Coding Style & Naming Conventions
- Go standards: `gofmt` formatting, idiomatic naming (`CamelCase` for exported, `camelCase` for unexported, packages lowercase).
- Errors: prefer sentinel vars `Err...` and `fmt.Errorf("...: %w", err)` wrapping.
- Linting: follow `.golangci.yml` rules; some linters are disabled by design (e.g., gofumpt, exhaustruct).
- Output: use `utils.StdErrorPrintln`/`utils.DebugPrintln` for consistent stderr and verbose logs.

## Testing Guidelines
- Frameworks: std `testing` with `testify/assert`; HTTP calls mocked via `jarcoal/httpmock`.
- Placement: co-locate tests as `*_test.go` within each package (see `parser/`, `presenter/`, `analyzer/`).
- Run: `go test ./...` (Air also runs it). Add focused tests for new logic and edge cases.

## Commit & Pull Request Guidelines
- Commits: short, imperative subject with a type prefix seen in history (e.g., `Fix: ...`, `Add: ...`, `Update: ...`, `Rename: ...`, `Upgrade: ...`).
  - Example: `Fix: include watchers in scoring`.
- PRs: include clear description, rationale, and linked issues; add tests for new behavior; include sample output when changing formats.
- Before opening: run `make lint` and `go test ./...` and ensure CLI help/flags remain accurate.

## Security & Configuration Tips
- Auth: set `GITHUB_TOKEN` (env or `-g`). Do not commit `.env` or tokens.
- Rate limits: prefer env var with a personal access token to avoid API throttling during analysis.

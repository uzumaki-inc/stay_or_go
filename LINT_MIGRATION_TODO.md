# Option B Lint Migration Plan

This document outlines concrete, reversible steps to move from Option A (lenient tests) to Option B (strict linting across all code, including tests).

## 1) Re-enable gofmt for tests

- [x] Update `.golangci.yml`: remove `gofmt` from the `_test.go` exclusion list.
- [x] Format code: `gofmt -s -w .`
- [x] Verify: `go run github.com/golangci/golangci-lint/cmd/golangci-lint run` passes.

## 2) Reduce complex tests (cyclomatic/cognitive)

- [x] Split large tests into smaller, focused tests (`presenter/analyzed_libinfo_test.go`).
- [x] Avoid long, chained condition checks; simplified per-assert checks.
- [x] Rename ultra-short variables in tests to self-descriptive names.
- [x] Verify: lints pass, `go test ./...` remains green.

## 3) Fix long lines (lll) in tests

- [x] Break long literals (expected outputs/JSON) into parts (done in `presenter/presenter_test.go`, `parser/go_parser_test.go`).
- [x] Prefer multi-line or split strings for long JSON and tables.
- [x] Remove temporary `//nolint:lll` where feasible.
- [x] Verify: lints pass, behavior unchanged (`golangci-lint run`, `go test ./...`).

## 4) Re-enable goimports/gci for tests

- [x] Update `.golangci.yml`: remove `goimports` and `gci` from the `_test.go` exclusion list.
- [x] Normalize imports in tests (fixed grouping/indent in `presenter/presenter_test.go`, `parser/go_parser_test.go`).
- [x] Verify: `golangci-lint run` passes import checks; `go test ./...` remains green.

## 5) Tighten parallel/testpackage rules

- [ ] Revisit `t.Parallel()` usage: avoid when mutating global state (stdout, env, cwd).
- [ ] For stdout-capturing tests, keep serial; otherwise add `t.Parallel()`.
- [ ] Where possible, move test packages to `*_test` (e.g., `presenter_test`) without breaking access patterns.
- [ ] Replace global env changes with per-test subprocess where safer.
- [ ] Verify: race-safe, lints pass.

## 6) Remove temporary suppressions for `cmd/root.go`

- [x] Update `.golangci.yml`: delete the `cmd/root.go`-specific exclusions (`gci`, `goimports`, `gofmt`, `gofumpt`, `whitespace`, `wsl`, `nlreturn`, `varnamelen`).
- [x] Adjust code style (whitespace) and add targeted function-level nolint for `wsl` to keep readability.
- [x] Keep error wrapping and sentinel errors (err113, wrapcheck) intact.
- [x] Verify: project-wide `golangci-lint run` passes.

## 7) Acceptance criteria (for each step)

- [x] `golangci-lint run` passes locally.
- [x] `go test ./...` passes.
- [x] Coverage stays ≥ 90% for core packages (`analyzer` 91.8%, `parser` 90.4%, `presenter` 97.1%, `cmd` 90.5%).
- [x] No behavior changes: CLI and test outputs remain consistent.

## 8) Rollback plan

- If a step causes friction, temporarily re-add that linter to the `_test.go` exclusion in `.golangci.yml` and open a follow-up task to address root causes.

## 9) Suggested command snippets

- Lint: `go run github.com/golangci/golangci-lint/cmd/golangci-lint run -v`
- Format: `gofmt -s -w .`
- Imports: `goimports -w .` (or configure in your IDE)
- Tests (with local cache): `GOCACHE="$(pwd)/.gocache" go test ./... -v`

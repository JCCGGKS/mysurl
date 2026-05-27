# Repository Guidelines

## Project Structure & Module Organization
This repository is currently a small Go module rooted at `/` with `go.mod`, `README.md`, and `LICENSE`. Treat the repository root as the module root for all Go code. Add application entrypoints under `cmd/<app>/` and keep reusable packages under `internal/` or `pkg/` as the codebase grows. The `vibecodinglearn/` directory is ignored by Git and should not be used for production code or tests.

## Build, Test, and Development Commands
Use standard Go tooling from the repository root:

- `go build ./...` builds all packages and catches compile errors.
- `go test ./...` runs the full test suite.
- `go test -cover ./...` checks package coverage before opening a PR.
- `gofmt -w .` formats Go source files in place.

Because the project is still scaffold-level, some commands may report "no packages" until source files are added.

## Coding Style & Naming Conventions
Follow idiomatic Go. Use tabs for indentation and let `gofmt` own formatting. Keep package names short and lowercase, for example `internal/store` or `pkg/shorturl`. Exported identifiers use `CamelCase`; unexported identifiers use `camelCase`. File names should be lowercase and descriptive, such as `handler.go`, `service_test.go`, or `memory_store.go`.

## Testing Guidelines
Write tests with Go’s built-in `testing` package in `*_test.go` files. Prefer table-driven tests for handlers, parsers, and validation logic. Keep tests close to the code they cover, and name them `Test<FunctionOrBehavior>`. Run `go test ./...` locally before committing; use `go test -cover ./...` when changing core logic.

## Commit & Pull Request Guidelines
Recent history uses short, imperative commit messages, for example `Initial commit` and `docs(.gitignore)`. Keep that style: one focused change per commit, with a concise subject line. For pull requests, include a short summary, test evidence (`go test ./...` output or equivalent), and any relevant issue link. Add request/response examples when API behavior changes.

## Configuration & Repo Hygiene
Do not commit `.env`, generated binaries, coverage artifacts, or local editor files; `.gitignore` already excludes common Go outputs. Keep README updates in the same PR when introducing new commands, directories, or runtime configuration.

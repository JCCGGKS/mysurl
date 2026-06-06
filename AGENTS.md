# Repository Guidelines

## Project Structure & Module Organization

This repository is a Go service for short links. The module root is the repository root.

- `mysurl1.go`: local entrypoint for running the API service
- `internal/config`: configuration structs
- `internal/dao`: MySQL and Redis access
- `internal/handler`: HTTP handlers
- `internal/logic`: business logic
- `internal/model`: table models
- `internal/schema`: request and response types
- `internal/svc`: service wiring and dependencies
- `internal/utils`: shared helpers and background workers
- `etc/`: runtime config, for example `etc/mysurl1-api.yaml`
- `wrk/`: benchmark scripts and result notes
- `internal/template/docs`: design and technical documents

Keep production code under `internal/`. Do not place new application code in `vibecodinglearn/`, `wiki/`, or other note directories.

## Build, Test, and Development Commands

- `go build ./...`: build all packages and catch compile issues
- `go test ./...`: run the full test suite
- `GOCACHE=/tmp/gocache go test ./internal/...`: run internal package tests in restricted environments
- `gofmt -w .`: format Go source files
- `go run mysurl1.go -f etc/mysurl1-api.yaml`: run the service locally
- `bash wrk/run_create.sh` or `bash wrk/run_get.sh`: run local benchmark scripts

## Coding Style & Naming Conventions

Follow idiomatic Go and let `gofmt` control formatting. Use tabs for indentation. Package names stay short and lowercase, such as `dao`, `logic`, and `utils`. Exported names use `CamelCase`; unexported names use `camelCase`.

Prefer clear file names like `redirectlogic.go`, `shortlinkdao.go`, or `visit_flush_worker.go`. If a function returns `error`, callers must handle it immediately.

## Testing Guidelines

Use Go’s built-in `testing` package in `*_test.go` files. Keep tests close to the code they cover. Name tests as `Test<Behavior>` and prefer table-driven tests for utilities, parsing, and strategy logic.

Run `go test ./...` before committing. For focused local checks, `go test ./internal/...` is acceptable.

## Commit & Pull Request Guidelines

Recent history uses concise Conventional Commit style, for example:

- `fix(gitignore)`
- `docs(wrk):新增压测结果报告`

Prefer one focused change per commit. PRs should include:

- a short summary
- test evidence, such as `go test ./...`
- config or API examples when behavior changes
- benchmark context when updating `wrk/` results

## Configuration & Benchmark Notes

Local runtime settings live in `etc/mysurl1-api.yaml`. Redis, MySQL, Bloom, and cache-expire experiments should be documented alongside code changes in `internal/template/docs` or `wrk/`.

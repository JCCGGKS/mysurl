# Repository Guidelines

## Project Structure & Module Organization

This repository is a Go service for short links. The module root is the repository root.

- `mysurl1.go`: local entrypoint for running the API service
- `internal/config`: configuration structs and config item definitions only
- `internal/dao`: MySQL and Redis access only; SQL, cache read/write, and persistence details belong here
- `internal/handler`: HTTP handlers only; parse requests, call logic, write responses, avoid business branching here
- `internal/logic`: business logic orchestration; decide workflow, validation flow, DAO calls, and response assembly
- `internal/model`: table models only; database row shapes and storage-oriented structs belong here
- `internal/schema`: request and response types only; HTTP/API wire structs belong here, not DB models
- `internal/middleware`: cross-cutting HTTP concerns such as auth and operation logging
- `internal/svc`: service wiring and shared dependencies; initialize DB, Redis, DAOs, managers, and workers here
- `internal/utils`: shared helpers and reusable infrastructure utilities; JWT, password hashing, response helpers, pagination, URL helpers, and background workers belong here
- `etc/`: runtime config, for example `etc/mysurl1-api.yaml`
- `wrk/`: benchmark scripts and result notes
- `internal/template/docs`: design and technical documents

Keep production code under `internal/`. Do not place new application code in `vibecodinglearn/`, `wiki/`, or other note directories.

### Layer Responsibilities

- `config` answers "what can be configured"
- `schema` answers "what does the API receive/return"
- `model` answers "what does one DB record look like"
- `dao` answers "how is data queried or persisted"
- `logic` answers "when and why is each operation performed"
- `handler` answers "how does HTTP enter the system"
- `middleware` answers "what concerns apply across many routes"
- `svc` answers "how are dependencies assembled"
- `utils` answers "what generic capability is reused across layers"

### Current Classification Notes

After checking the current `internal/` layout, the main packages are broadly in the right layer:

- auth helpers, response helpers, pagination, MySQL duplicate-key helpers, and visit flush worker are reasonable in `internal/utils`
- DAOs are correctly isolated in `internal/dao`
- request/response structs are correctly isolated in `internal/schema`
- auth/oplog interception remains correctly placed in `internal/middleware`

When adding new files, follow these boundaries:

- do not put SQL or Redis access in `logic`
- do not put HTTP request parsing or response writing in `logic` or `dao`
- do not put database row structs in `schema`
- do not put API response structs in `model`
- if a helper is reused by more than one logic or middleware and is not business-specific, prefer `utils`

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

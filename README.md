# mysurl-v1
短链系统

## Debugging

The service uses `go-zero`, which starts an internal dev server when `DevServer`
is enabled in config. In this repo it is exposed on `0.0.0.0:6060`.

- Health: `http://127.0.0.1:6060/healthz`
- Metrics: `http://127.0.0.1:6060/metrics`
- Pprof index: `http://127.0.0.1:6060/debug/pprof/`

Useful memory profiling commands:

```bash
go tool pprof http://127.0.0.1:6060/debug/pprof/heap
go tool pprof -alloc_space http://127.0.0.1:6060/debug/pprof/heap
go tool pprof -inuse_space http://127.0.0.1:6060/debug/pprof/heap
```

If you only want to suppress periodic `stat` logs such as
`CPU: ..., MEMORY: ...`, call `logx.DisableStat()` during startup.

This repo also supports disabling the go-zero stat background CPU sampler via
environment variable:

```bash
GOZERO_DISABLE_STAT_SAMPLER=1 go run mysurl1.go -f etc/mysurl1-api.yaml
```

# GC 排查记录

## 1. 现象

程序运行后持续出现：

```text
CPU: 4m, MEMORY: Alloc=3.3Mi, TotalAlloc=11.0Mi, Sys=13.3Mi, NumGC=3
```

表现为：
- 没有明显业务流量
- 仍然持续分配内存
- 到达阈值后触发 GC

## 2. 如何发现

先根据日志里的调用位置定位源码：

```text
stat/usage.go:82
```

确认日志来自 `go-zero` 依赖中的 `core/stat/usage.go`，不是业务代码主动打印。

随后通过 `pprof` 抓取两次 `allocs`：
- 第一次现象： [allocs1.md](/home/fanqicheng/project/jx/mysurl1/pprof/allocs1.md:1)
- 关闭日志打印后的第二次现象： [allocs2.md](/home/fanqicheng/project/jx/mysurl1/pprof/allocs2.md:1)

## 3. 如何确定问题

第一次结果 `allocs1` 的热点集中在：
- `ReadTextLines`
- `currentCgroupV2`
- `systemCpuUsage`
- `RefreshCpu`
- `core/stat.init`

关键片段：
- [allocs1.md](/home/fanqicheng/project/jx/mysurl1/pprof/allocs1.md:2)
- [allocs1.md](/home/fanqicheng/project/jx/mysurl1/pprof/allocs1.md:20)
- [allocs1.md](/home/fanqicheng/project/jx/mysurl1/pprof/allocs1.md:60)

为了排除“只是日志输出导致分配”，又关闭了：

```go
logx.DisableStat()
stat.DisableLog()
```

第二次结果 `allocs2` 的热点仍然是同一条链：
- `ReadTextLines -> currentCgroupV2 -> RefreshCpu -> stat.init`

关键片段：
- [allocs2.md](/home/fanqicheng/project/jx/mysurl1/pprof/allocs2.md:2)
- [allocs2.md](/home/fanqicheng/project/jx/mysurl1/pprof/allocs2.md:20)
- [allocs2.md](/home/fanqicheng/project/jx/mysurl1/pprof/allocs2.md:123)

再结合 `go-zero` 源码可以确认：
- `core/stat/usage.go` 在包 `init()` 中启动后台 goroutine
- 每 `250ms` 执行一次 CPU 采样
- 每次都会读取 cgroup 和 `/proc/stat`
- 过程中持续创建短命对象

## 4. 问题是什么

这次问题不是业务内存泄漏，而是 `go-zero/core/stat` 的后台 CPU 采样在持续制造短命对象，进而推动 GC 触发。

更准确地说：
- 不是“对象释放不掉”
- 而是“对象不断被短周期分配出来”

因此：
- `TotalAlloc` 持续上涨是正常现象
- GC 频繁触发是采样分配推动的结果
- 根因在框架后台任务，不在当前业务逻辑

## 5. 如何解决

短期处理：
- 关闭 `stat` 日志，减少观测噪音
- 排查时不要把 `/debug/pprof/trace` 的额外分配混进结论

根因处理：
- 不要在 `core/stat` 的包 `init()` 中直接启动 sampler
- 把采样逻辑改成显式启动，例如 `StartSampler()`
- 在业务 `main()` 中先加载配置，再决定是否启动 sampler

## 6. 结论

- 发现方式：先看日志来源，再用 `pprof allocs` 抓两次结果对比
- 确认方式：`allocs1` 和 `allocs2` 的热点一致，且都指向 `go-zero/core/stat`
- 问题本质：框架后台 CPU 采样导致的短命对象持续分配，不是业务内存泄漏
- 解决方向：把 stat sampler 从包 `init()` 自动启动改为配置加载后的显式启动

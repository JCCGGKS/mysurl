# pprof 使用手册

本文档用于团队内部排查 Go 服务的 CPU、内存和 goroutine 问题，适用于当前仓库。

参考资料：
- Go 官方 `net/http/pprof` 文档：<https://pkg.go.dev/net/http/pprof>
- 博客园文章《万字长文讲解Golang pprof 的使用》：<https://www.cnblogs.com/hobbybear/p/18059425>

## 1. 前置条件

当前项目已经通过 `go-zero` 的 dev server 暴露了 `pprof`。

配置位置：
- [etc/mysurl1-api.yaml](/home/fanqicheng/project/jx/mysurl1/etc/mysurl1-api.yaml:1)

关键配置：

```yaml
DevServer:
  Enabled: true
  Host: 0.0.0.0
  Port: 6060
  EnablePprof: true
```

默认访问地址：
- `http://127.0.0.1:6060/debug/pprof/`

补充说明：
- Go 1.22 及以后，`/debug/pprof/*` 需要使用 `GET` 请求。
- 当前项目使用 `go-zero` dev server，`pprof` handler 已由框架注册。

## 2. 启动服务

正常启动即可：

```bash
go run mysurl1.go -f etc/mysurl1-api.yaml
```

可先验证 `pprof` 是否已启动：

```bash
curl http://127.0.0.1:6060/debug/pprof/
```

## 3. 采样入口

当前常用入口包括：

- `allocs`：历史上所有内存分配的采样
- `block`：阻塞在同步原语上的调用栈
- `cmdline`：当前程序的启动命令行
- `goroutine`：当前所有 goroutine 的调用栈
- `heap`：当前仍然存活对象的内存分配采样
- `mutex`：竞争 mutex 的持有者调用栈
- `profile`：CPU profile
- `symbol`：将程序计数器映射到函数名
- `threadcreate`：创建新 OS 线程的调用栈
- `trace`：程序执行 trace

## 4. 命令行排查

### 4.1 堆内存

```bash
go tool pprof http://127.0.0.1:6060/debug/pprof/heap
go tool pprof -top -inuse_space http://127.0.0.1:6060/debug/pprof/heap
go tool pprof -top -alloc_space http://127.0.0.1:6060/debug/pprof/heap
```

说明：
- `inuse_space` 看当前仍然存活的内存。
- `alloc_space` 看累计分配过的内存。
- 排查泄漏优先看 `inuse_space`，不要只看 `alloc_space`。

如需先做一次 GC 再看：

```bash
go tool pprof http://127.0.0.1:6060/debug/pprof/heap?gc=1
```

### 4.2 CPU

```bash
go tool pprof http://127.0.0.1:6060/debug/pprof/profile?seconds=30
go tool pprof -top http://127.0.0.1:6060/debug/pprof/profile?seconds=30
```

### 4.3 goroutine

```bash
go tool pprof -top http://127.0.0.1:6060/debug/pprof/goroutine
```

### 4.4 常用交互命令

```text
top
top10
list <函数名>
tree
peek <关键字>
quit
```

含义：
- `top`：按热点排序显示当前 profile 的主要函数。
- `top10`：只显示前 10 个热点函数。
- `list <函数名>`：查看指定函数对应的源码位置和样本分布。
- `tree`：按调用树方式展示样本分布。
- `peek <关键字>`：按函数名或关键字搜索匹配项。
- `quit`：退出当前 pprof 交互界面。

### 4.5 推荐命令模板

抓当前堆内存热点：

```bash
go tool pprof -top -inuse_space http://127.0.0.1:6060/debug/pprof/heap
```

抓累计分配热点：

```bash
go tool pprof -top -alloc_space http://127.0.0.1:6060/debug/pprof/heap
```

抓 30 秒 CPU：

```bash
go tool pprof -top http://127.0.0.1:6060/debug/pprof/profile?seconds=30
```

保存结果到文件：

```bash
go tool pprof -top -inuse_space http://127.0.0.1:6060/debug/pprof/heap > pprof/heap_inuse_top.txt
go tool pprof -top -alloc_space http://127.0.0.1:6060/debug/pprof/heap > pprof/heap_alloc_top.txt
go tool pprof -top http://127.0.0.1:6060/debug/pprof/profile?seconds=30 > pprof/cpu_top.txt
```

## 5. 页面排查

### 5.1 直接访问页面

浏览器打开：

- `http://127.0.0.1:6060/debug/pprof/`
- `http://127.0.0.1:6060/debug/pprof/heap`
- `http://127.0.0.1:6060/debug/pprof/profile?seconds=30`
- `http://127.0.0.1:6060/debug/pprof/goroutine`

### 5.2 启动本地 Web UI

```bash
go tool pprof -http=:8082 http://127.0.0.1:6060/debug/pprof/heap
go tool pprof -http=:8082 http://127.0.0.1:6060/debug/pprof/profile?seconds=30
go tool pprof -http=:8082 http://127.0.0.1:6060/debug/pprof/goroutine
```

如果已经有本地 profile 文件：

```bash
go tool pprof -http=:8082 ./profile/heap.pprof
go tool pprof -http=:8082 ./profile/cpu.pprof
```

### 5.3 页面里重点看什么

- Flame Graph：看热点路径
- Graph：看调用关系
- Source：结合源码定位
- Peek：按函数名或关键字搜索

如果本机装了 Graphviz，交互命令里也可以直接执行：

```text
web
```

## 6. trace 文件排查

`trace` 适合排查：
- 请求偶发卡顿
- goroutine 阻塞
- 锁竞争
- syscall / 网络等待
- 调度延迟

### 6.1 抓取 trace

```bash
curl -o trace.out http://127.0.0.1:6060/debug/pprof/trace?seconds=5
```

### 6.2 打开 trace

```bash
go tool trace trace.out
```

### 6.3 trace 里重点看什么

- Goroutine analysis：看 goroutine 的运行、阻塞、唤醒
- Synchronization blocking profile：看锁、channel 等同步阻塞
- Syscall blocking profile：看 syscall 阻塞
- Network blocking profile：看网络等待
- Scheduler latency：看调度延迟

### 6.4 注意事项

- `trace` 开销比 `heap` / `profile` 更大，只在必要时抓
- 抓取时长通常 `5s` 到 `10s` 足够
- `trace` 本身会引入额外分配，不要把它混入普通内存泄漏判断

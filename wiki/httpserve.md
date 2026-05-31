# HTTP Serve 与优雅关闭

本文只回答两个问题：

- 当前项目的 HTTP 服务是怎么启动的
- `go-zero` 是怎么做优雅关闭的

## 1. 当前项目的启动方式

入口在 [mysurl1.go](/home/fanqicheng/project/jx/mysurl1/mysurl1.go:20)：

```go
server := rest.MustNewServer(c.RestConf)
defer server.Stop()
server.Start()
```

结论：

- `server.Start()` 是同步阻塞启动
- `defer server.Stop()` 不是主动关闭 HTTP 的主逻辑
- 真正的优雅关闭由 `go-zero` 内部完成

## 2. 为什么这里没有手写 goroutine 跑 HTTP

标准库自写优雅关闭时，常见写法是：

```go
go srv.ListenAndServe()

<-sig
srv.Shutdown(ctx)
```

参考：

- <https://blog.csdn.net/weixin_46290302/article/details/155200359>

这种写法的原因是：

- `ListenAndServe()` 会阻塞
- 主 goroutine 还要继续监听退出信号
- 所以要把 HTTP 服务放进 goroutine

但在 `go-zero` 里，这层编排已经被框架封装了：

- `main` goroutine 直接阻塞在 `server.Start()`
- `go-zero` 内部自己监听 `SIGTERM` / `SIGINT`
- 收到信号后，框架内部执行 shutdown listener

所以不是“优雅关闭不需要 goroutine”，而是“goroutine 已经在框架内部，不需要在 `main` 里再写一层”。

## 3. `go-zero` 的优雅关闭链路

### 3.1 信号监听

`go-zero` 在 `core/proc/signals.go` 里监听：

- `SIGTERM`
- `SIGINT`

收到信号后会调用 `gracefulStop(...)`。

### 3.2 关闭编排

`gracefulStop(...)` 在 `core/proc/shutdown.go` 里执行两类 listener：

- `wrapUpListeners`
- `shutdownListeners`

默认时间：

- `WrapUpTime = 1s`
- `WaitTime = 5.5s`

超过总等待时间还没退出，框架会强制结束进程。

### 3.3 HTTP 服务真正关闭的位置

REST 服务在 `rest/internal/starter.go` 里注册了 shutdown listener：

```go
waitForCalled := proc.AddShutdownListener(func() {
    healthManager.MarkNotReady()
    if e := server.Shutdown(context.Background()); e != nil {
        logx.Error(e)
    }
})
```

关键动作：

1. 先把服务标记为 `not ready`
2. 再调用 `http.Server.Shutdown(context.Background())`

`Shutdown()` 的作用是：

- 停止接收新连接
- 等待已有请求处理完
- 然后返回

## 4. `server.Stop()` 和 `Shutdown()` 的区别

- `http.Server.Shutdown(...)`：真正执行 HTTP 优雅关闭
- `server.Stop()`：`Start()` 返回后的补充清理，不是主关闭入口

## 5. 整体时序

1. `main` 调用 `server.Start()`
2. 内部进入 `ListenAndServe()`，主 goroutine 阻塞
3. `go-zero` 内部 goroutine 监听退出信号
4. 收到 `SIGTERM` / `SIGINT`
5. 框架执行 `gracefulStop()`
6. REST 服务的 shutdown listener 被触发
7. 调用 `http.Server.Shutdown(...)`
8. HTTP 服务排空后返回
9. `server.Start()` 返回
10. 执行 `defer server.Stop()`

## 6. 结论

- 当前项目的 HTTP 服务是同步阻塞启动
- `main` 里不需要手写 `go ListenAndServe()`
- 真正的优雅关闭由 `go-zero` 内部信号监听和 shutdown listener 完成
- HTTP 关闭的核心动作是 `http.Server.Shutdown(...)`

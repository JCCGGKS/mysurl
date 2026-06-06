# singleflight 源码梳理

## 1. 作用

官方实现位于：

- `golang.org/x/sync/singleflight`

它只解决一件事：

- 同一个 `key` 的多个并发请求
- 只允许一个请求执行 `fn`
- 其它请求等待并复用结果

它不是缓存，也不是分布式锁。

## 2. 核心结构

`Group` 维护一个 `key -> call` 的映射：

```go
type Group struct {
	mu sync.Mutex
	m  map[string]*call
}
```

`call` 表示某个 key 当前这一轮执行：

```go
type call struct {
	wg sync.WaitGroup

	val interface{}
	err error

	dups  int
	chans []chan<- Result
}
```

几个关键字段：

- `wg`：让重复请求等待首个请求完成
- `val` / `err`：保存本轮执行结果
- `dups`：记录重复请求数量
- `chans`：给 `DoChan` 的等待方发结果

## 3. Do

签名：

```go
func (g *Group) Do(key string, fn func() (any, error)) (v any, err error, shared bool)
```

关键实现：

```go
func (g *Group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		c.dups++
		g.mu.Unlock()
		c.wg.Wait()

		if e, ok := c.err.(*panicError); ok {
			panic(e)
		} else if c.err == errGoexit {
			runtime.Goexit()
		}
		return c.val, c.err, true
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	g.doCall(c, key, fn)
	return c.val, c.err, c.dups > 0
}
```

含义很直接：

- 如果 `g.m[key]` 已存在，说明已有请求在执行，当前请求只等待并复用结果
- 如果不存在，当前请求创建 `call` 并负责执行 `fn`

`shared` 的语义是：

- `true`：这份结果被多个请求共享
- 不是“当前请求是不是执行者”

## 4. DoChan 和 Forget

`DoChan` 把结果改成通过 channel 返回：

```go
func (g *Group) DoChan(key string, fn func() (interface{}, error)) <-chan Result {
	ch := make(chan Result, 1)
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		c.dups++
		c.chans = append(c.chans, ch)
		g.mu.Unlock()
		return ch
	}
	c := &call{chans: []chan<- Result{ch}}
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	go g.doCall(c, key, fn)
	return ch
}
```

`Forget` 只是把 key 从 map 删除：

```go
func (g *Group) Forget(key string) {
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()
}
```

注意：

- `Forget` 不是取消旧请求
- 它只是允许后续请求不再等待旧请求

## 5. doCall

真正执行 `fn` 的逻辑在 `doCall`：

```go
func (g *Group) doCall(c *call, key string, fn func() (interface{}, error)) {
	normalReturn := false
	recovered := false

	defer func() {
		if !normalReturn && !recovered {
			c.err = errGoexit
		}

		g.mu.Lock()
		defer g.mu.Unlock()
		c.wg.Done()
		if g.m[key] == c {
			delete(g.m, key)
		}

		if e, ok := c.err.(*panicError); ok {
			if len(c.chans) > 0 {
				go panic(e)
				select {} // 阻塞：保证panic顺利执行，且不再执行其余代码
			} else {
				panic(e)
			}
		} else if c.err == errGoexit {
		} else {
			for _, ch := range c.chans {
				ch <- Result{c.val, c.err, c.dups > 0}
			}
		}
	}()

	func() {
		defer func() {
			if !normalReturn {
				if r := recover(); r != nil {
					c.err = newPanicError(r)
				}
			}
		}()

		c.val, c.err = fn()
		normalReturn = true
	}()

	if !normalReturn {
		recovered = true
	}
}
```

这里要看三点：

- `c.wg.Done()` 放在 defer 里，避免等待方卡死
- `if g.m[key] == c { delete(g.m, key) }` 避免 `Forget` 后误删新请求
- panic / `runtime.Goexit` 都被显式处理，等待方语义和首个请求保持一致

## 6. go-zero 实现

本项目实际使用的是：

- `github.com/zeromicro/go-zero/core/syncx`

接口定义：

```go
type SingleFlight interface {
	Do(key string, fn func() (any, error)) (any, error)
	DoEx(key string, fn func() (any, error)) (any, bool, error)
}
```

核心结构更简单：

```go
type call struct {
	wg  sync.WaitGroup
	val any
	err error
}

type flightGroup struct {
	calls map[string]*call
	lock  sync.Mutex
}
```

`DoEx` 的关键实现：

```go
func (g *flightGroup) DoEx(key string, fn func() (any, error)) (val any, fresh bool, err error) {
	c, done := g.createCall(key)
	if done {
		return c.val, false, c.err
	}

	g.makeCall(c, key, fn)
	return c.val, true, c.err
}
```

`createCall` 和 `makeCall`：

```go
func (g *flightGroup) createCall(key string) (c *call, done bool) {
	g.lock.Lock()
	if c, ok := g.calls[key]; ok {
		g.lock.Unlock()
		c.wg.Wait()
		return c, true
	}

	c = new(call)
	c.wg.Add(1)
	g.calls[key] = c
	g.lock.Unlock()

	return c, false
}

func (g *flightGroup) makeCall(c *call, key string, fn func() (any, error)) {
	defer func() {
		g.lock.Lock()
		delete(g.calls, key)
		g.lock.Unlock()
		c.wg.Done()
	}()

	c.val, c.err = fn()
}
```

和官方版相比，go-zero 主要差异有三点：

- 没有 `DoChan`
- 没有 `Forget`
- 没有官方版对 panic / `runtime.Goexit` 的专门处理

它返回的 `fresh` 语义也不同：

- `fresh=true`：当前请求是实际执行 `fn` 的请求
- `fresh=false`：当前请求是等待并复用结果的请求

这比官方的 `shared` 更适合业务日志判断“我是不是首个执行者”。

## 7. 结论

官方 `singleflight` 很克制：

- 只解决“同 key 并发请求只执行一次”
- 不缓存历史结果
- 调用结束后立即删除 key
- 只适合放在缓存未命中后的回源阶段

对当前项目来说，理解这几点就够了：

- `redirect:{short_code}` 和 `create:{normalized_long_url}` 只是并发控制 key
- `singleflight` 不替代 Redis 缓存
- 它只能解决单进程内的重复调用合并

# wrk 压测工具

## 1. 简介

`wrk` 是一个高性能 HTTP 压测工具，常用于快速评估接口的：

- 吞吐量
- 延迟
- 并发连接表现
- 缓存命中 / 未命中差异

它的特点是：

- 基于多线程和事件循环
- 单机可以打出较高并发
- 支持 Lua 脚本自定义请求

适合当前项目这种 HTTP 短链服务。

## 2. 常用参数

最常用的几个参数：

- `-t`：线程数
- `-c`：连接数
- `-d`：压测时长
- `-H`：自定义 Header
- `-s`：Lua 脚本

示例：

```bash
wrk -t4 -c200 -d30s http://127.0.0.1:8888/abc123
```

含义：

- 4 个线程
- 200 个连接
- 持续压测 30 秒

## 3. 输出怎么看

`wrk` 常见输出里重点看这些：

- `Requests/sec`
  每秒请求数
- `Latency`
  平均延迟、抖动、最大值
- `Transfer/sec`
  每秒传输量

如果启用了分位数统计插件或脚本，还会关注：

- P50
- P90
- P99

对短链项目来说，核心通常看：

- QPS 是否稳定
- 平均延迟是否可接受
- 高分位延迟是否明显抖动

## 4. 当前项目压测对象

当前项目最适合压的两个接口：

### 4.1 创建短链

- `POST /api/v1/links`

作用：

- 测创建链路吞吐
- 测 Bloom / 缓存 / MySQL / singleflight 的影响

### 4.2 跳转短链

- `GET /{short_code}`

作用：

- 测跳转链路 QPS
- 测缓存命中场景
- 测缓存失效后的回源与 `singleflight` 合并效果

## 5. 压测跳转接口

最简单的命令：

```bash
wrk -t4 -c200 -d30s http://127.0.0.1:8888/2sfCSPGW2gE
```

这个场景适合验证：

- Redis short->long 缓存命中能力
- `visit_count` Redis 增量写入带来的额外开销

如果要测缓存未命中后的回源效果，通常需要：

- 先清理对应缓存
- 再执行压测

这样可以观察：

- 首次回源时的延迟抖动
- `singleflight` 是否有效减少重复查库

## 6. 压测创建接口

`wrk` 直接压 `POST` 一般要配合 Lua 脚本。

示例脚本：

```lua
wrk.method = "POST"
wrk.body   = '{"long_url":"https://example.com/article/1"}'
wrk.headers["Content-Type"] = "application/json"
```

执行方式：

```bash
wrk -t4 -c100 -d30s -s post_create.lua http://127.0.0.1:8888/api/v1/links
```

这个场景适合验证：

- 相同长链复用路径
- `create:{normalized_long_url}` singleflight 是否合并重复创建
- MySQL 发号 / 写库路径的吞吐表现

## 7. 用 Lua 模拟不同数据

如果你不想一直压同一个长链，可以用 Lua 动态生成请求体。

例如：

```lua
wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"

local counter = 0

request = function()
  counter = counter + 1
  local body = string.format(
    '{"long_url":"https://example.com/article/%d"}',
    counter
  )
  return wrk.format(nil, "/api/v1/links", nil, body)
end
```

这个脚本更适合验证：

- 真正的新建压力
- 不同长链下 MySQL / Redis / Bloom 的负载表现

## 8. 压测建议

### 8.1 先压跳转，再压创建

跳转接口通常是短链系统的核心热点路径，建议先测：

- 缓存命中跳转
- 缓存未命中跳转

再测创建接口。

### 8.2 分开测命中与未命中

不要把这些场景混在一起测：

- 缓存命中
- 缓存未命中
- 相同长链重复创建
- 不同长链持续新建

否则结果很难解释。

### 8.3 结合日志和指标看结果

只看 `wrk` 输出不够，建议同时结合：

- 服务日志
- Redis 命中情况
- MySQL QPS
- pprof / metrics

特别是当前项目已经有：

- Redis 精确缓存
- Bloom 优化
- `singleflight`
- `visit_count` 异步回刷

这些优化是否真正生效，通常要结合日志和依赖侧指标判断。

## 9. 当前项目建议关注点

对这个短链项目，压测时最值得观察的是：

### 9.1 跳转链路

- `shortlink:code:{short_code}` 命中后 QPS 能到多少
- 缓存未命中时 `singleflight` 是否抑制了重复查库
- `visit_count` Redis `INCR` 是否带来明显额外延迟

### 9.2 创建链路

- 相同长链并发创建时是否能稳定复用
- `create:{normalized_long_url}` 是否减少重复发号 / 重复新建
- Bloom 是否有效减少明显不存在数据的查库次数

## 10. 总结

`wrk` 适合当前项目做两类验证：

- 跳转链路热点压测
- 创建链路并发压测

它的价值不只是看 QPS，更重要的是帮助验证：

- 缓存是否生效
- `singleflight` 是否减少重复工作
- 异步 `visit_count` 是否降低了主链路写库压力

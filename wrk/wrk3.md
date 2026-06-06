# v3 singleflight 

## `/:code`接口压测结果 
为方便测试，将key的过期时间设置为1s

### snowflake-有singleflight-同步incr
Running 30s test @ http://127.0.0.1:8888
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   197.04ms   70.41ms 812.39ms   82.67%
    Req/Sec   127.83     23.92   210.00     72.67%
  15299 requests in 30.04s, 4.23MB read
Requests/sec:    509.30
Transfer/sec:    144.24KB
### snowflake-无singleflight-同步incr
Running 30s test @ http://127.0.0.1:8888
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   202.70ms   76.05ms   1.17s    82.30%
    Req/Sec   124.53     25.73   191.00     72.33%
  14906 requests in 30.05s, 4.12MB read
Requests/sec:    496.08
Transfer/sec:    140.49KB
### snowflake-有singleflight-异步incr
Running 30s test @ http://127.0.0.1:8888
  4 threads and 100 connections

  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     8.00ms    4.51ms  54.71ms   73.16%
    Req/Sec     3.22k   422.34     4.40k    70.67%
  385453 requests in 30.08s, 106.60MB read
Requests/sec:  12814.38
Transfer/sec:      3.54MB
### snowflake-无singleflight-异步incr
Running 30s test @ http://127.0.0.1:8888
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     8.67ms    5.38ms  87.37ms   77.64%
    Req/Sec     3.01k   469.05     4.07k    70.98%
  359952 requests in 30.09s, 99.55MB read
Requests/sec:  11962.53
Transfer/sec:      3.31MB

## 结果分析

### 1. `singleflight` 的收益受 `visit_count` 更新方式影响

同步 `incr`：

- 有 `singleflight`：`509.30 req/s`，`197.04ms`
- 无 `singleflight`：`496.08 req/s`，`202.70ms`

这一组提升很小，说明在同步更新 `visit_count` 的情况下，主链路主要还是被 MySQL 写入拖慢，`singleflight` 合并回源的收益会被掩盖。

异步 `incr`：

- 有 `singleflight`：`12814.38 req/s`，`8.00ms`
- 无 `singleflight`：`11962.53 req/s`，`8.67ms`

这一组提升就更明显，说明当 `visit_count` 不再阻塞主链路后，`singleflight` 才能更充分发挥作用。

结论：

- `singleflight` 不是主链路固定成本优化，而是并发回源优化
- 当同步 `incr` 仍在主链路中时，它的整体收益有限
- 当 `visit_count` 改为异步后，它的价值会明显放大

### 2. `singleflight` 主要优化的是缓存过期后的并发击穿

本轮测试将缓存过期时间设置为 `1s`，意味着：

- 缓存会频繁失效
- 同一个 `short_code` 在高并发下容易同时触发多次回源

这正是 `singleflight` 的作用场景：

- 只允许一个请求执行真实回源
- 其它并发请求复用本次结果
- 减少重复查库和重复回填缓存

因此 `singleflight` 的收益主要体现在：

- 缓存未命中
- 且同一 key 存在并发访问

如果请求已经命中缓存，`singleflight` 本身并不会带来收益。

### 3. `singleflight` 更适合搭配“缓存 + 异步 incr”

从这组数据看，较合理的组合是：

- 跳转结果走缓存
- `visit_count` 走 Redis `INCR` 聚合后异步回刷
- 缓存过期时用 `singleflight` 合并并发回源

这种组合下：

- 主链路避免同步 MySQL 写入
- 缓存失效时又能减少数据库瞬时压力

## 总结

- `singleflight` 的核心价值是抑制缓存失效时的并发击穿
- 它更适合与缓存和异步 `visit_count` 一起使用
- 如果主链路仍保留同步 MySQL `incr`，`singleflight` 的收益会被明显掩盖

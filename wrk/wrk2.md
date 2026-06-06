# v2 有无缓存性能对比

## `/:code` 接口压测结果

### snowflake-无缓存-有 incr
running: SHORT_CODE=2snRqpmjqNi wrk -t4 -c100 -d30s -s get_code.lua http://127.0.0.1:8888
Running 30s test @ http://127.0.0.1:8888
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   211.86ms   60.27ms 627.66ms   74.06%
    Req/Sec   118.19     25.86   200.00     66.58%
  14153 requests in 30.06s, 3.91MB read
Requests/sec:    470.82
Transfer/sec:    133.34KB

### snowflake-有缓存-有 incr
running: SHORT_CODE=2snRqpmjqNi wrk -t4 -c100 -d30s -s get_code.lua http://127.0.0.1:8888
Running 30s test @ http://127.0.0.1:8888
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   201.78ms   76.60ms 855.02ms   81.28%
    Req/Sec   125.09     25.75   202.00     71.58%
  14973 requests in 30.04s, 4.14MB read
Requests/sec:    498.41
Transfer/sec:    141.15KB

### snowflake-无缓存-无 incr
running: SHORT_CODE=2snRqpmjqNi wrk -t4 -c100 -d30s -s get_code.lua http://127.0.0.1:8888
Running 30s test @ http://127.0.0.1:8888
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    14.02ms   10.10ms 116.70ms   73.61%
    Req/Sec     1.92k   353.93     2.81k    68.33%
  230189 requests in 30.08s, 63.66MB read
Requests/sec:   7651.69
Transfer/sec:      2.12MB

### snowflake-有缓存-无 incr
running: SHORT_CODE=2snRqpmjqNi wrk -t4 -c100 -d30s -s get_code.lua http://127.0.0.1:8888
Running 30s test @ http://127.0.0.1:8888
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     8.62ms    6.38ms  89.64ms   74.46%
    Req/Sec     3.13k   662.21     6.29k    65.47%
  374797 requests in 30.10s, 103.66MB read
Requests/sec:  12451.98
Transfer/sec:      3.44MB

### snowflake-有缓存-无布隆过滤器-访问大量不存在的 code
Running 30s test @ http://127.0.0.1:8888
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    16.18ms    9.37ms 108.71ms   71.65%
    Req/Sec     1.59k   233.62     4.85k    74.19%
  190664 requests in 30.09s, 45.09MB read
  Non-2xx or 3xx responses: 190664
Requests/sec:   6336.36
Transfer/sec:      1.50MB

### snowflake-有缓存-有布隆过滤器-访问大量不存在 code
Running 30s test @ http://127.0.0.1:8888
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     6.16ms    4.05ms  65.47ms   71.10%
    Req/Sec     4.21k   516.68     5.45k    75.50%
  503950 requests in 30.07s, 119.19MB read
  Non-2xx or 3xx responses: 503950
Requests/sec:  16761.55
Transfer/sec:      3.96MB

### snowflake-有缓存-有布隆过滤器-访问存在 code
Running 30s test @ http://127.0.0.1:8888
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    15.73ms   43.67ms 627.48ms   97.32%
    Req/Sec     2.62k   455.28     3.73k    78.71%
  305688 requests in 30.20s, 84.54MB read
Requests/sec:  10121.30
Transfer/sec:      2.80MB

## 结果分析

### 1. `visit_count` 更新方式先决定性能上限

有 `incr` 时：

- 无缓存：`470.82 req/s`，`211.86ms`
- 有缓存：`498.41 req/s`，`201.78ms`

这一组提升很小，说明当时主链路的主要瓶颈不是 `short_code -> original_url` 查询，而是访问次数更新。

无 `incr` 时：

- 无缓存：`7651.69 req/s`，`14.02ms`
- 有缓存：`12451.98 req/s`，`8.62ms`

这一组缓存收益明显，说明当访问次数更新不再成为主要成本时，缓存才能真正减少回源 MySQL 带来的开销。

结论：

- `redirect` 接口性能首先受 `visit_count` 更新方式影响
- 当同步更新访问次数存在时，缓存收益会被显著掩盖
- 当访问次数更新成本下降后，缓存价值才会充分体现

### 2. 缓存在存在 code 场景下收益明确

对已存在 `code` 的请求，缓存可以减少一次 MySQL 查询，这在“无 `incr`”场景下已经体现得很清楚：

- 吞吐从 `7651.69 req/s` 提升到 `12451.98 req/s`
- 平均延迟从 `14.02ms` 降到 `8.62ms`

说明：

- 对正常跳转流量，缓存是有效优化
- 只要主链路里没有更重的固定成本，缓存命中能明显缩短请求路径

### 3. Bloom 适合拦截大量不存在的 code

访问大量不存在 `code`：

- 无 Bloom：`6336.36 req/s`，`16.18ms`
- 有 Bloom：`16761.55 req/s`，`6.16ms`

这一组提升非常明显，说明 Bloom 对“不存在短码请求”是高收益优化项。

原因是：

- 无 Bloom 时，不存在请求仍要继续走缓存、MySQL 或后续不存在判断
- 有 Bloom 时，可以在更前面直接返回不存在，减少无效查询

结论：

- Bloom 在 `redirect` 链路中的主要价值，是拦截大量不存在的 `short_code`
- 这类场景下，Bloom 的收益非常直接

### 4. Bloom 不适合放在高命中主链路最前面

访问存在 `code`：

- 有缓存、无 Bloom、无 `incr`：`12451.98 req/s`，`8.62ms`
- 有缓存、有 Bloom：`10121.30 req/s`，`15.73ms`

这说明在 `code` 已存在且缓存可命中的情况下，Bloom 会增加一次额外判断，反而拉长请求链路。

原因是：

- 本来请求可以直接命中缓存
- 加入 Bloom 后，需要先做一次 `BF.EXISTS`
- 这一步本身就是额外成本

结论：

- Bloom 不是对所有跳转请求都正收益
- 对已存在且高命中的 `code`，它会拖慢正常请求
- 因此 Bloom 更适合作为“不存在请求优化”，而不是默认放在所有高命中请求前面

## 总结

- `redirect` 接口性能先看 `visit_count` 更新方式，再看缓存与 Bloom
- 缓存在存在 `code` 场景下收益明确，但前提是访问次数更新不能成为主瓶颈
- Bloom 对大量不存在 `code` 的请求收益很高
- Bloom 对已存在且高命中的 `code` 会增加额外判断成本，不应简单视为通用优化

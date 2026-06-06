# v1 short_code 生成方式压测

## 1. 目标

比较三种 `short_code` 生成方式在创建接口上的性能：

- `mysql_auto_increment`
- `redis_incr`
- `snowflake`

压测接口：

- `POST /api/v1/links`

当前版本已注释相同长链复用逻辑和最后的落库逻辑，且固定同一个 `long_url`。

## 2. wrk 脚本

`wrk/post_create.lua`

```lua
wrk.method = "POST"
wrk.body   = '{"long_url":"https://example.com/test"}'
wrk.headers["Content-Type"] = "application/json"
```

## 3. 固定命令

三轮都使用同一条命令：

```bash
wrk -t4 -c100 -d30s -s wrk/post_create.lua http://127.0.0.1:8888/api/v1/links
```

## 4. 每轮步骤

1. 修改 `etc/mysurl1-api.yaml` 中的 `Short.Provider`
2. 清空 `short_links` 表
3. 清空 Redis 相关 key
4. 重启服务
5. 手工 `curl` 验证创建接口可用
6. 执行 `wrk`
7. 记录结果

## 5. 三轮配置

### mysql_auto_increment

```yaml
Short:
  Provider: mysql_auto_increment
```

### redis_incr

```yaml
Short:
  Provider: redis_incr
```

### snowflake

```yaml
Short:
  Provider: snowflake
  Snowflake:
    WorkerID: 1
```

## 6. 记录项

至少记录：

- Provider
- Requests/sec
- Avg latency
- Max latency
- Errors

建议表格：

| Provider | Req/s | Avg Latency | Max Latency | Errors |
| --- | --- | --- | --- | --- |
| mysql_auto_increment |  |  |  |  |
| redis_incr |  |  |  |  |
| snowflake |  |  |  |  |

## 7. 注意事项

- 三轮必须使用同一台机器、同一套依赖、同一条 `wrk` 命令
- 每轮压测前都要清理 MySQL 和 Redis 状态
- 第一轮先拿 baseline，不要一开始就把并发打太高

# 压测结果
## 机器配置

- 时间：`2026-06-05 15:14:44 +0800`
- 系统：`Ubuntu 24.04`，内核 `6.17.0-35-generic`
- 环境：`VMware` 虚拟机
- CPU：`Intel(R) Core(TM) i5-9300H CPU @ 2.40GHz`
- 逻辑核数：`8`
- 内存：`15 GiB`
- 当前可用内存：约 `9.7 GiB`
- 磁盘：系统盘 `196G`，可用 `121G`

## mysql_auto_increment
running: wrk -t4 -c100 -d30s -s post_create.lua http://127.0.0.1:8888/api/v1/links
Running 30s test @ http://127.0.0.1:8888/api/v1/links
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    23.66ms    9.96ms 176.87ms   76.93%
    Req/Sec     1.07k   167.41     1.87k    68.78%
  128181 requests in 30.10s, 43.63MB read
Requests/sec:   4258.36
Transfer/sec:      1.45MB

## redis_incr
running: wrk -t4 -c100 -d30s -s post_create.lua http://127.0.0.1:8888/api/v1/links
Running 30s test @ http://127.0.0.1:8888/api/v1/links
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     4.34ms    3.48ms  66.18ms   72.29%
    Req/Sec     6.22k     0.97k    8.39k    70.67%
  744059 requests in 30.07s, 254.74MB read
Requests/sec:  24743.30
Transfer/sec:      8.47MB

## snowflake
running: wrk -t4 -c100 -d30s -s post_create.lua http://127.0.0.1:8888/api/v1/links
Running 30s test @ http://127.0.0.1:8888/api/v1/links
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     3.76ms    4.40ms  94.36ms   92.24%
    Req/Sec     7.78k     1.46k   11.79k    73.00%
  931160 requests in 30.07s, 331.23MB read
Requests/sec:  30967.95
Transfer/sec:     11.02MB

## snowflake+最终的落库
running: wrk -t4 -c100 -d30s -s post_create.lua http://127.0.0.1:8888/api/v1/links
Running 30s test @ http://127.0.0.1:8888/api/v1/links
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    23.40ms   11.02ms 147.60ms   74.73%
    Req/Sec     1.09k   217.72     2.19k    67.34%
  130388 requests in 30.09s, 46.38MB read
Requests/sec:   4333.36
Transfer/sec:      1.54MB

## 横向对比

### 结论

- 不考虑最终落库时，性能排序为：`snowflake > redis_incr > mysql_auto_increment`
- 恢复最终 MySQL 落库后，创建接口吞吐从 `3w+ req/s` 下降到 `4.3k+ req/s`
- 说明当前创建链路的主瓶颈不是 `short_code` 生成方式，而是数据库写入

### 结果分析

`mysql_auto_increment`

- 吞吐 `4258.36 req/s`，平均延迟 `23.66ms`
- 这组最弱，说明短码生成强依赖 MySQL 时，创建链路成本最高
- 如果实现包含“插入取自增 ID，再回写 short_code”，本质上会放大数据库写压力

`redis_incr`

- 吞吐 `24743.30 req/s`，平均延迟 `4.34ms`
- 相比 `mysql_auto_increment` 提升明显，说明将全局递增序号生成下沉到 Redis 是有效的
- 但它仍然慢于 `snowflake`，主要差异来自 Redis 网络往返

`snowflake`

- 吞吐 `30967.95 req/s`，平均延迟 `3.76ms`
- 四组里最好，说明本地进程内生成短码的成本最低
- 这组数据可以近似视为“只比较生成策略本身”的上限结果

`snowflake + 最终的落库`

- 吞吐 `4333.36 req/s`，平均延迟 `23.40ms`
- 和纯 `snowflake` 相比，吞吐显著下降，平均延迟明显上升
- 这说明一旦恢复真实落库流程，前置的短码生成优化收益会被数据库写入成本覆盖

### 当前阶段判断

- `v1` 阶段若只比较短码生成能力，优先级应为：`snowflake`、`redis_incr`、`mysql_auto_increment`
- 若比较真实创建接口链路，优化重点应从“生成方式”转向“降低 MySQL 写入压力”
- 这也是后续继续做缓存、批量写、异步化优化的依据

## wrk 指标说明

- `Avg`：平均值，用于看整体水平
- `Stdev`：标准差，用于看波动大小，值越小越稳定
- `Max`：最大值，用于看最慢请求或最高瞬时值
- `+/- Stdev`：落在 1 个标准差范围内的数据占比，用于辅助判断分布是否集中

重点看两类 `Stdev`：

- `Latency Stdev`：请求耗时波动，越大说明延迟抖动越明显
- `Req/Sec Stdev`：每秒吞吐波动，越大说明吞吐越不稳定

因此压测不能只看 `Requests/sec`，还要结合：

- 平均延迟是否可接受
- 延迟波动是否过大
- 最大延迟是否异常
- 吞吐波动是否明显

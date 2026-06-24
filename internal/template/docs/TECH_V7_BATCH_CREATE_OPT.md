# 短链系统 V7 批量创建优化方案

## 1. 目标

优化 `POST /api/v1/links/batch` 的性能，重点降低批量创建链路中的 Redis 往返次数和重复工作。

本版优先做：

- 批量读缓存
- 批量回填缓存
- `redis_incr` 批量取号

## 2. 当前现状

当前批量创建逻辑已经具备这些能力：

- 请求内去重
- 批量查 MySQL 已有记录
- 批量写 MySQL 新记录
- `mysql_auto_increment` 模式下批量事务生成短码

说明当前主要瓶颈不在 MySQL 批量能力，而在 Redis 侧仍有串行开销。

## 3. 当前主要问题

### 3.1 缓存读取仍是逐条 `GET`

`loadExistingBatchResults` 当前会循环调用单条 `GetLongToShort`。

问题：

- 一个批量请求会产生多次 Redis RTT
- 批量接口没有用上 `MGET`

### 3.2 缓存回填仍是逐条写入

当前成功创建后，会逐条执行：

- `SetLongToShort`
- `ShortCodeBloomAdd`

问题：

- Redis 往返次数偏多
- 没有利用 pipeline

### 3.3 `redis_incr` 模式逐条发号

当前 `redis_incr` 模式会循环调用 `NextCode`。

问题：

- 一次批量请求会产生多次 Redis 发号 RTT

## 4. 为什么不先做多协程

当前 `long_urls` 上限只有 `20`，而且 MySQL 已经是批量 SQL。

如果直接改成“每个 URL 一个 goroutine”，问题通常会更多：

- 容易把批量 SQL 打散成并发小 SQL
- 冲突和重试逻辑更复杂
- Redis / MySQL 连接池竞争更重

结论：

- 当前阶段，多协程不是第一优先级
- 优先减少网络往返，比优先增加 goroutine 更有效

## 5. 优化方案

### 5.1 批量读取 `long -> short` 缓存

新增能力：

- `GetLongToShortBatch(ctx, userID, normalizedURLs []string) (map[string]string, error)`

建议实现：

- Redis `MGET`
- 返回 `normalizedURL -> short_code`

目标：

- `loadExistingBatchResults` 不再逐条单独 `GET`

### 5.2 批量回填缓存

新增能力：

- `FillCreateCachesBatch(ctx, userID, records []model.ShortLink) error`

建议实现：

- pipeline 批量 `SET long->short`
- pipeline 批量 `BF.ADD short_code`
- Bloom 不可用时沿用当前降级逻辑

目标：

- 创建成功后不再逐条回填缓存

### 5.3 `redis_incr` 批量取号

新增能力：

- `NextCodes(ctx, n int) ([]string, error)`

建议实现：

- 一次 `INCRBY sequence_key n`
- 本地回推连续 ID 区间
- 统一转 Base62

目标：

- 避免一次批量请求触发多次 Redis 发号 RTT

## 6. 推荐实施顺序

1. 批量读缓存
2. pipeline 回填缓存
3. `redis_incr` 批量取号
4. 只有未来批量上限明显增大时，再评估受控并发

## 7. 边界

- 不修改接口结构
- 不提高 `long_urls` 上限
- 不改变当前成功/失败语义
- 不引入分布式锁
- 不把 MySQL 批量写改成多协程小写入

## 8. 验收

- 批量创建接口返回结构不变
- 缓存读取支持批量 `MGET`
- 缓存回填支持 pipeline
- `redis_incr` 模式支持批量取号
- 压测下平均耗时和 P95 优于当前实现

## 9. 结论

批量创建接口可以做并发优化，但当前更值得优先落地的是：

- 批量 Redis 读取
- 批量 Redis 回填
- 批量发号

相比直接“每条 URL 起一个 goroutine”，这条路线改动更小，也更符合当前项目结构。

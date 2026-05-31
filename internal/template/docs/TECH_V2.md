# 短链系统 V2 Redis 优化建议

## 1. 目标

基于当前 V1 实现，梳理哪些步骤适合用 Redis 优化，哪些步骤不值得引入缓存。

前提：

- MySQL 仍是最终真值
- Redis 只负责加速与削峰
- V1 不启用过期逻辑
- 若使用布隆过滤器，需要 RedisBloom 或等价能力支持

## 2. 当前链路

创建短链：

1. 校验 `long_url`
2. 规范化 `long_url`
3. 计算 `url_hash`
4. 按 `url_hash` 查 MySQL
5. 命中相同长链则复用
6. 未命中则生成 `short_code`

跳转短链：

1. 校验 `short_code`
2. 按 `short_code` 查 MySQL
3. 累加 `visit_count`
4. 返回 `302`

## 3. 可优先优化的步骤

### 3.1 跳转查询缓存

优化点：

- `short_code -> original_url`

原因：

- 跳转是典型读多写少
- 热门短链会重复命中同一条记录
- 当前每次都查 MySQL，热点会直接压库

建议：

- Key：`shortlink:code:{short_code}`
- Value：跳转所需最小字段，如 `id`、`original_url`
- 不存在的短码可加短 TTL 空值缓存

结论：

- 这是 V2 收益最高的 Redis 优化点

### 3.2 访问次数异步累加

优化点：

- `visit_count = visit_count + 1`

原因：

- 当前每次跳转都同步写 MySQL
- 热门短链会形成单行更新热点

建议：

- 跳转成功后先写 Redis `INCR`
- 定时或批量刷回 MySQL
- 回刷时使用增量累加，不做覆盖写

结论：

- 应与跳转缓存一起落地

### 3.3 `url_hash` 布隆过滤器

优化点：

- `url_hash` 是否可能已存在

原因：

- 当前创建短链时，会先按 `url_hash` 查 MySQL
- 当大多数请求都是新长链时，这次查库通常无结果
- 布隆过滤器适合拦截这类“不存在”的查库请求

建议：

- 先对规范化长链计算 `url_hash`
- 先查布隆过滤器
- 若判断“不存在”，可跳过按 `url_hash` 查 MySQL 候选记录
- 若判断“可能存在”，再查 MySQL 候选记录并做精确比对
- 创建成功后，将 `url_hash` 写入布隆过滤器
- 若要严格保证“相同长链复用同一短链”，仍需结合精确缓存占位、分布式锁或数据库唯一约束等一致性手段

结论：

- 它优化的是“减少无效查库”
- 不能替代 MySQL 真值判断

### 3.4 长链复用缓存

优化点：

- `normalized_long_url -> short_code`

原因：

- 同一长链重复创建时，会反复走 MySQL 查重

建议：

- Key：`shortlink:long:{url_hash}:{normalized_long_url}`
- Value：`short_code`
- 创建成功后主动回填
- 正常命中缓存可直接返回；仅在缓存重建、失效恢复等场景再回源 MySQL 校验

结论：

- 这是创建链路的优化点，但优先级低于跳转链路

## 4. 风险与边界

- 布隆过滤器只能用于“是否可能存在”的预判，不能作为复用真值
- 当前系统没有针对规范化长链的数据库唯一约束，创建链路若要严格保证复用语义，需要额外一致性手段
- 跳转缓存与长链复用缓存都需要明确失效策略，否则会出现脏数据
- `visit_count` 异步回刷后，MySQL 中的访问次数将变为最终一致，而不是实时一致

## 5. 推荐落地顺序

第一阶段：

1. `short_code -> original_url` 跳转缓存
2. `visit_count` Redis 异步累加

第二阶段：

1. `url_hash` 布隆过滤器
2. `normalized_long_url -> short_code` 长链复用缓存

第三阶段：

1. 空值缓存
2. 热点 key 监控
3. 预热策略

## 6. 一致性要求

- MySQL 始终是最终真值
- 创建成功后主动回填缓存
- 后续若支持编辑、禁用、软删，必须同步删除相关缓存

涉及的关键 key：

- `shortlink:code:{short_code}`
- `shortlink:long:{url_hash}:{normalized_long_url}`
- `shortlink:visit:{id}` 或 `shortlink:visit:{short_code}`
- `shortlink:bloom:url_hash`

## 7. 结论

基于当前项目，最值得做的 Redis 优化是：

1. 跳转读缓存
2. 访问次数异步计数
3. `url_hash` 布隆过滤器
4. 长链复用缓存

其中前两项优先级最高，创建链路相关优化放在第二阶段即可。

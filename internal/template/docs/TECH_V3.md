# 短链系统 V3 技术方案

## 1. 目标

V3 在 V2 缓存链路基础上，继续降低高并发场景下的数据库压力。

本阶段只做三件事：

- `visit_count` 改为 Redis 聚合后异步回刷 MySQL
- 跳转链路缓存未命中时使用 `singleflight` 合并查库
- 创建链路首次新建时使用并发控制，减少重复发号和重复新建

## 2. 核心优化

### 2.1 visit_count 异步回刷

- 跳转成功后不再同步更新 MySQL
- 改为 Redis `INCR shortlink:visit:{id}`
- 后台任务批量回刷 MySQL
- 回刷使用增量更新
- Redis 中只保存增量，不保存最终总值
- 回刷流程完整扫描 `shortlink:visit:*`
- 每批使用 Redis `MGET` 批量读取增量
- MySQL 使用单条批量 `UPDATE ... CASE id WHEN ... THEN ... END`
- 回刷成功后批量 `DEL` 对应 key
- 回刷失败时保留增量，等待下次重试

收益：

- 跳转主链路去掉同步写库
- 热点短链不再频繁更新同一行
- Redis 和 MySQL 的往返次数进一步下降

### 2.2 跳转链路 `singleflight`

- `short_code -> original_url` 缓存未命中时进入 `singleflight`
- 直接使用 `github.com/zeromicro/go-zero/core/syncx`
- 并发控制 key 使用 `redirect:{short_code}`
- Redis 缓存 key 继续使用 `shortlink:code:{short_code}`，不与并发控制 key 混用
- 由一个请求回源 MySQL 并回填缓存，其余请求共享结果

收益：

- 避免热点短链缓存失效时并发击穿数据库

### 2.3 创建链路并发控制

- 以 `normalized_long_url` 作为控制粒度
- 优先使用进程内 `singleflight`
- key 使用 `create:{normalized_long_url}`
- 控制区内完成查缓存、查 MySQL、新建、回填

收益：

- 减少重复发号
- 减少重复新建

## 3. 接口链路

### 3.1 创建短链接口

建议顺序：

1. 校验并规范化 `long_url`
2. 查 Bloom
3. 查 `normalized_long_url -> short_code`
4. 若未命中，进入 `create:{normalized_long_url}` 的 `singleflight`
5. 控制区内再次查缓存和 MySQL
6. 若仍未命中，则生成 `short_code` 并写入 MySQL
7. 回填缓存和 Bloom

### 3.2 跳转短链接口

建议顺序：

1. 校验 `short_code`
2. 查 `shortlink:code:{short_code}`
3. 若未命中，进入 `redirect:{short_code}` 的 `singleflight`
4. 控制区内再次查 `shortlink:code:{short_code}`
5. 若仍未命中，则回源 MySQL 并回填缓存
6. 跳转成功后写 Redis 访问增量
7. 返回 `302`

补充：

- 主链路只写 Redis 增量，不再同步写 MySQL
- MySQL 中的 `visit_count` 由后台回刷任务更新

## 4. 关键组件

- `shortlink:visit:{id}`：访问次数增量
- `shortlink:code:{short_code}`：短链跳转精确缓存
- `singleflight.Group`：基于 `github.com/zeromicro/go-zero/core/syncx` 的统一并发控制组件
- `redirect:{short_code}` / `create:{normalized_long_url}`：通过不同 key 前缀区分跳转和创建链路
- 后台回刷任务：批量同步 `visit_count`

回刷任务建议职责：

- 周期性扫描 `shortlink:visit:*`
- 读取每批 key 的增量值
- 批量更新 `visit_count`
- 成功后批量删除对应 key
- 失败则保留 key，等待下次重试


建议执行方式：

1. 按固定周期执行，例如每 `5s`
2. 使用 `SCAN` 从游标 `0` 开始，循环扫描直到游标回到 `0`
3. `Batch` 作为每次 `SCAN` 的 `count`，例如 `100`
4. 对每批扫描结果执行一次 Redis `MGET`
5. 解析 `shortlink:visit:{id}` 中的 `id`
6. 生成单条批量 SQL：

```sql
UPDATE short_links
SET visit_count = visit_count + CASE id
  WHEN ? THEN ?
  WHEN ? THEN ?
  ...
  ELSE 0
END
WHERE id IN (?, ?, ...)
```

7. 在一个 MySQL 事务内执行该批量更新
8. 事务提交成功后，再批量 `DEL` 本轮回刷对应的 Redis key

说明：

- 第一版不使用 `KEYS`
- 第一版只启一个回刷 worker
- 第一版主链路只执行 Redis `INCR`，不再同步更新 MySQL
- 当前实现优先“不丢数据”，因此仍采用“先写 MySQL，再删 Redis”

## 5. 边界

- `visit_count` 从实时一致变为最终一致
- 当前方案在极端故障场景下可能存在少量重复累计风险
- 当前方案优先“不丢数据”，不保证精确一次回写
- `singleflight` 只能解决单实例并发
- `singleflight` key 只用于并发控制，不承担缓存语义
- 多实例场景下，创建链路若要求更强控制，后续再补 Redis 分布式锁
- `visit_count` 更适合作为统计口径，不适合作为强一致业务计费依据

## 6. 结论

V3 的重点不是增加新的缓存类型，而是减少重复工作：

- 跳转链路减少重复查库
- 访问次数去掉同步写库
- 创建链路减少重复新建

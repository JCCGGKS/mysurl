# 短链系统 V3 技术方案

## 1. 目标

V3 在 V2 缓存链路基础上，继续降低高并发场景下的数据库压力。

本阶段只做三件事：

- `visit_count` 改为 Redis 聚合后异步回刷 MySQL
- 跳转链路缓存未命中时使用 `singleflight` 合并查库
- 创建链路首次新建时使用并发控制，减少重复发号和唯一键冲突

默认约束：

- HTTP 接口不变
- MySQL 仍是最终真值
- V2 的精确缓存和 Bloom 继续保留

## 2. 核心优化

### 2.1 visit_count 异步回刷

- 跳转成功后不再同步更新 MySQL
- 改为 Redis `INCR shortlink:visit:{id}`
- 后台任务批量回刷 MySQL
- 回刷使用增量更新：`visit_count = visit_count + ?`

收益：

- 跳转主链路去掉同步写库
- 热点短链不再频繁更新同一行

### 2.2 跳转链路 `singleflight`

- `short_code -> original_url` 缓存未命中时进入 `singleflight`
- key 使用 `redirect:{short_code}`
- 由一个请求回源 MySQL 并回填缓存，其余请求共享结果

收益：

- 避免热点短链缓存失效时并发击穿数据库

### 2.3 创建链路并发控制

- 以 `normalized_long_url` 作为控制粒度
- 优先使用进程内 `singleflight`
- key 使用 `create:{normalized_long_url}`
- 控制区内完成查缓存、查 MySQL、新建、回填
- `uk_original_url` 继续保留为最终兜底

收益：

- 减少重复发号
- 减少唯一键冲突

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
2. 查 `short_code -> original_url` 缓存
3. 若未命中，进入 `redirect:{short_code}` 的 `singleflight`
4. 回源 MySQL 并回填缓存
5. 跳转成功后写 Redis 访问增量
6. 返回 `302`

## 4. 关键组件

- `shortlink:visit:{id}`：访问次数增量
- `singleflight.Group`：跳转链路查库合并
- `singleflight.Group`：创建链路并发控制
- 后台回刷任务：批量同步 `visit_count`

## 5. 边界

- `visit_count` 从实时一致变为最终一致
- `singleflight` 只能解决单实例并发
- 多实例场景下，创建链路若要求更强控制，后续再补 Redis 分布式锁
- `uk_original_url` 不能移除

## 6. 结论

V3 的重点不是增加新的缓存类型，而是减少重复工作：

- 跳转链路减少重复查库
- 访问次数去掉同步写库
- 创建链路减少重复新建

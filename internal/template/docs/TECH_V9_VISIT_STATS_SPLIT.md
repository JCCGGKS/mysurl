# 短链系统 V9 访问统计拆分与全量回刷方案

## 1. 目标

本次改造聚焦 `visit` 统计链路，目标是：

- 将访问统计从 `short_links` 主表中拆出
- Redis 不再记录增量，而是记录当前全量访问次数
- 后台 worker 回刷 MySQL 时直接覆盖统计值
- 回刷后不再删除 Redis 计数 key
- 短链列表读取 `visit` 时优先从缓存获取

核心思路是把“短链映射”和“访问统计”拆成两条职责不同的链路：

- `short_links` 只负责短码映射和核心元数据
- `visit_stats` 只负责访问统计

## 2. 改造前的问题

改造前方案的主要特征是：

- `short_links.visit_count` 存在于主表
- 跳转成功后 Redis 执行 `INCR shortlink:visit:{id}`
- Redis 中保存的是增量
- worker 周期性扫描 Redis，把增量累加回 `short_links.visit_count`
- 回刷成功后删除 Redis key

这个方案可以工作，但有几个问题：

- 主表承载了高频统计更新，热点短链会持续修改 `short_links`
- 统计字段和核心映射表耦合过深
- worker 需要处理“增量回放”和“回刷后删除 key”
- 一旦 MySQL 成功、Redis 删除失败，下一轮会重复累计
- 列表查询只能读主表中的旧值，不能优先使用 Redis 中的最新计数

## 3. 改造后的整体设计

### 3.1 表结构调整

主表 `short_links` 移除 `visit_count` 字段。

新增独立统计表：

```sql
CREATE TABLE `visit_stats` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `short_link_id` bigint unsigned NOT NULL COMMENT '短链ID',
  `visit_count` bigint unsigned NOT NULL DEFAULT 0 COMMENT '访问次数',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_short_link_id` (`short_link_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='短链访问统计表';
```

职责拆分后：

- `short_links`：短码、长链、归属用户、状态、创建时间
- `visit_stats`：访问次数

### 3.2 Redis 计数模型调整

Redis key 仍然沿用：

- `shortlink:visit:{id}`

但语义发生变化：

- 改造前：key 中保存的是“待回刷增量”
- 改造后：key 中保存的是“当前全量访问次数”

这意味着 worker 不再做：

- `visit_count = visit_count + delta`

而是做：

- `visit_count = redis_current_count`

### 3.3 跳转链路调整

跳转成功后，不再直接裸 `INCR`。

因为如果 Redis key 不存在，直接 `INCR` 会从 `1` 开始，可能把 MySQL 中历史计数覆盖掉。

因此当前实现改成：

1. 先读取 `visit_stats` 中该短链的当前计数作为基线
2. 对 Redis 执行原子逻辑：
   - 如果 key 已存在，直接 `INCR`
   - 如果 key 不存在，先以 MySQL 基线值初始化，再加 `1`

可理解为：

- 有缓存时，Redis 继续累加
- 冷 key 首次被访问时，以数据库当前值为起点继续累加

这样可以保证 Redis 中维护的是全量值，而不是脱离历史的局部值。

### 3.4 worker 回刷链路调整

worker 仍然会定时：

- `SCAN shortlink:visit:*`
- 对当前批次使用 `MGET`
- 解析 `short_link_id`

但回刷动作从“增量累加”改成了“全量覆盖”：

```sql
INSERT INTO visit_stats (short_link_id, visit_count)
VALUES ...
ON DUPLICATE KEY UPDATE
  visit_count = VALUES(visit_count)
```

因此：

- 不再对 `short_links` 主表做 `UPDATE`
- 不再在回刷成功后删除 Redis key

Redis key 会持续存在，并作为更实时的统计值来源。

## 4. 当前读取路径

### 4.1 跳转接口

跳转接口的核心职责没有变：

- 优先走短码缓存
- 缓存 miss 时回源 MySQL
- 成功跳转后更新访问统计

变化点只在统计写路径：

- 统计写入 Redis 的值变成“全量值”

### 4.2 短链列表接口

列表获取 `visit` 时，当前策略是：

1. 先批量查 `visit_stats`
2. 再批量查 Redis `shortlink:visit:{id}`
3. 如果 Redis 中存在值，则覆盖 DB 值
4. 返回给前端

这等价于“优先读缓存”：

- Redis 命中时，返回更实时的值
- Redis miss 时，回落到 `visit_stats`

这样列表接口不必等 worker 回刷完成，能尽量展示最新计数。

## 5. 一致性取舍

当前方案属于：

- 最终一致
- Redis 更实时
- MySQL 为持久化落点

一致性特征如下：

- 跳转刚发生后，Redis 中的计数会先变
- `visit_stats` 要等下一轮 worker 才会追平
- 列表接口因为优先读 Redis，通常能看到最新值
- 如果 Redis 丢失某个 key，下次访问会以 `visit_stats` 当前值作为基线重新建立

与旧方案相比，本次方案的优势是：

- 不再依赖“回刷成功后删除 key”
- 不再存在因为重复删除失败而导致的重复增量累计
- 主表不再承受高频 `visit_count` 更新

仍然保留的边界是：

- `visit_stats` 不是实时强一致
- 如果 Redis 整体丢失且某个短链长时间没有新访问，则该短链只能先显示数据库中的旧值

## 6. 收益

本次改造的主要收益有：

- 主表减压：`short_links` 不再承载高频统计写入
- 职责分离：映射数据与统计数据拆开
- 回刷逻辑简化：从“增量累计 + 删除 key”变成“全量覆盖”
- 统计读取更实时：列表优先使用 Redis 中的计数
- 为后续扩展统计维度预留空间，例如 UV、按天统计、来源统计

## 7. 上线与迁移注意事项

如果系统已经有历史数据，上线时需要关注以下事项：

### 7.1 建表

先创建 `visit_stats` 表。

### 7.2 历史数据迁移

如果旧版本的 `short_links.visit_count` 中已经有历史访问次数，需要在切换前迁移到 `visit_stats`。

迁移目标是：

- `visit_stats.short_link_id = short_links.id`
- `visit_stats.visit_count = short_links.visit_count`

否则切换后，Redis 冷 key 首次初始化时只能从 `0` 起步，历史值会丢失。

仓库中已提供批量迁移脚本：

- `internal/template/sqls/migrate_visit_stats.sh`

默认行为：

- 自动确保 `visit_stats` 表存在
- 按 `short_links.id` 分批迁移
- 默认每批 `1000` 条
- 默认只迁移 `visit_count > 0` 的记录
- 使用 `ON DUPLICATE KEY UPDATE`，可重复执行

示例：

```bash
bash internal/template/sqls/migrate_visit_stats.sh
```

指定批大小：

```bash
BATCH_SIZE=5000 bash internal/template/sqls/migrate_visit_stats.sh
```

指定迁移区间：

```bash
START_ID=1 END_ID=200000 bash internal/template/sqls/migrate_visit_stats.sh
```

### 7.3 Redis 旧 key 兼容

本次 Redis key 名称没有变化，但语义从“增量”改成了“全量”。

因此上线时要注意：

- 如果现网 Redis 中仍残留旧语义的增量 key，不能直接当成新语义的全量 key 使用
- 更稳妥的方式是切换前清理旧的 `shortlink:visit:*` key，或者先完成历史迁移再统一重建

### 7.4 回滚影响

如果需要回滚到旧方案，要特别注意：

- 新方案 Redis 中保存的是全量值
- 旧方案 Redis 中期望的是增量值

两者不能直接混用。

## 8. 后续可继续优化的方向

如果后续流量继续上涨，可以继续考虑：

- 跳转链路对 `visit_stats` 基线查询做本地短期缓存，减少冷 key 首访查库
- worker 做分片扫描，降低单实例回刷压力
- 对统计链路增加按天聚合表，而不是只有总计数
- 将超热点短链的访问统计进一步异步化或分桶化

## 9. 总结

这次改造的本质不是单纯换一张表，而是把 `visit` 从“主表上的附属字段”升级成“独立统计链路”。

新的统计链路具有三个关键特征：

- Redis 保存全量计数
- MySQL 持久化落到独立 `visit_stats` 表
- 列表读取优先使用缓存值

这样既降低了 `short_links` 主表压力，也让访问统计链路更符合后续扩展方向。

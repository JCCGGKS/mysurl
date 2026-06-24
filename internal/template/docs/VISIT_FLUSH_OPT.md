# visit_count 异步回写优化说明

## 1. 背景

短链跳转成功后，系统会记录访问次数 `visit_count`。

为了避免每次跳转都同步更新 MySQL，当前实现采用：

- 主链路只写 Redis 增量
- 后台 worker 周期性把增量回刷到 MySQL

Redis 增量 key 格式：

- `shortlink:visit:{id}`

## 2. 初始实现

最初版本的回写流程是：

1. 每个周期执行一次 `SCAN shortlink:visit:*`
2. 只扫描从游标 `0` 开始的一批 key
3. 对每个 key 单独执行一次 `GET`
4. 在 MySQL 事务内逐条执行：

```sql
UPDATE short_links SET visit_count = visit_count + ? WHERE id = ?
```

5. 事务成功后，对每个 key 单独执行一次 `DEL`

这个版本可以工作，但有几个明显问题：

- 单次只扫描一批 key，没有完整遍历整个游标
- Redis 读操作是 `SCAN + N 次 GET`
- MySQL 写操作是事务内多条 `UPDATE`
- Redis 删除操作是 `N 次 DEL`

## 3. 本次优化内容

### 3.1 完整扫描 Redis 游标

回写 worker 现在不再只调用一次 `SCAN`，而是：

- 从 `cursor=0` 开始
- 按 `Batch` 作为每次 `SCAN` 的 `count`
- 持续循环扫描
- 直到游标重新回到 `0`

这样单次 flush 周期会完整扫描当前 `shortlink:visit:*` 的 keyspace，不再只处理一小批开头结果。

### 3.2 批量获取增量

对每批扫描出的 key，不再逐个 `GET`，而是改成一次 Redis `MGET`。

优化前：

- `SCAN + N 次 GET`

优化后：

- `SCAN + 1 次 MGET`

收益：

- 降低 Redis 往返次数
- 批量回刷时吞吐更高

### 3.3 批量删除增量 key

MySQL 回写成功后，不再逐个 `DEL`，改为一次批量删除：

- `DEL key1 key2 key3 ...`

优化前：

- `N` 次删除请求

优化后：

- `1` 次删除请求

### 3.4 MySQL 批量更新

MySQL 事务内不再逐条执行：

```sql
UPDATE short_links SET visit_count = visit_count + ? WHERE id = ?
```

而是改成单条批量 SQL：

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

收益：

- 减少 MySQL 往返次数
- 降低事务内多次执行 SQL 的开销

## 4. 当前回写流程

当前 visit flush worker 的整体流程如下：

1. 定时触发，例如每 `5s`
2. 使用 `SCAN` 循环扫描 `shortlink:visit:*`
3. 每批 key 使用一次 `MGET` 读取增量
4. 从 key 中解析出短链 `id`
5. 组装本轮待回写列表
6. 在一个 MySQL 事务内执行单条批量 `UPDATE`
7. 事务成功后，一次批量 `DEL` 删除本轮 Redis key

## 5. 当前一致性取舍

当前方案采用的是：

- 先写 MySQL
- 再删 Redis

这个顺序的核心目标是：

- 优先不丢数据

含义是：

- 如果 MySQL 更新失败，Redis 增量仍然保留，下次还能重试
- 如果 MySQL 成功但 Redis 删除失败，下一轮可能重复回写同一批增量

因此当前方案的特征是：

- 最终一致
- 优先不丢数据
- 不保证精确一次回写
- 极端故障下可能存在少量重复累计

## 6. 为什么暂时不继续复杂化

更强的一致性方案通常需要：

- Redis Lua 脚本
- processing key
- 或更复杂的幂等恢复机制

这些方案可以进一步降低重复回写风险，但实现复杂度更高。

在当前项目阶段：

- 并发量不大
- `visit_count` 属于统计字段
- 不作为计费或强一致业务依据

因此当前实现优先选择更简单、偏保守的方案：

- 宁可极端情况下少量重复
- 也尽量避免访问计数丢失

## 7. 适用边界

当前方案适合：

- 中小规模流量
- 统计口径的访问计数
- 对实时一致性要求不高的短链系统

当前方案不适合直接用于：

- 计费
- 配额扣减
- 奖励发放
- 审计级精确统计

如果后续并发或一致性要求继续提高，可再考虑：

- 拆分统计表
- Redis Lua 原子迁移增量
- processing key 重试机制
- 更严格的幂等回写方案

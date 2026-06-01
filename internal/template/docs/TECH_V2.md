# 短链系统 V2 缓存优化方案

## 1. 目标
在不改变 V1 业务语义的前提下，用 Redis 降低 MySQL 查询压力。
## 2. 当前实现

### 2.1 创建短链

当前 `CreateLink` 的实际处理顺序是：

1. 校验 `long_url`
2. 规范化 `long_url`
3. 计算 `url_hash`
4. 按 `url_hash` 查询 MySQL 候选记录(可优化)
5. 对候选记录逐条比对规范化长链(可优化)
6. 命中则复用已有 `short_code`
7. 未命中则生成新的 `short_code`


### 2.2 跳转短链

当前 `Redirect` 的实际处理顺序是：

1. 校验 `short_code`
2. 按 `short_code` 查询 MySQL(可优化)
3. 更新 `visit_count`(可优化，放到V3)
4. 返回目标长链并发起 `302`

## 3. 优化思路

### 3.1 短链到长链

建议方案：

- 使用 `short_code -> original_url` 的映射缓存跳转结果
- Redis Key 使用 `shortlink:code:{short_code}`
- Value 只保留跳转所需最小字段，例如 `original_url`


### 3.2 长链到短链

- 直接缓存 `normalized_long_url -> short_code`
- Redis Key 使用 `shortlink:long:{normalized_long_url}`
- Value 为复用得到的 `short_code`
- 命中时可以直接返回已有短链，避免重复查库


创建成功后回填：

- `shortlink:long:{normalized_long_url}`

### 3.3 布隆过滤器减少不存在数据的查库次数

建议方案：

- 对规范化长链建立布隆过滤器
- 如果布隆判断不存在，可以跳过一次按 `url_hash` 的候选查库
- 如果布隆判断可能存在，仍然继续查 MySQL 并做精确比对
- 创建成功后回填布隆过滤器

收益：

- 新长链减少无结果查库
- 降低创建链路中不存在数据的查库次数


## 4. 优化后接口逻辑顺序

### 4.1 创建短链接口

优化后的建议顺序：

1. 校验 `long_url`
2. 规范化 `long_url`
3. 查询 `shortlink:long:{normalized_long_url}`
4. 若命中，直接返回已有 `short_code`
5. 若未命中，查询 `shortlink:bloom:normalized_long_url`
6. 若布隆判断可能存在，计算 `url_hash` 并查询 MySQL 候选记录
7. 对候选记录逐条比对规范化长链
8. 若命中已有记录，则返回已有 `short_code`，并回填 `shortlink:long:{normalized_long_url}`
9. 若仍未命中，则计算 `url_hash` 并生成新的 `short_code`
10. 写入 MySQL
11. 回填 `shortlink:long:{normalized_long_url}` 和 `shortlink:bloom:normalized_long_url`

### 4.2 跳转短链接口

优化后的建议顺序：

1. 校验 `short_code`
2. 查询 `shortlink:code:{short_code}`
3. 若命中，直接返回 `original_url`
4. 若未命中，查询 MySQL
5. 若命中记录，则回填 `shortlink:code:{short_code}`
6. 按现有逻辑更新 MySQL 中的 `visit_count`
7. 返回目标长链并发起 `302`

## 5. 关键 Key 设计

- `shortlink:code:{short_code}`
- `shortlink:long:{normalized_long_url}`
- `shortlink:bloom:normalized_long_url`

## 6. 结论

V2 优先优化三个点：

- 短链到长链：解决热点读问题
- 长链到短链：解决重复查重问题
- 布隆过滤器：减少不存在数据的查库次数

按当前项目形态，Redis 在 V2 的角色应当是：

- 为跳转提供读缓存
- 为创建提供精确复用缓存
- 为不存在数据提供查库前预判

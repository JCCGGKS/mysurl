# 短链系统 V1 技术方案

## 1. 目标

定义 V1 的数据模型、唯一短链生成方案和运行时处理约束。

## 2. 唯一短链生成方案

规则：

- `url_hash` 用于长链复用判断
- `short_code` 用于对外短链标识
- `url_hash` 命中后，必须再次比对规范化长链字符串
- `short_code` 冲突时重新生成并重试

步骤：

1. 校验 `long_url`
2. 对 `long_url` 做 V1 规范化
3. 计算 `url_hash`
4. 按 `url_hash` 查询候选记录
5. 若存在相同规范化长链，则直接返回已有短链
6. 若不存在，则生成新的 `short_code`

短码规则：

- 字符集：Base62
- 字符顺序固定为：`0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
- 长度：`Short.CodeLength`
- 唯一性：`uk_short_code`
- `short_url` 由 `BaseURL + short_code` 组成
- 编码与解码必须使用同一套字符表

## 3. `short_code` 生成方案

### 3.1 mysql-auto_increment

- 方式：使用 MySQL 自增 ID，转 Base62 生成短码
- 优点：实现简单，ID 单调递增
- 缺点：依赖数据库发号能力；先插入拿 ID，再回写 `short_code`，链路较重

### 3.2 redis-incr

- 方式：使用 Redis `INCR` 生成自增 ID，转 Base62 生成短码
- 优点：发号性能高于 MySQL
- 缺点：引入额外中间件依赖；仍然属于中心化发号

### 3.3 snowflake

- 方式：使用本地雪花算法生成 ID，转 Base62 生成短码
- 实现：使用 `bwmarrin/snowflake` 生成全局唯一 ID
- 配置：通过 `Short.Snowflake.WorkerID` 指定节点号
- 优点：不依赖外部发号器，性能高
- 缺点：需要处理机器号、时钟回拨等问题

### 3.4 code_manager-策略模式

- 三种方案同时实现：`mysql-auto_increment`、`redis-incr`、`snowflake`
- 运行时通过 `Short.Provider` 选择具体生成方案
- 配置值：`mysql_auto_increment`、`redis_incr`、`snowflake`
- 默认值：`mysql_auto_increment`
- 设计模式采用策略模式
- 引入统一的 `CodeManager` 作为短码生成策略的管理入口
- `CodeManager` 负责按配置选择具体 provider，并管理 provider 与生成器实现的映射关系
- 业务层仅依赖统一生成入口，不直接依赖具体生成方案
- 具体发号逻辑仍由各 provider 自己实现，便于后续扩展和替换
- `snowflake` 方案依赖 `Short.Snowflake.WorkerID` 配置节点号
- 通过集中管理策略注册与选择逻辑，降低业务层与具体实现的耦合

## 4. 数据模型

### 4.1 主表 `short_links`

| 字段 | 类型 | 约束/说明 |
| --- | --- | --- |
| `id` | bigint unsigned | 主键，自增 |
| `short_code` | varchar(16) | 短码，唯一 |
| `original_url` | varchar(2048) | 规范化后的长链 |
| `url_hash` | char(64) | 规范化长链的哈希值，用于辅助查重 |
| `visit_count` | bigint unsigned | 访问次数，默认 `0` |
| `expires_at` | datetime | 可空，过期时间预留字段，V1 不启用 |
| `status` | tinyint unsigned | 状态，默认 `1` |
| `created_at` | datetime | 创建时间 |
| `updated_at` | datetime | 更新时间 |
| `deleted_at` | datetime | 可空，软删除预留字段 |

### 4.2 索引

- `PRIMARY KEY (id)`
- `UNIQUE KEY uk_short_code (short_code)`
- `UNIQUE KEY uk_original_url (original_url)`
- `KEY idx_url_hash (url_hash)`

### 4.3 说明

- 相同长链以规范化后的 `original_url` 字符串一致为准
- `original_url` 保存规范化后的长链，并建立唯一约束
- `url_hash` 仅作为查重辅助键，不作为最终唯一真值
- `expires_at` 仅保留表结构，不参与 V1 创建和跳转逻辑
- 若 `url_hash` 相同但长链字符串不同，视为哈希冲突，继续生成新的 `short_code`

## 5. 运行时约束

- `long_url` 为空返回 `400`
- `long_url` 非法返回 `400`
- 短码冲突时重新生成并重试写入
- 成功跳转后 `visit_count` 累加 `1`
- 查询条件：`deleted_at IS NULL`
- 数据库异常返回 `500`

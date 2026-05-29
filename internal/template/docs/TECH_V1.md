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
7. 写入 `{short_code, original_url, url_hash}`

短码规则：

- 字符集：Base62
- 字符顺序固定为：`0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
- 长度：`Short.CodeLength`
- 唯一性：`uk_short_code`
- `short_url` 由 `BaseURL + short_code` 组成
- 编码与解码必须使用同一套字符表

## 3. 数据模型

### 3.1 主表 `short_links`

| 字段 | 类型 | 约束/说明 |
| --- | --- | --- |
| `id` | bigint unsigned | 主键，自增 |
| `short_code` | varchar(16) | 短码，唯一 |
| `original_url` | varchar(2048) | 原始长链 |
| `url_hash` | char(64) | `original_url` 去除尾部 `/` 后字符串的哈希值，用于查重 |
| `visit_count` | bigint unsigned | 访问次数，默认 `0` |
| `status` | tinyint unsigned | 状态，默认 `1` |
| `expires_at` | datetime | 过期时间，可空 |
| `created_at` | datetime | 创建时间 |
| `updated_at` | datetime | 更新时间 |
| `deleted_at` | datetime | 可空，软删除预留字段 |

### 3.2 索引

- `PRIMARY KEY (id)`
- `UNIQUE KEY uk_short_code (short_code)`
- `KEY idx_url_hash (url_hash)`

### 3.3 说明

- 相同长链以 `original_url` 去除尾部 `/` 后字符串一致为准
- `original_url` 不是数据库唯一键
- `url_hash` 仅作为查重辅助键，不作为最终唯一真值
- 若 `url_hash` 相同但长链字符串不同，视为哈希冲突，继续生成新的 `short_code`

## 4. 运行时约束

- `long_url` 为空返回 `400`
- `long_url` 非法返回 `400`
- `Short.ExpairedDays=0` 表示永久有效
- 短码冲突时重新生成并重试写入
- 成功跳转后 `visit_count` 累加 `1`
- 查询条件：`deleted_at IS NULL`，且 `expires_at IS NULL OR expires_at > 当前时间`
- 数据库异常返回 `500`

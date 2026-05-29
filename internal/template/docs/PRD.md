# 短链系统 V1 MVP PRD

## 1. 文档目的

定义短链系统 V1 的功能边界、接口约束和基础数据模型，作为后续服务实现与测试验收依据。

## 2. 目标与非目标

### 2.1 目标

- 创建短链
- 访问短链并返回 `302` 跳转
- 使用 MySQL 持久化短链映射

### 2.2 非目标

- 用户体系与权限控制
- 自定义短码
- 短链编辑、禁用、删除
- 统计报表、风控、二维码、管理后台

## 3. 功能需求

### 3.1 创建短链

输入：

- `long_url`：必填

规则：

- 只接受合法 `http://` 或 `https://` URL
- 相同长链重复提交时复用已有短链
- 短码由系统生成
- 短码必须全局唯一
- 第一版短链默认按配置项 `Short.ExpairedDays` 控制有效期，默认值 `0` 表示永久有效

输出：

- `short_code`
- `short_url`
- `original_url`

说明：

- `short_url` 由配置项 `Short.BaseURL` 与 `short_code` 拼接生成
- 短码长度由配置项 `Short.CodeLength` 控制，V1 默认值为 `4`
- 默认有效期由配置项 `Short.ExpairedDays` 控制，单位为天，默认值 `0` 表示永久有效

### 3.2 访问短链

输入：

- 路径参数 `code`

规则：

- 按 `short_code` 精确查询
- 仅返回 `deleted_at IS NULL` 的记录
- V1 不基于 `status` 做过滤
- 若 `expires_at` 非空，则仅返回 `expires_at > 当前时间` 的记录
- 命中返回 `302`
- 未命中返回 `404`

## 4. 数据模型

### 4.1 主表 `short_links`

| 字段 | 类型 | 约束/说明 |
| --- | --- | --- |
| `id` | bigint unsigned | 主键，自增 |
| `short_code` | varchar(16) | 短码，唯一 |
| `original_url` | varchar(2048) | 原始长链 |
| `url_hash` | char(64) | `original_url` 去除尾部 `/` 后字符串的哈希值，用于查重 |
| `visit_count` | bigint unsigned | 访问次数，默认 `0` |
| `status` | tinyint unsigned | 状态，默认 `1`，不使用类型零值 |
| `expires_at` | datetime | 过期时间，可空，预留字段 |
| `created_at` | datetime | 创建时间 |
| `updated_at` | datetime | 更新时间 |
| `deleted_at` | datetime | 可空，软删除预留字段 |

### 4.2 索引

- `PRIMARY KEY (id)`
- `UNIQUE KEY uk_short_code (short_code)`
- `KEY idx_url_hash (url_hash)`

### 4.3 说明

- V1 只使用主表
- 相同长链以 `original_url` 去除尾部 `/` 后字符串一致为准，重复创建时基于 `url_hash` 复用已有记录
- `url_hash` 由 `original_url` 去除尾部 `/` 后直接计算得到；除尾部 `/` 外，不做大小写、参数顺序等规范化处理
- `url_hash` 仅作为查重辅助键，不作为最终唯一真值
- 创建短链时，先按 `url_hash` 查询候选记录，再对比去除尾部 `/` 后的长链字符串；字符串一致才复用
- 若 `url_hash` 相同但长链字符串不同，视为哈希冲突，继续生成新的 `short_code` 并写入新记录
- `expires_at` 对应短链过期时间字段；创建时由 `Short.ExpairedDays` 决定是否写入，默认配置 `0` 表示永久有效
- `status` 不使用类型零值，V1 默认写入 `1`；预留枚举约定为 `1=active`、`2=disabled`，V1 不提供状态流转与状态控制接口
- `visit_count` 默认 `0`，用于累计成功跳转次数，但 V1 不提供统计报表能力
- 保留 `deleted_at` 作为软删除扩展位，V1 不提供删除接口
- 短链查询与跳转只处理 `deleted_at IS NULL` 的记录
- V1 不引入 `pv`、`uv`、访问明细表和统计聚合表

## 5. 接口定义

### 5.1 `POST /api/v1/links`

请求体：

```json
{
  "long_url": "https://xiaolincoding.com/other/offer.html"
}
```

响应体：

```json
{
  "short_code": "Jxts",
  "short_url": "http://127.0.0.1:8888/Jxts",
  "original_url": "https://xiaolincoding.com/other/offer.html"
}
```

### 5.2 `GET /:code`

行为：

- 命中返回 `302 Found`
- 响应头 `Location: <original_url>`
- 未命中返回 `404 Not Found`

说明：

- V1 固定使用 `302`，不使用 `301`
- `301` 可能被浏览器或中间层长期缓存，后续访问不再到达短链服务
- 使用 `302` 可以降低永久缓存带来的统计失真风险，并为后续目标地址调整、状态控制和风控预留空间

### 5.3 错误码

| 状态码 | 含义 |
| --- | --- |
| `400` | 参数非法 |
| `404` | 短码不存在 |
| `500` | 服务内部错误 |

## 6. 约束与验收

### 6.1 处理约束

- `long_url` 为空返回 `400`
- `long_url` 非法返回 `400`
- `url_hash` 命中后必须二次比对长链字符串，避免哈希冲突导致误复用
- 短码冲突时重新生成并重试写入
- 成功跳转后 `visit_count` 累加 `1`
- 已软删除记录不参与短链解析
- 已过期记录不参与短链解析
- 数据库异常返回 `500`

### 6.2 验收标准

- 可成功创建短链
- 相同长链重复创建返回同一短链
- `url_hash` 相同但长链字符串不同的情况下，不复用旧短链
- 访问有效短链返回 `302`，且 `Location` 正确
- 有效短链成功访问后 `visit_count` 增加
- 访问不存在短链返回 `404`
- 已软删除短链访问返回 `404`
- 已过期短链访问返回 `404`
- 服务重启后已创建短链仍可解析

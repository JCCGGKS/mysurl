# 短链系统 V6 认证刷新技术方案

## 1. 目标

V6 为短链系统补充 `access token + refresh token` 认证刷新能力，解决当前单 JWT 模式“过期后只能重新登录、无法主动吊销”的问题。

本版目标：

- 登录/注册后签发短期 `access token`
- 同时签发长期 `refresh token`
- 前端在 `access token` 失效后自动刷新
- 登出时服务端吊销当前 `refresh token`
- 使用 refresh token rotation 提升长期凭证安全性

## 2. 关键设计

### 2.1 双 token 模式

- `access token`
  - 使用 HS256 JWT
  - 仅用于访问业务接口
  - 默认有效期 `900` 秒
- `refresh token`
  - 使用高强度随机字符串
  - 仅用于刷新和登出
  - 默认有效期 `604800` 秒
  - 服务端持久化并支持吊销

核心原则：

- access token 保持无状态校验
- refresh token 改为有状态管理
- 刷新成功后旧 refresh token 立即失效
- 业务接口只接受 access token

### 2.2 配置调整

`Auth.ExpireSeconds` 拆分为：

- `Auth.AccessExpireSeconds`
- `Auth.RefreshExpireSeconds`

保留：

- `Auth.JWTSecret`
- `Auth.PasswordPepper`

### 2.3 数据存储

新增 `user_refresh_tokens` 表，至少包含：

- `id`
- `user_id`
- `token_hash`
- `expires_at`
- `revoked_at`
- `created_at`
- `updated_at`

要求：

- 数据库只保存 `token_hash`，不保存明文 refresh token
- `revoked_at IS NULL` 且 `expires_at > now()` 才算有效

建议索引：

- `uk_token_hash (token_hash)`
- `idx_user_id (user_id)`
- `idx_user_id_revoked_at (user_id, revoked_at)`

## 3. 接口设计

### 3.1 登录与注册响应

原响应：

- `token`
- `expires_at`
- `user`

升级为：

- `access_token`
- `access_expires_at`
- `refresh_token`
- `refresh_expires_at`
- `user`

### 3.2 刷新接口

- `POST /api/v1/auth/refresh`

请求参数：

- `refresh_token`

返回结构：

- `access_token`
- `access_expires_at`
- `refresh_token`
- `refresh_expires_at`

行为要求：

- 校验 refresh token 是否存在、未吊销、未过期
- 成功后吊销旧 refresh token
- 返回新的 access token 和 refresh token

### 3.3 登出接口

- `POST /api/v1/auth/logout`

请求参数：

- `refresh_token`

行为要求：

- 按 refresh token 吊销当前长期凭证
- 对不存在、已过期或已吊销的 refresh token 直接返回成功

## 4. 实现约束

### 4.1 后端

- access token 继续沿用当前 JWT 方案，claims 至少包含：
  - `user_id`
  - `username`
  - `iat`
  - `exp`
- 鉴权中间件继续只校验 access token
- refresh token 使用随机字符串生成，哈希后写库
- 登录/注册成功后同时签发双 token
- refresh 流程必须使用 rotation：
  1. 校验旧 refresh token
  2. 吊销旧 refresh token
  3. 生成新 refresh token
  4. 写入新 refresh token
  5. 返回新的 access token 和 refresh token
- 第 2 到第 4 步必须放在同一事务中

错误语义：

- access token 缺失：`authorization token is required`
- access token 非法或过期：`authorization token is invalid`
- refresh token 缺失：`refresh token is required`
- refresh token 非法、过期或已吊销：`refresh token is invalid`

### 4.2 前端

本地认证信息改为：

- `access_token`
- `refresh_token`
- `user`

请求规则：

- 业务接口统一携带 `Authorization: Bearer <access_token>`
- 收到 `401`` 时先尝试 refresh，再决定是否跳登录页

自动刷新流程：

1. 若无 refresh token，清理本地状态并跳登录页
2. 若有 refresh token，调用 `/api/v1/auth/refresh`
3. 刷新成功后更新本地双 token，并重试原请求一次
4. 刷新失败后清理本地状态并跳登录页

并发约束：

- 前端必须实现 refresh 单飞
- 同一时间只允许一个 refresh 请求执行
- 原业务请求最多重试一次
- refresh 接口自身失败时不递归重试

## 5. 验收

- 登录成功后返回 `access_token` 和 `refresh_token`
- 注册成功后返回 `access_token` 和 `refresh_token`
- 业务接口只接受 access token
- access token 过期后，前端可自动刷新并继续访问业务接口
- refresh 成功后旧 refresh token 立即失效
- 已吊销 refresh token 不能再次用于刷新
- logout 后原 refresh token 不可再次使用
- refresh token 在数据库中只保存 hash，不保存明文
- 多个请求同时触发 `401` 时，前端仅发起一次 refresh 请求

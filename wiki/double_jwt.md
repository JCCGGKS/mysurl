# 双 Token 机制整理

## 1. 先说结论

双 token 机制通常指：

- `access token`：短期有效，负责访问业务接口
- `refresh token`：长期有效，负责换新 access token

一句话理解：

`access token` 负责“当前能不能访问接口”，`refresh token` 负责“还能不能继续续期”。

## 2. 为什么需要双 token

如果系统只使用一个长效 JWT，通常会有这些问题：

- token 泄露后，在过期前都可能持续可用
- 无法方便地主动吊销
- 退出登录往往只是前端删本地状态
- token 过期后只能重新登录

双 token 的核心思路是：

- 把“访问接口”和“续期登录态”拆开
- 让高频使用的凭证短效化
- 让长期凭证可控、可吊销

所以：

- access token 追求轻量和无状态
- refresh token 追求可撤销和可管理

## 3. 两类 token 的职责

### 3.1 access token

特点：

- 生命周期短
- 请求频率高
- 一般通过 `Authorization: Bearer <token>` 传递
- 用于访问受保护业务接口

常见做法：

- 使用 JWT
- claims 中放 `user_id`、`username`、`exp`、`iat`

它的优点是无状态，缺点是通常不能方便地即时撤销，所以有效期应较短。

### 3.2 refresh token

特点：

- 生命周期长
- 使用频率低
- 不参与普通业务接口访问
- 只用于 refresh 和 logout

推荐做法：

- 使用高强度随机字符串
- 不建议直接复用 JWT
- 服务端持久化存储并支持吊销

它比 access token 更敏感，因为它可以换出新的 access token，本质上是长期续命凭证。

## 4. 典型流程

### 4.1 登录

1. 用户登录成功
2. 服务端签发 access token
3. 服务端生成 refresh token
4. 前端保存两者

### 4.2 正常访问

1. 前端只带 access token
2. 服务端校验 JWT
3. 校验通过后进入业务逻辑

### 4.3 access token 过期

1. 业务接口返回 `401`
2. 前端使用 refresh token 调用刷新接口
3. 服务端校验 refresh token
4. 校验通过后签发新的 access token
5. 前端更新本地 token
6. 重试原业务请求

### 4.4 登出

1. 前端调用 logout
2. 服务端吊销当前 refresh token
3. 前端清理本地状态

注意：

- logout 主要吊销 refresh token
- 已签发但未过期的 access token 一般依赖短生命周期自然失效

## 5. 为什么 refresh token 要持久化

如果 refresh token 不落库，服务端通常做不到：

- 主动吊销
- 退出登录后立即失效
- 多端会话管理
- refresh token rotation
- 安全审计

持久化之后，服务端才能判断：

- 这个 refresh token 是否存在
- 是否已吊销
- 是否已过期
- 是否还能继续换新 token

所以：

- access token 更适合无状态
- refresh token 更适合有状态

## 6. 为什么数据库里只存 hash

refresh token 是高价值凭证。

如果数据库直接保存明文，一旦库泄露，攻击者可以直接拿它去刷新新的 access token。

更稳妥的方式是：

1. 服务端生成 refresh token 明文
2. 返回给前端保存
3. 服务端只保存它的 hash

后续刷新时：

1. 前端提交明文 refresh token
2. 服务端做同样 hash
3. 用 hash 去数据库匹配

这样即使数据库泄露，也拿不到可直接使用的明文 token。

## 7. 什么是 refresh token rotation

rotation 的意思是：

- 每次 refresh 成功
- 旧 refresh token 立刻作废
- 服务端返回一个新的 refresh token


典型流程：

1. 校验旧 refresh token
2. 吊销旧 token
3. 生成新 token
4. 返回新 token

这样可以降低长期凭证重复使用和泄露后的风险。

## 8. 前端要注意什么

- 普通业务请求只带 access token
- refresh token 只用于 `/auth/refresh` 和 `/auth/logout`
- 业务接口返回 `401` 后，应先尝试 refresh，再决定是否跳登录页
- refresh 成功后重试原请求一次
- refresh 失败后清理登录态并跳登录页

还要做 refresh 单飞：

- 同一时间只允许一个 refresh 请求执行
- 其他因为 `401` 失败的请求等待该结果

否则多个 refresh 并发时，旧 refresh token 可能被重复使用而全部失败。

## 9. 当前项目怎么落

按我们讨论的设计，当前项目适合这样做：

- `access token`
  - 继续使用 JWT
  - 放在 `Authorization` 头里
  - 生命周期短，例如 15 分钟

- `refresh token`
  - 使用随机字符串
  - 服务端持久化到 MySQL
  - 数据库存 hash
  - 生命周期长，例如 7 天

- 服务端新增接口：
  - `POST /api/v1/auth/refresh`
  - `POST /api/v1/auth/logout`

- 前端本地保存：
  - `access_token`
  - `refresh_token`
  - `user`
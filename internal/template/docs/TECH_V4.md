# 短链系统 V4 技术方案

## 1. 目标

V4 为短链系统补充基础用户体系，支持：

- 注册
- 登录
- JWT 鉴权
- 登录后创建短链
- 短链记录创建者归属

本版不修改短链接口路径。前端页面要求登录后才能进入创建流程；后端在创建请求携带有效 token 时记录创建者归属。

## 2. 数据变更

### 2.1 users 表

新增 `users` 表，至少包含：

- `id`
- `username`
- `password_hash`
- `created_at`
- `updated_at`
- `deleted_at`

约束：

- `username` 唯一
- 密码只存哈希，默认使用 `bcrypt`

### 2.2 short_links 表

新增字段：

- `user_id bigint unsigned DEFAULT NULL`

说明：

- 创建请求带有效 token 时写入当前用户 `id`
- 历史数据允许 `user_id = NULL`
- 跳转逻辑不受影响

## 3. 接口设计

### 3.1 注册

- `POST /api/v1/auth/register`

请求：

- `username`
- `password`
- `confirm_password`

规则：

- 用户名非空，限制长度
- 只允许字母、数字、下划线
- 密码最小长度默认 `8`
- 两次密码必须一致

行为：

- 校验参数
- 检查用户名是否存在
- 创建用户
- 直接返回 token 和用户信息

### 3.2 登录

- `POST /api/v1/auth/login`

请求：

- `username`
- `password`

行为：

- 校验用户名和密码
- 返回 token 和用户信息

失败返回 `401`。

### 3.3 创建短链

- `POST /api/v1/links`

说明：

- 接口路径和返回结构保持不变
- 请求可继续匿名调用
- 如果请求携带有效 Bearer token，则解析当前用户并写入 `short_links.user_id`
- 如果没有 token，则按匿名请求处理，`user_id` 为空

### 3.4 跳转短链

- `GET /:code`

说明：

- 路径和行为保持不变
- 继续匿名可访问

## 4. 鉴权方案

- 使用 JWT Bearer
- 请求头：`Authorization: Bearer <token>`
- token 至少包含：
  - `user_id`
  - `username`
  - `exp`

新增配置：

- `Auth.JWTSecret`
- `Auth.ExpireSeconds`

鉴权中间件负责：

- 读取和校验 token
- 将当前用户写入上下文
- 用于注册、登录外的受保护接口或创建接口中的“可选用户识别”

## 5. 前端改造

Vue 端新增：

- 注册页
- 登录页
- 登录后创建页

前端规则：

- 登录或注册成功后保存 token
- 应用启动时按本地 token 判断登录态
- 创建请求自动携带 Bearer token
- token 失效时清理本地登录态并跳回登录页

保留当前创建页已有能力：

- 长链输入
- URL 校验
- 结果展示
- 复制短链
- 打开短链

## 6. 边界

本版保持当前全局复用语义：

- 相同规范化长链仍复用同一条短链
- 不按用户拆分独立短链

本版不做：

- refresh token
- 登出黑名单
- 找回密码
- 邮箱或短信验证
- 角色权限
- 我的短链列表

## 7. 验收

- 注册成功并返回 token
- 登录成功并返回 token
- 重复用户名注册失败
- 非法 token 登录后创建时不会写入 `user_id`
- 已登录创建短链成功并写入 `user_id`
- 匿名创建仍可成功，但 `user_id = NULL`
- 历史 `user_id = NULL` 的短链仍可正常跳转

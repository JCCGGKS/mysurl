# 短链前端认证方案

## 目标

在现有创建页基础上补充前端认证流，支持：

- 登录页
- 注册页
- 登录后进入创建页
- 创建请求自动携带 token
- 创建页显示最小用户区

## 页面结构

前端新增三条路由：

- `/login`
- `/register`
- `/create`

默认行为：

- 有 token 时进入 `/create`
- 无 token 时进入 `/login`

## 登录态

- 使用 `localStorage` 保存 token
- 同时保存最小用户信息，至少包含 `username`
- 访问 `/create` 时检查本地 token
- 没有 token 时跳转 `/login`
- 点击退出时清空本地 token 和用户信息

## 接口联动

### 注册

- `POST /api/v1/auth/register`

请求：

```json
{
  "username": "demo_user",
  "password": "password123",
  "confirm_password": "password123"
}
```

行为：

- 注册成功后跳转 `/login`
- 不自动登录
- 给出“注册成功，请登录”提示

### 登录

- `POST /api/v1/auth/login`

请求：

```json
{
  "username": "demo_user",
  "password": "password123"
}
```

行为：

- 登录成功后保存 token 和用户信息
- 跳转 `/create`

### 创建短链

- `POST /api/v1/links`

请求头：

```http
Authorization: Bearer <token>
```

请求体：

```json
{
  "long_url": "https://example.com/article/123"
}
```

## 创建页行为

保留当前能力：

- 输入长链
- URL 校验
- 展示 `short_url`、`short_code`、`original_url`
- 复制短链
- 打开短链

新增能力：

- 顶部显示当前用户名
- 提供退出按钮
- 请求自动携带 token

## 异常处理

- 登录失败：展示后端错误
- 注册失败：展示后端错误
- 创建接口返回 `401`：清空本地登录态并跳回 `/login`
- 网络失败：展示统一错误提示

## 实现建议

- 引入 `vue-router`
- `App.vue` 改为路由容器
- 登录页、注册页、创建页拆成独立视图组件
- 使用轻量认证工具模块统一处理：
  - token 读写
  - 用户信息读写
  - logout
  - 创建请求头注入

## 验收标准

- 无 token 不能进入 `/create`
- 注册成功后跳转 `/login`
- 登录成功后跳转 `/create`
- 创建请求能自动带上 Bearer token
- 创建页显示用户名和退出按钮
- 退出后回到登录页
- `401` 时能清理登录态并重新跳转登录页

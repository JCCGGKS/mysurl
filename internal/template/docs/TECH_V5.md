# 短链系统 V5 技术方案

## 1. 目标

V5 为短链系统补充用户操作日志能力，支持：

- 记录登录成功日志
- 记录创建短链成功日志
- 前端展示当前用户的操作日志列表
- 使用操作日志中间件统一落日志

本版中间件放在鉴权中间件之前，但只记录成功业务动作，不记录未授权请求到用户操作日志表。

## 2. 数据变更

### 2.1 user_operation_logs 表

新增 `user_operation_logs` 表，至少包含：

- `id`
- `user_id`
- `action`
- `result`
- `target_id`
- `target_code`
- `created_at`

字段说明：

- `action` 第一版只支持：
  - `login`
  - `create_link`
- `result` 第一版只记录：
  - `success`
- `target_id` 记录被操作对象的数据库主键
- `target_code` 记录被操作对象的业务标识，例如短码
- 登录成功日志允许 `target_id`、`target_code` 为空

建议索引：

- `idx_user_id_id (user_id, id)`
- `idx_user_id_action_id (user_id, action, id)`

本版日志默认永久保留，不做清理或归档。

## 3. 中间件方案

### 3.1 操作日志中间件

新增操作日志中间件，职责：

- 作为最外层包装请求
- 在请求返回后读取最终响应状态
- 读取上下文中的操作日志摘要
- 满足条件时写入 `user_operation_logs`

路由链路顺序：

- `operationLogMiddleware`
- `authMiddleware`
- `handler`

说明：

- 操作日志中间件虽然放在鉴权之前，但不直接记录所有请求
- 缺 token、坏 token、401 请求不写用户操作日志表
- 中间件本身不解析业务响应 JSON
- 日志写入失败不影响主业务成功返回

### 3.2 日志写入条件

第一版仅在以下条件同时满足时写入日志：

- HTTP 响应成功
- 业务链路已写入操作日志摘要

未满足条件时：

- 不写 `user_operation_logs`
- 保留服务日志用于排查

## 4. 上下文摘要传递

为了让前置中间件拿到业务结果，业务链路需要向上下文写入操作日志摘要。

### 4.1 登录成功摘要

登录成功后写入：

- `user_id`
- `action=login`
- `result=success`

### 4.2 创建短链成功摘要

创建短链成功后写入：

- `user_id`
- `action=create_link`
- `result=success`
- `target_id=short_links.id`
- `target_code=short_links.short_code`

说明：

- 创建新短链成功时记录
- 复用已有短链成功时也记录

## 5. 接口设计

### 5.1 用户操作日志列表

- `GET /api/v1/user-operation-logs`

请求参数：

- `last_id`
- `limit`

说明：

- 按当前登录用户查询
- 使用基于 `id` 的游标分页
- 第一版不做动作筛选

返回结构：

- `items`
- `total`
- `limit`
- `has_more`
- `next_last_id`

每条日志项至少包含：

- `id`
- `action`
- `result`
- `target_id`
- `target_code`
- `created_at`

## 6. 前端改造

前端“用户操作日志”页面改为真实列表页，展示当前登录用户的日志记录。

页面内容：

- 页面标题
- 日志表格
- 分页栏

表格列：

- `ID`
- `时间`
- `动作`
- `结果`
- `对象ID`
- `短码`

前端文案映射：

- `login` -> `登录`
- `create_link` -> `创建短链`
- `success` -> `成功`

本版不做：

- 动作筛选
- 关键字搜索
- 日志详情弹窗

## 7. 边界

本版不记录：

- 登录失败
- 注册成功
- 短链列表查询
- 短链跳转
- 未授权访问

本版不做：

- 日志清理
- TTL
- 归档
- 审计报表

## 8. 验收

- 登录成功后写入一条 `login/success` 日志
- 创建短链成功后写入一条 `create_link/success` 日志
- 创建短链复用已有短链时也写入日志
- 缺 token、坏 token、401 请求不写用户操作日志表
- 日志写入失败不影响主业务成功
- `GET /api/v1/user-operation-logs` 支持分页查询
- 前端页面可展示真实日志记录

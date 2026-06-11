# 操作日志重构方案

## 目标

- 业务层不再显式写操作日志上下文
- 操作日志尽量在中间件统一收口
- 成功/失败结果由统一响应推导
- 创建短链成功时补充命中来源，例如 `cache_hit`、`database_hit`、`created`

## 现状方案

### 1. 路由匹配

操作日志中间件维护一份 `method + path -> process` 的映射表。  
当前已接入：

- `POST /api/v1/auth/login`
- `POST /api/v1/links`
- `POST /api/v1/links/batch`

`process` 只负责两件事：

- 声明操作类型 `action`
- 提供成功/失败时的 `reason` 提取逻辑

### 2. 中间件职责

操作日志中间件统一完成以下工作：

- 根据 `method + path` 找到对应 `process`
- 包装 `ResponseWriter`，捕获最终响应体
- 解析统一响应结构
- 根据响应 `code` 推导 `result`
  - `code == 0` 记为 `success`
  - 其他情况记为 `failed`
- 提取当前用户 `user_id`
- 组装并落库操作日志

中间件统一组装的字段：

- `user_id`
- `action`
- `result`
- `reason`

### 3. 用户 ID 获取

鉴权接口继续使用 `context` 传递 claims。  
因此受保护接口的中间件顺序必须是：

- `auth -> operationLog -> handler`

这样操作日志中间件才能从 `request context` 中拿到 `user_id`。  
登录接口没有鉴权信息，登录成功时从响应 `data.user.id` 兜底提取 `user_id`。

### 4. 统一响应扩展字段

统一响应结构 `utils.Response` 增加：

- `extdata`

该字段用于传递仅后置链路关心的扩展信息。  
当前主要用于创建短链成功场景，记录成功来源：

- `cache_hit`
- `database_hit`
- `created`

### 5. 创建短链成功来源

`CreateLinkLogic` 在成功返回时额外产出来源信息：

- 命中缓存：`cache_hit`
- 命中数据库：`database_hit`
- 新建成功：`created`

`CreateLinkHandler` 在写统一成功响应时，将该值放入响应 `extdata`。  
操作日志中间件从响应体解析 `extdata`，并写入操作日志 `reason` 字段。

## 失败原因规则

- 失败场景：优先使用统一响应 `msg` 作为 `reason`
- 成功场景：默认空字符串
- 创建短链成功：使用响应 `extdata` 作为 `reason`

## 当前收益

- 去掉了业务代码中的日志埋点分支
- 操作日志写入逻辑集中在中间件
- 成功/失败判定统一依赖响应结果
- 创建短链可区分缓存命中、数据库命中和新建成功

## 后续扩展方式

新增接口接入操作日志时，按以下步骤处理：

1. 在操作日志中间件的 `process map` 中注册 `method + path`
2. 配置该接口的 `action`
3. 如有特殊成功原因，从统一响应 `extdata` 或响应内容中提取
4. 如有特殊失败原因，从统一响应 `msg` 或响应内容中提取

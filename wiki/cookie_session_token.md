# Cookie / Session / Token 简要整理

## 一、先记结论

- `Cookie`：浏览器侧的存储与自动携带机制
- `Session`：服务端保存的会话状态
- `Token`：客户端持有的身份凭证

三者不是同一层概念：

- `Cookie` 更像“载体”
- `Session` 和 `Token` 更像“登录态方案”

一句话记忆：

`Cookie` 是“怎么带”，`Session` 是“服务端怎么记”，`Token` 是“客户端拿什么证明自己”。

## 二、三者分别是什么

### 1. Cookie

Cookie 是浏览器保存的一小段数据。服务端可以通过 `Set-Cookie` 写入，浏览器后续访问同域名时会自动带上。

它常用来存：

- `session_id`
- `token`
- 语言、主题等轻量配置

常见属性：

- `HttpOnly`
- `Secure`
- `SameSite`
- `Expires` / `Max-Age`

Cookie 本身不等于登录态，它只是传输和存储容器。

### 2. Session

Session 是服务端维护的一段会话数据。典型模式是：

1. 用户登录成功
2. 服务端创建 session
3. 返回一个 `session_id`
4. 浏览器通过 Cookie 自动带上 `session_id`
5. 服务端根据 `session_id` 查出用户会话

特点：

- 用户状态主要在服务端
- 方便主动失效、踢下线、改密码后失效
- 多实例部署时通常要把 session 放到 Redis 等共享存储

### 3. Token

Token 是客户端持有的凭证。典型模式是：

1. 用户登录成功
2. 服务端签发 token
3. 客户端保存 token
4. 请求时主动带上 token
5. 服务端校验 token 合法性和过期时间

常见传输方式：

- `Authorization: Bearer <token>`
- 放在 Cookie 中

常见形式：

- 随机字符串
- JWT

特点：

- 更适合前后端分离和 API
- 更容易做无状态扩容
- 主动失效通常比 session 麻烦

## 三、它们之间的关系

### 1. Cookie + Session

最经典的组合：

- Cookie 里放 `session_id`
- 服务端保存 session 内容

也就是：

- Cookie 是门票号码
- Session 是服务端档案

### 2. Cookie + Token

Token 也可以放进 Cookie。

这时本质上仍然是 token 认证，只是传输载体从 `Authorization` 头换成了 Cookie。

### 3. Session + Token

它们都是登录态方案，但思路不同：

- Session：状态主要在服务端
- Token：状态主要在客户端凭证里

生产里也可能并存，例如：

- 短期 access token
- 长期 refresh token

## 四、当前项目属于哪种

当前项目使用的是：

- JWT token
- 前端本地保存 token
- 请求头 `Authorization: Bearer <token>`
- 服务端中间件解析 token 得到 `user_id`、`username`

所以它不是 session 模式，因为：

- 服务端没有维护登录 session 表
- 登录后不依赖 `session_id` 查会话
- 主要依赖 JWT 本身完成身份校验

当前链路大致是：

1. 登录成功
2. 服务端签发 JWT
3. 前端保存 token
4. 请求业务接口时带上 `Authorization`
5. 后端中间件验签并解析 claims

## 五、怎么选

### 1. Session 更适合

- 传统服务端渲染网站
- 强依赖服务端会话控制
- 需要方便地下线、踢人、强制失效

### 2. Token 更适合

- 前后端分离
- APP / 小程序 / 多端 API
- 微服务或网关鉴权
- 更强调无状态扩容

不要机械认为 JWT 一定更高级：

- Session 更直观
- 主动失效更容易
- 小系统里往往更省心

## 六、安全上最该关注什么

### 1. Cookie

风险：

- 被窃取
- `CSRF`
- 未设置 `HttpOnly` 被脚本读取

常见防护：

- `HttpOnly`
- `Secure`
- `SameSite`

### 2. Session

风险：

- `session_id` 泄漏
- 长时间不失效

常见防护：

- 登录后重新生成 session id
- 设置合理过期时间
- 共享存储统一管理

### 3. Token

风险：

- token 泄漏
- 有效期过长
- 难以快速撤销

常见防护：

- 短有效期
- refresh token 机制
- HTTPS
- 不在 token 中放敏感明文

## 七、压缩成一句

如果只记一句：

`Cookie` 负责带数据，`Session` 负责在服务端存会话，`Token` 负责让客户端证明身份。

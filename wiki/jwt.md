# JWT 学习整理

## 一、JWT 是什么

JWT 全称是 `JSON Web Token`，本质上是一种“可签名的、自包含的令牌格式”。

常见用途：

- 登录后给客户端发访问令牌
- 服务间传递身份信息
- 无状态鉴权

JWT 解决的问题不是“加密存储用户信息”，而是“让接收方能验证这份身份声明是不是可信、有没有过期、有没有被篡改”。

## 二、JWT 的基本结构

一个 JWT 通常由三段组成，中间用 `.` 分隔：

`Header.Payload.Signature`

例如：

`xxxxx.yyyyy.zzzzz`

### 1. Header

Header 描述令牌元信息，最常见的是：

- `typ`: 类型，通常是 `JWT`
- `alg`: 签名算法，比如 `HS256`、`RS256`

示例：

```json
{
  "typ": "JWT",
  "alg": "HS256"
}
```

### 2. Payload

Payload 是声明集合，也叫 `claims`。

它通常包含：

- 业务字段：如 `user_id`、`username`、`role`
- 标准字段：如 `exp`、`iat`、`sub`

示例：

```json
{
  "user_id": 123,
  "username": "alice",
  "sub": "alice",
  "iat": 1710000000,
  "exp": 1710086400
}
```

### 3. Signature

签名部分用来证明：

- Header 和 Payload 没被篡改
- Token 的确由持有密钥的一方签发

如果使用 `HS256`，签名依赖一个对称密钥 `secret`。

## 三、JWT 不是加密格式

这是最容易混淆的一点。

普通 JWT 默认只是：

1. 把 Header 和 Payload 做 Base64URL 编码
2. 再做签名

所以：

- JWT 内容通常可以被解码查看
- 不能把密码、银行卡号、身份证号等敏感信息直接塞进 Payload
- JWT 的安全性主要来自“签名校验”，不是“内容保密”

如果需要“既能验证，又能加密内容”，那已经不是普通 JWT 的使用范畴，通常会进入 `JWE/JOSE` 体系。

## 四、常见标准 Claims

JWT 里最常见的标准字段如下：

| 字段 | 含义 | 典型用途 |
| --- | --- | --- |
| `sub` | Subject | 标识 token 主题，通常是用户名或用户 ID |
| `exp` | Expiration Time | 过期时间，超过后 token 无效 |
| `iat` | Issued At | 签发时间 |
| `nbf` | Not Before | 生效时间，早于该时间不可用 |
| `iss` | Issuer | 签发者 |
| `aud` | Audience | 目标受众 |
| `jti` | JWT ID | token 唯一编号，常用于撤销或去重 |

要点：

- `exp` 几乎是生产必备
- `iat` 便于排查问题
- `iss`、`aud` 在多系统对接时更重要
- `jti` 适合做黑名单、单点踢出、刷新令牌控制

## 五、常见签名算法

### 1. HS256 / HS384 / HS512

HMAC 对称签名算法。

特点：

- 签发和校验都使用同一个 secret
- 实现简单，性能好
- 适合单体服务、内部服务、自主签发场景

风险：

- 任何能校验 token 的服务，如果拿到 secret，也能签发 token
- secret 泄漏后风险很大

### 2. RS256 / RS512

RSA 非对称签名算法。

特点：

- 私钥签发
- 公钥校验
- 更适合多服务、多消费方场景

适合：

- 单独认证中心签发
- 多个服务只负责验签
- 对接第三方身份系统

### 3. ES256

ECDSA 非对称算法。

特点：

- 同属非对称签名
- 密钥更短
- 实现和兼容性通常比 HMAC 复杂一些

## 六、JWT 的典型工作流程

登录型系统通常是：

1. 用户提交用户名和密码
2. 服务端校验账号密码
3. 服务端生成 JWT
4. 客户端保存 token
5. 后续请求通过 `Authorization: Bearer <token>` 携带 token
6. 服务端验签并解析 claims
7. 服务端根据 claims 识别用户身份

它的核心思想是：

- 登录态从“查 session”转为“校验 token”
- 减少对服务端内存会话存储的依赖
- 用 token 自带信息表达身份

## 七、JWT 的优点和限制

### 优点

- 无状态，横向扩容简单
- 跨服务、跨语言传递方便
- 可以在 token 中携带必要身份信息
- 很适合前后端分离和 API 鉴权

### 限制

- 发出去后天然难以主动撤销
- token 一旦泄漏，在过期前可能持续可用
- Payload 可读，不适合放敏感明文
- 如果 claims 过多，请求头会变大

## 八、生产实践里最容易踩的坑

### 1. 把 JWT 当加密容器

错误理解。JWT 默认不是加密格式。

### 2. secret 太短或手写

例如：

- `123456`
- `admin`
- `change-me`

这类都不合格。

### 3. 不校验过期时间

不校验 `exp`，token 就可能长期有效。

### 4. token 有效期太长

有效期太长会放大泄漏风险。

### 5. 完全不考虑撤销机制

如果有“退出登录、密码修改后失效、管理员踢下线”这类需求，单靠 JWT 本身不够，需要额外机制。

## 九、JWT 撤销与续期的常见策略

JWT 常被说成“无状态”，但生产里通常还是会配合状态化策略。

常见做法：

### 1. 短生命周期 Access Token

例如：

- Access Token：15 分钟
- Refresh Token：7 天

这是最常见的折中方案。

### 2. Refresh Token 续期

Access Token 过期后，客户端用 Refresh Token 申请新 token。

好处：

- 降低 Access Token 泄漏窗口
- 保持用户体验

### 3. 黑名单

把已撤销的 `jti` 或 token 指纹放到 Redis 之类存储中，校验时额外判断。

适合：

- 强制退出
- 密码修改后失效
- 管理员封禁

### 4. 版本号策略

给用户表增加 token 版本号或登录版本号，签发时写入 claims，校验时和数据库当前版本比对。

适合：

- 单用户全端下线
- 修改密码后使旧 token 整体失效

## 十、Go 里常见 JWT 相关库

### 1. `github.com/golang-jwt/jwt`

最常见的纯 JWT 库，适合：

- 自己签发 token
- 自己验签
- `HS256/RS256/ES256` 等常见场景

特点：

- 社区最常见
- API 直接
- 上手成本低

说明：

- 老项目里常见的 `dgrijalva/jwt-go` 已经迁移到这里

### 2. `github.com/go-jose/go-jose`

适合 JOSE 全家桶场景，包括：

- JWT
- JWS
- JWE

适合更复杂的认证集成，而不只是简单登录 token。

### 3. `github.com/lestrrat-go/jwx`

适合：

- JWT
- JWK/JWKS
- OIDC
- 第三方签发 token 校验

如果你要对接外部身份平台，通常会比纯 `jwt` 库更顺手。

## 十一、当前项目 `mysurl1` 的 JWT 实现

### 1. 当前使用的库

项目当前使用：

`github.com/golang-jwt/jwt/v4`

代码位置在 [internal/utils/auth.go](/home/fanqicheng/project/jx/mysurl1/internal/utils/auth.go:11)。

### 2. 当前 claims 结构

项目里定义的是：

```go
type AuthClaims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
```

含义：

- 自定义业务字段：`user_id`、`username`
- 标准 claims：通过 `jwt.RegisteredClaims` 承载

### 3. 当前签名算法

当前项目在 [internal/utils/auth.go](/home/fanqicheng/project/jx/mysurl1/internal/utils/auth.go:113) 使用：

```go
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
```

说明当前采用的是：

- `HS256`
- 对称密钥签名
- 使用配置中的 `Auth.JWTSecret`

### 4. 当前配置项

配置定义在 [internal/config/config.go](/home/fanqicheng/project/jx/mysurl1/internal/config/config.go:42)：

```go
type AuthConf struct {
	JWTSecret      string `json:",optional"`
	ExpireSeconds  int64  `json:",default=86400"`
	PasswordPepper string `json:",optional"`
}
```

对应 YAML 在 [etc/mysurl1-api.yaml](/home/fanqicheng/project/jx/mysurl1/etc/mysurl1-api.yaml:26)：

```yaml
Auth:
  JWTSecret: change-me-in-production
  ExpireSeconds: 86400
  PasswordPepper: ""
```

### 5. 当前过期策略

项目默认：

- `ExpireSeconds = 86400`
- 即默认 24 小时过期

生成 token 时会自动补齐：

- `ExpiresAt`
- `IssuedAt`
- `Subject`

逻辑在 [internal/utils/auth.go](/home/fanqicheng/project/jx/mysurl1/internal/utils/auth.go:57)。

## 十二、适合当前项目的 `JWTSecret` 生成方式

因为当前项目使用的是 `HS256`，所以最适合的是生成一个高熵随机对称密钥。

推荐命令：

```bash
openssl rand -hex 32
```

这会生成：

- 32 字节随机值
- 64 位十六进制字符串输出

为什么适合当前项目：

- 强度足够支撑 `HS256`
- 只有十六进制字符，直接写 YAML 最稳妥
- 不会出现空格、换行、`+`、`/` 等转义和解析问题

示例：

```yaml
Auth:
  JWTSecret: "8f2d7d2a4d2e0f9b0c1a6c8d4e5f7a9b1c3d5e7f9a0b2c4d6e8f1a3b5c7d9e1"
  ExpireSeconds: 86400
  PasswordPepper: ""
```

## 十三、当前项目的改进建议

### 1. 不要把生产 secret 提交到仓库

当前示例配置里的：

`change-me-in-production`

只能用于占位，不能用于真实环境。

### 2. 优先从环境变量注入 secret

比把真实密钥写死在配置文件里更稳妥。

### 3. 后续可补充 token 校验中间件

当前项目已经有签发逻辑，但如果要形成完整鉴权闭环，还需要：

- Bearer Token 解析
- 签名校验
- claims 注入上下文
- 路由级鉴权

### 4. 如有退出登录需求，需要补撤销机制

如果只是登录后发 token，短期可用；但要做完整账户体系，通常还需要：

- 黑名单
- token 版本号
- refresh token

## 十四、结论

- JWT 的核心价值是“可验证的身份声明”
- 普通 JWT 默认不是加密容器
- 当前项目使用的是 `HS256 + 对称密钥`
- 对这个项目最合适的 secret 生成方式是 `openssl rand -hex 32`
- 如果系统后续要做更完整鉴权，应继续补充验签中间件、撤销策略和密钥管理方案

# JWT 学习整理

## 1. JWT 是什么

JWT 全称是 `JSON Web Token`，本质上是一种可签名的令牌格式。

常见用途：

- 登录后发访问令牌
- 服务间传递身份信息
- API 无状态鉴权

JWT 解决的是：

- 这份身份声明是不是可信
- 有没有过期
- 有没有被篡改

JWT 不是用来加密存储敏感信息的。

## 2. JWT 结构

JWT 由三段组成：

`Header.Payload.Signature`

### Header

描述令牌元信息，例如：

- `typ`: `JWT`
- `alg`: `HS256`、`RS256`

### Payload

存放声明 `claims`，常见字段：

- 业务字段：`user_id`、`username`
- 标准字段：`exp`、`iat`、`sub`

### Signature

签名部分用来证明：

- Header 和 Payload 没被篡改
- Token 确实由签发方生成

## 3. JWT 不是加密

普通 JWT 默认只是：

1. Header 和 Payload 做 Base64URL 编码
2. 再做签名

所以：

- 内容通常可被解码查看
- 不要放密码、身份证号、银行卡号等敏感明文
- 安全性主要来自“签名校验”，不是“内容保密”

## 4. 常见 Claims

| 字段 | 含义 |
| --- | --- |
| `sub` | 主题 |
| `exp` | 过期时间 |
| `iat` | 签发时间 |
| `nbf` | 生效时间 |
| `iss` | 签发者 |
| `aud` | 受众 |
| `jti` | 唯一编号 |

最常用的是：

- `exp`
- `iat`
- 业务身份字段

## 5. JWT 全部常见签名算法

JWT 常见算法可以分为 4 类。

### 5.1 HMAC 对称签名

| 算法 | 说明 | 服务端密钥配置 |
| --- | --- | --- |
| `HS256` | HMAC + SHA-256 | 一个共享 `secret` |
| `HS384` | HMAC + SHA-384 | 一个共享 `secret` |
| `HS512` | HMAC + SHA-512 | 一个共享 `secret` |

特点：

- 签发和验签都使用同一个密钥
- 实现简单，性能好
- 适合单体服务、内部服务、自签发场景

注意：

- 能验签的一方，只要拿到 secret，也能签发 token
- 所以 secret 泄漏风险很高

### 5.2 RSA 非对称签名

| 算法 | 说明 | 服务端密钥配置 |
| --- | --- | --- |
| `RS256` | RSA + SHA-256 | 私钥签发，公钥验签 |
| `RS384` | RSA + SHA-384 | 私钥签发，公钥验签 |
| `RS512` | RSA + SHA-512 | 私钥签发，公钥验签 |

特点：

- 认证中心可独立持有私钥
- 业务服务只持有公钥也能验签
- 适合多服务、多消费方场景

### 5.3 RSA-PSS 非对称签名

| 算法 | 说明 | 服务端密钥配置 |
| --- | --- | --- |
| `PS256` | RSA-PSS + SHA-256 | 私钥签发，公钥验签 |
| `PS384` | RSA-PSS + SHA-384 | 私钥签发，公钥验签 |
| `PS512` | RSA-PSS + SHA-512 | 私钥签发，公钥验签 |

特点：

- 也是 RSA 非对称方案
- 比传统 `RS*` 在签名填充方式上更现代

### 5.4 ECDSA 非对称签名

| 算法 | 说明 | 服务端密钥配置 |
| --- | --- | --- |
| `ES256` | ECDSA P-256 + SHA-256 | 私钥签发，公钥验签 |
| `ES384` | ECDSA P-384 + SHA-384 | 私钥签发，公钥验签 |
| `ES512` | ECDSA P-521 + SHA-512 | 私钥签发，公钥验签 |

特点：

- 同样是非对称签名
- 密钥通常更短
- 配置和兼容性一般比 HMAC 复杂

### 5.5 EdDSA

| 算法 | 说明 | 服务端密钥配置 |
| --- | --- | --- |
| `EdDSA` | Edwards 曲线签名 | 私钥签发，公钥验签 |

特点：

- 也是非对称签名
- 现代实现里越来越常见

## 6. 各算法对应的服务端密钥配置思路

### 6.1 HMAC

服务端只需要一个共享密钥：

```yaml
Auth:
  JWTSecret: your-long-random-secret
```

适合当前项目这种：

- 单个 API 服务签发
- 同一个服务负责验签

### 6.2 RSA / RSA-PSS / ECDSA / EdDSA

服务端通常要配置：

- 签发服务：私钥
- 验签服务：公钥

例如可以抽象成：

```yaml
Auth:
  JWTPrivateKeyPath: /path/to/private.pem
  JWTPublicKeyPath: /path/to/public.pem
```

或者：

```yaml
Auth:
  JWTPrivateKey: |
    -----BEGIN PRIVATE KEY-----
    ...
    -----END PRIVATE KEY-----
  JWTPublicKey: |
    -----BEGIN PUBLIC KEY-----
    ...
    -----END PUBLIC KEY-----
```

当前项目还没有做这套配置，当前只实现了 `HS256 + JWTSecret`。

## 7. 典型流程

1. 用户提交用户名和密码
2. 服务端校验成功
3. 服务端生成 JWT
4. 客户端保存 token
5. 后续请求带 `Authorization: Bearer <token>`
6. 服务端验签并解析 claims
7. 根据 claims 识别用户

核心思想：

- 不再靠服务端 session 查登录态
- 而是靠 token 自带身份信息 + 签名校验

## 8. 优点和限制

### 优点

- 无状态，扩容方便
- 适合前后端分离
- 跨服务传递身份方便

### 限制

- 主动撤销麻烦
- token 泄漏后，在过期前仍可能可用
- Payload 可读

## 9. 当前项目里的 JWT 配置

当前项目 JWT 配置定义在：

- [internal/config/config.go](/home/fanqicheng/project/jx/mysurl1/internal/config/config.go:34)

当前配置结构是：

```go
type AuthConf struct {
    JWTSecret      string
    ExpireSeconds  int64
    PasswordPepper string
}
```

示例配置在：

- [etc/mysurl1-api.yaml](/home/fanqicheng/project/jx/mysurl1/etc/mysurl1-api.yaml:23)

当前项目实际配置方式：

```yaml
Auth:
  JWTSecret: change-me-in-production
  ExpireSeconds: 86400
  PasswordPepper: ""
```

说明：

- `JWTSecret`：JWT HMAC 签名密钥
- `ExpireSeconds`：token 过期秒数
- `PasswordPepper`：密码哈希时附加的 pepper，不属于 JWT，但和登录认证链路相关

## 10. 当前项目的结论

当前项目是典型的：

- `JWT + Authorization Header`
- `HS256 + JWTSecret`
- 登录成功后签发
- 中间件统一验签

JWT 是一种可签名的身份令牌格式。当前项目使用的是最简单直接的 `HS256` 方案：服务端用 `JWTSecret` 签发和验签，前端通过 `Authorization: Bearer <token>` 传递登录态。

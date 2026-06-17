# 从 0 到 1 实现一个 Go 短链系统：`mysurl1` 的架构设计、性能优化与工程取舍

## 引言

短链系统几乎是后端学习阶段绕不开的经典项目。

表面上看，它只有两个核心动作：

- 把一个长链接转换成短链接
- 根据短链接跳转回原始地址

但当你把它当成一个真正的服务来设计时，问题会迅速变得更工程化：

- 相同链接是否应该复用已有短链
- 跳转流量远高于创建流量时，数据库如何扛住
- 大量随机短码请求如何避免穿透数据库
- 热点短码缓存失效时，如何避免击穿
- 访问统计应该同步写库还是异步聚合
- 用户登录态如何设计，才能既便于接入前端，又具备吊销能力

`mysurl1` 这个项目，正好就是一个围绕这些问题逐步补齐的 Go 短链系统。它基于 `go-zero` 实现，覆盖了短链生成、短链跳转、缓存、Bloom Filter、`singleflight`、异步访问统计，以及 `access token + refresh token` 认证等关键能力。

本文会从架构、核心链路和工程取舍三个角度，系统介绍这个项目。

---

## 一、项目定位：它不是一个 CRUD Demo

很多短链示例项目，本质上只是下面这套流程：

1. 收到长链接
2. 生成一个短码
3. 写数据库
4. 跳转时按短码查库

这种实现当然能跑，但只适合演示功能，不适合讨论系统设计。

`mysurl1` 的价值在于，它在功能之外，补上了几个真正影响服务质量的点：

- 创建链路支持去重复用，而不是每次都新建
- 跳转链路引入 Redis 精确缓存
- 不存在短码通过 Bloom Filter 预过滤
- 缓存 miss 时通过 `singleflight` 合并回源
- 访问次数先写 Redis，再后台回刷 MySQL
- 用户态接口采用双 token 认证，而不是单 JWT

所以更准确地说，`mysurl1` 是一个适合做“中小型后端架构拆解”的项目，而不只是一个练手接口集合。

---

## 二、整体结构：典型但干净的 Go 分层

这个项目的目录结构比较清晰：

- `mysurl1.go`：服务入口
- `internal/config`：配置定义
- `internal/handler`：HTTP 入口
- `internal/logic`：业务编排
- `internal/dao`：MySQL / Redis 访问
- `internal/model`：数据库表结构映射
- `internal/schema`：接口请求与响应结构
- `internal/middleware`：认证、操作日志等横切逻辑
- `internal/svc`：依赖装配
- `internal/utils`：复用工具与后台 worker

这类结构的价值，不是“看起来规范”，而是它天然约束了代码职责：

- `handler` 只负责收参与响应，不处理复杂业务
- `logic` 负责流程编排和判断
- `dao` 只负责存储访问
- `utils` 放可复用基础能力，不放领域逻辑

对于一个会持续演进的后端项目来说，这种边界很重要。因为只要边界开始混乱，后续每增加一个功能，维护成本都会上升。

服务入口本身也很克制：

```go
var c config.Config
conf.MustLoad(*configFile, &c)

server := rest.MustNewServer(c.RestConf)
defer server.Stop()

ctx := svc.NewServiceContext(c)
handler.RegisterHandlers(server, ctx)

server.Start()
```

代码位置：[mysurl1.go](/home/fanqicheng/project/jx/mysurl1/mysurl1.go:1)

入口只做三件事：

1. 读取配置
2. 初始化依赖
3. 启动 HTTP 服务

这是一个健康项目该有的样子。

---

## 三、依赖装配：`ServiceContext` 是系统骨架

项目所有共享依赖都通过 `ServiceContext` 装起来，包括：

- MySQL 连接
- Redis 客户端
- ShortLink DAO
- 用户 DAO
- refresh token DAO
- 短链缓存
- 短码生成器管理器
- `singleflight` 组

核心装配代码如下：

```go
serviceContext = &ServiceContext{
	Config:      c,
	DB:          newMySQL(c.MySQL),
	Redis:       newRedis(c.Redis),
	FlightGroup: syncx.NewSingleFlight(),
}
serviceContext.ShortLinkCache = dao.NewShortLinkCache(serviceContext.Redis)
serviceContext.ShortLinkDAO = dao.NewShortLinkDAO(serviceContext.DB)
serviceContext.UserDAO = dao.NewUserDAO(serviceContext.DB)
serviceContext.UserRefreshTokenDAO = dao.NewUserRefreshTokenDAO(serviceContext.DB)
serviceContext.UserOperationLogDAO = dao.NewUserOperationLogDAO(serviceContext.DB)
serviceContext.CodeManager = mustNewCodeManager(c.Short, serviceContext.ShortLinkDAO)
utils.StartVisitFlushWorker(serviceContext.DB, serviceContext.ShortLinkCache, c.VisitFlush)
```

代码位置：[internal/svc/servicecontext.go](/home/fanqicheng/project/jx/mysurl1/internal/svc/servicecontext.go:1)

这里有一个值得注意的点：访问统计后台 worker 也是在依赖装配阶段启动的。这意味着它不是一个零散的 goroutine，而是系统能力的一部分。

---

## 四、创建短链：先复用，再生成

短链系统的第一个常见误区，是把“创建短链”理解成“无脑生成一个新码并写库”。

但真实场景里，用户可能会多次提交相同链接。如果每次都生成新短码，会带来几个问题：

- 同一用户数据膨胀
- 缓存和数据库产生重复映射
- 用户体验不稳定

`mysurl1` 的创建逻辑更合理。它的处理顺序是：

1. 校验长链接格式
2. 规范化 URL
3. 先查 Redis 是否已有短链映射
4. Redis miss 后查 MySQL
5. 命中则复用
6. 未命中才真正创建
7. 最后回填缓存和 Bloom

关键代码如下：

```go
shortCode, cacheErr := l.svcCtx.ShortLinkCache.GetLongToShort(l.ctx, userID, normalizedURL)
if cacheErr == nil && shortCode != "" {
	return &createLinkResult{
		ShortCode:   shortCode,
		OriginalURL: normalizedURL,
		Source:      createLinkSourceCacheHit,
	}, nil
}

record, err := l.svcCtx.ShortLinkDAO.FindAvailableByOriginalURL(l.ctx, userID, normalizedURL)
if err == nil && record != nil {
	l.fillCreateCaches(userID, normalizedURL, record.ShortCode)
	return &createLinkResult{
		ShortCode:   record.ShortCode,
		OriginalURL: record.OriginalURL,
		Source:      createLinkSourceDBHit,
	}, nil
}
```

代码位置：[internal/logic/createlinklogic.go](/home/fanqicheng/project/jx/mysurl1/internal/logic/createlinklogic.go:1)

这段逻辑背后体现的是两个原则：

- 优先复用，避免重复创建
- 缓存优先，数据库兜底

这比“功能上能返回一个短链接”更接近真实服务设计。

---

## 五、短码生成：为什么要做成策略模式

短码的生成方式，直接决定了系统在可扩展性和部署方式上的上限。

项目中把发号逻辑抽象成了统一接口：

```go
type Generator interface {
	Provider() string
	NextCode(ctx context.Context) (string, error)
}
```

代码位置：[internal/logic/code_strategy/strategy.go](/home/fanqicheng/project/jx/mysurl1/internal/logic/code_strategy/strategy.go:1)

当前支持三种策略。

### 1. MySQL 自增

这是最直接的方案：先插入一条记录，再把自增主键转成 Base62 短码。

```go
result, err := session.ExecCtx(ctx, insertQuery, userID, normalizedURL, urlHash)
id, err := result.LastInsertId()
shortCode = codestrategy.BuildCodeFromID(uint64(id))
updateQuery := "UPDATE short_links SET short_code = ? WHERE id = ?"
```

代码位置：[internal/dao/shortlinkdao.go](/home/fanqicheng/project/jx/mysurl1/internal/dao/shortlinkdao.go:120)

优点是简单稳定，缺点是生成动作和数据库事务耦合较深。

### 2. Redis 自增

把 Redis 当成一个全局序列号生成器：

```go
sequenceID, err := g.redis.Incr(ctx, redisShortCodeSequenceKey).Uint64()
shortCode := utils.EncodeBase62(sequenceID)
```

代码位置：[internal/logic/code_strategy/redis_incr.go](/home/fanqicheng/project/jx/mysurl1/internal/logic/code_strategy/redis_incr.go:1)

这种方式比 MySQL 自增更轻量，但会依赖 Redis 作为发号中心。

### 3. Snowflake

使用雪花算法生成分布式 ID，再转 Base62：

```go
id := g.node.Generate().Int64()
shortCode := utils.EncodeBase62(uint64(id))
```

代码位置：[internal/logic/code_strategy/snowflake.go](/home/fanqicheng/project/jx/mysurl1/internal/logic/code_strategy/snowflake.go:1)

这种方式适合未来横向扩展场景，但会引入 worker ID 管理问题。

把这些方案抽象成统一策略后，系统可以在“创建流程不改”的前提下切换短码生成方式。这是一个很典型、也很值得保留的设计点。

---

## 六、跳转链路：真正体现系统价值的地方

短链系统里最核心的不是创建，而是跳转。

因为跳转请求通常远高于创建请求，而且读流量集中、热点明显、攻击面也更大。

`mysurl1` 的跳转链路按这个顺序处理：

1. 先查 Redis 精确缓存
2. 缓存 miss 后用 Bloom Filter 预判是否存在
3. Bloom 明确判不存在，则直接返回 404
4. Bloom 认为可能存在时，再查 MySQL
5. 查到后回填缓存
6. 回源过程通过 `singleflight` 合并
7. 访问统计只记 Redis，不直接写库

核心代码如下：

```go
cacheValue, cacheErr := l.svcCtx.ShortLinkCache.GetShortToLong(l.ctx, code)
if cacheErr == nil && cacheValue != nil {
	return l.returnRedirectTarget(code, cacheValue.ID, cacheValue.OriginalURL, "short->long cache")
}

bloomExists, bloomErr := l.svcCtx.ShortLinkCache.ShortCodeBloomExists(l.ctx, code)
if bloomErr == nil && !bloomExists {
	return "", utils.NotFound("short link not found")
}

result, fresh, err := l.svcCtx.FlightGroup.DoEx(redirectSingleflightKey(code), func() (any, error) {
	record, err := l.svcCtx.ShortLinkDAO.FindAvailableByCode(l.ctx, code)
	...
})
```

代码位置：[internal/logic/redirectlogic.go](/home/fanqicheng/project/jx/mysurl1/internal/logic/redirectlogic.go:1)

这套流程很有代表性，因为它同时覆盖了两个缓存领域的典型问题：

- 缓存穿透
- 缓存击穿

### 1. 用 Bloom Filter 抗穿透

如果有人不断请求随机短码，而这些短码都不存在，那么只依赖“缓存 miss 后查库”会导致数据库被持续打空。

Bloom Filter 的作用就是让这类请求尽量在缓存层结束。

在这个项目里，Bloom 的命令是直接走 RedisBloom 模块：

```go
result, err := c.redis.Do(ctx, "BF.EXISTS", bloomShortCodeKey, shortCode).Bool()
```

代码位置：[internal/dao/shortlinkcache.go](/home/fanqicheng/project/jx/mysurl1/internal/dao/shortlinkcache.go:1)

### 2. 用 `singleflight` 抗击穿

如果某个热点短码缓存刚好过期，瞬间有很多请求同时进来，那么所有请求都可能一起打到数据库。

这是典型的缓存击穿。

这个项目通过 `singleflight` 把同一短码的并发回源合并成一次查询，其余请求复用结果。这在高并发热点访问场景里非常有效。

---

## 七、Bloom Filter 的降级设计：比“依赖新组件”更重要

在很多示例项目里，一旦你引入了 RedisBloom、消息队列或分布式锁组件，项目就会默认这些依赖永远存在。

这种写法的一个问题是：环境稍微不完整，系统就起不来，或者逻辑直接报错。

`mysurl1` 在 Bloom 这件事上做了更实用的处理。如果 Redis 没有加载 RedisBloom 模块，它会自动降级：

```go
if err != nil {
	if c.handleBloomUnavailable(err) {
		return true, nil
	}

	return false, err
}
```

代码位置：[internal/dao/shortlinkcache.go](/home/fanqicheng/project/jx/mysurl1/internal/dao/shortlinkcache.go:1)

也就是说：

- 有 RedisBloom 时，走 Bloom 优化路径
- 没有 RedisBloom 时，退回“缓存 + MySQL”逻辑

这不是一个复杂设计，但非常体现工程意识。因为“能降级”往往比“引入了高级组件”更重要。

---

## 八、访问统计：为什么不在请求链路里直写 MySQL

短链跳转通常会伴随访问计数更新。

最直观的做法是，每次跳转执行一次：

```sql
UPDATE short_links
SET visit_count = visit_count + 1
WHERE id = ?
```

但这种方案的问题很明显：

- 热点短链会形成高频写热点
- MySQL 承受大量细粒度 update
- 请求延迟受数据库写性能影响

`mysurl1` 没有在请求链路里直接写库，而是：

1. 先在 Redis 里 `INCR`
2. 后台周期性扫描访问计数键
3. 批量事务回刷 MySQL
4. 删除已刷盘的 Redis 临时键

请求侧代码非常轻：

```go
func (c *ShortLinkCache) IncrVisitCount(ctx context.Context, id uint64) error {
	return c.redis.Incr(ctx, visitCountKey(id)).Err()
}
```

后台回刷逻辑：

```go
err = db.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
	for _, item := range items {
		if _, err := session.ExecCtx(ctx, "UPDATE short_links SET visit_count = visit_count + ? WHERE id = ?", item.delta, item.id); err != nil {
			return err
		}
	}

	return nil
})
```

代码位置：

- [internal/dao/shortlinkcache.go](/home/fanqicheng/project/jx/mysurl1/internal/dao/shortlinkcache.go:1)
- [internal/utils/visit_flush_worker.go](/home/fanqicheng/project/jx/mysurl1/internal/utils/visit_flush_worker.go:1)

这种模式适用于“允许短暂延迟一致”的计数字段，是典型的用缓存换写放大的思路。

---

## 九、认证体系：为什么双 token 比单 JWT 更合适

项目的认证部分也不是最简单的“登录发一个 JWT”。

它采用的是：

- 短期有效的 `access token`
- 长期有效的 `refresh token`
- refresh token 服务端持久化
- refresh token rotation
- 登出可吊销 refresh token

登录成功后的逻辑如下：

```go
authResp, err := utils.BuildAuthResponse(authConf, utils.AuthClaims{
	UserID:   user.ID,
	Username: user.Username,
})

if err := l.svcCtx.UserRefreshTokenDAO.Insert(
	l.ctx,
	user.ID,
	utils.HashRefreshToken(authResp.RefreshToken),
	time.Unix(authResp.RefreshExpiresAt, 0),
); err != nil {
	return nil, utils.InternalError("save refresh token failed: " + err.Error())
}
```

代码位置：[internal/logic/loginlogic.go](/home/fanqicheng/project/jx/mysurl1/internal/logic/loginlogic.go:1)

双 token 的生成逻辑如下：

```go
refreshToken, err := GenerateRefreshToken()
refreshExpiresAt := time.Now().Add(time.Duration(ensureRefreshExpireSeconds(auth)) * time.Second)

return &TokenPair{
	AccessToken:      accessToken,
	AccessExpiresAt:  accessExpiresAt,
	RefreshToken:     refreshToken,
	RefreshTokenHash: HashRefreshToken(refreshToken),
	RefreshExpiresAt: refreshExpiresAt,
}, nil
```

代码位置：[internal/utils/auth.go](/home/fanqicheng/project/jx/mysurl1/internal/utils/auth.go:1)

这种设计相比单 JWT 的优势在于：

- access token 可设置得更短，降低泄露风险
- refresh token 可服务端吊销，支持真正登出
- 前端可以自动刷新，无需频繁重新登录
- 更适合前后端分离应用

它的复杂度确实比单 JWT 高一点，但在真实项目里，这种复杂度是值得的。

---

## 十、接口边界：路由层清楚，业务层干净

从路由定义也能看出这个项目边界比较明确：

```go
server.AddRoutes(
	[]rest.Route{
		{
			Method:  http.MethodGet,
			Path:    "/:code",
			Handler: RedirectHandler(serverCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/auth/login",
			Handler: operationLogMiddleware.Handle(LoginHandler(serverCtx)),
		},
	},
)
```

代码位置：[internal/handler/routes.go](/home/fanqicheng/project/jx/mysurl1/internal/handler/routes.go:1)

可以看到：

- 跳转接口公开暴露
- 用户态接口通过鉴权中间件保护
- 操作日志作为中间件统一接入

这让业务逻辑不需要关心“请求从哪里来、日志怎么记、token 怎么校验”这些横切问题。

---

## 十一、这个项目的工程价值到底在哪里

如果只从功能数量看，`mysurl1` 并不复杂。

但从后端设计质量看，它具备几个值得单独拿出来讲的点：

- 有清晰的分层结构，职责边界明确
- 创建链路考虑了去重与复用
- 跳转链路同时处理了穿透与击穿
- 统计链路用异步聚合降低数据库写压力
- 短码生成策略可切换
- 认证体系不是玩具级实现

这意味着它已经超出了“会写接口”的阶段，而进入了“开始考虑服务在真实运行下会遇到什么问题”的阶段。

对于学习 Go Web 的开发者来说，这种项目比单纯的 CRUD 更有参考价值。

---

## 十二、如果继续演进，我会优先补什么

如果把 `mysurl1` 继续往更接近生产环境的方向推进，我会优先考虑这些内容：

### 1. 短链生命周期管理

- 支持过期时间
- 支持主动失效
- 支持软删除与恢复

### 2. 统计维度扩展

- UV
- IP 分布
- User-Agent / Referer
- 时间窗口聚合

### 3. 缓存治理

- TTL 策略
- 热点预热
- 缓存重建机制
- Bloom 初始化与重建流程

### 4. 可观测性

- 更细的业务指标
- tracing
- 慢查询与缓存命中率监控

### 5. 测试补齐

- 核心逻辑单元测试
- DAO 层集成测试
- 并发场景测试
- 缓存与 worker 协同测试

这些工作不会改变项目主线，但会显著提升它的长期维护价值。

---

## 结语

短链系统常常被当作入门项目，但真正值得学习的地方，从来不是“把 URL 变短”这件事本身，而是你如何围绕它处理缓存、并发、认证、统计和扩展性。

`mysurl1` 的可取之处在于，它没有停留在“功能跑通”，而是继续往前走了一步，把这些现实问题逐步纳入了设计。

如果你正在学习 Go 后端、缓存优化、接口分层，或者想找一个适合做架构拆解和项目复盘的案例，这个仓库是一个不错的起点。

---

## 附：相关源码入口

- 服务入口：[mysurl1.go](/home/fanqicheng/project/jx/mysurl1/mysurl1.go:1)
- 路由注册：[internal/handler/routes.go](/home/fanqicheng/project/jx/mysurl1/internal/handler/routes.go:1)
- 依赖装配：[internal/svc/servicecontext.go](/home/fanqicheng/project/jx/mysurl1/internal/svc/servicecontext.go:1)
- 创建短链逻辑：[internal/logic/createlinklogic.go](/home/fanqicheng/project/jx/mysurl1/internal/logic/createlinklogic.go:1)
- 跳转逻辑：[internal/logic/redirectlogic.go](/home/fanqicheng/project/jx/mysurl1/internal/logic/redirectlogic.go:1)
- 短链缓存：[internal/dao/shortlinkcache.go](/home/fanqicheng/project/jx/mysurl1/internal/dao/shortlinkcache.go:1)
- DAO 实现：[internal/dao/shortlinkdao.go](/home/fanqicheng/project/jx/mysurl1/internal/dao/shortlinkdao.go:1)
- 访问统计回刷 worker：[internal/utils/visit_flush_worker.go](/home/fanqicheng/project/jx/mysurl1/internal/utils/visit_flush_worker.go:1)
- 登录逻辑：[internal/logic/loginlogic.go](/home/fanqicheng/project/jx/mysurl1/internal/logic/loginlogic.go:1)
- 认证工具：[internal/utils/auth.go](/home/fanqicheng/project/jx/mysurl1/internal/utils/auth.go:1)

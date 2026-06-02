# mysurl-v1

一个基于 `go-zero` 实现的短链服务第一版。

当前版本已经支持：
- 创建短链
- 相同长链复用已有短链
- 访问短链并返回 `302` 跳转
- 基于 MySQL 持久化短链映射
- 支持 `mysql_auto_increment`、`redis_incr`、`snowflake` 三种短码生成策略
- V2 接口链路已接入 Redis 精确缓存，并支持 RedisBloom 缺失时自动降级
- V3 已将 `visit_count` 调整为 Redis 聚合后异步回刷 MySQL


## 运行依赖

- Go `1.24.4`
- MySQL
- Redis

说明：
- `mysql_auto_increment` 方案依赖 MySQL
- `redis_incr` 方案依赖 Redis
- `snowflake` 方案本地发号
- 若使用 V2 中的 Bloom 优化，Redis 需安装 `RedisBloom` 模块；未安装时系统自动降级到“精确缓存 + MySQL”路径

## 配置说明

默认配置文件是 [etc/mysurl1-api.yaml](/home/fanqicheng/project/jx/mysurl1/etc/mysurl1-api.yaml)。

关键配置项：

```yaml
MySQL:
  Host: 127.0.0.1
  Port: 3306
  User: root
  Password: root
  Database: short

Redis:
  Host: 127.0.0.1
  Port: 6379
  Password: ""
  DB: 0

Short:
  BaseURL: http://127.0.0.1:8888
  Provider: snowflake
  Snowflake:
    WorkerID: 1
```

`Short.Provider` 可选值：
- `mysql_auto_increment`
- `redis_incr`
- `snowflake`

当 `Short.Provider=snowflake` 时，需要配置：
- `Short.Snowflake.WorkerID`

## 数据库初始化

表结构模板见：
- [internal/template/sqls/short_links.sql](/home/fanqicheng/project/jx/mysurl1/internal/template/sqls/short_links.sql)

初始化后，确保配置中的 `MySQL.Database` 指向对应库。

## 启动方式

在项目根目录执行：

```bash
go run mysurl1.go -f etc/mysurl1-api.yaml
```

也可以先编译：

```bash
go build ./...
./mysurl1 -f etc/mysurl1-api.yaml
```

## 接口使用

### 1. 创建短链

请求：

```bash
curl -X POST 'http://127.0.0.1:8888/api/v1/links' \
  -H 'Content-Type: application/json' \
  -d '{
    "long_url": "https://github.com/JCCGGKS/mysurl"
  }'
```

成功响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "short_code": "2sfCSPGW2gE",
    "short_url": "http://127.0.0.1:8888/2sfCSPGW2gE",
    "original_url": "https://github.com/JCCGGKS/mysurl"
  },
  "timestamp": 1748620800
}
```

说明：
- 相同长链重复提交时会直接复用已有短链
- 创建接口返回统一响应体：`code`、`msg`、`data`、`timestamp`

### 2. 访问短链

请求：

```bash
curl -i 'http://127.0.0.1:8888/2sfCSPGW2gE'
```

成功行为：
- 返回 `302 Found`
- 响应头里带 `Location: <original_url>`

说明：
- 跳转接口保持原生重定向语义，不包装统一 JSON 响应体

## 开发验证

常用命令：

```bash
go build ./...
go test ./...
go test -cover ./...
gofmt -w .
```

## 调试说明

服务基于 `go-zero`，可通过配置开启内部 `DevServer`。

当前配置文件里默认关闭：

```yaml
DevServer:
  Enabled: false
```

如果开启，可使用：
- Health: `http://127.0.0.1:6060/healthz`
- Metrics: `http://127.0.0.1:6060/metrics`
- Pprof: `http://127.0.0.1:6060/debug/pprof/`

常用 pprof 命令：

```bash
go tool pprof http://127.0.0.1:6060/debug/pprof/heap
go tool pprof -alloc_space http://127.0.0.1:6060/debug/pprof/heap
go tool pprof -inuse_space http://127.0.0.1:6060/debug/pprof/heap
```

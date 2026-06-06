# mysurl1

一个基于 `go-zero` 实现的短链服务示例项目，当前版本覆盖了短链创建、跳转、缓存优化、Bloom 过滤器和并发回源控制等核心链路。

## 功能概览

- `POST /api/v1/links`：创建短链
- `GET /:code`：302 跳转到原始长链
- 相同规范化长链复用已有短链
- 支持三种短码生成策略：
  - `mysql_auto_increment`
  - `redis_incr`
  - `snowflake`
- Redis 精确缓存：
  - `normalized_url -> short_code`
  - `short_code -> original_url`
- RedisBloom 优化长链复用判断
- `singleflight` 合并缓存失效后的并发回源
- `visit_count` 走 Redis `INCR` 聚合，后台异步回刷 MySQL

## 项目结构

- `mysurl1.go`：服务入口
- `etc/mysurl1-api.yaml`：本地配置
- `internal/config`：配置结构体
- `internal/dao`：MySQL / Redis 访问层
- `internal/logic`：创建、跳转等业务逻辑
- `internal/handler`：HTTP 处理器
- `internal/svc`：依赖注入与服务上下文
- `internal/utils`：工具方法、统一响应、后台 worker
- `internal/template/sqls`：表结构模板
- `internal/template/docs`：PRD、技术设计文档
- `wrk`：压测脚本与压测记录

## 依赖

- Go `1.24.x`
- MySQL
- Redis

说明：

- `redis_incr` 依赖 Redis 发号
- Bloom 优化依赖 RedisBloom 模块；未安装时会自动降级到“缓存 + MySQL”路径

## 快速开始

### 1. 初始化表结构

执行 [internal/template/sqls/short_links.sql](/home/fanqicheng/project/jx/mysurl1/internal/template/sqls/short_links.sql)。

### 2. 配置服务

默认配置文件是 [etc/mysurl1-api.yaml](/home/fanqicheng/project/jx/mysurl1/etc/mysurl1-api.yaml)：

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

### 3. 启动服务

```bash
go run mysurl1.go -f etc/mysurl1-api.yaml
```

## 接口示例

### 创建短链

```bash
curl -X POST 'http://127.0.0.1:8888/api/v1/links' \
  -H 'Content-Type: application/json' \
  -d '{"long_url":"https://github.com/JCCGGKS/mysurl"}'
```

返回统一 JSON 结构：

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

### 跳转短链

```bash
curl -i 'http://127.0.0.1:8888/2sfCSPGW2gE'
```

成功时返回 `302 Found`，并在 `Location` 头中带上目标长链。

## 开发与测试

```bash
go build ./...
go test ./...
go test -cover ./...
gofmt -w .
```

受限环境下可使用：

```bash
GOCACHE=/tmp/gocache go test ./internal/...
```

## 压测与文档

- `wrk/run_create.sh`：创建接口压测
- `wrk/run_get.sh`：单短码跳转压测
- `wrk/run_create_urls.sh`：批量长链创建压测
- `wrk/run_get_noexists_code.sh`：大量不存在短码压测

相关设计与分析文档：

- `internal/template/docs/PRD.md`
- `internal/template/docs/TECH_V1.md`
- `internal/template/docs/TECH_V2.md`
- `internal/template/docs/TECH_V3.md`
- `wrk/wrk1.md`
- `wrk/wrk2.md`
- `wrk/wrk3.md`

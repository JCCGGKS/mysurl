# QUICKSTART

本文档用于把 `mysurl1` 在本地快速跑起来，覆盖前置准备、数据库初始化、后端启动、前端启动和基础自检。

## 1. 项目概览

- 后端技术栈：Go + go-zero + MySQL + Redis
- 前端技术栈：Vue 3 + Vite
- 后端默认地址：`http://127.0.0.1:8888`
- 前端默认地址：`http://127.0.0.1:5173`

后端入口：

```bash
go run mysurl1.go -f etc/mysurl1-api.yaml
```

前端开发入口：

```bash
cd vue
npm install
npm run dev
```

## 2. 前置准备

本项目本地运行至少需要以下依赖：

- Go `1.24.4` 或更高版本
- MySQL `8.x`
- Redis `6.x` 或更高版本
- Node.js `18+`
- npm

先确认本机环境：

```bash
go version
mysql --version
redis-server --version
node -v
npm -v
```

## 3. 配置说明

默认配置文件是 [etc/mysurl1-api.yaml](/home/fanqicheng/project/jx/mysurl1/etc/mysurl1-api.yaml:1)：

```yaml
Host: 0.0.0.0
Port: 8888
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
```

本地最重要的是这几项：

- `MySQL`：要和你本地数据库账号、端口保持一致
- `Redis`：要和你本地 Redis 实例保持一致
- `Short.BaseURL`：决定返回的短链域名，本地一般保持 `http://127.0.0.1:8888`
- `Short.Provider`：默认使用 `snowflake`

可选的短码策略：

- `mysql_auto_increment`
- `redis_incr`
- `snowflake`

如果你只是本地启动验证，保持默认 `snowflake` 即可。

## 4. 初始化 MySQL

### 4.1 创建数据库和表

仓库自带初始化脚本 [internal/template/sqls/install_db.sh](/home/fanqicheng/project/jx/mysurl1/internal/template/sqls/install_db.sh:1)，会自动：

- 创建数据库 `short`
- 执行 `internal/template/sqls/*.sql`

直接运行：

```bash
bash internal/template/sqls/install_db.sh
```

如果你的本地 MySQL 不是默认账号密码，可以覆盖环境变量：

```bash
MYSQL_HOST=127.0.0.1 \
MYSQL_PORT=3306 \
MYSQL_USER=root \
MYSQL_PASSWORD=root \
MYSQL_DATABASE=short \
bash internal/template/sqls/install_db.sh
```

### 4.2 会创建哪些表

初始化脚本会导入这些表结构：

- `users`
- `short_links`
- `user_refresh_tokens`
- `user_operation_logs`

SQL 文件位于 [internal/template/sqls](/home/fanqicheng/project/jx/mysurl1/internal/template/sqls)。

## 5. 启动 Redis

确保本地 Redis 已启动，并且与配置文件一致。

例如本地直接启动：

```bash
redis-server
```

如果你已经有 Redis 服务在后台运行，可以直接验证：

```bash
redis-cli -h 127.0.0.1 -p 6379 ping
```

返回 `PONG` 即可。

### 5.1 RedisBloom 说明

项目支持使用 RedisBloom 做不存在短码拦截，但它不是本地启动的硬依赖。

- 如果 Redis 已安装 RedisBloom，项目会自动使用
- 如果没有安装，代码会自动降级为 `Redis 缓存 + MySQL` 路径

所以本地快速启动时，不需要先解决 RedisBloom。

## 6. 启动后端

在仓库根目录执行：

```bash
go run mysurl1.go -f etc/mysurl1-api.yaml
```

看到类似输出说明启动成功：

```text
Starting server at 0.0.0.0:8888...
```

如果你只想先做编译检查，也可以执行：

```bash
go build ./...
```

## 7. 启动前端

前端目录在 [vue](/home/fanqicheng/project/jx/mysurl1/vue)。

首次启动：

```bash
cd vue
npm install
npm run dev
```

默认访问地址：

```text
http://127.0.0.1:5173
```

前端开发服务器已经配置了代理 [vue/vite.config.js](/home/fanqicheng/project/jx/mysurl1/vue/vite.config.js:1)：

- `/api` 会转发到 `http://127.0.0.1:8888`

这意味着本地联调时，通常只需要：

1. 启动 MySQL
2. 启动 Redis
3. 启动后端
4. 启动前端

## 8. 最小启动步骤

如果你只想最快跑起来，按这个顺序执行：

```bash
# 1) 初始化数据库
bash internal/template/sqls/install_db.sh

# 2) 启动 Redis
redis-server

# 3) 启动后端
go run mysurl1.go -f etc/mysurl1-api.yaml

# 4) 新开一个终端启动前端
cd vue
npm install
npm run dev
```

## 9. 基础自检

### 9.1 注册用户

```bash
curl -X POST 'http://127.0.0.1:8888/api/v1/auth/register' \
  -H 'Content-Type: application/json' \
  -d '{
    "username": "demo",
    "password": "123456",
    "confirm_password": "123456"
  }'
```

### 9.2 登录

```bash
curl -X POST 'http://127.0.0.1:8888/api/v1/auth/login' \
  -H 'Content-Type: application/json' \
  -d '{
    "username": "demo",
    "password": "123456"
  }'
```

返回体是统一结构：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "access_token": "...",
    "refresh_token": "...",
    "user": {
      "id": 1,
      "username": "demo"
    }
  }
}
```

### 9.3 创建短链

把上一步返回的 `access_token` 替换到下面命令里：

```bash
curl -X POST 'http://127.0.0.1:8888/api/v1/links' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer <access_token>' \
  -d '{
    "long_url": "https://example.com/article/1"
  }'
```

### 9.4 访问短链

创建成功后，返回值里会包含：

- `short_code`
- `short_url`

直接访问 `short_url` 即可验证跳转能力。

## 10. 常见问题

### 10.1 后端启动时报 MySQL 连接错误

先检查：

- MySQL 是否已启动
- [etc/mysurl1-api.yaml](/home/fanqicheng/project/jx/mysurl1/etc/mysurl1-api.yaml:1) 中的账号密码是否正确
- 数据库 `short` 是否已创建

如果没初始化过表，重新执行：

```bash
bash internal/template/sqls/install_db.sh
```

### 10.2 后端启动时报 Redis 连接错误

先检查：

- Redis 是否已启动
- 配置文件中的 `Host`、`Port`、`Password` 是否正确

### 10.3 短链创建成功，但访问时报 404

优先检查：

- `Short.BaseURL` 是否仍然是本地服务地址
- 你访问的服务是否就是启动中的 `8888` 端口
- 后端服务是否正常运行

### 10.4 前端请求接口失败

先确认：

- 前端是否通过 `npm run dev` 在 `5173` 端口启动
- 后端是否在 `8888` 端口启动
- 本地是否修改过 [vue/vite.config.js](/home/fanqicheng/project/jx/mysurl1/vue/vite.config.js:1) 的代理配置

## 11. 常用命令

```bash
# 编译全部 Go 代码
go build ./...

# 运行全部测试
go test ./...

# 运行 internal 下测试
GOCACHE=/tmp/gocache go test ./internal/...

# 启动后端
go run mysurl1.go -f etc/mysurl1-api.yaml

# 启动前端
cd vue && npm run dev
```

## 12. 相关文件

- [mysurl1.go](/home/fanqicheng/project/jx/mysurl1/mysurl1.go:1)：后端启动入口
- [etc/mysurl1-api.yaml](/home/fanqicheng/project/jx/mysurl1/etc/mysurl1-api.yaml:1)：本地配置
- [mysurl1.api](/home/fanqicheng/project/jx/mysurl1/mysurl1.api:1)：接口定义
- [internal/template/sqls/install_db.sh](/home/fanqicheng/project/jx/mysurl1/internal/template/sqls/install_db.sh:1)：数据库初始化脚本
- [vue/package.json](/home/fanqicheng/project/jx/mysurl1/vue/package.json:1)：前端启动脚本

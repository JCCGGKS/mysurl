# 配置加载链路

## 概览

当前项目不使用 `viper`。

配置加载入口在 [mysurl1.go](/home/fanqicheng/project/jx/mysurl1/mysurl1.go:26)：

```go
var c config.Config
conf.MustLoad(*configFile, &c)
```

这里使用的是 `github.com/zeromicro/go-zero/core/conf`。

## 加载流程

配置加载主链路如下：

1. `conf.MustLoad(path, &c)`
2. `conf.Load(path, &c)`
3. 根据文件扩展名选择 loader
4. `.yaml/.yml` 走 `LoadFromYamlBytes`
5. YAML 内容先转换为 JSON
6. JSON 再通过 `mapping.UnmarshalJsonMap` 映射到 Go struct
7. 如果配置 struct 实现了 `Validate()`，最后执行校验

## 关键实现

### 1. `MustLoad`

`MustLoad` 本身很薄，只是包装 `Load`，失败时直接退出进程。

```go
func MustLoad(path string, v any, opts ...Option) {
	if err := Load(path, v, opts...); err != nil {
		log.Fatalf("error: config file %s, %s", path, err.Error())
	}
}
```

### 2. `Load`

`Load` 的职责：

- 读取配置文件内容
- 根据扩展名选择解析器
- 可选地展开环境变量

支持的格式包括：

- `.json`
- `.json5`
- `.toml`
- `.yaml`
- `.yml`

### 3. YAML 实际处理方式

YAML 不是直接反序列化到 struct，而是先转成 JSON：

```go
func LoadFromYamlBytes(content []byte, v any) error {
	b, err := encoding.YamlToJson(content)
	if err != nil {
		return err
	}

	return LoadFromJsonBytes(b, v)
}
```

所以当前项目虽然配置文件格式是 YAML，但内部统一走 JSON 映射逻辑。

### 4. JSON 到 struct 的映射

核心逻辑：

```go
mapping.UnmarshalJsonMap(lowerCaseKeyMap, v,
	mapping.WithCanonicalKeyFunc(toLowerCase))
```

这一步使用的是 `github.com/zeromicro/go-zero/core/mapping`。

## 为什么看的是 `json` tag

`go-zero/conf` 在内部固定使用 `json` 作为 tag key：

```go
const jsonTagKey = "json"
```

因此：

- 配置文件是 YAML
- 但 struct 字段映射看的是 `json:"..."` tag
- 不是 `yaml:"..."` tag

这也是当前配置 struct 里使用 `json:",optional"` 的原因。

## `optional` 在哪里生效

`optional` 不是 YAML 解析器处理的，也不是 `viper` 处理的，而是在 `core/mapping` 层处理的。

处理链路：

1. 解析 struct tag
2. 提取 `optional`、`default` 等字段选项
3. 映射配置值到字段
4. 字段缺失时，根据 `optional` 决定是否报错

如果字段不是 optional，缺失时会走到类似下面的逻辑：

```go
if !opts.optional() {
	return newInitError(fullName)
}
```

含义：

- 有 `optional`：字段缺失允许，使用零值
- 无 `optional`：字段缺失可能报错

## `optional` 的实际意义

以当前项目的配置为例：

```go
type MySQLConf struct {
	Host     string `json:",optional"`
	Port     int    `json:",optional"`
	User     string `json:",optional"`
	Password string `json:",optional"`
	Database string `json:",optional"`
}
```

这表示：

- `MySQL.Host` 可以不在 YAML 中出现
- 不出现时，`Host` 使用 Go 零值 `""`
- 服务是否允许这种情况，取决于业务代码是否额外做启动校验

## `stat` 相关配置

当前项目有 3 个相关开关：

| 配置项 | 作用 | 是否影响采样器 |
| --- | --- | --- |
| `Log.Stat` | 控制 `logx` 的 stat 日志输出 | 否 |
| `Stat.DisableLog` | 控制 `core/stat` 的周期资源日志输出 | 否 |
| `Stat.DisableSampler` | 控制 `core/stat` 后台采样协程是否启动 | 是 |

说明：

- `Log.Stat` 是 go-zero 自带的 `logx.LogConf` 配置
- `Stat.DisableLog` 和 `Stat.DisableSampler` 是项目自定义配置
- `Log.Stat` 只管日志输出，不替代采样器开关

当前项目默认配置：

```yaml
Log:
  Stat: false

Stat:
  DisableLog: true
  DisableSampler: true
```

等价于：

- 不输出 `logx` stat 日志
- 不输出 `core/stat` 周期资源日志
- 不启动 `core/stat` 采样协程

## 结论

- 当前项目配置加载框架是 `go-zero/core/conf`
- 配置文件格式是 YAML
- YAML 会先被转换成 JSON
- struct 映射由 `go-zero/core/mapping` 完成
- 字段约束主要看 `json` tag
- `optional` 在 `mapping` 层决定字段是否允许缺失

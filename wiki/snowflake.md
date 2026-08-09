# Snowflake 学习整理

## 一、Twitter 官方原版 Snowflake

2010 年 Twitter 公开了经典 Snowflake 思路，核心目标是在分布式环境中生成全局唯一、趋势递增的 64 位整数 ID。它是后续大量“雪花类”算法和实现的参考原型，但不同语言、不同库的具体位分配和时钟处理策略并不完全一致。

### 1. 原版标准结构

经典结构通常写作：

`0 | 41bit时间 | 10bit节点 | 12bit序列号`

- 有效时长：约 69 年
- 最大集群节点：1024 个
- 单节点毫秒最大 ID：4096 个

### 2. 原生核心生成逻辑

原生思路非常明确，核心逻辑通常是：

1. 获取当前时间戳
2. 同一毫秒内序列号自增
3. 序列号用尽则等待下一毫秒
4. 新毫秒到来后重置序列号
5. 通过位运算拼接时间、节点号和序列号，返回唯一 ID

### 3. 原版核心缺陷

- 机器 ID 需要保证全局唯一
- 依赖系统时钟、有时钟回拨的风险
- 只能保证趋势递增，不能当作严格连续递增的业务编号使用

+ 改良版雪花算法：https://www.cnblogs.com/thisiswhy/p/17611163.html

## 二、主流 Go 雪花库核心信息

### 1. `bwmarrin/snowflake`（当前项目使用）

- 仓库：`github.com/bwmarrin/snowflake`
- 位结构：默认 `1bit符号位 + 41bit时间 + 10bit机器 + 12bit序列号`，对应源码里的 `NodeBits = 10`、`StepBits = 12`
- 机器号：手动分配，默认范围 `0~1023`，由 `NewNode(node int64)` 校验
- 回拨策略：当前项目使用的 `v0.3.0` 实现并没有在 `Generate()` 中显式写出“检测回拨后直接 panic”的分支；它基于 `time.Since(n.epoch)` 计算相对时间，同毫秒内序列号递增，序列号用尽后等待下一毫秒
- 源码位置：当前项目锁定版本 `v0.3.0` 的源码文件为 <https://github.com/bwmarrin/snowflake/blob/v0.3.0/snowflake.go>
- 重点关注：`NodeBits`、`StepBits` 的默认值，`NewNode()` 对节点号范围的校验，`Generate()` 中同毫秒序列递增和序列耗尽后的等待逻辑
- 核心特点：零依赖、接口简单、社区常见、默认位分配接近经典 Snowflake
- 适用场景：通用业务、传统微服务、需要简单稳定本地发号的线上系统

### 2. `sony/sonyflake`

- 仓库：`github.com/sony/sonyflake`
- 位结构：官方源码明确为 `39bit时间(单位 10ms) + 8bit序列号 + 16bit机器号`
- 机器号：支持通过 `Settings.MachineID` 自定义；未配置时默认取私网 IPv4 地址的低 16 位，`CheckMachineID` 可额外校验唯一性
- 回拨策略：`NextID()` 通过比较 `elapsedTime` 与当前时间片推进发号；当时间片未前进且序列耗尽时，会推进 `elapsedTime` 并 `sleep` 到对应时间片，不是旧文档里那种固定“10ms 阈值报错”规则
- 源码位置：<https://github.com/sony/sonyflake/blob/master/sonyflake.go>
- 重点关注：默认位分配不是经典 `41/10/12`，`Settings` 中机器号获取与校验的扩展方式，`NextID()` 的时间推进和序列控制逻辑
- 核心特点：可配置性更强，适合结合部署环境定制机器号获取方式
- 适用场景：分布式服务、需要自定义机器号策略的部署环境

### 3. `influxdata/snowflake`

- 典型来源：`github.com/influxdata/influxdb/pkg/snowflake`
- 位结构：旧文档引用的是 InfluxDB 历史代码路径中的业务内嵌 Snowflake 生成器；当前官方主仓库 `main` 分支已无法直接按原路径定位到该实现，因此不能继续把位分配写死成固定结论
- 机器号：依赖对应历史版本实现与业务配置，不属于像 `bwmarrin`、`sonyflake` 这样稳定暴露的独立库接口
- 回拨策略：历史实现通常把时钟倒退视为错误返回，而不是静默吞掉；但具体行为仍应以实际引用的 InfluxDB 版本源码为准
- 源码位置：历史路径是 `github.com/influxdata/influxdb/pkg/snowflake/generator.go`；当前官方仓库主线已切到 InfluxDB 3，`main` 分支无法直接定位该文件，如需追溯需切到对应历史分支或版本标签
- 重点关注：具体版本里的位分配定义，时钟回退时的错误处理方式，序列耗尽时如何推进时间片
- 核心特点：更偏工程内嵌实现，适合参考其源码理解雪花生成器在实际项目中的写法
- 适用场景：学习业务工程内雪花实现，或处于相关技术栈上下文中时参考

### 4. `ppzz/golang-snowflake-id`

- 仓库：`github.com/ppzz/golang-snowflake-id`
- 位结构：官方实现的原始拼接是 `时间戳 << (serverIdLength + counterLength) | serverId << counterLength | counter`，随后再经过一次 `chaos()` 位重排；因此它不是文档里适合直接写成单一“41/10/12”口径的实现
- 机器号：使用 `serverId`，由调用方提供并参与最终位拼接
- 回拨策略：官方实现采用等待补齐时间片的思路，例如 `toNextMillisecond()` 计算需要睡眠到下一毫秒，而不是直接报错退出
- 源码位置：仓库根目录 `internal.go`
- 重点关注：时间小于上次时间时的处理方式，节点号、序列号和时间片的拼接逻辑，以及 `chaos()` 位重排
- 核心特点：实现轻量、代码直观、上手成本低
- 适用场景：小型内部系统、学习雪花算法实现、对接入复杂度要求低的服务

### 5. `yitter/IdGenerator`

- 仓库：`github.com/yitter/IdGenerator-Go`
- 位结构：官方默认选项是 `WorkerIdBitLength = 6`、`SeqBitLength = 6`，并要求 `WorkerIdBitLength + SeqBitLength <= 22`；也就是说它的机器位和序列位是可配置的，不是固定单一位宽
- 机器号：使用 `WorkerId`，必须由外部设置，最大值取决于 `2^WorkerIdBitLength - 1`
- 回拨策略：官方默认 `Method = 1` 使用漂移算法；源码里为回拨预留了序列号低位，时间回退时走 `CalcTurnBackId()` 和回拨次序控制，不是简单“报错”或单纯 `sleep` 等待
- 源码位置：
  - 真实仓库为 `github.com/yitter/IdGenerator`，Go 配置定义在 <https://github.com/yitter/IdGenerator/blob/master/Go/source/idgen/IdGeneratorOptions.go>
  - 回拨与发号主逻辑在 <https://github.com/yitter/IdGenerator/blob/master/Go/source/idgen/SnowWorkerM1.go>
- 重点关注：参数配置项、回拨处理相关增强策略、高并发场景下的序列控制方式
- 核心特点：配置项较多，更偏增强型发号器，强调复杂场景下的可用性设计
- 适用场景：对可用性、回拨处理策略和工程增强能力要求更高的系统

## 三、生产选型结论

- 通用业务系统：优先 `bwmarrin/snowflake`（结构经典、实现直接、社区常见）
- 需要更强机器号配置能力的分布式部署：优先 `sony/sonyflake`
- 想学习业务工程里的雪花实现变体：可参考 `influxdata/snowflake`
- 小型内部低复杂度系统：可选 `ppzz/golang-snowflake-id`
- 对增强型回拨处理和工程可用性更敏感的系统：优先 `yitter/IdGenerator`

## 四、生产最佳实践

1. 统一纪元：全业务线应明确统一纪元和 ID 结构，避免不同实现混用导致解释不一致
2. 机器 ID 唯一：物理机、虚拟机、容器实例都必须保证节点号全局不重复
3. 时钟强制同步：集群节点必须开启时间同步机制，并监控时钟漂移
4. 二次封装兜底：业务不要把第三方雪花库当作“无条件可靠黑盒”，应统一封装配置、错误处理和监控
5. 规范存储类型：Go 使用 `int64`，MySQL 使用 `bigint`，不要把雪花 ID 当字符串主键存储
6. 数据库唯一约束不能省：即使发号器理论上保证唯一，落库时仍应保留唯一索引兜底

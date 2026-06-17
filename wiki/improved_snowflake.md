# 改良版雪花算法笔记

## 核心结论

标准雪花算法能保证全局唯一，也具备趋势递增特性。常见 bit 布局是：

```text
0 | timestamp | workerId | sequence
```

标准 Snowflake 对物理时钟过于敏感：一旦运行期发生时钟回拨，传统实现要么拒绝继续发号，要么面临重复 ID 风险。除此之外，标准实现同一毫秒内还有严格的 `4096` 个序列号上限。

Seata 这类改良版雪花的核心思路，是解除生成器与操作系统时间戳的强绑定。为此，它把 `timestamp` 和 `workerId` 对调，让 `timestamp` 与 `sequence` 在低位连续排列：

```text
0 | workerId | timestamp | sequence
```

这样 `timestamp` 和 `sequence` 可以被视为一个连续的 53 位状态，便于直接用一个原子整数整体推进。初始化后，运行期不再持续依赖物理时钟逐毫秒前进，所以小幅时钟回拨不会像标准 Snowflake 那样立刻影响发号；同时，序列号耗尽时也可以直接进位到内部维护的下一个时间片，不再受标准实现严格的 `4096/ms` 限制。

## Seata 的补充

Seata 的实现补充了两个很关键的工程点：

- 标准雪花对时钟回拨敏感，而且单毫秒突发上限本质上是 `4096`
- Seata 不只是调换 bit 位，还把“时间戳 + 序列号”作为连续状态推进，运行期不再强依赖物理时钟逐毫秒前进

这样做的结果是：

- 序列号满了以后，可以直接进位到内部维护的下一个时间片，不必立刻依赖物理时钟推进
- 运行期对时钟漂移或回拨弱依赖，因为初始化后不再持续同步系统时钟
- 但重启后仍然存在重复风险

Seata 还强调了 `workerId` 分配策略本身也是生成器设计的一部分。节点号如果冲突，雪花算法一样会失去全局唯一性。

参考文章里最关键的实现片段大致如下。

把“时间戳 + 序列号”放进同一个原子状态：

```java
/**
 * timestamp and sequence mix in one Long
 * highest 11 bit: not used
 * middle  41 bit: timestamp
 * lowest  12 bit: sequence
 */
private AtomicLong timestampAndSequence;
```

把固定的 `workerId` 作为高位预先放好：

```java
/**
 * business meaning: machine ID (0 ~ 1023)
 * actual layout in memory:
 * highest 1 bit: 0
 * middle 10 bit: workerId
 * lowest 53 bit: all 0
 */
private long workerId;
```

生成 ID 时直接推进这个连续状态：

```java
public long nextId() {
   long next = timestampAndSequence.incrementAndGet();
   long timestampWithSequence = next & timestampAndSequenceMask;
   return workerId | timestampWithSequence;
}
```

这一版代码最值得注意的点不是语法，而是思路：

- `timestamp` 和 `sequence` 已经被当成一个整体递增
- 初始化后，发号过程不再每次都读取并校验当前系统时钟
- 序列溢出时会自然进位到时间戳
- 因而不再有标准雪花那种严格的 `4096/ms` 突发限制
- 但如果内部“超前时间”跑得比物理时间更快，仍可能在物理时间追平前进入阻塞

## 参考文章

- 博客园 why技术：《在开源项目中看到一个改良版的雪花算法，现在它是你的了。》<https://www.cnblogs.com/thisiswhy/p/17611163.html>
- Apache Seata：《Seata基于改良版雪花算法的分布式UUID生成器分析》<https://seata.apache.org/zh-cn/blog/seata-analysis-UUID-generator/>

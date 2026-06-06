# 布隆过滤器

## 1. 定义

布隆过滤器，Bloom Filter，是一种空间效率很高的概率型数据结构，用来判断一个元素是否“可能存在”于集合中。

它只能给出两种结论：

- 一定不存在
- 可能存在

它不能保证“命中就一定存在”。

## 2. 工作原理

布隆过滤器底层通常是一段 bit 数组，加上多个 hash 函数。

假设：

- 有一个长度为 `m` 的 bit 数组，初始全为 `0`
- 有 `k` 个 hash 函数

插入元素 `x`：

1. 计算 `k` 个 hash 位置
2. 把这些位置上的 bit 置为 `1`

查询元素 `x`：

1. 计算同样的 `k` 个 hash 位置
2. 只要有一个 bit 为 `0`，则一定不存在
3. 如果全部为 `1`，则可能存在

## 3. 误判与边界

误判来源：

- 不同元素可能落到同一组 bit 上
- 某个未插入元素查询时，所需 bit 恰好都被其他元素置成了 `1`

所以布隆过滤器的特点是：

- 不会把“已存在”判断成“不存在”
- 会把“某些不存在”误判成“可能存在”

它适合做前置过滤，不适合做最终真值存储。

## 4. 为什么不方便删除

普通布隆过滤器只记录 bit 是否为 `1`，不记录这个 bit 是谁置上的。

如果直接清零：

- 可能把别的元素共享的 bit 一起删掉
- 导致本来存在的元素被误判为不存在

所以标准 Bloom 一般只支持：

- 插入
- 查询

如果要求删除，通常要考虑：

- Counting Bloom Filter
- 或 Cuckoo Filter

## 5. 调参

三个核心参数：

- `n`：预计插入元素数量
- `m`：bit 数组长度
- `k`：hash 函数数量

核心公式：

```text
m = -n * ln(p) / (ln2)^2
k = ln2 * m / n
p = (1 - e^(-kn/m))^k
```

常见经验值：

- `p = 1% (0.01)`：约 `9.6 bit/元素`，`k ≈ 7`
- `p = 0.1% (0.001)`：约 `14.4 bit/元素`，`k ≈ 10`
- `p = 0.01% (0.0001)`：约 `19.2 bit/元素`，`k ≈ 14`

常用在线计算器：

- `https://hur.st/bloomfilter/`
- `https://krisives.github.io/bloom-calculator/`
- `https://codingace.net/maths/bloom_filter.html`

## 6. 典型场景

- 缓存穿透防护
- 短链系统里的长链存在性预判
- 爬虫去重 / 消息去重

在当前短链项目里，它的定位应该是：

- 前置优化组件
- 减少无意义查库和查缓存
- 不参与最终正确性判断

## 7. RedisBloom 源码入口

当前仓库里的 RedisBloom Bloom 实现主要在：

- `RedisBloom/src/rebloom.c`
- `RedisBloom/src/sb.c`
- `RedisBloom/deps/bloom/bloom.c`

职责分层：

- `rebloom.c`：Redis 命令入口
- `sb.c`：可扩容 Bloom chain
- `bloom.c`：底层 bitset / hash / check / add

## 8. RedisBloom 命令层

Bloom 命令在 `rebloom.c` 注册：

```c
RegisterCommand(ctx, "bf.reserve", BFReserve_RedisCommand, "write deny-oom", "write fast");
RegisterCommand(ctx, "bf.add", BFAdd_RedisCommand, "write deny-oom", "write");
RegisterCommand(ctx, "bf.madd", BFAdd_RedisCommand, "write deny-oom", "write");
RegisterCommand(ctx, "bf.insert", BFInsert_RedisCommand, "write deny-oom", "write");
RegisterCommand(ctx, "bf.exists", BFCheck_RedisCommand, "readonly fast", "read");
RegisterCommand(ctx, "bf.mexists", BFCheck_RedisCommand, "readonly fast", "read");
```

`BF.EXISTS` 入口：

```c
int exists = SBChain_Check(sb, s, n);
```

`BF.ADD` 入口：

```c
rv = SBChain_Add(sb, s, n);
```

命令层本身很薄，真正逻辑分别在 `SBChain_Check` 和 `SBChain_Add`。

## 9. RedisBloom 的可扩容 Bloom Chain

RedisBloom 不是只维护一个 Bloom filter，而是维护一个 `SBChain`。

查询：

```c
int SBChain_Check(const SBChain *sb, const void *data, size_t len) {
    bloom_hashval hv = SBChain_GetHash(sb, data, len);
    for (int ii = sb->nfilters - 1; ii >= 0; --ii) {
        if (bloom_check_h(&sb->filters[ii].inner, hv)) {
            return 1;
        }
    }
    return 0;
}
```

新增：

```c
int SBChain_Add(SBChain *sb, const void *data, size_t len) {
    bloom_hashval h = SBChain_GetHash(sb, data, len);
    for (int ii = sb->nfilters - 1; ii >= 0; --ii) {
        if (bloom_check_h(&sb->filters[ii].inner, h)) {
            return 0;
        }
    }
    ...
}
```

关键点：

- 先查后加
- 任意一层命中就认为可能存在
- 当前层满了时可以继续扩容出新的 filter

所以 RedisBloom 的 Bloom 实现不是单块位图，而是一条可扩展的 filter 链。

## 10. 底层 bloom.c

hash 计算：

```c
bloom_hashval bloom_calc_hash64(const void *buffer, int len) {
    bloom_hashval rv;
    rv.a = MurmurHash64A_Bloom(buffer, len, 0xc6a4a7935bd1e995ULL);
    rv.b = MurmurHash64A_Bloom(buffer, len, rv.a);
    return rv;
}
```

多个 hash 位置由这两个基值派生：

- `x = (a + i * b) % mod`

check / add 共用一套核心逻辑：

```c
#define CHECK_ADD_FUNC(T, modExp) \
    ... \
    for (i = 0; i < bloom->hashes; i++) { \
        T x = ((hashval.a + i * hashval.b)) % mod; \
        if (!test_bit_set_bit(bloom->bf, x, mode)) { \
            if (mode == MODE_READ) { \
                return 0; \
            } \
            found_unset = 1; \
        } \
    } \
    if (mode == MODE_READ) { \
        return 1; \
    } \
    return found_unset;
```

含义：

- 读模式：只要发现一个 bit 没置位，就直接返回不存在
- 写模式：边检查边置位，并判断这次是否为新加入

## 11. 返回值语义

`bloom_check_h`：

- `1`：可能存在
- `0`：一定不存在

`bloom_add_h`：

- `0`：原来不存在，本次已新增
- `1`：原来就存在，或发生碰撞

这也说明：

- `BF.EXISTS` 只能给出概率判断
- `BF.ADD` 自身就带“是否已存在”的判断能力

## 12. 总结

Bloom Filter 的核心价值是：

- 用很小的内存
- 很快过滤掉大量一定不存在的数据

RedisBloom 在标准 Bloom 思路上又补了一层工程实现：

- Redis 命令封装
- 可扩容 Bloom chain
- 底层 hash / bitset 统一 check/add 逻辑

对当前项目来说，最重要的理解仍然是：

- Bloom 未命中可以快速判定不存在
- Bloom 命中不能直接当成存在，仍然要查精确缓存或数据库

# RedisBloom Bloom Filter 源码梳理

## 1. 入口

这份文档分析的是仓库里的 RedisBloom 源码：

- [RedisBloom](/home/fanqicheng/project/jx/mysurl1/RedisBloom)

Bloom Filter 相关核心代码主要在三层：

- `src/rebloom.c`
  负责 Redis 命令注册和命令入口，如 `BF.ADD`、`BF.EXISTS`
- `src/sb.c`
  负责可扩容 Bloom Chain，处理“多段 filter”逻辑
- `deps/bloom/bloom.c`
  负责底层 bitset / hash / check / add 实现

可以理解为：

- `rebloom.c`：Redis 模块层
- `sb.c`：结构组织层
- `bloom.c`：位图算法层

## 2. 命令注册

Bloom 命令在 `src/rebloom.c` 注册：

```c
RegisterCommand(ctx, "bf.reserve", BFReserve_RedisCommand, "write deny-oom", "write fast");
RegisterCommand(ctx, "bf.add", BFAdd_RedisCommand, "write deny-oom", "write");
RegisterCommand(ctx, "bf.madd", BFAdd_RedisCommand, "write deny-oom", "write");
RegisterCommand(ctx, "bf.insert", BFInsert_RedisCommand, "write deny-oom", "write");
RegisterCommand(ctx, "bf.exists", BFCheck_RedisCommand, "readonly fast", "read");
RegisterCommand(ctx, "bf.mexists", BFCheck_RedisCommand, "readonly fast", "read");
RegisterCommand(ctx, "bf.info", BFInfo_RedisCommand, "readonly fast", "read fast");
RegisterCommand(ctx, "bf.card", BFCard_RedisCommand, "readonly fast", "read fast");
```

当前项目实际使用的就是：

- `BF.ADD`
- `BF.EXISTS`

## 3. BF.EXISTS

命令入口在 `BFCheck_RedisCommand`：

```c
static int BFCheck_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc) {
    RedisModule_AutoMemory(ctx);
    int is_multi = isMulti(argv[0]);

    RedisModuleKey *key = RedisModule_OpenKey(ctx, argv[1], REDISMODULE_READ);
    SBChain *sb;
    int status = bfGetChain(key, &sb);

    int is_empty = 0;
    if (status != SB_OK) {
        is_empty = 1;
    }

    for (size_t ii = 2; ii < argc; ++ii) {
        if (is_empty == 1) {
            reply = false;
        } else {
            size_t n;
            const char *s = RedisModule_StringPtrLen(argv[ii], &n);
            int exists = SBChain_Check(sb, s, n);
            reply = !!exists;
        }
        ...
    }
}
```

关键点：

- 先从 Redis key 里取出 `SBChain`
- 如果 key 不存在或类型不对，直接视为不存在
- 存在时调用 `SBChain_Check(sb, s, n)`

所以 `BF.EXISTS` 的命令层本身很薄，真正逻辑在 `SBChain_Check`。

## 4. BF.ADD

`BF.ADD` 和 `BF.MADD` 都走 `BFAdd_RedisCommand`：

```c
static int BFAdd_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc) {
    BFInsertOptions options = {
        .capacity = rm_config.bf_initial_size.value,
        .error_rate = rm_config.bf_error_rate.value,
        .autocreate = 1,
        .expansion = rm_config.bf_expansion_factor.value,
        .nonScaling = rm_config.bf_expansion_factor.value == 0 ? BLOOM_OPT_NO_SCALING : 0,
    };
    options.is_multi = isMulti(argv[0]);

    return bfInsertCommon(ctx, argv[1], argv + 2, argc - 2, &options);
}
```

真正处理在 `bfInsertCommon`：

```c
static int bfInsertCommon(...) {
    RedisModuleKey *key = RedisModule_OpenKey(ctx, keystr, REDISMODULE_READ | REDISMODULE_WRITE);
    SBChain *sb;
    const int status = bfGetChain(key, &sb);

    if (status == SB_EMPTY && options->autocreate) {
        sb = bfCreateChain(...);
    } else if (status != SB_OK) {
        return RedisModule_ReplyWithError(ctx, statusStrerror(status));
    }

    for (size_t ii = 0; ii < nitems && rv != -2; ++ii) {
        const char *s = RedisModule_StringPtrLen(items[ii], &n);
        rv = SBChain_Add(sb, s, n);
        ...
    }
}
```

关键点：

- `BF.ADD` 支持自动建 filter
- 每个元素真正写入都走 `SBChain_Add`
- 返回值语义是：
  - `1`：新加入
  - `0`：之前已存在或发生碰撞

## 5. SBChain 这一层

RedisBloom 不是只维护一个 Bloom filter，而是维护一个 `SBChain`。

`SBChain_Check`：

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

含义：

- 先算 hash
- 再从后往前遍历链上的每个 filter
- 任意一层命中就返回存在

`SBChain_Add`：

```c
int SBChain_Add(SBChain *sb, const void *data, size_t len) {
    bloom_hashval h = SBChain_GetHash(sb, data, len);
    for (int ii = sb->nfilters - 1; ii >= 0; --ii) {
        if (bloom_check_h(&sb->filters[ii].inner, h)) {
            return 0;
        }
    }

    SBLink *cur = CUR_FILTER(sb);
    if (cur->size >= cur->inner.entries) {
        if (sb->options & BLOOM_OPT_NO_SCALING) {
            return -2;
        }
        double error = cur->inner.error * ERROR_TIGHTENING_RATIO;
        if (SBChain_AddLink(sb, cur->inner.entries * (size_t)sb->growth, error) != 0) {
            return -1;
        }
        cur = CUR_FILTER(sb);
    }

    int rv = SBChain_AddToLink(cur, h);
    if (rv) {
        sb->size++;
    }
    return rv;
}
```

这层代码体现了 RedisBloom 的两个设计点：

1. 先查后加
- 如果链上任何一层已经命中，直接返回 0
- 不会重复写入

2. 支持扩容
- 当前 filter 满了时，不是直接失败
- 而是新增一个更大的 filter 挂到链尾
- 新层 error rate 会继续收紧：`cur->inner.error * 0.5`

所以 RedisBloom 的 Bloom Filter 实际是：

- 一个可扩展的 filter 链
- 不是单块固定大小位图

## 6. 底层 bloom.c

真正的布隆判断和写位在 `deps/bloom/bloom.c`。

### 6.1 hash 计算

```c
bloom_hashval bloom_calc_hash64(const void *buffer, int len) {
    bloom_hashval rv;
    rv.a = MurmurHash64A_Bloom(buffer, len, 0xc6a4a7935bd1e995ULL);
    rv.b = MurmurHash64A_Bloom(buffer, len, rv.a);
    return rv;
}
```

它不是为每个 hash 函数单独算完整 hash，而是先算出两个基值 `a` 和 `b`，后面用：

- `x = (a + i * b) % mod`

派生出多个 hash 位置。

### 6.2 check / add 统一核心

核心宏：

```c
#define CHECK_ADD_FUNC(T, modExp) \
    T i; \
    int found_unset = 0; \
    const register T mod = modExp; \
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

这段代码很关键：

- `MODE_READ` 下：
  - 只要有一个 bit 没置位，就立即返回不存在
  - 全部 bit 都已置位，返回存在

- `MODE_WRITE` 下：
  - 检查每个 bit，并顺手置位
  - 只要发现过未置位 bit，说明本次是新加入

也就是说：

- check 和 add 复用了同一套核心逻辑
- add 本身就天然包含“是否已存在”的判断

### 6.3 bloom_check_h / bloom_add_h

```c
int bloom_check_h(const struct bloom *bloom, bloom_hashval hash) {
    ...
    return bloom_check_add64((void *)bloom, hash, MODE_READ);
}

int bloom_add_h(struct bloom *bloom, bloom_hashval hash) {
    ...
    return !bloom_check_add64(bloom, hash, MODE_WRITE);
}
```

语义要特别注意：

- `bloom_check_h(...) == 1`：可能存在
- `bloom_add_h(...) == 0`：原来不存在，本次已新增
- `bloom_add_h(...) == 1`：原来就存在，或发生碰撞

这个返回值语义也对应了 `bloom.h` 的说明：

- `0` - element was not present and was added
- `1` - element (or a collision) had already been added previously

## 7. bloom_init

初始化逻辑在 `bloom_init`：

```c
int bloom_init(struct bloom *bloom, uint64_t entries, double error, unsigned options) {
    ...
    bloom->bpe = calc_bpe(error);
    ...
    bloom->hashes = (int)ceil(LN2 * bloom->bpe);
    bloom->bf = (unsigned char *)BLOOM_TRYCALLOC(bloom->bytes, sizeof(unsigned char));
    ...
}
```

这里主要做三件事：

- 根据 error rate 算每个元素需要的 bit 数 `bpe`
- 根据容量和策略算总 bit 数 / byte 数
- 计算需要多少个 hash 函数 `hashes`

所以 `BF.RESERVE error_rate capacity` 最终会落到：

- 位图大小
- hash 函数数量
- 可容纳元素规模

## 8. 结论

如果只看 `BF.ADD` / `BF.EXISTS`，RedisBloom 的调用链可以压缩成：

1. `rebloom.c`
- 注册命令
- 解析 Redis 参数
- 调用 `SBChain_Add` / `SBChain_Check`

2. `sb.c`
- 维护一条可扩容 Bloom 链
- 负责“先查后加”和扩容策略

3. `bloom.c`
- 计算 hash
- 映射 bit 位
- 执行真正的 check / add

对当前项目最重要的理解只有两点：

- `BF.EXISTS` 返回的只是“可能存在”
- `BF.ADD` / `BF.EXISTS` 背后不是单一 filter，而是 RedisBloom 自己维护的一条可扩容 Bloom chain

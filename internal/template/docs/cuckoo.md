# 布谷鸟过滤器

## 1. 为什么会引入布谷鸟过滤器

布隆过滤器的主要问题之一是：

- 普通实现不支持安全删除

原因是它只记录 bit 是否为 `1`，不知道这些 bit 是由哪些元素置上的。

如果直接清零某些 bit，就可能把别的元素一起“误删”。

为了解决“支持删除”的问题，工程上常见方案之一就是：

- 布谷鸟过滤器，Cuckoo Filter

## 2. 核心思路

布谷鸟过滤器不存完整元素，也不存共享 bit，而是：

- 对元素计算一个短指纹 `fingerprint`
- 为它计算两个候选桶位置
- 把指纹放进两个桶中的一个

如果两个桶都满了，就触发“踢出重排”：

- 先把某个桶里的旧指纹踢出来
- 把新指纹放进去
- 再把被踢出的指纹尝试放到它的另一个候选桶
- 如此循环，直到插入成功或达到最大尝试次数

这就是“布谷鸟”名字的来源。

## 3. 与布隆过滤器的区别

布隆过滤器：

- 记录 bit 是否被置位
- 适合只增不删
- 实现简单

布谷鸟过滤器：

- 记录元素的短指纹
- 支持查询、插入、删除
- 插入逻辑更复杂

所以两者的本质区别是：

- Bloom 是“位图命中”
- Cuckoo 是“桶里是否有这个指纹”

## 4. 为什么能删除

布谷鸟过滤器能删除，是因为它存的是具体的指纹条目，而不是共享 bit。

删除时只需要：

- 找到两个候选桶
- 在桶里定位到目标指纹
- 把该槽位置空

这样不会像 Bloom 那样影响其它元素。

## 5. RedisBloom 源码位置

当前仓库里 RedisBloom 的布谷鸟实现主要在：

- `RedisBloom/src/cuckoo.h`
- `RedisBloom/src/cuckoo.c`
- `RedisBloom/src/rebloom.c`

其中：

- `rebloom.c` 负责 Redis 命令层
- `cuckoo.c` 负责布谷鸟过滤器核心逻辑

## 6. RedisBloom 命令入口

布谷鸟命令同样在 `rebloom.c` 注册：

```c
RegisterCommand(ctx, "cf.reserve", CFReserve_RedisCommand, "write deny-oom", "write fast");
RegisterCommand(ctx, "cf.add", CFAdd_RedisCommand, "write deny-oom", "write");
RegisterCommand(ctx, "cf.addnx", CFAdd_RedisCommand, "write deny-oom", "write");
RegisterCommand(ctx, "cf.insert", CFInsert_RedisCommand, "write deny-oom", "write");
RegisterCommand(ctx, "cf.insertnx", CFInsert_RedisCommand, "write deny-oom", "write");
RegisterCommand(ctx, "cf.exists", CFCheck_RedisCommand, "readonly fast", "read");
RegisterCommand(ctx, "cf.mexists", CFCheck_RedisCommand, "readonly fast", "read");
RegisterCommand(ctx, "cf.count", CFCheck_RedisCommand, "readonly fast", "read");
RegisterCommand(ctx, "cf.del", CFDel_RedisCommand, "write fast", "write");
```

所以 Cuckoo 这一层对 Redis 暴露的是：

- `CF.ADD` / `CF.INSERT`
- `CF.EXISTS`
- `CF.DEL`

这也正对应了它相比 Bloom 多出来的能力：

- 删除

## 7. 候选桶位置如何计算

RedisBloom 里的关键结构在 [cuckoo.h](/home/fanqicheng/project/jx/mysurl1/RedisBloom/src/cuckoo.h:57)：

```c
typedef struct {
    uint64_t i1;
    uint64_t i2;
    CuckooFingerprint fp;
} CuckooKey;
```

真正的计算逻辑在 [cuckoo.c](/home/fanqicheng/project/jx/mysurl1/RedisBloom/src/cuckoo.c:110)：

```c
static CuckooHash getAltHash(CuckooFingerprint fp, CuckooHash index) {
    return ((CuckooHash)(index ^ ((CuckooHash)fp * 0x5bd1e995)));
}

static void getLookupParams(CuckooHash hash, LookupParams *params) {
    params->fp = (CuckooFingerprint)(hash % 255 + 1);
    params->h1 = hash;
    params->h2 = getAltHash(params->fp, params->h1);
}
```

这里的逻辑是：

1. 先对原始元素算一个大 hash
2. 从 hash 中截出一个 8 位指纹：
   - `fp = hash % 255 + 1`
3. 第一个候选位置来自原始 hash：
   - `h1 = hash`
4. 第二个候选位置由“第一个位置 + 指纹”推导出来：
   - `h2 = h1 ^ (fp * 0x5bd1e995)`

真正映射到桶时，再对桶数取模：

```c
static uint32_t SubCF_GetIndex(const SubCF *subCF, CuckooHash hash) {
    return (hash % subCF->numBuckets) * subCF->bucketSize;
}
```

所以可以理解成：

- `bucket1 = h1 % numBuckets`
- `bucket2 = h2 % numBuckets`

## 8. 为什么这个公式适合布谷鸟

关键在于可逆性。

因为：

```text
h2 = h1 ^ delta
h1 = h2 ^ delta
```

只要知道：

- 当前桶位置
- 指纹 `fp`

就能算出另一个候选桶。

这对“踢出重排”非常重要，因为被踢出来的元素手里通常只剩：

- 当前桶
- 它自己的指纹

## 9. 查询

查询逻辑很直接，在 [cuckoo.c](/home/fanqicheng/project/jx/mysurl1/RedisBloom/src/cuckoo.c:145)：

```c
static int Filter_Find(const SubCF *filter, const LookupParams *params) {
    uint8_t bucketSize = filter->bucketSize;
    uint64_t loc1 = SubCF_GetIndex(filter, params->h1);
    uint64_t loc2 = SubCF_GetIndex(filter, params->h2);
    return Bucket_Find(&filter->data[loc1], bucketSize, params->fp) != NULL ||
           Bucket_Find(&filter->data[loc2], bucketSize, params->fp) != NULL;
}
```

也就是：

- 只查两个候选桶
- 看指纹是否在其中任意一个桶里

## 10. 删除

删除逻辑也只需要查这两个桶，在 [cuckoo.c](/home/fanqicheng/project/jx/mysurl1/RedisBloom/src/cuckoo.c:163)：

```c
static int Filter_Delete(const SubCF *filter, const LookupParams *params) {
    uint8_t bucketSize = filter->bucketSize;
    uint64_t loc1 = SubCF_GetIndex(filter, params->h1);
    uint64_t loc2 = SubCF_GetIndex(filter, params->h2);
    return Bucket_Delete(&filter->data[loc1], bucketSize, params->fp) ||
           Bucket_Delete(&filter->data[loc2], bucketSize, params->fp);
}
```

桶里找到指纹后，直接置空：

```c
bucket[ii] = CUCKOO_NULLFP;
```

这就是它能支持删除的根本原因。

## 11. 插入与踢出

插入时会先尝试在两个候选桶里找空槽：

```c
if ((slot = Bucket_FindAvailable(&filter->data[loc1], bucketSize)) ||
    (slot = Bucket_FindAvailable(&filter->data[loc2], bucketSize))) {
    return slot;
}
```

如果两个桶都没有空位，就进入踢出逻辑 `Filter_KOInsert`：

```c
while (counter++ < maxIterations) {
    uint8_t *bucket = &curFilter->data[ii * bucketSize];
    swapFPs(bucket + victimIx, &fp);
    ii = getAltHash(fp, ii) % numBuckets;
    uint8_t *empty = Bucket_FindAvailable(&curFilter->data[ii * bucketSize], bucketSize);
    if (empty) {
        *empty = fp;
        return CuckooInsert_Inserted;
    }
    victimIx = (victimIx + 1) % bucketSize;
}
```

含义：

1. 从当前桶里选一个 victim
2. 把新指纹塞进去
3. 被踢出的旧指纹根据 `getAltHash(fp, ii)` 算出另一个桶
4. 去新桶找空位
5. 如果还满，就继续踢

如果一直踢不进去，就认为当前 filter 放不下，后续再走扩容。

## 12. RedisBloom 里的扩容

RedisBloom 的 Cuckoo Filter 也不是单块固定结构，而是支持扩容。

初始化时：

```c
filter->numBuckets = getNextN2(capacity / bucketSize);
```

扩容时：

```c
size_t growth = pow(filter->expansion, filter->numFilters);
currentFilter->numBuckets = filter->numBuckets * growth;
```

也就是说：

- 底层会维护多个 `SubCF`
- 新的子过滤器会更大
- 插入失败时可以增长出新的 sub filter

## 13. 优点和代价

优点：

- 支持删除
- 查询只看两个桶，速度快
- 空间效率高

代价：

- 实现比 Bloom 更复杂
- 插入可能需要多次踢出
- 高负载时插入失败处理更复杂

## 14. 总结

如果压成一句话，布谷鸟过滤器的底层原理就是：

- 用短指纹代替完整元素
- 用两个候选桶定位元素
- 用踢出重排解决冲突

RedisBloom 这份实现又在此基础上加了：

- Redis 命令封装
- 多个子过滤器扩容
- 删除与压缩逻辑

所以它比布隆过滤器更复杂，但换来了一个很重要的能力：

- 可以删除

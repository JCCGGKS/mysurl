# 游标分页方案梳理

## 1. 目标

当前项目的短链列表和用户操作日志列表都没有使用 `page + page_size` 偏移分页，而是使用 `last_id + limit` 的游标分页。

这样做的目标是：

- 避免深分页时 `offset` 越来越大带来的性能问题
- 在按主键顺序翻页时，查询条件更简单
- 前后端只需要围绕 `last_id`、`has_more`、`next_last_id` 协作

## 2. 当前接口约定

### 2.1 短链列表

请求定义见 `internal/schema/short.go`：

```go
type ListUserLinksRequest struct {
	LastID      uint64 `form:"last_id,optional"`
	Limit       int    `form:"limit,optional"`
	ShortCode   string `form:"short_code,optional"`
	OriginalURL string `form:"original_url,optional"`
}
```

响应定义：

```go
type UserLinkListResponse struct {
	Items      []UserLinkItem `json:"items"`
	Total      int64          `json:"total"`
	Limit      int            `json:"limit"`
	HasMore    bool           `json:"has_more"`
	NextLastID uint64         `json:"next_last_id"`
}
```

### 2.2 用户操作日志列表

请求定义见 `internal/schema/short.go`：

```go
type ListUserOperationLogsRequest struct {
	LastID uint64 `form:"last_id,optional"`
	Limit  int    `form:"limit,optional"`
	Action string `form:"action,optional"`
}
```

响应定义：

```go
type UserOperationLogListResponse struct {
	Items      []UserOperationLogItem `json:"items"`
	Total      int64                  `json:"total"`
	Limit      int                    `json:"limit"`
	HasMore    bool                   `json:"has_more"`
	NextLastID uint64                 `json:"next_last_id"`
}
```

结论很简单：前端真正参与分页的字段只有 3 个：

- `last_id`
- `has_more`
- `next_last_id`

## 3. 后端实现方式

### 3.1 逻辑层职责

短链列表逻辑见 `internal/logic/listuserlinkslogic.go`，操作日志逻辑见 `internal/logic/listuseroperationlogslogic.go`。

逻辑层主要做 4 件事：

1. 读取筛选条件和 `limit`
2. 先查询总数 `total`
3. 再按游标查列表数据，实际查询 `limit + 1`
4. 组装 `has_more` 和 `next_last_id`

### 3.2 为什么查 `limit + 1`

以短链列表为例：

```go
records, err := l.svcCtx.ShortLinkDAO.ListByUserIDWithCursor(
	l.ctx,
	claims.UserID,
	shortCode,
	originalURL,
	req.LastID,
	limit+1,
)
```

作用是探测是否还有下一页。

例如：

- 前端要求每页 10 条
- 后端实际查 11 条
- 如果只查到 8 条，说明没有下一页
- 如果查到 11 条，说明至少还有下一页

随后逻辑层会：

- 用 `len(records) > limit` 判断 `has_more`
- 如果有多出来的 1 条，只返回前 `limit` 条给前端
- 把当前返回结果最后一条记录的 `id` 作为 `next_last_id`

### 3.3 DAO / SQL 配合

短链 DAO 见 `internal/dao/shortlinkdao.go`，操作日志 DAO 见 `internal/dao/useroperationlogdao.go`。

操作日志查询的核心 SQL 结构最直观：

```sql
SELECT ...
FROM user_operation_logs
WHERE user_id = ?
  AND action = ?
  AND id > ?
ORDER BY id ASC
LIMIT ?
```

短链列表也是同样思路：

- 先按 `user_id` 过滤
- 如果前端传了 `short_code`，就追加这个筛选
- 如果前端传了 `original_url`，就追加这个筛选
- 如果前端传了 `last_id`，就追加 `id > ?`
- 最后 `ORDER BY id ASC LIMIT ?`

这里的 `id > last_id` 就是游标分页的核心。

## 4. 前端如何传递分页参数

### 4.1 前端不会把“页码”传给后端

前端页码按钮只是 UI 展示，真正发给后端的仍然是 `last_id`。

以短链列表页 `vue/src/views/links/LinkListView.vue` 为例，请求代码是：

```js
const params = new URLSearchParams({
  limit: String(limit.value),
})

if (currentCursor.value > 0) {
  params.set('last_id', String(currentCursor.value))
}

if (shortCode.value.trim()) {
  params.set('short_code', shortCode.value.trim())
}

if (originalUrl.value.trim()) {
  params.set('original_url', originalUrl.value.trim())
}
```

所以后端根本不知道“用户点击的是第几页”，后端只知道：

- 这次查多少条 `limit`
- 从哪个游标之后开始查 `last_id`

操作日志页 `vue/src/views/users/UserOperationLogView.vue` 也是同样结构，只是筛选参数变成了 `action`。

## 5. 前端状态设计

当前实现不是“双栈”，而是：

- 一个当前游标：`currentCursor`
- 一个下一页游标：`nextCursor`
- 一个历史游标栈：`cursorHistory`
- 一个仅用于显示的页码：`currentPage`

以短链列表页 `vue/src/views/links/LinkListView.vue` 为例：

```js
const currentPage = ref(1)
const limit = ref(10)
const currentCursor = ref(0)
const nextCursor = ref(0)
const hasMore = ref(false)
const cursorHistory = ref([])
```

### 5.1 首次加载

第一页时：

- `currentCursor = 0`
- 不传 `last_id`
- `cursorHistory = []`

### 5.2 点击下一页

实现代码：

```js
function goNextPage() {
  if (!hasMore.value || loading.value) return
  cursorHistory.value.push(currentCursor.value)
  currentCursor.value = nextCursor.value
  currentPage.value += 1
  loadLinks()
}
```

过程是：

1. 先把当前页游标压入历史栈
2. 再把后端返回的 `nextCursor` 赋给 `currentCursor`
3. 用新的 `currentCursor` 请求下一页

### 5.3 点击上一页

实现代码：

```js
function goPrevPage() {
  if (currentPage.value <= 1 || loading.value) return
  currentCursor.value = cursorHistory.value.pop() ?? 0
  currentPage.value -= 1
  loadLinks()
}
```

过程是：

1. 从 `cursorHistory` 弹出上一个游标
2. 回填到 `currentCursor`
3. 重新请求上一页

## 6. 页码按钮是怎么工作的

页码按钮不是把 `page=2` 直接传给后端，而是靠前端维护过的游标状态来模拟。

当前逻辑见 `vue/src/views/links/LinkListView.vue` 和 `vue/src/views/users/UserOperationLogView.vue`。

关键点：

- 可以跳到当前页前一页：直接走 `goPrevPage()`
- 可以跳到当前页后一页：直接走 `goNextPage()`
- 可以跳回更早访问过的页：从 `cursorHistory` 中取对应游标
- 不能直接跳到一个从未访问过的更远页，因为前端并不知道那一页对应的 `last_id`

因此当前页码按钮本质上只是：

- 已访问页的回跳入口
- 相邻下一页的前进入口

不是真正意义上的“任意页码跳转”。

## 7. 一个具体例子

假设每页 10 条。

### 第 1 页

- `currentCursor = 0`
- 请求不带 `last_id`
- 后端返回 `next_last_id = 10`
- `cursorHistory = []`

### 点下一页到第 2 页

- 先 `cursorHistory.push(0)`
- 再 `currentCursor = 10`
- 请求带 `last_id=10`
- 后端返回 `next_last_id = 20`
- 此时 `cursorHistory = [0]`

### 再点下一页到第 3 页

- 先 `cursorHistory.push(10)`
- 再 `currentCursor = 20`
- 请求带 `last_id=20`
- 后端返回 `next_last_id = 30`
- 此时 `cursorHistory = [0, 10]`

### 从第 3 页回第 2 页

- `pop()` 得到 `10`
- `currentCursor = 10`
- 再请求 `last_id=10`

### 从第 3 页点击页码 1

前端会取：

```js
currentCursor.value = page === 1 ? 0 : (cursorHistory.value[page - 1] ?? 0)
```

也就是：

- 第 1 页对应游标 `0`
- 第 2 页对应游标 `10`

这里能成立，是因为这些游标都已经在前面的访问过程中被前端记住了。

## 8. 前后端配合的边界

### 8.1 适合的场景

- 按创建顺序或主键顺序翻页
- 只需要上一页、下一页
- 数据量较大，希望避免深分页性能问题

### 8.2 当前限制

- 不适合直接跳任意页
- 页码按钮只能基于“已访问页 + 紧邻下一页”工作
- 如果后续要支持 `1, 2, 3, ... 100` 任意跳转，当前纯游标分页方案不够直接

## 9. 如果后续要增强

有两个常见方向：

### 9.1 保持游标分页

前端增加更明确的页码到游标映射，例如：

- `pageToCursor: Map<number, number>`

这样可以比现在的数组索引方式更清晰，但本质限制不变：

- 只能跳到已经拿到过游标的页

### 9.2 改成偏移分页

后端改成支持：

- `page`
- `page_size`

SQL 变为：

```sql
LIMIT ? OFFSET ?
```

这样可以直接跳任意页，但深分页性能会比游标分页差。

## 10. 当前项目的结论

当前项目的短链列表和操作日志列表，采用的是：

- 后端：`last_id + limit` 游标分页
- 前端：`currentCursor + nextCursor + cursorHistory` 管理翻页状态
- 页码按钮：基于历史游标模拟，不是把页码直接传给后端

所以这套方案更准确的描述应当是：

**游标分页 + 前端历史栈回退**

而不是传统意义上的“完整页码分页”。

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { getJson } from '../../services/api'

const loading = ref(false)
const errorMessage = ref('')
const logs = ref([])
const total = ref(0)
const currentPage = ref(1)
const limit = ref(10)
const currentCursor = ref(0)
const nextCursor = ref(0)
const hasMore = ref(false)
const cursorHistory = ref([])
const pageSizeOptions = [10, 20, 50]
const selectedAction = ref('')
const selectedResult = ref('')

const actionOptions = [
  { value: '', label: '全部类型' },
  { value: 'login', label: '登录' },
  { value: 'create_link', label: '创建短链' },
  { value: 'create_link_batch', label: '批量创建短链' },
]

const resultOptions = [
  { value: '', label: '全部结果' },
  { value: 'success', label: '成功' },
  { value: 'partial_success', label: '部分成功' },
  { value: 'failed', label: '失败' },
]

const pageLabel = computed(() => {
  if (total.value === 0) return '0 / 0'
  const start = (currentPage.value - 1) * limit.value + 1
  const end = Math.min(start + logs.value.length - 1, total.value)
  return `${start}-${end} / ${total.value}`
})

const totalPages = computed(() => {
  if (total.value === 0) return 0
  return Math.ceil(total.value / limit.value)
})

const pageButtons = computed(() => {
  const reachablePages = Math.min(totalPages.value, currentPage.value + (hasMore.value ? 1 : 0))
  return Array.from({ length: reachablePages }, (_, index) => index + 1)
})

onMounted(() => {
  loadLogs()
})

watch(limit, (value, oldValue) => {
  if (value === oldValue) return
  resetPagination()
  loadLogs()
})

async function loadLogs() {
  loading.value = true
  errorMessage.value = ''

  try {
    const params = new URLSearchParams({
      limit: String(limit.value),
    })
    if (currentCursor.value > 0) {
      params.set('last_id', String(currentCursor.value))
    }
    if (selectedAction.value) {
      params.set('action', selectedAction.value)
    }
    if (selectedResult.value) {
      params.set('result', selectedResult.value)
    }

    const data = await getJson(`/api/v1/user-operation-logs?${params.toString()}`, { auth: true })
    logs.value = Array.isArray(data.items) ? data.items : []
    total.value = Number(data.total || 0)
    limit.value = Number(data.limit || limit.value)
    hasMore.value = Boolean(data.has_more)
    nextCursor.value = Number(data.next_last_id || 0)
  } catch (error) {
    errorMessage.value = error.message || '加载操作日志失败，请稍后重试'
    logs.value = []
    total.value = 0
    hasMore.value = false
    nextCursor.value = 0
  } finally {
    loading.value = false
  }
}

function resetPagination() {
  currentPage.value = 1
  currentCursor.value = 0
  nextCursor.value = 0
  hasMore.value = false
  cursorHistory.value = []
}

function refreshLogs() {
  resetPagination()
  loadLogs()
}

function applyFilters() {
  resetPagination()
  loadLogs()
}

function resetFilters() {
  selectedAction.value = ''
  selectedResult.value = ''
  limit.value = 10
  resetPagination()
  loadLogs()
}

function goPrevPage() {
  if (currentPage.value <= 1 || loading.value) return
  currentCursor.value = cursorHistory.value.pop() ?? 0
  currentPage.value -= 1
  loadLogs()
}

function goNextPage() {
  if (!hasMore.value || loading.value) return
  cursorHistory.value.push(currentCursor.value)
  currentCursor.value = nextCursor.value
  currentPage.value += 1
  loadLogs()
}

function goToPage(page) {
  if (loading.value || page === currentPage.value) return
  if (page < 1 || page > totalPages.value) return

  if (page === currentPage.value - 1) {
    goPrevPage()
    return
  }

  if (page === currentPage.value + 1) {
    goNextPage()
    return
  }

  if (page < currentPage.value) {
    currentCursor.value = page === 1 ? 0 : (cursorHistory.value[page - 1] ?? 0)
    currentPage.value = page
    cursorHistory.value = cursorHistory.value.slice(0, Math.max(0, page - 1))
    loadLogs()
  }
}

function formatAction(action) {
  if (action === 'login') return '登录'
  if (action === 'create_link') return '创建短链'
  if (action === 'create_link_batch') return '批量创建短链'
  return action || '--'
}

function formatResult(result) {
  if (result === 'success') return '成功'
  if (result === 'partial_success') return '部分成功'
  if (result === 'failed') return '失败'
  return result || '--'
}

function formatDate(timestamp) {
  if (!timestamp) return '--'

  const date = new Date(timestamp * 1000)
  if (Number.isNaN(date.getTime())) {
    return '--'
  }

  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  }).format(date)
}
</script>

<template>
  <section class="dashboard-page dashboard-page-wide">
    <header class="dashboard-page-head">
      <div class="workspace-head">
        <div>
          <h2>用户操作日志</h2>
        </div>
        <p class="workspace-note">总数：{{ total }}</p>
      </div>
    </header>

    <section class="filter-panel operation-log-toolbar">
      <div class="operation-log-toolbar-row">
        <label class="filter-field operation-log-size-field">
          <span class="field-label">类型</span>
          <select v-model="selectedAction" class="text-input filter-select" :disabled="loading">
            <option v-for="option in actionOptions" :key="option.value" :value="option.value">
              {{ option.label }}
            </option>
          </select>
        </label>

        <label class="filter-field operation-log-size-field">
          <span class="field-label">结果</span>
          <select v-model="selectedResult" class="text-input filter-select" :disabled="loading">
            <option v-for="option in resultOptions" :key="option.value" :value="option.value">
              {{ option.label }}
            </option>
          </select>
        </label>

        <div class="filter-actions">
          <button class="secondary-button" type="button" @click="applyFilters" :disabled="loading">
            {{ loading ? '查询中...' : '查询' }}
          </button>
          <button class="ghost-link" type="button" @click="resetFilters" :disabled="loading">
            重置
          </button>
          <button class="ghost-link" type="button" @click="refreshLogs" :disabled="loading">
            {{ loading ? '刷新中...' : '刷新列表' }}
          </button>
        </div>
      </div>
    </section>

    <p v-if="errorMessage" class="feedback feedback-error" role="alert">
      {{ errorMessage }}
    </p>

    <div v-if="loading" class="list-empty">
      <h3>正在加载用户操作日志。</h3>
    </div>

    <div v-else-if="logs.length === 0" class="list-empty">
      <h3>当前还没有操作日志。</h3>
      <p>登录成功和创建短链成功后，这里会出现最新记录。</p>
    </div>

    <section v-else class="list-content">
      <div class="link-table-shell">
        <table class="link-table operation-log-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>时间</th>
              <th>类型</th>
              <th>结果</th>
              <th>备注</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in logs" :key="item.id">
              <td>{{ item.id }}</td>
              <td>{{ formatDate(item.created_at) }}</td>
              <td>{{ formatAction(item.action) }}</td>
              <td>
                <span class="log-result-badge">{{ formatResult(item.result) }}</span>
              </td>
              <td>
                <span class="batch-error-text">{{ item.reason || '--' }}</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="pagination-bar">
        <p class="pagination-meta">显示 {{ pageLabel }}</p>
        <div class="pagination-actions">
          <span class="pagination-current">总: {{ totalPages }} 页</span>
          <button
            class="ghost-link pagination-button pagination-arrow"
            type="button"
            @click="goPrevPage"
            :disabled="currentPage <= 1 || loading"
            aria-label="上一页"
          >
            &lt;
          </button>
          <div class="pagination-page-list">
            <button
              v-for="page in pageButtons"
              :key="page"
              class="ghost-link pagination-button pagination-page-number"
              :class="{ 'pagination-page-number-active': page === currentPage }"
              type="button"
              @click="goToPage(page)"
              :disabled="loading || page > currentPage + 1"
            >
              {{ page }}
            </button>
          </div>
          <button
            class="ghost-link pagination-button pagination-arrow"
            type="button"
            @click="goNextPage"
            :disabled="!hasMore || loading"
            aria-label="下一页"
          >
            &gt;
          </button>
          <label class="filter-field operation-log-page-size-field">
            <select v-model.number="limit" class="text-input filter-select pagination-size-select" :disabled="loading">
              <option v-for="size in pageSizeOptions" :key="size" :value="size">每页条数：{{ size }}</option>
            </select>
          </label>
        </div>
      </div>
    </section>
  </section>
</template>

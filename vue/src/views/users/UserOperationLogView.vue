<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { getJson } from '../../services/api'
import { handleUnauthorized } from '../../router'

const router = useRouter()

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

const pageLabel = computed(() => {
  if (total.value === 0) return '0 / 0'
  const start = (currentPage.value - 1) * limit.value + 1
  const end = Math.min(start + logs.value.length - 1, total.value)
  return `${start}-${end} / ${total.value}`
})

const latestActionTime = computed(() => {
  if (logs.value.length === 0) return '--'
  return formatDate(logs.value[logs.value.length - 1]?.created_at || 0)
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

    const data = await getJson(`/api/v1/user-operation-logs?${params.toString()}`, { auth: true })
    logs.value = Array.isArray(data.items) ? data.items : []
    total.value = Number(data.total || 0)
    limit.value = Number(data.limit || limit.value)
    hasMore.value = Boolean(data.has_more)
    nextCursor.value = Number(data.next_last_id || 0)
  } catch (error) {
    if (error.status === 401) {
      handleUnauthorized(router)
      return
    }
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

function formatResult(result) {
  if (result === 'success') return '成功'
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

    <section class="user-grid operation-log-summary-grid">
      <article class="user-card operation-log-summary-card">
        <span class="signal-title">当前页记录</span>
        <strong>{{ logs.length }}</strong>
      </article>
      <article class="user-card operation-log-summary-card">
        <span class="signal-title">当前页结果</span>
        <strong>{{ logs.length > 0 ? '成功' : '--' }}</strong>
      </article>
      <article class="user-card operation-log-summary-card">
        <span class="signal-title">最新时间</span>
        <strong>{{ latestActionTime }}</strong>
      </article>
    </section>

    <section class="filter-panel operation-log-toolbar">
      <div class="operation-log-toolbar-row">
        <label class="filter-field operation-log-size-field">
          <span class="field-label">每页条数</span>
          <select v-model.number="limit" class="text-input filter-select" :disabled="loading">
            <option v-for="size in pageSizeOptions" :key="size" :value="size">{{ size }} 条</option>
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
              <th>结果</th>
              <th>备注</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in logs" :key="item.id">
              <td>{{ item.id }}</td>
              <td>{{ formatDate(item.created_at) }}</td>
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
        <div class="pagination-actions">
          <button
            class="ghost-link pagination-button"
            type="button"
            @click="goPrevPage"
            :disabled="currentPage <= 1 || loading"
          >
            上一页
          </button>
          <span class="pagination-current">第 {{ currentPage }} 页</span>
          <button
            class="ghost-link pagination-button"
            type="button"
            @click="goNextPage"
            :disabled="!hasMore || loading"
          >
            下一页
          </button>
        </div>
        <p class="pagination-meta">显示 {{ pageLabel }}</p>
      </div>
    </section>
  </section>
</template>

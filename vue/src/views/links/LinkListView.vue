<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { getJson } from '../../services/api'
import { handleUnauthorized } from '../../router'

const router = useRouter()

const loading = ref(false)
const copiedId = ref(0)
const errorMessage = ref('')
const links = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const shortCode = ref('')
const originalUrl = ref('')
const pageSizeOptions = [10, 20, 50]

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize.value)))
const pageLabel = computed(() => {
  if (total.value === 0) return '0 / 0'
  const start = (page.value - 1) * pageSize.value + 1
  const end = Math.min(page.value * pageSize.value, total.value)
  return `${start}-${end} / ${total.value}`
})

onMounted(() => {
  loadLinks()
})

async function loadLinks() {
  loading.value = true
  errorMessage.value = ''

  try {
    const params = new URLSearchParams({
      page: String(page.value),
      page_size: String(pageSize.value),
    })
    if (shortCode.value.trim()) {
      params.set('short_code', shortCode.value.trim())
    }
    if (originalUrl.value.trim()) {
      params.set('original_url', originalUrl.value.trim())
    }

    const data = await getJson(`/api/v1/links/mine?${params.toString()}`, { auth: true })
    links.value = Array.isArray(data.items) ? data.items : []
    total.value = Number(data.total || 0)
    page.value = Number(data.page || page.value)
    pageSize.value = Number(data.page_size || pageSize.value)
  } catch (error) {
    if (error.status === 401) {
      handleUnauthorized(router)
      return
    }
    errorMessage.value = error.message || '加载短链列表失败，请稍后重试'
    links.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

function applyFilters() {
  page.value = 1
  loadLinks()
}

function resetFilters() {
  shortCode.value = ''
  originalUrl.value = ''
  page.value = 1
  pageSize.value = 10
  loadLinks()
}

function goPrevPage() {
  if (page.value <= 1 || loading.value) return
  page.value -= 1
  loadLinks()
}

function goNextPage() {
  if (page.value >= totalPages.value || loading.value) return
  page.value += 1
  loadLinks()
}

watch(pageSize, (value, oldValue) => {
  if (value === oldValue) return
  page.value = 1
  loadLinks()
})

async function copyShortUrl(item) {
  errorMessage.value = ''

  try {
    await writeToClipboard(item.short_url)
    copiedId.value = item.id
    window.setTimeout(() => {
      if (copiedId.value === item.id) {
        copiedId.value = 0
      }
    }, 1600)
  } catch {
    errorMessage.value = '复制失败，请手动复制短链'
  }
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
  }).format(date)
}

async function writeToClipboard(value) {
  if (navigator.clipboard?.writeText && window.isSecureContext) {
    await navigator.clipboard.writeText(value)
    return
  }

  const textarea = document.createElement('textarea')
  textarea.value = value
  textarea.setAttribute('readonly', 'true')
  textarea.style.position = 'fixed'
  textarea.style.top = '-9999px'
  textarea.style.left = '-9999px'

  document.body.appendChild(textarea)
  textarea.focus()
  textarea.select()
  textarea.setSelectionRange(0, textarea.value.length)

  const copiedWithExecCommand = document.execCommand('copy')
  document.body.removeChild(textarea)

  if (!copiedWithExecCommand) {
    throw new Error('copy failed')
  }
}
</script>

<template>
  <section class="dashboard-page dashboard-page-wide">
    <header class="dashboard-page-head">
      <h2>短链列表</h2>
    </header>

    <section class="filter-panel">
      <div class="filter-row">
        <label class="filter-field">
          <span class="field-label">短码</span>
          <input
            v-model="shortCode"
            class="text-input filter-input"
            type="text"
            placeholder="输入短码筛选"
            :disabled="loading"
            @keyup.enter="applyFilters"
          />
        </label>

        <label class="filter-field">
          <span class="field-label">原始链接</span>
          <input
            v-model="originalUrl"
            class="text-input filter-input"
            type="text"
            placeholder="输入原始链接筛选"
            :disabled="loading"
            @keyup.enter="applyFilters"
          />
        </label>

        <label class="filter-field">
          <span class="field-label">每页条数</span>
          <select v-model.number="pageSize" class="text-input filter-select" :disabled="loading">
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
          <button class="ghost-link" type="button" @click="loadLinks" :disabled="loading">
          {{ loading ? '刷新中...' : '刷新列表' }}
          </button>
        </div>
      </div>
    </section>

    <p v-if="errorMessage" class="feedback feedback-error" role="alert">
      {{ errorMessage }}
    </p>

    <div v-if="loading" class="list-empty">
      <h3>正在加载你的短链列表。</h3>
    </div>

    <div v-else-if="links.length === 0" class="list-empty">
      <h3>当前筛选条件下没有短链。</h3>
    </div>

    <section v-else class="list-content">
      <div class="link-table-shell">
        <table class="link-table">
          <thead>
            <tr>
              <th>短码</th>
              <th>短链</th>
              <th>原始链接</th>
              <th>访问量</th>
              <th>创建时间</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in links" :key="item.id">
              <td>
                <strong class="table-code">{{ item.short_code }}</strong>
              </td>
              <td>
                <a class="table-link" :href="item.short_url" target="_blank" rel="noreferrer">
                  {{ item.short_url }}
                </a>
              </td>
              <td>
                <span class="table-original" :title="item.original_url">{{ item.original_url }}</span>
              </td>
              <td>{{ item.visit_count }}</td>
              <td>{{ formatDate(item.created_at) }}</td>
              <td>
                <div class="table-actions">
                  <button class="secondary-button table-button" type="button" @click="copyShortUrl(item)">
                    {{ copiedId === item.id ? '已复制' : '复制' }}
                  </button>
                  <a class="ghost-link table-button" :href="item.short_url" target="_blank" rel="noreferrer">
                    打开
                  </a>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="pagination-bar">
        <p class="pagination-meta">显示 {{ pageLabel }}</p>
        <div class="pagination-actions">
          <button class="ghost-link pagination-button" type="button" @click="goPrevPage" :disabled="page <= 1 || loading">
            上一页
          </button>
          <span class="pagination-current">第 {{ page }} / {{ totalPages }} 页</span>
          <button
            class="ghost-link pagination-button"
            type="button"
            @click="goNextPage"
            :disabled="page >= totalPages || loading"
          >
            下一页
          </button>
        </div>
      </div>
    </section>
  </section>
</template>

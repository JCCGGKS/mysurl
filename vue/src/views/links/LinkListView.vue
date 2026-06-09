<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { getJson } from '../../services/api'
import { handleUnauthorized } from '../../router'

const router = useRouter()

const loading = ref(false)
const copiedId = ref(0)
const errorMessage = ref('')
const links = ref([])

onMounted(() => {
  loadLinks()
})

async function loadLinks() {
  loading.value = true
  errorMessage.value = ''

  try {
    const data = await getJson('/api/v1/links/mine', { auth: true })
    links.value = Array.isArray(data.items) ? data.items : []
  } catch (error) {
    if (error.status === 401) {
      handleUnauthorized(router)
      return
    }
    errorMessage.value = error.message || '加载短链列表失败，请稍后重试'
    links.value = []
  } finally {
    loading.value = false
  }
}

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
  <section class="dashboard-page">
    <header class="dashboard-page-head">
      <div>
        <p class="section-kicker">Links</p>
        <h2>短链列表</h2>
      </div>
      <p class="workspace-note">GET /api/v1/links/mine</p>
    </header>

    <section class="workspace-card workspace-card-dashboard">
      <div class="page-intro">
        <p class="summary">
          这里只展示当前登录用户创建过的短链。列表数据来自受保护接口，并按创建时间倒序返回。
        </p>
        <div class="inline-signal">
          <span class="signal-title">列表数量</span>
          <strong>{{ loading ? '加载中' : `${links.length} 条` }}</strong>
        </div>
      </div>

      <div class="list-toolbar">
        <button class="secondary-button" type="button" @click="loadLinks" :disabled="loading">
          {{ loading ? '刷新中...' : '刷新列表' }}
        </button>
      </div>

      <p v-if="errorMessage" class="feedback feedback-error" role="alert">
        {{ errorMessage }}
      </p>

      <div v-if="loading" class="list-empty">
        <p class="eyebrow">Loading</p>
        <h3>正在加载你的短链列表。</h3>
      </div>

      <div v-else-if="links.length === 0" class="list-empty">
        <p class="eyebrow">No Links</p>
        <h3>当前账号下还没有短链。</h3>
        <p>先去“创建短链”页面生成第一条短链，随后这里会显示对应记录。</p>
      </div>

      <div v-else class="link-table-shell">
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
    </section>
  </section>
</template>

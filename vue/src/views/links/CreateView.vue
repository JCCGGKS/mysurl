<script setup>
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { postJson } from '../../services/api'
import { handleUnauthorized } from '../../router'

const router = useRouter()

const longUrl = ref('')
const loading = ref(false)
const errorMessage = ref('')
const copied = ref(false)
const result = ref(null)

const examples = [
  'https://github.com/JCCGGKS/mysurl',
  'https://go-zero.dev/zh-cn/',
  'https://github.com/golang-jwt/jwt',
]

const statusLabel = computed(() => {
  if (loading.value) return '生成中'
  if (errorMessage.value) return '请求失败'
  if (result.value) return '已生成'
  return '准备就绪'
})

function validateUrl(value) {
  if (!value.trim()) return '请输入长链接'
  if (!/^https?:\/\//i.test(value)) return '链接必须以 http:// 或 https:// 开头'

  try {
    const parsed = new URL(value)
    if (!['http:', 'https:'].includes(parsed.protocol)) {
      return '当前仅支持 http 和 https 链接'
    }
  } catch {
    return '请输入合法的 URL'
  }

  return ''
}

async function submit() {
  copied.value = false
  errorMessage.value = validateUrl(longUrl.value)
  if (errorMessage.value) return

  loading.value = true
  errorMessage.value = ''

  try {
    result.value = await postJson(
      '/api/v1/links',
      {
        long_url: longUrl.value.trim(),
      },
      { auth: true },
    )
  } catch (error) {
    if (error.status === 401) {
      handleUnauthorized(router)
      return
    }
    errorMessage.value = error.message || '请求失败，请稍后重试'
    result.value = null
  } finally {
    loading.value = false
  }
}

async function copyShortUrl() {
  if (!result.value?.short_url) return
  errorMessage.value = ''

  try {
    await writeToClipboard(result.value.short_url)
    copied.value = true
  } catch {
    copied.value = false
    errorMessage.value = '复制失败，请手动复制短链'
  }
}

function fillExample(example) {
  longUrl.value = example
  errorMessage.value = ''
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
      <div class="workspace-head">
        <div>
          <h2>输入长链并生成短链</h2>
        </div>
        <p class="workspace-note">状态：{{ statusLabel }}</p>
      </div>
    </header>

    <section class="workspace-card workspace-card-dashboard">
      <div class="page-intro">
        <div class="inline-signal">
          <strong>创建短链</strong>
        </div>
      </div>

      <form class="composer" @submit.prevent="submit">
        <label class="field-label" for="long-url">原始长链接</label>
        <textarea
          id="long-url"
          v-model="longUrl"
          class="url-input"
          rows="5"
          placeholder="https://example.com/article/123?from=campaign"
          :disabled="loading"
        ></textarea>

        <div class="example-row">
          <span class="example-title">示例：</span>
          <button
            v-for="example in examples"
            :key="example"
            class="example-chip"
            type="button"
            @click="fillExample(example)"
          >
            {{ example.replace(/^https?:\/\//, '') }}
          </button>
        </div>

        <div class="action-row">
          <button class="submit-button" type="submit" :disabled="loading">
            {{ loading ? '生成中...' : '生成短链' }}
          </button>
        </div>
      </form>

      <p v-if="errorMessage" class="feedback feedback-error" role="alert">
        {{ errorMessage }}
      </p>

      <section v-if="result" class="result-card" aria-live="polite">
        <div class="result-head">
          <h3>短链已生成</h3>
          <span class="result-badge">Ready to share</span>
        </div>

        <dl class="result-grid">
          <div class="result-item result-item-primary">
            <dt>short_url</dt>
            <dd>{{ result.short_url }}</dd>
          </div>
          <div class="result-item">
            <dt>short_code</dt>
            <dd>{{ result.short_code }}</dd>
          </div>
          <div class="result-item">
            <dt>original_url</dt>
            <dd>{{ result.original_url }}</dd>
          </div>
        </dl>

        <div class="result-actions">
          <button class="secondary-button" type="button" @click="copyShortUrl">
            {{ copied ? '复制成功' : '复制短链' }}
          </button>
          <a class="ghost-link" :href="result.short_url" target="_blank" rel="noreferrer">
            打开短链
          </a>
        </div>
      </section>
    </section>
  </section>
</template>

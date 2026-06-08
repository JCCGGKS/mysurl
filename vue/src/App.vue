<script setup>
import { computed, ref } from 'vue'

const longUrl = ref('')
const loading = ref(false)
const errorMessage = ref('')
const copied = ref(false)
const result = ref(null)

const examples = [
  'https://github.com/JCCGGKS/mysurl',
  'https://xiaolincoding.com/other/offer.html',
  'https://go-zero.dev/docs/tasks',
]

const statusLabel = computed(() => {
  if (loading.value) return '生成中'
  if (errorMessage.value) return '请求失败'
  if (result.value) return '已生成'
  return '准备就绪'
})

function validateUrl(value) {
  if (!value.trim()) {
    return '请输入长链接'
  }

  if (!/^https?:\/\//i.test(value)) {
    return '链接必须以 http:// 或 https:// 开头'
  }

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

  if (errorMessage.value) {
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    const response = await fetch('/api/v1/links', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        long_url: longUrl.value.trim(),
      }),
    })

    const payload = await response.json().catch(() => null)

    if (!response.ok) {
      errorMessage.value =
        payload?.msg || mapErrorMessage(response.status) || '请求失败，请稍后重试'
      result.value = null
      return
    }

    if (payload?.code !== 0 || !payload?.data) {
      errorMessage.value = payload?.msg || '服务返回格式异常'
      result.value = null
      return
    }

    result.value = payload.data
  } catch {
    errorMessage.value = '网络异常或服务不可用，请确认后端已启动'
    result.value = null
  } finally {
    loading.value = false
  }
}

function mapErrorMessage(status) {
  if (status === 400) return '参数非法，请检查链接格式'
  if (status === 404) return '接口不存在，请确认服务路由配置'
  if (status >= 500) return '服务内部错误，请稍后重试'
  return ''
}

async function copyShortUrl() {
  if (!result.value?.short_url) {
    return
  }

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
  <main class="page-shell">
    <div class="ambient ambient-left"></div>
    <div class="ambient ambient-right"></div>

    <section class="hero-panel">
      <header class="hero-copy">
        <p class="eyebrow">mysurl1 / Shorten with precision</p>
        <h1>把冗长链接压缩成一枚清晰、可分享的短入口。</h1>
        <p class="summary">
          这个页面直接对接当前短链服务。输入合法长链，立即生成短链，并保留结果供复制和验证。
        </p>
      </header>

      <div class="signal-board" aria-hidden="true">
        <div class="signal-card">
          <span class="signal-title">系统状态</span>
          <strong>{{ statusLabel }}</strong>
        </div>
        <div class="signal-card">
          <span class="signal-title">接口地址</span>
          <strong>POST /api/v1/links</strong>
        </div>
        <div class="signal-card">
          <span class="signal-title">跳转机制</span>
          <strong>GET /:code → 302</strong>
        </div>
      </div>
    </section>

    <section class="workspace-card">
      <div class="workspace-head">
        <div>
          <p class="section-kicker">Create Link</p>
          <h2>输入长链并生成短链</h2>
        </div>
        <p class="workspace-note">仅支持 `http://` 或 `https://`</p>
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
            <span v-if="loading">生成中...</span>
            <span v-else>生成短链</span>
          </button>
          <p class="action-hint">重复提交同一长链时，后端将复用已有短链。</p>
        </div>
      </form>

      <p v-if="errorMessage" class="feedback feedback-error" role="alert">
        {{ errorMessage }}
      </p>

      <section v-if="result" class="result-card" aria-live="polite">
        <div class="result-head">
          <div>
            <p class="section-kicker">Result</p>
            <h3>短链已生成</h3>
          </div>
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
          <a
            class="ghost-link"
            :href="result.short_url"
            target="_blank"
            rel="noreferrer"
          >
            打开短链
          </a>
        </div>
      </section>
    </section>
  </main>
</template>

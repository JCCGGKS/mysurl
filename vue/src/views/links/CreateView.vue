<script setup>
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { postJson } from '../../services/api'
import { handleUnauthorized } from '../../router'

const router = useRouter()

const mode = ref('single')
const longUrl = ref('')
const batchLongUrls = ref('')
const loading = ref(false)
const errorMessage = ref('')
const copied = ref(false)
const copiedBatchIndex = ref(-1)
const copiedAllBatch = ref(false)
const result = ref(null)
const batchResult = ref(null)

const examples = [
  'https://github.com/JCCGGKS/mysurl',
  'https://go-zero.dev/zh-cn/',
  'https://github.com/golang-jwt/jwt',
]

const statusLabel = computed(() => {
  if (loading.value) return '生成中'
  if (errorMessage.value) return '请求失败'
  if (mode.value === 'single' && result.value) return '已生成'
  if (mode.value === 'batch' && batchResult.value) return '已生成'
  return '准备就绪'
})

const batchLines = computed(() =>
  batchLongUrls.value
    .split('\n')
    .map((item) => item.trim())
    .filter(Boolean),
)

const singleLines = computed(() =>
  longUrl.value
    .split('\n')
    .map((item) => item.trim())
    .filter(Boolean),
)

const batchSuccessItems = computed(() =>
  Array.isArray(batchResult.value?.items)
    ? batchResult.value.items.filter((item) => item.success && item.short_url)
    : [],
)

function switchMode(nextMode) {
  if (loading.value || mode.value === nextMode) return
  mode.value = nextMode
  errorMessage.value = ''
  copied.value = false
  copiedBatchIndex.value = -1
  copiedAllBatch.value = false
}

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

function validateSingleUrl(value) {
  if (singleLines.value.length > 1) {
    return '单条创建最多支持 1 条长链接，请切换到批量创建'
  }

  return validateUrl(value)
}

function validateBatchUrls(values) {
  if (values.length === 0) return '请输入至少一条长链接'
  if (values.length > 20) return '单次最多支持 20 条长链接'

  for (const value of values) {
    const err = validateUrl(value)
    if (err) return err
  }

  return ''
}

async function submit() {
  copied.value = false
  copiedBatchIndex.value = -1
  copiedAllBatch.value = false
  errorMessage.value = ''

  if (mode.value === 'single') {
    await submitSingle()
    return
  }

  await submitBatch()
}

async function submitSingle() {
  errorMessage.value = validateSingleUrl(longUrl.value)
  if (errorMessage.value) return

  loading.value = true
  batchResult.value = null
  copiedAllBatch.value = false

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

async function submitBatch() {
  const urls = batchLines.value
  errorMessage.value = validateBatchUrls(urls)
  if (errorMessage.value) return

  loading.value = true
  result.value = null
  copiedAllBatch.value = false

  try {
    batchResult.value = await postJson(
      '/api/v1/links/batch',
      {
        long_urls: urls,
      },
      { auth: true },
    )
  } catch (error) {
    if (error.status === 401) {
      handleUnauthorized(router)
      return
    }
    errorMessage.value = error.message || '批量创建失败，请稍后重试'
    batchResult.value = null
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

async function copyBatchShortUrl(item, index) {
  if (!item?.short_url) return
  errorMessage.value = ''

  try {
    await writeToClipboard(item.short_url)
    copiedBatchIndex.value = index
    window.setTimeout(() => {
      if (copiedBatchIndex.value === index) {
        copiedBatchIndex.value = -1
      }
    }, 1600)
  } catch {
    copiedBatchIndex.value = -1
    errorMessage.value = '复制失败，请手动复制短链'
  }
}

async function copyAllBatchShortUrls() {
  if (batchSuccessItems.value.length === 0) return
  errorMessage.value = ''

  try {
    await writeToClipboard(batchSuccessItems.value.map((item) => item.short_url).join('\n'))
    copiedAllBatch.value = true
    window.setTimeout(() => {
      copiedAllBatch.value = false
    }, 1800)
  } catch {
    copiedAllBatch.value = false
    errorMessage.value = '批量复制失败，请手动复制短链'
  }
}

function fillExample(example) {
  if (mode.value === 'single') {
    longUrl.value = example
  } else {
    batchLongUrls.value = batchLongUrls.value.trim()
      ? `${batchLongUrls.value.trim()}\n${example}`
      : example
  }
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

      <div class="mode-switch" role="tablist" aria-label="创建模式">
        <button
          class="mode-switch-button"
          :class="{ 'is-active': mode === 'single' }"
          type="button"
          @click="switchMode('single')"
        >
          单条创建
        </button>
        <button
          class="mode-switch-button"
          :class="{ 'is-active': mode === 'batch' }"
          type="button"
          @click="switchMode('batch')"
        >
          批量创建
        </button>
      </div>

      <form class="composer" @submit.prevent="submit">
        <template v-if="mode === 'single'">
          <label class="field-label" for="long-url">原始长链接</label>
          <textarea
            id="long-url"
            v-model="longUrl"
            class="url-input"
            rows="4"
            placeholder="https://example.com/article/123?from=campaign"
            :disabled="loading"
          ></textarea>
          <p class="batch-help-text">当前 {{ singleLines.length }} 条，单条创建最多 1 条。</p>
        </template>

        <template v-else>
          <label class="field-label" for="batch-long-urls">原始长链接</label>
          <textarea
            id="batch-long-urls"
            v-model="batchLongUrls"
            class="url-input batch-url-input"
            rows="8"
            placeholder="每行输入一条长链接，单次最多 20 条"
            :disabled="loading"
          ></textarea>
          <p class="batch-help-text">当前 {{ batchLines.length }} 条，按换行分隔。</p>
        </template>

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
            {{
              loading
                ? '生成中...'
                : mode === 'single'
                  ? '生成短链'
                  : '批量生成短链'
            }}
          </button>
        </div>
      </form>

      <p v-if="errorMessage" class="feedback feedback-error" role="alert">
        {{ errorMessage }}
      </p>

      <section v-if="mode === 'single' && result" class="result-card" aria-live="polite">
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

      <section v-if="mode === 'batch' && batchResult" class="result-card" aria-live="polite">
        <div class="result-head">
          <h3>批量创建结果</h3>
          <span class="result-badge">
            成功 {{ batchResult.success_count }} / {{ batchResult.total }}
          </span>
        </div>

        <div class="batch-result-toolbar">
          <p class="batch-result-note">
            已生成 {{ batchResult.success_count }} 条可用短链，失败 {{ batchResult.failed_count }} 条。
          </p>
          <button
            class="secondary-button"
            type="button"
            @click="copyAllBatchShortUrls"
            :disabled="batchSuccessItems.length === 0"
          >
            {{ copiedAllBatch ? '已全部复制' : '复制全部成功短链' }}
          </button>
        </div>

        <div class="batch-summary-grid">
          <div class="result-item">
            <dt>total</dt>
            <dd>{{ batchResult.total }}</dd>
          </div>
          <div class="result-item">
            <dt>success_count</dt>
            <dd>{{ batchResult.success_count }}</dd>
          </div>
          <div class="result-item">
            <dt>failed_count</dt>
            <dd>{{ batchResult.failed_count }}</dd>
          </div>
        </div>

        <div class="link-table-shell batch-result-table-shell">
          <table class="link-table batch-result-table">
            <thead>
              <tr>
                <th>#</th>
                <th>原始链接</th>
                <th>结果</th>
                <th>短链</th>
                <th>短码</th>
                <th>备注</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in batchResult.items" :key="`${item.index}-${item.long_url}`">
                <td>{{ item.index + 1 }}</td>
                <td>
                  <span class="table-original" :title="item.long_url">{{ item.long_url }}</span>
                </td>
                <td>
                  <span class="batch-status" :class="{ 'is-success': item.success, 'is-failed': !item.success }">
                    {{ item.success ? '成功' : '失败' }}
                  </span>
                </td>
                <td>
                  <a
                    v-if="item.short_url"
                    class="table-link"
                    :href="item.short_url"
                    target="_blank"
                    rel="noreferrer"
                  >
                    {{ item.short_url }}
                  </a>
                  <span v-else>--</span>
                </td>
                <td>
                  <strong v-if="item.short_code" class="table-code">{{ item.short_code }}</strong>
                  <span v-else>--</span>
                </td>
                <td>
                  <span class="batch-error-text">{{ item.error || '--' }}</span>
                </td>
                <td>
                  <div v-if="item.short_url" class="table-actions">
                    <button class="secondary-button table-button" type="button" @click="copyBatchShortUrl(item, item.index)">
                      {{ copiedBatchIndex === item.index ? '已复制' : '复制' }}
                    </button>
                  </div>
                  <span v-else>--</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </section>
  </section>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { postJson } from '../services/api'
import { setAuth } from '../services/auth'

const router = useRouter()
const route = useRoute()

const username = ref('')
const password = ref('')
const loading = ref(false)
const errorMessage = ref('')

const successMessage = computed(() =>
  route.query.registered === '1' ? '注册成功，请登录后继续创建短链。' : '',
)

async function submit() {
  errorMessage.value = ''

  if (!username.value.trim()) {
    errorMessage.value = '请输入用户名'
    return
  }
  if (!password.value) {
    errorMessage.value = '请输入密码'
    return
  }

  loading.value = true
  try {
    const data = await postJson('/api/v1/auth/login', {
      username: username.value.trim(),
      password: password.value,
    })

    setAuth({
      token: data.token,
      user: data.user,
    })

    router.push({ name: 'create' })
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <main class="auth-page">
    <div class="auth-panel auth-panel-brand">
      <p class="eyebrow">mysurl1 / Access</p>
      <h1>登录后开始创建属于你的短链入口。</h1>
      <p class="summary">
        当前创建接口要求携带认证 token。登录成功后，前端会自动附带 Bearer token 发起创建请求。
      </p>
    </div>

    <section class="auth-panel auth-panel-form">
      <div class="auth-head">
        <p class="section-kicker">Login</p>
        <h2>账号登录</h2>
        <p class="auth-note">没有账号？<RouterLink class="auth-link" to="/register">去注册</RouterLink></p>
      </div>

      <p v-if="successMessage" class="feedback feedback-success">{{ successMessage }}</p>
      <p v-if="errorMessage" class="feedback feedback-error">{{ errorMessage }}</p>

      <form class="auth-form" @submit.prevent="submit">
        <label class="field-label" for="login-username">用户名</label>
        <input
          id="login-username"
          v-model="username"
          class="text-input"
          type="text"
          autocomplete="username"
          placeholder="demo_user"
          :disabled="loading"
        />

        <label class="field-label" for="login-password">密码</label>
        <input
          id="login-password"
          v-model="password"
          class="text-input"
          type="password"
          autocomplete="current-password"
          placeholder="请输入密码"
          :disabled="loading"
        />

        <button class="submit-button auth-submit" type="submit" :disabled="loading">
          {{ loading ? '登录中...' : '登录' }}
        </button>
      </form>
    </section>
  </main>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { postJson } from '../../services/api'

const router = useRouter()

const username = ref('')
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const errorMessage = ref('')

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
  if (password.value.length < 8) {
    errorMessage.value = '密码长度至少为 8 位'
    return
  }
  if (password.value !== confirmPassword.value) {
    errorMessage.value = '两次输入的密码不一致'
    return
  }

  loading.value = true
  try {
    await postJson('/api/v1/auth/register', {
      username: username.value.trim(),
      password: password.value,
      confirm_password: confirmPassword.value,
    })

    router.push({ name: 'login', query: { registered: '1' } })
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <main class="auth-page">
    <section class="auth-panel auth-panel-form auth-panel-single">
      <div class="auth-orbit auth-orbit-one"></div>
      <div class="auth-orbit auth-orbit-two"></div>

      <div class="auth-head">
        <p class="auth-chip">mysurl1 / Join</p>
        <h2>创建账号</h2>
        <p class="auth-note">已有账号？<RouterLink class="auth-link" to="/login">去登录</RouterLink></p>
      </div>

      <p v-if="errorMessage" class="feedback feedback-error">{{ errorMessage }}</p>

      <form class="auth-form" @submit.prevent="submit">
        <label class="field-label" for="register-username">用户名</label>
        <input
          id="register-username"
          v-model="username"
          class="text-input"
          type="text"
          autocomplete="username"
          placeholder="demo_user"
          :disabled="loading"
        />

        <label class="field-label" for="register-password">密码</label>
        <input
          id="register-password"
          v-model="password"
          class="text-input"
          type="password"
          autocomplete="new-password"
          placeholder="至少 8 位"
          :disabled="loading"
        />

        <label class="field-label" for="register-confirm">确认密码</label>
        <input
          id="register-confirm"
          v-model="confirmPassword"
          class="text-input"
          type="password"
          autocomplete="new-password"
          placeholder="再次输入密码"
          :disabled="loading"
        />

        <button class="submit-button auth-submit" type="submit" :disabled="loading">
          {{ loading ? '注册中...' : '注册' }}
        </button>
      </form>
    </section>
  </main>
</template>

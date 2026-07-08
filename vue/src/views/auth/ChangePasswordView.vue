<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { postJson } from '../../services/api'
import { getAccessToken, getUser } from '../../services/auth'

const router = useRouter()
const isLoggedIn = computed(() => !!getAccessToken())

const username = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

async function submit() {
  errorMessage.value = ''
  successMessage.value = ''

  if (!isLoggedIn.value && !username.value.trim()) {
    errorMessage.value = '请输入用户名'
    return
  }
  if (!newPassword.value) {
    errorMessage.value = '请输入新密码'
    return
  }
  if (newPassword.value.length < 8) {
    errorMessage.value = '新密码长度至少8位'
    return
  }
  if (newPassword.value !== confirmPassword.value) {
    errorMessage.value = '两次输入的新密码不一致'
    return
  }

  loading.value = true
  try {
    const payload = {
      new_password: newPassword.value,
      confirm_password: confirmPassword.value,
    }
    if (!isLoggedIn.value) {
      payload.username = username.value.trim()
    }

    await postJson('/api/v1/auth/change-password', payload)

    successMessage.value = '密码修改成功'
    username.value = ''
    newPassword.value = ''
    confirmPassword.value = ''

    setTimeout(() => {
      router.push({ name: 'login' })
    }, 1500)
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
        <p class="auth-chip">mysurl1 / Security</p>
        <h2>修改密码</h2>
        <p class="auth-note">修改后需重新登录</p>
      </div>

      <p v-if="successMessage" class="feedback feedback-success">{{ successMessage }}</p>
      <p v-if="errorMessage" class="feedback feedback-error">{{ errorMessage }}</p>

      <form class="auth-form" @submit.prevent="submit">
        <template v-if="!isLoggedIn">
          <label class="field-label" for="username">用户名</label>
          <input
            id="username"
            v-model="username"
            class="text-input"
            type="text"
            autocomplete="username"
            placeholder="请输入用户名"
            :disabled="loading"
          />
        </template>

        <label class="field-label" for="new-password">新密码</label>
        <input
          id="new-password"
          v-model="newPassword"
          class="text-input"
          type="password"
          autocomplete="new-password"
          placeholder="至少8位字符"
          :disabled="loading"
        />

        <label class="field-label" for="confirm-password">确认新密码</label>
        <input
          id="confirm-password"
          v-model="confirmPassword"
          class="text-input"
          type="password"
          autocomplete="new-password"
          placeholder="请再次输入新密码"
          :disabled="loading"
        />

        <button class="submit-button auth-submit" type="submit" :disabled="loading">
          {{ loading ? '提交中...' : '修改密码' }}
        </button>
      </form>

      <p class="auth-note" style="margin-top: 1rem;">
        <RouterLink class="auth-link" to="/login">返回登录</RouterLink>
      </p>
    </section>
  </main>
</template>

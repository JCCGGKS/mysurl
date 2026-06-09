<script setup>
import { computed } from 'vue'
import { getUser } from '../../services/auth'

const currentUser = computed(() => getUser())
</script>

<template>
  <section class="dashboard-page">
    <header class="dashboard-page-head">
      <div>
        <p class="section-kicker">User</p>
        <h2>用户管理</h2>
      </div>
      <p class="workspace-note">基础信息</p>
    </header>

    <section class="user-grid">
      <article class="user-card">
        <span class="signal-title">Username</span>
        <strong>{{ currentUser?.username || 'unknown' }}</strong>
        <p>当前登录用户名称，来源于本地保存的认证信息。</p>
      </article>
      <article class="user-card">
        <span class="signal-title">User ID</span>
        <strong>{{ currentUser?.user_id || '--' }}</strong>
        <p>登录接口返回的用户标识，可用于和后端归属逻辑对应。</p>
      </article>
      <article class="user-card">
        <span class="signal-title">Status</span>
        <strong>Authenticated</strong>
        <p>当前页面仅登录态可访问，退出后会回到登录页面。</p>
      </article>
    </section>
  </section>
</template>

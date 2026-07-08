<script setup>
import { computed, ref } from 'vue'
import { RouterLink, RouterView, useRoute, useRouter } from 'vue-router'
import { postJson } from '../../services/api'
import { clearAuth, getRefreshToken, getUser } from '../../services/auth'

const route = useRoute()
const router = useRouter()
const currentUser = computed(() => getUser())

const directItems = [
  { name: 'home', label: '首页' },
]

const navGroups = [
  {
    key: 'links',
    label: '短链管理',
    children: [
      { name: 'links-create', label: '创建短链' },
      { name: 'links-list', label: '短链列表' },
    ],
  },
  {
    key: 'users',
    label: '用户管理',
    children: [
      { name: 'users-operation-log', label: '用户操作日志' },
    ],
  },
]

const expandedGroups = ref({
  links: route.name?.toString().startsWith('links-') ?? false,
  users: route.name?.toString().startsWith('users-') ?? false,
})

function isGroupActive(group) {
  return group.children.some((item) => item.name === route.name)
}

function toggleGroup(groupKey) {
  expandedGroups.value[groupKey] = !expandedGroups.value[groupKey]
}

async function logout() {
  try {
    const refreshToken = getRefreshToken()
    if (refreshToken) {
      await postJson('/api/v1/auth/logout', { refresh_token: refreshToken })
    }
  } catch {
    // ignore logout errors
  } finally {
    clearAuth()
    router.push({ name: 'login' })
  }
}
</script>

<template>
  <main class="dashboard-shell">
    <aside class="dashboard-sidebar">
      <div class="sidebar-brand">
        <p class="sidebar-chip">mysurl1</p>
        <h1>Short Link Console</h1>
      </div>

      <nav class="sidebar-nav" aria-label="主导航">
        <RouterLink
          v-for="item in directItems"
          :key="item.name"
          :to="{ name: item.name }"
          class="sidebar-link"
          :class="{ 'is-active': route.name === item.name }"
        >
          <span class="sidebar-link-label">{{ item.label }}</span>
        </RouterLink>

        <section
          v-for="group in navGroups"
          :key="group.key"
          class="sidebar-group"
          :class="{ 'is-open': expandedGroups[group.key], 'is-active': isGroupActive(group) }"
        >
          <button class="sidebar-group-trigger" type="button" @click="toggleGroup(group.key)">
            <span class="sidebar-group-copy">
              <span class="sidebar-link-label">{{ group.label }}</span>
            </span>
            <span class="sidebar-group-arrow" aria-hidden="true">▾</span>
          </button>

          <div v-if="expandedGroups[group.key]" class="sidebar-subnav">
            <RouterLink
              v-for="item in group.children"
              :key="item.name"
              :to="{ name: item.name }"
              class="sidebar-sublink"
              :class="{ 'is-active': route.name === item.name }"
            >
              <span class="sidebar-sublink-label">{{ item.label }}</span>
            </RouterLink>
          </div>
        </section>
      </nav>

      <div class="sidebar-user">
        <p class="sidebar-user-label">当前账号</p>
        <strong>{{ currentUser?.username || 'unknown' }}</strong>
        <RouterLink class="ghost-link sidebar-link" to="/change-password">修改密码</RouterLink>
        <button class="ghost-link sidebar-logout" type="button" @click="logout">退出登录</button>
      </div>
    </aside>

    <section class="dashboard-main">
      <RouterView />
    </section>
  </main>
</template>

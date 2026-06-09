<script setup>
import { computed, ref } from 'vue'
import { RouterLink, RouterView, useRoute, useRouter } from 'vue-router'
import { clearAuth, getUser } from '../../services/auth'

const route = useRoute()
const router = useRouter()
const currentUser = computed(() => getUser())

const directItems = [
  { name: 'home', label: '首页', caption: '项目介绍' },
]

const navGroups = [
  {
    key: 'links',
    label: '短链管理',
    caption: '创建、列表与维护',
    children: [
      { name: 'links-create', label: '创建短链', caption: '新增短链' },
      { name: 'links-list', label: '短链列表', caption: '查看记录' },
    ],
  },
  {
    key: 'users',
    label: '用户管理',
    caption: '资料与安全设置',
    children: [
      { name: 'users-profile', label: '用户资料', caption: '账号信息' },
      { name: 'users-security', label: '安全设置', caption: '密码与会话' },
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

function logout() {
  clearAuth()
  router.push({ name: 'login' })
}
</script>

<template>
  <main class="dashboard-shell">
    <aside class="dashboard-sidebar">
      <div class="sidebar-brand">
        <p class="sidebar-chip">mysurl1</p>
        <h1>Short Link Console</h1>
        <p class="sidebar-summary">登录后的所有操作都从这里进入，短链创建会自动关联当前账号。</p>
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
          <span class="sidebar-link-caption">{{ item.caption }}</span>
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
              <span class="sidebar-link-caption">{{ group.caption }}</span>
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
              <span class="sidebar-sublink-caption">{{ item.caption }}</span>
            </RouterLink>
          </div>
        </section>
      </nav>

      <div class="sidebar-user">
        <p class="sidebar-user-label">当前账号</p>
        <strong>{{ currentUser?.username || 'unknown' }}</strong>
        <button class="ghost-link sidebar-logout" type="button" @click="logout">退出登录</button>
      </div>
    </aside>

    <section class="dashboard-main">
      <RouterView />
    </section>
  </main>
</template>

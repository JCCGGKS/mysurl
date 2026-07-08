import { createRouter, createWebHistory } from 'vue-router'
import { clearAuth, getAccessToken } from './services/auth'
import LoginView from './views/auth/LoginView.vue'
import RegisterView from './views/auth/RegisterView.vue'
import ChangePasswordView from './views/auth/ChangePasswordView.vue'
import HomeView from './views/home/HomeView.vue'
import DashboardLayout from './views/layout/DashboardLayout.vue'
import CreateView from './views/links/CreateView.vue'
import LinkListView from './views/links/LinkListView.vue'
import UserOperationLogView from './views/users/UserOperationLogView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: () => (getAccessToken() ? '/app/home' : '/login'),
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView,
      meta: { guestOnly: true },
    },
    {
      path: '/register',
      name: 'register',
      component: RegisterView,
      meta: { guestOnly: true },
    },
    {
      path: '/change-password',
      name: 'change-password',
      component: ChangePasswordView,
      meta: { guestOnly: true },
    },
    {
      path: '/app',
      component: DashboardLayout,
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          redirect: { name: 'home' },
        },
        {
          path: 'home',
          name: 'home',
          component: HomeView,
        },
        {
          path: 'links',
          redirect: { name: 'links-create' },
        },
        {
          path: 'links/create',
          name: 'links-create',
          component: CreateView,
        },
        {
          path: 'links/list',
          name: 'links-list',
          component: LinkListView,
        },
        {
          path: 'users',
          redirect: { name: 'users-operation-log' },
        },
        {
          path: 'users/operation-log',
          name: 'users-operation-log',
          component: UserOperationLogView,
        },
        {
          path: ':pathMatch(.*)*',
          redirect: { name: 'home' },
        },
      ],
    },
  ],
})

router.beforeEach((to) => {
  const token = getAccessToken()

  if (to.meta.requiresAuth && !token) {
    return { name: 'login' }
  }

  if (to.meta.guestOnly && token) {
    return { name: 'home' }
  }

  return true
})

export function handleUnauthorized(routerInstance = router) {
  clearAuth()
  routerInstance.push({ name: 'login' })
}

export default router

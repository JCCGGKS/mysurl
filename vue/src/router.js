import { createRouter, createWebHistory } from 'vue-router'
import { clearAuth, getToken } from './services/auth'
import CreateView from './views/CreateView.vue'
import LoginView from './views/LoginView.vue'
import RegisterView from './views/RegisterView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: () => (getToken() ? '/create' : '/login'),
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
      path: '/create',
      name: 'create',
      component: CreateView,
      meta: { requiresAuth: true },
    },
  ],
})

router.beforeEach((to) => {
  const token = getToken()

  if (to.meta.requiresAuth && !token) {
    return { name: 'login' }
  }

  if (to.meta.guestOnly && token) {
    return { name: 'create' }
  }

  return true
})

export function handleUnauthorized(routerInstance = router) {
  clearAuth()
  routerInstance.push({ name: 'login' })
}

export default router

import { createRouter, createWebHistory } from 'vue-router'
import { clearAuth, getToken } from './services/auth'
import LoginView from './views/auth/LoginView.vue'
import RegisterView from './views/auth/RegisterView.vue'
import HomeView from './views/home/HomeView.vue'
import DashboardLayout from './views/layout/DashboardLayout.vue'
import CreateView from './views/links/CreateView.vue'
import LinkListView from './views/links/LinkListView.vue'
import UserView from './views/users/UserView.vue'
import UserSecurityView from './views/users/UserSecurityView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: () => (getToken() ? '/app/home' : '/login'),
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
          redirect: { name: 'users-profile' },
        },
        {
          path: 'users/profile',
          name: 'users-profile',
          component: UserView,
        },
        {
          path: 'users/security',
          name: 'users-security',
          component: UserSecurityView,
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
  const token = getToken()

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

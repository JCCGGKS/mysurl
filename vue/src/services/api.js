import router from '../router'
import { clearAuth, getAccessToken, getRefreshToken, getUser, setAuth } from './auth'

let refreshPromise = null

export async function getJson(url, options = {}) {
  return requestJson(url, {
    method: 'GET',
    ...options,
  })
}

export async function postJson(url, body, options = {}) {
  return requestJson(url, {
    method: 'POST',
    body,
    ...options,
  })
}

async function requestJson(url, { method = 'GET', body, auth = false, retry = true } = {}) {
  const headers = {}
  if (body !== undefined) {
    headers['Content-Type'] = 'application/json'
  }

  if (auth) {
    const token = getAccessToken()
    if (token) {
      headers.Authorization = `Bearer ${token}`
    }
  }

  const response = await fetch(url, {
    method,
    headers,
    body: body === undefined ? undefined : JSON.stringify(body),
  })

  if (auth && response.status === 401 && retry) {
    const refreshed = await tryRefresh()
    if (refreshed) {
      return requestJson(url, { method, body, auth, retry: false })
    }
  }

  return handleJsonResponse(response)
}

async function tryRefresh() {
  const refreshToken = getRefreshToken()
  if (!refreshToken) {
    clearAndRedirectToLogin()
    return false
  }

  if (!refreshPromise) {
    refreshPromise = refreshAuth(refreshToken).finally(() => {
      refreshPromise = null
    })
  }

  return refreshPromise
}

async function refreshAuth(refreshToken) {
  const response = await fetch('/api/v1/auth/refresh', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      refresh_token: refreshToken,
    }),
  })

  const payload = await response.json().catch(() => null)
  if (!response.ok || payload?.code !== 0 || !payload?.data) {
    clearAndRedirectToLogin()
    return false
  }

  const data = payload.data
  setAuth({
    accessToken: data.access_token,
    refreshToken: data.refresh_token,
    user: data.user || getUser(),
  })

  return true
}

function clearAndRedirectToLogin() {
  clearAuth()
  if (router.currentRoute.value.name !== 'login') {
    router.push({ name: 'login' })
  }
}

async function handleJsonResponse(response) {
  const payload = await response.json().catch(() => null)

  if (!response.ok || payload?.code !== 0 || !payload?.data) {
    const error = new Error(
      payload?.msg || mapErrorMessage(response.status) || '请求失败，请稍后重试',
    )
    error.status = response.status
    error.payload = payload
    throw error
  }

  return payload.data
}

function mapErrorMessage(status) {
  if (status === 400) return '参数非法，请检查输入'
  if (status === 401) return '登录状态失效，请重新登录'
  if (status === 404) return '接口不存在，请确认服务路由配置'
  if (status >= 500) return '服务内部错误，请稍后重试'
  return ''
}

import { getToken } from './auth'

export async function postJson(url, body, { auth = false } = {}) {
  const headers = {
    'Content-Type': 'application/json',
  }

  if (auth) {
    const token = getToken()
    if (token) {
      headers.Authorization = `Bearer ${token}`
    }
  }

  const response = await fetch(url, {
    method: 'POST',
    headers,
    body: JSON.stringify(body),
  })

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

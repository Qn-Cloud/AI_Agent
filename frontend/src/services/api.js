import axios from 'axios'
import { ElMessage } from 'element-plus'
import config from '../config'

// 创建API客户端实例
const api = axios.create({
  baseURL: config.apiBaseUrl,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    // 从localStorage获取token
    const token = localStorage.getItem('auth_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    const { data } = response
    
    // 后端统一响应格式：{ code, msg, data }
    if (data.code === 200 || data.code === 0) {
      return data
    } else if (data.code === 401) {
      // token失效，清除本地存储并跳转登录
      localStorage.removeItem('auth_token')
      localStorage.removeItem('user_info')
      window.location.href = '/login'
      return Promise.reject(new Error('登录已过期，请重新登录'))
    } else {
      // 其他错误
      ElMessage.error(data.msg || '请求失败')
      return Promise.reject(new Error(data.msg || '请求失败'))
    }
  },
  (error) => {
    if (error.response) {
      const { status, data } = error.response
      
      if (status === 401) {
        localStorage.removeItem('auth_token')
        localStorage.removeItem('user_info')
        window.location.href = '/login'
        ElMessage.error('登录已过期，请重新登录')
      } else if (status === 403) {
        ElMessage.error('没有权限访问该资源')
      } else if (status === 404) {
        ElMessage.error('请求的资源不存在')
      } else if (status === 500) {
        ElMessage.error('服务器内部错误')
      } else {
        ElMessage.error(data?.msg || error.message || '网络错误')
      }
    } else if (error.code === 'ECONNABORTED') {
      ElMessage.error('请求超时，请重试')
    } else {
      ElMessage.error('网络错误，请检查网络连接')
    }
    
    return Promise.reject(error)
  }
)

export default api

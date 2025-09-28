import axios from 'axios'
import { ElMessage } from 'element-plus'
import config from '../config'

// 创建API实例的通用函数
const createApiInstance = (baseURL, serviceName = '') => {
  console.log(`🔧 创建API实例 [${serviceName}]:`, baseURL)
  
  const api = axios.create({
    baseURL,
    timeout: 30000,
    headers: {
      'Content-Type': 'application/json'
    }
  })

  // 请求拦截器
  api.interceptors.request.use(
    (config) => {
      console.log(`🚀 [${serviceName}] 发起请求:`, config.method?.toUpperCase(), config.url, config.baseURL)
      
      // 从localStorage获取token
      const token = localStorage.getItem('auth_token')
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }
      
      // 添加服务标识头部（可选）
      if (serviceName) {
        // 避免中文字符导致的编码问题，使用英文
        const serviceNameMap = {
          '聊天服务': 'chat-service',
          '角色服务': 'character-service',
          '用户服务': 'user-service',
          'AI服务': 'ai-service',
          '语音服务': 'speech-service',
          '存储服务': 'storage-service',
          '网关服务': 'gateway-service',
          '默认服务': 'default-service'
        }
        config.headers['X-Service-Name'] = serviceNameMap[serviceName] || serviceName
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
      console.log(`✅ [${serviceName}] 请求成功:`, response.status, response.config.url)
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
        ElMessage.error(data.msg || `${serviceName}服务请求失败`)
        return Promise.reject(new Error(data.msg || `${serviceName}服务请求失败`))
      }
    },
    (error) => {
      console.error(`❌ [${serviceName}] 请求失败:`, error.message, error.config?.url)
      
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
          ElMessage.error(`${serviceName}服务器内部错误`)
        } else {
          ElMessage.error(data?.msg || error.message || `${serviceName}网络错误`)
        }
      } else if (error.code === 'ECONNABORTED') {
        ElMessage.error(`${serviceName}请求超时，请重试`)
      } else {
        ElMessage.error(`${serviceName}网络错误，请检查网络连接`)
      }
      
      return Promise.reject(error)
    }
  )

  return api
}

// 打印配置信息用于调试
console.log('🔧 API配置信息:', {
  chatBaseUrl: config.api.chatBaseUrl,
  characterBaseUrl: config.api.characterBaseUrl,
  环境变量: {
    VITE_CHAT_API_URL: import.meta.env.VITE_CHAT_API_URL,
    VITE_CHARACTER_API_URL: import.meta.env.VITE_CHARACTER_API_URL
  }
})

// 创建各个微服务的API实例
export const chatApi = createApiInstance(config.api.chatBaseUrl, '聊天服务')
export const characterApi = createApiInstance(config.api.characterBaseUrl, '角色服务')
export const userApi = createApiInstance(config.api.userBaseUrl, '用户服务')
export const aiApi = createApiInstance(config.api.aiBaseUrl, 'AI服务')
export const speechApi = createApiInstance(config.api.speechBaseUrl, '语音服务')
export const storageApi = createApiInstance(config.api.storageBaseUrl, '存储服务')
export const gatewayApi = createApiInstance(config.api.gatewayBaseUrl, '网关服务')

// 保留原有的默认API实例以兼容现有代码
export const defaultApi = createApiInstance(config.apiBaseUrl, '默认服务')

export default defaultApi 
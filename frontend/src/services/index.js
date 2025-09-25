// API服务统一导出
export { default as api } from './api'
export { userApi } from './user'
export { characterApi } from './character'
export { chatApi } from './chat'
export { speechApi, VoiceRecorder, VoicePlayer } from './speech'
export { aiApi, ChatSession } from './ai'

// 健康检查和系统状态
export const systemApi = {
  // 健康检查
  healthCheck() {
    return api.get('/api/health')
  },

  // 服务状态
  getServiceStatus() {
    return api.get('/api/status')
  },

  // 版本信息
  getVersionInfo() {
    return api.get('/api/version')
  }
}

// 存储服务API
export const storageApi = {
  // 上传文件
  uploadFile(file, options = {}) {
    const formData = new FormData()
    formData.append('file', file)
    
    if (options.path) {
      formData.append('path', options.path)
    }
    if (options.public !== undefined) {
      formData.append('public', options.public)
    }
    
    return api.post('/api/storage/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      },
      onUploadProgress: options.onProgress
    })
  },

  // 删除文件
  deleteFile(fileId) {
    return api.delete(`/api/storage/file/${fileId}`)
  },

  // 获取文件信息
  getFileInfo(fileId) {
    return api.get(`/api/storage/file/${fileId}`)
  }
}

// 工具函数
export const utils = {
  // 格式化文件大小
  formatFileSize(bytes) {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  },

  // 格式化时间
  formatTime(date) {
    if (!date) return ''
    const d = new Date(date)
    return d.toLocaleString('zh-CN')
  },

  // 检查文件类型
  isImageFile(file) {
    return file.type.startsWith('image/')
  },

  isAudioFile(file) {
    return file.type.startsWith('audio/')
  },

  isVideoFile(file) {
    return file.type.startsWith('video/')
  },

  // 生成唯一ID
  generateId() {
    return Date.now().toString(36) + Math.random().toString(36).substr(2)
  },

  // 防抖函数
  debounce(func, wait) {
    let timeout
    return function executedFunction(...args) {
      const later = () => {
        clearTimeout(timeout)
        func(...args)
      }
      clearTimeout(timeout)
      timeout = setTimeout(later, wait)
    }
  },

  // 节流函数
  throttle(func, limit) {
    let inThrottle
    return function() {
      const args = arguments
      const context = this
      if (!inThrottle) {
        func.apply(context, args)
        inThrottle = true
        setTimeout(() => inThrottle = false, limit)
      }
    }
  }
}

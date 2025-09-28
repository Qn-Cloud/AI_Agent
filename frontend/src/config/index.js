// 应用配置
export const config = {
  // API配置 - 支持多个微服务
  api: {
    // 开发环境使用代理，生产环境使用直接地址
    chatBaseUrl: import.meta.env.DEV ? '' : (import.meta.env.VITE_CHAT_API_URL || 'http://192.168.23.188:7001'),
    characterBaseUrl: import.meta.env.DEV ? '' : (import.meta.env.VITE_CHARACTER_API_URL || 'http://192.168.23.188:7002'),
    userBaseUrl: import.meta.env.DEV ? '' : (import.meta.env.VITE_USER_API_URL || 'http://192.168.23.188:7003'),
    aiBaseUrl: import.meta.env.DEV ? '' : (import.meta.env.VITE_AI_API_URL || 'http://192.168.23.188:7004'),
    speechBaseUrl: import.meta.env.DEV ? '' : (import.meta.env.VITE_SPEECH_API_URL || 'http://192.168.23.188:7005'),
    storageBaseUrl: import.meta.env.DEV ? '' : (import.meta.env.VITE_STORAGE_API_URL || 'http://192.168.23.188:7006'),
    gatewayBaseUrl: import.meta.env.DEV ? '' : (import.meta.env.VITE_GATEWAY_API_URL || 'http://192.168.23.188:8888')
  },
  
  // 保留旧的配置以兼容现有代码
  apiBaseUrl: import.meta.env.VITE_API_BASE_URL || 'http://192.168.23.188:8888',
  
  // 应用信息
  appTitle: import.meta.env.VITE_APP_TITLE || 'AI角色扮演语音交互产品',
  appVersion: import.meta.env.VITE_APP_VERSION || '1.0.0',
  
  // 语音配置
  speech: {
    // 支持的语音格式
    supportedFormats: ['mp3', 'wav', 'ogg'],
    // 默认语音设置
    defaultVoiceSettings: {
      rate: 1.0,
      pitch: 1.0,
      volume: 0.8
    },
    // 录音配置
    recordingConfig: {
      sampleRate: 16000,
      channelCount: 1,
      echoCancellation: true,
      noiseSuppression: true
    }
  },
  
  // 聊天配置
  chat: {
    // 最大消息长度
    maxMessageLength: 2000,
    // 消息历史保留数量
    maxHistoryMessages: 100,
    // 默认每页消息数
    defaultPageSize: 50
  },
  
  // 文件上传配置
  upload: {
    // 支持的图片格式
    imageFormats: ['jpg', 'jpeg', 'png', 'gif', 'webp'],
    // 支持的音频格式
    audioFormats: ['mp3', 'wav', 'ogg', 'webm'],
    // 最大文件大小 (MB)
    maxFileSize: 10,
    // 最大图片大小 (MB)
    maxImageSize: 5
  },
  
  // 缓存配置
  cache: {
    // 角色列表缓存时间 (分钟)
    characterListTTL: 30,
    // 用户信息缓存时间 (分钟)
    userInfoTTL: 60
  },
  
  // 开发配置
  isDevelopment: import.meta.env.DEV,
  isProduction: import.meta.env.PROD
}

export default config

import { defineStore } from 'pinia'
import { chatApiService, aiApiService as aiApi, speechApiService as speechApi, VoiceRecorder, VoicePlayer } from '../services'

// 角色ID映射
const characterMap = {
  1: '哈利·波特',
  2: '苏格拉底', 
  3: '莎士比亚',
  4: '爱因斯坦',
  5: '夏洛克·福尔摩斯',
  6: '赫敏·格兰杰'
}

// 获取角色名称
const getCharacterName = (id) => {
  return characterMap[id] || `角色${id}`
}

export const useChatStore = defineStore('chat', {
  state: () => ({
    conversations: [],
    currentConversation: null,
    isLoading: false,
    error: null,
    
    // 语音相关状态
    isRecording: false,
    isPlaying: false,
    currentAudioUrl: null,
    
    // 统计数据
    stats: {
      messageCount: 0,
      characterCount: 0,
      activeDays: 0
    },
    
    // 分组历史数据
    groupedHistory: {
      todays: [],
      yesterdays: [],
      befores: []
    }
  }),

  getters: {
    // 获取当前对话的消息
    currentMessages: (state) => {
      return state.currentConversation?.messages || []
    },
    
    // 按角色分组的对话
    conversationsByCharacter: (state) => {
      const grouped = {}
      state.conversations.forEach(conv => {
        if (!grouped[conv.characterId]) {
          grouped[conv.characterId] = []
        }
        grouped[conv.characterId].push(conv)
      })
      return grouped
    },
    
    // 最近的对话
    recentConversations: (state) => {
      return state.conversations
        .slice()
        .sort((a, b) => new Date(b.lastUpdate) - new Date(a.lastUpdate))
        .slice(0, 10)
    },
    
    // 统计信息
    stats: (state) => {
      return {
        totalConversations: state.conversations.length,
        totalMessages: state.conversations.reduce((sum, conv) => sum + (conv.messageCount || conv.messages?.length || 0), 0),
        uniqueCharacters: new Set(state.conversations.map(conv => conv.characterId)).size
      }
    }
  },

  actions: {
    // 初始化语音控制器
    initVoiceControllers() {
      if (!this.voiceRecorder) {
        this.voiceRecorder = new VoiceRecorder()
        // 设置语音转文字的回调函数
        this.voiceRecorder.setTranscriptCallback((text) => {
          this.handleTranscript(text)
        })
      }
      if (!this.voicePlayer) {
        this.voicePlayer = new VoicePlayer()
      }
    },

    // 处理语音转文字结果
    handleTranscript(text) {
      console.log('🔤 语音转文字结果:', text)
      if (text && text.trim() && !text.includes('[语音转文字失败')) {
        // 设置转录文本到输入框
        this.transcript = text.trim()
        console.log('📝 语音转录完成，文本已设置到输入框')
      } else if (text.includes('[语音转文字失败')) {
        console.error('语音转文字失败')
        this.transcript = ''
      }
    },

    // 加载分组的对话历史
    async loadGroupedHistory() {
      try {
        this.isLoading = true
        
        const response = await chatApiService.getChatHistoryBefore(1)
        
        if (response && response.data) {
          this.groupedHistory = {
            todays: Array.isArray(response.data.todays) ? response.data.todays : [],
            yesterdays: Array.isArray(response.data.yesterdays) ? response.data.yesterdays : [],
            befores: Array.isArray(response.data.befores) ? response.data.befores : []
          }
        }
        
      } catch (error) {
        this.error = error.message
        console.error('❌ 加载分组对话历史失败:', error)
      } finally {
        this.isLoading = false
      }
    },

    // 加载对话历史
    async loadConversationHistory(params = {}) {
      try {
        this.isLoading = true
        this.error = null
        
        console.log('正在从后端加载对话历史...', params)
        
        const response = await chatApiService.getConversationHistory({
          page: params.page || 1,
          pageSize: params.pageSize || 20,
          characterId: params.characterId,
          userId: 1, // 明确设置为1
          startTime: params.startTime,
          endTime: params.endTime
        })
        
        console.log('🔍 API响应:', response)
        console.log('🔍 响应数据结构:', response?.data)
        
        if (response && response.data) {
          // 清空当前对话列表（如果是第一页）
          if (!params.page || params.page === 1) {
            this.conversations = []
          }
          
          // 转换后端数据格式到前端格式
          const conversations = response.data.list?.map(conv => ({
            id: conv.conversation_id,
            characterId: conv.character_id,
            title: `与${getCharacterName(conv.character_id)}的对话`, // 使用角色名称生成标题
            startTime: new Date(conv.LastMessageTime), // 使用最后消息时间作为开始时间
            lastUpdate: new Date(conv.LastMessageTime),
            messageCount: conv.message_count,
            status: 1, // 默认状态为正常
            lastMessageContent: conv.LastMessageContent, // 添加最后消息内容
            conversationDuration: conv.conversation_duration,
            messages: [] // 消息会在需要时单独加载
          })) || []
          
          // 添加到对话列表
          this.conversations.push(...conversations)
          
          // 更新统计数据
          this.stats = {
            conversationTotal: response.data.ConversationTotal || 0,
            messageCount: response.data.message_count || 0,
            characterCount: response.data.character_count || 0,
            activeDays: response.data.active_days || 0
          }
          
          console.log('成功加载对话历史:', conversations.length, '条对话')
          console.log('🔍 转换后的对话数据:', conversations)
          console.log('🔍 统计数据:', this.stats)
          return {
            conversations,
            total: response.data.ConversationTotal || 0,
            hasMore: false // 后端没有返回has_more字段，暂时设为false
          }
        }
      } catch (error) {
        this.error = error.message
        console.error('加载对话历史失败:', error)
        
        // 如果API失败，使用默认的模拟数据
        console.log('API失败，使用模拟数据')
        if (this.conversations.length === 0) {
          // Fallback to mock data if API fails and no conversations are loaded yet
          this.conversations = [
            {
              id: 'conv-1',
              characterId: 1,
              title: '与哈利·波特的魔法之旅',
              startTime: new Date('2025-09-23T10:30:00'),
              lastUpdate: new Date('2025-09-23T11:15:00'),
              messages: [
                {
                  id: 'msg-1',
                  type: 'user',
                  content: '你好，哈利！你能告诉我一些关于霍格沃茨的事情吗？',
                  timestamp: new Date('2025-09-23T10:30:00')
                },
                {
                  id: 'msg-2',
                  type: 'ai',
                  content: '你好！当然可以！霍格沃茨是世界上最棒的魔法学校之一。那里有四个学院：格兰芬多、赫奇帕奇、拉文克劳和斯莱特林。我在格兰芬多学院，我们以勇敢和正义著称。学校里有很多神奇的地方，比如会动的楼梯、会说话的肖像画，还有禁林里的神奇动物们。你想了解哪个方面呢？',
                  timestamp: new Date('2025-09-23T10:31:00')
                },
                {
                  id: 'msg-3',
                  type: 'user',
                  content: '格兰芬多学院听起来很棒！你能说说你的朋友们吗？',
                  timestamp: new Date('2025-09-23T11:10:00')
                },
                {
                  id: 'msg-4',
                  type: 'ai',
                  content: '我有两个最好的朋友，罗恩·韦斯莱和赫敏·格兰杰。罗恩来自一个纯血统的巫师家庭，他很忠诚，虽然有时候会害怕，但在关键时刻总是很勇敢。赫敏是麻瓜出身，但她是我们年级最聪明的学生，几乎知道所有的魔法知识。我们三个一起经历了很多冒险，他们就像我的家人一样重要。',
                  timestamp: new Date('2025-09-23T11:15:00')
                }
              ]
            },
            {
              id: 'conv-2',
              characterId: 2,
              title: '与苏格拉底探讨智慧',
              startTime: new Date('2025-09-22T14:20:00'),
              lastUpdate: new Date('2025-09-22T15:05:00'),
              messages: [
                {
                  id: 'msg-5',
                  type: 'user',
                  content: '苏格拉底，什么是真正的智慧？',
                  timestamp: new Date('2025-09-22T14:20:00')
                },
                {
                  id: 'msg-6',
                  type: 'ai',
                  content: '我的朋友，这是一个深刻的问题。我曾说过"我知道我什么都不知道"。真正的智慧，不是拥有大量的知识，而是认识到自己无知的局限性。那么，你认为知识和智慧有什么区别呢？',
                  timestamp: new Date('2025-09-22T14:22:00')
                },
                {
                  id: 'msg-7',
                  type: 'user',
                  content: '我觉得知识是信息，而智慧是如何运用这些信息？',
                  timestamp: new Date('2025-09-22T15:00:00')
                },
                {
                  id: 'msg-8',
                  type: 'ai',
                  content: '很好的思考！你提到了运用，这让我想到另一个问题：如果一个人拥有丰富的知识，但从未将其用于善的行为，我们还能称他为智者吗？智慧是否必须与美德相结合？',
                  timestamp: new Date('2025-09-22T15:05:00')
                }
              ]
            }
          ]
        }
        return { conversations: this.conversations, total: this.conversations.length, hasMore: false }
      } finally {
        this.isLoading = false
      }
    },

    // 初始化数据（页面加载时调用）
    async initializeData() {
      console.log('🔄 chatStore.initializeData() 被调用')
      console.log('🔄 当前对话数量:', this.conversations.length)
      
      if (this.conversations.length === 0) {
        console.log('🔄 对话列表为空，开始加载...')
        await this.loadConversationHistory()
      } else {
        console.log('🔄 对话列表已存在，跳过加载')
      }
    },

    // 手动测试API连接
    async testApiConnection() {
      console.log('🧪 开始手动测试API连接...')
      try {
        const response = await chatApiService.getConversationHistory({
          page: 1,
          pageSize: 5,
          userId: 1
        })
        console.log('🧪 API测试成功:', response)
        return response
      } catch (error) {
        console.error('🧪 API测试失败:', error)
        throw error
      }
    },

    // 开始新对话
    async startNewConversation(characterId) {
      try {
        this.isLoading = true
        this.error = null
        
        const response = await chatApiService.createConversation({
          characterId,
          title: `与${getCharacterName(characterId)}的对话`
        })
        
        if (response && response.data) {
          const newConversation = {
            id: response.data.id,
            characterId,
            title: response.data.title,
            startTime: new Date(),
            lastUpdate: new Date(),
            messages: [],
            status: 1
          }
          
          this.conversations.unshift(newConversation)
          this.currentConversation = newConversation
          
          return newConversation
        }
      } catch (error) {
        this.error = error.message
        console.error('创建对话失败:', error)
        
        // 创建临时对话
        const tempConversation = {
          id: `temp-${Date.now()}`,
          characterId,
          title: `与${getCharacterName(characterId)}的对话`,
          startTime: new Date(),
          lastUpdate: new Date(),
          messages: [],
          status: 1
        }
        
        this.conversations.unshift(tempConversation)
        this.currentConversation = tempConversation
        
        return tempConversation
      } finally {
        this.isLoading = false
      }
    },

    // 选择对话
    selectConversation(conversationId) {
      const conversation = this.conversations.find(conv => conv.id === conversationId)
      if (conversation) {
        this.currentConversation = conversation
      }
    },

    // 发送消息 (支持SSE实时更新)
    async sendMessage(content, type = 'text') {
      if (!this.currentConversation) {
        throw new Error('没有选中的对话')
      }

      try {
        this.isLoading = true
        
        const userMessage = {
          id: `msg-${Date.now()}`,
          type: 'user',
          content,
          timestamp: new Date()
        }
        
        // 添加用户消息
        this.currentConversation.messages.push(userMessage)
        
        // 创建AI回复消息占位符
        const aiMessage = {
          id: `ai-${Date.now()}`,
          type: 'ai',
          content: '',
          timestamp: new Date(),
          isStreaming: true // 标记为流式回复
        }
        
        this.currentConversation.messages.push(aiMessage)
        
        // 发送到后端 (使用Fetch流式读取，更稳定)
        const response = await this.sendMessageFetchStream({
          conversationId: this.currentConversation.id,
          characterId: this.currentConversation.characterId,
          content,
          type
        }, (partialContent) => {
          // 实时更新AI回复内容
          const messageIndex = this.currentConversation.messages.findIndex(msg => msg.id === aiMessage.id)
          if (messageIndex !== -1) {
            this.currentConversation.messages[messageIndex].content = partialContent
            console.log(`🔄 实时更新消息内容，长度: ${partialContent.length}, isStreaming: ${this.currentConversation.messages[messageIndex].isStreaming}`)
          }
        })
        
        console.log('🎯 sendMessageFetchStream 返回结果:', response)
        
        // 无论是否收到完整响应，都要停止流式状态
        if (response && response.data) {
          // 更新最终的AI回复，但保持原有的ID
          const messageIndex = this.currentConversation.messages.findIndex(msg => msg.id === aiMessage.id)
          if (messageIndex !== -1) {
            // 保留原有消息的所有属性，只更新必要的字段
            const currentMessage = this.currentConversation.messages[messageIndex]
            console.log(`🔧 更新前消息状态: isStreaming=${currentMessage.isStreaming}, 内容长度=${currentMessage.content.length}`)
            
            this.currentConversation.messages[messageIndex] = {
              ...currentMessage, // 保留所有原有属性
              content: response.data.ai_message.content || currentMessage.content, // 保留已有内容
              timestamp: response.data.ai_message.timestamp ? new Date(response.data.ai_message.timestamp) : currentMessage.timestamp,
              isStreaming: false // 明确标记为完成
            }
            
            console.log(`✅ 更新后消息状态: isStreaming=${this.currentConversation.messages[messageIndex].isStreaming}, 内容长度=${this.currentConversation.messages[messageIndex].content.length}`)
            console.log('✅ AI消息流式传输完成，最终内容长度:', this.currentConversation.messages[messageIndex].content.length)
            console.log('✅ 消息ID保持不变:', this.currentConversation.messages[messageIndex].id)
          }
          
          this.currentConversation.lastUpdate = new Date()
          return this.currentConversation.messages[messageIndex]
        } else {
          // 如果没有收到完整响应，也要停止流式状态
          const messageIndex = this.currentConversation.messages.findIndex(msg => msg.id === aiMessage.id)
          if (messageIndex !== -1) {
            console.log(`⚠️ 没有收到完整响应，强制停止流式状态`)
            console.log(`⚠️ 更新前: isStreaming=${this.currentConversation.messages[messageIndex].isStreaming}`)
            this.currentConversation.messages[messageIndex].isStreaming = false
            console.log(`⚠️ 更新后: isStreaming=${this.currentConversation.messages[messageIndex].isStreaming}`)
            console.log('⚠️ 没有收到完整响应，但停止流式状态，保留内容长度:', this.currentConversation.messages[messageIndex].content.length)
          }
        }
      } catch (error) {
        this.error = error.message
        console.error('发送消息失败:', error)
        
        // 移除失败的AI消息占位符
        const aiMessageIndex = this.currentConversation.messages.findIndex(
          msg => msg.type === 'ai' && msg.isStreaming
        )
        if (aiMessageIndex !== -1) {
          this.currentConversation.messages.splice(aiMessageIndex, 1)
        }
        
        throw error
      } finally {
        this.isLoading = false
      }
    },

    // SSE发送消息的具体实现 (带重试机制)
    async sendMessageSSE(data, onUpdate, retryCount = 0) {
      const maxRetries = 2 // 最多重试2次
      
      return new Promise((resolve, reject) => {
        // 使用前端代理，避免CORS问题
        const backendURL = chatApiService.defaults?.baseURL || ''
        
        const requestData = {
          conversation_id: data.conversationId,
          character_id: data.characterId,
          content: data.content,
          message_type: data.type === 'text' ? 1 : (data.type === 'voice' ? 2 : 1), // 1=文本, 2=语音
          user_id: 1 // 暂时固定为1
        }
        
        console.log(`📤 发送SSE聊天请求到: ${backendURL}/api/chat/send (尝试 ${retryCount + 1}/${maxRetries + 1})`)
        console.log('📤 请求参数:', requestData)
        
        // 由于EventSource只支持GET请求，我们需要将参数作为查询参数
        const queryParams = new URLSearchParams(requestData).toString()
        const url = `${backendURL}/api/chat/send?${queryParams}`
        
        console.log('🔗 完整SSE URL:', url)
        
        // 创建EventSource连接
        const eventSource = new EventSource(url)
        
        // 连接保活 - 防止浏览器关闭长连接
        let keepAliveInterval
        const startKeepAlive = () => {
          keepAliveInterval = setInterval(() => {
            if (eventSource.readyState === 1) {
              console.log('🔄 SSE连接保活检查 - 连接正常')
            } else {
              console.warn('⚠️ SSE连接保活检查 - 连接异常:', eventSource.readyState)
              clearInterval(keepAliveInterval)
            }
          }, 5000) // 每5秒检查一次
        }
        
        let aiResponse = ''
        let messageId = null
        let isComplete = false
        let hasReceivedData = false // 标记是否收到过数据
        
                  eventSource.onopen = () => {
            console.log('✅ SSE连接已建立')
            console.log('🔍 EventSource readyState:', eventSource.readyState)
            
            // 监控连接状态
            const monitorInterval = setInterval(() => {
              if (eventSource.readyState !== 1) {
                console.warn(`⚠️ SSE连接状态变化: ${eventSource.readyState} (0=CONNECTING, 1=OPEN, 2=CLOSED)`)
                clearInterval(monitorInterval)
              }
              if (isComplete) {
                clearInterval(monitorInterval)
              }
            }, 500) // 每0.5秒检查一次
          }
        
        eventSource.onmessage = (event) => {
          try {
            hasReceivedData = true
            console.log('📨 收到原始SSE数据:', event.data)
            const responseData = JSON.parse(event.data)
            console.log('📨 解析后的SSE消息:', responseData)
            
            if (responseData.type === 'message') {
              // 接收到AI回复的片段
              aiResponse += responseData.content || ''
              messageId = responseData.message_id
              
              console.log('📝 累积回复内容:', aiResponse)
              
              // 触发实时更新回调
              if (onUpdate) {
                onUpdate(aiResponse)
              }
            } else if (responseData.type === 'complete') {
              // 回复完成
              isComplete = true
              eventSource.close()
              
              console.log('✅ SSE回复完成，最终内容:', aiResponse)
              
              resolve({
                data: {
                  ai_message: {
                    id: messageId || `ai-${Date.now()}`,
                    content: aiResponse,
                    timestamp: new Date().toISOString()
                  }
                }
              })
            } else if (responseData.type === 'error') {
              // 发生错误
              console.error('❌ 服务端返回错误:', responseData.message)
              eventSource.close()
              
              // 如果是可重试的错误且还有重试次数，则重试
              if (retryCount < maxRetries && !responseData.message.includes('context canceled')) {
                console.log(`🔄 准备重试... (${retryCount + 1}/${maxRetries})`)
                setTimeout(() => {
                  this.sendMessageSSE(data, onUpdate, retryCount + 1).then(resolve).catch(reject)
                }, 1000) // 1秒后重试
              } else {
                reject(new Error(responseData.message || '聊天请求失败'))
              }
            }
          } catch (error) {
            console.error('❌ 解析SSE消息失败:', error)
            console.error('❌ 原始数据:', event.data)
          }
        }
        
        eventSource.onerror = (error) => {
          console.error('❌ SSE连接错误:', error)
          console.error('❌ EventSource readyState:', eventSource.readyState)
          
          eventSource.close()
          
          if (!isComplete) {
            // 如果没有收到任何数据且还有重试次数，则重试
            if (!hasReceivedData && retryCount < maxRetries) {
              console.log(`🔄 连接失败，准备重试... (${retryCount + 1}/${maxRetries})`)
              setTimeout(() => {
                this.sendMessageSSE(data, onUpdate, retryCount + 1).then(resolve).catch(reject)
              }, 2000) // 2秒后重试
            } else {
              reject(new Error(`SSE连接中断，readyState: ${eventSource.readyState}`))
            }
          }
        }
        
        // 设置超时
        const timeoutId = setTimeout(() => {
          if (!isComplete) {
            console.warn('⏰ SSE请求超时')
            eventSource.close()
            
            // 如果没有收到数据且还有重试次数，则重试
            if (!hasReceivedData && retryCount < maxRetries) {
              console.log(`🔄 超时重试... (${retryCount + 1}/${maxRetries})`)
              this.sendMessageSSE(data, onUpdate, retryCount + 1).then(resolve).catch(reject)
            } else {
              reject(new Error('请求超时'))
            }
          }
                 }, 180000) // 180秒超时
        
        // 清理超时
        const cleanup = () => {
          clearTimeout(timeoutId)
        }
        
        const originalResolve = resolve
        const originalReject = reject
        resolve = (...args) => {
          cleanup()
          originalResolve(...args)
        }
        reject = (...args) => {
          cleanup()
          originalReject(...args)
        }
      })
    },

    // 备用的fetch流式实现 (使用前端代理避免CORS)
    async sendMessageFetchStream(data, onUpdate, retryCount = 0) {
      const maxRetries = 2
      let abortController = null
      let timeoutId = null
      let dataTimeoutId = null
      
      try {
        const backendURL = chatApiService.defaults?.baseURL || '' // 使用前端代理，避免CORS问题
        
        const requestData = {
          conversation_id: data.conversationId,
          character_id: data.characterId,
          content: data.content,
          message_type: data.type === 'text' ? 1 : (data.type === 'voice' ? 2 : 1),
          user_id: 1
        }
        
        console.log(`📤 发送Fetch流式请求到: ${backendURL}/api/chat/send (尝试 ${retryCount + 1}/${maxRetries + 1})`)
        console.log('📤 请求参数:', requestData)
        console.log('🔍 使用代理模式，baseURL:', backendURL)
        
        const queryParams = new URLSearchParams(requestData).toString()
        const url = `${backendURL}/api/chat/send?${queryParams}`
        
        // 创建手动控制的AbortController，更稳定
        abortController = new AbortController()
        timeoutId = setTimeout(() => {
          console.warn('⏰ Fetch请求超时(180秒)，取消请求')
          if (abortController) {
            abortController.abort()
          }
        }, 180000) // 180秒超时
        
        const response = await fetch(url, {
          method: 'GET',
          headers: {
            'Accept': 'text/event-stream',
            'Cache-Control': 'no-cache',
            'Connection': 'keep-alive', // 明确要求保持连接
          },
          signal: abortController.signal
        })
        
        // 请求成功，清除超时
        if (timeoutId) {
          clearTimeout(timeoutId)
          timeoutId = null
        }
        
        if (!response.ok) {
          throw new Error(`HTTP ${response.status}: ${response.statusText}`)
        }
        
        console.log('✅ Fetch流式连接已建立')
        console.log('🔍 响应头:', {
          'content-type': response.headers.get('content-type'),
          'cache-control': response.headers.get('cache-control'),
          'connection': response.headers.get('connection')
        })
        
        const reader = response.body.getReader()
        const decoder = new TextDecoder()
        
        // 改进数据超时检查逻辑
        let lastDataTime = Date.now()
        let isStreamActive = true
        
        dataTimeoutId = setInterval(() => {
          if (!isStreamActive) {
            clearInterval(dataTimeoutId)
            return
          }
          
          const timeSinceLastData = Date.now() - lastDataTime
          if (timeSinceLastData > 60000) { // 增加到60秒
            console.warn(`⚠️ ${timeSinceLastData/1000}秒内没有收到数据，连接可能异常`)
            isStreamActive = false
            clearInterval(dataTimeoutId)
            if (abortController && !abortController.signal.aborted) {
              console.log('🔄 数据超时，但不中止请求，等待自然结束')
              // 不调用 abort()，让连接自然结束
            }
          } else if (timeSinceLastData > 30000) {
            console.log(`🔍 ${timeSinceLastData/1000}秒没有收到新数据，连接仍在等待...`)
          }
        }, 10000) // 每10秒检查一次
        
        let aiResponse = ''
        let messageId = null
        let buffer = ''
        let hasReceivedData = false
        let processedMessageIds = new Set() // 用于去重
        
        try {
          while (true) {
            const { done, value } = await reader.read()
            
            if (done) {
              console.log('✅ 流式读取完成')
              isStreamActive = false
              break
            }
            
            hasReceivedData = true
            lastDataTime = Date.now() // 更新最后接收数据的时间
            
            // 解码数据
            const chunk = decoder.decode(value, { stream: true })
            buffer += chunk
            
            // 处理SSE数据格式
            const lines = buffer.split('\n')
            buffer = lines.pop() || '' // 保留最后一行（可能不完整）
            
            for (const line of lines) {
              if (line.startsWith('data: ')) {
                const data = line.substring(6)
                if (data.trim()) {
                  try {
                    const responseData = JSON.parse(data)
                    
                    // 去重检查：如果是message类型，检查是否已经处理过
                    if (responseData.type === 'message' && responseData.message_id) {
                      const uniqueKey = `${responseData.message_id}-${responseData.content?.length || 0}`
                      if (processedMessageIds.has(uniqueKey)) {
                        console.log('⚠️ 跳过重复消息:', uniqueKey)
                        continue
                      }
                      processedMessageIds.add(uniqueKey)
                    }
                    
                    console.log('📨 收到流式数据:', responseData)
                    
                    if (responseData.type === 'message') {
                      // 直接累加增量内容（后端现在只发送增量）
                      const newContent = responseData.content || ''
                      if (newContent) {
                        aiResponse += newContent
                        console.log(`📨 累加增量内容: +${newContent.length}字符, 总长度: ${aiResponse.length}`)
                      }
                      messageId = responseData.message_id
                      
                      if (onUpdate) {
                        onUpdate(aiResponse)
                      }
                    } else if (responseData.type === 'complete' || responseData.type === 'done') {
                      console.log(`✅ 收到${responseData.type}事件，流式回复完成`)
                      console.log('✅ 事件数据:', responseData)
                      
                      // 如果done/complete事件包含完整内容，优先使用它
                      const finalContent = responseData.content || aiResponse
                      console.log('✅ 最终AI回复内容长度:', finalContent.length)
                      console.log('✅ 最终AI回复内容预览:', finalContent.substring(0, 100) + (finalContent.length > 100 ? '...' : ''))
                      
                      isStreamActive = false
                      if (dataTimeoutId) {
                        clearInterval(dataTimeoutId)
                        dataTimeoutId = null
                      }
                      
                      return {
                        data: {
                          ai_message: {
                            content: finalContent, // 使用最终内容
                            timestamp: new Date().toISOString()
                          }
                        }
                      }
                    } else if (responseData.type === 'error') {
                      isStreamActive = false
                      if (dataTimeoutId) {
                        clearInterval(dataTimeoutId)
                        dataTimeoutId = null
                      }
                      throw new Error(responseData.message || '聊天请求失败')
                    }
                  } catch (parseError) {
                    console.warn('⚠️ 解析SSE数据失败:', parseError, data)
                  }
                }
              }
            }
          }
          
          // 如果循环结束但没有收到complete，返回当前内容
          if (aiResponse) {
            console.log('⚠️ 流结束但没有收到complete事件，返回已接收内容:', aiResponse.length)
            return {
              data: {
                ai_message: {
                  content: aiResponse,
                  timestamp: new Date().toISOString()
                }
              }
            }
          } else if (hasReceivedData) {
            console.warn('⚠️ 收到数据但没有有效内容')
            throw new Error('收到数据但没有有效的回复内容')
          } else {
            console.error('❌ 没有收到任何数据')
            throw new Error('没有收到任何回复数据')
          }
          
        } finally {
          isStreamActive = false
          if (dataTimeoutId) {
            clearInterval(dataTimeoutId)
            dataTimeoutId = null
          }
          try {
            reader.releaseLock()
          } catch (e) {
            console.warn('⚠️ 释放reader锁失败:', e.message)
          }
        }
        
      } catch (error) {
        console.error('❌ Fetch流式请求失败:', error)
        
        // 清理定时器
        if (timeoutId) {
          clearTimeout(timeoutId)
          timeoutId = null
        }
        if (dataTimeoutId) {
          clearInterval(dataTimeoutId)
          dataTimeoutId = null
        }
        
        // 改进的重试逻辑
        const shouldRetry = retryCount < maxRetries && 
                           !error.name?.includes('AbortError') && 
                           !error.message?.includes('timeout') &&
                           !error.message?.includes('aborted')
        
        if (shouldRetry) {
          console.log(`🔄 准备重试... (${retryCount + 1}/${maxRetries})，错误类型: ${error.name}`)
          await new Promise(resolve => setTimeout(resolve, 2000 + retryCount * 1000)) // 递增延迟
          return this.sendMessageFetchStream(data, onUpdate, retryCount + 1)
        }
        
        // 如果是AbortError，提供更友好的错误信息
        if (error.name === 'AbortError') {
          throw new Error('请求被中断，可能是网络超时或连接问题')
        }
        
        throw error
      }
    },

    // 开始录音
    async startRecording() {
      try {
        this.initVoiceControllers()
        
        const success = await this.voiceRecorder.startRecording()
        if (success) {
          this.isRecording = true
          this.transcript = ''
        }
        return success
      } catch (error) {
        this.error = error.message
        console.error('开始录音失败:', error)
        return false
      }
    },

    // 停止录音
    async stopRecording() {
      try {
        if (this.voiceRecorder && this.isRecording) {
          this.voiceRecorder.stopRecording()
          this.isRecording = false
          
          // 语音转文字由VoiceRecorder类自动处理
          console.log('🎤 录音已停止，等待语音转文字...')
        }
      } catch (error) {
        this.error = error.message
        console.error('停止录音失败:', error)
      }
    },

    // 播放语音
    async playVoice(text) {
      try {
        this.initVoiceControllers()
        
        // 这里可以添加TTS逻辑
        // const audioUrl = await speechApi.textToSpeech({ text, characterId: this.currentConversation?.characterId })
        
        const success = await this.voicePlayer.play(text) // 临时使用浏览器TTS
        if (success) {
          this.isSpeaking = true
        }
        return success
      } catch (error) {
        this.error = error.message
        console.error('播放语音失败:', error)
        return false
      }
    },

    // 停止播放
    stopPlaying() {
      if (this.voicePlayer && this.isSpeaking) {
        this.voicePlayer.stop()
        this.isSpeaking = false
      }
    },

    // 更新语音设置
    updateVoiceSettings(settings) {
      this.voiceSettings = { ...this.voiceSettings, ...settings }
    },

    // 清除错误
    clearError() {
      this.error = null
    },

    // 删除对话
    async deleteConversation(conversationId) {
      try {
        await chatApiService.deleteConversation(conversationId)
        
        this.conversations = this.conversations.filter(conv => conv.id !== conversationId)
        
        if (this.currentConversation?.id === conversationId) {
          this.currentConversation = null
        }
      } catch (error) {
        this.error = error.message
        console.error('删除对话失败:', error)
        throw error
      }
    },

    // 强制停止所有流式回复
    forceStopStreaming() {
      console.log('🛑 强制停止所有流式回复')
      if (this.currentConversation && this.currentConversation.messages) {
        this.currentConversation.messages.forEach((message, index) => {
          if (message.isStreaming) {
            console.log(`🛑 停止消息 ${message.id} 的流式状态`)
            message.isStreaming = false
          }
        })
      }
      this.isLoading = false
      console.log('�� 所有流式回复已停止')
    }
  }
})

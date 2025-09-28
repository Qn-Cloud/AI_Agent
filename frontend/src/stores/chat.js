import { defineStore } from 'pinia'
import { chatApiService as chatApi, aiApiService as aiApi, speechApiService as speechApi, VoiceRecorder, VoicePlayer } from '../services'

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
    // 当前对话
    currentConversation: null,
    // 所有对话历史（初始为空，从后端加载）
    conversations: [],
    // 统计数据（从后端获取）
    stats: {
      conversationTotal: 0,
      messageCount: 0,
      characterCount: 0,
      activeDays: 0
    },
    // 语音交互状态
    isRecording: false,
    isSpeaking: false,
    isLoading: false,
    // 临时录音文本
    transcript: '',
    // 语音设置
    voiceSettings: {
      volume: 80,
      rate: 1.0,
      pitch: 1.0,
      voice: null
    },
    // 语音控制器实例
    voiceRecorder: null,
    voicePlayer: null,
    // 错误状态
    error: null
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
      }
      if (!this.voicePlayer) {
        this.voicePlayer = new VoicePlayer()
      }
    },

    // 从后端加载对话历史
    async loadConversationHistory(params = {}) {
      try {
        this.isLoading = true
        this.error = null
        
        console.log('正在从后端加载对话历史...', params)
        
        const response = await chatApi.getConversationHistory({
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
        const response = await chatApi.getConversationHistory({
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
        
        const response = await chatApi.createConversation({
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

    // 发送消息
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
        
        // 发送到后端
        const response = await chatApi.sendMessage({
          conversationId: this.currentConversation.id,
          characterId: this.currentConversation.characterId,
          content,
          type
        })
        
        if (response && response.data) {
          // 添加AI回复
          const aiMessage = {
            id: response.data.ai_message.id,
            type: 'ai',
            content: response.data.ai_message.content,
            timestamp: new Date(response.data.ai_message.timestamp)
          }
          
          this.currentConversation.messages.push(aiMessage)
          this.currentConversation.lastUpdate = new Date()
          
          return aiMessage
        }
      } catch (error) {
        this.error = error.message
        console.error('发送消息失败:', error)
        throw error
      } finally {
        this.isLoading = false
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
          
          // 这里可以添加语音识别逻辑
          // const transcript = await speechApi.speechToText(audioData)
          // this.transcript = transcript
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
        await chatApi.deleteConversation(conversationId)
        
        this.conversations = this.conversations.filter(conv => conv.id !== conversationId)
        
        if (this.currentConversation?.id === conversationId) {
          this.currentConversation = null
        }
      } catch (error) {
        this.error = error.message
        console.error('删除对话失败:', error)
        throw error
      }
    }
  }
})

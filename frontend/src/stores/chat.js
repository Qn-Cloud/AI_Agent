import { defineStore } from 'pinia'
import { chatApiService as chatApi, aiApiService as aiApi, speechApiService as speechApi, VoiceRecorder, VoicePlayer } from '../services'

export const useChatStore = defineStore('chat', {
  state: () => ({
    // 当前对话
    currentConversation: null,
    // 所有对话历史（初始为空，从后端加载）
    conversations: [
      {
        id: 'conv-1',
        characterId: 'harry-potter',
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
        characterId: 'socrates',
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
    ],
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
      autoSpeak: true
    },
    // 语音控制器
    voiceRecorder: null,
    voicePlayer: null,
    // 错误信息
    error: null
  }),

  getters: {
    // 按时间排序的对话
    sortedConversations: (state) => {
      return [...state.conversations].sort((a, b) => 
        new Date(b.lastUpdate) - new Date(a.lastUpdate)
      )
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

    // 当前对话的消息
    currentMessages: (state) => {
      return state.currentConversation?.messages || []
    },

    // 最后一条AI消息
    lastAIMessage: (state) => {
      if (!state.currentConversation) return null
      const messages = state.currentConversation.messages
      for (let i = messages.length - 1; i >= 0; i--) {
        if (messages[i].type === 'ai') {
          return messages[i]
        }
      }
      return null
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
        
        if (response && response.data) {
          // 清空当前对话列表（如果是第一页）
          if (!params.page || params.page === 1) {
            this.conversations = []
          }
          
          // 转换后端数据格式到前端格式
          const conversations = response.data.list?.map(conv => ({
            id: conv.id,
            characterId: conv.character_id,
            title: conv.title,
            startTime: new Date(conv.start_time),
            lastUpdate: new Date(conv.last_message_time),
            messageCount: conv.message_count,
            status: conv.status,
            messages: [] // 消息会在需要时单独加载
          })) || []
          
          // 添加到对话列表
          this.conversations.push(...conversations)
          
          console.log('成功加载对话历史:', conversations.length, '条对话')
          return {
            conversations,
            total: response.data.total || 0,
            hasMore: response.data.has_more || false
          }
        }
      } catch (error) {
        this.error = error.message
        console.error('加载对话历史失败:', error)
        
        // 如果API失败，使用默认的模拟数据
        console.log('API失败，使用模拟数据')
        if (this.conversations.length === 0) {
          // 保持原有的模拟数据
          this.conversations = [
            {
              id: 'conv-1',
              characterId: 'harry-potter',
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
              characterId: 'socrates',
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
          title: `新对话 ${new Date().toLocaleString()}`
        })
        
        if (response.data) {
          const newConversation = {
            id: response.data.conversation_id,
            characterId,
            title: response.data.title,
            startTime: new Date(response.data.created_at),
            lastUpdate: new Date(response.data.updated_at),
            messages: []
          }
          
          this.conversations.unshift(newConversation)
          this.currentConversation = newConversation
          
          return newConversation
        }
      } catch (error) {
        this.error = error.message
        console.error('创建对话失败:', error)
        
        // 如果API失败，创建本地对话
        const localConversation = {
          id: `local-conv-${Date.now()}`,
          characterId,
          title: `新对话 ${new Date().toLocaleString()}`,
          startTime: new Date(),
          lastUpdate: new Date(),
          messages: [],
          isLocal: true
        }
        
        this.conversations.unshift(localConversation)
        this.currentConversation = localConversation
        return localConversation
      } finally {
        this.isLoading = false
      }
    },

    // 选择现有对话
    async selectConversation(conversationId) {
      try {
        // 先从本地查找
        let conversation = this.conversations.find(c => c.id === conversationId)
        
        if (conversation) {
          this.currentConversation = conversation
          
          // 如果消息为空，从API加载消息
          if (conversation.messages.length === 0 && !conversation.isLocal) {
            await this.loadMessages(conversationId)
          }
        } else {
          // 从API获取对话详情
          const response = await chatApi.getConversation(conversationId)
          if (response.data) {
            conversation = response.data.conversation
            this.conversations.unshift(conversation)
            this.currentConversation = conversation
            await this.loadMessages(conversationId)
          }
        }
      } catch (error) {
        this.error = error.message
        console.error('选择对话失败:', error)
      }
    },

    // 加载对话消息
    async loadMessages(conversationId) {
      try {
        const response = await chatApi.getMessages(conversationId)
        if (response.data && response.data.messages) {
          const conversation = this.conversations.find(c => c.id === conversationId)
          if (conversation) {
            conversation.messages = response.data.messages
          }
        }
      } catch (error) {
        console.error('加载消息失败:', error)
      }
    },

    // 发送消息
    async sendMessage(content, type = 'user') {
      if (!this.currentConversation) return

      try {
        this.isLoading = true
        this.error = null

        const userMessage = {
          id: `msg-${Date.now()}`,
          type,
          content,
          timestamp: new Date()
        }

        this.currentConversation.messages.push(userMessage)
        this.currentConversation.lastUpdate = new Date()

        // 如果是用户消息，发送到后端获取AI回复
        if (type === 'user') {
          await this.getAIReply(content)
        }

        return userMessage
      } catch (error) {
        this.error = error.message
        console.error('发送消息失败:', error)
      } finally {
        this.isLoading = false
      }
    },

    // 获取AI回复
    async getAIReply(userMessage) {
      try {
        const response = await chatApi.sendMessage({
          conversationId: this.currentConversation.id,
          characterId: this.currentConversation.characterId,
          content: userMessage
        })

        if (response.data) {
          const aiMessage = {
            id: response.data.message_id || `msg-${Date.now() + 1}`,
            type: 'ai',
            content: response.data.reply,
            timestamp: new Date()
          }

          this.currentConversation.messages.push(aiMessage)
          this.currentConversation.lastUpdate = new Date()

          // 如果开启自动播放，则播放语音
          if (this.voiceSettings.autoSpeak) {
            await this.speakMessage(response.data.reply)
          }

          return aiMessage
        }
      } catch (error) {
        console.error('获取AI回复失败:', error)
        // 使用模拟回复作为降级方案
        await this.simulateAIReply(userMessage)
      }
    },

    // 模拟AI回复（使用假数据）
    async simulateAIReply(userMessage) {
      this.isLoading = true
      
      // 模拟网络延迟
      await new Promise(resolve => setTimeout(resolve, 1000 + Math.random() * 2000))

      // 生成模拟回复
      const replies = [
        '这是一个很有趣的问题！让我来为你详细解答...',
        '根据我的理解，这个话题涉及到很多方面...',
        '你提出了一个很深刻的观点，我想分享一些我的想法...',
        '这让我想起了一个相关的故事...',
        '从另一个角度来看，我们可以这样思考...'
      ]

      const randomReply = replies[Math.floor(Math.random() * replies.length)]
      
      this.sendMessage(randomReply, 'ai')
      this.isLoading = false

      // 如果开启自动播放，则播放语音
      if (this.voiceSettings.autoSpeak) {
        this.speakMessage(randomReply)
      }
    },

    // 删除对话
    deleteConversation(conversationId) {
      const index = this.conversations.findIndex(c => c.id === conversationId)
      if (index > -1) {
        this.conversations.splice(index, 1)
        if (this.currentConversation?.id === conversationId) {
          this.currentConversation = null
        }
      }
    },

    // 清空当前对话
    clearCurrentConversation() {
      if (this.currentConversation) {
        this.currentConversation.messages = []
        this.currentConversation.lastUpdate = new Date()
      }
    },

    // 语音识别控制
    async startRecording() {
      try {
        this.initVoiceControllers()
        this.isRecording = true
        this.transcript = ''
        this.error = null
        
        await this.voiceRecorder.startRecording()
      } catch (error) {
        this.isRecording = false
        this.error = error.message
        console.error('录音启动失败:', error)
        throw error
      }
    },

    async stopRecording() {
      if (!this.isRecording || !this.voiceRecorder) return
      
      try {
        this.isRecording = false
        const audioBlob = await this.voiceRecorder.stopRecording()
        
        if (audioBlob) {
          // 调用语音识别API
          const response = await speechApi.speechToText(audioBlob, {
            language: 'zh-CN'
          })
          
          if (response.data && response.data.text) {
            this.transcript = response.data.text
            // 自动发送识别的文本
            await this.sendMessage(this.transcript)
            this.transcript = ''
          }
        }
      } catch (error) {
        this.error = error.message
        console.error('语音识别失败:', error)
        
        // 使用模拟识别作为降级方案
        this.simulateVoiceRecognition()
      }
    },

    // 模拟语音识别（降级方案）
    simulateVoiceRecognition() {
      const mockPhrases = [
        '你好，很高兴见到你',
        '今天天气真不错',
        '你能帮我解答一个问题吗',
        '谢谢你的帮助',
        '这个话题很有趣'
      ]
      
      setTimeout(() => {
        this.transcript = mockPhrases[Math.floor(Math.random() * mockPhrases.length)]
        this.sendMessage(this.transcript)
        this.transcript = ''
      }, 2000)
    },

    // 语音播放
    async speakMessage(text, characterId = null) {
      try {
        this.initVoiceControllers()
        this.isSpeaking = true
        this.error = null
        
        // 如果有角色ID，使用角色的语音设置
        const voiceSettings = characterId ? 
          this.getCharacterVoiceSettings(characterId) : 
          this.voiceSettings
        
        // 调用TTS API
        const response = await speechApi.textToSpeech({
          text,
          characterId,
          voiceSettings: {
            rate: voiceSettings.rate,
            pitch: voiceSettings.pitch,
            volume: voiceSettings.volume / 100
          }
        })
        
        if (response.data && response.data.audio_url) {
          // 播放语音
          await this.voicePlayer.play(response.data.audio_url, {
            volume: voiceSettings.volume / 100,
            onEnded: () => {
              this.isSpeaking = false
            },
            onError: (error) => {
              console.error('语音播放失败:', error)
              this.isSpeaking = false
            }
          })
        }
      } catch (error) {
        this.isSpeaking = false
        this.error = error.message
        console.error('语音合成失败:', error)
        
        // 使用Web Speech API作为降级方案
        this.speakWithWebAPI(text)
      }
    },

    // 使用Web Speech API播放语音（降级方案）
    speakWithWebAPI(text) {
      if ('speechSynthesis' in window) {
        const utterance = new SpeechSynthesisUtterance(text)
        utterance.lang = 'zh-CN'
        utterance.rate = this.voiceSettings.rate
        utterance.pitch = this.voiceSettings.pitch
        utterance.volume = this.voiceSettings.volume / 100
        
        utterance.onend = () => {
          this.isSpeaking = false
        }
        
        utterance.onerror = () => {
          this.isSpeaking = false
        }
        
        speechSynthesis.speak(utterance)
      } else {
        // 如果都不支持，使用计时器模拟
        setTimeout(() => {
          this.isSpeaking = false
        }, Math.max(1000, text.length * 50))
      }
    },

    stopSpeaking() {
      this.isSpeaking = false
      
      if (this.voicePlayer) {
        this.voicePlayer.stop()
      }
      
      if ('speechSynthesis' in window) {
        speechSynthesis.cancel()
      }
    },

    // 获取角色的语音设置
    getCharacterVoiceSettings(characterId) {
      // 这里需要从character store获取角色的语音设置
      // 暂时返回默认设置
      return this.voiceSettings
    },

    // 更新语音设置
    updateVoiceSettings(settings) {
      this.voiceSettings = { ...this.voiceSettings, ...settings }
    },

    // 搜索对话
    searchConversations(query) {
      if (!query.trim()) return this.conversations

      return this.conversations.filter(conversation => {
        const titleMatch = conversation.title.toLowerCase().includes(query.toLowerCase())
        const contentMatch = conversation.messages.some(msg => 
          msg.content.toLowerCase().includes(query.toLowerCase())
        )
        return titleMatch || contentMatch
      })
    }
  }
})

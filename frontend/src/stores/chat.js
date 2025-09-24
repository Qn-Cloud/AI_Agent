import { defineStore } from 'pinia'

export const useChatStore = defineStore('chat', {
  state: () => ({
    // 当前对话
    currentConversation: null,
    // 所有对话历史
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
    }
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
    // 开始新对话
    startNewConversation(characterId) {
      const newConversation = {
        id: `conv-${Date.now()}`,
        characterId,
        title: `新对话 ${new Date().toLocaleString()}`,
        startTime: new Date(),
        lastUpdate: new Date(),
        messages: []
      }
      this.conversations.push(newConversation)
      this.currentConversation = newConversation
      return newConversation
    },

    // 选择现有对话
    selectConversation(conversationId) {
      const conversation = this.conversations.find(c => c.id === conversationId)
      if (conversation) {
        this.currentConversation = conversation
      }
    },

    // 发送消息
    sendMessage(content, type = 'user') {
      if (!this.currentConversation) return

      const message = {
        id: `msg-${Date.now()}`,
        type,
        content,
        timestamp: new Date()
      }

      this.currentConversation.messages.push(message)
      this.currentConversation.lastUpdate = new Date()

      // 如果是用户消息，模拟AI回复
      if (type === 'user') {
        this.simulateAIReply(content)
      }

      return message
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
    startRecording() {
      this.isRecording = true
      this.transcript = ''
      // 这里会调用实际的语音识别API
      this.simulateVoiceRecognition()
    },

    stopRecording() {
      this.isRecording = false
      if (this.transcript.trim()) {
        this.sendMessage(this.transcript)
        this.transcript = ''
      }
    },

    // 模拟语音识别
    simulateVoiceRecognition() {
      const mockPhrases = [
        '你好，很高兴见到你',
        '今天天气真不错',
        '你能帮我解答一个问题吗',
        '谢谢你的帮助',
        '这个话题很有趣'
      ]
      
      setTimeout(() => {
        if (this.isRecording) {
          this.transcript = mockPhrases[Math.floor(Math.random() * mockPhrases.length)]
        }
      }, 2000)
    },

    // 语音播放
    speakMessage(text) {
      this.isSpeaking = true
      
      // 这里会调用实际的TTS API
      setTimeout(() => {
        this.isSpeaking = false
      }, Math.max(1000, text.length * 50)) // 模拟播放时间
    },

    stopSpeaking() {
      this.isSpeaking = false
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

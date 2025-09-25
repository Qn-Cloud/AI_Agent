import api from './api'

export const aiApi = {
  // AI对话
  aiChat(data) {
    return api.post('/api/ai/chat', {
      character_id: data.characterId,
      conversation_id: data.conversationId,
      message: data.message,
      context: data.context || [],
      model: data.model || 'gpt-3.5-turbo',
      options: data.options || {
        temperature: 0.7,
        max_tokens: 2048
      }
    })
  },

  // 获取AI模型列表
  getAIModels() {
    return api.get('/api/ai/models')
  },

  // 获取用户使用统计
  getUsageStats(params = {}) {
    return api.get('/api/ai/usage', {
      params: {
        start_date: params.startDate,
        end_date: params.endDate,
        period: params.period || 'day'
      }
    })
  }
}

// 聊天会话管理器
export class ChatSession {
  constructor(characterId, conversationId = null) {
    this.characterId = characterId
    this.conversationId = conversationId
    this.messageHistory = []
    this.isProcessing = false
  }

  // 发送消息并获取AI回复
  async sendMessage(content, options = {}) {
    if (this.isProcessing) {
      throw new Error('正在处理中，请稍后再试')
    }

    try {
      this.isProcessing = true

      // 构建上下文
      const context = this.messageHistory.slice(-10) // 保留最近10条消息作为上下文

      // 调用AI接口
      const response = await aiApi.aiChat({
        characterId: this.characterId,
        conversationId: this.conversationId,
        message: content,
        context,
        ...options
      })

      // 更新消息历史
      const userMessage = {
        id: `msg-${Date.now()}`,
        type: 'user',
        content,
        timestamp: new Date()
      }

      const aiMessage = {
        id: `msg-${Date.now() + 1}`,
        type: 'ai',
        content: response.data.reply,
        timestamp: new Date()
      }

      this.messageHistory.push(userMessage, aiMessage)

      return {
        userMessage,
        aiMessage,
        conversationId: response.data.conversation_id
      }
    } finally {
      this.isProcessing = false
    }
  }

  // 更新对话ID
  setConversationId(conversationId) {
    this.conversationId = conversationId
  }

  // 获取消息历史
  getMessageHistory() {
    return [...this.messageHistory]
  }

  // 清空消息历史
  clearHistory() {
    this.messageHistory = []
  }

  // 获取处理状态
  getProcessingState() {
    return this.isProcessing
  }
}

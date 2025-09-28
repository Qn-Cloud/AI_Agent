import { chatApi } from './apiFactory'

export const chatApiService = {
  // 发送消息 (SSE方式)
  sendMessageSSE(data) {
    return new Promise((resolve, reject) => {
      const baseURL = chatApi.defaults?.baseURL || ''
      const url = `${baseURL}/api/chat/send`
      
      const requestData = {
        conversation_id: data.conversationId,
        character_id: data.characterId,
        content: data.content,
        type: data.type || 'text',
        user_id: 1 // 暂时固定为1
      }
      
      console.log('📤 发送SSE聊天请求:', requestData)
      
      // 创建EventSource连接
      const eventSource = new EventSource(url + '?' + new URLSearchParams(requestData))
      
      let aiResponse = ''
      let messageId = null
      let isComplete = false
      
      eventSource.onopen = () => {
        console.log('🔗 SSE连接已建立')
      }
      
      eventSource.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          console.log('📨 收到SSE消息:', data)
          
          if (data.type === 'message') {
            // 接收到AI回复的片段
            aiResponse += data.content
            messageId = data.message_id
            
            // 触发实时更新回调
            if (data.onUpdate) {
              data.onUpdate(aiResponse)
            }
          } else if (data.type === 'complete') {
            // 回复完成
            isComplete = true
            eventSource.close()
            
            resolve({
              data: {
                ai_message: {
                  id: messageId,
                  content: aiResponse,
                  timestamp: new Date().toISOString()
                }
              }
            })
          } else if (data.type === 'error') {
            // 发生错误
            eventSource.close()
            reject(new Error(data.message || '聊天请求失败'))
          }
        } catch (error) {
          console.error('❌ 解析SSE消息失败:', error, event.data)
        }
      }
      
      eventSource.onerror = (error) => {
        console.error('❌ SSE连接错误:', error)
        eventSource.close()
        
        if (!isComplete) {
          reject(new Error('连接中断'))
        }
      }
      
      // 设置超时
      setTimeout(() => {
        if (!isComplete) {
          eventSource.close()
          reject(new Error('请求超时'))
        }
      }, 60000) // 60秒超时
    })
  },

  // 保留原有方法作为备用
  sendMessage(data) {
    return this.sendMessageSSE(data)
  },

  // 创建新对话
  createConversation(data) {
    return chatApi.post('/api/chat/conversation', {
      character_id: data.characterId,
      title: data.title || '新对话'
    })
  },

  // 获取对话详情
  getConversation(id) {
    return chatApi.get(`/api/chat/conversation/${id}`)
  },

  // 获取对话历史
  getConversationHistory(params = {}) {
    const requestData = {
      page: params.page || 1,
      page_size: params.pageSize || 20,
      character_id: params.characterId,
      user_id: 1, // 固定设置为1，后续从JWT token获取
      start_time: params.startTime,
      end_time: params.endTime
    }
    
    console.log('📤 发送getConversationHistory请求:', requestData)
    
    return chatApi.post('/api/chat/history', requestData)
  },

  // 保留原有方法名作为别名，避免前端代码大量修改
  getConversationList(params = {}) {
    return this.getConversationHistory(params)
  },

  // 获取对话消息历史
  getMessages(conversationId, params = {}) {
    const requestParams = {
      conversation_id: conversationId,
      page: params.page || 1,
      page_size: params.pageSize || 50,
      before_id: params.beforeId // 用于分页
    }
    
    console.log('📤 发送getMessages请求:', {
      url: '/api/chat/messages',
      params: requestParams
    })
    
    return chatApi.get('/api/chat/messages', {
      params: requestParams
    })
  },

  // 删除对话
  deleteConversation(id) {
    return chatApi.delete(`/api/chat/conversation/${id}`)
  },

  // 清空对话消息
  clearMessages(id) {
    return chatApi.delete(`/api/chat/conversation/${id}/messages`)
  },

  // 更新对话标题
  updateConversationTitle(id, title) {
    return chatApi.put(`/api/chat/conversation/${id}/title`, {
      title
    })
  },

  // 搜索对话
  searchConversations(params) {
    return chatApi.get('/api/chat/search', {
      params: {
        keyword: params.keyword,
        page: params.page || 1,
        page_size: params.pageSize || 20
      }
    })
  },

  // 导出对话记录
  exportConversation(id, format = 'json') {
    return chatApi.get(`/api/chat/conversation/${id}/export`, {
      params: { format }
    })
  },

  // 批量删除对话
  batchDeleteConversations(ids) {
    return chatApi.post('/api/chat/conversations/batch-delete', {
      conversation_ids: ids
    })
  }
}

// 保持向后兼容
export const chatApi_old = chatApiService
export default chatApiService

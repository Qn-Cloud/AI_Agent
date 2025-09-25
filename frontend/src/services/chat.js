import api from './api'

export const chatApi = {
  // 发送消息
  sendMessage(data) {
    return api.post('/api/chat/message', {
      conversation_id: data.conversationId,
      character_id: data.characterId,
      content: data.content,
      type: data.type || 'text'
    })
  },

  // 创建新对话
  createConversation(data) {
    return api.post('/api/chat/conversation', {
      character_id: data.characterId,
      title: data.title || '新对话'
    })
  },

  // 获取对话详情
  getConversation(id) {
    return api.get(`/api/chat/conversation/${id}`)
  },

  // 获取对话列表
  getConversationList(params = {}) {
    return api.get('/api/chat/conversations', {
      params: {
        page: params.page || 1,
        page_size: params.pageSize || 20,
        character_id: params.characterId
      }
    })
  },

  // 获取对话消息历史
  getMessages(conversationId, params = {}) {
    return api.get(`/api/chat/conversation/${conversationId}/messages`, {
      params: {
        page: params.page || 1,
        page_size: params.pageSize || 50,
        before_id: params.beforeId // 用于分页
      }
    })
  },

  // 删除对话
  deleteConversation(id) {
    return api.delete(`/api/chat/conversation/${id}`)
  },

  // 清空对话消息
  clearMessages(id) {
    return api.delete(`/api/chat/conversation/${id}/messages`)
  },

  // 更新对话标题
  updateConversationTitle(id, title) {
    return api.put(`/api/chat/conversation/${id}/title`, {
      title
    })
  },

  // 搜索对话
  searchConversations(params) {
    return api.get('/api/chat/search', {
      params: {
        keyword: params.keyword,
        page: params.page || 1,
        page_size: params.pageSize || 20
      }
    })
  },

  // 导出对话记录
  exportConversation(id, format = 'json') {
    return api.get(`/api/chat/conversation/${id}/export`, {
      params: { format }
    })
  },

  // 批量删除对话
  batchDeleteConversations(ids) {
    return api.post('/api/chat/conversations/batch-delete', {
      conversation_ids: ids
    })
  }
}

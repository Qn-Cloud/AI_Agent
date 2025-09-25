import api from './api'

export const characterApi = {
  // 获取角色列表
  getCharacterList(params = {}) {
    return api.get('/api/character/list', {
      params: {
        page: params.page || 1,
        page_size: params.pageSize || 20,
        category: params.category,
        tags: params.tags
      }
    })
  },

  // 搜索角色
  searchCharacters(params) {
    return api.get('/api/character/search', {
      params: {
        keyword: params.keyword,
        page: params.page || 1,
        page_size: params.pageSize || 20
      }
    })
  },

  // 获取角色详情
  getCharacterDetail(id) {
    return api.get(`/api/character/${id}`)
  },

  // 获取推荐角色
  getRecommendedCharacters(params = {}) {
    return api.get('/api/character/recommended', {
      params: {
        count: params.count || 10
      }
    })
  },

  // 获取角色分类
  getCharacterCategories() {
    return api.get('/api/character/categories')
  },

  // 获取热门角色
  getPopularCharacters(params = {}) {
    return api.get('/api/character/popular', {
      params: {
        page: params.page || 1,
        page_size: params.pageSize || 20
      }
    })
  },

  // 获取角色标签
  getCharacterTags() {
    return api.get('/api/character/tags')
  },

  // 创建自定义角色（需要认证）
  createCharacter(data) {
    return api.post('/api/character', {
      name: data.name,
      description: data.description,
      avatar: data.avatar,
      tags: data.tags,
      prompt: data.prompt,
      personality: data.personality,
      voice_settings: data.voiceSettings
    })
  },

  // 更新角色信息（需要认证）
  updateCharacter(id, data) {
    return api.put(`/api/character/${id}`, data)
  },

  // 删除角色（需要认证）
  deleteCharacter(id) {
    return api.delete(`/api/character/${id}`)
  },

  // 收藏/取消收藏角色（需要认证）
  toggleFavorite(id) {
    return api.post(`/api/character/${id}/favorite`)
  },

  // 获取我的收藏角色（需要认证）
  getMyFavorites(params = {}) {
    return api.get('/api/character/favorites', {
      params: {
        page: params.page || 1,
        page_size: params.pageSize || 20
      }
    })
  },

  // 获取我创建的角色（需要认证）
  getMyCharacters(params = {}) {
    return api.get('/api/character/my', {
      params: {
        page: params.page || 1,
        page_size: params.pageSize || 20
      }
    })
  },

  // 更新角色提示词（需要认证）
  updatePrompt(id, prompt) {
    return api.put(`/api/character/${id}/prompt`, {
      prompt
    })
  },

  // 更新角色性格设置（需要认证）
  updatePersonality(id, personality) {
    return api.put(`/api/character/${id}/personality`, {
      personality
    })
  },

  // 更新语音设置（需要认证）
  updateVoiceSettings(id, voiceSettings) {
    return api.put(`/api/character/${id}/voice`, {
      voice_settings: voiceSettings
    })
  }
}

import { characterApi } from './apiFactory'

export const characterApiService = {
  // 获取角色列表
  getCharacterList(params = {}) {
    return characterApi.get('/api/character/list', {
      params: {
        page: params.page || 1,
        page_size: params.pageSize || 20,
        category_id: params.categoryId,
        sort_by: params.sortBy || 'created_at',
        order: params.order || 'desc'
      }
    })
  },

  // 搜索角色
  searchCharacters(params) {
    return characterApi.get('/api/character/search', {
      params: {
        keyword: params.keyword,
        category_id: params.categoryId,
        page: params.page || 1,
        page_size: params.pageSize || 20
      }
    })
  },

  // 获取角色详情
  getCharacterDetail(id) {
    return characterApi.get(`/api/character/${id}`)
  },

  // 获取推荐角色
  getRecommendedCharacters(params = {}) {
    return characterApi.get('/api/character/recommended', {
      params: {
        count: params.count || 10,
        exclude_ids: params.excludeIds
      }
    })
  },

  // 获取角色分类
  getCharacterCategories() {
    return characterApi.get('/api/character/categories')
  },

  // 获取热门角色
  getPopularCharacters(params = {}) {
    return characterApi.get('/api/character/popular', {
      params: {
        page: params.page || 1,
        page_size: params.pageSize || 20
      }
    })
  },

  // 获取角色标签
  getCharacterTags() {
    return characterApi.get('/api/character/tags')
  },

  // 创建角色（需要认证）
  createCharacter(data) {
    return characterApi.post('/api/character', {
      name: data.name,
      avatar: data.avatar,
      description: data.description,
      short_desc: data.shortDesc,
      category_id: data.categoryId,
      tags: data.tags,
      prompt: data.prompt,
      personality: data.personality,
      voice_settings: data.voiceSettings,
      is_public: data.isPublic
    })
  },

  // 更新角色（需要认证）
  updateCharacter(id, data) {
    return characterApi.put(`/api/character/${id}`, {
      name: data.name,
      avatar: data.avatar,
      description: data.description,
      short_desc: data.shortDesc,
      category_id: data.categoryId,
      tags: data.tags,
      prompt: data.prompt,
      personality: data.personality,
      voice_settings: data.voiceSettings,
      status: data.status,
      is_public: data.isPublic
    })
  },

  // 删除角色（需要认证）
  deleteCharacter(id) {
    return characterApi.delete(`/api/character/${id}`)
  },

  // 收藏/取消收藏角色（需要认证）
  toggleFavorite(id) {
    return characterApi.post(`/api/character/${id}/favorite`)
  },

  // 获取我的收藏角色（需要认证）
  getMyFavorites(params = {}) {
    return characterApi.get('/api/character/favorites', {
      params: {
        page: params.page || 1,
        page_size: params.pageSize || 20
      }
    })
  },

  // 获取我创建的角色（需要认证）
  getMyCharacters(params = {}) {
    return characterApi.get('/api/character/my', {
      params: {
        page: params.page || 1,
        page_size: params.pageSize || 20,
        status: params.status,
        is_public: params.isPublic
      }
    })
  },

  // 更新角色提示词（需要认证）
  updatePrompt(id, prompt) {
    return characterApi.put(`/api/character/${id}/prompt`, {
      prompt
    })
  },

  // 更新角色性格设置（需要认证）
  updatePersonality(id, personality) {
    return characterApi.put(`/api/character/${id}/personality`, {
      personality
    })
  },

  // 更新语音设置（需要认证）
  updateVoiceSettings(id, voiceSettings) {
    return characterApi.put(`/api/character/${id}/voice`, {
      voice_settings: voiceSettings
    })
  }
}

// 保持向后兼容
export const characterApi_old = characterApiService
export default characterApiService

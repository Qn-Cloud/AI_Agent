import { defineStore } from 'pinia'
import { characterApiService as characterApi } from '../services'

export const useCharacterStore = defineStore('character', {
  state: () => ({
    // 角色列表
    characters: [
      {
        id: 'harry-potter',
        name: '哈利·波特',
        avatar: '/images/avatars/harry-potter.jpg',
        description: '霍格沃茨魔法学校的学生，拥有闪电疤痕的男孩巫师。勇敢、善良，擅长魁地奇运动。',
        tags: ['魔法', '勇敢', '冒险', '友谊'],
        status: 'online',
        personality: {
          friendliness: 85,
          humor: 70,
          intelligence: 80,
          creativity: 75
        },
        prompt: '你是哈利·波特，霍格沃茨的学生。你勇敢善良，有着丰富的魔法世界冒险经历。请用哈利的语气和视角来回答问题。',
        voiceSettings: {
          rate: 1.0,
          pitch: 1.0,
          volume: 0.8
        }
      },
      {
        id: 'socrates',
        name: '苏格拉底',
        avatar: '/images/avatars/socrates.jpg',
        description: '古希腊哲学家，以苏格拉底式问答法闻名。追求智慧与真理，善于启发式教学。',
        tags: ['哲学', '智慧', '思辨', '教育'],
        status: 'online',
        personality: {
          friendliness: 75,
          humor: 60,
          intelligence: 95,
          creativity: 85
        },
        prompt: '你是苏格拉底，古希腊的哲学家。你善于通过提问来启发他人思考，追求智慧和真理。请用苏格拉底的方式来对话。',
        voiceSettings: {
          rate: 0.9,
          pitch: 0.8,
          volume: 0.7
        }
      },
      {
        id: 'shakespeare',
        name: '莎士比亚',
        avatar: '/images/avatars/shakespeare.jpg',
        description: '英国文艺复兴时期的伟大剧作家和诗人，创作了众多不朽的戏剧和十四行诗。',
        tags: ['文学', '戏剧', '诗歌', '创作'],
        status: 'online',
        personality: {
          friendliness: 80,
          humor: 90,
          intelligence: 92,
          creativity: 98
        },
        prompt: '你是威廉·莎士比亚，伟大的剧作家和诗人。你富有创造力，语言优美，善于用戏剧性的方式表达。',
        voiceSettings: {
          rate: 1.1,
          pitch: 1.2,
          volume: 0.9
        }
      },
      {
        id: 'einstein',
        name: '爱因斯坦',
        avatar: '/images/avatars/einstein.jpg',
        description: '20世纪最伟大的物理学家之一，相对论的提出者，诺贝尔物理学奖获得者。',
        tags: ['科学', '物理', '相对论', '思考'],
        status: 'online',
        personality: {
          friendliness: 85,
          humor: 75,
          intelligence: 98,
          creativity: 90
        },
        prompt: '你是阿尔伯特·爱因斯坦，著名的物理学家。你善于用简单的方式解释复杂的科学概念，充满好奇心和想象力。',
        voiceSettings: {
          rate: 0.95,
          pitch: 0.9,
          volume: 0.8
        }
      },
      {
        id: 'sherlock',
        name: '夏洛克·福尔摩斯',
        avatar: '/images/avatars/sherlock.jpg',
        description: '世界著名的咨询侦探，居住在贝克街221B号，擅长演绎推理和观察细节。',
        tags: ['推理', '侦探', '观察', '逻辑'],
        status: 'online',
        personality: {
          friendliness: 60,
          humor: 65,
          intelligence: 96,
          creativity: 85
        },
        prompt: '你是夏洛克·福尔摩斯，世界上最优秀的咨询侦探。你善于观察细节，进行逻辑推理，有时显得冷漠但内心正义。',
        voiceSettings: {
          rate: 1.2,
          pitch: 1.0,
          volume: 0.85
        }
      },
      {
        id: 'hermione',
        name: '赫敏·格兰杰',
        avatar: '/images/avatars/hermione.jpg',
        description: '霍格沃茨最聪明的学生之一，博学多才，热爱读书，是哈利和罗恩的好友。',
        tags: ['魔法', '学霸', '聪明', '正义'],
        status: 'online',
        personality: {
          friendliness: 80,
          humor: 55,
          intelligence: 95,
          creativity: 75
        },
        prompt: '你是赫敏·格兰杰，霍格沃茨的优秀学生。你博学多才，逻辑清晰，总是能找到解决问题的方法。',
        voiceSettings: {
          rate: 1.1,
          pitch: 1.3,
          volume: 0.8
        }
      }
    ],
    currentCharacter: null,
    favorites: [],
    searchQuery: '',
    selectedTags: [],
    loading: false,
    error: null,
    categories: [],
    tags: []
  }),

  getters: {
    // 过滤后的角色列表
    filteredCharacters: (state) => {
      let filtered = state.characters

      // 按搜索关键词过滤
      if (state.searchQuery) {
        const query = state.searchQuery.toLowerCase()
        filtered = filtered.filter(character => 
          character.name.toLowerCase().includes(query) ||
          character.description.toLowerCase().includes(query) ||
          character.tags.some(tag => tag.toLowerCase().includes(query))
        )
      }

      // 按标签过滤
      if (state.selectedTags.length > 0) {
        filtered = filtered.filter(character =>
          state.selectedTags.some(tag => character.tags.includes(tag))
        )
      }

      return filtered
    },

    // 获取所有标签
    allTags: (state) => {
      const tags = new Set()
      state.characters.forEach(character => {
        character.tags.forEach(tag => tags.add(tag))
      })
      return Array.from(tags)
    },

    // 收藏的角色
    favoriteCharacters: (state) => {
      return state.characters.filter(character => 
        state.favorites.includes(character.id)
      )
    }
  },

  actions: {
    // 从API加载角色列表
    async loadCharacters(params = {}) {
      try {
        this.loading = true
        this.error = null
        const response = await characterApi.getCharacterList(params)
        if (response.data && response.data.characters) {
          this.characters = response.data.characters
        }
      } catch (error) {
        this.error = error.message
        console.error('加载角色列表失败:', error)
      } finally {
        this.loading = false
      }
    },

    // 搜索角色
    async searchCharacters(keyword) {
      try {
        this.loading = true
        this.error = null
        const response = await characterApi.searchCharacters({ keyword })
        if (response.data && response.data.characters) {
          this.characters = response.data.characters
        }
      } catch (error) {
        this.error = error.message
        console.error('搜索角色失败:', error)
      } finally {
        this.loading = false
      }
    },

    // 获取推荐角色
    async loadRecommendedCharacters() {
      try {
        this.loading = true
        const response = await characterApi.getRecommendedCharacters()
        if (response.data && response.data.characters) {
          this.characters = response.data.characters
        }
      } catch (error) {
        this.error = error.message
        console.error('加载推荐角色失败:', error)
      } finally {
        this.loading = false
      }
    },

    // 获取角色分类
    async loadCategories() {
      try {
        const response = await characterApi.getCharacterCategories()
        if (response.data) {
          this.categories = response.data.categories || []
        }
      } catch (error) {
        console.error('加载角色分类失败:', error)
      }
    },

    // 获取角色标签
    async loadTags() {
      try {
        const response = await characterApi.getCharacterTags()
        if (response.data) {
          this.tags = response.data.tags || []
        }
      } catch (error) {
        console.error('加载角色标签失败:', error)
      }
    },

    // 选择角色
    async selectCharacter(characterId) {
      // 先从本地查找
      let character = this.characters.find(c => c.id === characterId)
      
      // 如果本地没有，从API获取详情
      if (!character) {
        try {
          const response = await characterApi.getCharacterDetail(characterId)
          if (response.data) {
            character = response.data.character
            // 添加到本地列表
            this.characters.push(character)
          }
        } catch (error) {
          console.error('获取角色详情失败:', error)
          return
        }
      }
      
      if (character) {
        this.currentCharacter = character
      }
    },

    // 设置搜索关键词
    setSearchQuery(query) {
      this.searchQuery = query
    },

    // 设置选中的标签
    setSelectedTags(tags) {
      this.selectedTags = tags
    },

    // 添加/移除收藏
    async toggleFavorite(characterId) {
      try {
        await characterApi.toggleFavorite(characterId)
        
        const index = this.favorites.indexOf(characterId)
        if (index > -1) {
          this.favorites.splice(index, 1)
        } else {
          this.favorites.push(characterId)
        }
      } catch (error) {
        console.error('收藏操作失败:', error)
        throw error
      }
    },

    // 加载我的收藏
    async loadMyFavorites() {
      try {
        const response = await characterApi.getMyFavorites()
        if (response.data && response.data.characters) {
          this.favorites = response.data.characters.map(c => c.id)
        }
      } catch (error) {
        console.error('加载收藏列表失败:', error)
      }
    },

    // 更新角色设定
    async updateCharacterPrompt(characterId, prompt) {
      try {
        await characterApi.updatePrompt(characterId, prompt)
        
        const character = this.characters.find(c => c.id === characterId)
        if (character) {
          character.prompt = prompt
        }
      } catch (error) {
        console.error('更新角色设定失败:', error)
        throw error
      }
    },

    // 更新角色性格
    async updateCharacterPersonality(characterId, personality) {
      try {
        await characterApi.updatePersonality(characterId, personality)
        
        const character = this.characters.find(c => c.id === characterId)
        if (character) {
          character.personality = { ...character.personality, ...personality }
        }
      } catch (error) {
        console.error('更新角色性格失败:', error)
        throw error
      }
    },

    // 更新语音设置
    async updateVoiceSettings(characterId, settings) {
      try {
        await characterApi.updateVoiceSettings(characterId, settings)
        
        const character = this.characters.find(c => c.id === characterId)
        if (character) {
          character.voiceSettings = { ...character.voiceSettings, ...settings }
        }
      } catch (error) {
        console.error('更新语音设置失败:', error)
        throw error
      }
    },

    // 创建自定义角色
    async createCharacter(characterData) {
      try {
        const response = await characterApi.createCharacter(characterData)
        if (response.data && response.data.character) {
          this.characters.push(response.data.character)
          return response.data.character
        }
      } catch (error) {
        console.error('创建角色失败:', error)
        throw error
      }
    }
  }
})

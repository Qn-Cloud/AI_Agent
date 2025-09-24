<template>
  <div class="home-page">
    <!-- 搜索框 -->
    <div class="search-section">
      <el-input
        v-model="searchQuery"
        placeholder="搜索你想要聊天的角色..."
        size="large"
        clearable
        @input="handleSearch"
        class="search-input"
      >
        <template #suffix>
          <el-icon class="search-icon"><Search /></el-icon>
        </template>
      </el-input>
    </div>

    <!-- 为您推荐 -->
    <div class="section">
      <h2 class="section-title">为您推荐</h2>
      <div class="character-grid">
        <div
          v-for="character in recommendedCharacters"
          :key="character.id"
          class="character-card"
          @click="startChat(character.id)"
        >
          <div class="character-avatar">
            <img 
              :src="character.avatar" 
              :alt="character.name"
              class="avatar-image"
              @error="handleImageError"
            />
          </div>
          <div class="character-info">
            <h3 class="character-name">{{ character.name }}</h3>
            <p class="character-desc">{{ character.shortDesc }}</p>
            <div class="character-stats">
              <span class="stat-item">1.1k <el-icon><Star /></el-icon></span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 历史人物 -->
    <div class="section">
      <h2 class="section-title">历史人物</h2>
      <div class="character-grid">
        <div
          v-for="character in historicalCharacters"
          :key="character.id"
          class="character-card"
          @click="startChat(character.id)"
        >
          <div class="character-avatar">
            <img 
              :src="character.avatar" 
              :alt="character.name"
              class="avatar-image"
              @error="handleImageError"
            />
          </div>
          <div class="character-info">
            <h3 class="character-name">{{ character.name }}</h3>
            <p class="character-desc">{{ character.shortDesc }}</p>
            <div class="character-stats">
              <span class="stat-item">1.1k <el-icon><Star /></el-icon></span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 经典IP -->
    <div class="section">
      <h2 class="section-title">经典IP</h2>
      <div class="character-grid">
        <div
          v-for="character in classicCharacters"
          :key="character.id"
          class="character-card"
          @click="startChat(character.id)"
        >
          <div class="character-avatar">
            <img 
              :src="character.avatar" 
              :alt="character.name"
              class="avatar-image"
              @error="handleImageError"
            />
          </div>
          <div class="character-info">
            <h3 class="character-name">{{ character.name }}</h3>
            <p class="character-desc">{{ character.shortDesc }}</p>
            <div class="character-stats">
              <span class="stat-item">1.1k <el-icon><Star /></el-icon></span>
            </div>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useCharacterStore } from '../stores/character'
import { useChatStore } from '../stores/chat'
import { Search, User, Star } from '@element-plus/icons-vue'

const router = useRouter()
const characterStore = useCharacterStore()
const chatStore = useChatStore()

// 响应式数据
const searchQuery = ref('')

// 计算属性
const allCharacters = computed(() => characterStore.characters)

// 为每个角色添加短描述
const charactersWithShortDesc = computed(() => {
  return allCharacters.value.map(char => ({
    ...char,
    shortDesc: char.description.slice(0, 50) + (char.description.length > 50 ? '...' : '')
  }))
})

// 推荐角色（前4个）
const recommendedCharacters = computed(() => {
  return charactersWithShortDesc.value.slice(0, 4)
})

// 历史人物（爱因斯坦、苏格拉底、莎士比亚）
const historicalCharacters = computed(() => {
  return charactersWithShortDesc.value.filter(char => 
    ['einstein', 'socrates', 'shakespeare'].includes(char.id)
  )
})

// 经典IP角色（哈利波特、福尔摩斯、赫敏）
const classicCharacters = computed(() => {
  return charactersWithShortDesc.value.filter(char => 
    ['harry-potter', 'sherlock', 'hermione'].includes(char.id)
  )
})

// 方法
const handleSearch = (value) => {
  characterStore.setSearchQuery(value)
}

const startChat = (characterId) => {
  characterStore.selectCharacter(characterId)
  chatStore.startNewConversation(characterId)
  router.push(`/chat/${characterId}`)
}

const handleImageError = (event) => {
  // 图片加载失败时显示默认头像
  event.target.style.display = 'none'
  const placeholder = document.createElement('div')
  placeholder.className = 'avatar-placeholder'
  placeholder.innerHTML = '<el-icon><User /></el-icon>'
  event.target.parentNode.appendChild(placeholder)
}
</script>

<style lang="scss" scoped>
.home-page {
  padding: 40px 60px;
  background: #f5f5f5;
  min-height: 100vh;
}

.search-section {
  display: flex;
  justify-content: center;
  margin-bottom: 48px;
  
  .search-input {
    max-width: 600px;
    width: 100%;
    
    .el-input__wrapper {
      border-radius: 24px;
      height: 48px;
      background: white;
      border: 1px solid #e5e5e5;
      box-shadow: none;
      
      &:hover {
        border-color: #d0d0d0;
        box-shadow: none;
      }
      
      &.is-focus {
        border-color: #4A90E2;
        box-shadow: none;
      }
    }
    
    .el-input__inner {
      font-size: 16px;
      padding-left: 20px;
      color: #333;
      
      &::placeholder {
        color: #999;
      }
    }
    
    .search-icon {
      color: #999;
      font-size: 18px;
      margin-right: 16px;
    }
  }
}

.section {
  margin-bottom: 48px;
  
  .section-title {
    font-size: 20px;
    font-weight: 600;
    color: #333;
    margin-bottom: 24px;
  }
}

.character-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 20px;
}

.character-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.2s ease;
  border: 1px solid #f0f0f0;
  
  &:hover {
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
    transform: translateY(-2px);
  }
  
  .character-avatar {
    display: flex;
    justify-content: center;
    margin-bottom: 16px;
    
    .avatar-image {
      width: 80px;
      height: 80px;
      border-radius: 12px;
      object-fit: cover;
      object-position: center;
      transition: transform 0.3s ease;
      
      &:hover {
        transform: scale(1.05);
      }
    }
    
    .avatar-placeholder {
      width: 80px;
      height: 80px;
      background: #4A90E2;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;
      font-size: 32px;
    }
  }
  
  .character-info {
    text-align: center;
    
    .character-name {
      font-size: 16px;
      font-weight: 600;
      color: #4A90E2;
      margin: 0 0 8px 0;
    }
    
    .character-desc {
      font-size: 14px;
      color: #666;
      line-height: 1.4;
      margin: 0 0 12px 0;
      min-height: 40px;
    }
    
    .character-stats {
      .stat-item {
        display: inline-flex;
        align-items: center;
        gap: 4px;
        font-size: 13px;
        color: #999;
        
        .el-icon {
          font-size: 14px;
        }
      }
    }
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .home-page {
    padding: 20px 20px;
  }
  
  .search-section {
    margin-bottom: 32px;
    
    .search-input {
      .el-input__wrapper {
        height: 44px;
      }
      
      .el-input__inner {
        font-size: 15px;
        padding-left: 16px;
      }
    }
  }
  
  .section {
    margin-bottom: 36px;
    
    .section-title {
      font-size: 18px;
      margin-bottom: 20px;
    }
  }
  
  .character-grid {
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 16px;
  }
  
  .character-card {
    padding: 16px;
    
    .character-avatar {
      margin-bottom: 12px;
      
      .avatar-image {
        width: 64px;
        height: 64px;
      }
      
      .avatar-placeholder {
        width: 64px;
        height: 64px;
        font-size: 24px;
      }
    }
    
    .character-info {
      .character-name {
        font-size: 15px;
      }
      
      .character-desc {
        font-size: 13px;
        min-height: 32px;
      }
    }
  }
}
</style>

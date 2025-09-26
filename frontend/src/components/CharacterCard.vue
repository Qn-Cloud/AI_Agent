<template>
  <div class="character-card" @click="selectCharacter">
    <div class="card-header">
      <div class="avatar-section">
        <img :src="character.avatar" :alt="character.name" class="avatar" />
        <div class="status-indicator" :class="character.status"></div>
      </div>
      <div class="favorite-btn" @click.stop="toggleFavorite">
        <el-icon :class="{ active: isFavorite }">
          <StarFilled v-if="isFavorite" />
          <Star v-else />
        </el-icon>
      </div>
    </div>

    <div class="card-content">
      <h3 class="character-name">{{ character.name }}</h3>
      <p class="character-description">{{ character.description }}</p>
      
      <div class="tags-section">
        <el-tag
          v-for="tag in character.tags"
          :key="tag"
          size="small"
          class="tag-item"
          :type="getTagType(tag)"
        >
          {{ tag }}
        </el-tag>
      </div>

      <div class="personality-bars">
        <div class="personality-item">
          <span class="label">友善度</span>
          <el-progress
            :percentage="character.personality.friendliness"
            :show-text="false"
            :stroke-width="4"
            color="#67C23A"
          />
        </div>
        <div class="personality-item">
          <span class="label">智慧</span>
          <el-progress
            :percentage="character.personality.intelligence"
            :show-text="false"
            :stroke-width="4"
            color="#409EFF"
          />
        </div>
      </div>
    </div>

    <div class="card-actions">
      <el-button
        type="primary"
        size="default"
        @click.stop="startChat"
        class="chat-btn"
      >
        <el-icon><ChatDotRound /></el-icon>
        开始对话
      </el-button>
      <el-button
        size="default"
        @click.stop="openSettings"
        class="settings-btn"
      >
        <el-icon><Setting /></el-icon>
        设置
      </el-button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useCharacterStore } from '../stores/character'
import { useChatStore } from '../stores/chat'
import { ChatDotRound, Setting } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const props = defineProps({
  character: {
    type: Object,
    required: true
  }
})

const router = useRouter()
const characterStore = useCharacterStore()
const chatStore = useChatStore()

const isFavorite = computed(() => {
  return characterStore.favorites.includes(props.character.id)
})

const selectCharacter = () => {
  characterStore.selectCharacter(props.character.id)
}

const toggleFavorite = async () => {
  try {
    await characterStore.toggleFavorite(props.character.id)
  } catch (error) {
    ElMessage.error('收藏操作失败: ' + error.message)
  }
}

const startChat = async () => {
  try {
    await characterStore.selectCharacter(props.character.id)
    await chatStore.startNewConversation(props.character.id)
    router.push(`/chat/${props.character.id}`)
  } catch (error) {
    ElMessage.error('启动对话失败: ' + error.message)
  }
}

const openSettings = async () => {
  try {
    await characterStore.selectCharacter(props.character.id)
    router.push('/settings')
  } catch (error) {
    ElMessage.error('打开设置失败: ' + error.message)
  }
}

const getTagType = (tag) => {
  const tagTypes = {
    '魔法': '',
    '哲学': 'success',
    '科学': 'info',
    '文学': 'warning',
    '推理': 'danger',
    '勇敢': 'success',
    '智慧': 'info',
    '创作': 'warning'
  }
  return tagTypes[tag] || ''
}
</script>

<style lang="scss" scoped>
.character-card {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  height: 100%;
  display: flex;
  flex-direction: column;

  &:hover {
    transform: translateY(-8px);
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
    background: rgba(255, 255, 255, 0.98);
  }
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
}

.avatar-section {
  position: relative;
  
  .avatar {
    width: 60px;
    height: 60px;
    border-radius: 50%;
    object-fit: cover;
    border: 3px solid rgba(255, 255, 255, 0.5);
  }
  
  .status-indicator {
    position: absolute;
    bottom: 2px;
    right: 2px;
    width: 16px;
    height: 16px;
    border-radius: 50%;
    border: 2px solid white;
    
    &.online {
      background: #67C23A;
    }
    
    &.offline {
      background: #909399;
    }
    
    &.busy {
      background: #F56C6C;
    }
  }
}

.favorite-btn {
  padding: 4px;
  border-radius: 50%;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &:hover {
    background: rgba(0, 0, 0, 0.05);
  }
  
  .el-icon {
    font-size: 20px;
    color: #DCDFE6;
    transition: color 0.2s ease;
    
    &.active {
      color: #F7BA2A;
    }
  }
}

.card-content {
  flex: 1;
  margin-bottom: 16px;
}

.character-name {
  font-size: 18px;
  font-weight: bold;
  margin: 0 0 8px 0;
  color: #303133;
}

.character-description {
  font-size: 14px;
  color: #606266;
  line-height: 1.5;
  margin: 0 0 12px 0;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.tags-section {
  margin-bottom: 16px;
  
  .tag-item {
    margin-right: 6px;
    margin-bottom: 6px;
    border-radius: 12px;
  }
}

.personality-bars {
  .personality-item {
    display: flex;
    align-items: center;
    margin-bottom: 8px;
    
    .label {
      font-size: 12px;
      color: #909399;
      width: 50px;
      margin-right: 12px;
    }
    
    .el-progress {
      flex: 1;
    }
  }
}

.card-actions {
  display: flex;
  gap: 8px;
  
  .chat-btn {
    flex: 1;
    border-radius: 8px;
  }
  
  .settings-btn {
    border-radius: 8px;
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .character-card {
    padding: 16px;
  }
  
  .avatar-section .avatar {
    width: 50px;
    height: 50px;
  }
  
  .character-name {
    font-size: 16px;
  }
  
  .character-description {
    font-size: 13px;
    -webkit-line-clamp: 2;
  }
  
  .card-actions {
    flex-direction: column;
    
    .chat-btn,
    .settings-btn {
      width: 100%;
    }
  }
}
</style>

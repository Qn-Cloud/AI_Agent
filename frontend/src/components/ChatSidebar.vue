<template>
  <div class="chat-sidebar">
    <div class="sidebar-header">
      <h2>对话历史</h2>
      <el-button @click="refreshHistory" :icon="Refresh" size="small" text>
        刷新
      </el-button>
    </div>

    <div class="sidebar-content">
      <!-- 加载状态 -->
      <div v-if="chatStore.isLoading" class="loading-state">
        <el-skeleton :rows="5" animated />
      </div>

      <!-- 历史记录列表 -->
      <div v-else class="history-groups">
        <!-- 今天 -->
        <div v-if="chatStore.groupedHistory.todays.length > 0" class="history-group">
          <div class="group-header">
            <h3>今天</h3>
            <span class="group-count">{{ chatStore.groupedHistory.todays.length }}</span>
          </div>
          <div class="group-items">
            <div
              v-for="item in chatStore.groupedHistory.todays"
              :key="item.conversation_id"
              class="history-item"
              @click="selectConversation(item)"
            >
              <div class="item-avatar">
                <img
                  :src="getCharacterAvatar(item.character_id)"
                  :alt="item.character_name"
                  class="character-avatar"
                  @error="handleImageError"
                />
              </div>
              <div class="item-content">
                <div class="item-title">{{ item.character_name }}</div>
                <div class="item-time">{{ formatTime(item.created_at) }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- 昨天 -->
        <div v-if="chatStore.groupedHistory.yesterdays.length > 0" class="history-group">
          <div class="group-header">
            <h3>昨天</h3>
            <span class="group-count">{{ chatStore.groupedHistory.yesterdays.length }}</span>
          </div>
          <div class="group-items">
            <div
              v-for="item in chatStore.groupedHistory.yesterdays"
              :key="item.conversation_id"
              class="history-item"
              @click="selectConversation(item)"
            >
              <div class="item-avatar">
                <img
                  :src="getCharacterAvatar(item.character_id)"
                  :alt="item.character_name"
                  class="character-avatar"
                  @error="handleImageError"
                />
              </div>
              <div class="item-content">
                <div class="item-title">{{ item.character_name }}</div>
                <div class="item-time">{{ formatTime(item.created_at) }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- 更久之前 -->
        <div v-if="chatStore.groupedHistory.befores.length > 0" class="history-group">
          <div class="group-header">
            <h3>更久之前</h3>
            <span class="group-count">{{ chatStore.groupedHistory.befores.length }}</span>
          </div>
          <div class="group-items">
            <div
              v-for="item in chatStore.groupedHistory.befores"
              :key="item.conversation_id"
              class="history-item"
              @click="selectConversation(item)"
            >
              <div class="item-avatar">
                <img
                  :src="getCharacterAvatar(item.character_id)"
                  :alt="item.character_name"
                  class="character-avatar"
                  @error="handleImageError"
                />
              </div>
              <div class="item-content">
                <div class="item-title">{{ item.character_name }}</div>
                <div class="item-time">{{ formatTime(item.created_at) }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- 空状态 -->
        <div v-if="!hasAnyHistory" class="empty-state">
          <el-empty description="暂无对话历史">
            <el-button type="primary" @click="goToDiscover">
              开始新对话
            </el-button>
          </el-empty>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useChatStore } from '../stores/chat'
import { useCharacterStore } from '../stores/character'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'

const router = useRouter()
const chatStore = useChatStore()
const characterStore = useCharacterStore()

// 角色映射表（与History.vue保持一致）
const characterMap = {
  1: { id: 1, name: '哈利·波特', avatar: 'https://example.com/harry.jpg' },
  2: { id: 2, name: '苏格拉底', avatar: 'https://example.com/socrates.jpg' },
  4: { id: 4, name: '爱因斯坦', avatar: 'https://example.com/einstein.jpg' },
  6: { id: 6, name: '赫敏·格兰杰', avatar: 'https://example.com/hermione.jpg' }
}

// 计算属性
const hasAnyHistory = computed(() => {
  return chatStore.groupedHistory.todays.length > 0 ||
         chatStore.groupedHistory.yesterdays.length > 0 ||
         chatStore.groupedHistory.befores.length > 0
})

// 方法
const refreshHistory = async () => {
  try {
    await chatStore.loadGroupedHistory()
    ElMessage.success('历史记录已刷新')
  } catch (error) {
    ElMessage.error('刷新失败: ' + error.message)
  }
}

const selectConversation = (item) => {
  console.log('选择对话:', item)
  // 跳转到聊天页面并选择对应的角色和对话
  router.push({
    name: 'Chat',
    params: {
      characterId: item.character_id
    },
    query: {
      conversationId: item.conversation_id
    }
  })
}

const getCharacterAvatar = (characterId) => {
  return characterMap[characterId]?.avatar || 'https://example.com/default.jpg'
}

const handleImageError = (event) => {
  event.target.src = 'https://example.com/default.jpg'
}

const formatTime = (timeString) => {
  try {
    const date = new Date(timeString)
    const now = new Date()
    const diff = now - date
    
    // 小于1小时显示分钟
    if (diff < 60 * 60 * 1000) {
      const minutes = Math.floor(diff / (60 * 1000))
      return minutes <= 0 ? '刚刚' : `${minutes}分钟前`
    }
    
    // 小于24小时显示小时
    if (diff < 24 * 60 * 60 * 1000) {
      const hours = Math.floor(diff / (60 * 60 * 1000))
      return `${hours}小时前`
    }
    
    // 否则显示具体时间
    return date.toLocaleString('zh-CN', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch (error) {
    return timeString
  }
}

const goToDiscover = () => {
  router.push('/discover')
}

// 生命周期
onMounted(() => {
  chatStore.loadGroupedHistory()
})
</script>

<style lang="scss" scoped>
.chat-sidebar {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: transparent;
  border: none;

  .sidebar-header {
    padding: 16px 20px 12px;
    border-bottom: 1px solid #f0f0f0;
    display: flex;
    align-items: center;
    justify-content: space-between;

    h2 {
      margin: 0;
      font-size: 16px;
      font-weight: 600;
      color: #2c3e50;
    }
  }

  .sidebar-content {
    flex: 1;
    overflow-y: auto;
    padding: 12px 0;

    .loading-state {
      padding: 0 20px;
    }

    .history-groups {
      .history-group {
        margin-bottom: 20px;

        .group-header {
          padding: 0 20px 8px;
          display: flex;
          align-items: center;
          justify-content: space-between;

          h3 {
            margin: 0;
            font-size: 13px;
            font-weight: 600;
            color: #6c757d;
            text-transform: uppercase;
            letter-spacing: 0.5px;
          }

          .group-count {
            font-size: 11px;
            color: #adb5bd;
            background: #f8f9fa;
            padding: 2px 6px;
            border-radius: 8px;
          }
        }

        .group-items {
          .history-item {
            padding: 10px 20px;
            display: flex;
            align-items: center;
            cursor: pointer;
            transition: all 0.2s ease;
            border-radius: 0;

            &:hover {
              background: #f8f9fa;
            }

            .item-avatar {
              width: 32px;
              height: 32px;
              margin-right: 12px;
              flex-shrink: 0;

              .character-avatar {
                width: 100%;
                height: 100%;
                border-radius: 50%;
                object-fit: cover;
              }
            }

            .item-content {
              flex: 1;
              min-width: 0;

              .item-title {
                font-size: 13px;
                font-weight: 500;
                color: #2c3e50;
                margin-bottom: 2px;
                white-space: nowrap;
                overflow: hidden;
                text-overflow: ellipsis;
              }

              .item-time {
                font-size: 11px;
                color: #6c757d;
              }
            }
          }
        }
      }

      .empty-state {
        padding: 30px 20px;
        text-align: center;
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .chat-sidebar {
    width: 100%;
    height: auto;
    border-right: none;
    border-bottom: 1px solid #e9ecef;
  }
}
</style> 
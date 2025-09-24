<template>
  <div class="history-page">
    <div class="history-container">
      <!-- 页面标题 -->
      <div class="page-header">
        <h1>对话历史</h1>
        <p>查看和管理您的所有对话记录</p>
      </div>

      <!-- 搜索和筛选 -->
      <div class="search-section">
        <div class="search-controls">
          <el-input
            v-model="searchQuery"
            placeholder="搜索对话标题或内容..."
            clearable
            @input="handleSearch"
            class="search-input"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>

          <el-select
            v-model="selectedCharacter"
            placeholder="选择角色"
            clearable
            @change="handleCharacterFilter"
            class="character-filter"
          >
            <el-option label="所有角色" value="" />
            <el-option
              v-for="character in characters"
              :key="character.id"
              :label="character.name"
              :value="character.id"
            />
          </el-select>

          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            @change="handleDateFilter"
            class="date-filter"
          />
        </div>

        <div class="action-controls">
          <el-button @click="exportAll" :icon="Download">
            导出全部
          </el-button>
          <el-button @click="deleteSelected" :icon="Delete" type="danger" :disabled="selectedConversations.length === 0">
            删除选中 ({{ selectedConversations.length }})
          </el-button>
        </div>
      </div>

      <!-- 统计信息 -->
      <div class="stats-section">
        <div class="stat-card">
          <div class="stat-icon">
            <el-icon><ChatDotRound /></el-icon>
          </div>
          <div class="stat-content">
            <span class="stat-number">{{ filteredConversations.length }}</span>
            <span class="stat-label">对话总数</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">
            <el-icon><Timer /></el-icon>
          </div>
          <div class="stat-content">
            <span class="stat-number">{{ totalMessages }}</span>
            <span class="stat-label">消息总数</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">
            <el-icon><User /></el-icon>
          </div>
          <div class="stat-content">
            <span class="stat-number">{{ uniqueCharacters }}</span>
            <span class="stat-label">互动角色</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">
            <el-icon><Calendar /></el-icon>
          </div>
          <div class="stat-content">
            <span class="stat-number">{{ activeDays }}</span>
            <span class="stat-label">活跃天数</span>
          </div>
        </div>
      </div>

      <!-- 排序和视图选项 -->
      <div class="control-section">
        <div class="sort-controls">
          <span class="control-label">排序：</span>
          <el-radio-group v-model="sortBy" @change="handleSort">
            <el-radio-button label="time">时间</el-radio-button>
            <el-radio-button label="messages">消息数</el-radio-button>
            <el-radio-button label="character">角色</el-radio-button>
          </el-radio-group>
        </div>

        <div class="view-controls">
          <span class="control-label">视图：</span>
          <el-radio-group v-model="viewMode" @change="handleViewChange">
            <el-radio-button label="list">列表</el-radio-button>
            <el-radio-button label="grid">网格</el-radio-button>
            <el-radio-button label="timeline">时间线</el-radio-button>
          </el-radio-group>
        </div>

        <div class="select-controls">
          <el-checkbox
            v-model="selectAll"
            @change="handleSelectAll"
            :indeterminate="isIndeterminate"
          >
            全选
          </el-checkbox>
        </div>
      </div>

      <!-- 对话列表 -->
      <div class="conversations-section">
        <!-- 空状态 -->
        <div v-if="filteredConversations.length === 0" class="empty-state">
          <el-empty description="没有找到对话记录">
            <el-button type="primary" @click="goToHome">开始新对话</el-button>
          </el-empty>
        </div>

        <!-- 列表视图 -->
        <div v-else-if="viewMode === 'list'" class="list-view">
          <div
            v-for="conversation in paginatedConversations"
            :key="conversation.id"
            class="conversation-item"
            :class="{ selected: selectedConversations.includes(conversation.id) }"
          >
            <div class="conversation-select">
              <el-checkbox
                :model-value="selectedConversations.includes(conversation.id)"
                @change="toggleConversationSelect(conversation.id)"
              />
            </div>

            <div class="conversation-avatar">
              <img
                :src="getCharacterById(conversation.characterId)?.avatar"
                :alt="getCharacterById(conversation.characterId)?.name"
                class="character-avatar"
              />
            </div>

            <div class="conversation-content" @click="openConversation(conversation)">
              <div class="conversation-header">
                <h3 class="conversation-title">{{ conversation.title }}</h3>
                <span class="conversation-time">{{ formatTime(conversation.lastUpdate) }}</span>
              </div>
              <div class="conversation-meta">
                <span class="character-name">与{{ getCharacterById(conversation.characterId)?.name }}的对话</span>
                <span class="message-count">{{ conversation.messages.length }} 条消息</span>
                <span class="duration">{{ formatDuration(conversation) }}</span>
              </div>
              <div class="conversation-preview">
                {{ getLastMessage(conversation) }}
              </div>
            </div>

            <div class="conversation-actions">
              <el-button @click.stop="continueConversation(conversation)" :icon="ChatDotRound" size="small">
                继续
              </el-button>
              <el-button @click.stop="exportConversation(conversation)" :icon="Download" size="small">
                导出
              </el-button>
              <el-button @click.stop="deleteConversation(conversation.id)" :icon="Delete" size="small" type="danger">
                删除
              </el-button>
            </div>
          </div>
        </div>

        <!-- 网格视图 -->
        <div v-else-if="viewMode === 'grid'" class="grid-view">
          <div class="grid-container">
            <div
              v-for="conversation in paginatedConversations"
              :key="conversation.id"
              class="conversation-card"
              :class="{ selected: selectedConversations.includes(conversation.id) }"
            >
              <div class="card-header">
                <el-checkbox
                  :model-value="selectedConversations.includes(conversation.id)"
                  @change="toggleConversationSelect(conversation.id)"
                  class="card-checkbox"
                />
                <img
                  :src="getCharacterById(conversation.characterId)?.avatar"
                  :alt="getCharacterById(conversation.characterId)?.name"
                  class="card-avatar"
                />
              </div>

              <div class="card-content" @click="openConversation(conversation)">
                <h4 class="card-title">{{ conversation.title }}</h4>
                <p class="card-character">{{ getCharacterById(conversation.characterId)?.name }}</p>
                <p class="card-preview">{{ getLastMessage(conversation) }}</p>
              </div>

              <div class="card-footer">
                <div class="card-meta">
                  <span class="card-messages">{{ conversation.messages.length }} 条</span>
                  <span class="card-time">{{ formatRelativeTime(conversation.lastUpdate) }}</span>
                </div>
                <div class="card-actions">
                  <el-button @click.stop="continueConversation(conversation)" :icon="ChatDotRound" size="small" text />
                  <el-button @click.stop="deleteConversation(conversation.id)" :icon="Delete" size="small" text type="danger" />
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 时间线视图 -->
        <div v-else class="timeline-view">
          <div class="timeline-container">
            <div
              v-for="group in groupedConversations"
              :key="group.date"
              class="timeline-group"
            >
              <div class="timeline-date">
                <h3>{{ group.date }}</h3>
                <span class="group-count">{{ group.conversations.length }} 个对话</span>
              </div>
              <div class="timeline-items">
                <div
                  v-for="conversation in group.conversations"
                  :key="conversation.id"
                  class="timeline-item"
                  :class="{ selected: selectedConversations.includes(conversation.id) }"
                >
                  <div class="timeline-marker"></div>
                  <div class="timeline-content">
                    <div class="timeline-header">
                      <el-checkbox
                        :model-value="selectedConversations.includes(conversation.id)"
                        @change="toggleConversationSelect(conversation.id)"
                      />
                      <img
                        :src="getCharacterById(conversation.characterId)?.avatar"
                        :alt="getCharacterById(conversation.characterId)?.name"
                        class="timeline-avatar"
                      />
                      <div class="timeline-info">
                        <h4 @click="openConversation(conversation)">{{ conversation.title }}</h4>
                        <p>与{{ getCharacterById(conversation.characterId)?.name }}的对话</p>
                        <span>{{ formatTime(conversation.lastUpdate) }}</span>
                      </div>
                    </div>
                    <div class="timeline-actions">
                      <el-button @click="continueConversation(conversation)" size="small">继续对话</el-button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div v-if="filteredConversations.length > pageSize" class="pagination-section">
        <el-pagination
          v-model:current-page="currentPage"
          :page-size="pageSize"
          :total="filteredConversations.length"
          layout="total, prev, pager, next, jumper"
          @current-change="handlePageChange"
        />
      </div>
    </div>

    <!-- 对话详情弹窗 -->
    <el-dialog
      v-model="showDetailDialog"
      title="对话详情"
      width="800px"
      :before-close="closeDetailDialog"
    >
      <div v-if="selectedConversationDetail" class="conversation-detail">
        <div class="detail-header">
          <h3>{{ selectedConversationDetail.title }}</h3>
          <div class="detail-meta">
            <span>与{{ getCharacterById(selectedConversationDetail.characterId)?.name }}的对话</span>
            <span>{{ selectedConversationDetail.messages.length }} 条消息</span>
            <span>{{ formatTime(selectedConversationDetail.startTime) }} - {{ formatTime(selectedConversationDetail.lastUpdate) }}</span>
          </div>
        </div>
        
        <div class="detail-messages">
          <div
            v-for="message in selectedConversationDetail.messages"
            :key="message.id"
            class="detail-message"
            :class="{ 'user-message': message.type === 'user' }"
          >
            <div class="message-sender">
              {{ message.type === 'user' ? '用户' : getCharacterById(selectedConversationDetail.characterId)?.name }}
            </div>
            <div class="message-content">{{ message.content }}</div>
            <div class="message-time">{{ formatTime(message.timestamp) }}</div>
          </div>
        </div>
      </div>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDetailDialog">关闭</el-button>
          <el-button @click="exportConversation(selectedConversationDetail)" type="primary">导出对话</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useCharacterStore } from '../stores/character'
import { useChatStore } from '../stores/chat'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, ChatDotRound, Timer, User, Calendar } from '@element-plus/icons-vue'

const router = useRouter()
const characterStore = useCharacterStore()
const chatStore = useChatStore()

// 响应式数据
const searchQuery = ref('')
const selectedCharacter = ref('')
const dateRange = ref([])
const sortBy = ref('time')
const viewMode = ref('list')
const selectedConversations = ref([])
const selectAll = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const showDetailDialog = ref(false)
const selectedConversationDetail = ref(null)

// 计算属性
const characters = computed(() => characterStore.characters)
const allConversations = computed(() => chatStore.conversations)

const filteredConversations = computed(() => {
  let conversations = [...allConversations.value]

  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    conversations = conversations.filter(conv =>
      conv.title.toLowerCase().includes(query) ||
      conv.messages.some(msg => msg.content.toLowerCase().includes(query))
    )
  }

  // 角色过滤
  if (selectedCharacter.value) {
    conversations = conversations.filter(conv => conv.characterId === selectedCharacter.value)
  }

  // 日期过滤
  if (dateRange.value && dateRange.value.length === 2) {
    const [start, end] = dateRange.value
    conversations = conversations.filter(conv => {
      const convDate = new Date(conv.lastUpdate)
      return convDate >= start && convDate <= end
    })
  }

  // 排序
  conversations.sort((a, b) => {
    switch (sortBy.value) {
      case 'time':
        return new Date(b.lastUpdate) - new Date(a.lastUpdate)
      case 'messages':
        return b.messages.length - a.messages.length
      case 'character':
        return getCharacterById(a.characterId)?.name.localeCompare(getCharacterById(b.characterId)?.name) || 0
      default:
        return 0
    }
  })

  return conversations
})

const paginatedConversations = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredConversations.value.slice(start, end)
})

const groupedConversations = computed(() => {
  const groups = {}
  
  filteredConversations.value.forEach(conv => {
    const date = new Date(conv.lastUpdate).toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    })
    
    if (!groups[date]) {
      groups[date] = []
    }
    groups[date].push(conv)
  })

  return Object.keys(groups).map(date => ({
    date,
    conversations: groups[date]
  }))
})

const totalMessages = computed(() => {
  return filteredConversations.value.reduce((total, conv) => total + conv.messages.length, 0)
})

const uniqueCharacters = computed(() => {
  const characterIds = new Set(filteredConversations.value.map(conv => conv.characterId))
  return characterIds.size
})

const activeDays = computed(() => {
  const dates = new Set(
    filteredConversations.value.map(conv =>
      new Date(conv.lastUpdate).toDateString()
    )
  )
  return dates.size
})

const isIndeterminate = computed(() => {
  const selectedCount = selectedConversations.value.length
  const totalCount = paginatedConversations.value.length
  return selectedCount > 0 && selectedCount < totalCount
})

// 监听器
watch(selectAll, (val) => {
  if (val) {
    selectedConversations.value = paginatedConversations.value.map(conv => conv.id)
  } else {
    selectedConversations.value = []
  }
})

watch(paginatedConversations, () => {
  // 当分页数据变化时，重置选择状态
  const validSelections = selectedConversations.value.filter(id =>
    paginatedConversations.value.some(conv => conv.id === id)
  )
  selectedConversations.value = validSelections
  selectAll.value = validSelections.length === paginatedConversations.value.length
})

// 方法
const getCharacterById = (id) => {
  return characters.value.find(char => char.id === id)
}

const handleSearch = () => {
  currentPage.value = 1
}

const handleCharacterFilter = () => {
  currentPage.value = 1
}

const handleDateFilter = () => {
  currentPage.value = 1
}

const handleSort = () => {
  currentPage.value = 1
}

const handleViewChange = () => {
  selectedConversations.value = []
  selectAll.value = false
}

const handleSelectAll = () => {
  // selectAll的watch会处理具体逻辑
}

const toggleConversationSelect = (conversationId) => {
  const index = selectedConversations.value.indexOf(conversationId)
  if (index > -1) {
    selectedConversations.value.splice(index, 1)
  } else {
    selectedConversations.value.push(conversationId)
  }
  
  const selectedCount = selectedConversations.value.length
  const totalCount = paginatedConversations.value.length
  selectAll.value = selectedCount === totalCount
}

const handlePageChange = (page) => {
  currentPage.value = page
  selectedConversations.value = []
  selectAll.value = false
}

const openConversation = (conversation) => {
  selectedConversationDetail.value = conversation
  showDetailDialog.value = true
}

const closeDetailDialog = () => {
  showDetailDialog.value = false
  selectedConversationDetail.value = null
}

const continueConversation = (conversation) => {
  chatStore.selectConversation(conversation.id)
  characterStore.selectCharacter(conversation.characterId)
  router.push(`/chat/${conversation.characterId}`)
}

const exportConversation = (conversation) => {
  const character = getCharacterById(conversation.characterId)
  const content = conversation.messages.map(msg =>
    `${msg.type === 'user' ? '用户' : character?.name}: ${msg.content}`
  ).join('\n\n')

  const blob = new Blob([content], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${conversation.title}.txt`
  a.click()
  URL.revokeObjectURL(url)

  ElMessage.success('对话已导出')
}

const exportAll = () => {
  if (filteredConversations.value.length === 0) {
    ElMessage.warning('没有对话可导出')
    return
  }

  const allContent = filteredConversations.value.map(conversation => {
    const character = getCharacterById(conversation.characterId)
    const header = `=== ${conversation.title} ===\n时间: ${formatTime(conversation.startTime)} - ${formatTime(conversation.lastUpdate)}\n角色: ${character?.name}\n\n`
    const messages = conversation.messages.map(msg =>
      `${msg.type === 'user' ? '用户' : character?.name}: ${msg.content}`
    ).join('\n\n')
    return header + messages
  }).join('\n\n' + '='.repeat(50) + '\n\n')

  const blob = new Blob([allContent], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `对话历史_${new Date().toLocaleDateString()}.txt`
  a.click()
  URL.revokeObjectURL(url)

  ElMessage.success('所有对话已导出')
}

const deleteSelected = async () => {
  if (selectedConversations.value.length === 0) return

  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedConversations.value.length} 个对话吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    selectedConversations.value.forEach(id => {
      chatStore.deleteConversation(id)
    })

    selectedConversations.value = []
    selectAll.value = false
    ElMessage.success('选中的对话已删除')
  } catch {
    // 用户取消
  }
}

const deleteConversation = async (conversationId) => {
  try {
    await ElMessageBox.confirm('确定要删除这个对话吗？', '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    chatStore.deleteConversation(conversationId)
    
    // 从选中列表中移除
    const index = selectedConversations.value.indexOf(conversationId)
    if (index > -1) {
      selectedConversations.value.splice(index, 1)
    }

    ElMessage.success('对话已删除')
  } catch {
    // 用户取消
  }
}

const goToHome = () => {
  router.push('/')
}

const getLastMessage = (conversation) => {
  const lastMessage = conversation.messages[conversation.messages.length - 1]
  return lastMessage ? lastMessage.content.slice(0, 100) + (lastMessage.content.length > 100 ? '...' : '') : ''
}

const formatTime = (timestamp) => {
  return new Date(timestamp).toLocaleString('zh-CN')
}

const formatRelativeTime = (timestamp) => {
  const now = new Date()
  const time = new Date(timestamp)
  const diff = now - time
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  
  if (days === 0) return '今天'
  if (days === 1) return '昨天'
  if (days < 7) return `${days}天前`
  return time.toLocaleDateString('zh-CN')
}

const formatDuration = (conversation) => {
  const start = new Date(conversation.startTime)
  const end = new Date(conversation.lastUpdate)
  const diff = end - start
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
  
  if (hours > 0) {
    return `${hours}小时${minutes}分钟`
  }
  return `${minutes}分钟`
}
</script>

<style lang="scss" scoped>
.history-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.history-container {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  padding: 32px;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.page-header {
  text-align: center;
  margin-bottom: 32px;
  
  h1 {
    font-size: 32px;
    font-weight: bold;
    color: #303133;
    margin: 0 0 12px 0;
  }
  
  p {
    font-size: 16px;
    color: #606266;
    margin: 0;
  }
}

.search-section {
  margin-bottom: 24px;
  
  .search-controls {
    display: flex;
    gap: 16px;
    margin-bottom: 16px;
    flex-wrap: wrap;
    
    .search-input {
      flex: 1;
      min-width: 200px;
    }
    
    .character-filter,
    .date-filter {
      min-width: 150px;
    }
  }
  
  .action-controls {
    display: flex;
    gap: 8px;
    justify-content: flex-end;
  }
}

.stats-section {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  background: linear-gradient(135deg, #f6f9fc 0%, #e9ecef 100%);
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  
  .stat-icon {
    font-size: 24px;
    color: #667eea;
    background: rgba(102, 126, 234, 0.1);
    padding: 12px;
    border-radius: 50%;
  }
  
  .stat-content {
    .stat-number {
      display: block;
      font-size: 24px;
      font-weight: bold;
      color: #303133;
    }
    
    .stat-label {
      font-size: 14px;
      color: #909399;
    }
  }
}

.control-section {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
  
  .control-label {
    font-size: 14px;
    color: #606266;
    margin-right: 8px;
  }
  
  .sort-controls,
  .view-controls,
  .select-controls {
    display: flex;
    align-items: center;
  }
}

.conversations-section {
  margin-bottom: 24px;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
}

// 列表视图样式
.list-view {
  .conversation-item {
    display: flex;
    align-items: center;
    padding: 16px;
    border: 1px solid transparent;
    border-radius: 12px;
    margin-bottom: 12px;
    transition: all 0.3s ease;
    background: rgba(0, 0, 0, 0.02);
    
    &:hover {
      background: rgba(0, 0, 0, 0.05);
      transform: translateY(-2px);
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    }
    
    &.selected {
      border-color: #409EFF;
      background: rgba(64, 158, 255, 0.1);
    }
    
    .conversation-select {
      margin-right: 12px;
    }
    
    .conversation-avatar {
      margin-right: 16px;
      
      .character-avatar {
        width: 48px;
        height: 48px;
        border-radius: 50%;
        object-fit: cover;
      }
    }
    
    .conversation-content {
      flex: 1;
      cursor: pointer;
      
      .conversation-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 8px;
        
        .conversation-title {
          font-size: 16px;
          font-weight: bold;
          color: #303133;
          margin: 0;
        }
        
        .conversation-time {
          font-size: 13px;
          color: #909399;
        }
      }
      
      .conversation-meta {
        display: flex;
        gap: 16px;
        margin-bottom: 8px;
        font-size: 14px;
        color: #606266;
      }
      
      .conversation-preview {
        font-size: 14px;
        color: #909399;
        line-height: 1.4;
      }
    }
    
    .conversation-actions {
      display: flex;
      gap: 8px;
    }
  }
}

// 网格视图样式
.grid-view {
  .grid-container {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 20px;
  }
  
  .conversation-card {
    background: rgba(0, 0, 0, 0.02);
    border: 1px solid transparent;
    border-radius: 12px;
    padding: 16px;
    transition: all 0.3s ease;
    
    &:hover {
      background: rgba(0, 0, 0, 0.05);
      transform: translateY(-4px);
      box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
    }
    
    &.selected {
      border-color: #409EFF;
      background: rgba(64, 158, 255, 0.1);
    }
    
    .card-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 16px;
      
      .card-avatar {
        width: 40px;
        height: 40px;
        border-radius: 50%;
        object-fit: cover;
      }
    }
    
    .card-content {
      cursor: pointer;
      margin-bottom: 16px;
      
      .card-title {
        font-size: 16px;
        font-weight: bold;
        color: #303133;
        margin: 0 0 8px 0;
      }
      
      .card-character {
        font-size: 14px;
        color: #606266;
        margin: 0 0 8px 0;
      }
      
      .card-preview {
        font-size: 13px;
        color: #909399;
        line-height: 1.4;
        margin: 0;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
      }
    }
    
    .card-footer {
      display: flex;
      align-items: center;
      justify-content: space-between;
      
      .card-meta {
        font-size: 12px;
        color: #909399;
        
        .card-messages {
          margin-right: 8px;
        }
      }
      
      .card-actions {
        display: flex;
        gap: 4px;
      }
    }
  }
}

// 时间线视图样式
.timeline-view {
  .timeline-group {
    margin-bottom: 32px;
    
    .timeline-date {
      display: flex;
      align-items: center;
      gap: 12px;
      margin-bottom: 16px;
      
      h3 {
        font-size: 18px;
        font-weight: bold;
        color: #303133;
        margin: 0;
      }
      
      .group-count {
        font-size: 14px;
        color: #909399;
      }
    }
    
    .timeline-items {
      position: relative;
      
      &::before {
        content: '';
        position: absolute;
        left: 20px;
        top: 0;
        bottom: 0;
        width: 2px;
        background: #E4E7ED;
      }
    }
    
    .timeline-item {
      position: relative;
      padding-left: 60px;
      margin-bottom: 20px;
      
      &.selected .timeline-content {
        border-color: #409EFF;
        background: rgba(64, 158, 255, 0.05);
      }
      
      .timeline-marker {
        position: absolute;
        left: 12px;
        top: 16px;
        width: 16px;
        height: 16px;
        border-radius: 50%;
        background: #409EFF;
        border: 3px solid white;
        box-shadow: 0 0 0 1px #E4E7ED;
      }
      
      .timeline-content {
        background: rgba(0, 0, 0, 0.02);
        border: 1px solid transparent;
        border-radius: 12px;
        padding: 16px;
        transition: all 0.3s ease;
        
        &:hover {
          background: rgba(0, 0, 0, 0.05);
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
        }
        
        .timeline-header {
          display: flex;
          align-items: center;
          gap: 12px;
          margin-bottom: 12px;
          
          .timeline-avatar {
            width: 40px;
            height: 40px;
            border-radius: 50%;
            object-fit: cover;
          }
          
          .timeline-info {
            flex: 1;
            
            h4 {
              font-size: 16px;
              font-weight: bold;
              color: #303133;
              margin: 0 0 4px 0;
              cursor: pointer;
              
              &:hover {
                color: #409EFF;
              }
            }
            
            p {
              font-size: 14px;
              color: #606266;
              margin: 0 0 4px 0;
            }
            
            span {
              font-size: 13px;
              color: #909399;
            }
          }
        }
        
        .timeline-actions {
          text-align: right;
        }
      }
    }
  }
}

.pagination-section {
  display: flex;
  justify-content: center;
  margin-top: 32px;
}

// 对话详情弹窗样式
.conversation-detail {
  .detail-header {
    margin-bottom: 24px;
    
    h3 {
      font-size: 20px;
      font-weight: bold;
      color: #303133;
      margin: 0 0 8px 0;
    }
    
    .detail-meta {
      display: flex;
      gap: 16px;
      font-size: 14px;
      color: #606266;
      flex-wrap: wrap;
    }
  }
  
  .detail-messages {
    max-height: 400px;
    overflow-y: auto;
    
    .detail-message {
      margin-bottom: 16px;
      padding: 12px;
      border-radius: 8px;
      background: rgba(0, 0, 0, 0.02);
      
      &.user-message {
        background: rgba(64, 158, 255, 0.1);
        margin-left: 20px;
      }
      
      .message-sender {
        font-size: 12px;
        font-weight: bold;
        color: #606266;
        margin-bottom: 4px;
      }
      
      .message-content {
        font-size: 14px;
        color: #303133;
        line-height: 1.5;
        margin-bottom: 4px;
      }
      
      .message-time {
        font-size: 11px;
        color: #909399;
      }
    }
  }
}

.dialog-footer {
  text-align: right;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .history-page {
    padding: 0 10px;
  }
  
  .history-container {
    padding: 20px;
  }
  
  .page-header h1 {
    font-size: 24px;
  }
  
  .search-controls {
    flex-direction: column;
    
    .search-input,
    .character-filter,
    .date-filter {
      min-width: 100%;
    }
  }
  
  .control-section {
    flex-direction: column;
    align-items: stretch;
    
    .sort-controls,
    .view-controls,
    .select-controls {
      justify-content: center;
      margin: 8px 0;
    }
  }
  
  .stats-section {
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
  }
  
  .grid-view .grid-container {
    grid-template-columns: 1fr;
  }
  
  .conversation-item {
    flex-direction: column;
    align-items: stretch;
    
    .conversation-content {
      margin: 12px 0;
    }
    
    .conversation-actions {
      justify-content: center;
    }
  }
  
  .timeline-item {
    padding-left: 40px;
    
    .timeline-content .timeline-header {
      flex-direction: column;
      align-items: flex-start;
      
      .timeline-info {
        margin-top: 8px;
      }
    }
  }
}
</style>

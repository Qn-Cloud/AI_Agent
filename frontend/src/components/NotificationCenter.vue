<template>
  <div class="notification-center">
    <!-- 通知触发按钮 -->
    <el-badge :value="unreadCount" :hidden="unreadCount === 0" class="notification-badge">
      <el-button
        @click="togglePanel"
        :icon="Bell"
        circle
        size="large"
        class="notification-button"
        :class="{ active: showPanel }"
      />
    </el-badge>

    <!-- 通知面板 -->
    <transition name="slide-fade">
      <div v-if="showPanel" class="notification-panel">
        <div class="panel-header">
          <h4>通知中心</h4>
          <div class="header-actions">
            <el-button @click="markAllAsRead" text size="small" type="primary">
              全部已读
            </el-button>
            <el-button @click="clearAll" text size="small" type="danger">
              清空
            </el-button>
          </div>
        </div>

        <div class="panel-content">
          <!-- 通知筛选 -->
          <div class="notification-filters">
            <el-radio-group v-model="activeFilter" size="small">
              <el-radio-button label="all">全部</el-radio-button>
              <el-radio-button label="unread">未读</el-radio-button>
              <el-radio-button label="system">系统</el-radio-button>
              <el-radio-button label="chat">对话</el-radio-button>
            </el-radio-group>
          </div>

          <!-- 通知列表 -->
          <div class="notification-list">
            <div v-if="filteredNotifications.length === 0" class="empty-state">
              <el-icon><BellFilled /></el-icon>
              <p>暂无通知</p>
            </div>
            
            <div
              v-for="notification in filteredNotifications"
              :key="notification.id"
              class="notification-item"
              :class="{
                unread: !notification.read,
                [notification.type]: true
              }"
              @click="handleNotificationClick(notification)"
            >
              <div class="notification-icon">
                <el-icon>
                  <InfoFilled v-if="notification.type === 'info'" />
                  <SuccessFilled v-else-if="notification.type === 'success'" />
                  <WarningFilled v-else-if="notification.type === 'warning'" />
                  <CircleCloseFilled v-else-if="notification.type === 'error'" />
                  <ChatDotRound v-else-if="notification.type === 'chat'" />
                  <Setting v-else-if="notification.type === 'system'" />
                  <Bell v-else />
                </el-icon>
              </div>

              <div class="notification-content">
                <h5 class="notification-title">{{ notification.title }}</h5>
                <p class="notification-message">{{ notification.message }}</p>
                <div class="notification-meta">
                  <span class="notification-time">{{ formatTime(notification.createdAt) }}</span>
                  <span v-if="notification.category" class="notification-category">
                    {{ getCategoryName(notification.category) }}
                  </span>
                </div>
              </div>

              <div class="notification-actions">
                <el-button
                  v-if="!notification.read"
                  @click.stop="markAsRead(notification.id)"
                  size="small"
                  text
                  type="primary"
                >
                  标记已读
                </el-button>
                <el-button
                  @click.stop="removeNotification(notification.id)"
                  size="small"
                  text
                  type="danger"
                >
                  删除
                </el-button>
              </div>
            </div>
          </div>
        </div>

        <div class="panel-footer">
          <el-button @click="showPanel = false" size="small" style="width: 100%">
            关闭
          </el-button>
        </div>
      </div>
    </transition>

    <!-- 遮罩层 -->
    <div
      v-if="showPanel"
      class="notification-overlay"
      @click="showPanel = false"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Bell, BellFilled, InfoFilled, SuccessFilled, WarningFilled,
  CircleCloseFilled, ChatDotRound, Setting
} from '@element-plus/icons-vue'

// 响应式数据
const showPanel = ref(false)
const activeFilter = ref('all')
const notifications = ref([])

// 计算属性
const unreadCount = computed(() => {
  return notifications.value.filter(n => !n.read).length
})

const filteredNotifications = computed(() => {
  let filtered = notifications.value

  switch (activeFilter.value) {
    case 'unread':
      filtered = filtered.filter(n => !n.read)
      break
    case 'system':
      filtered = filtered.filter(n => n.category === 'system')
      break
    case 'chat':
      filtered = filtered.filter(n => n.category === 'chat')
      break
    default:
      // 'all' - 显示所有通知
      break
  }

  return filtered.sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt))
})

// 方法
const togglePanel = () => {
  showPanel.value = !showPanel.value
}

const addNotification = (notification) => {
  const newNotification = {
    id: Date.now().toString(),
    read: false,
    createdAt: new Date(),
    ...notification
  }
  
  notifications.value.unshift(newNotification)
  
  // 限制通知数量
  if (notifications.value.length > 50) {
    notifications.value = notifications.value.slice(0, 50)
  }
  
  // 保存到本地存储
  saveToLocalStorage()
}

const markAsRead = (notificationId) => {
  const notification = notifications.value.find(n => n.id === notificationId)
  if (notification) {
    notification.read = true
    saveToLocalStorage()
  }
}

const markAllAsRead = () => {
  notifications.value.forEach(n => n.read = true)
  saveToLocalStorage()
  ElMessage.success('所有通知已标记为已读')
}

const removeNotification = (notificationId) => {
  const index = notifications.value.findIndex(n => n.id === notificationId)
  if (index > -1) {
    notifications.value.splice(index, 1)
    saveToLocalStorage()
  }
}

const clearAll = () => {
  notifications.value = []
  saveToLocalStorage()
  ElMessage.success('所有通知已清空')
}

const handleNotificationClick = (notification) => {
  // 标记为已读
  if (!notification.read) {
    markAsRead(notification.id)
  }
  
  // 执行通知的动作
  if (notification.action) {
    notification.action()
  }
  
  // 关闭面板
  showPanel.value = false
}

const formatTime = (timestamp) => {
  const now = new Date()
  const time = new Date(timestamp)
  const diff = now - time
  
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)
  
  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 7) return `${days}天前`
  
  return time.toLocaleDateString('zh-CN')
}

const getCategoryName = (category) => {
  const categoryNames = {
    system: '系统',
    chat: '对话',
    user: '用户',
    error: '错误',
    update: '更新'
  }
  return categoryNames[category] || category
}

const saveToLocalStorage = () => {
  localStorage.setItem('notifications', JSON.stringify(notifications.value))
}

const loadFromLocalStorage = () => {
  try {
    const saved = localStorage.getItem('notifications')
    if (saved) {
      notifications.value = JSON.parse(saved).map(n => ({
        ...n,
        createdAt: new Date(n.createdAt)
      }))
    }
  } catch (error) {
    console.error('加载通知失败:', error)
  }
}

// 全局通知方法
const showNotification = (options) => {
  addNotification(options)
}

// 预设通知类型
const notify = {
  success: (title, message, action = null) => {
    showNotification({
      type: 'success',
      category: 'system',
      title,
      message,
      action
    })
  },
  
  error: (title, message, action = null) => {
    showNotification({
      type: 'error',
      category: 'system',
      title,
      message,
      action
    })
  },
  
  warning: (title, message, action = null) => {
    showNotification({
      type: 'warning',
      category: 'system',
      title,
      message,
      action
    })
  },
  
  info: (title, message, action = null) => {
    showNotification({
      type: 'info',
      category: 'system',
      title,
      message,
      action
    })
  },
  
  chat: (title, message, action = null) => {
    showNotification({
      type: 'chat',
      category: 'chat',
      title,
      message,
      action
    })
  }
}

// 监听点击外部关闭面板
const handleClickOutside = (event) => {
  if (showPanel.value && !event.target.closest('.notification-center')) {
    showPanel.value = false
  }
}

// 生命周期
onMounted(() => {
  loadFromLocalStorage()
  document.addEventListener('click', handleClickOutside)
  
  // 添加一些示例通知（开发时可以移除）
  if (notifications.value.length === 0) {
    addNotification({
      type: 'info',
      category: 'system',
      title: '欢迎使用AI角色扮演',
      message: '开始与您喜欢的AI角色对话吧！'
    })
  }
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})

// 暴露方法给父组件
defineExpose({
  showNotification,
  notify,
  addNotification,
  markAsRead,
  markAllAsRead,
  clearAll
})
</script>

<style lang="scss" scoped>
.notification-center {
  position: relative;
  display: inline-block;
}

.notification-badge {
  .notification-button {
    background: rgba(255, 255, 255, 0.9);
    border: 1px solid rgba(64, 158, 255, 0.2);
    color: #409EFF;
    transition: all 0.3s ease;
    
    &:hover {
      background: rgba(64, 158, 255, 0.1);
      border-color: #409EFF;
    }
    
    &.active {
      background: #409EFF;
      color: white;
    }
  }
}

.notification-panel {
  position: absolute;
  top: 100%;
  right: 0;
  width: 380px;
  max-height: 600px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
  border: 1px solid rgba(0, 0, 0, 0.1);
  z-index: 1000;
  overflow: hidden;
  margin-top: 8px;
}

.panel-header {
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #fafafa;
  
  h4 {
    margin: 0;
    font-size: 16px;
    font-weight: 600;
    color: #303133;
  }
  
  .header-actions {
    display: flex;
    gap: 8px;
  }
}

.panel-content {
  max-height: 480px;
  overflow-y: auto;
}

.notification-filters {
  padding: 12px 20px;
  border-bottom: 1px solid #f0f0f0;
  background: #fafafa;
}

.notification-list {
  .empty-state {
    text-align: center;
    padding: 40px 20px;
    color: #909399;
    
    .el-icon {
      font-size: 48px;
      margin-bottom: 12px;
    }
    
    p {
      margin: 0;
      font-size: 14px;
    }
  }
}

.notification-item {
  display: flex;
  align-items: flex-start;
  padding: 16px 20px;
  border-bottom: 1px solid #f5f5f5;
  cursor: pointer;
  transition: background-color 0.2s ease;
  
  &:hover {
    background: #f9f9f9;
  }
  
  &.unread {
    background: rgba(64, 158, 255, 0.02);
    border-left: 3px solid #409EFF;
    
    .notification-title {
      font-weight: 600;
    }
  }
  
  &:last-child {
    border-bottom: none;
  }
}

.notification-icon {
  margin-right: 12px;
  margin-top: 2px;
  
  .el-icon {
    font-size: 20px;
    
    .notification-item.success & {
      color: #67C23A;
    }
    
    .notification-item.warning & {
      color: #E6A23C;
    }
    
    .notification-item.error & {
      color: #F56C6C;
    }
    
    .notification-item.info & {
      color: #409EFF;
    }
    
    .notification-item.chat & {
      color: #909399;
    }
    
    .notification-item.system & {
      color: #606266;
    }
  }
}

.notification-content {
  flex: 1;
  min-width: 0;
  
  .notification-title {
    font-size: 14px;
    font-weight: 500;
    color: #303133;
    margin: 0 0 4px 0;
    line-height: 1.4;
  }
  
  .notification-message {
    font-size: 13px;
    color: #606266;
    line-height: 1.4;
    margin: 0 0 8px 0;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
  
  .notification-meta {
    display: flex;
    align-items: center;
    gap: 12px;
    
    .notification-time {
      font-size: 12px;
      color: #909399;
    }
    
    .notification-category {
      font-size: 12px;
      color: #409EFF;
      background: rgba(64, 158, 255, 0.1);
      padding: 2px 6px;
      border-radius: 4px;
    }
  }
}

.notification-actions {
  margin-left: 8px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.2s ease;
  
  .notification-item:hover & {
    opacity: 1;
  }
}

.panel-footer {
  padding: 12px 20px;
  border-top: 1px solid #f0f0f0;
  background: #fafafa;
}

.notification-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 999;
}

// 动画
.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: all 0.3s ease;
}

.slide-fade-enter-from {
  opacity: 0;
  transform: translateY(-10px);
}

.slide-fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

// 响应式设计
@media (max-width: 768px) {
  .notification-panel {
    width: 320px;
    max-height: 500px;
  }
  
  .notification-item {
    padding: 12px 16px;
    
    .notification-content {
      .notification-title {
        font-size: 13px;
      }
      
      .notification-message {
        font-size: 12px;
      }
    }
  }
}
</style>

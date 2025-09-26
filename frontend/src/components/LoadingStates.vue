<template>
  <!-- 全屏加载 -->
  <div v-if="type === 'fullscreen'" class="loading-fullscreen">
    <div class="loading-content">
      <div class="loading-spinner">
        <el-icon class="is-loading"><Loading /></el-icon>
      </div>
      <h3 class="loading-title">{{ title || '正在加载...' }}</h3>
      <p v-if="message" class="loading-message">{{ message }}</p>
      <div v-if="showProgress" class="loading-progress">
        <el-progress :percentage="progress" :show-text="false" />
        <span class="progress-text">{{ progress }}%</span>
      </div>
    </div>
  </div>

  <!-- 内联加载 -->
  <div v-else-if="type === 'inline'" class="loading-inline">
    <el-icon class="is-loading loading-icon"><Loading /></el-icon>
    <span class="loading-text">{{ title || '加载中...' }}</span>
  </div>

  <!-- 骨架屏 - 角色卡片 -->
  <div v-else-if="type === 'skeleton-card'" class="skeleton-card">
    <div class="skeleton-header">
      <el-skeleton-item variant="circle" style="width: 60px; height: 60px;" />
      <div class="skeleton-info">
        <el-skeleton-item variant="text" style="width: 120px;" />
        <el-skeleton-item variant="text" style="width: 80px;" />
      </div>
    </div>
    <div class="skeleton-content">
      <el-skeleton-item variant="text" style="width: 100%;" />
      <el-skeleton-item variant="text" style="width: 85%;" />
      <el-skeleton-item variant="text" style="width: 92%;" />
    </div>
    <div class="skeleton-tags">
      <el-skeleton-item variant="button" style="width: 60px; height: 24px; margin-right: 8px;" />
      <el-skeleton-item variant="button" style="width: 45px; height: 24px; margin-right: 8px;" />
      <el-skeleton-item variant="button" style="width: 55px; height: 24px;" />
    </div>
    <div class="skeleton-actions">
      <el-skeleton-item variant="button" style="width: 100px; height: 36px;" />
      <el-skeleton-item variant="button" style="width: 80px; height: 36px;" />
    </div>
  </div>

  <!-- 骨架屏 - 消息列表 -->
  <div v-else-if="type === 'skeleton-messages'" class="skeleton-messages">
    <div v-for="i in count" :key="i" class="skeleton-message" :class="{ 'user-message': i % 3 === 0 }">
      <div v-if="i % 3 !== 0" class="skeleton-avatar">
        <el-skeleton-item variant="circle" style="width: 32px; height: 32px;" />
      </div>
      <div class="skeleton-bubble">
        <el-skeleton-item variant="text" :style="{ width: getRandomWidth() }" />
        <el-skeleton-item variant="text" :style="{ width: getRandomWidth() }" />
        <el-skeleton-item v-if="Math.random() > 0.5" variant="text" :style="{ width: getRandomWidth() }" />
      </div>
    </div>
  </div>

  <!-- 骨架屏 - 对话历史列表 -->
  <div v-else-if="type === 'skeleton-history'" class="skeleton-history">
    <div v-for="i in count" :key="i" class="skeleton-history-item">
      <el-skeleton-item variant="circle" style="width: 48px; height: 48px;" />
      <div class="skeleton-history-content">
        <div class="skeleton-history-header">
          <el-skeleton-item variant="text" style="width: 200px;" />
          <el-skeleton-item variant="text" style="width: 100px;" />
        </div>
        <el-skeleton-item variant="text" style="width: 300px;" />
        <el-skeleton-item variant="text" style="width: 250px;" />
      </div>
    </div>
  </div>

  <!-- 自定义骨架屏 -->
  <div v-else-if="type === 'skeleton-custom'" class="skeleton-custom">
    <slot name="skeleton" />
  </div>

  <!-- 空状态加载 -->
  <div v-else-if="type === 'empty'" class="loading-empty">
    <div class="empty-icon">
      <el-icon><Box /></el-icon>
    </div>
    <h4 class="empty-title">{{ title || '暂无数据' }}</h4>
    <p v-if="message" class="empty-message">{{ message }}</p>
    <slot name="actions" />
  </div>

  <!-- 默认加载动画 -->
  <div v-else class="loading-default">
    <el-loading :loading="true" :text="title" />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { Loading, Box } from '@element-plus/icons-vue'

const props = defineProps({
  type: {
    type: String,
    default: 'default',
    validator: (value) => [
      'fullscreen', 'inline', 'skeleton-card', 'skeleton-messages', 
      'skeleton-history', 'skeleton-custom', 'empty', 'default'
    ].includes(value)
  },
  title: {
    type: String,
    default: ''
  },
  message: {
    type: String,
    default: ''
  },
  progress: {
    type: Number,
    default: 0
  },
  showProgress: {
    type: Boolean,
    default: false
  },
  count: {
    type: Number,
    default: 5
  }
})

const getRandomWidth = () => {
  const widths = ['60%', '75%', '85%', '90%', '95%', '100%']
  return widths[Math.floor(Math.random() * widths.length)]
}
</script>

<style lang="scss" scoped>
// 全屏加载
.loading-fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;

  .loading-content {
    text-align: center;
    max-width: 400px;
    padding: 40px;

    .loading-spinner {
      font-size: 48px;
      color: #409EFF;
      margin-bottom: 24px;
    }

    .loading-title {
      font-size: 20px;
      font-weight: 600;
      color: #303133;
      margin: 0 0 12px 0;
    }

    .loading-message {
      font-size: 14px;
      color: #606266;
      margin: 0 0 24px 0;
      line-height: 1.5;
    }

    .loading-progress {
      .progress-text {
        font-size: 14px;
        color: #909399;
        margin-top: 8px;
        display: block;
      }
    }
  }
}

// 内联加载
.loading-inline {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;

  .loading-icon {
    font-size: 20px;
    color: #409EFF;
    margin-right: 8px;
  }

  .loading-text {
    font-size: 14px;
    color: #606266;
  }
}

// 骨架屏 - 角色卡片
.skeleton-card {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  padding: 20px;
  border: 1px solid rgba(255, 255, 255, 0.2);

  .skeleton-header {
    display: flex;
    align-items: center;
    margin-bottom: 16px;

    .skeleton-info {
      margin-left: 16px;
      flex: 1;

      .el-skeleton-item {
        margin-bottom: 8px;
      }
    }
  }

  .skeleton-content {
    margin-bottom: 16px;

    .el-skeleton-item {
      margin-bottom: 8px;
    }
  }

  .skeleton-tags {
    display: flex;
    margin-bottom: 16px;
  }

  .skeleton-actions {
    display: flex;
    gap: 8px;
  }
}

// 骨架屏 - 消息列表
.skeleton-messages {
  padding: 20px;

  .skeleton-message {
    display: flex;
    margin-bottom: 20px;

    &.user-message {
      justify-content: flex-end;

      .skeleton-bubble {
        background: rgba(64, 158, 255, 0.1);
        border-radius: 16px 16px 4px 16px;
      }
    }

    .skeleton-avatar {
      margin-right: 12px;
    }

    .skeleton-bubble {
      background: rgba(0, 0, 0, 0.02);
      border-radius: 16px;
      padding: 12px 16px;
      max-width: 70%;

      .el-skeleton-item {
        margin-bottom: 4px;

        &:last-child {
          margin-bottom: 0;
        }
      }
    }
  }
}

// 骨架屏 - 对话历史
.skeleton-history {
  .skeleton-history-item {
    display: flex;
    align-items: center;
    padding: 16px;
    border-radius: 12px;
    background: rgba(0, 0, 0, 0.02);
    margin-bottom: 12px;

    .skeleton-history-content {
      margin-left: 16px;
      flex: 1;

      .skeleton-history-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 8px;
      }

      .el-skeleton-item {
        margin-bottom: 4px;

        &:last-child {
          margin-bottom: 0;
        }
      }
    }
  }
}

// 空状态
.loading-empty {
  text-align: center;
  padding: 60px 20px;

  .empty-icon {
    font-size: 64px;
    color: #DCDFE6;
    margin-bottom: 24px;
  }

  .empty-title {
    font-size: 18px;
    font-weight: 600;
    color: #303133;
    margin: 0 0 12px 0;
  }

  .empty-message {
    font-size: 14px;
    color: #606266;
    margin: 0 0 24px 0;
    line-height: 1.5;
  }
}

// 默认加载
.loading-default {
  min-height: 200px;
  position: relative;
}

// 动画
@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.4;
  }
}

.el-skeleton-item {
  animation: pulse 1.5s ease-in-out infinite;
}

// 响应式设计
@media (max-width: 768px) {
  .loading-fullscreen .loading-content {
    padding: 20px;

    .loading-title {
      font-size: 18px;
    }

    .loading-message {
      font-size: 13px;
    }
  }

  .skeleton-card {
    padding: 16px;
  }

  .skeleton-messages {
    padding: 16px;
  }

  .skeleton-message .skeleton-bubble {
    max-width: 85%;
  }
}
</style>

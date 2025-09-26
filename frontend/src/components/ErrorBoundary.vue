<template>
  <div v-if="hasError" class="error-boundary">
    <div class="error-container">
      <div class="error-icon">
        <el-icon><WarningFilled /></el-icon>
      </div>
      <h3 class="error-title">{{ errorTitle }}</h3>
      <p class="error-message">{{ errorMessage }}</p>
      <div class="error-actions">
        <el-button @click="retry" type="primary">重试</el-button>
        <el-button @click="goHome" type="default">返回首页</el-button>
        <el-button @click="reportError" type="info" text>报告问题</el-button>
      </div>
      <details v-if="isDevelopment" class="error-details">
        <summary>错误详情</summary>
        <pre>{{ errorDetails }}</pre>
      </details>
    </div>
  </div>
  <slot v-else />
</template>

<script setup>
import { ref, onErrorCaptured } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { WarningFilled } from '@element-plus/icons-vue'
import config from '../config'

const router = useRouter()

const hasError = ref(false)
const errorTitle = ref('')
const errorMessage = ref('')
const errorDetails = ref('')
const isDevelopment = config.isDevelopment

const errorTypes = {
  network: {
    title: '网络连接错误',
    message: '无法连接到服务器，请检查网络连接后重试'
  },
  auth: {
    title: '身份验证失败',
    message: '登录状态已过期，请重新登录'
  },
  permission: {
    title: '权限不足',
    message: '您没有权限访问此功能'
  },
  notFound: {
    title: '页面未找到',
    message: '抱歉，您访问的页面不存在'
  },
  server: {
    title: '服务器错误',
    message: '服务器遇到了一些问题，请稍后重试'
  },
  default: {
    title: '出现了一些问题',
    message: '应用程序遇到了意外错误，我们正在努力修复'
  }
}

const setError = (error, type = 'default') => {
  hasError.value = true
  
  const errorType = errorTypes[type] || errorTypes.default
  errorTitle.value = errorType.title
  errorMessage.value = errorType.message
  
  if (isDevelopment) {
    errorDetails.value = error.stack || error.message || JSON.stringify(error, null, 2)
  }
  
  console.error('ErrorBoundary caught error:', error)
}

const retry = () => {
  hasError.value = false
  errorTitle.value = ''
  errorMessage.value = ''
  errorDetails.value = ''
  
  // 重新加载当前路由
  router.go(0)
}

const goHome = () => {
  hasError.value = false
  router.push('/')
}

const reportError = () => {
  ElMessage.info('错误报告功能开发中...')
}

// 捕获子组件错误
onErrorCaptured((error, instance, info) => {
  console.error('Error captured:', error, info)
  
  // 根据错误类型设置不同的错误信息
  let errorType = 'default'
  
  if (error.message?.includes('Network Error') || error.code === 'NETWORK_ERROR') {
    errorType = 'network'
  } else if (error.response?.status === 401) {
    errorType = 'auth'
  } else if (error.response?.status === 403) {
    errorType = 'permission'
  } else if (error.response?.status === 404) {
    errorType = 'notFound'
  } else if (error.response?.status >= 500) {
    errorType = 'server'
  }
  
  setError(error, errorType)
  
  // 阻止错误继续传播
  return false
})

// 暴露方法给父组件
defineExpose({
  setError
})
</script>

<style lang="scss" scoped>
.error-boundary {
  min-height: 400px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
}

.error-container {
  text-align: center;
  max-width: 500px;
  width: 100%;
}

.error-icon {
  font-size: 64px;
  color: #F56C6C;
  margin-bottom: 24px;
}

.error-title {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
  margin: 0 0 16px 0;
}

.error-message {
  font-size: 16px;
  color: #606266;
  line-height: 1.6;
  margin: 0 0 32px 0;
}

.error-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
  flex-wrap: wrap;
  margin-bottom: 24px;
}

.error-details {
  text-align: left;
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 16px;
  margin-top: 24px;
  
  summary {
    cursor: pointer;
    font-weight: bold;
    color: #606266;
    margin-bottom: 12px;
    
    &:hover {
      color: #409EFF;
    }
  }
  
  pre {
    font-size: 12px;
    color: #303133;
    line-height: 1.4;
    white-space: pre-wrap;
    word-break: break-all;
    margin: 0;
  }
}

@media (max-width: 768px) {
  .error-container {
    padding: 0 20px;
  }
  
  .error-title {
    font-size: 20px;
  }
  
  .error-message {
    font-size: 14px;
  }
  
  .error-actions {
    flex-direction: column;
    
    .el-button {
      width: 100%;
    }
  }
}</style>

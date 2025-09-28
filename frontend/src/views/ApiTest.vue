<template>
  <div style="padding: 20px;">
    <h2>API测试页面</h2>
    
    <div style="margin: 20px 0;">
      <h3>配置信息</h3>
      <pre>{{ configInfo }}</pre>
    </div>
    
    <div style="margin: 20px 0;">
      <el-button @click="testChatApi" type="primary">测试聊天API</el-button>
      <el-button @click="testCharacterApi" type="success">测试角色API</el-button>
      <el-button @click="testChatStore" type="warning">测试Chat Store</el-button>
      <el-button @click="clearResults" type="info">清空结果</el-button>
    </div>
    
    <div style="margin: 20px 0;">
      <h3>测试结果</h3>
      <pre style="max-height: 400px; overflow-y: auto;">{{ testResults }}</pre>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { chatApiService, characterApiService } from '../services'
import { useChatStore } from '../stores/chat'
import config from '../config'

const configInfo = ref('')
const testResults = ref('')
const chatStore = useChatStore()

onMounted(() => {
  configInfo.value = JSON.stringify({
    chatBaseUrl: config.api.chatBaseUrl,
    characterBaseUrl: config.api.characterBaseUrl,
    环境变量: {
      VITE_CHAT_API_URL: import.meta.env.VITE_CHAT_API_URL,
      VITE_CHARACTER_API_URL: import.meta.env.VITE_CHARACTER_API_URL,
      NODE_ENV: import.meta.env.NODE_ENV,
      DEV: import.meta.env.DEV
    }
  }, null, 2)
})

const clearResults = () => {
  testResults.value = ''
}

const testChatApi = async () => {
  testResults.value = '正在测试聊天API...\n'
  
  try {
    console.log('🧪 开始测试聊天API')
    const result = await chatApiService.getConversationHistory({
      page: 1,
      pageSize: 10
    })
    
    testResults.value += `✅ 聊天API测试成功:\n${JSON.stringify(result, null, 2)}\n`
  } catch (error) {
    testResults.value += `❌ 聊天API测试失败:\n${error.message}\n${error.stack}\n`
    console.error('聊天API测试失败:', error)
  }
}

const testCharacterApi = async () => {
  testResults.value += '\n正在测试角色API...\n'
  
  try {
    console.log('🧪 开始测试角色API')
    const result = await characterApiService.getCharacterList({
      page: 1,
      pageSize: 10
    })
    
    testResults.value += `✅ 角色API测试成功:\n${JSON.stringify(result, null, 2)}\n`
  } catch (error) {
    testResults.value += `❌ 角色API测试失败:\n${error.message}\n${error.stack}\n`
    console.error('角色API测试失败:', error)
  }
}

const testChatStore = async () => {
  testResults.value += '\n正在测试Chat Store...\n'
  
  try {
    console.log('🧪 开始测试Chat Store')
    
    // 测试 testApiConnection 方法
    const result = await chatStore.testApiConnection()
    
    testResults.value += `✅ Chat Store测试成功:\n${JSON.stringify(result, null, 2)}\n`
  } catch (error) {
    testResults.value += `❌ Chat Store测试失败:\n${error.message}\n${error.stack}\n`
    console.error('Chat Store测试失败:', error)
  }
}
</script> 
<template>
  <div class="chat-page">
    <!-- 聊天头部 -->
    <div class="chat-header">
      <div class="header-left">
        <el-button @click="goBack" :icon="ArrowLeft" size="large" text>
          返回
        </el-button>
      </div>
      <div class="header-center" v-if="currentCharacter">
        <div class="character-avatar">
          <img 
            :src="currentCharacter.avatar" 
            :alt="currentCharacter.name"
            class="avatar-image"
            @error="handleImageError"
          />
        </div>
        <div class="character-info">
          <h3 class="character-name">{{ currentCharacter.name }}</h3>
          <span class="character-status">在线</span>
        </div>
      </div>
      <div class="header-right">
        <el-button @click="testSSEConnection" type="success" size="small">
          测试SSE
        </el-button>
        <el-button @click="testSpeechService" type="primary" size="small">
          测试语音
        </el-button>
        <el-button @click="openSettings" :icon="Setting" size="large" text>
          设置
        </el-button>
      </div>
    </div>

    <!-- 聊天内容区域 -->
    <div class="chat-content">
      <!-- 消息列表 -->
      <div class="messages-container" ref="messagesContainer">
        <!-- 空状态 -->
        <div v-if="currentMessages.length === 0" class="empty-chat">
          <div class="empty-content">
            <div class="empty-avatar">
              <img 
                v-if="currentCharacter"
                :src="currentCharacter.avatar" 
                :alt="currentCharacter.name"
                class="avatar-image"
                @error="handleImageError"
              />
              <div v-else class="avatar-placeholder">
                <el-icon><ChatDotRound /></el-icon>
              </div>
            </div>
            <h3>开始与{{ currentCharacter?.name }}对话</h3>
            <p>你可以通过语音或文字与AI角色进行对话</p>
            <div class="quick-start-tips">
              <div
                v-for="tip in quickTips"
                :key="tip"
                class="tip-button"
                @click="sendQuickMessage(tip)"
              >
                {{ tip }}
              </div>
            </div>
          </div>
        </div>

        <!-- 加载状态 -->
        <div v-else-if="isLoading && currentMessages.length === 0" class="loading-messages">
          <LoadingStates type="skeleton-messages" :count="3" />
        </div>

        <!-- 消息列表 -->
        <div v-else class="messages-list">
          <div
            v-for="message in currentMessages"
            :key="message.id"
            class="message-item"
            :class="{ 'user-message': message.type === 'user', 'ai-message': message.type === 'ai' }"
          >
            <div class="message-avatar" v-if="message.type === 'ai'">
              <img 
                v-if="currentCharacter"
                :src="currentCharacter.avatar" 
                :alt="currentCharacter.name"
                class="avatar-image"
                @error="handleImageError"
              />
              <div v-else class="avatar-placeholder">
                <el-icon><User /></el-icon>
              </div>
            </div>
            <div class="message-content">
              <div class="message-bubble" :class="{ 'streaming': message.isStreaming }">
                <div class="message-text">
                  {{ message.content }}
                  <span v-if="message.isStreaming" class="streaming-cursor">|</span>
                </div>
              </div>
              <div class="message-time">
                {{ formatMessageTime(message.timestamp) }}
                <span v-if="message.isStreaming" class="streaming-status">正在输入...</span>
              </div>
            </div>
          </div>

          <!-- 加载状态 -->
          <div v-if="isLoading" class="loading-message">
            <div class="message-avatar">
              <img 
                v-if="currentCharacter"
                :src="currentCharacter.avatar" 
                :alt="currentCharacter.name"
                class="avatar-image"
                @error="handleImageError"
              />
              <div v-else class="avatar-placeholder">
                <el-icon><User /></el-icon>
              </div>
            </div>
            <div class="message-content">
              <div class="message-bubble loading-bubble">
                <div class="typing-indicator">
                  <span></span>
                  <span></span>
                  <span></span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 输入区域 -->
    <div class="input-area">
      <!-- 语音录制按钮 -->
      <div class="voice-section">
        <button
          class="voice-btn"
          :class="{ recording: isRecording, disabled: isLoading }"
          @mousedown="startRecording"
          @mouseup="stopRecording"
          @touchstart="startRecording"
          @touchend="stopRecording"
          :disabled="isLoading"
        >
          <el-icon v-if="!isRecording"><Microphone /></el-icon>
          <el-icon v-else class="recording-icon"><Microphone /></el-icon>
        </button>
        <span class="voice-tip">{{ voiceTipText }}</span>
      </div>

      <!-- 文字输入 -->
      <div class="text-input-section">
        <div class="input-wrapper">
          <el-input
            v-model="textInput"
            placeholder="输入消息..."
            @keydown.enter.exact="handleEnterSend"
            class="text-input"
            :disabled="isLoading"
          />
          <el-button
            @click="sendTextMessage"
            type="primary"
            :disabled="!textInput.trim() || isLoading"
            :loading="isLoading"
            class="send-btn"
            size="large"
          >
            <el-icon><ArrowRight /></el-icon>
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useCharacterStore } from '../stores/character'
import { useChatStore } from '../stores/chat'
import { ElMessage } from 'element-plus'
import { ArrowLeft, Setting, User, ChatDotRound, Microphone, ArrowRight } from '@element-plus/icons-vue'
import LoadingStates from '../components/LoadingStates.vue'
import { speechApiService } from '../services/speech'

const route = useRoute()
const router = useRouter()
const characterStore = useCharacterStore()
const chatStore = useChatStore()

// 响应式数据
const textInput = ref('')
const messagesContainer = ref(null)

// 计算属性
const currentCharacter = computed(() => characterStore.currentCharacter)
const currentConversation = computed(() => chatStore.currentConversation)
const currentMessages = computed(() => chatStore.currentMessages)
const isLoading = computed(() => chatStore.isLoading)
const isRecording = computed(() => chatStore.isRecording)
const transcript = computed(() => chatStore.transcript)

const voiceTipText = computed(() => {
  if (isLoading.value) return '处理中...'
  if (isRecording.value) return '松开发送'
  return '按住说话'
})

const quickTips = computed(() => {
  if (!currentCharacter.value) return []
  
  const tips = {
    'harry-potter': ['你好，哈利！', '告诉我霍格沃茨的故事', '你的朋友们怎么样？'],
    'socrates': ['什么是智慧？', '如何获得真正的知识？', '美德与知识的关系是什么？'],
    'shakespeare': ['谈谈爱情与人生', '创作的灵感来自哪里？', '如何理解人性？'],
    'einstein': ['相对论是什么？', '科学与哲学的关系', '想象力的重要性'],
    'sherlock': ['如何进行推理？', '观察的技巧有哪些？', '分析这个案例'],
    'hermione': ['学习方法有哪些？', '魔法世界的知识', '如何解决问题？']
  }
  
  return tips[currentCharacter.value.id] || ['你好！', '介绍一下自己', '分享一些想法']
})

// 方法
const initializeChat = async (characterId) => {
  // 如果当前角色已经设置且匹配，则跳过选择角色的步骤
  if (!characterStore.currentCharacter || 
      (characterStore.currentCharacter.id != characterId && 
       characterStore.currentCharacter.id !== String(characterId))) {
    // 只有在当前角色不匹配时才重新选择角色
  if (!characterStore.currentCharacter || characterStore.currentCharacter.id != characterId) {
    await characterStore.selectCharacter(characterId)
  }
  }
  
  if (!currentConversation.value) {
    await chatStore.startNewConversation(characterId)
  }
}

// 监听器
watch(() => route.params.characterId, async (newId) => {
  if (newId) {
    try {
      await initializeChat(newId)
    } catch (error) {
      ElMessage.error('初始化聊天失败: ' + error.message)
    }
  }
}, { immediate: true })

watch(currentMessages, () => {
  nextTick(() => {
    scrollToBottom()
  })
}, { deep: true })

// 监听语音转文字结果
watch(() => chatStore.transcript, (newTranscript) => {
  if (newTranscript && newTranscript.trim()) {
    textInput.value = newTranscript
    console.log('📝 语音转文字结果已填入输入框:', newTranscript)
    // 清空transcript，避免重复设置
    chatStore.transcript = ''
  }
})

const goBack = () => {
  router.push('/')
}

const openSettings = () => {
  router.push('/settings')
}

const sendQuickMessage = async (message) => {
  textInput.value = message
  await sendTextMessage()
}

const sendTextMessage = async () => {
  if (!textInput.value.trim()) return
  
  console.log('🚀 开始发送消息:', textInput.value.trim())
  console.log('📋 当前对话:', chatStore.currentConversation)
  console.log('🎭 当前角色:', currentCharacter.value)
  
  try {
    await chatStore.sendMessage(textInput.value.trim())
    textInput.value = ''
    console.log('✅ 消息发送成功')
  } catch (error) {
    console.error('❌ 消息发送失败:', error)
    
    // 如果是context canceled错误，提示用户重试
    if (error.message.includes('context canceled') || error.message.includes('连接中断')) {
      ElMessage({
        message: '服务器响应超时，正在自动重试...',
        type: 'warning',
        duration: 2000
      })
      
      // 自动重试一次
      setTimeout(async () => {
        try {
          await chatStore.sendMessage(textInput.value.trim())
          textInput.value = ''
          ElMessage.success('重试成功！')
        } catch (retryError) {
          ElMessage.error('重试失败，请稍后再试: ' + retryError.message)
        }
      }, 3000)
    } else {
      ElMessage.error('发送消息失败: ' + error.message)
    }
  }
}

const handleEnterSend = (event) => {
  event.preventDefault()
  sendTextMessage()
}

const startRecording = async () => {
  if (isLoading.value) return
  try {
    await chatStore.startRecording()
    ElMessage.info('开始录音，松开按钮发送')
  } catch (error) {
    ElMessage.error('录音启动失败: ' + error.message)
  }
}

const stopRecording = async () => {
  if (!isRecording.value) return
  try {
    await chatStore.stopRecording()
  } catch (error) {
    ElMessage.error('录音处理失败: ' + error.message)
  }
}

const scrollToBottom = () => {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

const formatMessageTime = (timestamp) => {
  return new Date(timestamp).toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

const handleImageError = (event) => {
  // 图片加载失败时显示默认头像
  event.target.style.display = 'none'
  const placeholder = document.createElement('div')
  placeholder.className = 'avatar-placeholder'
  placeholder.innerHTML = '<el-icon><User /></el-icon>'
  event.target.parentNode.appendChild(placeholder)
}

// 测试语音服务连接
const testSpeechService = async () => {
  try {
    console.log('🔍 开始测试语音服务连接...')
    ElMessage.info('正在测试语音服务连接...')
    
    const response = await speechApiService.healthCheck()
    console.log('✅ 语音服务连接成功:', response)
    ElMessage.success('语音服务连接正常！')
  } catch (error) {
    console.error('❌ 语音服务连接失败:', error)
    ElMessage.error(`语音服务连接失败: ${error.message}`)
  }
}

// 测试SSE连接
const testSSEConnection = async () => {
  try {
    console.log('�� 开始测试SSE连接...')
    ElMessage.info('正在测试SSE连接...')
    
    // 模拟一个简单的SSE请求
    const response = await fetch('http://localhost:8000/sse/test') // 替换为实际的SSE端点
    if (response.ok) {
      console.log('✅ SSE连接成功')
      ElMessage.success('SSE连接正常！')
    } else {
      console.error('❌ SSE连接失败:', response.status, response.statusText)
      ElMessage.error(`SSE连接失败: ${response.status} ${response.statusText}`)
    }
  } catch (error) {
    console.error('❌ SSE连接失败:', error)
    ElMessage.error(`SSE连接失败: ${error.message}`)
  }
}

// 生命周期
onMounted(() => {
  scrollToBottom()
})
</script>

<style lang="scss" scoped>
.chat-page {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f5f5;
}

// 聊天头部
.chat-header {
  background: white;
  border-bottom: 1px solid #e5e5e5;
  padding: 16px 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;

  .header-left,
  .header-right {
    flex: 1;
  }

  .header-center {
    display: flex;
    align-items: center;
    gap: 12px;
    flex: 2;
    justify-content: center;

    .character-avatar {
      .avatar-image {
        width: 40px;
        height: 40px;
        border-radius: 50%;
        object-fit: cover;
        object-position: center;
      }
      
      .avatar-placeholder {
        width: 40px;
        height: 40px;
        background: #4A90E2;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        color: white;
        font-size: 18px;
      }
    }

    .character-info {
      .character-name {
        font-size: 16px;
        font-weight: 600;
        color: #333;
        margin: 0 0 2px 0;
      }

      .character-status {
        font-size: 12px;
        color: #67C23A;
      }
    }
  }

  .header-right {
    text-align: right;
  }
}

// 聊天内容区域
.chat-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;

  .messages-container {
    flex: 1;
    overflow-y: auto;
    padding: 20px;
    background: white;
    margin: 20px 20px 0 20px;
    border-radius: 12px 12px 0 0;

    // 空状态
    .empty-chat {
      height: 100%;
      display: flex;
      align-items: center;
      justify-content: center;

      .empty-content {
        text-align: center;
        max-width: 400px;
        padding: 40px;

        .empty-avatar {
          margin-bottom: 24px;

          .avatar-image {
            width: 80px;
            height: 80px;
            border-radius: 50%;
            object-fit: cover;
            object-position: center;
            margin: 0 auto;
          }

          .avatar-placeholder {
            width: 80px;
            height: 80px;
            background: #4A90E2;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
            font-size: 32px;
            margin: 0 auto;
          }
        }

        h3 {
          font-size: 20px;
          color: #333;
          margin: 0 0 8px 0;
        }

        p {
          color: #666;
          margin: 0 0 24px 0;
          line-height: 1.5;
        }

        .quick-start-tips {
          display: flex;
          flex-wrap: wrap;
          gap: 8px;
          justify-content: center;

          .tip-button {
            padding: 8px 16px;
            background: #f0f9ff;
            border: 1px solid #4A90E2;
            border-radius: 20px;
            color: #4A90E2;
            font-size: 14px;
            cursor: pointer;
            transition: all 0.2s ease;

            &:hover {
              background: #4A90E2;
              color: white;
            }
          }
        }
      }
    }

    // 消息列表
    .messages-list {
      .message-item {
        display: flex;
        margin-bottom: 20px;

        &.user-message {
          justify-content: flex-end;

          .message-content {
            .message-bubble {
              background: #4A90E2;
              color: white;
            }
          }
        }

        &.ai-message {
          justify-content: flex-start;

          .message-avatar {
            margin-right: 12px;

            .avatar-image {
              width: 32px;
              height: 32px;
              border-radius: 50%;
              object-fit: cover;
              object-position: center;
            }

            .avatar-placeholder {
              width: 32px;
              height: 32px;
              background: #4A90E2;
              border-radius: 50%;
              display: flex;
              align-items: center;
              justify-content: center;
              color: white;
              font-size: 14px;
            }
          }
        }

        .message-content {
          max-width: 70%;

          .message-bubble {
            background: #f5f5f5;
            border-radius: 16px;
            padding: 12px 16px;
            
            &.streaming {
              background: linear-gradient(90deg, #f5f5f5 0%, #e8f4fd 50%, #f5f5f5 100%);
              background-size: 200% 100%;
              animation: streamingGlow 2s ease-in-out infinite;
            }

            .message-text {
              line-height: 1.5;
              word-break: break-word;
              
              .streaming-cursor {
                display: inline-block;
                animation: blink 1s infinite;
                font-weight: bold;
                color: #4A90E2;
              }
            }
          }

          .message-time {
            font-size: 12px;
            color: #999;
            margin-top: 4px;
            text-align: center;
            
            .streaming-status {
              color: #4A90E2;
              font-style: italic;
              margin-left: 8px;
            }
          }
        }
      }

      .loading-message {
        display: flex;
        margin-bottom: 20px;

        .message-avatar {
          margin-right: 12px;

          .avatar-image {
            width: 32px;
            height: 32px;
            border-radius: 50%;
            object-fit: cover;
            object-position: center;
          }

          .avatar-placeholder {
            width: 32px;
            height: 32px;
            background: #4A90E2;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
            font-size: 14px;
          }
        }

        .loading-bubble {
          background: #f5f5f5;

          .typing-indicator {
            display: flex;
            gap: 4px;

            span {
              width: 8px;
              height: 8px;
              border-radius: 50%;
              background: #999;
              animation: typing 1.4s infinite ease-in-out;

              &:nth-child(2) {
                animation-delay: 0.2s;
              }

              &:nth-child(3) {
                animation-delay: 0.4s;
              }
            }
          }
        }
      }
    }
  }
}

// 输入区域
.input-area {
  background: white;
  padding: 20px;
  margin: 0 20px 20px 20px;
  border-radius: 0 0 12px 12px;
  border-top: 1px solid #f0f0f0;

  .voice-section {
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 16px;

    .voice-btn {
      width: 60px;
      height: 60px;
      border-radius: 50%;
      background: #4A90E2;
      border: none;
      color: white;
      font-size: 24px;
      cursor: pointer;
      transition: all 0.3s ease;
      margin-right: 12px;

      &:hover {
        background: #357abd;
        transform: scale(1.05);
      }

      &.recording {
        background: #ff4757;
        animation: pulse 1s infinite;

        .recording-icon {
          animation: recording-pulse 1s infinite;
        }
      }

      &.disabled {
        background: #ccc;
        cursor: not-allowed;
        transform: none;
      }
    }

    .voice-tip {
      font-size: 14px;
      color: #666;
    }
  }

  .text-input-section {
    .input-wrapper {
      display: flex;
      gap: 12px;
      align-items: flex-end;

      .text-input {
        flex: 1;

        :deep(.el-input__wrapper) {
          border-radius: 20px;
          min-height: 40px;
          border: 1px solid #e5e5e5;

          &:hover {
            border-color: #4A90E2;
          }

          &.is-focus {
            border-color: #4A90E2;
            box-shadow: 0 0 0 2px rgba(74, 144, 226, 0.2);
          }
        }

        :deep(.el-input__inner) {
          padding: 8px 16px;
          font-size: 14px;
        }
      }

      .send-btn {
        border-radius: 20px;
        padding: 8px 16px;
        height: 40px;
      }
    }
  }
}

// 动画
@keyframes typing {
  0%, 60%, 100% {
    transform: translateY(0);
    opacity: 0.3;
  }
  30% {
    transform: translateY(-10px);
    opacity: 1;
  }
}

@keyframes pulse {
  0% { transform: scale(1); }
  50% { transform: scale(1.1); }
  100% { transform: scale(1); }
}

@keyframes recording-pulse {
  0% { opacity: 1; }
  50% { opacity: 0.5; }
  100% { opacity: 1; }
}

// 响应式设计
@media (max-width: 768px) {
  .chat-header {
    padding: 12px 16px;

    .header-center {
      .character-avatar .avatar-placeholder {
        width: 32px;
        height: 32px;
        font-size: 14px;
      }

      .character-info .character-name {
        font-size: 14px;
      }
    }
  }

  .chat-content {
    .messages-container {
      padding: 16px;
      margin: 16px 16px 0 16px;

      .empty-chat .empty-content {
        padding: 24px;

        .empty-avatar .avatar-placeholder {
          width: 60px;
          height: 60px;
          font-size: 24px;
        }

        h3 {
          font-size: 18px;
        }
      }

      .messages-list .message-item .message-content {
        max-width: 85%;
      }
    }
  }

  .input-area {
    padding: 16px;
    margin: 0 16px 16px 16px;

    .voice-section .voice-btn {
      width: 50px;
      height: 50px;
      font-size: 20px;
    }
  }

// 流式回复动画效果
@keyframes streamingGlow {
  0% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0% 50%;
  }
}

@keyframes blink {
  0%, 50% {
    opacity: 1;
  }
  51%, 100% {
    opacity: 0;
  }
}

// 流式回复样式
.message-bubble.streaming {
  background: linear-gradient(90deg, #f5f5f5 0%, #e8f4fd 50%, #f5f5f5 100%) !important;
  background-size: 200% 100%;
  animation: streamingGlow 2s ease-in-out infinite;
}

.streaming-cursor {
  display: inline-block;
  animation: blink 1s infinite;
  font-weight: bold;
  color: #4A90E2;
}

.streaming-status {
  color: #4A90E2;
  font-style: italic;
  margin-left: 8px;
}
}
</style>
<template>
  <div class="voice-controller">
    <!-- 语音输入区域 -->
    <div class="voice-input-section">
      <div class="record-area">
        <button
          class="record-btn"
          :class="{ 
            active: isRecording, 
            disabled: isLoading || isSpeaking 
          }"
          @mousedown="startRecording"
          @mouseup="stopRecording"
          @touchstart="startRecording"
          @touchend="stopRecording"
          :disabled="isLoading || isSpeaking"
        >
          <div class="record-icon">
            <el-icon v-if="!isRecording">
              <Microphone />
            </el-icon>
            <div v-else class="recording-animation">
              <div class="pulse"></div>
              <el-icon><Microphone /></el-icon>
            </div>
          </div>
          <span class="record-text">
            {{ recordButtonText }}
          </span>
        </button>

        <!-- 语音波形显示 -->
        <div v-if="isRecording" class="wave-form">
          <div 
            v-for="i in 20" 
            :key="i" 
            class="wave-bar"
            :style="{ 
              height: waveData[i - 1] + '%',
              animationDelay: (i * 0.1) + 's'
            }"
          ></div>
        </div>
      </div>

      <!-- 识别结果显示 -->
      <div v-if="transcript" class="transcript-display">
        <div class="transcript-content">
          <el-icon class="transcript-icon"><ChatLineRound /></el-icon>
          <span class="transcript-text">{{ transcript }}</span>
        </div>
        <el-button
          type="primary"
          size="small"
          @click="sendTranscript"
          :disabled="!transcript.trim()"
        >
          发送
        </el-button>
      </div>
    </div>

    <!-- 语音输出区域 -->
    <div class="voice-output-section">
      <div class="playback-controls">
        <el-button
          :type="isSpeaking ? 'danger' : 'success'"
          @click="toggleSpeaking"
          :disabled="!lastMessage || isLoading"
          class="speak-btn"
        >
          <el-icon>
            <VideoPause v-if="isSpeaking" />
            <VideoPlay v-else />
          </el-icon>
          {{ isSpeaking ? '停止播放' : '重播最后消息' }}
        </el-button>

        <!-- 音量控制 -->
        <div class="volume-control">
          <el-icon class="volume-icon"><Headset /></el-icon>
          <el-slider
            v-model="volume"
            :min="0"
            :max="100"
            @change="updateVolume"
            class="volume-slider"
            :show-tooltip="false"
          />
          <span class="volume-text">{{ volume }}%</span>
        </div>
      </div>

      <!-- 语音设置 -->
      <div class="voice-settings">
        <el-collapse accordion>
          <el-collapse-item title="语音设置" name="voice-settings">
            <div class="setting-group">
              <div class="setting-item">
                <label>语速</label>
                <el-slider
                  v-model="voiceRate"
                  :min="0.5"
                  :max="2"
                  :step="0.1"
                  @change="updateVoiceSettings"
                  :show-tooltip="false"
                />
                <span class="setting-value">{{ voiceRate.toFixed(1) }}x</span>
              </div>

              <div class="setting-item">
                <label>音调</label>
                <el-slider
                  v-model="voicePitch"
                  :min="0.5"
                  :max="2"
                  :step="0.1"
                  @change="updateVoiceSettings"
                  :show-tooltip="false"
                />
                <span class="setting-value">{{ voicePitch.toFixed(1) }}</span>
              </div>

              <div class="setting-item">
                <el-checkbox
                  v-model="autoSpeak"
                  @change="updateAutoSpeak"
                >
                  自动播放AI回复
                </el-checkbox>
              </div>
            </div>
          </el-collapse-item>
        </el-collapse>
      </div>
    </div>

    <!-- 状态指示器 -->
    <div v-if="isLoading" class="loading-indicator">
      <el-icon class="loading-icon"><Loading /></el-icon>
      <span>AI正在思考中...</span>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useChatStore } from '../stores/chat'
import { ElMessage } from 'element-plus'
import { Microphone, ChatLineRound, Loading, VideoPause, VideoPlay, Headset } from '@element-plus/icons-vue'

const chatStore = useChatStore()

// 响应式数据
const volume = ref(chatStore.voiceSettings.volume)
const voiceRate = ref(chatStore.voiceSettings.rate)
const voicePitch = ref(chatStore.voiceSettings.pitch)
const autoSpeak = ref(chatStore.voiceSettings.autoSpeak)
const waveData = ref(Array.from({ length: 20 }, () => Math.random() * 100))

// 计算属性
const isRecording = computed(() => chatStore.isRecording)
const isSpeaking = computed(() => chatStore.isSpeaking)
const isLoading = computed(() => chatStore.isLoading)
const transcript = computed(() => chatStore.transcript)
const lastMessage = computed(() => chatStore.lastAIMessage?.content)

const recordButtonText = computed(() => {
  if (isLoading.value) return '处理中...'
  if (isSpeaking.value) return '播放中...'
  if (isRecording.value) return '松开发送'
  return '按住说话'
})

// 动画定时器
let waveTimer = null

// 方法
const startRecording = () => {
  if (isLoading.value || isSpeaking.value) return
  
  chatStore.startRecording()
  startWaveAnimation()
  ElMessage.info('开始录音，松开按钮发送')
}

const stopRecording = () => {
  if (!isRecording.value) return
  
  chatStore.stopRecording()
  stopWaveAnimation()
}

const sendTranscript = () => {
  if (transcript.value.trim()) {
    chatStore.sendMessage(transcript.value)
  }
}

const toggleSpeaking = () => {
  if (isSpeaking.value) {
    chatStore.stopSpeaking()
  } else if (lastMessage.value) {
    chatStore.speakMessage(lastMessage.value)
  }
}

const updateVolume = (value) => {
  chatStore.updateVoiceSettings({ volume: value })
}

const updateVoiceSettings = () => {
  chatStore.updateVoiceSettings({
    rate: voiceRate.value,
    pitch: voicePitch.value
  })
}

const updateAutoSpeak = (value) => {
  chatStore.updateVoiceSettings({ autoSpeak: value })
}

const startWaveAnimation = () => {
  waveTimer = setInterval(() => {
    waveData.value = waveData.value.map(() => 20 + Math.random() * 80)
  }, 100)
}

const stopWaveAnimation = () => {
  if (waveTimer) {
    clearInterval(waveTimer)
    waveTimer = null
  }
}

// 监听器
watch(volume, (newValue) => {
  updateVolume(newValue)
})

// 生命周期
onMounted(() => {
  // 检查浏览器语音支持
  if (!('webkitSpeechRecognition' in window) && !('SpeechRecognition' in window)) {
    ElMessage.warning('您的浏览器不支持语音识别功能')
  }
  
  if (!('speechSynthesis' in window)) {
    ElMessage.warning('您的浏览器不支持语音合成功能')
  }
})

onUnmounted(() => {
  stopWaveAnimation()
})
</script>

<style lang="scss" scoped>
.voice-controller {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  padding: 24px;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.voice-input-section {
  margin-bottom: 24px;
}

.record-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 16px;
}

.record-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: 50%;
  width: 120px;
  height: 120px;
  color: white;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.3);
  
  &:hover:not(.disabled) {
    transform: scale(1.05);
    box-shadow: 0 12px 24px rgba(102, 126, 234, 0.4);
  }
  
  &.active {
    background: linear-gradient(135deg, #f56c6c 0%, #e91e63 100%);
    animation: pulse 1.5s infinite;
  }
  
  &.disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  
  .record-icon {
    font-size: 32px;
    position: relative;
  }
  
  .record-text {
    font-size: 14px;
    font-weight: 500;
  }
}

.recording-animation {
  position: relative;
  
  .pulse {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 60px;
    height: 60px;
    border: 2px solid rgba(255, 255, 255, 0.6);
    border-radius: 50%;
    animation: pulse-ring 1.5s infinite;
  }
}

.wave-form {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 3px;
  height: 60px;
  margin-top: 20px;
}

.wave-bar {
  width: 4px;
  background: linear-gradient(to top, #667eea, #764ba2);
  border-radius: 2px;
  animation: wave 0.5s infinite alternate;
  min-height: 10px;
}

.transcript-display {
  background: rgba(103, 194, 58, 0.1);
  border: 1px solid rgba(103, 194, 58, 0.2);
  border-radius: 12px;
  padding: 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  
  .transcript-content {
    flex: 1;
    display: flex;
    align-items: center;
    gap: 8px;
    
    .transcript-icon {
      color: #67C23A;
      font-size: 18px;
    }
    
    .transcript-text {
      color: #529b2e;
      font-weight: 500;
    }
  }
}

.voice-output-section {
  .playback-controls {
    display: flex;
    align-items: center;
    gap: 16px;
    margin-bottom: 20px;
    flex-wrap: wrap;
  }
  
  .speak-btn {
    border-radius: 8px;
  }
  
  .volume-control {
    display: flex;
    align-items: center;
    gap: 8px;
    flex: 1;
    min-width: 200px;
    
    .volume-icon {
      color: #909399;
      font-size: 18px;
    }
    
    .volume-slider {
      flex: 1;
    }
    
    .volume-text {
      font-size: 12px;
      color: #909399;
      min-width: 35px;
    }
  }
}

.voice-settings {
  .setting-group {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }
  
  .setting-item {
    display: flex;
    align-items: center;
    gap: 12px;
    
    label {
      font-size: 14px;
      color: #606266;
      min-width: 40px;
    }
    
    .el-slider {
      flex: 1;
    }
    
    .setting-value {
      font-size: 12px;
      color: #909399;
      min-width: 35px;
    }
  }
}

.loading-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 16px;
  background: rgba(64, 158, 255, 0.1);
  border: 1px solid rgba(64, 158, 255, 0.2);
  border-radius: 12px;
  color: #409EFF;
  margin-top: 16px;
  
  .loading-icon {
    animation: spin 1s linear infinite;
  }
}

/* 动画 */
@keyframes pulse {
  0% { transform: scale(1); }
  50% { transform: scale(1.1); }
  100% { transform: scale(1); }
}

@keyframes pulse-ring {
  0% {
    transform: translate(-50%, -50%) scale(0.8);
    opacity: 1;
  }
  100% {
    transform: translate(-50%, -50%) scale(1.2);
    opacity: 0;
  }
}

@keyframes wave {
  0% { height: 20%; }
  100% { height: 100%; }
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .voice-controller {
    padding: 16px;
  }
  
  .record-btn {
    width: 100px;
    height: 100px;
    
    .record-icon {
      font-size: 28px;
    }
    
    .record-text {
      font-size: 12px;
    }
  }
  
  .playback-controls {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
    
    .volume-control {
      min-width: 100%;
    }
  }
  
  .transcript-display {
    flex-direction: column;
    align-items: stretch;
    gap: 8px;
  }
}
</style>

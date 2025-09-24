<template>
  <div class="settings-page">
    <div class="settings-container">
      <div class="settings-header">
        <h1>角色设定</h1>
        <p>自定义您的AI角色，打造独特的对话体验</p>
      </div>

      <!-- 角色选择 -->
      <div class="character-selector">
        <h3>选择要设置的角色</h3>
        <div class="character-grid">
          <div
            v-for="character in characters"
            :key="character.id"
            class="character-option"
            :class="{ active: selectedCharacter?.id === character.id }"
            @click="selectCharacter(character)"
          >
            <img :src="character.avatar" :alt="character.name" />
            <span>{{ character.name }}</span>
          </div>
        </div>
      </div>

      <!-- 设置面板 -->
      <div v-if="selectedCharacter" class="settings-panels">
        <el-tabs v-model="activeTab" type="border-card">
          <!-- 基础信息 -->
          <el-tab-pane label="基础信息" name="basic">
            <div class="setting-section">
              <h4>角色基础信息</h4>
              <el-form :model="characterForm" label-width="100px">
                <el-form-item label="角色名称">
                  <el-input v-model="characterForm.name" />
                </el-form-item>
                <el-form-item label="角色描述">
                  <el-input
                    v-model="characterForm.description"
                    type="textarea"
                    :rows="4"
                    placeholder="描述角色的背景、特点等..."
                  />
                </el-form-item>
                <el-form-item label="角色标签">
                  <div class="tags-input">
                    <el-tag
                      v-for="tag in characterForm.tags"
                      :key="tag"
                      closable
                      @close="removeTag(tag)"
                      class="tag-item"
                    >
                      {{ tag }}
                    </el-tag>
                    <el-input
                      v-if="tagInputVisible"
                      ref="tagInput"
                      v-model="tagInputValue"
                      size="small"
                      @keyup.enter="addTag"
                      @blur="addTag"
                      class="tag-input"
                    />
                    <el-button
                      v-else
                      @click="showTagInput"
                      size="small"
                      type="primary"
                      plain
                    >
                      + 添加标签
                    </el-button>
                  </div>
                </el-form-item>
              </el-form>
            </div>
          </el-tab-pane>

          <!-- 性格设定 -->
          <el-tab-pane label="性格特征" name="personality">
            <div class="setting-section">
              <h4>性格特征调整</h4>
              <div class="personality-sliders">
                <div class="slider-item">
                  <label>友善度</label>
                  <el-slider
                    v-model="characterForm.personality.friendliness"
                    :min="0"
                    :max="100"
                    show-tooltip
                    :format-tooltip="formatTooltip"
                  />
                  <span class="value">{{ characterForm.personality.friendliness }}</span>
                </div>
                <div class="slider-item">
                  <label>幽默感</label>
                  <el-slider
                    v-model="characterForm.personality.humor"
                    :min="0"
                    :max="100"
                    show-tooltip
                    :format-tooltip="formatTooltip"
                  />
                  <span class="value">{{ characterForm.personality.humor }}</span>
                </div>
                <div class="slider-item">
                  <label>智慧</label>
                  <el-slider
                    v-model="characterForm.personality.intelligence"
                    :min="0"
                    :max="100"
                    show-tooltip
                    :format-tooltip="formatTooltip"
                  />
                  <span class="value">{{ characterForm.personality.intelligence }}</span>
                </div>
                <div class="slider-item">
                  <label>创造力</label>
                  <el-slider
                    v-model="characterForm.personality.creativity"
                    :min="0"
                    :max="100"
                    show-tooltip
                    :format-tooltip="formatTooltip"
                  />
                  <span class="value">{{ characterForm.personality.creativity }}</span>
                </div>
              </div>

              <div class="personality-preview">
                <h5>性格预览</h5>
                <div class="preview-radar">
                  <!-- 简单的雷达图预览 -->
                  <div class="radar-chart">
                    <div class="radar-bg"></div>
                    <div class="radar-data">
                      <div
                        class="radar-point"
                        v-for="(value, key) in characterForm.personality"
                        :key="key"
                        :style="getRadarPointStyle(key, value)"
                      ></div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </el-tab-pane>

          <!-- 对话风格 -->
          <el-tab-pane label="对话风格" name="style">
            <div class="setting-section">
              <h4>对话风格设置</h4>
              <el-form :model="characterForm" label-width="120px">
                <el-form-item label="正式程度">
                  <el-radio-group v-model="characterForm.style.formality">
                    <el-radio label="formal">正式</el-radio>
                    <el-radio label="casual">随意</el-radio>
                    <el-radio label="mixed">混合</el-radio>
                  </el-radio-group>
                </el-form-item>
                <el-form-item label="回答长度">
                  <el-radio-group v-model="characterForm.style.responseLength">
                    <el-radio label="brief">简洁</el-radio>
                    <el-radio label="detailed">详细</el-radio>
                    <el-radio label="adaptive">自适应</el-radio>
                  </el-radio-group>
                </el-form-item>
                <el-form-item label="话题倾向">
                  <el-checkbox-group v-model="characterForm.style.topics">
                    <el-checkbox label="philosophy">哲学思考</el-checkbox>
                    <el-checkbox label="science">科学探索</el-checkbox>
                    <el-checkbox label="art">艺术文化</el-checkbox>
                    <el-checkbox label="life">生活感悟</el-checkbox>
                    <el-checkbox label="history">历史知识</el-checkbox>
                  </el-checkbox-group>
                </el-form-item>
              </el-form>
            </div>
          </el-tab-pane>

          <!-- 自定义Prompt -->
          <el-tab-pane label="自定义Prompt" name="prompt">
            <div class="setting-section">
              <h4>自定义角色提示词</h4>
              
              <!-- 模板选择 -->
              <div class="template-section">
                <h5>选择模板</h5>
                <el-select v-model="selectedTemplate" @change="loadTemplate">
                  <el-option
                    v-for="template in promptTemplates"
                    :key="template.id"
                    :label="template.name"
                    :value="template.id"
                  />
                </el-select>
              </div>

              <!-- Prompt编辑器 -->
              <div class="prompt-editor">
                <h5>自定义提示词</h5>
                <el-input
                  v-model="characterForm.prompt"
                  type="textarea"
                  :rows="8"
                  placeholder="在这里输入自定义的角色设定..."
                />
                
                <div class="prompt-tips">
                  <h6>提示词建议：</h6>
                  <ul>
                    <li>明确角色身份和背景</li>
                    <li>描述说话风格和用词习惯</li>
                    <li>设定知识领域和兴趣爱好</li>
                    <li>定义交互方式和回应模式</li>
                  </ul>
                </div>
              </div>

              <!-- 效果预览 -->
              <div class="prompt-preview">
                <h5>效果预览</h5>
                <div class="preview-content">
                  {{ generatePromptPreview() }}
                </div>
                <el-button @click="testPrompt" type="primary">测试效果</el-button>
              </div>
            </div>
          </el-tab-pane>

          <!-- 语音设置 -->
          <el-tab-pane label="语音设置" name="voice">
            <div class="setting-section">
              <h4>语音合成设置</h4>
              <el-form :model="characterForm.voiceSettings" label-width="120px">
                <el-form-item label="语速">
                  <el-slider
                    v-model="characterForm.voiceSettings.rate"
                    :min="0.5"
                    :max="2"
                    :step="0.1"
                    show-tooltip
                    :format-tooltip="(val) => val.toFixed(1) + 'x'"
                  />
                </el-form-item>
                <el-form-item label="音调">
                  <el-slider
                    v-model="characterForm.voiceSettings.pitch"
                    :min="0.5"
                    :max="2"
                    :step="0.1"
                    show-tooltip
                    :format-tooltip="(val) => val.toFixed(1)"
                  />
                </el-form-item>
                <el-form-item label="音量">
                  <el-slider
                    v-model="characterForm.voiceSettings.volume"
                    :min="0"
                    :max="1"
                    :step="0.1"
                    show-tooltip
                    :format-tooltip="(val) => Math.round(val * 100) + '%'"
                  />
                </el-form-item>
              </el-form>
              
              <div class="voice-test">
                <el-button @click="testVoice" type="primary">
                  <el-icon><Headset /></el-icon>
                  测试语音
                </el-button>
              </div>
            </div>
          </el-tab-pane>
        </el-tabs>

        <!-- 操作按钮 -->
        <div class="settings-actions">
          <el-button @click="resetSettings">重置</el-button>
          <el-button @click="previewChanges" type="info">预览</el-button>
          <el-button @click="saveSettings" type="primary" :loading="saving">
            保存设置
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import { useCharacterStore } from '../stores/character'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Headset } from '@element-plus/icons-vue'

const characterStore = useCharacterStore()

// 响应式数据
const activeTab = ref('basic')
const selectedCharacter = ref(null)
const saving = ref(false)
const tagInputVisible = ref(false)
const tagInputValue = ref('')
const tagInput = ref(null)
const selectedTemplate = ref('')

// 表单数据
const characterForm = ref({
  name: '',
  description: '',
  tags: [],
  personality: {
    friendliness: 50,
    humor: 50,
    intelligence: 50,
    creativity: 50
  },
  style: {
    formality: 'casual',
    responseLength: 'adaptive',
    topics: []
  },
  prompt: '',
  voiceSettings: {
    rate: 1.0,
    pitch: 1.0,
    volume: 0.8
  }
})

// Prompt模板
const promptTemplates = ref([
  {
    id: 'default',
    name: '默认模板',
    content: '你是{{name}}，{{description}}。请用{{name}}的语气和视角来回答问题。'
  },
  {
    id: 'detailed',
    name: '详细模板',
    content: '你是{{name}}，{{description}}。你的性格特点是友善度{{friendliness}}%，幽默感{{humor}}%，智慧{{intelligence}}%，创造力{{creativity}}%。请始终保持角色设定，用{{name}}的语气、知识背景和思维方式来回答问题。'
  },
  {
    id: 'educational',
    name: '教育模板',
    content: '你是{{name}}，一位知识渊博的{{description}}。你善于教学和启发他人思考。请用简单易懂的方式解答问题，并适时提出引导性问题来促进深入思考。'
  },
  {
    id: 'conversational',
    name: '对话模板',
    content: '你是{{name}}，{{description}}。你喜欢与人交流，善于倾听。请用温和友善的语气回答问题，并表现出对对话者的关心和兴趣。'
  }
])

// 计算属性
const characters = computed(() => characterStore.characters)

// 监听器
watch(selectedCharacter, (character) => {
  if (character) {
    loadCharacterData(character)
  }
})

// 方法
const selectCharacter = (character) => {
  selectedCharacter.value = character
  characterStore.selectCharacter(character.id)
}

const loadCharacterData = (character) => {
  characterForm.value = {
    name: character.name,
    description: character.description,
    tags: [...character.tags],
    personality: { ...character.personality },
    style: {
      formality: character.style?.formality || 'casual',
      responseLength: character.style?.responseLength || 'adaptive',
      topics: character.style?.topics || []
    },
    prompt: character.prompt,
    voiceSettings: { ...character.voiceSettings }
  }
}

const removeTag = (tag) => {
  const index = characterForm.value.tags.indexOf(tag)
  if (index > -1) {
    characterForm.value.tags.splice(index, 1)
  }
}

const showTagInput = () => {
  tagInputVisible.value = true
  nextTick(() => {
    tagInput.value?.focus()
  })
}

const addTag = () => {
  const tag = tagInputValue.value.trim()
  if (tag && !characterForm.value.tags.includes(tag)) {
    characterForm.value.tags.push(tag)
  }
  tagInputVisible.value = false
  tagInputValue.value = ''
}

const formatTooltip = (value) => {
  if (value <= 20) return '很低'
  if (value <= 40) return '较低'
  if (value <= 60) return '中等'
  if (value <= 80) return '较高'
  return '很高'
}

const getRadarPointStyle = (key, value) => {
  const angles = {
    friendliness: 0,
    humor: 90,
    intelligence: 180,
    creativity: 270
  }
  
  const angle = angles[key] || 0
  const radius = (value / 100) * 80 // 最大半径80px
  const x = 50 + radius * Math.cos(angle * Math.PI / 180)
  const y = 50 + radius * Math.sin(angle * Math.PI / 180)
  
  return {
    left: x + '%',
    top: y + '%'
  }
}

const loadTemplate = (templateId) => {
  const template = promptTemplates.value.find(t => t.id === templateId)
  if (template) {
    characterForm.value.prompt = template.content
  }
}

const generatePromptPreview = () => {
  let preview = characterForm.value.prompt
  
  // 替换模板变量
  preview = preview.replace(/\{\{name\}\}/g, characterForm.value.name)
  preview = preview.replace(/\{\{description\}\}/g, characterForm.value.description)
  preview = preview.replace(/\{\{friendliness\}\}/g, characterForm.value.personality.friendliness)
  preview = preview.replace(/\{\{humor\}\}/g, characterForm.value.personality.humor)
  preview = preview.replace(/\{\{intelligence\}\}/g, characterForm.value.personality.intelligence)
  preview = preview.replace(/\{\{creativity\}\}/g, characterForm.value.personality.creativity)
  
  return preview || '请输入自定义提示词...'
}

const testPrompt = () => {
  ElMessage.info('提示词测试功能开发中...')
}

const testVoice = () => {
  const testText = `你好，我是${characterForm.value.name}。这是语音测试。`
  
  if ('speechSynthesis' in window) {
    const utterance = new SpeechSynthesisUtterance(testText)
    utterance.rate = characterForm.value.voiceSettings.rate
    utterance.pitch = characterForm.value.voiceSettings.pitch
    utterance.volume = characterForm.value.voiceSettings.volume
    
    speechSynthesis.speak(utterance)
  } else {
    ElMessage.warning('您的浏览器不支持语音合成功能')
  }
}

const resetSettings = async () => {
  try {
    await ElMessageBox.confirm('确定要重置设置吗？这将恢复到默认配置。', '确认重置', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    if (selectedCharacter.value) {
      loadCharacterData(selectedCharacter.value)
      ElMessage.success('设置已重置')
    }
  } catch {
    // 用户取消
  }
}

const previewChanges = () => {
  ElMessage.info('预览功能开发中...')
}

const saveSettings = async () => {
  if (!selectedCharacter.value) return
  
  saving.value = true
  
  try {
    // 更新角色数据
    characterStore.updateCharacterPrompt(selectedCharacter.value.id, characterForm.value.prompt)
    characterStore.updateCharacterPersonality(selectedCharacter.value.id, characterForm.value.personality)
    characterStore.updateVoiceSettings(selectedCharacter.value.id, characterForm.value.voiceSettings)
    
    // 模拟保存延迟
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    ElMessage.success('设置已保存')
  } catch (error) {
    ElMessage.error('保存失败，请重试')
  } finally {
    saving.value = false
  }
}

// 初始化
if (characterStore.currentCharacter) {
  selectCharacter(characterStore.currentCharacter)
} else if (characters.value.length > 0) {
  selectCharacter(characters.value[0])
}
</script>

<style lang="scss" scoped>
.settings-page {
  max-width: 1000px;
  margin: 0 auto;
  padding: 0 20px;
}

.settings-container {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  padding: 32px;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.settings-header {
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

.character-selector {
  margin-bottom: 32px;
  
  h3 {
    margin: 0 0 16px 0;
    color: #303133;
  }
  
  .character-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
    gap: 16px;
  }
  
  .character-option {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 16px;
    border: 2px solid transparent;
    border-radius: 12px;
    cursor: pointer;
    transition: all 0.3s ease;
    background: rgba(0, 0, 0, 0.02);
    
    &:hover {
      background: rgba(0, 0, 0, 0.05);
      transform: translateY(-2px);
    }
    
    &.active {
      border-color: #409EFF;
      background: rgba(64, 158, 255, 0.1);
    }
    
    img {
      width: 60px;
      height: 60px;
      border-radius: 50%;
      object-fit: cover;
      margin-bottom: 8px;
    }
    
    span {
      font-size: 14px;
      color: #303133;
      text-align: center;
    }
  }
}

.settings-panels {
  .el-tabs {
    border: none;
    
    .el-tabs__header {
      margin-bottom: 24px;
    }
    
    .el-tabs__content {
      padding: 0;
    }
  }
}

.setting-section {
  h4 {
    margin: 0 0 24px 0;
    color: #303133;
    font-size: 18px;
    border-bottom: 1px solid rgba(0, 0, 0, 0.1);
    padding-bottom: 8px;
  }
  
  h5 {
    margin: 20px 0 12px 0;
    color: #606266;
    font-size: 16px;
  }
  
  h6 {
    margin: 12px 0 8px 0;
    color: #909399;
    font-size: 14px;
  }
}

.tags-input {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
  
  .tag-item {
    border-radius: 12px;
  }
  
  .tag-input {
    width: 100px;
  }
}

.personality-sliders {
  .slider-item {
    display: flex;
    align-items: center;
    margin-bottom: 20px;
    
    label {
      font-size: 14px;
      color: #606266;
      width: 80px;
      margin-right: 16px;
    }
    
    .el-slider {
      flex: 1;
      margin-right: 16px;
    }
    
    .value {
      font-size: 14px;
      color: #303133;
      font-weight: bold;
      width: 30px;
    }
  }
}

.personality-preview {
  margin-top: 32px;
  
  .radar-chart {
    position: relative;
    width: 200px;
    height: 200px;
    margin: 0 auto;
    
    .radar-bg {
      width: 100%;
      height: 100%;
      border: 2px solid #DCDFE6;
      border-radius: 50%;
      position: relative;
      
      &::before,
      &::after {
        content: '';
        position: absolute;
        border: 1px solid #E4E7ED;
        border-radius: 50%;
      }
      
      &::before {
        width: 66%;
        height: 66%;
        top: 17%;
        left: 17%;
      }
      
      &::after {
        width: 33%;
        height: 33%;
        top: 33.5%;
        left: 33.5%;
      }
    }
    
    .radar-data {
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
    }
    
    .radar-point {
      position: absolute;
      width: 8px;
      height: 8px;
      background: #409EFF;
      border-radius: 50%;
      transform: translate(-50%, -50%);
      border: 2px solid white;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
    }
  }
}

.template-section {
  margin-bottom: 24px;
  
  .el-select {
    width: 200px;
  }
}

.prompt-editor {
  margin-bottom: 24px;
  
  .el-textarea {
    .el-textarea__inner {
      font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
      line-height: 1.5;
    }
  }
  
  .prompt-tips {
    margin-top: 16px;
    padding: 16px;
    background: rgba(103, 194, 58, 0.1);
    border-radius: 8px;
    border: 1px solid rgba(103, 194, 58, 0.2);
    
    ul {
      margin: 8px 0 0 0;
      padding-left: 20px;
      
      li {
        color: #529b2e;
        font-size: 14px;
        line-height: 1.6;
      }
    }
  }
}

.prompt-preview {
  .preview-content {
    background: #f6f9fc;
    border: 1px solid #E4E7ED;
    border-radius: 8px;
    padding: 16px;
    font-size: 14px;
    line-height: 1.5;
    color: #303133;
    margin-bottom: 16px;
    min-height: 80px;
  }
}

.voice-test {
  margin-top: 24px;
  text-align: center;
}

.settings-actions {
  margin-top: 32px;
  display: flex;
  justify-content: center;
  gap: 16px;
  padding-top: 24px;
  border-top: 1px solid rgba(0, 0, 0, 0.1);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .settings-page {
    padding: 0 10px;
  }
  
  .settings-container {
    padding: 20px;
  }
  
  .settings-header {
    h1 {
      font-size: 24px;
    }
    
    p {
      font-size: 14px;
    }
  }
  
  .character-grid {
    grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
    gap: 12px;
  }
  
  .character-option {
    padding: 12px;
    
    img {
      width: 50px;
      height: 50px;
    }
    
    span {
      font-size: 12px;
    }
  }
  
  .personality-sliders .slider-item {
    flex-direction: column;
    align-items: stretch;
    
    label {
      width: auto;
      margin-bottom: 8px;
    }
    
    .el-slider {
      margin: 0 0 8px 0;
    }
    
    .value {
      text-align: center;
      width: auto;
    }
  }
  
  .settings-actions {
    flex-direction: column;
    
    .el-button {
      width: 100%;
    }
  }
}
</style>

<template>
  <div class="character-editor">
    <div class="editor-header">
      <h2>{{ isEdit ? '编辑角色' : '自定义您的AI角色，打造独特的对话体验' }}</h2>
    </div>

    <div class="editor-content">
      <!-- 角色选择 -->
      <div v-if="!isEdit" class="character-selection">
        <h3>选择要设置的角色</h3>
        <div class="character-grid">
          <div
            v-for="character in availableCharacters"
            :key="character.id"
            class="character-card"
            :class="{ active: selectedCharacterId === character.id }"
            @click="selectCharacter(character.id)"
          >
            <img :src="character.avatar" :alt="character.name" class="character-avatar" />
            <span class="character-name">{{ character.name }}</span>
          </div>
        </div>
      </div>

      <!-- 编辑表单 -->
      <div v-if="selectedCharacterId || isEdit" class="editor-form">
        <el-tabs v-model="activeTab" type="border-card">
          <!-- 基础信息 -->
          <el-tab-pane label="基础信息" name="basic">
            <div class="form-section">
              <h4>角色基础信息</h4>
              <el-form :model="characterForm" :rules="formRules" ref="characterFormRef" label-width="100px">
                <el-form-item label="角色名称" prop="name">
                  <el-input v-model="characterForm.name" placeholder="请输入角色名称" />
                </el-form-item>
                
                <el-form-item label="角色描述" prop="description">
                  <el-input
                    v-model="characterForm.description"
                    type="textarea"
                    :rows="4"
                    placeholder="请输入角色描述"
                  />
                </el-form-item>
                
                <el-form-item label="角色标签" prop="tags">
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
                    v-if="inputVisible"
                    ref="inputRef"
                    v-model="inputValue"
                    class="tag-input"
                    size="small"
                    @keyup.enter="handleInputConfirm"
                    @blur="handleInputConfirm"
                  />
                  <el-button v-else class="button-new-tag" size="small" @click="showInput">
                    + 添加标签
                  </el-button>
                </el-form-item>
              </el-form>
            </div>
          </el-tab-pane>

          <!-- 性格特征 -->
          <el-tab-pane label="性格特征" name="personality">
            <div class="form-section">
              <h4>性格设置</h4>
              <el-form :model="characterForm.personality" label-width="120px">
                <el-form-item label="外向性">
                  <el-slider v-model="characterForm.personality.extroversion" :min="0" :max="100" />
                </el-form-item>
                <el-form-item label="开放性">
                  <el-slider v-model="characterForm.personality.openness" :min="0" :max="100" />
                </el-form-item>
                <el-form-item label="情绪稳定性">
                  <el-slider v-model="characterForm.personality.emotional_stability" :min="0" :max="100" />
                </el-form-item>
                <el-form-item label="亲和性">
                  <el-slider v-model="characterForm.personality.agreeableness" :min="0" :max="100" />
                </el-form-item>
                <el-form-item label="责任心">
                  <el-slider v-model="characterForm.personality.conscientiousness" :min="0" :max="100" />
                </el-form-item>
              </el-form>
            </div>
          </el-tab-pane>

          <!-- 对话风格 -->
          <el-tab-pane label="对话风格" name="dialogue">
            <div class="form-section">
              <h4>对话风格设置</h4>
              <el-form label-width="120px">
                <el-form-item label="语言风格">
                  <el-select v-model="characterForm.dialogue_style.language_style" placeholder="选择语言风格">
                    <el-option label="正式" value="formal" />
                    <el-option label="随意" value="casual" />
                    <el-option label="幽默" value="humorous" />
                    <el-option label="严肃" value="serious" />
                  </el-select>
                </el-form-item>
                <el-form-item label="回复长度">
                  <el-radio-group v-model="characterForm.dialogue_style.response_length">
                    <el-radio label="short">简短</el-radio>
                    <el-radio label="medium">中等</el-radio>
                    <el-radio label="long">详细</el-radio>
                  </el-radio-group>
                </el-form-item>
              </el-form>
            </div>
          </el-tab-pane>

          <!-- 自定义Prompt -->
          <el-tab-pane label="自定义Prompt" name="prompt">
            <div class="form-section">
              <h4>角色提示词</h4>
              <el-input
                v-model="characterForm.prompt"
                type="textarea"
                :rows="10"
                placeholder="请输入角色的详细提示词，这将决定AI如何扮演这个角色..."
              />
            </div>
          </el-tab-pane>

          <!-- 语音设置 -->
          <el-tab-pane label="语音设置" name="voice">
            <div class="form-section">
              <h4>语音配置</h4>
              <el-form :model="characterForm.voice_settings" label-width="120px">
                <el-form-item label="语音类型">
                  <el-select v-model="characterForm.voice_settings.voice_type" placeholder="选择语音类型">
                    <el-option label="男性" value="male" />
                    <el-option label="女性" value="female" />
                    <el-option label="中性" value="neutral" />
                  </el-select>
                </el-form-item>
                <el-form-item label="语速">
                  <el-slider v-model="characterForm.voice_settings.speed" :min="0.5" :max="2" :step="0.1" />
                </el-form-item>
                <el-form-item label="音调">
                  <el-slider v-model="characterForm.voice_settings.pitch" :min="0.5" :max="2" :step="0.1" />
                </el-form-item>
              </el-form>
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>
    </div>

    <!-- 操作按钮 -->
    <div class="editor-actions">
      <el-button @click="handleCancel">重置</el-button>
      <el-button @click="handlePreview">预览</el-button>
      <el-button type="primary" @click="handleSave" :loading="saving">
        {{ isEdit ? '保存设置' : '保存设置' }}
      </el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouter } from 'vue-router'
import { characterApiService } from '../services/character'

const router = useRouter()
const props = defineProps({
  characterId: {
    type: [String, Number],
    default: null
  },
  isEdit: {
    type: Boolean,
    default: false
  }
})

// 响应式数据
const activeTab = ref('basic')
const selectedCharacterId = ref(props.characterId)
const saving = ref(false)
const inputVisible = ref(false)
const inputValue = ref('')
const inputRef = ref()
const characterFormRef = ref()

// 可选角色列表
const availableCharacters = ref([
  { id: 1, name: '哈利·波特', avatar: '/images/avatars/harry-potter.jpg' },
  { id: 2, name: '苏格拉底', avatar: '/images/avatars/socrates.jpg' },
  { id: 3, name: '莎士比亚', avatar: '/images/avatars/shakespeare.jpg' },
  { id: 4, name: '爱因斯坦', avatar: '/images/avatars/einstein.jpg' },
  { id: 5, name: '夏洛克·福尔摩斯', avatar: '/images/avatars/sherlock.jpg' },
  { id: 6, name: '赫敏·格兰杰', avatar: '/images/avatars/hermione.jpg' }
])

// 表单数据
const characterForm = reactive({
  name: '',
  description: '',
  tags: [],
  prompt: '',
  personality: {
    extroversion: 50,
    openness: 50,
    emotional_stability: 50,
    agreeableness: 50,
    conscientiousness: 50
  },
  dialogue_style: {
    language_style: 'casual',
    response_length: 'medium'
  },
  voice_settings: {
    voice_type: 'neutral',
    speed: 1.0,
    pitch: 1.0
  },
  is_public: false
})

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入角色名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  description: [
    { required: true, message: '请输入角色描述', trigger: 'blur' },
    { min: 10, max: 500, message: '长度在 10 到 500 个字符', trigger: 'blur' }
  ]
}

// 方法
const selectCharacter = (characterId) => {
  selectedCharacterId.value = characterId
  // 根据选择的角色初始化表单
  const character = availableCharacters.value.find(c => c.id === characterId)
  if (character) {
    characterForm.name = character.name
    characterForm.description = `与${character.name}的对话`
  }
}

const removeTag = (tag) => {
  characterForm.tags.splice(characterForm.tags.indexOf(tag), 1)
}

const showInput = () => {
  inputVisible.value = true
  nextTick(() => {
    inputRef.value.input.focus()
  })
}

const handleInputConfirm = () => {
  if (inputValue.value && !characterForm.tags.includes(inputValue.value)) {
    characterForm.tags.push(inputValue.value)
  }
  inputVisible.value = false
  inputValue.value = ''
}

const handleCancel = () => {
  ElMessageBox.confirm('确定要重置所有设置吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    // 重置表单
    Object.assign(characterForm, {
      name: '',
      description: '',
      tags: [],
      prompt: '',
      personality: {
        extroversion: 50,
        openness: 50,
        emotional_stability: 50,
        agreeableness: 50,
        conscientiousness: 50
      },
      dialogue_style: {
        language_style: 'casual',
        response_length: 'medium'
      },
      voice_settings: {
        voice_type: 'neutral',
        speed: 1.0,
        pitch: 1.0
      }
    })
    ElMessage.success('已重置')
  }).catch(() => {})
}

const handlePreview = () => {
  // 预览功能
  ElMessage.info('预览功能开发中...')
}

const handleSave = async () => {
  try {
    // 验证表单
    await characterFormRef.value.validate()
    
    saving.value = true
    
    const requestData = {
      ...characterForm,
      category_id: 1, // 默认分类
      avatar: availableCharacters.value.find(c => c.id === selectedCharacterId.value)?.avatar || ''
    }
    
    let response
    if (props.isEdit && props.characterId) {
      // 更新角色
      response = await characterApiService.updateCharacter(props.characterId, requestData)
    } else {
      // 创建角色
      response = await characterApiService.createCharacter(requestData)
    }
    
    if (response.code === 0) {
      ElMessage.success(props.isEdit ? '角色更新成功！' : '角色创建成功！')
      // 跳转到角色详情或列表页
      router.push('/characters')
    } else {
      ElMessage.error(response.msg || '操作失败')
    }
    
  } catch (error) {
    console.error('保存角色失败:', error)
    ElMessage.error('保存失败: ' + error.message)
  } finally {
    saving.value = false
  }
}

// 加载角色数据（编辑模式）
const loadCharacterData = async () => {
  if (props.isEdit && props.characterId) {
    try {
      const response = await characterApiService.getCharacterDetail(props.characterId)
      if (response.code === 0) {
        Object.assign(characterForm, response.character)
      }
    } catch (error) {
      console.error('加载角色数据失败:', error)
      ElMessage.error('加载角色数据失败')
    }
  }
}

// 生命周期
onMounted(() => {
  loadCharacterData()
})
</script>

<style lang="scss" scoped>
.character-editor {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;

  .editor-header {
    text-align: center;
    margin-bottom: 30px;

    h2 {
      color: #2c3e50;
      margin: 0;
    }
  }

  .character-selection {
    margin-bottom: 30px;

    h3 {
      margin-bottom: 20px;
      color: #2c3e50;
    }

    .character-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
      gap: 20px;
      margin-bottom: 30px;

      .character-card {
        display: flex;
        flex-direction: column;
        align-items: center;
        padding: 20px;
        border: 2px solid #e9ecef;
        border-radius: 12px;
        cursor: pointer;
        transition: all 0.3s ease;

        &:hover {
          border-color: #4A90E2;
          transform: translateY(-2px);
        }

        &.active {
          border-color: #4A90E2;
          background: #f0f8ff;
        }

        .character-avatar {
          width: 80px;
          height: 80px;
          border-radius: 50%;
          object-fit: cover;
          margin-bottom: 10px;
        }

        .character-name {
          font-size: 14px;
          font-weight: 500;
          color: #2c3e50;
          text-align: center;
        }
      }
    }
  }

  .editor-form {
    margin-bottom: 30px;

    .form-section {
      h4 {
        margin-bottom: 20px;
        color: #2c3e50;
      }
    }

    .tag-item {
      margin-right: 8px;
      margin-bottom: 8px;
    }

    .tag-input {
      width: 120px;
      margin-right: 8px;
    }

    .button-new-tag {
      border-style: dashed;
    }
  }

  .editor-actions {
    display: flex;
    justify-content: center;
    gap: 16px;
    padding: 20px 0;
    border-top: 1px solid #e9ecef;
  }
}

// 响应式设计
@media (max-width: 768px) {
  .character-editor {
    padding: 10px;

    .character-selection .character-grid {
      grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
      gap: 15px;

      .character-card {
        padding: 15px;

        .character-avatar {
          width: 60px;
          height: 60px;
        }
      }
    }
  }
}
</style> 
<template>
  <div class="search-bar">
    <div class="search-input-section">
      <el-input
        v-model="localSearchQuery"
        placeholder="搜索角色名称、描述或标签..."
        size="large"
        clearable
        @input="handleSearchInput"
        @clear="handleClear"
        class="search-input"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
    </div>

    <div class="filter-section">
      <div class="tags-filter">
        <span class="filter-label">标签筛选：</span>
        <div class="tags-container">
          <el-tag
            v-for="tag in allTags"
            :key="tag"
            :type="isTagSelected(tag) ? 'primary' : 'info'"
            :effect="isTagSelected(tag) ? 'dark' : 'plain'"
            @click="toggleTag(tag)"
            class="tag-filter"
            size="default"
          >
            {{ tag }}
          </el-tag>
        </div>
        <el-button
          v-if="selectedTags.length > 0"
          type="text"
          size="small"
          @click="clearAllTags"
          class="clear-tags-btn"
        >
          清除所有
        </el-button>
      </div>

      <div class="sort-section">
        <span class="filter-label">排序方式：</span>
        <el-select
          v-model="sortBy"
          placeholder="选择排序方式"
          size="default"
          @change="handleSortChange"
          class="sort-select"
        >
          <el-option label="默认排序" value="default" />
          <el-option label="按名称排序" value="name" />
          <el-option label="按智慧排序" value="intelligence" />
          <el-option label="按友善度排序" value="friendliness" />
          <el-option label="最近使用" value="recent" />
        </el-select>
      </div>

      <div class="view-toggle">
        <el-button-group>
          <el-button
            :type="viewMode === 'grid' ? 'primary' : 'default'"
            @click="setViewMode('grid')"
            size="default"
          >
            <el-icon><Grid /></el-icon>
          </el-button>
          <el-button
            :type="viewMode === 'list' ? 'primary' : 'default'"
            @click="setViewMode('list')"
            size="default"
          >
            <el-icon><List /></el-icon>
          </el-button>
        </el-button-group>
      </div>
    </div>

    <div v-if="hasActiveFilters" class="active-filters">
      <span class="filters-label">当前筛选：</span>
      <el-tag
        v-if="localSearchQuery"
        closable
        @close="clearSearch"
        type="warning"
        class="active-filter"
      >
        搜索: {{ localSearchQuery }}
      </el-tag>
      <el-tag
        v-for="tag in selectedTags"
        :key="tag"
        closable
        @close="removeTag(tag)"
        type="primary"
        class="active-filter"
      >
        {{ tag }}
      </el-tag>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useCharacterStore } from '../stores/character'
import { Search, Grid, List } from '@element-plus/icons-vue'

const emit = defineEmits(['viewModeChange', 'sortChange'])

const characterStore = useCharacterStore()

// 本地搜索查询
const localSearchQuery = ref(characterStore.searchQuery)
// 排序方式
const sortBy = ref('default')
// 视图模式
const viewMode = ref('grid')

// 计算属性
const allTags = computed(() => characterStore.allTags)
const selectedTags = computed(() => characterStore.selectedTags)

const hasActiveFilters = computed(() => {
  return localSearchQuery.value || selectedTags.value.length > 0
})

// 监听器
watch(localSearchQuery, (newValue) => {
  characterStore.setSearchQuery(newValue)
})

// 方法
const handleSearchInput = (value) => {
  localSearchQuery.value = value
}

const handleClear = () => {
  localSearchQuery.value = ''
  characterStore.setSearchQuery('')
}

const clearSearch = () => {
  localSearchQuery.value = ''
  characterStore.setSearchQuery('')
}

const isTagSelected = (tag) => {
  return selectedTags.value.includes(tag)
}

const toggleTag = (tag) => {
  const currentTags = [...selectedTags.value]
  const index = currentTags.indexOf(tag)
  
  if (index > -1) {
    currentTags.splice(index, 1)
  } else {
    currentTags.push(tag)
  }
  
  characterStore.setSelectedTags(currentTags)
}

const removeTag = (tag) => {
  const currentTags = selectedTags.value.filter(t => t !== tag)
  characterStore.setSelectedTags(currentTags)
}

const clearAllTags = () => {
  characterStore.setSelectedTags([])
}

const handleSortChange = (value) => {
  emit('sortChange', value)
}

const setViewMode = (mode) => {
  viewMode.value = mode
  emit('viewModeChange', mode)
}
</script>

<style lang="scss" scoped>
.search-bar {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  padding: 24px;
  margin-bottom: 24px;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.search-input-section {
  margin-bottom: 20px;
  
  .search-input {
    .el-input__wrapper {
      border-radius: 12px;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    }
  }
}

.filter-section {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  align-items: center;
  margin-bottom: 16px;
}

.filter-label {
  font-weight: 500;
  color: #606266;
  font-size: 14px;
  white-space: nowrap;
}

.tags-filter {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  flex: 1;
  min-width: 300px;
  
  .tags-container {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
  }
  
  .tag-filter {
    cursor: pointer;
    transition: all 0.2s ease;
    border-radius: 12px;
    
    &:hover {
      transform: scale(1.05);
    }
  }
  
  .clear-tags-btn {
    margin-left: 8px;
    color: #909399;
    
    &:hover {
      color: #F56C6C;
    }
  }
}

.sort-section {
  display: flex;
  align-items: center;
  gap: 8px;
  
  .sort-select {
    width: 120px;
  }
}

.view-toggle {
  .el-button-group .el-button {
    border-radius: 8px;
    
    &:first-child {
      border-top-right-radius: 0;
      border-bottom-right-radius: 0;
    }
    
    &:last-child {
      border-top-left-radius: 0;
      border-bottom-left-radius: 0;
    }
  }
}

.active-filters {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
  padding-top: 16px;
  border-top: 1px solid rgba(0, 0, 0, 0.06);
}

.filters-label {
  font-size: 13px;
  color: #909399;
  white-space: nowrap;
}

.active-filter {
  border-radius: 12px;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .filter-section {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
  
  .tags-filter {
    min-width: 100%;
    flex-direction: column;
    align-items: flex-start;
    
    .tags-container {
      width: 100%;
    }
  }
  
  .sort-section,
  .view-toggle {
    width: 100%;
    justify-content: flex-start;
  }
}

@media (max-width: 768px) {
  .search-bar {
    padding: 16px;
  }
  
  .filter-section {
    gap: 12px;
  }
  
  .tags-filter .tags-container {
    justify-content: flex-start;
  }
  
  .sort-section .sort-select {
    width: 100px;
  }
  
  .active-filters {
    gap: 6px;
  }
}
</style>

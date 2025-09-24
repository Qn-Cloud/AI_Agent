<template>
  <div id="app" class="app-container">
    <el-container class="layout-container">
      <!-- 左侧导航栏 -->
      <el-aside class="sidebar" width="220px">
        <div class="sidebar-content">
          <!-- Logo -->
          <div class="logo-section">
            <div class="logo">LOGO</div>
          </div>

          <!-- 导航菜单 -->
          <el-menu
            :default-active="$route.path"
            class="sidebar-menu"
            @select="handleMenuSelect"
          >
            <el-menu-item index="/" class="menu-item">
              <span>发现</span>
            </el-menu-item>
            <el-menu-item index="/history" class="menu-item">
              <span>历史</span>
            </el-menu-item>
          </el-menu>

          <!-- 历史记录分组 -->
          <div class="history-section">
            <div class="section-title">今天</div>
            <div class="history-item">
              <div class="item-icon"></div>
              <span class="item-name">角色名称</span>
            </div>
            <div class="history-item">
              <div class="item-icon"></div>
              <span class="item-name">角色名称</span>
            </div>

            <div class="section-title">昨天</div>
            <div class="history-item">
              <div class="item-icon"></div>
              <span class="item-name">角色名称</span>
            </div>
            <div class="history-item">
              <div class="item-icon"></div>
              <span class="item-name">角色名称</span>
            </div>

            <div class="section-title">更久之前</div>
            <div class="history-item">
              <div class="item-icon"></div>
              <span class="item-name">角色名称</span>
            </div>
            <div class="history-item">
              <div class="item-icon"></div>
              <span class="item-name">角色名称</span>
            </div>
          </div>
        </div>
      </el-aside>

      <!-- 主要内容区域 -->
      <el-main class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const handleMenuSelect = (index) => {
  router.push(index)
}
</script>

<style lang="scss">
.app-container {
  height: 100vh;
  background: #f5f5f5;
}

.layout-container {
  height: 100vh;
}

.sidebar {
  background: white;
  border-right: 1px solid #e5e5e5;
  height: 100vh;
  
  .sidebar-content {
    padding: 24px 16px;
    height: 100%;
    display: flex;
    flex-direction: column;
  }
  
  .logo-section {
    margin-bottom: 40px;
    
    .logo {
      font-size: 18px;
      font-weight: bold;
      color: #333;
      padding: 8px 12px;
    }
  }
  
  .sidebar-menu {
    border: none;
    background: transparent;
    margin-bottom: 32px;
    
    .menu-item {
      height: 44px;
      line-height: 44px;
      border-radius: 8px;
      margin-bottom: 8px;
      padding: 0 12px;
      
      &:hover {
        background: #f5f5f5;
      }
      
      &.is-active {
        background: #f0f0f0;
        font-weight: 500;
      }
      
      span {
        font-size: 15px;
        color: #333;
      }
    }
  }
  
  .history-section {
    flex: 1;
    overflow-y: auto;
    
    .section-title {
      font-size: 13px;
      color: #666;
      margin: 20px 0 12px 0;
      font-weight: 500;
      
      &:first-child {
        margin-top: 0;
      }
    }
    
    .history-item {
      display: flex;
      align-items: center;
      padding: 8px 12px;
      border-radius: 6px;
      cursor: pointer;
      margin-bottom: 4px;
      
      &:hover {
        background: #f5f5f5;
      }
      
      .item-icon {
        width: 16px;
        height: 16px;
        background: #4A90E2;
        border-radius: 4px;
        margin-right: 12px;
        flex-shrink: 0;
      }
      
      .item-name {
        font-size: 14px;
        color: #666;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }
  }
}

.main-content {
  background: #f5f5f5;
  padding: 0;
  height: 100vh;
  overflow-y: auto;
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    width: 180px !important;
    
    .sidebar-content {
      padding: 20px 12px;
    }
  }
}
</style>

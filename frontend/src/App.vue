<template>
  <div id="app" class="app-container">
    <el-container class="layout-container">
      <!-- 左侧导航栏 -->
      <el-aside class="sidebar" width="280px">
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
            <el-menu-item index="/api-test" class="menu-item">
              <span>API测试</span>
            </el-menu-item>
          </el-menu>

          <!-- 动态历史记录分组 -->
          <ChatSidebar />
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
import ChatSidebar from './components/ChatSidebar.vue'

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
  overflow: hidden;
  
  .sidebar-content {
    height: 100%;
    display: flex;
    flex-direction: column;
  }
  
  .logo-section {
    padding: 20px;
    text-align: center;
    border-bottom: 1px solid #f0f0f0;
    
    .logo {
      font-size: 24px;
      font-weight: bold;
      color: #333;
    }
  }
  
  .sidebar-menu {
    border: none;
    
    .menu-item {
      margin: 8px 12px;
      border-radius: 8px;
      
      &.is-active {
        background-color: #f0f8ff;
        color: #4A90E2;
      }
      
      &:hover {
        background-color: #f8f9fa;
      }
      
      span {
        font-size: 14px;
        font-weight: 500;
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

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'

import App from './App.vue'
import router from './router'
import config from './config'

// 导入全局组件
import ErrorBoundary from './components/ErrorBoundary.vue'
import LoadingStates from './components/LoadingStates.vue'
import NotificationCenter from './components/NotificationCenter.vue'

// 导入工具和指令
import { lazyLoad, performanceMonitor } from './utils/performance'

// 创建应用实例
const app = createApp(App)

// 使用插件
app.use(createPinia())
app.use(router)
app.use(ElementPlus)

// 注册全局组件
app.component('ErrorBoundary', ErrorBoundary)
app.component('LoadingStates', LoadingStates)
app.component('NotificationCenter', NotificationCenter)

// 注册全局指令
app.directive('lazy', lazyLoad)

// 全局属性
app.config.globalProperties.$config = config

// 全局错误处理
app.config.errorHandler = (error, instance, info) => {
  console.error('Global error:', error, info)
  
  // 在生产环境中，可以将错误发送到错误监控服务
  if (config.isProduction) {
    // 发送错误报告到监控服务
    console.log('Sending error report to monitoring service...')
  }
}

// 性能监控
if (config.isDevelopment) {
  performanceMonitor.start('app-init')
  
  app.config.performance = true
  
  // 监听路由变化的性能
  router.beforeEach((to, from, next) => {
    performanceMonitor.start(`route-${to.name}`)
    next()
  })
  
  router.afterEach((to) => {
    performanceMonitor.end(`route-${to.name}`)
  })
}

// 应用挂载
const mountApp = () => {
  app.mount('#app')
  
  if (config.isDevelopment) {
    performanceMonitor.end('app-init')
    console.log('🚀 App initialized successfully!')
    console.log('📊 Performance metrics:', performanceMonitor.getAllMetrics())
  }
}

// 检查浏览器兼容性
const checkBrowserSupport = () => {
  const requiredFeatures = [
    'Promise',
    'fetch',
    'localStorage',
    'sessionStorage'
  ]
  
  const unsupported = requiredFeatures.filter(feature => !(feature in window))
  
  if (unsupported.length > 0) {
    console.error('Unsupported browser features:', unsupported)
    alert('您的浏览器版本过低，请升级到最新版本以获得最佳体验。')
    return false
  }
  
  return true
}

// 初始化应用
if (checkBrowserSupport()) {
  mountApp()
} else {
  document.getElementById('app').innerHTML = `
    <div style="text-align: center; padding: 50px; font-family: Arial, sans-serif;">
      <h2>浏览器不兼容</h2>
      <p>请使用现代浏览器访问本应用，推荐使用 Chrome、Firefox、Safari 或 Edge 的最新版本。</p>
    </div>
  `
}

// 导出应用实例（用于测试）
export default app
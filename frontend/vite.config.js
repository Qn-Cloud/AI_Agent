import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
      'vue': 'vue/dist/vue.esm-bundler.js'
    }
  },
  server: {
    port: 3000,
    open: true,
    host: '0.0.0.0',
    proxy: {
      // 代理所有以 /api 开头的请求
      '/api': {
        target: 'http://192.168.23.188:7001', // 聊天服务
        changeOrigin: true,
        secure: false,
        rewrite: (path) => {
          console.log('🔄 代理请求:', path)
          return path
        }
      },
      // 如果需要代理到不同的服务，可以添加更多规则
      '/api/character': {
        target: 'http://192.168.23.188:7002', // 角色服务
        changeOrigin: true,
        secure: false,
        rewrite: (path) => {
          console.log('🔄 代理角色请求:', path)
          return path
        }
      }
    }
  },
  build: {
    outDir: 'dist',
    sourcemap: false,
    chunkSizeWarningLimit: 1500
  }
})

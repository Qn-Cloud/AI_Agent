import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src')
    }
  },
  server: {
    port: 3000,
    open: true,
    host: '0.0.0.0',
    proxy: {
      // 注意：更具体的路径规则要放在前面，通用规则放在后面
      '/api/speech': {
        target: 'http://192.168.23.188:7005', // 语音服务
        changeOrigin: true,
        secure: false,
        rewrite: (path) => {
          console.log('🔄 代理语音请求:', path)
          return path
        }
      },
      '/api/character': {
        target: 'http://192.168.23.188:7002', // 角色服务
        changeOrigin: true,
        secure: false,
        rewrite: (path) => {
          console.log('🔄 代理角色请求:', path)
          return path
        }
      },
      // 代理其他以 /api 开头的请求到聊天服务
      '/api': {
        target: 'http://192.168.23.188:7001', // 聊天服务
        changeOrigin: true,
        secure: false,
        rewrite: (path) => {
          console.log('🔄 代理请求:', path)
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

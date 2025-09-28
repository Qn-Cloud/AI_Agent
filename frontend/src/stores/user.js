import { defineStore } from 'pinia'
import { userApiService as userApi } from '../services'

export const useUserStore = defineStore('user', {
  state: () => ({
    // 用户信息
    userInfo: null,
    // 认证状态
    isAuthenticated: false,
    // 认证token
    token: null,
    refreshToken: null,
    // 加载状态
    loading: false,
    error: null
  }),

  getters: {
    // 是否已登录
    isLoggedIn: (state) => {
      return state.isAuthenticated && state.token && state.userInfo
    },

    // 用户名
    username: (state) => {
      return state.userInfo?.username || ''
    },

    // 用户头像
    avatar: (state) => {
      return state.userInfo?.avatar || ''
    },

    // 用户昵称
    nickname: (state) => {
      return state.userInfo?.nickname || state.userInfo?.username || ''
    }
  },

  actions: {
    // 初始化用户状态
    initializeAuth() {
      const token = localStorage.getItem('auth_token')
      const refreshToken = localStorage.getItem('refresh_token')
      const userInfo = localStorage.getItem('user_info')

      if (token && userInfo) {
        try {
          this.token = token
          this.refreshToken = refreshToken
          this.userInfo = JSON.parse(userInfo)
          this.isAuthenticated = true
        } catch (error) {
          console.error('解析用户信息失败:', error)
          this.clearAuth()
        }
      }
    },

    // 用户登录
    async login(credentials) {
      try {
        this.loading = true
        this.error = null

        const response = await userApi.login(credentials)
        
        if (response.data) {
          this.setAuthData(response.data, response.token, response.refresh_token)
          return response.data
        }
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    // 用户注册
    async register(userData) {
      try {
        this.loading = true
        this.error = null

        const response = await userApi.register(userData)
        
        if (response.data) {
          this.setAuthData(response.data, response.token)
          return response.data
        }
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    // 用户登出
    async logout() {
      try {
        // 调用后端登出接口
        await userApi.logout()
      } catch (error) {
        console.error('登出接口调用失败:', error)
      } finally {
        this.clearAuth()
      }
    },

    // 获取用户信息
    async fetchUserInfo() {
      try {
        this.loading = true
        const response = await userApi.getUserInfo()
        
        if (response.data) {
          this.userInfo = response.data
          this.updateLocalStorage()
        }
      } catch (error) {
        console.error('获取用户信息失败:', error)
        if (error.response?.status === 401) {
          this.clearAuth()
        }
      } finally {
        this.loading = false
      }
    },

    // 更新用户信息
    async updateUserInfo(userData) {
      try {
        this.loading = true
        this.error = null

        await userApi.updateUserInfo(userData)
        
        // 更新本地用户信息
        this.userInfo = { ...this.userInfo, ...userData }
        this.updateLocalStorage()
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    // 修改密码
    async changePassword(passwordData) {
      try {
        this.loading = true
        this.error = null

        await userApi.changePassword(passwordData)
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    // 上传头像
    async uploadAvatar(file) {
      try {
        this.loading = true
        this.error = null

        const response = await userApi.uploadAvatar(file)
        
        if (response.data && response.data.avatar_url) {
          this.userInfo.avatar = response.data.avatar_url
          this.updateLocalStorage()
          return response.data.avatar_url
        }
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    // 刷新token
    async refreshAuthToken() {
      try {
        if (!this.refreshToken) {
          throw new Error('没有refresh token')
        }

        const response = await userApi.refreshToken(this.refreshToken)
        
        if (response.token) {
          this.token = response.token
          if (response.refresh_token) {
            this.refreshToken = response.refresh_token
          }
          this.updateLocalStorage()
          return response.token
        }
      } catch (error) {
        console.error('刷新token失败:', error)
        this.clearAuth()
        throw error
      }
    },

    // 设置认证数据
    setAuthData(userInfo, token, refreshToken = null) {
      this.userInfo = userInfo
      this.token = token
      this.refreshToken = refreshToken
      this.isAuthenticated = true
      this.updateLocalStorage()
    },

    // 更新本地存储
    updateLocalStorage() {
      if (this.token) {
        localStorage.setItem('auth_token', this.token)
      }
      if (this.refreshToken) {
        localStorage.setItem('refresh_token', this.refreshToken)
      }
      if (this.userInfo) {
        localStorage.setItem('user_info', JSON.stringify(this.userInfo))
      }
    },

    // 清除认证信息
    clearAuth() {
      this.userInfo = null
      this.token = null
      this.refreshToken = null
      this.isAuthenticated = false
      this.error = null

      // 清除本地存储
      localStorage.removeItem('auth_token')
      localStorage.removeItem('refresh_token')
      localStorage.removeItem('user_info')
    },

    // 检查用户名可用性
    async checkUsername(username) {
      try {
        const response = await userApi.checkUsername(username)
        return response.data?.available || false
      } catch (error) {
        console.error('检查用户名失败:', error)
        return false
      }
    },

    // 检查邮箱可用性
    async checkEmail(email) {
      try {
        const response = await userApi.checkEmail(email)
        return response.data?.available || false
      } catch (error) {
        console.error('检查邮箱失败:', error)
        return false
      }
    }
  }
})

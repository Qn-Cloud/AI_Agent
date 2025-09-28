import { userApi } from './apiFactory'

export const userApiService = {
  // 用户注册
  register(data) {
    return userApi.post('/api/user/register', {
      username: data.username,
      email: data.email,
      password: data.password
    })
  },

  // 用户登录
  login(data) {
    return userApi.post('/api/user/login', {
      username: data.username,
      password: data.password,
      remember: data.remember || false
    })
  },

  // 刷新token
  refreshToken(refreshToken) {
    return userApi.post('/api/user/refresh', {
      refresh_token: refreshToken
    })
  },

  // 用户退出
  logout() {
    return userApi.post('/api/user/logout')
  },

  // 获取用户信息
  getUserInfo() {
    return userApi.get('/api/user/info')
  },

  // 更新用户信息
  updateUserInfo(data) {
    return userApi.put('/api/user/info', data)
  },

  // 修改密码
  changePassword(data) {
    return userApi.put('/api/user/password', {
      old_password: data.oldPassword,
      new_password: data.newPassword
    })
  },

  // 上传头像
  uploadAvatar(file) {
    const formData = new FormData()
    formData.append('avatar', file)
    return userApi.post('/api/user/avatar', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  // 检查用户名是否可用
  checkUsername(username) {
    return userApi.get('/api/user/check-username', {
      params: { username }
    })
  },

  // 检查邮箱是否可用
  checkEmail(email) {
    return userApi.get('/api/user/check-email', {
      params: { email }
    })
  }
}

// 保持向后兼容
export const userApi_old = userApiService
export default userApiService

import api from './api'

export const userApi = {
  // 用户注册
  register(data) {
    return api.post('/api/user/register', {
      username: data.username,
      email: data.email,
      password: data.password
    })
  },

  // 用户登录
  login(data) {
    return api.post('/api/user/login', {
      username: data.username,
      password: data.password,
      remember: data.remember || false
    })
  },

  // 刷新token
  refreshToken(refreshToken) {
    return api.post('/api/user/refresh', {
      refresh_token: refreshToken
    })
  },

  // 用户退出
  logout() {
    return api.post('/api/user/logout')
  },

  // 获取用户信息
  getUserInfo() {
    return api.get('/api/user/info')
  },

  // 更新用户信息
  updateUserInfo(data) {
    return api.put('/api/user/info', data)
  },

  // 修改密码
  changePassword(data) {
    return api.put('/api/user/password', {
      old_password: data.oldPassword,
      new_password: data.newPassword
    })
  },

  // 上传头像
  uploadAvatar(file) {
    const formData = new FormData()
    formData.append('avatar', file)
    return api.post('/api/user/avatar', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  // 检查用户名是否可用
  checkUsername(username) {
    return api.get('/api/user/check-username', {
      params: { username }
    })
  },

  // 检查邮箱是否可用
  checkEmail(email) {
    return api.get('/api/user/check-email', {
      params: { email }
    })
  }
}

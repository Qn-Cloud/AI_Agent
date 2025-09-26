// 数据验证工具类

export class Validator {
  constructor() {
    this.errors = {}
  }

  // 清除所有错误
  clearErrors() {
    this.errors = {}
    return this
  }

  // 清除特定字段错误
  clearError(field) {
    delete this.errors[field]
    return this
  }

  // 添加错误
  addError(field, message) {
    if (!this.errors[field]) {
      this.errors[field] = []
    }
    this.errors[field].push(message)
    return this
  }

  // 获取错误
  getErrors() {
    return this.errors
  }

  // 获取特定字段错误
  getError(field) {
    return this.errors[field] || []
  }

  // 检查是否有错误
  hasErrors() {
    return Object.keys(this.errors).length > 0
  }

  // 检查特定字段是否有错误
  hasError(field) {
    return this.errors[field] && this.errors[field].length > 0
  }

  // 获取第一个错误信息
  getFirstError(field) {
    const errors = this.getError(field)
    return errors.length > 0 ? errors[0] : null
  }

  // 验证必填字段
  required(value, field, message = `${field}不能为空`) {
    if (value === null || value === undefined || value === '' || 
        (Array.isArray(value) && value.length === 0)) {
      this.addError(field, message)
    }
    return this
  }

  // 验证字符串长度
  length(value, field, min, max, message = null) {
    if (value && typeof value === 'string') {
      if (min && value.length < min) {
        this.addError(field, message || `${field}长度不能少于${min}个字符`)
      }
      if (max && value.length > max) {
        this.addError(field, message || `${field}长度不能超过${max}个字符`)
      }
    }
    return this
  }

  // 验证邮箱格式
  email(value, field, message = '请输入有效的邮箱地址') {
    if (value) {
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
      if (!emailRegex.test(value)) {
        this.addError(field, message)
      }
    }
    return this
  }

  // 验证手机号格式
  phone(value, field, message = '请输入有效的手机号码') {
    if (value) {
      const phoneRegex = /^1[3-9]\d{9}$/
      if (!phoneRegex.test(value)) {
        this.addError(field, message)
      }
    }
    return this
  }

  // 验证密码强度
  password(value, field, options = {}) {
    if (value) {
      const {
        minLength = 6,
        requireUppercase = false,
        requireLowercase = false,
        requireNumbers = false,
        requireSpecialChars = false
      } = options

      if (value.length < minLength) {
        this.addError(field, `密码长度不能少于${minLength}个字符`)
      }

      if (requireUppercase && !/[A-Z]/.test(value)) {
        this.addError(field, '密码必须包含大写字母')
      }

      if (requireLowercase && !/[a-z]/.test(value)) {
        this.addError(field, '密码必须包含小写字母')
      }

      if (requireNumbers && !/\d/.test(value)) {
        this.addError(field, '密码必须包含数字')
      }

      if (requireSpecialChars && !/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(value)) {
        this.addError(field, '密码必须包含特殊字符')
      }
    }
    return this
  }

  // 验证数字范围
  range(value, field, min, max, message = null) {
    if (value !== null && value !== undefined) {
      const num = Number(value)
      if (isNaN(num)) {
        this.addError(field, `${field}必须是有效数字`)
      } else {
        if (min !== null && num < min) {
          this.addError(field, message || `${field}不能小于${min}`)
        }
        if (max !== null && num > max) {
          this.addError(field, message || `${field}不能大于${max}`)
        }
      }
    }
    return this
  }

  // 验证URL格式
  url(value, field, message = '请输入有效的URL地址') {
    if (value) {
      try {
        new URL(value)
      } catch {
        this.addError(field, message)
      }
    }
    return this
  }

  // 自定义验证规则
  custom(value, field, validator, message) {
    if (!validator(value)) {
      this.addError(field, message)
    }
    return this
  }

  // 验证两个字段是否相等（如确认密码）
  equals(value1, value2, field, message = '两次输入不一致') {
    if (value1 !== value2) {
      this.addError(field, message)
    }
    return this
  }

  // 验证数组长度
  arrayLength(value, field, min, max, message = null) {
    if (Array.isArray(value)) {
      if (min && value.length < min) {
        this.addError(field, message || `${field}至少需要${min}个项目`)
      }
      if (max && value.length > max) {
        this.addError(field, message || `${field}最多只能有${max}个项目`)
      }
    }
    return this
  }

  // 验证文件类型
  fileType(file, field, allowedTypes, message = null) {
    if (file) {
      const fileType = file.type || ''
      const fileName = file.name || ''
      const fileExt = fileName.split('.').pop()?.toLowerCase()
      
      const isAllowed = allowedTypes.some(type => {
        if (type.includes('/')) {
          // MIME type check
          return fileType === type || fileType.startsWith(type.split('/')[0] + '/')
        } else {
          // File extension check
          return fileExt === type.toLowerCase()
        }
      })

      if (!isAllowed) {
        this.addError(field, message || `文件类型不支持，仅支持：${allowedTypes.join(', ')}`)
      }
    }
    return this
  }

  // 验证文件大小
  fileSize(file, field, maxSize, message = null) {
    if (file && file.size > maxSize) {
      const maxSizeMB = (maxSize / 1024 / 1024).toFixed(1)
      this.addError(field, message || `文件大小不能超过${maxSizeMB}MB`)
    }
    return this
  }
}

// 常用验证规则
export const ValidationRules = {
  // 用户注册验证
  userRegistration: {
    username: (value) => {
      const validator = new Validator()
      return validator
        .required(value, 'username', '用户名不能为空')
        .length(value, 'username', 3, 20)
        .custom(value, 'username', (v) => /^[a-zA-Z0-9_]+$/.test(v), '用户名只能包含字母、数字和下划线')
    },
    
    email: (value) => {
      const validator = new Validator()
      return validator
        .required(value, 'email', '邮箱不能为空')
        .email(value, 'email')
    },
    
    password: (value) => {
      const validator = new Validator()
      return validator
        .required(value, 'password', '密码不能为空')
        .password(value, 'password', {
          minLength: 6,
          requireNumbers: true
        })
    }
  },

  // 角色设置验证
  characterSettings: {
    name: (value) => {
      const validator = new Validator()
      return validator
        .required(value, 'name', '角色名称不能为空')
        .length(value, 'name', 1, 50)
    },
    
    description: (value) => {
      const validator = new Validator()
      return validator
        .required(value, 'description', '角色描述不能为空')
        .length(value, 'description', 10, 500)
    },
    
    tags: (value) => {
      const validator = new Validator()
      return validator
        .arrayLength(value, 'tags', 1, 10, '请选择1-10个标签')
    },
    
    prompt: (value) => {
      const validator = new Validator()
      return validator
        .required(value, 'prompt', '角色提示词不能为空')
        .length(value, 'prompt', 20, 2000)
    }
  },

  // 文件上传验证
  fileUpload: {
    avatar: (file) => {
      const validator = new Validator()
      return validator
        .required(file, 'avatar', '请选择头像文件')
        .fileType(file, 'avatar', ['image/jpeg', 'image/png', 'image/gif', 'image/webp'])
        .fileSize(file, 'avatar', 5 * 1024 * 1024) // 5MB
    },
    
    audio: (file) => {
      const validator = new Validator()
      return validator
        .required(file, 'audio', '请选择音频文件')
        .fileType(file, 'audio', ['audio/mp3', 'audio/wav', 'audio/ogg', 'audio/webm'])
        .fileSize(file, 'audio', 10 * 1024 * 1024) // 10MB
    }
  }
}

// 表单验证助手
export class FormValidator {
  constructor(rules = {}) {
    this.rules = rules
    this.errors = {}
  }

  // 验证单个字段
  validateField(field, value) {
    const rule = this.rules[field]
    if (rule) {
      const validator = rule(value)
      if (validator.hasError(field)) {
        this.errors[field] = validator.getError(field)
      } else {
        delete this.errors[field]
      }
    }
    return !this.hasError(field)
  }

  // 验证所有字段
  validate(data) {
    this.errors = {}
    
    Object.keys(this.rules).forEach(field => {
      this.validateField(field, data[field])
    })
    
    return !this.hasErrors()
  }

  // 获取错误
  getErrors() {
    return this.errors
  }

  // 获取特定字段错误
  getError(field) {
    return this.errors[field] || []
  }

  // 获取第一个错误
  getFirstError(field) {
    const errors = this.getError(field)
    return errors.length > 0 ? errors[0] : null
  }

  // 检查是否有错误
  hasErrors() {
    return Object.keys(this.errors).length > 0
  }

  // 检查特定字段是否有错误
  hasError(field) {
    return this.errors[field] && this.errors[field].length > 0
  }

  // 清除错误
  clearErrors() {
    this.errors = {}
  }

  // 清除特定字段错误
  clearError(field) {
    delete this.errors[field]
  }
}

// 导出默认验证器
export default Validator

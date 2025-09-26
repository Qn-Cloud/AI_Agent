// 性能优化工具类

export class PerformanceMonitor {
  constructor() {
    this.metrics = new Map()
    this.observers = []
  }

  // 开始性能监控
  start(name) {
    this.metrics.set(name, {
      startTime: performance.now(),
      endTime: null,
      duration: null
    })
  }

  // 结束性能监控
  end(name) {
    const metric = this.metrics.get(name)
    if (metric) {
      metric.endTime = performance.now()
      metric.duration = metric.endTime - metric.startTime
      
      // 通知观察者
      this.notifyObservers(name, metric)
      
      // 开发环境下输出到控制台
      if (import.meta.env.DEV) {
        console.log(`⏱️ ${name}: ${metric.duration.toFixed(2)}ms`)
      }
    }
  }

  // 获取性能指标
  getMetric(name) {
    return this.metrics.get(name)
  }

  // 获取所有性能指标
  getAllMetrics() {
    return Object.fromEntries(this.metrics)
  }

  // 添加观察者
  addObserver(callback) {
    this.observers.push(callback)
  }

  // 通知观察者
  notifyObservers(name, metric) {
    this.observers.forEach(callback => callback(name, metric))
  }

  // 清除指标
  clear() {
    this.metrics.clear()
  }
}

// 图片懒加载指令
export const lazyLoad = {
  mounted(el, binding) {
    const observer = new IntersectionObserver((entries) => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          const img = entry.target
          img.src = binding.value
          img.classList.add('loaded')
          observer.unobserve(img)
        }
      })
    }, {
      threshold: 0.1
    })

    observer.observe(el)
  }
}

// 防抖函数
export function debounce(func, wait) {
  let timeout
  return function executedFunction(...args) {
    const later = () => {
      clearTimeout(timeout)
      func(...args)
    }
    clearTimeout(timeout)
    timeout = setTimeout(later, wait)
  }
}

// 节流函数
export function throttle(func, limit) {
  let inThrottle
  return function() {
    const args = arguments
    const context = this
    if (!inThrottle) {
      func.apply(context, args)
      inThrottle = true
      setTimeout(() => inThrottle = false, limit)
    }
  }
}

// 内存使用监控
export class MemoryMonitor {
  constructor() {
    this.isSupported = 'memory' in performance
  }

  getMemoryInfo() {
    if (!this.isSupported) {
      return null
    }

    const memory = performance.memory
    return {
      usedJSHeapSize: this.formatBytes(memory.usedJSHeapSize),
      totalJSHeapSize: this.formatBytes(memory.totalJSHeapSize),
      jsHeapSizeLimit: this.formatBytes(memory.jsHeapSizeLimit),
      usedPercentage: ((memory.usedJSHeapSize / memory.jsHeapSizeLimit) * 100).toFixed(2)
    }
  }

  formatBytes(bytes) {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }

  startMonitoring(interval = 5000) {
    if (!this.isSupported) {
      console.warn('Memory API is not supported in this browser')
      return
    }

    const monitor = () => {
      const memInfo = this.getMemoryInfo()
      if (import.meta.env.DEV) {
        console.log('🧠 Memory Usage:', memInfo)
      }
      
      // 内存使用率超过80%时发出警告
      if (parseFloat(memInfo.usedPercentage) > 80) {
        console.warn('⚠️ High memory usage detected:', memInfo.usedPercentage + '%')
      }
    }

    monitor() // 立即执行一次
    return setInterval(monitor, interval)
  }
}

// 组件渲染性能监控
export function withPerformanceTracking(component, name) {
  const monitor = new PerformanceMonitor()
  
  return {
    ...component,
    beforeMount() {
      monitor.start(`${name}-mount`)
      if (component.beforeMount) {
        component.beforeMount.call(this)
      }
    },
    mounted() {
      monitor.end(`${name}-mount`)
      if (component.mounted) {
        component.mounted.call(this)
      }
    },
    beforeUpdate() {
      monitor.start(`${name}-update`)
      if (component.beforeUpdate) {
        component.beforeUpdate.call(this)
      }
    },
    updated() {
      monitor.end(`${name}-update`)
      if (component.updated) {
        component.updated.call(this)
      }
    }
  }
}

// 网络请求性能监控
export class NetworkMonitor {
  constructor() {
    this.requests = new Map()
  }

  trackRequest(url, options = {}) {
    const requestId = this.generateId()
    const startTime = performance.now()
    
    this.requests.set(requestId, {
      url,
      method: options.method || 'GET',
      startTime,
      endTime: null,
      duration: null,
      status: null,
      size: null
    })

    return {
      requestId,
      complete: (response) => {
        this.completeRequest(requestId, response)
      },
      error: (error) => {
        this.errorRequest(requestId, error)
      }
    }
  }

  completeRequest(requestId, response) {
    const request = this.requests.get(requestId)
    if (request) {
      request.endTime = performance.now()
      request.duration = request.endTime - request.startTime
      request.status = response.status
      request.size = this.getResponseSize(response)

      if (import.meta.env.DEV) {
        console.log(`🌐 ${request.method} ${request.url}: ${request.duration.toFixed(2)}ms (${request.status})`)
      }
    }
  }

  errorRequest(requestId, error) {
    const request = this.requests.get(requestId)
    if (request) {
      request.endTime = performance.now()
      request.duration = request.endTime - request.startTime
      request.error = error.message

      if (import.meta.env.DEV) {
        console.error(`🌐 ${request.method} ${request.url}: ${request.duration.toFixed(2)}ms (ERROR: ${error.message})`)
      }
    }
  }

  getResponseSize(response) {
    const contentLength = response.headers.get('content-length')
    return contentLength ? parseInt(contentLength) : null
  }

  generateId() {
    return Date.now().toString(36) + Math.random().toString(36).substr(2)
  }

  getStats() {
    const requests = Array.from(this.requests.values())
    const completed = requests.filter(r => r.endTime && !r.error)
    
    if (completed.length === 0) return null

    const totalDuration = completed.reduce((sum, r) => sum + r.duration, 0)
    const avgDuration = totalDuration / completed.length
    const slowest = Math.max(...completed.map(r => r.duration))
    const fastest = Math.min(...completed.map(r => r.duration))

    return {
      total: requests.length,
      completed: completed.length,
      errors: requests.filter(r => r.error).length,
      avgDuration: avgDuration.toFixed(2),
      slowest: slowest.toFixed(2),
      fastest: fastest.toFixed(2)
    }
  }
}

// 虚拟滚动实现
export class VirtualScroll {
  constructor(options) {
    this.itemHeight = options.itemHeight || 50
    this.containerHeight = options.containerHeight || 300
    this.buffer = options.buffer || 5
    this.items = options.items || []
    
    this.visibleStart = 0
    this.visibleEnd = 0
    this.scrollTop = 0
    
    this.updateVisibleRange()
  }

  updateVisibleRange() {
    const visibleCount = Math.ceil(this.containerHeight / this.itemHeight)
    this.visibleStart = Math.max(0, Math.floor(this.scrollTop / this.itemHeight) - this.buffer)
    this.visibleEnd = Math.min(this.items.length, this.visibleStart + visibleCount + this.buffer * 2)
  }

  onScroll(scrollTop) {
    this.scrollTop = scrollTop
    this.updateVisibleRange()
  }

  getVisibleItems() {
    return this.items.slice(this.visibleStart, this.visibleEnd)
  }

  getOffsetY() {
    return this.visibleStart * this.itemHeight
  }

  getTotalHeight() {
    return this.items.length * this.itemHeight
  }

  updateItems(newItems) {
    this.items = newItems
    this.updateVisibleRange()
  }
}

// 缓存管理器
export class CacheManager {
  constructor(maxSize = 50) {
    this.cache = new Map()
    this.maxSize = maxSize
  }

  set(key, value, ttl = null) {
    // 如果缓存已满，删除最旧的项
    if (this.cache.size >= this.maxSize) {
      const firstKey = this.cache.keys().next().value
      this.cache.delete(firstKey)
    }

    const item = {
      value,
      timestamp: Date.now(),
      ttl
    }

    this.cache.set(key, item)
  }

  get(key) {
    const item = this.cache.get(key)
    
    if (!item) return null

    // 检查是否过期
    if (item.ttl && Date.now() - item.timestamp > item.ttl) {
      this.cache.delete(key)
      return null
    }

    // 更新访问时间（LRU）
    this.cache.delete(key)
    this.cache.set(key, item)

    return item.value
  }

  has(key) {
    return this.get(key) !== null
  }

  delete(key) {
    return this.cache.delete(key)
  }

  clear() {
    this.cache.clear()
  }

  size() {
    return this.cache.size
  }

  // 清理过期项
  cleanup() {
    const now = Date.now()
    for (const [key, item] of this.cache.entries()) {
      if (item.ttl && now - item.timestamp > item.ttl) {
        this.cache.delete(key)
      }
    }
  }
}

// 全局性能监控实例
export const performanceMonitor = new PerformanceMonitor()
export const memoryMonitor = new MemoryMonitor()
export const networkMonitor = new NetworkMonitor()
export const cacheManager = new CacheManager()

// 自动启动内存监控（仅在开发环境）
if (import.meta.env.DEV) {
  memoryMonitor.startMonitoring()
}

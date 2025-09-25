import api from './api'

export const speechApi = {
  // 语音转文字
  speechToText(audioData, options = {}) {
    const formData = new FormData()
    
    // 如果是文件对象
    if (audioData instanceof File) {
      formData.append('audio', audioData)
    } else if (audioData instanceof Blob) {
      // 如果是Blob对象
      formData.append('audio', audioData, 'audio.wav')
    } else {
      // 如果是base64字符串
      formData.append('audio_data', audioData)
      formData.append('format', options.format || 'wav')
    }
    
    // 添加其他参数
    if (options.language) {
      formData.append('language', options.language)
    }
    if (options.sampleRate) {
      formData.append('sample_rate', options.sampleRate)
    }
    
    return api.post('/api/speech/stt', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  // 文字转语音
  textToSpeech(data) {
    return api.post('/api/speech/tts', {
      text: data.text,
      character_id: data.characterId,
      voice_settings: data.voiceSettings || {
        rate: 1.0,
        pitch: 1.0,
        volume: 0.8
      },
      format: data.format || 'mp3'
    })
  },

  // 上传语音文件
  uploadAudio(file) {
    const formData = new FormData()
    formData.append('audio', file)
    
    return api.post('/api/speech/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  // 获取语音文件
  getAudio(id) {
    return api.get(`/api/speech/audio/${id}`, {
      responseType: 'blob'
    })
  }
}

// 语音录制相关工具函数
export class VoiceRecorder {
  constructor() {
    this.mediaRecorder = null
    this.stream = null
    this.chunks = []
    this.isRecording = false
  }

  // 开始录音
  async startRecording() {
    try {
      this.stream = await navigator.mediaDevices.getUserMedia({ 
        audio: {
          sampleRate: 16000,
          channelCount: 1,
          echoCancellation: true,
          noiseSuppression: true
        }
      })
      
      this.mediaRecorder = new MediaRecorder(this.stream, {
        mimeType: 'audio/webm;codecs=opus'
      })
      
      this.chunks = []
      
      this.mediaRecorder.ondataavailable = (event) => {
        if (event.data.size > 0) {
          this.chunks.push(event.data)
        }
      }
      
      this.mediaRecorder.start()
      this.isRecording = true
      
      return true
    } catch (error) {
      console.error('录音启动失败:', error)
      throw new Error('无法访问麦克风，请检查权限设置')
    }
  }

  // 停止录音
  stopRecording() {
    return new Promise((resolve) => {
      if (!this.mediaRecorder || !this.isRecording) {
        resolve(null)
        return
      }

      this.mediaRecorder.onstop = () => {
        const blob = new Blob(this.chunks, { type: 'audio/webm' })
        this.cleanup()
        resolve(blob)
      }

      this.mediaRecorder.stop()
      this.isRecording = false
    })
  }

  // 清理资源
  cleanup() {
    if (this.stream) {
      this.stream.getTracks().forEach(track => track.stop())
      this.stream = null
    }
    this.mediaRecorder = null
    this.chunks = []
    this.isRecording = false
  }

  // 获取录音状态
  getRecordingState() {
    return this.isRecording
  }
}

// 语音播放相关工具函数
export class VoicePlayer {
  constructor() {
    this.audio = new Audio()
    this.isPlaying = false
    this.currentSource = null
  }

  // 播放语音
  async play(audioData, options = {}) {
    try {
      // 停止当前播放
      this.stop()

      let audioUrl
      if (typeof audioData === 'string') {
        // 如果是URL或base64
        audioUrl = audioData.startsWith('data:') ? audioData : audioData
      } else if (audioData instanceof Blob) {
        // 如果是Blob对象
        audioUrl = URL.createObjectURL(audioData)
        this.currentSource = audioUrl
      }

      this.audio.src = audioUrl
      
      // 设置音频参数
      if (options.volume !== undefined) {
        this.audio.volume = Math.max(0, Math.min(1, options.volume))
      }
      if (options.playbackRate !== undefined) {
        this.audio.playbackRate = Math.max(0.25, Math.min(4, options.playbackRate))
      }

      // 播放完成后的回调
      this.audio.onended = () => {
        this.isPlaying = false
        this.cleanup()
        if (options.onEnded) {
          options.onEnded()
        }
      }

      // 播放错误处理
      this.audio.onerror = (error) => {
        console.error('音频播放失败:', error)
        this.isPlaying = false
        this.cleanup()
        if (options.onError) {
          options.onError(error)
        }
      }

      await this.audio.play()
      this.isPlaying = true

      return true
    } catch (error) {
      console.error('语音播放失败:', error)
      this.cleanup()
      throw error
    }
  }

  // 停止播放
  stop() {
    if (this.audio && !this.audio.paused) {
      this.audio.pause()
      this.audio.currentTime = 0
    }
    this.isPlaying = false
    this.cleanup()
  }

  // 暂停播放
  pause() {
    if (this.audio && !this.audio.paused) {
      this.audio.pause()
      this.isPlaying = false
    }
  }

  // 恢复播放
  resume() {
    if (this.audio && this.audio.paused) {
      this.audio.play()
      this.isPlaying = true
    }
  }

  // 清理资源
  cleanup() {
    if (this.currentSource) {
      URL.revokeObjectURL(this.currentSource)
      this.currentSource = null
    }
  }

  // 获取播放状态
  getPlayingState() {
    return this.isPlaying
  }

  // 获取播放进度
  getProgress() {
    if (!this.audio) return 0
    return this.audio.duration ? this.audio.currentTime / this.audio.duration : 0
  }
}

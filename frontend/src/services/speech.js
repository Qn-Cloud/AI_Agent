import { speechApi } from './apiFactory'

export const speechApiService = {
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
    
    return speechApi.post('/api/speech/stt', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  // 文字转语音
  textToSpeech(data) {
    return speechApi.post('/api/speech/tts', {
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
    
    return speechApi.post('/api/speech/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  // 获取语音文件
  getAudio(id) {
    return speechApi.get(`/api/speech/audio/${id}`, {
      responseType: 'blob'
    })
  }
}

// 完整的VoiceRecorder类，支持语音转文字
export class VoiceRecorder {
  constructor() {
    this.isRecording = false
    this.mediaRecorder = null
    this.stream = null
    this.audioChunks = []
    this.onTranscriptCallback = null
  }

  async startRecording() {
    try {
      this.stream = await navigator.mediaDevices.getUserMedia({ audio: true })
      this.mediaRecorder = new MediaRecorder(this.stream)
      this.audioChunks = []
      this.isRecording = true
      
      // 收集音频数据
      this.mediaRecorder.ondataavailable = (event) => {
        if (event.data.size > 0) {
          this.audioChunks.push(event.data)
        }
      }
      
      // 录音结束时处理
      this.mediaRecorder.onstop = async () => {
        await this.processRecording()
      }
      
      this.mediaRecorder.start()
      console.log('🎤 开始录音...')
      return true
    } catch (error) {
      console.error('录音启动失败:', error)
      return false
    }
  }

  stopRecording() {
    if (this.mediaRecorder && this.isRecording) {
      console.log('🎤 停止录音...')
      this.mediaRecorder.stop()
      this.isRecording = false
      if (this.stream) {
        this.stream.getTracks().forEach(track => track.stop())
      }
    }
  }

  async processRecording() {
    try {
      if (this.audioChunks.length === 0) {
        console.warn('没有录音数据')
        return
      }

      // 创建音频Blob
      const audioBlob = new Blob(this.audioChunks, { type: 'audio/wav' })
      console.log('🎤 录音完成，大小:', (audioBlob.size / 1024).toFixed(2), 'KB')

      // 调用语音转文字API
      if (this.onTranscriptCallback) {
        try {
          console.log('🔄 正在转换语音为文字...')
          const response = await speechApiService.speechToText(audioBlob, {
            language: 'zh-CN',
            format: 'wav'
          })
          
          if (response && response.data && response.data.text) {
            console.log('✅ 语音转文字成功:', response.data.text)
            this.onTranscriptCallback(response.data.text)
          } else {
            console.warn('语音转文字返回空结果')
            this.onTranscriptCallback('')
          }
        } catch (error) {
          console.error('❌ 语音转文字失败:', error)
          // 降级方案：显示提示信息
          this.onTranscriptCallback('[语音转文字失败，请重试]')
        }
      }
    } catch (error) {
      console.error('处理录音失败:', error)
    }
  }

  // 设置转录回调函数
  setTranscriptCallback(callback) {
    this.onTranscriptCallback = callback
  }
}

// 简化的VoicePlayer类
export class VoicePlayer {
  constructor() {
    this.audio = new Audio()
    this.isPlaying = false
  }

  async play(audioUrl) {
    try {
      this.audio.src = audioUrl
      await this.audio.play()
      this.isPlaying = true
      return true
    } catch (error) {
      console.error('音频播放失败:', error)
      return false
    }
  }

  stop() {
    this.audio.pause()
    this.audio.currentTime = 0
    this.isPlaying = false
  }
}

// 保持向后兼容
export const speechApi_old = speechApiService
export default speechApiService 
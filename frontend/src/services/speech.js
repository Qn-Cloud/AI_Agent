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

// 简化的VoiceRecorder类
export class VoiceRecorder {
  constructor() {
    this.isRecording = false
    this.mediaRecorder = null
    this.stream = null
  }

  async startRecording() {
    try {
      this.stream = await navigator.mediaDevices.getUserMedia({ audio: true })
      this.mediaRecorder = new MediaRecorder(this.stream)
      this.isRecording = true
      this.mediaRecorder.start()
      return true
    } catch (error) {
      console.error('录音启动失败:', error)
      return false
    }
  }

  stopRecording() {
    if (this.mediaRecorder && this.isRecording) {
      this.mediaRecorder.stop()
      this.isRecording = false
      if (this.stream) {
        this.stream.getTracks().forEach(track => track.stop())
      }
    }
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
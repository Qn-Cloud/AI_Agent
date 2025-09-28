# AI 角色扮演 - 语音服务 (Speech Service)

## 🎤 服务简介

语音服务提供语音转文字(ASR)和文字转语音(TTS)功能，支持AI角色扮演系统的语音交互需求。

## 🚀 快速启动

### 启动服务

**Linux/Mac:**
```bash
cd backend/services/speech
chmod +x start.sh
./start.sh
```

**Windows:**
```cmd
cd backend\services\speech
start.bat
```

### 服务信息
- **端口**: 7005
- **健康检查**: http://localhost:7005/api/speech/health
- **API前缀**: /api/speech

## 📡 API 接口

### 1. 语音转文字 (STT)

**POST** `/api/speech/stt`

**请求方式1: Multipart Form Data (推荐)**
```bash
curl -X POST "http://localhost:7005/api/speech/stt" \
  -F "audio=@audio.wav" \
  -F "language=zh-CN" \
  -F "format=wav"
```

**请求方式2: JSON with Base64**
```bash
curl -X POST "http://localhost:7005/api/speech/stt" \
  -H "Content-Type: application/json" \
  -d '{
    "audio_data": "base64_encoded_audio_data",
    "language": "zh-CN",
    "format": "wav",
    "sample_rate": 16000
  }'
```

**参数说明:**
- `audio`: 音频文件 (multipart)
- `audio_data`: Base64编码的音频数据 (JSON)
- `language`: 语言代码，默认 "zh-CN"
- `format`: 音频格式，默认 "wav"
- `sample_rate`: 采样率，默认 16000

**响应示例:**
```json
{
  "code": 200,
  "msg": "语音转文字成功",
  "data": {
    "text": "你好，我是人工智能助手",
    "confidence": 0.92,
    "duration": 2.5
  }
}
```

### 2. 文字转语音 (TTS)

**POST** `/api/speech/tts`

**请求示例:**
```bash
curl -X POST "http://localhost:7005/api/speech/tts" \
  -H "Content-Type: application/json" \
  -d '{
    "text": "你好，我是人工智能助手",
    "character_id": 1,
    "voice_settings": {
      "rate": 1.0,
      "pitch": 1.0,
      "volume": 0.8
    },
    "format": "mp3"
  }'
```

**参数说明:**
- `text`: 要转换的文本
- `character_id`: 角色ID (可选)
- `voice_settings`: 语音设置
  - `rate`: 语速，默认 1.0
  - `pitch`: 音调，默认 1.0  
  - `volume`: 音量，默认 0.8
- `format`: 输出格式，默认 "mp3"

**响应示例:**
```json
{
  "code": 200,
  "msg": "文字转语音成功",
  "data": {
    "audio_url": "https://cdn.example.com/audio/tts/abc123.mp3",
    "duration": 3.2,
    "size": 51200
  }
}
```

### 3. 健康检查

**GET** `/api/speech/health`

**响应示例:**
```json
{
  "code": 200,
  "msg": "语音服务运行正常"
}
```

## ⚙️ 配置说明

配置文件位置: `api/etc/speech.yaml`

```yaml
Name: Speech
Host: 0.0.0.0
Port: 7005
MaxBytes: 33554432  # 32MB

Speech:
  ASR:
    Provider: "mock"  # 语音识别提供商
    Timeout: 30s
    MaxDuration: 60s
    SupportedFormats: ["wav", "mp3", "ogg", "m4a"]
    
  TTS:
    Provider: "mock"  # 语音合成提供商
    Timeout: 30s
    MaxTextLength: 1000
    DefaultFormat: "mp3"
    SupportedFormats: ["mp3", "wav", "ogg"]
```

## 🔧 集成第三方服务

当前使用模拟实现，您可以集成以下真实服务:

### 百度AI
```yaml
External:
  Baidu:
    AppID: "your_app_id"
    APIKey: "your_api_key"
    SecretKey: "your_secret_key"
```

### 腾讯云
```yaml
External:
  Tencent:
    SecretID: "your_secret_id"
    SecretKey: "your_secret_key"
    Region: "ap-beijing"
```

### 阿里云
```yaml
External:
  Aliyun:
    AccessKeyID: "your_access_key_id"
    AccessKeySecret: "your_access_key_secret"
    Region: "cn-hangzhou"
```

## 📝 开发说明

### 目录结构
```
speech/
├── api/                    # API服务
│   ├── etc/               # 配置文件
│   ├── internal/          # 内部实现
│   │   ├── handler/       # HTTP处理器
│   │   ├── logic/         # 业务逻辑
│   │   ├── svc/           # 服务上下文
│   │   └── types/         # 类型定义
│   ├── speech.api         # API定义
│   ├── speech_types.api   # 类型定义
│   └── speech.go          # 入口文件
├── model/                 # 数据模型
├── rpc/                   # RPC服务
├── start.sh               # Linux启动脚本
├── start.bat              # Windows启动脚本
└── README.md              # 文档
```

### 添加新的ASR/TTS服务

1. 在 `logic/speechtotextlogic.go` 中添加新的实现函数
2. 在 `logic/texttospeechlogic.go` 中添加新的实现函数
3. 更新配置文件支持新的提供商
4. 在业务逻辑中根据配置选择对应的实现

### 测试

```bash
# 测试健康检查
curl http://localhost:7005/api/speech/health

# 测试语音转文字
curl -X POST "http://localhost:7005/api/speech/stt" \
  -F "audio=@test.wav"

# 测试文字转语音  
curl -X POST "http://localhost:7005/api/speech/stt" \
  -H "Content-Type: application/json" \
  -d '{"text": "测试语音合成"}'
```

## 🔍 常见问题

### 1. 音频格式支持
- 支持格式: WAV, MP3, OGG, M4A
- 推荐格式: WAV (16kHz, 16bit)
- 最大文件大小: 32MB

### 2. 文本长度限制
- STT: 建议60秒内的音频
- TTS: 最大1000字符

### 3. 性能优化
- 使用音频压缩减少传输时间
- 启用缓存避免重复转换
- 异步处理长音频文件

## 📞 技术支持

如有问题，请查看日志或联系开发团队。

---

**版本**: v1.0  
**更新时间**: 2025-09-28  
**负责团队**: AI Agent Team 
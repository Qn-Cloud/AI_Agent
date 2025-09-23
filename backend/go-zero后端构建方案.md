# AI 角色扮演语音交互产品 - 后端构建方案（Go-Zero）

## 1. 项目概述

### 1.1 产品定位
基于 Go-Zero 微服务框架构建的AI角色扮演语音交互产品后端系统，为前端提供稳定、高性能的API服务，支持角色管理、实时语音处理、AI对话和用户数据管理等核心功能。

### 1.2 核心功能
- **角色管理服务**：预设角色数据管理、自定义角色创建
- **AI对话服务**：智能对话生成、上下文管理
- **语音处理服务**：ASR语音识别、TTS语音合成
- **用户服务**：用户认证、个人设置、对话历史
- **存储服务**：文件上传、语音文件管理

## 2. 技术选型

### 2.1 核心框架：Go-Zero
**选择理由：**
- **微服务架构**：天然支持微服务拆分，便于横向扩展
- **代码生成**：通过 API 定义自动生成代码，提高开发效率
- **内置中间件**：集成限流、熔断、链路追踪等功能
- **高性能**：基于 Go 语言，并发性能优异
- **生态完善**：配套工具链完整，文档丰富

### 2.2 数据库选择
**主数据库：MySQL 8.0**
- 用户数据、角色信息、对话记录等结构化数据存储
- 支持事务，保证数据一致性

**缓存数据库：Redis 6.0**
- 用户会话缓存、热点数据缓存
- 分布式锁、限流计数器

**向量数据库：Milvus/Qdrant**
- 角色向量存储，支持语义搜索
- 对话历史向量化检索

### 2.3 消息队列：NATS/RabbitMQ
- 异步任务处理（语音转换、AI推理）
- 服务间解耦通信
- 事件驱动架构支持

### 2.4 AI服务集成
**大语言模型：**
- **OpenAI GPT-4/3.5**：主要对话生成
- **备选方案**：Claude、智谱AI、百度文心一言

**语音服务：**
- **ASR**：阿里云语音识别、腾讯云语音识别
- **TTS**：Azure语音服务、科大讯飞

### 2.5 其他技术栈
- **配置管理**：go-zero config
- **日志管理**：go-zero logx + ELK Stack
- **监控告警**：Prometheus + Grafana
- **链路追踪**：Jaeger
- **容器化**：Docker + Kubernetes
- **CI/CD**：GitLab CI/CD

## 3. 微服务架构设计

### 3.1 服务拆分

```
backend/
├── gateway/                    # API网关
│   ├── etc/
│   ├── internal/
│   └── gateway.go
├── services/
│   ├── user/                   # 用户服务
│   │   ├── api/               # HTTP API
│   │   ├── rpc/               # gRPC服务
│   │   └── model/             # 数据模型
│   ├── character/             # 角色服务
│   │   ├── api/
│   │   ├── rpc/
│   │   └── model/
│   ├── chat/                  # 对话服务
│   │   ├── api/
│   │   ├── rpc/
│   │   └── model/
│   ├── ai/                    # AI服务
│   │   ├── api/
│   │   ├── rpc/
│   │   └── model/
│   ├── speech/                # 语音服务
│   │   ├── api/
│   │   ├── rpc/
│   │   └── model/
│   └── storage/               # 存储服务
│       ├── api/
│       ├── rpc/
│       └── model/
├── common/                    # 公共库
│   ├── model/                 # 通用数据模型
│   ├── utils/                 # 工具函数
│   ├── middleware/            # 中间件
│   └── response/              # 统一响应格式
├── deploy/                    # 部署配置
│   ├── docker/
│   ├── k8s/
│   └── scripts/
└── docs/                      # 文档
```

### 3.2 服务职责划分

#### 3.2.1 API网关（Gateway）
- 统一入口，路由分发
- 身份认证和授权
- 限流熔断
- 请求日志记录

#### 3.2.2 用户服务（User Service）
- 用户注册、登录、注销
- 用户信息管理
- 个人设置配置
- 对话历史管理

#### 3.2.3 角色服务（Character Service）
- 预设角色数据管理
- 角色搜索和推荐
- 自定义角色创建
- 角色评分和收藏

#### 3.2.4 对话服务（Chat Service）
- 对话会话管理
- 消息记录存储
- 上下文维护
- 对话历史查询

#### 3.2.5 AI服务（AI Service）
- LLM模型调用
- Prompt工程
- 对话生成
- 内容安全检查

#### 3.2.6 语音服务（Speech Service）
- ASR语音识别
- TTS语音合成
- 音频文件处理
- 语音质量优化

#### 3.2.7 存储服务（Storage Service）
- 文件上传下载
- 音频文件管理
- 静态资源服务
- 文件清理策略

## 4. API接口设计

### 4.1 用户服务API

```go
// user.api
syntax = "v1"

info(
    title: "用户服务"
    desc: "用户注册登录、个人信息管理"
    author: "backend team"
    version: "1.0"
)

// 用户注册
type (
    RegisterRequest {
        Username string `json:"username"`
        Password string `json:"password"`
        Email    string `json:"email"`
    }
    
    RegisterResponse {
        UserId int64  `json:"user_id"`
        Token  string `json:"token"`
    }
)

// 用户登录
type (
    LoginRequest {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    
    LoginResponse {
        UserId   int64  `json:"user_id"`
        Token    string `json:"token"`
        Username string `json:"username"`
        Avatar   string `json:"avatar"`
    }
)

// 用户信息
type (
    UserInfoResponse {
        UserId     int64  `json:"user_id"`
        Username   string `json:"username"`
        Email      string `json:"email"`
        Avatar     string `json:"avatar"`
        CreatedAt  string `json:"created_at"`
    }
)

service user-api {
    @handler registerHandler
    post /api/user/register (RegisterRequest) returns (RegisterResponse)
    
    @handler loginHandler
    post /api/user/login (LoginRequest) returns (LoginResponse)
    
    @handler userInfoHandler
    get /api/user/info returns (UserInfoResponse)
}
```

### 4.2 角色服务API

```go
// character.api
syntax = "v1"

info(
    title: "角色服务"
    desc: "角色管理、搜索、自定义"
    version: "1.0"
)

// 角色信息
type (
    Character {
        Id          int64    `json:"id"`
        Name        string   `json:"name"`
        Avatar      string   `json:"avatar"`
        Description string   `json:"description"`
        Tags        []string `json:"tags"`
        Prompt      string   `json:"prompt"`
        Status      string   `json:"status"`
        Rating      float64  `json:"rating"`
        CreatedAt   string   `json:"created_at"`
    }
    
    CharacterListRequest {
        Page     int    `form:"page,default=1"`
        PageSize int    `form:"page_size,default=20"`
        Category string `form:"category,optional"`
        Search   string `form:"search,optional"`
    }
    
    CharacterListResponse {
        List  []Character `json:"list"`
        Total int64       `json:"total"`
    }
    
    CharacterDetailResponse {
        Character Character `json:"character"`
    }
    
    CreateCharacterRequest {
        Name        string   `json:"name"`
        Avatar      string   `json:"avatar"`
        Description string   `json:"description"`
        Tags        []string `json:"tags"`
        Prompt      string   `json:"prompt"`
    }
    
    CreateCharacterResponse {
        CharacterId int64 `json:"character_id"`
    }
)

service character-api {
    @handler listHandler
    get /api/character/list (CharacterListRequest) returns (CharacterListResponse)
    
    @handler detailHandler
    get /api/character/:id returns (CharacterDetailResponse)
    
    @handler createHandler
    post /api/character (CreateCharacterRequest) returns (CreateCharacterResponse)
    
    @handler favoriteHandler
    post /api/character/:id/favorite
    
    @handler searchHandler
    get /api/character/search (CharacterListRequest) returns (CharacterListResponse)
}
```

### 4.3 对话服务API

```go
// chat.api
syntax = "v1"

info(
    title: "对话服务"
    desc: "对话管理、消息处理"
    version: "1.0"
)

type (
    Message {
        Id        int64  `json:"id"`
        Type      string `json:"type"` // user, ai
        Content   string `json:"content"`
        AudioUrl  string `json:"audio_url,optional"`
        Timestamp string `json:"timestamp"`
    }
    
    Conversation {
        Id          int64     `json:"id"`
        CharacterId int64     `json:"character_id"`
        Title       string    `json:"title"`
        StartTime   string    `json:"start_time"`
        LastUpdate  string    `json:"last_update"`
        Messages    []Message `json:"messages"`
    }
    
    SendMessageRequest {
        ConversationId int64  `json:"conversation_id"`
        Content        string `json:"content"`
        Type          string `json:"type"` // text, audio
        AudioUrl      string `json:"audio_url,optional"`
    }
    
    SendMessageResponse {
        Message Message `json:"message"`
        AiReply Message `json:"ai_reply"`
    }
    
    ConversationListRequest {
        Page     int   `form:"page,default=1"`
        PageSize int   `form:"page_size,default=20"`
        CharacterId int64 `form:"character_id,optional"`
    }
    
    ConversationListResponse {
        List  []Conversation `json:"list"`
        Total int64          `json:"total"`
    }
    
    CreateConversationRequest {
        CharacterId int64  `json:"character_id"`
        Title       string `json:"title,optional"`
    }
    
    CreateConversationResponse {
        ConversationId int64 `json:"conversation_id"`
    }
)

service chat-api {
    @handler sendMessageHandler
    post /api/chat/message (SendMessageRequest) returns (SendMessageResponse)
    
    @handler listConversationsHandler
    get /api/chat/conversations (ConversationListRequest) returns (ConversationListResponse)
    
    @handler createConversationHandler
    post /api/chat/conversation (CreateConversationRequest) returns (CreateConversationResponse)
    
    @handler getConversationHandler
    get /api/chat/conversation/:id returns (Conversation)
    
    @handler deleteConversationHandler
    delete /api/chat/conversation/:id
}
```

### 4.4 语音服务API

```go
// speech.api
syntax = "v1"

info(
    title: "语音服务"
    desc: "语音识别和合成"
    version: "1.0"
)

type (
    SpeechToTextRequest {
        AudioUrl string `json:"audio_url"`
        Language string `json:"language,default=zh-CN"`
    }
    
    SpeechToTextResponse {
        Text       string  `json:"text"`
        Confidence float64 `json:"confidence"`
    }
    
    TextToSpeechRequest {
        Text     string `json:"text"`
        Voice    string `json:"voice,optional"`
        Speed    float64 `json:"speed,default=1.0"`
        Pitch    float64 `json:"pitch,default=1.0"`
        Volume   float64 `json:"volume,default=1.0"`
    }
    
    TextToSpeechResponse {
        AudioUrl string `json:"audio_url"`
        Duration int64  `json:"duration"`
    }
)

service speech-api {
    @handler speechToTextHandler
    post /api/speech/stt (SpeechToTextRequest) returns (SpeechToTextResponse)
    
    @handler textToSpeechHandler
    post /api/speech/tts (TextToSpeechRequest) returns (TextToSpeechResponse)
    
    @handler uploadAudioHandler
    post /api/speech/upload returns (string)
}
```

## 5. 数据库设计

### 5.1 用户表（users）

```sql
CREATE TABLE `users` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL UNIQUE,
  `email` varchar(100) NOT NULL UNIQUE,
  `password_hash` varchar(255) NOT NULL,
  `avatar` varchar(255) DEFAULT '',
  `nickname` varchar(50) DEFAULT '',
  `bio` text,
  `status` tinyint DEFAULT '1' COMMENT '1:正常 0:禁用',
  `last_login_at` timestamp NULL,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_username` (`username`),
  KEY `idx_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 5.2 角色表（characters）

```sql
CREATE TABLE `characters` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `avatar` varchar(255) DEFAULT '',
  `description` text,
  `prompt` text NOT NULL,
  `tags` json,
  `category` varchar(50) DEFAULT '',
  `rating` decimal(3,2) DEFAULT '0.00',
  `rating_count` int DEFAULT '0',
  `favorite_count` int DEFAULT '0',
  `chat_count` int DEFAULT '0',
  `status` tinyint DEFAULT '1' COMMENT '1:启用 0:禁用',
  `is_public` tinyint DEFAULT '1' COMMENT '1:公开 0:私有',
  `creator_id` bigint DEFAULT '0',
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_name` (`name`),
  KEY `idx_category` (`category`),
  KEY `idx_creator` (`creator_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 5.3 对话表（conversations）

```sql
CREATE TABLE `conversations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `character_id` bigint NOT NULL,
  `title` varchar(200) DEFAULT '',
  `start_time` timestamp DEFAULT CURRENT_TIMESTAMP,
  `last_message_time` timestamp NULL,
  `message_count` int DEFAULT '0',
  `status` tinyint DEFAULT '1' COMMENT '1:活跃 0:结束',
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_character_id` (`character_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 5.4 消息表（messages）

```sql
CREATE TABLE `messages` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `conversation_id` bigint NOT NULL,
  `user_id` bigint NOT NULL,
  `type` varchar(10) NOT NULL COMMENT 'user/ai',
  `content` text NOT NULL,
  `audio_url` varchar(255) DEFAULT '',
  `audio_duration` int DEFAULT '0',
  `metadata` json,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_conversation_id` (`conversation_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_type` (`type`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 5.5 用户收藏表（user_favorites）

```sql
CREATE TABLE `user_favorites` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `character_id` bigint NOT NULL,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_character` (`user_id`, `character_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_character_id` (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## 6. 核心服务实现

### 6.1 AI对话服务实现

```go
// services/ai/internal/logic/chatlogic.go
package logic

import (
    "context"
    "encoding/json"
    
    "ai-roleplay/services/ai/internal/svc"
    "ai-roleplay/services/ai/internal/types"
    
    "github.com/zeromicro/go-zero/core/logx"
)

type ChatLogic struct {
    logx.Logger
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func NewChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatLogic {
    return &ChatLogic{
        Logger: logx.WithContext(ctx),
        ctx:    ctx,
        svcCtx: svcCtx,
    }
}

func (l *ChatLogic) Chat(req *types.ChatRequest) (resp *types.ChatResponse, err error) {
    // 1. 获取角色信息和Prompt
    character, err := l.svcCtx.CharacterRpc.GetCharacter(l.ctx, &characterrpc.GetCharacterRequest{
        Id: req.CharacterId,
    })
    if err != nil {
        return nil, err
    }

    // 2. 构建对话上下文
    messages := l.buildChatContext(req.ConversationId, character.Prompt, req.Content)

    // 3. 调用LLM生成回复
    aiResponse, err := l.callLLM(messages)
    if err != nil {
        return nil, err
    }

    // 4. 内容安全检查
    if err := l.contentSafetyCheck(aiResponse); err != nil {
        aiResponse = "抱歉，我无法回答这个问题。"
    }

    // 5. 保存对话记录
    err = l.saveMessage(req.ConversationId, req.UserId, "ai", aiResponse)
    if err != nil {
        logx.Error("保存AI回复失败:", err)
    }

    return &types.ChatResponse{
        Content:   aiResponse,
        MessageId: l.generateMessageId(),
    }, nil
}

func (l *ChatLogic) buildChatContext(conversationId int64, prompt string, userInput string) []ChatMessage {
    // 从数据库获取最近的对话历史
    history, _ := l.svcCtx.ChatRpc.GetRecentMessages(l.ctx, &chatrpc.GetRecentMessagesRequest{
        ConversationId: conversationId,
        Limit:         10,
    })

    messages := []ChatMessage{
        {Role: "system", Content: prompt},
    }

    // 添加历史对话
    for _, msg := range history.Messages {
        messages = append(messages, ChatMessage{
            Role:    msg.Type,
            Content: msg.Content,
        })
    }

    // 添加用户当前输入
    messages = append(messages, ChatMessage{
        Role:    "user",
        Content: userInput,
    })

    return messages
}

func (l *ChatLogic) callLLM(messages []ChatMessage) (string, error) {
    // 这里可以集成OpenAI、Claude等LLM服务
    // 示例使用OpenAI
    client := openai.NewClient(l.svcCtx.Config.OpenAI.ApiKey)
    
    resp, err := client.CreateChatCompletion(
        context.Background(),
        openai.ChatCompletionRequest{
            Model:    openai.GPT3Dot5Turbo,
            Messages: convertToOpenAIMessages(messages),
            MaxTokens: 1000,
            Temperature: 0.7,
        },
    )
    
    if err != nil {
        return "", err
    }
    
    return resp.Choices[0].Message.Content, nil
}
```

### 6.2 语音处理服务

```go
// services/speech/internal/logic/sttlogic.go
package logic

import (
    "context"
    "fmt"
    "io"
    "net/http"
    
    "ai-roleplay/services/speech/internal/svc"
    "ai-roleplay/services/speech/internal/types"
    
    "github.com/zeromicro/go-zero/core/logx"
)

type SttLogic struct {
    logx.Logger
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func NewSttLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SttLogic {
    return &SttLogic{
        Logger: logx.WithContext(ctx),
        ctx:    ctx,
        svcCtx: svcCtx,
    }
}

func (l *SttLogic) SpeechToText(req *types.SttRequest) (resp *types.SttResponse, err error) {
    // 1. 下载音频文件
    audioData, err := l.downloadAudio(req.AudioUrl)
    if err != nil {
        return nil, fmt.Errorf("下载音频失败: %w", err)
    }

    // 2. 调用ASR服务
    text, confidence, err := l.callASRService(audioData, req.Language)
    if err != nil {
        return nil, fmt.Errorf("语音识别失败: %w", err)
    }

    // 3. 内容过滤
    filteredText := l.contentFilter(text)

    return &types.SttResponse{
        Text:       filteredText,
        Confidence: confidence,
    }, nil
}

func (l *SttLogic) downloadAudio(url string) ([]byte, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}

func (l *SttLogic) callASRService(audioData []byte, language string) (string, float64, error) {
    // 这里可以集成阿里云、腾讯云等ASR服务
    // 示例代码
    asrClient := l.svcCtx.ASRClient
    
    result, err := asrClient.Recognize(l.ctx, &ASRRequest{
        AudioData: audioData,
        Language:  language,
        Format:    "wav",
    })
    
    if err != nil {
        return "", 0, err
    }
    
    return result.Text, result.Confidence, nil
}
```

### 6.3 TTS语音合成服务

```go
// services/speech/internal/logic/ttslogic.go
package logic

import (
    "context"
    "fmt"
    
    "ai-roleplay/services/speech/internal/svc"
    "ai-roleplay/services/speech/internal/types"
    
    "github.com/zeromicro/go-zero/core/logx"
)

type TtsLogic struct {
    logx.Logger
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func NewTtsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TtsLogic {
    return &TtsLogic{
        Logger: logx.WithContext(ctx),
        ctx:    ctx,
        svcCtx: svcCtx,
    }
}

func (l *TtsLogic) TextToSpeech(req *types.TtsRequest) (resp *types.TtsResponse, err error) {
    // 1. 文本预处理
    processedText := l.preprocessText(req.Text)

    // 2. 调用TTS服务
    audioData, duration, err := l.callTTSService(processedText, req.Voice, req.Speed, req.Pitch)
    if err != nil {
        return nil, fmt.Errorf("语音合成失败: %w", err)
    }

    // 3. 上传音频文件
    audioUrl, err := l.uploadAudio(audioData)
    if err != nil {
        return nil, fmt.Errorf("上传音频失败: %w", err)
    }

    return &types.TtsResponse{
        AudioUrl: audioUrl,
        Duration: duration,
    }, nil
}

func (l *TtsLogic) preprocessText(text string) string {
    // 文本清理、特殊字符处理等
    // 移除不适合语音合成的内容
    return text
}

func (l *TtsLogic) callTTSService(text, voice string, speed, pitch float64) ([]byte, int64, error) {
    // 这里可以集成Azure、科大讯飞等TTS服务
    ttsClient := l.svcCtx.TTSClient
    
    result, err := ttsClient.Synthesize(l.ctx, &TTSRequest{
        Text:   text,
        Voice:  voice,
        Speed:  speed,
        Pitch:  pitch,
        Format: "mp3",
    })
    
    if err != nil {
        return nil, 0, err
    }
    
    return result.AudioData, result.Duration, nil
}

func (l *TtsLogic) uploadAudio(audioData []byte) (string, error) {
    // 调用存储服务上传音频文件
    storageResp, err := l.svcCtx.StorageRpc.UploadFile(l.ctx, &storagerpc.UploadFileRequest{
        FileData: audioData,
        FileType: "audio/mp3",
        Category: "tts",
    })
    
    if err != nil {
        return "", err
    }
    
    return storageResp.FileUrl, nil
}
```

## 7. 中间件和工具

### 7.1 认证中间件

```go
// common/middleware/auth.go
package middleware

import (
    "net/http"
    "strings"
    
    "github.com/zeromicro/go-zero/rest/httpx"
    "github.com/zeromicro/go-zero/core/logx"
)

type AuthMiddleware struct {
    Secret string
}

func NewAuthMiddleware(secret string) *AuthMiddleware {
    return &AuthMiddleware{
        Secret: secret,
    }
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 1. 获取Token
        token := r.Header.Get("Authorization")
        if token == "" {
            httpx.Error(w, errors.New("missing authorization header"))
            return
        }
        
        if !strings.HasPrefix(token, "Bearer ") {
            httpx.Error(w, errors.New("invalid authorization format"))
            return
        }
        
        token = strings.TrimPrefix(token, "Bearer ")
        
        // 2. 验证Token
        claims, err := m.parseToken(token)
        if err != nil {
            httpx.Error(w, errors.New("invalid token"))
            return
        }
        
        // 3. 设置用户信息到上下文
        ctx := context.WithValue(r.Context(), "userId", claims.UserId)
        ctx = context.WithValue(ctx, "username", claims.Username)
        
        next(w, r.WithContext(ctx))
    })
}

func (m *AuthMiddleware) parseToken(tokenString string) (*Claims, error) {
    // JWT Token解析逻辑
    // ...
}
```

### 7.2 限流中间件

```go
// common/middleware/ratelimit.go
package middleware

import (
    "net/http"
    "time"
    
    "github.com/zeromicro/go-zero/core/limit"
    "github.com/zeromicro/go-zero/rest/httpx"
)

type RateLimitMiddleware struct {
    limiter *limit.TokenLimiter
}

func NewRateLimitMiddleware(rate int, burst int) *RateLimitMiddleware {
    limiter := limit.NewTokenLimiter(rate, burst, time.Minute, "rate-limit")
    return &RateLimitMiddleware{
        limiter: limiter,
    }
}

func (m *RateLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 获取客户端标识（IP或用户ID）
        clientId := m.getClientId(r)
        
        // 检查限流
        if m.limiter.Allow(clientId) {
            next(w, r)
        } else {
            httpx.Error(w, errors.New("rate limit exceeded"))
        }
    })
}

func (m *RateLimitMiddleware) getClientId(r *http.Request) string {
    // 优先使用用户ID，否则使用IP
    userId := r.Context().Value("userId")
    if userId != nil {
        return fmt.Sprintf("user:%v", userId)
    }
    
    return fmt.Sprintf("ip:%s", r.RemoteAddr)
}
```

### 7.3 统一响应格式

```go
// common/response/response.go
package response

import (
    "net/http"
    
    "github.com/zeromicro/go-zero/rest/httpx"
)

type Body struct {
    Code int         `json:"code"`
    Msg  string      `json:"msg"`
    Data interface{} `json:"data,omitempty"`
}

func Response(w http.ResponseWriter, resp interface{}, err error) {
    var body Body
    
    if err != nil {
        body.Code = -1
        body.Msg = err.Error()
    } else {
        body.Code = 0
        body.Msg = "success"
        body.Data = resp
    }
    
    httpx.OkJson(w, body)
}

func Success(w http.ResponseWriter, data interface{}) {
    Response(w, data, nil)
}

func Error(w http.ResponseWriter, err error) {
    Response(w, nil, err)
}
```

## 8. 配置管理

### 8.1 API网关配置

```yaml
# gateway/etc/gateway-api.yaml
Name: gateway-api
Host: 0.0.0.0
Port: 8888

Log:
  ServiceName: gateway-api
  Mode: file
  Path: logs
  Level: info

# 服务发现
Etcd:
  Hosts:
    - etcd:2379
  Key: gateway-api

# 上游服务
Upstreams:
  - Name: user-api
    Uris:
      - http://user-api:8001
  - Name: character-api
    Uris:
      - http://character-api:8002
  - Name: chat-api
    Uris:
      - http://chat-api:8003
  - Name: ai-api
    Uris:
      - http://ai-api:8004
  - Name: speech-api
    Uris:
      - http://speech-api:8005

# 认证配置
Auth:
  AccessSecret: your-access-secret
  AccessExpire: 86400

# 限流配置
RateLimit:
  Rate: 100
  Burst: 200

# CORS配置
CORS:
  AllowOrigins:
    - "*"
  AllowMethods:
    - GET
    - POST
    - PUT
    - DELETE
    - OPTIONS
  AllowHeaders:
    - Content-Type
    - Authorization
```

### 8.2 AI服务配置

```yaml
# services/ai/etc/ai-api.yaml
Name: ai-api
Host: 0.0.0.0
Port: 8004

Log:
  ServiceName: ai-api
  Mode: file
  Path: logs
  Level: info

# 数据库配置
DataSource: user:password@tcp(mysql:3306)/ai_roleplay?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai

# Redis配置
RedisConf:
  Host: redis:6379
  Type: node

# OpenAI配置
OpenAI:
  ApiKey: your-openai-api-key
  BaseUrl: https://api.openai.com/v1
  Model: gpt-3.5-turbo
  MaxTokens: 1000
  Temperature: 0.7

# 内容安全
ContentSafety:
  Enable: true
  Provider: aliyun # aliyun, tencent, baidu
  ApiKey: your-content-safety-api-key

# RPC服务发现
CharacterRpc:
  Etcd:
    Hosts:
      - etcd:2379
    Key: character-rpc

ChatRpc:
  Etcd:
    Hosts:
      - etcd:2379
    Key: chat-rpc
```

## 9. 部署方案

### 9.1 Docker容器化

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/etc ./etc

CMD ["./main"]
```

### 9.2 Docker Compose部署

```yaml
# docker-compose.yml
version: '3.8'

services:
  # 基础设施
  etcd:
    image: quay.io/coreos/etcd:v3.5.0
    environment:
      - ETCD_ADVERTISE_CLIENT_URLS=http://etcd:2379
      - ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379
    ports:
      - "2379:2379"

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root123
      MYSQL_DATABASE: ai_roleplay
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql

  redis:
    image: redis:6.2-alpine
    ports:
      - "6379:6379"

  # API网关
  gateway:
    build: ./gateway
    ports:
      - "8888:8888"
    depends_on:
      - etcd
    environment:
      - ETCD_HOSTS=etcd:2379

  # 用户服务
  user-api:
    build: ./services/user/api
    ports:
      - "8001:8001"
    depends_on:
      - mysql
      - redis
      - etcd

  user-rpc:
    build: ./services/user/rpc
    depends_on:
      - mysql
      - redis
      - etcd

  # 角色服务
  character-api:
    build: ./services/character/api
    ports:
      - "8002:8002"
    depends_on:
      - mysql
      - redis
      - etcd

  character-rpc:
    build: ./services/character/rpc
    depends_on:
      - mysql
      - redis
      - etcd

  # 对话服务
  chat-api:
    build: ./services/chat/api
    ports:
      - "8003:8003"
    depends_on:
      - mysql
      - redis
      - etcd

  chat-rpc:
    build: ./services/chat/rpc
    depends_on:
      - mysql
      - redis
      - etcd

  # AI服务
  ai-api:
    build: ./services/ai/api
    ports:
      - "8004:8004"
    depends_on:
      - mysql
      - redis
      - etcd

  # 语音服务
  speech-api:
    build: ./services/speech/api
    ports:
      - "8005:8005"
    depends_on:
      - mysql
      - redis
      - etcd

volumes:
  mysql_data:
```

### 9.3 Kubernetes部署

```yaml
# deploy/k8s/ai-roleplay.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gateway
spec:
  replicas: 3
  selector:
    matchLabels:
      app: gateway
  template:
    metadata:
      labels:
        app: gateway
    spec:
      containers:
      - name: gateway
        image: ai-roleplay/gateway:latest
        ports:
        - containerPort: 8888
        env:
        - name: ETCD_HOSTS
          value: "etcd:2379"
        resources:
          requests:
            memory: "64Mi"
            cpu: "250m"
          limits:
            memory: "128Mi"
            cpu: "500m"

---
apiVersion: v1
kind: Service
metadata:
  name: gateway-service
spec:
  selector:
    app: gateway
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8888
  type: LoadBalancer

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ai-api
spec:
  replicas: 2
  selector:
    matchLabels:
      app: ai-api
  template:
    metadata:
      labels:
        app: ai-api
    spec:
      containers:
      - name: ai-api
        image: ai-roleplay/ai-api:latest
        ports:
        - containerPort: 8004
        env:
        - name: OPENAI_API_KEY
          valueFrom:
            secretKeyRef:
              name: ai-secrets
              key: openai-api-key
        resources:
          requests:
            memory: "128Mi"
            cpu: "500m"
          limits:
            memory: "256Mi"
            cpu: "1000m"
```

## 10. 监控与运维

### 10.1 监控指标

```go
// common/metrics/metrics.go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    // HTTP请求指标
    HttpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )

    HttpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "Duration of HTTP requests",
        },
        []string{"method", "endpoint"},
    )

    // AI服务指标
    AIRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "ai_requests_total",
            Help: "Total number of AI requests",
        },
        []string{"model", "status"},
    )

    AIRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "ai_request_duration_seconds",
            Help: "Duration of AI requests",
        },
        []string{"model"},
    )

    // 语音服务指标
    SpeechRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "speech_requests_total",
            Help: "Total number of speech requests",
        },
        []string{"type", "status"}, // type: tts/stt
    )
)
```

### 10.2 健康检查

```go
// common/health/health.go
package health

import (
    "context"
    "database/sql"
    "fmt"
    "time"
    
    "github.com/go-redis/redis/v8"
)

type HealthChecker struct {
    db    *sql.DB
    redis *redis.Client
}

func NewHealthChecker(db *sql.DB, redis *redis.Client) *HealthChecker {
    return &HealthChecker{
        db:    db,
        redis: redis,
    }
}

func (h *HealthChecker) Check() error {
    // 检查数据库连接
    if err := h.checkDatabase(); err != nil {
        return fmt.Errorf("database health check failed: %w", err)
    }

    // 检查Redis连接
    if err := h.checkRedis(); err != nil {
        return fmt.Errorf("redis health check failed: %w", err)
    }

    return nil
}

func (h *HealthChecker) checkDatabase() error {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    return h.db.PingContext(ctx)
}

func (h *HealthChecker) checkRedis() error {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    return h.redis.Ping(ctx).Err()
}
```

### 10.3 日志配置

```yaml
# 日志配置
Log:
  ServiceName: ai-roleplay
  Mode: file
  Path: /var/log/ai-roleplay
  Level: info
  MaxFiles: 7
  MaxSize: 100
  Compress: true
  KeepDays: 30
  
  # 结构化日志
  Structured: true
  
  # 链路追踪
  Trace:
    Enable: true
    Jaeger:
      Endpoint: http://jaeger:14268/api/traces
```

## 11. 安全策略

### 11.1 数据加密

```go
// common/crypto/crypto.go
package crypto

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "errors"
    "io"
    
    "golang.org/x/crypto/bcrypt"
)

// 密码哈希
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

// AES加密
func Encrypt(plaintext, key string) (string, error) {
    block, err := aes.NewCipher([]byte(createHash(key)))
    if err != nil {
        return "", err
    }

    plainTextByte := []byte(plaintext)
    cfb := cipher.NewCFBEncrypter(block, iv)
    cipherText := make([]byte, len(plainTextByte))
    cfb.XORKeyStream(cipherText, plainTextByte)
    
    return base64.StdEncoding.EncodeToString(cipherText), nil
}

func Decrypt(ciphertext, key string) (string, error) {
    block, err := aes.NewCipher([]byte(createHash(key)))
    if err != nil {
        return "", err
    }

    cipherTextByte, err := base64.StdEncoding.DecodeString(ciphertext)
    if err != nil {
        return "", err
    }

    cfb := cipher.NewCFBDecrypter(block, iv)
    plainTextByte := make([]byte, len(cipherTextByte))
    cfb.XORKeyStream(plainTextByte, cipherTextByte)
    
    return string(plainTextByte), nil
}

func createHash(key string) string {
    hasher := sha256.New()
    hasher.Write([]byte(key))
    return string(hasher.Sum(nil)[:32])
}
```

### 11.2 输入验证

```go
// common/validator/validator.go
package validator

import (
    "regexp"
    "strings"
    "unicode/utf8"
)

// 用户名验证
func ValidateUsername(username string) error {
    if len(username) < 3 || len(username) > 20 {
        return errors.New("用户名长度必须在3-20个字符之间")
    }
    
    matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
    if !matched {
        return errors.New("用户名只能包含字母、数字和下划线")
    }
    
    return nil
}

// 邮箱验证
func ValidateEmail(email string) error {
    matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
    if !matched {
        return errors.New("邮箱格式不正确")
    }
    
    return nil
}

// 内容长度验证
func ValidateTextLength(text string, maxLength int) error {
    if utf8.RuneCountInString(text) > maxLength {
        return fmt.Errorf("内容长度不能超过%d个字符", maxLength)
    }
    
    return nil
}

// 内容安全检查
func ContentSafetyCheck(content string) error {
    // 敏感词检查
    if containsSensitiveWords(content) {
        return errors.New("内容包含敏感词汇")
    }
    
    // 恶意脚本检查
    if containsScript(content) {
        return errors.New("内容包含恶意脚本")
    }
    
    return nil
}

func containsSensitiveWords(content string) bool {
    // 敏感词列表检查逻辑
    sensitiveWords := []string{"敏感词1", "敏感词2"}
    content = strings.ToLower(content)
    
    for _, word := range sensitiveWords {
        if strings.Contains(content, word) {
            return true
        }
    }
    
    return false
}

func containsScript(content string) bool {
    scriptPatterns := []string{
        `<script\b[^<]*(?:(?!<\/script>)<[^<]*)*<\/script>`,
        `javascript:`,
        `on\w+\s*=`,
    }
    
    for _, pattern := range scriptPatterns {
        matched, _ := regexp.MatchString(pattern, content)
        if matched {
            return true
        }
    }
    
    return false
}
```

## 12. 性能优化

### 12.1 缓存策略

```go
// common/cache/cache.go
package cache

import (
    "context"
    "encoding/json"
    "time"
    
    "github.com/go-redis/redis/v8"
)

type CacheManager struct {
    redis *redis.Client
}

func NewCacheManager(redis *redis.Client) *CacheManager {
    return &CacheManager{redis: redis}
}

// 角色信息缓存
func (c *CacheManager) GetCharacter(ctx context.Context, characterId int64) (*Character, error) {
    key := fmt.Sprintf("character:%d", characterId)
    
    val, err := c.redis.Get(ctx, key).Result()
    if err == redis.Nil {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    
    var character Character
    err = json.Unmarshal([]byte(val), &character)
    return &character, err
}

func (c *CacheManager) SetCharacter(ctx context.Context, character *Character) error {
    key := fmt.Sprintf("character:%d", character.Id)
    
    data, err := json.Marshal(character)
    if err != nil {
        return err
    }
    
    return c.redis.Set(ctx, key, data, 30*time.Minute).Err()
}

// 用户会话缓存
func (c *CacheManager) GetUserSession(ctx context.Context, token string) (*UserSession, error) {
    key := fmt.Sprintf("session:%s", token)
    
    val, err := c.redis.Get(ctx, key).Result()
    if err == redis.Nil {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    
    var session UserSession
    err = json.Unmarshal([]byte(val), &session)
    return &session, err
}

func (c *CacheManager) SetUserSession(ctx context.Context, token string, session *UserSession) error {
    key := fmt.Sprintf("session:%s", token)
    
    data, err := json.Marshal(session)
    if err != nil {
        return err
    }
    
    return c.redis.Set(ctx, key, data, 24*time.Hour).Err()
}

// 对话上下文缓存
func (c *CacheManager) GetConversationContext(ctx context.Context, conversationId int64) ([]Message, error) {
    key := fmt.Sprintf("conversation:context:%d", conversationId)
    
    val, err := c.redis.Get(ctx, key).Result()
    if err == redis.Nil {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    
    var messages []Message
    err = json.Unmarshal([]byte(val), &messages)
    return messages, err
}

func (c *CacheManager) SetConversationContext(ctx context.Context, conversationId int64, messages []Message) error {
    key := fmt.Sprintf("conversation:context:%d", conversationId)
    
    data, err := json.Marshal(messages)
    if err != nil {
        return err
    }
    
    return c.redis.Set(ctx, key, data, 2*time.Hour).Err()
}
```

### 12.2 数据库优化

```sql
-- 数据库索引优化
-- 角色表索引
CREATE INDEX idx_characters_category_status ON characters(category, status);
CREATE INDEX idx_characters_rating ON characters(rating DESC);
CREATE INDEX idx_characters_created_at ON characters(created_at DESC);

-- 对话表索引
CREATE INDEX idx_conversations_user_character ON conversations(user_id, character_id);
CREATE INDEX idx_conversations_last_message_time ON conversations(last_message_time DESC);

-- 消息表索引
CREATE INDEX idx_messages_conversation_created ON messages(conversation_id, created_at);
CREATE INDEX idx_messages_user_created ON messages(user_id, created_at DESC);

-- 查询优化
-- 分页查询优化
SELECT id, name, avatar, description, rating 
FROM characters 
WHERE status = 1 
  AND category = ? 
ORDER BY rating DESC, id DESC 
LIMIT ? OFFSET ?;

-- 对话历史查询优化
SELECT m.id, m.type, m.content, m.created_at
FROM messages m
WHERE m.conversation_id = ?
ORDER BY m.created_at DESC
LIMIT 20;
```

## 13. 测试策略

### 13.1 单元测试

```go
// services/user/internal/logic/loginlogic_test.go
package logic

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestLoginLogic_Login(t *testing.T) {
    // Mock dependencies
    mockUserModel := new(MockUserModel)
    mockUserModel.On("FindOneByUsername", mock.Anything, "testuser").Return(&User{
        Id:           1,
        Username:     "testuser",
        PasswordHash: "$2a$14$...", // bcrypt hash of "password"
        Status:       1,
    }, nil)

    // Create logic instance
    svcCtx := &svc.ServiceContext{
        UserModel: mockUserModel,
    }
    logic := NewLoginLogic(context.Background(), svcCtx)

    // Test case
    req := &types.LoginRequest{
        Username: "testuser",
        Password: "password",
    }

    resp, err := logic.Login(req)

    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, int64(1), resp.UserId)
    assert.NotEmpty(t, resp.Token)
    assert.Equal(t, "testuser", resp.Username)

    mockUserModel.AssertExpectations(t)
}
```

### 13.2 集成测试

```go
// tests/integration/user_test.go
package integration

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/stretchr/testify/assert"
)

func TestUserRegistrationFlow(t *testing.T) {
    // Setup test server
    server := setupTestServer()
    defer server.Close()

    // Test user registration
    regReq := map[string]interface{}{
        "username": "testuser",
        "email":    "test@example.com",
        "password": "password123",
    }
    
    regBody, _ := json.Marshal(regReq)
    regResp := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(regBody))
    req.Header.Set("Content-Type", "application/json")
    
    server.ServeHTTP(regResp, req)
    
    assert.Equal(t, http.StatusOK, regResp.Code)
    
    var regResult map[string]interface{}
    json.Unmarshal(regResp.Body.Bytes(), &regResult)
    
    assert.Equal(t, float64(0), regResult["code"])
    assert.NotNil(t, regResult["data"])
    
    // Test user login
    loginReq := map[string]interface{}{
        "username": "testuser",
        "password": "password123",
    }
    
    loginBody, _ := json.Marshal(loginReq)
    loginResp := httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(loginBody))
    req.Header.Set("Content-Type", "application/json")
    
    server.ServeHTTP(loginResp, req)
    
    assert.Equal(t, http.StatusOK, loginResp.Code)
    
    var loginResult map[string]interface{}
    json.Unmarshal(loginResp.Body.Bytes(), &loginResult)
    
    assert.Equal(t, float64(0), loginResult["code"])
    
    data := loginResult["data"].(map[string]interface{})
    token := data["token"].(string)
    assert.NotEmpty(t, token)
}
```

## 14. 发布计划

### 14.1 MVP版本（第一版）
- ✅ 用户注册登录系统
- ✅ 基础角色管理功能
- ✅ 简单文字对话功能
- ✅ 对话历史存储

### 14.2 增强版本（第二版）
- 🔄 语音识别和合成服务
- 🔄 AI角色智能对话
- 🔄 角色自定义功能
- 🔄 高级缓存策略

### 14.3 完整版本（第三版）
- 📋 实时语音流处理
- 📋 多模态AI交互
- 📋 分布式部署方案
- 📋 大规模用户支持

## 15. 总结

本方案采用 Go-Zero 微服务架构，具有以下优势：

1. **高性能**：Go语言天然高并发，框架优化完善
2. **易扩展**：微服务架构便于水平扩展和功能迭代
3. **高可用**：内置熔断、限流、负载均衡等保障机制
4. **易维护**：代码生成、统一规范降低维护成本
5. **生产就绪**：完整的监控、日志、部署方案

通过合理的架构设计和技术选型，能够快速构建出稳定、高效的AI角色扮演语音交互产品后端服务。

---

*本文档版本：v1.0*  
*创建日期：2025-09-23*  
*作者：后端开发团队* 
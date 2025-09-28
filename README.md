# AI角色扮演语音交互产品

<div align="center">

![AI Agent](https://img.shields.io/badge/AI-Agent-blue)
![Go-Zero](https://img.shields.io/badge/Go--Zero-v1.9.0-green)
![Vue3](https://img.shields.io/badge/Vue3-v3.4.0-brightgreen)
![License](https://img.shields.io/badge/License-MIT-yellow)

一个基于AI的多角色扮演智能对话系统，支持语音交互、角色自定义和实时流式对话

[功能特性](#功能特性) • [技术架构](#技术架构) • [快速开始](#快速开始) • [API文档](#api文档) • [部署指南](#部署指南)

</div>

## 📖 项目介绍

AI角色扮演语音交互产品是一个创新的多模态AI对话系统，用户可以与不同性格的AI角色进行自然对话。系统采用微服务架构，支持文本和语音双模式交互，提供沉浸式的角色扮演体验。

### 🎯 核心亮点

- 🤖 **多角色系统** - 预设哈利·波特、爱因斯坦、苏格拉底等经典角色
- 🎙️ **语音交互** - 支持语音输入输出，真实的对话体验
- ⚡ **实时流式对话** - SSE技术实现流式AI回复，响应更自然
- 🎨 **角色自定义** - 支持创建和编辑自定义AI角色
- 📱 **响应式设计** - 适配桌面端和移动端设备
- 🔄 **对话历史管理** - 智能分组管理历史对话记录

## ✨ 功能特性

### 🎭 角色管理
- **预设角色库** - 包含文学、历史、科学等领域的经典角色
- **角色自定义** - 支持设置角色性格、对话风格、提示词等
- **角色收藏** - 收藏喜爱的角色，快速访问
- **角色分类** - 按类别浏览和搜索角色

### 💬 对话系统
- **实时对话** - 基于SSE的流式AI回复
- **多模态交互** - 支持文本输入和语音输入
- **上下文记忆** - 保持对话连贯性和角色一致性
- **对话导出** - 支持导出对话记录为文本文件

### 🎵 语音功能
- **语音识别** - 支持多种音频格式的语音转文字
- **语音合成** - 将AI回复转换为语音输出
- **语音设置** - 可调节语速、音调等参数

### 📊 数据管理
- **对话历史** - 按时间分组的对话历史记录
- **数据统计** - 对话数量、角色互动等统计信息
- **数据导入导出** - 支持对话数据的备份和迁移

## 🏗️ 技术架构

### 后端架构
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   API Gateway   │────│  Load Balancer  │────│     Nginx       │
│    (8888)       │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │
         ├─────────────────────────────────────────────────────────┐
         │                                                         │
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐      │
│  Chat Service   │  │Character Service│  │ Speech Service  │      │
│    (7001)       │  │    (7002)       │  │    (7005)       │      │
└─────────────────┘  └─────────────────┘  └─────────────────┘      │
         │                    │                    │               │
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐      │
│  User Service   │  │   AI Service    │  │ Storage Service │      │
│    (7003)       │  │    (7004)       │  │    (7006)       │      │
└─────────────────┘  └─────────────────┘  └─────────────────┘      │
         │                    │                    │               │
         └────────────────────┼────────────────────┘               │
                              │                                    │
┌─────────────────────────────┼────────────────────────────────────┘
│                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│  │     MySQL       │  │     Redis       │  │     Etcd        │
│  │   (Database)    │  │    (Cache)      │  │ (Service Disc.) │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘
│
│  ┌─────────────────┐  ┌─────────────────┐
│  │     MinIO       │  │   External AI   │
│  │ (Object Store)  │  │  (DeepSeek/GPT) │
│  └─────────────────┘  └─────────────────┘
└─────────────────────────────────────────────────────────────────
```

### 前端架构
```
┌─────────────────────────────────────────────────────────────────┐
│                        Vue 3 + Vite                            │
├─────────────────┬─────────────────┬─────────────────┬───────────┤
│   Pages/Views   │   Components    │     Stores      │ Services  │
│                 │                 │                 │           │
│ • Home.vue      │ • ChatSidebar   │ • chatStore     │ • chatApi │
│ • Chat.vue      │ • VoiceRecorder │ • characterStore│ • charApi │
│ • History.vue   │ • LoadingStates │ • userStore     │ • aiApi   │
│ • Settings.vue  │ • CharacterCard │ • speechStore   │ • speechApi│
└─────────────────┴─────────────────┴─────────────────┴───────────┘
```

### 核心技术栈

#### 后端技术
- **Go-Zero** - 微服务框架，提供API网关、服务发现、负载均衡
- **MySQL** - 主数据库，存储用户、角色、对话等数据
- **Redis** - 缓存和会话管理
- **Etcd** - 服务注册与发现
- **MinIO** - 对象存储，处理文件上传
- **Docker** - 容器化部署

#### 前端技术
- **Vue 3** - 现代化的前端框架
- **Vite** - 快速的构建工具
- **Pinia** - 状态管理
- **Element Plus** - UI组件库
- **Axios** - HTTP客户端
- **Vue Router** - 路由管理

#### AI集成
- **DeepSeek API** - 主要的AI对话服务
- **腾讯云ASR** - 语音识别服务
- **腾讯云TTS** - 语音合成服务

## 🚀 快速开始

### 环境要求
- Go 1.21+
- Node.js 18+
- MySQL 8.0+
- Redis 6.0+
- Docker & Docker Compose

### 1. 克隆项目
```bash
git clone https://github.com/your-username/AI_Agent.git
cd AI_Agent
```

### 2. 后端部署

#### 使用Docker Compose（推荐）
```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps
```

#### 手动部署
```bash
# 进入后端目录
cd backend

# 启动各个微服务
cd services/chat/api && go run chat.go &
cd services/character/api && go run character.go &
cd services/speech/api && go run speech.go &
cd services/user/api && go run user.go &
cd services/ai/api && go run ai.go &
cd services/storage/api && go run storage.go &

# 启动API网关
cd gateway && go run gateway.go
```

### 3. 前端部署
```bash
# 进入前端目录
cd frontend

# 安装依赖
npm install

# 创建环境配置
cp .env.example .env

# 编辑环境变量
vim .env

# 启动开发服务器
npm run dev

# 或构建生产版本
npm run build
```

### 4. 数据库初始化
```bash
# 导入数据库结构
mysql -u root -p < backend/docs/sql/init.sql

# 导入示例数据
mysql -u root -p < backend/docs/sql/data.sql
```

### 5. 配置文件设置

#### 后端配置示例
```yaml
# backend/services/chat/api/etc/chat.yaml
Name: chat-api
Host: 0.0.0.0
Port: 7001

Mysql:
  DataSource: user:password@tcp(localhost:3306)/ai_roleplay?charset=utf8mb4&parseTime=true

Redis:
  Host: localhost:6379
  Type: node

AI:
  Provider: deepseek
  ApiKey: your-deepseek-api-key
  BaseUrl: https://api.deepseek.com
```

#### 前端配置示例
```env
# frontend/.env
VITE_API_BASE_URL=http://localhost:8888
VITE_CHAT_API_URL=http://localhost:7001
VITE_CHARACTER_API_URL=http://localhost:7002
VITE_SPEECH_API_URL=http://localhost:7005
```

## 📋 API文档

### 主要API端点

#### 聊天服务 (7001)
```
GET  /api/chat/history          # 获取对话历史
GET  /api/chat/before           # 获取分组对话历史
GET  /api/chat/send             # SSE流式对话
POST /api/chat/conversation     # 创建对话
GET  /api/chat/messages         # 获取消息列表
```

#### 角色服务 (7002)
```
GET  /api/character/list        # 获取角色列表
GET  /api/character/:id         # 获取角色详情
POST /api/character             # 创建角色
PUT  /api/character/:id         # 更新角色
GET  /api/character/search      # 搜索角色
```

#### 语音服务 (7005)
```
POST /api/speech/stt            # 语音转文字
POST /api/speech/tts            # 文字转语音
GET  /api/speech/health         # 健康检查
```

### API响应格式
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    // 响应数据
  }
}
```

## 🐳 部署指南

### Docker部署

#### 1. 构建镜像
```bash
# 构建所有服务镜像
make docker-build

# 或单独构建
docker build -t ai-agent/chat-service -f backend/services/chat/Dockerfile .
docker build -t ai-agent/frontend -f frontend/Dockerfile .
```

#### 2. 使用Docker Compose
```yaml
version: '3.8'
services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: your-password
      MYSQL_DATABASE: ai_roleplay
    ports:
      - "3306:3306"
    
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
      
  chat-service:
    image: ai-agent/chat-service
    ports:
      - "7001:7001"
    depends_on:
      - mysql
      - redis
      
  frontend:
    image: ai-agent/frontend
    ports:
      - "3000:80"
    depends_on:
      - chat-service
```

### Kubernetes部署

```bash
# 应用Kubernetes配置
kubectl apply -f k8s/

# 查看部署状态
kubectl get pods -n ai-agent
```

### 生产环境配置

#### Nginx反向代理
```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
    
    location /api/ {
        proxy_pass http://localhost:8888;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 🔧 开发指南

### 后端开发
```bash
# 生成API代码
goctl api go -api services/chat/api/chat.api -dir services/chat/api

# 生成模型代码
goctl model mysql ddl -src docs/sql/schema.sql -dir model

# 运行测试
go test ./...
```

### 前端开发
```bash
# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 代码检查
npm run lint

# 类型检查
npm run type-check
```

### 代码规范
- 后端遵循Go标准代码规范
- 前端使用ESLint + Prettier
- 提交信息遵循Conventional Commits

## 📊 监控和日志

### 日志配置
```yaml
Log:
  ServiceName: chat-api
  Mode: file
  Path: logs
  Level: info
  Compress: true
  KeepDays: 7
```

### 监控指标
- API响应时间
- 服务可用性
- 数据库连接池状态
- Redis缓存命中率
- AI服务调用统计

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交Pull Request

## 📄 许可证

本项目采用MIT许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 🙋‍♂️ 支持

- 📖 分工: 
    小组：说的道理
    产品设计：杨艺
    前端设计：石家豪
    后端设计：陈浩钦
- 🐛 问题反馈: [GitHub Issues](https://github.com/your-username/AI_Agent/issues)

## 🔄 更新日志

### v1.0.0 (2025-09-28)
- ✨ 初始版本发布
- 🎭 支持多角色对话系统
- 🎙️ 集成语音识别和合成
- ⚡ 实现SSE流式对话
- 📱 响应式前端界面
- 🐳 Docker容器化部署

---

<div align="center">
  <p>如果这个项目对您有帮助，请给我们一个⭐</p>
  <p>Made with ❤️ by AI Agent Team</p>
</div>

# AI 角色扮演语音交互产品 - 后端服务

基于 Go-Zero 微服务框架构建的 AI 角色扮演语音交互产品后端系统，为前端提供稳定、高性能的 API 服务。

## 🏗️ 架构设计

### 微服务架构
- **API 网关**：统一入口，路由分发，认证授权
- **用户服务**：用户注册登录、个人信息管理
- **角色服务**：AI角色管理、搜索、自定义
- **对话服务**：对话会话管理、消息存储
- **AI 服务**：大语言模型集成、智能对话生成
- **语音服务**：ASR语音识别、TTS语音合成
- **存储服务**：文件上传下载、静态资源管理

### 技术栈
- **框架**：Go-Zero（微服务框架）
- **数据库**：MySQL 8.0（主数据库）+ Redis 6.0（缓存）
- **消息队列**：NATS/RabbitMQ
- **服务发现**：Etcd
- **AI服务**：OpenAI GPT、Azure OpenAI
- **语音服务**：Azure Speech Services、阿里云、腾讯云
- **监控**：Prometheus + Grafana + Jaeger
- **部署**：Docker + Kubernetes

## 📂 项目结构

```
backend/
├── gateway/                    # API网关
│   ├── etc/                   # 配置文件
│   ├── internal/              # 内部逻辑
│   └── gateway.api            # API定义
├── services/                  # 微服务
│   ├── user/                  # 用户服务
│   │   ├── api/               # HTTP API
│   │   ├── rpc/               # gRPC服务
│   │   └── model/             # 数据模型
│   ├── character/             # 角色服务
│   ├── chat/                  # 对话服务
│   ├── ai/                    # AI服务
│   ├── speech/                # 语音服务
│   └── storage/               # 存储服务
├── common/                    # 公共库
│   ├── middleware/            # 中间件
│   ├── response/              # 统一响应
│   ├── utils/                 # 工具函数
│   └── config/                # 配置管理
├── deploy/                    # 部署配置
│   ├── docker/                # Docker配置
│   ├── k8s/                   # Kubernetes配置
│   └── scripts/               # 部署脚本
├── docs/                      # 文档
│   ├── api/                   # API文档
│   ├── sql/                   # 数据库脚本
│   └── design/                # 设计文档
└── scripts/                   # 工具脚本
```

## 🚀 快速开始

### 1. 环境准备

**基础环境：**
- Go 1.19+
- MySQL 8.0+
- Redis 6.0+
- Docker & Docker Compose

**安装 Go-Zero 工具：**
```bash
go install github.com/zeromicro/go-zero/tools/goctl@latest
```

### 2. 项目初始化

```bash
# 克隆项目
git clone <repository-url>
cd AI_Agent/backend

# 运行初始化脚本
./scripts/init-project.sh

# 安装依赖
go mod tidy
```

### 3. 配置环境

```bash
# 复制环境配置文件
cp .env.example .env

# 编辑配置文件，填入相关参数
vim .env
```

**主要配置项：**
- 数据库连接信息
- Redis连接信息
- JWT密钥
- OpenAI API Key
- 语音服务配置

### 4. 启动基础设施

```bash
# 启动 MySQL、Redis、Etcd、MinIO
make infra

# 等待服务启动完成
docker-compose ps
```

### 5. 初始化数据库

```bash
# 创建数据库表结构和初始数据
make init-db
```

### 6. 生成服务代码

```bash
# 生成 API 代码
make api

# 生成 RPC 代码
make rpc

# 生成数据模型
make model
```

### 7. 构建和启动服务

```bash
# 构建所有服务
make build

# 启动 API 网关
./bin/gateway &

# 启动各个微服务
./bin/user-api &
./bin/character-api &
./bin/chat-api &
./bin/ai-api &
./bin/speech-api &
./bin/storage-api &
```

## 🔧 开发指南

### API 开发

1. **定义 API 接口**
   ```bash
   # 编辑 API 文件
   vim services/user/api/user.api
   
   # 生成代码
   cd services/user/api
   goctl api go -api user.api -dir .
   ```

2. **实现业务逻辑**
   ```bash
   # 编辑 logic 文件
   vim services/user/api/internal/logic/loginlogic.go
   ```

3. **配置路由和中间件**
   ```bash
   # 编辑路由配置
   vim services/user/api/internal/handler/routes.go
   ```

### RPC 开发

1. **定义 Proto 文件**
   ```bash
   # 编辑 proto 文件
   vim services/user/rpc/user.proto
   
   # 生成代码
   cd services/user/rpc
   goctl rpc protoc user.proto --go_out=. --go-grpc_out=. --zrpc_out=.
   ```

2. **实现 RPC 服务**
   ```bash
   # 编辑 logic 文件
   vim services/user/rpc/internal/logic/getuserlogic.go
   ```

### 数据模型

```bash
# 从 SQL 生成模型
goctl model mysql ddl -src user.sql -dir ./model

# 从数据库生成模型
goctl model mysql datasource -url="user:password@tcp(127.0.0.1:3306)/database" -table="users" -dir="./model"
```

## 📊 核心业务流程

### 用户注册登录
1. 前端提交注册/登录信息
2. API网关路由到用户服务
3. 用户服务验证信息并生成JWT
4. 返回用户信息和token

### AI 对话流程
1. 前端发送消息到对话服务
2. 对话服务调用AI服务生成回复
3. AI服务调用大语言模型
4. 保存对话记录并返回结果

### 语音处理流程
1. 前端上传语音数据
2. 语音服务调用ASR识别文字
3. 文字发送到AI服务生成回复
4. AI回复通过TTS合成语音
5. 返回文字和语音结果

## 🔒 安全机制

### 身份认证
- JWT Token认证
- Token刷新机制
- 多端登录控制

### 接口安全
- API频率限制
- 参数验证
- SQL注入防护
- XSS攻击防护

### 数据安全
- 密码加密存储
- 敏感数据加密
- 数据传输HTTPS
- 内容安全检查

## 📈 监控和运维

### 健康检查
```bash
# 服务健康检查
curl http://localhost:8888/api/health

# 服务状态检查
curl http://localhost:8888/api/status
```

### 监控指标
- 请求响应时间
- 错误率统计
- 并发连接数
- 资源使用情况

### 日志管理
- 结构化日志输出
- 日志级别控制
- 日志轮转和清理
- 链路追踪

## 🐳 Docker 部署

### 本地开发环境
```bash
# 启动基础设施
docker-compose up -d

# 构建服务镜像
make docker

# 启动所有服务
docker-compose -f docker-compose.full.yml up -d
```

### 生产环境
```bash
# 使用生产配置
docker-compose -f docker-compose.prod.yml up -d
```

## ☁️ Kubernetes 部署

### 部署到 K8s
```bash
# 部署所有服务
kubectl apply -f deploy/k8s/

# 检查部署状态
kubectl get pods -n ai-roleplay

# 查看服务
kubectl get svc -n ai-roleplay
```

### 水平扩展
```bash
# 扩展用户服务
kubectl scale deployment user-api --replicas=3 -n ai-roleplay

# 扩展AI服务
kubectl scale deployment ai-api --replicas=2 -n ai-roleplay
```

## 🧪 测试

### 单元测试
```bash
# 运行所有测试
make test

# 运行特定服务测试
cd services/user/api && go test ./...

# 测试覆盖率
go test -cover ./...
```

### 集成测试
```bash
# 启动测试环境
docker-compose -f docker-compose.test.yml up -d

# 运行集成测试
go test -tags=integration ./tests/...
```

### 性能测试
```bash
# 使用 wrk 进行压力测试
wrk -t12 -c400 -d30s http://localhost:8888/api/character/list

# 使用 ab 进行基准测试
ab -n 1000 -c 10 http://localhost:8888/api/health
```

## 🔧 常用命令

### 开发命令
```bash
make help          # 查看所有可用命令
make build          # 构建所有服务
make run            # 运行服务
make test           # 运行测试
make clean          # 清理构建文件
make docker         # 构建Docker镜像
```

### 代码生成
```bash
make api            # 生成API代码
make rpc            # 生成RPC代码
make model          # 生成数据模型
```

### 数据库操作
```bash
make init-db        # 初始化数据库
make migrate        # 数据库迁移
make seed           # 插入测试数据
```

## 📚 相关文档

- [Go-Zero 官方文档](https://go-zero.dev/)
- [API 接口文档](./docs/api/)
- [数据库设计文档](./docs/database.md)
- [部署指南](./docs/deployment.md)
- [开发规范](./docs/development.md)

## 🤝 贡献指南

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📝 更新日志

### v1.0.0 (2025-09-23)
- ✅ 初始版本发布
- ✅ 完整的微服务架构
- ✅ 用户认证系统
- ✅ AI对话功能
- ✅ 语音处理功能
- ✅ 角色管理功能

## 📄 许可证

本项目基于 MIT 许可证开源 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 📞 联系方式

- 项目维护者：后端开发团队
- 邮箱：backend-team@example.com
- 问题反馈：[GitHub Issues](https://github.com/your-org/ai-roleplay/issues)

---

> 💡 **提示**：如果在使用过程中遇到问题，请先查看[常见问题文档](./docs/faq.md)或在 Issues 中搜索相关问题。 
#!/bin/bash

# AI角色扮演语音交互产品 - 后端项目初始化脚本
# 使用 go-zero 框架生成项目结构

set -e

echo "🚀 开始初始化 AI 角色扮演语音交互产品后端项目..."

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_ROOT=$(pwd)
echo "📁 项目根目录: $PROJECT_ROOT"

# 检查 go-zero 工具是否安装
check_tools() {
    echo "🔍 检查必要工具..."
    
    if ! command -v goctl &> /dev/null; then
        echo -e "${RED}❌ goctl 未安装，请先安装 go-zero 工具${NC}"
        echo "安装命令: go install github.com/zeromicro/go-zero/tools/goctl@latest"
        exit 1
    fi
    
    if ! command -v protoc &> /dev/null; then
        echo -e "${YELLOW}⚠️  protoc 未安装，RPC 服务生成可能失败${NC}"
        echo "安装命令: brew install protobuf (macOS) 或 apt-get install protobuf-compiler (Ubuntu)"
    fi
    
    echo -e "${GREEN}✅ 工具检查完成${NC}"
}

# 创建项目目录结构
create_directories() {
    echo "📂 创建项目目录结构..."
    
    mkdir -p $PROJECT_ROOT/services/{user,character,chat,ai,speech,storage}/{api,rpc,model}
    mkdir -p $PROJECT_ROOT/common/{middleware,response,utils,config}
    mkdir -p $PROJECT_ROOT/gateway
    mkdir -p $PROJECT_ROOT/deploy/{docker,k8s,scripts}
    mkdir -p $PROJECT_ROOT/docs/{api,design}
    mkdir -p $PROJECT_ROOT/scripts
    
    echo -e "${GREEN}✅ 目录结构创建完成${NC}"
}

# 生成 API 服务
generate_api_services() {
    echo "🔧 生成 API 服务..."
    
    services=("user" "character" "chat" "ai" "speech" "storage")
    
    for service in "${services[@]}"; do
        echo "生成 ${service} API 服务..."
        
        cd $PROJECT_ROOT/services/${service}/api
        
        if [ -f "${service}.api" ]; then
            goctl api go -api ${service}.api -dir . --style go_zero
            echo -e "${GREEN}✅ ${service} API 服务生成完成${NC}"
        else
            echo -e "${YELLOW}⚠️  ${service}.api 文件不存在，跳过生成${NC}"
        fi
    done
    
    cd $PROJECT_ROOT
}

# 生成 RPC 服务
generate_rpc_services() {
    echo "🔧 生成 RPC 服务..."
    
    services=("user" "character" "chat" "ai" "speech" "storage")
    
    for service in "${services[@]}"; do
        echo "生成 ${service} RPC 服务..."
        
        cd $PROJECT_ROOT/services/${service}/rpc
        
        # 创建基础 proto 文件
        if [ ! -f "${service}.proto" ]; then
            cat > ${service}.proto << EOF
syntax = "proto3";

package ${service};

option go_package = "./pb";

// ${service} 服务
service ${service^}Service {
  // 健康检查
  rpc Ping(PingRequest) returns (PingResponse);
}

message PingRequest {
  string ping = 1;
}

message PingResponse {
  string pong = 1;
}
EOF
        fi
        
        # 生成 RPC 代码
        goctl rpc protoc ${service}.proto --go_out=. --go-grpc_out=. --zrpc_out=. --style go_zero
        
        echo -e "${GREEN}✅ ${service} RPC 服务生成完成${NC}"
    done
    
    cd $PROJECT_ROOT
}

# 生成数据模型
generate_models() {
    echo "🗄️  生成数据模型..."
    
    # 创建模型目录
    mkdir -p $PROJECT_ROOT/services/common/model
    
    cd $PROJECT_ROOT/services/common/model
    
    # 用户表模型
    goctl model mysql ddl -src ../../../docs/sql/users.sql -dir . --style go_zero || echo "用户表 SQL 文件不存在"
    
    # 角色表模型
    goctl model mysql ddl -src ../../../docs/sql/characters.sql -dir . --style go_zero || echo "角色表 SQL 文件不存在"
    
    # 对话表模型
    goctl model mysql ddl -src ../../../docs/sql/conversations.sql -dir . --style go_zero || echo "对话表 SQL 文件不存在"
    
    # 消息表模型
    goctl model mysql ddl -src ../../../docs/sql/messages.sql -dir . --style go_zero || echo "消息表 SQL 文件不存在"
    
    echo -e "${GREEN}✅ 数据模型生成完成${NC}"
    
    cd $PROJECT_ROOT
}

# 生成网关
generate_gateway() {
    echo "🌐 生成 API 网关..."
    
    cd $PROJECT_ROOT/gateway
    
    if [ -f "gateway.api" ]; then
        goctl api go -api gateway.api -dir . --style go_zero
        echo -e "${GREEN}✅ API 网关生成完成${NC}"
    else
        echo -e "${YELLOW}⚠️  gateway.api 文件不存在，跳过生成${NC}"
    fi
    
    cd $PROJECT_ROOT
}

# 创建配置文件
create_configs() {
    echo "⚙️  创建配置文件..."
    
    # 创建环境配置
    cat > $PROJECT_ROOT/.env.example << 'EOF'
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=123456
DB_NAME=ai_roleplay

# Redis 配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT 配置
JWT_SECRET=your-jwt-secret-key
JWT_EXPIRE=86400

# OpenAI 配置
OPENAI_API_KEY=your-openai-api-key
OPENAI_BASE_URL=https://api.openai.com/v1
OPENAI_MODEL=gpt-3.5-turbo

# 语音服务配置
AZURE_SPEECH_KEY=your-azure-speech-key
AZURE_SPEECH_REGION=your-region

# 存储配置
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
EOF

    # 创建 Docker 配置
    cat > $PROJECT_ROOT/docker-compose.yml << 'EOF'
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: 123456
      MYSQL_DATABASE: ai_roleplay
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./docs/sql:/docker-entrypoint-initdb.d

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data

  etcd:
    image: quay.io/coreos/etcd:v3.5.0
    environment:
      - ETCD_ADVERTISE_CLIENT_URLS=http://etcd:2379
      - ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379
      - ETCD_INITIAL_CLUSTER_STATE=new
    ports:
      - "2379:2379"

  minio:
    image: minio/minio
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    ports:
      - "9000:9000"
      - "9001:9001"
    volumes:
      - minio_data:/data
    command: server /data --console-address ":9001"

volumes:
  mysql_data:
  redis_data:
  minio_data:
EOF

    echo -e "${GREEN}✅ 配置文件创建完成${NC}"
}

# 创建 Makefile
create_makefile() {
    echo "🔨 创建 Makefile..."
    
    cat > $PROJECT_ROOT/Makefile << 'EOF'
.PHONY: help build run test clean docker api rpc model

# 默认目标
help:
	@echo "AI 角色扮演语音交互产品 - 后端项目"
	@echo ""
	@echo "可用命令:"
	@echo "  make build     - 构建所有服务"
	@echo "  make run       - 运行所有服务"
	@echo "  make test      - 运行测试"
	@echo "  make clean     - 清理构建文件"
	@echo "  make docker    - 构建 Docker 镜像"
	@echo "  make api       - 生成 API 代码"
	@echo "  make rpc       - 生成 RPC 代码"
	@echo "  make model     - 生成数据模型"

# 构建所有服务
build:
	@echo "构建所有服务..."
	cd gateway && go build -o ../bin/gateway .
	cd services/user/api && go build -o ../../../bin/user-api .
	cd services/user/rpc && go build -o ../../../bin/user-rpc .
	cd services/character/api && go build -o ../../../bin/character-api .
	cd services/character/rpc && go build -o ../../../bin/character-rpc .
	cd services/chat/api && go build -o ../../../bin/chat-api .
	cd services/chat/rpc && go build -o ../../../bin/chat-rpc .
	cd services/ai/api && go build -o ../../../bin/ai-api .
	cd services/speech/api && go build -o ../../../bin/speech-api .
	cd services/storage/api && go build -o ../../../bin/storage-api .

# 运行基础设施
infra:
	docker-compose up -d mysql redis etcd minio

# 生成 API 代码
api:
	./scripts/generate-api.sh

# 生成 RPC 代码
rpc:
	./scripts/generate-rpc.sh

# 生成数据模型
model:
	./scripts/generate-model.sh

# 运行测试
test:
	go test ./...

# 清理构建文件
clean:
	rm -rf bin/
	go clean ./...

# 构建 Docker 镜像
docker:
	docker build -t ai-roleplay/gateway -f deploy/docker/Dockerfile.gateway .
	docker build -t ai-roleplay/user-api -f deploy/docker/Dockerfile.user-api .
	docker build -t ai-roleplay/character-api -f deploy/docker/Dockerfile.character-api .
	docker build -t ai-roleplay/chat-api -f deploy/docker/Dockerfile.chat-api .
	docker build -t ai-roleplay/ai-api -f deploy/docker/Dockerfile.ai-api .
	docker build -t ai-roleplay/speech-api -f deploy/docker/Dockerfile.speech-api .
	docker build -t ai-roleplay/storage-api -f deploy/docker/Dockerfile.storage-api .

# 初始化数据库
init-db:
	mysql -h localhost -u root -p123456 ai_roleplay < docs/sql/init.sql
EOF

    echo -e "${GREEN}✅ Makefile 创建完成${NC}"
}

# 创建 SQL 初始化脚本
create_sql_scripts() {
    echo "📊 创建 SQL 初始化脚本..."
    
    mkdir -p $PROJECT_ROOT/docs/sql
    
    # 创建初始化 SQL
    cat > $PROJECT_ROOT/docs/sql/init.sql << 'EOF'
-- AI 角色扮演语音交互产品数据库初始化脚本
SET NAMES utf8mb4;

-- 用户表
CREATE TABLE IF NOT EXISTS `users` (
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
  KEY `idx_email` (`email`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 角色表
CREATE TABLE IF NOT EXISTS `characters` (
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
  KEY `idx_status` (`status`),
  KEY `idx_public` (`is_public`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

-- 对话表
CREATE TABLE IF NOT EXISTS `conversations` (
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
  KEY `idx_status` (`status`),
  KEY `idx_last_message_time` (`last_message_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='对话表';

-- 消息表
CREATE TABLE IF NOT EXISTS `messages` (
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息表';

-- 用户收藏表
CREATE TABLE IF NOT EXISTS `user_favorites` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `character_id` bigint NOT NULL,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_character` (`user_id`, `character_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_character_id` (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户收藏表';

-- 插入初始数据
INSERT IGNORE INTO `characters` (`id`, `name`, `avatar`, `description`, `prompt`, `tags`, `category`, `creator_id`) VALUES
(1, '哈利·波特', '/images/avatars/harry-potter.jpg', '霍格沃茨魔法学校的学生，拥有闪电疤痕的男孩巫师。勇敢、善良，擅长魁地奇运动。', '你是哈利·波特，霍格沃茨的学生。你勇敢善良，有着丰富的魔法世界冒险经历。请用哈利的语气和视角来回答问题。', '["魔法", "勇敢", "冒险", "友谊"]', '经典IP', 0),
(2, '苏格拉底', '/images/avatars/socrates.jpg', '古希腊哲学家，以苏格拉底式问答法闻名。追求智慧与真理，善于启发式教学。', '你是苏格拉底，古希腊的哲学家。你善于通过提问来启发他人思考，追求智慧和真理。请用苏格拉底的方式来对话。', '["哲学", "智慧", "思辨", "教育"]', '历史人物', 0),
(3, '莎士比亚', '/images/avatars/shakespeare.jpg', '英国文艺复兴时期的伟大剧作家和诗人，创作了众多不朽的戏剧和十四行诗。', '你是威廉·莎士比亚，伟大的剧作家和诗人。你富有创造力，语言优美，善于用戏剧性的方式表达。', '["文学", "戏剧", "诗歌", "创作"]', '历史人物', 0),
(4, '爱因斯坦', '/images/avatars/einstein.jpg', '20世纪最伟大的物理学家之一，相对论的提出者，诺贝尔物理学奖获得者。', '你是阿尔伯特·爱因斯坦，著名的物理学家。你善于用简单的方式解释复杂的科学概念，充满好奇心和想象力。', '["科学", "物理", "相对论", "思考"]', '历史人物', 0),
(5, '夏洛克·福尔摩斯', '/images/avatars/sherlock.jpg', '世界著名的咨询侦探，居住在贝克街221B号，擅长演绎推理和观察细节。', '你是夏洛克·福尔摩斯，世界上最优秀的咨询侦探。你善于观察细节，进行逻辑推理，有时显得冷漠但内心正义。', '["推理", "侦探", "观察", "逻辑"]', '经典IP', 0),
(6, '赫敏·格兰杰', '/images/avatars/hermione.jpg', '霍格沃茨最聪明的学生之一，博学多才，热爱读书，是哈利和罗恩的好友。', '你是赫敏·格兰杰，霍格沃茨的优秀学生。你博学多才，逻辑清晰，总是能找到解决问题的方法。', '["魔法", "学霸", "聪明", "正义"]', '经典IP', 0);
EOF

    echo -e "${GREEN}✅ SQL 脚本创建完成${NC}"
}

# 创建 README
create_readme() {
    echo "📖 创建 README 文档..."
    
    cat > $PROJECT_ROOT/README.md << 'EOF'
# AI 角色扮演语音交互产品 - 后端服务

基于 Go-Zero 微服务框架构建的 AI 角色扮演语音交互产品后端系统。

## 🚀 快速开始

### 环境要求

- Go 1.19+
- MySQL 8.0+
- Redis 6.0+
- Docker & Docker Compose

### 安装依赖

```bash
# 安装 go-zero 工具
go install github.com/zeromicro/go-zero/tools/goctl@latest

# 安装项目依赖
go mod tidy
```

### 启动基础设施

```bash
# 启动 MySQL、Redis、Etcd、MinIO
make infra
```

### 初始化数据库

```bash
# 初始化数据库表结构和基础数据
make init-db
```

### 启动服务

```bash
# 构建所有服务
make build

# 启动网关
./bin/gateway

# 启动各个微服务
./bin/user-api &
./bin/character-api &
./bin/chat-api &
./bin/ai-api &
./bin/speech-api &
./bin/storage-api &
```

## 📁 项目结构

```
backend/
├── gateway/                    # API网关
├── services/                   # 微服务
│   ├── user/                   # 用户服务
│   ├── character/              # 角色服务
│   ├── chat/                   # 对话服务
│   ├── ai/                     # AI服务
│   ├── speech/                 # 语音服务
│   └── storage/                # 存储服务
├── common/                     # 公共库
├── deploy/                     # 部署配置
├── docs/                       # 文档
└── scripts/                    # 脚本
```

## 🔧 开发指南

### 生成代码

```bash
# 生成 API 代码
make api

# 生成 RPC 代码
make rpc

# 生成数据模型
make model
```

### 测试

```bash
# 运行测试
make test
```

### 构建

```bash
# 构建所有服务
make build

# 构建 Docker 镜像
make docker
```

## 📚 API 文档

- 网关地址：http://localhost:8888
- API 文档：http://localhost:8888/swagger
- 健康检查：http://localhost:8888/api/health

## 🎯 核心功能

- ✅ 用户注册登录系统
- ✅ 角色管理和搜索
- ✅ AI 智能对话
- ✅ 语音识别和合成
- ✅ 对话历史管理
- ✅ 文件存储服务

## 🔒 安全特性

- JWT 身份认证
- API 频率限制
- 内容安全检查
- 数据加密存储

## 📈 监控告警

- Prometheus 指标收集
- Grafana 仪表盘
- Jaeger 链路追踪
- ELK 日志分析

## 🐳 Docker 部署

```bash
# 使用 Docker Compose 部署
docker-compose up -d
```

## ☁️ Kubernetes 部署

```bash
# 部署到 Kubernetes
kubectl apply -f deploy/k8s/
```

## 📄 License

MIT License
EOF

    echo -e "${GREEN}✅ README 文档创建完成${NC}"
}

# 主函数
main() {
    echo "🎉 AI 角色扮演语音交互产品后端项目初始化"
    echo "================================================"
    
    check_tools
    create_directories
    create_configs
    create_sql_scripts
    create_makefile
    create_readme
    
    # 如果 API 文件存在，则生成代码
    if [ -f "$PROJECT_ROOT/gateway/gateway.api" ]; then
        generate_gateway
    fi
    
    echo ""
    echo "================================================"
    echo -e "${GREEN}🎉 项目初始化完成！${NC}"
    echo ""
    echo "📋 下一步操作："
    echo "1. 复制 .env.example 为 .env 并配置相关参数"
    echo "2. 启动基础设施：make infra"
    echo "3. 初始化数据库：make init-db"
    echo "4. 生成服务代码：make api && make rpc"
    echo "5. 构建项目：make build"
    echo "6. 启动服务：./bin/gateway"
    echo ""
    echo -e "${YELLOW}💡 提示：请确保已安装 MySQL、Redis 等依赖服务${NC}"
}

# 执行主函数
main "$@" 
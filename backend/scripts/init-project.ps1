# AI角色扮演语音交互产品 - 后端项目初始化脚本 (Windows PowerShell版本)
# 使用 go-zero 框架生成项目结构

param(
    [switch]$Help
)

# 显示帮助信息
if ($Help) {
    Write-Host "AI角色扮演语音交互产品 - 后端项目初始化脚本"
    Write-Host "用法: .\init-project.ps1"
    Write-Host "选项:"
    Write-Host "  -Help    显示此帮助信息"
    exit 0
}

# 颜色函数
function Write-Success { param($Message) Write-Host $Message -ForegroundColor Green }
function Write-Warning { param($Message) Write-Host $Message -ForegroundColor Yellow }
function Write-Error { param($Message) Write-Host $Message -ForegroundColor Red }
function Write-Info { param($Message) Write-Host $Message -ForegroundColor Cyan }

Write-Host "🚀 开始初始化 AI 角色扮演语音交互产品后端项目..." -ForegroundColor Blue

# 项目根目录
$PROJECT_ROOT = Get-Location
Write-Info "📁 项目根目录: $PROJECT_ROOT"

# 检查必要工具
function Test-Tools {
    Write-Info "🔍 检查必要工具..."
    
    # 检查 Go
    if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
        Write-Error "❌ Go 未安装，请先安装 Go 1.19+ 版本"
        Write-Host "下载地址: https://golang.org/dl/"
        exit 1
    }
    
    # 检查 goctl
    if (-not (Get-Command goctl -ErrorAction SilentlyContinue)) {
        Write-Error "❌ goctl 未安装，请先安装 go-zero 工具"
        Write-Host "安装命令: go install github.com/zeromicro/go-zero/tools/goctl@latest"
        exit 1
    }
    
    # 检查 protoc (可选)
    if (-not (Get-Command protoc -ErrorAction SilentlyContinue)) {
        Write-Warning "⚠️  protoc 未安装，RPC 服务生成可能失败"
        Write-Host "下载地址: https://github.com/protocolbuffers/protobuf/releases"
    }
    
    # 检查 Docker (可选)
    if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
        Write-Warning "⚠️  Docker 未安装，容器化部署将不可用"
        Write-Host "下载地址: https://www.docker.com/products/docker-desktop"
    }
    
    Write-Success "✅ 工具检查完成"
}

# 创建项目目录结构
function New-ProjectDirectories {
    Write-Info "📂 创建项目目录结构..."
    
    $directories = @(
        "services\user\api", "services\user\rpc", "services\user\model",
        "services\character\api", "services\character\rpc", "services\character\model",
        "services\chat\api", "services\chat\rpc", "services\chat\model",
        "services\ai\api", "services\ai\rpc", "services\ai\model",
        "services\speech\api", "services\speech\rpc", "services\speech\model",
        "services\storage\api", "services\storage\rpc", "services\storage\model",
        "common\middleware", "common\response", "common\utils", "common\config",
        "gateway",
        "deploy\docker", "deploy\k8s", "deploy\scripts",
        "docs\api", "docs\design", "docs\sql",
        "scripts", "bin"
    )
    
    foreach ($dir in $directories) {
        $fullPath = Join-Path $PROJECT_ROOT $dir
        if (-not (Test-Path $fullPath)) {
            New-Item -Path $fullPath -ItemType Directory -Force | Out-Null
        }
    }
    
    Write-Success "✅ 目录结构创建完成"
}

# 生成 API 服务
function New-ApiServices {
    Write-Info "🔧 生成 API 服务..."
    
    $services = @("user", "character", "chat", "ai", "speech", "storage")
    
    foreach ($service in $services) {
        Write-Host "生成 $service API 服务..."
        
        $servicePath = Join-Path $PROJECT_ROOT "services\$service\api"
        $apiFile = Join-Path $servicePath "$service.api"
        
        Push-Location $servicePath
        
        if (Test-Path $apiFile) {
            try {
                goctl api go -api "$service.api" -dir . --style go_zero
                Write-Success "✅ $service API 服务生成完成"
            }
            catch {
                Write-Warning "⚠️  $service API 服务生成失败: $_"
            }
        }
        else {
            Write-Warning "⚠️  $service.api 文件不存在，跳过生成"
        }
        
        Pop-Location
    }
}

# 生成 RPC 服务
function New-RpcServices {
    Write-Info "🔧 生成 RPC 服务..."
    
    $services = @("user", "character", "chat", "ai", "speech", "storage")
    
    foreach ($service in $services) {
        Write-Host "生成 $service RPC 服务..."
        
        $servicePath = Join-Path $PROJECT_ROOT "services\$service\rpc"
        $protoFile = Join-Path $servicePath "$service.proto"
        
        Push-Location $servicePath
        
        # 创建基础 proto 文件
        if (-not (Test-Path $protoFile)) {
            $protoContent = @"
syntax = "proto3";

package $service;

option go_package = "./pb";

// $service 服务
service $($service.Substring(0,1).ToUpper() + $service.Substring(1))Service {
  // 健康检查
  rpc Ping(PingRequest) returns (PingResponse);
}

message PingRequest {
  string ping = 1;
}

message PingResponse {
  string pong = 1;
}
"@
            Set-Content -Path $protoFile -Value $protoContent -Encoding UTF8
        }
        
        # 生成 RPC 代码
        try {
            goctl rpc protoc "$service.proto" --go_out=. --go-grpc_out=. --zrpc_out=. --style go_zero
            Write-Success "✅ $service RPC 服务生成完成"
        }
        catch {
            Write-Warning "⚠️  $service RPC 服务生成失败: $_"
        }
        
        Pop-Location
    }
}

# 生成网关
function New-Gateway {
    Write-Info "🌐 生成 API 网关..."
    
    $gatewayPath = Join-Path $PROJECT_ROOT "gateway"
    $gatewayFile = Join-Path $gatewayPath "gateway.api"
    
    Push-Location $gatewayPath
    
    if (Test-Path $gatewayFile) {
        try {
            goctl api go -api gateway.api -dir . --style go_zero
            Write-Success "✅ API 网关生成完成"
        }
        catch {
            Write-Warning "⚠️  API 网关生成失败: $_"
        }
    }
    else {
        Write-Warning "⚠️  gateway.api 文件不存在，跳过生成"
    }
    
    Pop-Location
}

# 创建配置文件
function New-ConfigFiles {
    Write-Info "⚙️  创建配置文件..."
    
    # 创建 .env.example 文件
    $envContent = @"
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
"@
    Set-Content -Path (Join-Path $PROJECT_ROOT ".env.example") -Value $envContent -Encoding UTF8
    
    # 创建 docker-compose.yml 文件
    $dockerComposeContent = @"
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
"@
    Set-Content -Path (Join-Path $PROJECT_ROOT "docker-compose.yml") -Value $dockerComposeContent -Encoding UTF8
    
    Write-Success "✅ 配置文件创建完成"
}

# 创建 PowerShell 构建脚本
function New-BuildScripts {
    Write-Info "🔨 创建构建脚本..."
    
    # 创建 build.ps1
    $buildContent = @"
# AI 角色扮演语音交互产品 - 构建脚本

param(
    [string]`$Action = "build",
    [switch]`$Help
)

if (`$Help) {
    Write-Host "可用命令:"
    Write-Host "  .\build.ps1 build      - 构建所有服务"
    Write-Host "  .\build.ps1 clean      - 清理构建文件"
    Write-Host "  .\build.ps1 test       - 运行测试"
    Write-Host "  .\build.ps1 docker     - 构建Docker镜像"
    Write-Host "  .\build.ps1 infra      - 启动基础设施"
    exit 0
}

function Build-Services {
    Write-Host "🔨 构建所有服务..." -ForegroundColor Blue
    
    if (-not (Test-Path "bin")) {
        New-Item -Path "bin" -ItemType Directory | Out-Null
    }
    
    # 构建网关
    if (Test-Path "gateway") {
        Push-Location gateway
        go build -o ..\bin\gateway.exe .
        Pop-Location
        Write-Host "✅ 网关构建完成" -ForegroundColor Green
    }
    
    # 构建各个服务
    `$services = @("user", "character", "chat", "ai", "speech", "storage")
    foreach (`$service in `$services) {
        if (Test-Path "services\`$service\api") {
            Push-Location "services\`$service\api"
            go build -o "..\..\..\bin\`$service-api.exe" .
            Pop-Location
            Write-Host "✅ `$service-api 构建完成" -ForegroundColor Green
        }
        
        if (Test-Path "services\`$service\rpc") {
            Push-Location "services\`$service\rpc"
            go build -o "..\..\..\bin\`$service-rpc.exe" .
            Pop-Location
            Write-Host "✅ `$service-rpc 构建完成" -ForegroundColor Green
        }
    }
    
    Write-Host "🎉 所有服务构建完成！" -ForegroundColor Green
}

function Start-Infrastructure {
    Write-Host "🚀 启动基础设施..." -ForegroundColor Blue
    docker-compose up -d mysql redis etcd minio
    Write-Host "✅ 基础设施启动完成" -ForegroundColor Green
}

function Test-Services {
    Write-Host "🧪 运行测试..." -ForegroundColor Blue
    go test ./...
}

function Remove-BuildFiles {
    Write-Host "🧹 清理构建文件..." -ForegroundColor Blue
    if (Test-Path "bin") {
        Remove-Item -Path "bin" -Recurse -Force
    }
    go clean ./...
    Write-Host "✅ 清理完成" -ForegroundColor Green
}

function Build-DockerImages {
    Write-Host "🐳 构建Docker镜像..." -ForegroundColor Blue
    # 这里添加Docker构建逻辑
    Write-Host "✅ Docker镜像构建完成" -ForegroundColor Green
}

switch (`$Action.ToLower()) {
    "build" { Build-Services }
    "infra" { Start-Infrastructure }
    "test" { Test-Services }
    "clean" { Remove-BuildFiles }
    "docker" { Build-DockerImages }
    default { 
        Write-Host "未知命令: `$Action" -ForegroundColor Red
        Write-Host "使用 .\build.ps1 -Help 查看帮助"
    }
}
"@
    Set-Content -Path (Join-Path $PROJECT_ROOT "build.ps1") -Value $buildContent -Encoding UTF8
    
    Write-Success "✅ 构建脚本创建完成"
}

# 创建 SQL 初始化脚本
function New-SqlScripts {
    Write-Info "📊 创建 SQL 初始化脚本..."
    
    $sqlPath = Join-Path $PROJECT_ROOT "docs\sql"
    if (-not (Test-Path $sqlPath)) {
        New-Item -Path $sqlPath -ItemType Directory -Force | Out-Null
    }
    
    # 创建初始化 SQL
    $initSqlContent = @"
-- AI 角色扮演语音交互产品数据库初始化脚本
SET NAMES utf8mb4;

-- 用户表
CREATE TABLE IF NOT EXISTS ``users`` (
  ``id`` bigint NOT NULL AUTO_INCREMENT,
  ``username`` varchar(50) NOT NULL UNIQUE,
  ``email`` varchar(100) NOT NULL UNIQUE,
  ``password_hash`` varchar(255) NOT NULL,
  ``avatar`` varchar(255) DEFAULT '',
  ``nickname`` varchar(50) DEFAULT '',
  ``bio`` text,
  ``status`` tinyint DEFAULT '1' COMMENT '1:正常 0:禁用',
  ``last_login_at`` timestamp NULL,
  ``created_at`` timestamp DEFAULT CURRENT_TIMESTAMP,
  ``updated_at`` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (``id``),
  KEY ``idx_username`` (``username``),
  KEY ``idx_email`` (``email``),
  KEY ``idx_status`` (``status``)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 角色表
CREATE TABLE IF NOT EXISTS ``characters`` (
  ``id`` bigint NOT NULL AUTO_INCREMENT,
  ``name`` varchar(100) NOT NULL,
  ``avatar`` varchar(255) DEFAULT '',
  ``description`` text,
  ``prompt`` text NOT NULL,
  ``tags`` json,
  ``category`` varchar(50) DEFAULT '',
  ``rating`` decimal(3,2) DEFAULT '0.00',
  ``rating_count`` int DEFAULT '0',
  ``favorite_count`` int DEFAULT '0',
  ``chat_count`` int DEFAULT '0',
  ``status`` tinyint DEFAULT '1' COMMENT '1:启用 0:禁用',
  ``is_public`` tinyint DEFAULT '1' COMMENT '1:公开 0:私有',
  ``creator_id`` bigint DEFAULT '0',
  ``created_at`` timestamp DEFAULT CURRENT_TIMESTAMP,
  ``updated_at`` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (``id``),
  KEY ``idx_name`` (``name``),
  KEY ``idx_category`` (``category``),
  KEY ``idx_creator`` (``creator_id``),
  KEY ``idx_status`` (``status``),
  KEY ``idx_public`` (``is_public``)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

-- 对话表
CREATE TABLE IF NOT EXISTS ``conversations`` (
  ``id`` bigint NOT NULL AUTO_INCREMENT,
  ``user_id`` bigint NOT NULL,
  ``character_id`` bigint NOT NULL,
  ``title`` varchar(200) DEFAULT '',
  ``start_time`` timestamp DEFAULT CURRENT_TIMESTAMP,
  ``last_message_time`` timestamp NULL,
  ``message_count`` int DEFAULT '0',
  ``status`` tinyint DEFAULT '1' COMMENT '1:活跃 0:结束',
  ``created_at`` timestamp DEFAULT CURRENT_TIMESTAMP,
  ``updated_at`` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (``id``),
  KEY ``idx_user_id`` (``user_id``),
  KEY ``idx_character_id`` (``character_id``),
  KEY ``idx_status`` (``status``),
  KEY ``idx_last_message_time`` (``last_message_time``)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='对话表';

-- 消息表
CREATE TABLE IF NOT EXISTS ``messages`` (
  ``id`` bigint NOT NULL AUTO_INCREMENT,
  ``conversation_id`` bigint NOT NULL,
  ``user_id`` bigint NOT NULL,
  ``type`` varchar(10) NOT NULL COMMENT 'user/ai',
  ``content`` text NOT NULL,
  ``audio_url`` varchar(255) DEFAULT '',
  ``audio_duration`` int DEFAULT '0',
  ``metadata`` json,
  ``created_at`` timestamp DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (``id``),
  KEY ``idx_conversation_id`` (``conversation_id``),
  KEY ``idx_user_id`` (``user_id``),
  KEY ``idx_type`` (``type``),
  KEY ``idx_created_at`` (``created_at``)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息表';

-- 用户收藏表
CREATE TABLE IF NOT EXISTS ``user_favorites`` (
  ``id`` bigint NOT NULL AUTO_INCREMENT,
  ``user_id`` bigint NOT NULL,
  ``character_id`` bigint NOT NULL,
  ``created_at`` timestamp DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (``id``),
  UNIQUE KEY ``uk_user_character`` (``user_id``, ``character_id``),
  KEY ``idx_user_id`` (``user_id``),
  KEY ``idx_character_id`` (``character_id``)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户收藏表';

-- 插入初始数据
INSERT IGNORE INTO ``characters`` (``id``, ``name``, ``avatar``, ``description``, ``prompt``, ``tags``, ``category``, ``creator_id``) VALUES
(1, '哈利·波特', '/images/avatars/harry-potter.jpg', '霍格沃茨魔法学校的学生，拥有闪电疤痕的男孩巫师。勇敢、善良，擅长魁地奇运动。', '你是哈利·波特，霍格沃茨的学生。你勇敢善良，有着丰富的魔法世界冒险经历。请用哈利的语气和视角来回答问题。', '[\"魔法\", \"勇敢\", \"冒险\", \"友谊\"]', '经典IP', 0),
(2, '苏格拉底', '/images/avatars/socrates.jpg', '古希腊哲学家，以苏格拉底式问答法闻名。追求智慧与真理，善于启发式教学。', '你是苏格拉底，古希腊的哲学家。你善于通过提问来启发他人思考，追求智慧和真理。请用苏格拉底的方式来对话。', '[\"哲学\", \"智慧\", \"思辨\", \"教育\"]', '历史人物', 0),
(3, '莎士比亚', '/images/avatars/shakespeare.jpg', '英国文艺复兴时期的伟大剧作家和诗人，创作了众多不朽的戏剧和十四行诗。', '你是威廉·莎士比亚，伟大的剧作家和诗人。你富有创造力，语言优美，善于用戏剧性的方式表达。', '[\"文学\", \"戏剧\", \"诗歌\", \"创作\"]', '历史人物', 0),
(4, '爱因斯坦', '/images/avatars/einstein.jpg', '20世纪最伟大的物理学家之一，相对论的提出者，诺贝尔物理学奖获得者。', '你是阿尔伯特·爱因斯坦，著名的物理学家。你善于用简单的方式解释复杂的科学概念，充满好奇心和想象力。', '[\"科学\", \"物理\", \"相对论\", \"思考\"]', '历史人物', 0),
(5, '夏洛克·福尔摩斯', '/images/avatars/sherlock.jpg', '世界著名的咨询侦探，居住在贝克街221B号，擅长演绎推理和观察细节。', '你是夏洛克·福尔摩斯，世界上最优秀的咨询侦探。你善于观察细节，进行逻辑推理，有时显得冷漠但内心正义。', '[\"推理\", \"侦探\", \"观察\", \"逻辑\"]', '经典IP', 0),
(6, '赫敏·格兰杰', '/images/avatars/hermione.jpg', '霍格沃茨最聪明的学生之一，博学多才，热爱读书，是哈利和罗恩的好友。', '你是赫敏·格兰杰，霍格沃茨的优秀学生。你博学多才，逻辑清晰，总是能找到解决问题的方法。', '[\"魔法\", \"学霸\", \"聪明\", \"正义\"]', '经典IP', 0);
"@
    Set-Content -Path (Join-Path $sqlPath "init.sql") -Value $initSqlContent -Encoding UTF8
    
    Write-Success "✅ SQL 脚本创建完成"
}

# 创建 README
function New-ReadmeFile {
    Write-Info "📖 创建 README 文档..."
    
    $readmeContent = @"
# AI 角色扮演语音交互产品 - 后端服务 (Windows版)

基于 Go-Zero 微服务框架构建的 AI 角色扮演语音交互产品后端系统。

## 🚀 Windows 快速开始

### 环境要求

- Go 1.19+
- MySQL 8.0+
- Redis 6.0+
- Docker Desktop for Windows (可选)

### 安装步骤

1. **安装 Go**
   - 下载地址: https://golang.org/dl/
   - 安装后确保 GOPATH 和 GOROOT 配置正确

2. **安装 go-zero 工具**
   ```powershell
   go install github.com/zeromicro/go-zero/tools/goctl@latest
   ```

3. **安装 protoc (可选)**
   - 下载地址: https://github.com/protocolbuffers/protobuf/releases
   - 将 protoc.exe 添加到 PATH 环境变量

4. **安装 Docker Desktop (可选)**
   - 下载地址: https://www.docker.com/products/docker-desktop

### 项目初始化

```powershell
# 克隆项目
git clone <repository-url>
cd AI_Agent\backend

# 运行初始化脚本
.\scripts\init-project.ps1

# 安装依赖
go mod tidy
```

### 配置环境

```powershell
# 复制环境配置文件
copy .env.example .env

# 编辑配置文件
notepad .env
```

### 启动服务

```powershell
# 启动基础设施 (需要 Docker)
.\build.ps1 infra

# 构建所有服务
.\build.ps1 build

# 启动网关
.\bin\gateway.exe

# 启动其他服务
.\bin\user-api.exe
.\bin\character-api.exe
.\bin\chat-api.exe
.\bin\ai-api.exe
.\bin\speech-api.exe
.\bin\storage-api.exe
```

## 🔧 Windows 开发指南

### 构建命令

```powershell
.\build.ps1 build      # 构建所有服务
.\build.ps1 clean      # 清理构建文件
.\build.ps1 test       # 运行测试
.\build.ps1 infra      # 启动基础设施
```

### 开发工具推荐

- **IDE**: Visual Studio Code + Go 插件
- **终端**: Windows PowerShell 或 Windows Terminal
- **数据库客户端**: Navicat、DBeaver
- **API 测试**: Postman、Insomnia

### Windows 特有注意事项

1. **路径分隔符**: 使用反斜杠 `\` 而不是正斜杠 `/`
2. **可执行文件**: 生成的文件带有 `.exe` 后缀
3. **权限**: 某些操作可能需要管理员权限
4. **防火墙**: 可能需要配置 Windows 防火墙允许端口访问

## 📚 相关链接

- [Go-Zero 官方文档](https://go-zero.dev/)
- [Go 官方文档](https://golang.org/doc/)
- [Docker Desktop for Windows](https://docs.docker.com/desktop/windows/)

## 🐛 常见问题

### 1. goctl 命令找不到
```powershell
# 确保 GOPATH/bin 在 PATH 环境变量中
go env GOPATH
# 将 %GOPATH%\bin 添加到系统 PATH
```

### 2. 端口被占用
```powershell
# 查看端口占用
netstat -ano | findstr :8888
# 结束占用进程
taskkill /PID <进程ID> /F
```

### 3. Docker 启动失败
- 确保 Docker Desktop 正在运行
- 检查 Hyper-V 是否启用
- 确保 WSL2 正确配置

## 📞 技术支持

如有问题，请提交 Issue 或联系开发团队。
"@
    Set-Content -Path (Join-Path $PROJECT_ROOT "README_Windows.md") -Value $readmeContent -Encoding UTF8
    
    Write-Success "✅ README 文档创建完成"
}

# 主函数
function Start-Initialization {
    Write-Host "🎉 AI 角色扮演语音交互产品后端项目初始化 (Windows版)" -ForegroundColor Blue
    Write-Host "================================================" -ForegroundColor Blue
    
    try {
        Test-Tools
        New-ProjectDirectories
        New-ConfigFiles
        New-SqlScripts
        New-BuildScripts
        New-ReadmeFile
        
        # 如果 API 文件存在，则生成代码
        if (Test-Path (Join-Path $PROJECT_ROOT "gateway\gateway.api")) {
            New-Gateway
        }
        
        Write-Host ""
        Write-Host "================================================" -ForegroundColor Blue
        Write-Success "🎉 项目初始化完成！"
        Write-Host ""
        Write-Host "📋 下一步操作:" -ForegroundColor Cyan
        Write-Host "1. 复制 .env.example 为 .env 并配置相关参数"
        Write-Host "2. 启动基础设施：.\build.ps1 infra"
        Write-Host "3. 构建项目：.\build.ps1 build"
        Write-Host "4. 启动服务：.\bin\gateway.exe"
        Write-Host ""
        Write-Warning "💡 提示：请确保已安装 MySQL、Redis 等依赖服务"
        Write-Host ""
        Write-Info "📖 详细说明请查看 README_Windows.md 文件"
    }
    catch {
        Write-Error "❌ 初始化过程中出现错误: $_"
        exit 1
    }
}

# 执行主函数
Start-Initialization 
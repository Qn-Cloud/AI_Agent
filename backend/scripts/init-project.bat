@echo off
setlocal enabledelayedexpansion

REM AI角色扮演语音交互产品 - 后端项目初始化脚本 (Windows批处理版本)
REM 使用 go-zero 框架生成项目结构

echo 🚀 开始初始化 AI 角色扮演语音交互产品后端项目...
echo.

REM 设置项目根目录
set PROJECT_ROOT=%cd%
echo 📁 项目根目录: %PROJECT_ROOT%
echo.

REM 检查必要工具
echo 🔍 检查必要工具...

REM 检查 Go
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Go 未安装，请先安装 Go 1.19+ 版本
    echo 下载地址: https://golang.org/dl/
    pause
    exit /b 1
)

REM 检查 goctl
goctl --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ goctl 未安装，请先安装 go-zero 工具
    echo 安装命令: go install github.com/zeromicro/go-zero/tools/goctl@latest
    pause
    exit /b 1
)

REM 检查 protoc (可选)
protoc --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ⚠️  protoc 未安装，RPC 服务生成可能失败
    echo 下载地址: https://github.com/protocolbuffers/protobuf/releases
)

REM 检查 Docker (可选)
docker --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ⚠️  Docker 未安装，容器化部署将不可用
    echo 下载地址: https://www.docker.com/products/docker-desktop
)

echo ✅ 工具检查完成
echo.

REM 创建项目目录结构
echo 📂 创建项目目录结构...

set directories=services\user\api services\user\rpc services\user\model^
 services\character\api services\character\rpc services\character\model^
 services\chat\api services\chat\rpc services\chat\model^
 services\ai\api services\ai\rpc services\ai\model^
 services\speech\api services\speech\rpc services\speech\model^
 services\storage\api services\storage\rpc services\storage\model^
 common\middleware common\response common\utils common\config^
 gateway^
 deploy\docker deploy\k8s deploy\scripts^
 docs\api docs\design docs\sql^
 scripts bin

for %%d in (%directories%) do (
    if not exist "%%d" (
        mkdir "%%d" 2>nul
    )
)

echo ✅ 目录结构创建完成
echo.

REM 创建配置文件
echo ⚙️  创建配置文件...

REM 创建 .env.example 文件
(
echo # 数据库配置
echo DB_HOST=localhost
echo DB_PORT=3306
echo DB_USER=root
echo DB_PASSWORD=123456
echo DB_NAME=ai_roleplay
echo.
echo # Redis 配置
echo REDIS_HOST=localhost
echo REDIS_PORT=6379
echo REDIS_PASSWORD=
echo.
echo # JWT 配置
echo JWT_SECRET=your-jwt-secret-key
echo JWT_EXPIRE=86400
echo.
echo # OpenAI 配置
echo OPENAI_API_KEY=your-openai-api-key
echo OPENAI_BASE_URL=https://api.openai.com/v1
echo OPENAI_MODEL=gpt-3.5-turbo
echo.
echo # 语音服务配置
echo AZURE_SPEECH_KEY=your-azure-speech-key
echo AZURE_SPEECH_REGION=your-region
echo.
echo # 存储配置
echo MINIO_ENDPOINT=localhost:9000
echo MINIO_ACCESS_KEY=minioadmin
echo MINIO_SECRET_KEY=minioadmin
) > .env.example

REM 创建 docker-compose.yml 文件
(
echo version: '3.8'
echo.
echo services:
echo   mysql:
echo     image: mysql:8.0
echo     environment:
echo       MYSQL_ROOT_PASSWORD: 123456
echo       MYSQL_DATABASE: ai_roleplay
echo     ports:
echo       - "3306:3306"
echo     volumes:
echo       - mysql_data:/var/lib/mysql
echo       - ./docs/sql:/docker-entrypoint-initdb.d
echo.
echo   redis:
echo     image: redis:7-alpine
echo     ports:
echo       - "6379:6379"
echo     command: redis-server --appendonly yes
echo     volumes:
echo       - redis_data:/data
echo.
echo   etcd:
echo     image: quay.io/coreos/etcd:v3.5.0
echo     environment:
echo       - ETCD_ADVERTISE_CLIENT_URLS=http://etcd:2379
echo       - ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379
echo       - ETCD_INITIAL_CLUSTER_STATE=new
echo     ports:
echo       - "2379:2379"
echo.
echo   minio:
echo     image: minio/minio
echo     environment:
echo       MINIO_ROOT_USER: minioadmin
echo       MINIO_ROOT_PASSWORD: minioadmin
echo     ports:
echo       - "9000:9000"
echo       - "9001:9001"
echo     volumes:
echo       - minio_data:/data
echo     command: server /data --console-address ":9001"
echo.
echo volumes:
echo   mysql_data:
echo   redis_data:
echo   minio_data:
) > docker-compose.yml

echo ✅ 配置文件创建完成
echo.

REM 创建构建脚本
echo 🔨 创建构建脚本...

REM 创建 build.bat
(
echo @echo off
echo setlocal
echo.
echo if "%%1"=="help" ^(
echo     echo 可用命令:
echo     echo   build.bat build      - 构建所有服务
echo     echo   build.bat clean      - 清理构建文件
echo     echo   build.bat test       - 运行测试
echo     echo   build.bat infra      - 启动基础设施
echo     goto :eof
echo ^)
echo.
echo if "%%1"=="build" ^(
echo     echo 🔨 构建所有服务...
echo     if not exist bin mkdir bin
echo     
echo     REM 构建网关
echo     if exist gateway ^(
echo         cd gateway
echo         go build -o ..\bin\gateway.exe .
echo         cd ..
echo         echo ✅ 网关构建完成
echo     ^)
echo     
echo     REM 构建各个服务
echo     for %%%%s in ^(user character chat ai speech storage^) do ^(
echo         if exist services\%%%%s\api ^(
echo             cd services\%%%%s\api
echo             go build -o ..\..\..\bin\%%%%s-api.exe .
echo             cd ..\..\..
echo             echo ✅ %%%%s-api 构建完成
echo         ^)
echo         if exist services\%%%%s\rpc ^(
echo             cd services\%%%%s\rpc
echo             go build -o ..\..\..\bin\%%%%s-rpc.exe .
echo             cd ..\..\..
echo             echo ✅ %%%%s-rpc 构建完成
echo         ^)
echo     ^)
echo     echo 🎉 所有服务构建完成！
echo ^) else if "%%1"=="infra" ^(
echo     echo 🚀 启动基础设施...
echo     docker-compose up -d mysql redis etcd minio
echo     echo ✅ 基础设施启动完成
echo ^) else if "%%1"=="test" ^(
echo     echo 🧪 运行测试...
echo     go test ./...
echo ^) else if "%%1"=="clean" ^(
echo     echo 🧹 清理构建文件...
echo     if exist bin rmdir /s /q bin
echo     go clean ./...
echo     echo ✅ 清理完成
echo ^) else ^(
echo     echo 未知命令: %%1
echo     echo 使用 build.bat help 查看帮助
echo ^)
) > build.bat

echo ✅ 构建脚本创建完成
echo.

REM 创建 SQL 初始化脚本
echo 📊 创建 SQL 初始化脚本...

(
echo -- AI 角色扮演语音交互产品数据库初始化脚本
echo SET NAMES utf8mb4;
echo.
echo -- 用户表
echo CREATE TABLE IF NOT EXISTS `users` ^(
echo   `id` bigint NOT NULL AUTO_INCREMENT,
echo   `username` varchar^(50^) NOT NULL UNIQUE,
echo   `email` varchar^(100^) NOT NULL UNIQUE,
echo   `password_hash` varchar^(255^) NOT NULL,
echo   `avatar` varchar^(255^) DEFAULT '',
echo   `nickname` varchar^(50^) DEFAULT '',
echo   `bio` text,
echo   `status` tinyint DEFAULT '1' COMMENT '1:正常 0:禁用',
echo   `last_login_at` timestamp NULL,
echo   `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
echo   `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
echo   PRIMARY KEY ^(`id`^),
echo   KEY `idx_username` ^(`username`^),
echo   KEY `idx_email` ^(`email`^),
echo   KEY `idx_status` ^(`status`^)
echo ^) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
echo.
echo -- 角色表
echo CREATE TABLE IF NOT EXISTS `characters` ^(
echo   `id` bigint NOT NULL AUTO_INCREMENT,
echo   `name` varchar^(100^) NOT NULL,
echo   `avatar` varchar^(255^) DEFAULT '',
echo   `description` text,
echo   `prompt` text NOT NULL,
echo   `tags` json,
echo   `category` varchar^(50^) DEFAULT '',
echo   `rating` decimal^(3,2^) DEFAULT '0.00',
echo   `rating_count` int DEFAULT '0',
echo   `favorite_count` int DEFAULT '0',
echo   `chat_count` int DEFAULT '0',
echo   `status` tinyint DEFAULT '1' COMMENT '1:启用 0:禁用',
echo   `is_public` tinyint DEFAULT '1' COMMENT '1:公开 0:私有',
echo   `creator_id` bigint DEFAULT '0',
echo   `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
echo   `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
echo   PRIMARY KEY ^(`id`^),
echo   KEY `idx_name` ^(`name`^),
echo   KEY `idx_category` ^(`category`^),
echo   KEY `idx_creator` ^(`creator_id`^),
echo   KEY `idx_status` ^(`status`^),
echo   KEY `idx_public` ^(`is_public`^)
echo ^) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';
echo.
echo -- 插入初始数据
echo INSERT IGNORE INTO `characters` ^(`id`, `name`, `avatar`, `description`, `prompt`, `tags`, `category`, `creator_id`^) VALUES
echo ^(1, '哈利·波特', '/images/avatars/harry-potter.jpg', '霍格沃茨魔法学校的学生，拥有闪电疤痕的男孩巫师。勇敢、善良，擅长魁地奇运动。', '你是哈利·波特，霍格沃茨的学生。你勇敢善良，有着丰富的魔法世界冒险经历。请用哈利的语气和视角来回答问题。', '["魔法", "勇敢", "冒险", "友谊"]', '经典IP', 0^),
echo ^(2, '苏格拉底', '/images/avatars/socrates.jpg', '古希腊哲学家，以苏格拉底式问答法闻名。追求智慧与真理，善于启发式教学。', '你是苏格拉底，古希腊的哲学家。你善于通过提问来启发他人思考，追求智慧和真理。请用苏格拉底的方式来对话。', '["哲学", "智慧", "思辨", "教育"]', '历史人物', 0^);
) > docs\sql\init.sql

echo ✅ SQL 脚本创建完成
echo.

REM 创建 README
echo 📖 创建 README 文档...

(
echo # AI 角色扮演语音交互产品 - 后端服务 ^(Windows版^)
echo.
echo 基于 Go-Zero 微服务框架构建的 AI 角色扮演语音交互产品后端系统。
echo.
echo ## 🚀 Windows 快速开始
echo.
echo ### 环境要求
echo.
echo - Go 1.19+
echo - MySQL 8.0+
echo - Redis 6.0+
echo - Docker Desktop for Windows ^(可选^)
echo.
echo ### 项目初始化
echo.
echo ```batch
echo REM 运行初始化脚本
echo scripts\init-project.bat
echo.
echo REM 安装依赖
echo go mod tidy
echo ```
echo.
echo ### 启动服务
echo.
echo ```batch
echo REM 启动基础设施
echo build.bat infra
echo.
echo REM 构建所有服务
echo build.bat build
echo.
echo REM 启动网关
echo bin\gateway.exe
echo ```
echo.
echo ## 🔧 开发指南
echo.
echo ```batch
echo build.bat build      # 构建所有服务
echo build.bat clean      # 清理构建文件
echo build.bat test       # 运行测试
echo build.bat infra      # 启动基础设施
echo build.bat help       # 查看帮助
echo ```
echo.
echo ## 📞 技术支持
echo.
echo 如有问题，请提交 Issue 或联系开发团队。
) > README_Windows.md

echo ✅ README 文档创建完成
echo.

echo ================================================
echo 🎉 项目初始化完成！
echo.
echo 📋 下一步操作:
echo 1. 复制 .env.example 为 .env 并配置相关参数
echo 2. 启动基础设施：build.bat infra
echo 3. 构建项目：build.bat build
echo 4. 启动服务：bin\gateway.exe
echo.
echo 💡 提示：请确保已安装 MySQL、Redis 等依赖服务
echo.
echo 📖 详细说明请查看 README_Windows.md 文件
echo.
echo ================================================

pause 
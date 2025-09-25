#!/bin/bash

# AI角色扮演语音交互产品 - 开发环境启动脚本

echo "🚀 启动AI角色扮演语音交互产品开发环境..."

# 检查Node.js和Go环境
check_prerequisites() {
    echo "📋 检查前置条件..."
    
    # 检查Node.js
    if ! command -v node &> /dev/null; then
        echo "❌ Node.js 未安装，请先安装Node.js 16+版本"
        exit 1
    fi
    
    NODE_VERSION=$(node -v | sed 's/v//')
    echo "✅ Node.js版本: $NODE_VERSION"
    
    # 检查Go
    if ! command -v go &> /dev/null; then
        echo "❌ Go 未安装，请先安装Go 1.19+版本"
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}')
    echo "✅ Go版本: $GO_VERSION"
    
    # 检查Docker
    if ! command -v docker &> /dev/null; then
        echo "⚠️  Docker 未安装，部分功能可能无法使用"
    else
        echo "✅ Docker已安装"
    fi
}

# 启动基础设施服务
start_infrastructure() {
    echo "🔧 启动基础设施服务..."
    
    cd backend
    
    # 检查docker-compose文件
    if [ -f "docker-compose.yml" ]; then
        echo "📦 启动MySQL、Redis、Etcd、MinIO..."
        docker-compose up -d
        
        # 等待服务启动
        echo "⏳ 等待服务启动完成..."
        sleep 10
        
        # 检查服务状态
        docker-compose ps
    else
        echo "⚠️  docker-compose.yml 文件不存在，跳过基础设施启动"
    fi
    
    cd ..
}

# 启动后端服务
start_backend() {
    echo "🖥️  准备启动后端服务..."
    
    cd backend
    
    # 检查go.mod文件
    if [ -f "go.mod" ]; then
        echo "📦 安装Go依赖..."
        go mod tidy
        
        echo "🔨 构建后端服务..."
        # 这里可以添加具体的后端启动命令
        # 由于后端是微服务架构，这里只做准备工作
        echo "✅ 后端准备完成"
        echo "💡 请手动启动各个微服务："
        echo "   - 用户服务: go run services/user/api/user.go"
        echo "   - 角色服务: go run services/character/api/character.go"
        echo "   - 聊天服务: go run services/chat/api/chat.go"
        echo "   - 语音服务: go run services/speech/api/speech.go"
        echo "   - AI服务: go run services/ai/api/ai.go"
        echo "   - 网关服务: go run gateway/gateway.go"
    else
        echo "❌ go.mod 文件不存在"
        exit 1
    fi
    
    cd ..
}

# 启动前端服务
start_frontend() {
    echo "🌐 启动前端服务..."
    
    cd frontend
    
    # 检查package.json文件
    if [ -f "package.json" ]; then
        echo "📦 检查并安装前端依赖..."
        
        # 检查node_modules是否存在
        if [ ! -d "node_modules" ]; then
            echo "📥 安装前端依赖..."
            npm install
        else
            echo "✅ 前端依赖已存在"
        fi
        
        # 创建环境配置文件
        if [ ! -f ".env.development" ]; then
            echo "📝 创建开发环境配置文件..."
            cat > .env.development << EOF
# 开发环境配置
VITE_API_BASE_URL=http://localhost:8888
VITE_APP_TITLE=AI角色扮演语音交互产品
VITE_APP_VERSION=1.0.0
EOF
        fi
        
        echo "🚀 启动前端开发服务器..."
        npm run dev &
        FRONTEND_PID=$!
        
        echo "✅ 前端服务启动完成，PID: $FRONTEND_PID"
        echo "🌐 前端地址: http://localhost:3000"
        
    else
        echo "❌ package.json 文件不存在"
        exit 1
    fi
    
    cd ..
}

# 显示服务信息
show_service_info() {
    echo ""
    echo "🎉 开发环境启动完成！"
    echo ""
    echo "📋 服务信息:"
    echo "   前端服务: http://localhost:3000"
    echo "   后端网关: http://localhost:8888"
    echo "   MySQL: localhost:3306 (root/123456)"
    echo "   Redis: localhost:6379"
    echo "   MinIO: http://localhost:9001 (minioadmin/minioadmin)"
    echo ""
    echo "📖 API文档:"
    echo "   网关API: http://localhost:8888/api/health"
    echo ""
    echo "🛠️  开发工具:"
    echo "   停止前端: Ctrl+C"
    echo "   查看容器: docker-compose ps"
    echo "   停止容器: docker-compose down"
    echo ""
    echo "💡 提示:"
    echo "   1. 确保后端微服务都已启动"
    echo "   2. 检查API接口是否正常工作"
    echo "   3. 查看浏览器控制台是否有错误"
    echo ""
}

# 清理函数
cleanup() {
    echo ""
    echo "🧹 清理资源..."
    
    # 停止前端服务
    if [ ! -z "$FRONTEND_PID" ]; then
        kill $FRONTEND_PID 2>/dev/null
        echo "🔴 前端服务已停止"
    fi
    
    echo "👋 退出开发环境"
    exit 0
}

# 设置信号处理
trap cleanup SIGINT SIGTERM

# 主流程
main() {
    check_prerequisites
    start_infrastructure
    start_backend
    start_frontend
    show_service_info
    
    # 保持脚本运行
    echo "按 Ctrl+C 停止服务"
    while true; do
        sleep 1
    done
}

# 运行主流程
main

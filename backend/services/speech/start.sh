#!/bin/bash

# 语音服务启动脚本
SERVICE_NAME="speech-api"
API_DIR="./api"
CONFIG_FILE="./api/etc/speech.yaml"

echo "🎤 启动AI角色扮演 - 语音服务"
echo "========================================"

# 检查配置文件
if [ ! -f "$CONFIG_FILE" ]; then
    echo "❌ 配置文件不存在: $CONFIG_FILE"
    exit 1
fi

# 检查端口是否被占用
PORT=$(grep -E "^Port:" $CONFIG_FILE | awk '{print $2}')
if [ -n "$PORT" ]; then
    if lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null ; then
        echo "⚠️  端口 $PORT 已被占用，正在尝试停止..."
        pkill -f "$SERVICE_NAME"
        sleep 2
    fi
fi

# 进入API目录
cd $API_DIR

# 构建并启动服务
echo "🔨 构建服务..."
go build -o $SERVICE_NAME speech.go

if [ $? -ne 0 ]; then
    echo "❌ 构建失败"
    exit 1
fi

echo "🚀 启动服务..."
echo "📍 配置文件: $CONFIG_FILE"
echo "🌐 监听端口: $PORT"
echo "📚 API文档: http://localhost:$PORT/api/speech/health"
echo "========================================"

# 启动服务
./$SERVICE_NAME -f etc/speech.yaml

echo "✅ 语音服务已启动" 
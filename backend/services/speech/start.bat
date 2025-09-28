@echo off
chcp 65001 > nul

REM 语音服务启动脚本 (Windows)
set SERVICE_NAME=speech-api.exe
set API_DIR=.\api
set CONFIG_FILE=.\api\etc\speech.yaml

echo 🎤 启动AI角色扮演 - 语音服务
echo ========================================

REM 检查配置文件
if not exist "%CONFIG_FILE%" (
    echo ❌ 配置文件不存在: %CONFIG_FILE%
    pause
    exit /b 1
)

REM 进入API目录
cd %API_DIR%

REM 构建服务
echo 🔨 构建服务...
go build -o %SERVICE_NAME% speech.go

if %errorlevel% neq 0 (
    echo ❌ 构建失败
    pause
    exit /b 1
)

echo 🚀 启动服务...
echo 📍 配置文件: %CONFIG_FILE%
echo 🌐 监听端口: 7005
echo 📚 API文档: http://localhost:7005/api/speech/health
echo ========================================

REM 启动服务
%SERVICE_NAME% -f etc/speech.yaml

echo ✅ 语音服务已启动
pause 
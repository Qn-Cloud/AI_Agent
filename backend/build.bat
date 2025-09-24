@echo off
setlocal

if "%1"=="help" (
    echo 可用命令:
    echo   build.bat build      - 构建所有服�?
  build.bat clean      - 清理构建文件
    echo   build.bat test       - 运行测试
    echo   build.bat infra      - 启动基础设施
    goto :eof
)

if "%1"=="build" (
    echo 🔨 构建所有服�?..
    if not exist bin mkdir bin
ECHO ���ڹر�״̬��
    REM 构建网关
    if exist gateway (
        cd gateway
        go build -o ..\bin\gateway.exe .
        cd ..
        echo �?网关构建完成
    )
ECHO ���ڹر�״̬��
    REM 构建各个服务
    for %%s in (user character chat ai speech storage) do (
        if exist services\%%s\api (
            cd services\%%s\api
            go build -o ..\..\..\bin\%%s-api.exe .
            cd ..\..\..
            echo �?%%s-api 构建完成
        )
        if exist services\%%s\rpc (
            cd services\%%s\rpc
            go build -o ..\..\..\bin\%%s-rpc.exe .
            cd ..\..\..
            echo �?%%s-rpc 构建完成
        )
    )
    echo 🎉 所有服务构建完成！
) else if "%1"=="infra" (
    echo 🚀 启动基础设施...
    docker-compose up -d mysql redis etcd minio
    echo �?基础设施启动完成
) else if "%1"=="test" (
    echo 🧪 运行测试...
    go test ./...
) else if "%1"=="clean" (
    echo 🧹 清理构建文件...
    if exist bin rmdir /s /q bin
    go clean ./...
    echo �?清理完成
) else (
    echo 未知命令: %1
    echo 使用 build.bat help 查看帮助
)

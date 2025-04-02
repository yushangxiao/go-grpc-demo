@echo off
setlocal
cd ..
echo ======================================
echo gRPC 构建和运行脚本
echo ======================================


echo [1/4] 生成 Protocol Buffers 代码...
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/*.proto
if %ERRORLEVEL% neq 0 (
    echo 错误: Protocol Buffers 代码生成失败!
    pause
    exit /b 1
)
echo 成功: Protocol Buffers 代码已生成

echo [2/4] 构建服务端...
if not exist "bin" mkdir bin
go build -o bin/server.exe ./server
if %ERRORLEVEL% neq 0 (
    echo 错误: 服务端构建失败!
    pause
    exit /b 1
)
echo 成功: 服务端已构建

echo [3/4] 构建客户端...
go build -o bin/client.exe ./client
if %ERRORLEVEL% neq 0 (
    echo 错误: 客户端构建失败!
    pause
    exit /b 1
)
echo 成功: 客户端已构建

echo [4/4] 启动服务端...
echo 服务器正在运行，按 Ctrl+C 停止服务器
echo ======================================
bin\server.exe

echo 服务器已停止
pause

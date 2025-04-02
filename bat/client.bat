@echo off
setlocal
cd ..
echo ======================================
echo 运行 Hello 客户端
echo ======================================


:: 检查客户端可执行文件是否存在
if not exist "bin\client.exe" (
    echo 错误: 客户端可执行文件未找到!
    echo 请先运行 build_and_run.bat 构建项目
    pause
    exit /b 1
)

:: 运行客户端
bin\client.exe


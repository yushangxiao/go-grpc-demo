# gRPC 示例 - Go

这是一个简单的 Go gRPC 示例项目，包含服务端和客户端实现。

## 项目结构

```
.
├── proto/         # Protocol Buffers 定义
├── server/        # gRPC 服务端
├── client/        # gRPC 客户端
├── bin/           # 编译后的可执行文件
├── Makefile       # 构建脚本
├── bat/*.bat          # Windows批处理脚本
└── README.md      # 说明文档
```

## 准备工作

1. 安装 Go (1.16+)
2. 安装 Protocol Buffers 编译器 (protoc)
3. 安装 Go 的 Protocol Buffers 插件:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

4. 安装依赖:

```bash
go mod tidy
```

## Windows 一键操作脚本

本项目提供了几个Windows批处理脚本，方便快速构建和运行：

### 构建和启动服务器

双击运行 `build_and_run.bat`，该脚本将：
1. 生成Protocol Buffers代码
2. 构建服务端
3. 构建客户端
4. 启动服务端

### 运行客户端


### 依赖问题解决

如果在构建过程中遇到类似 "missing go.sum entry for module" 的错误，请执行以下命令来解决：

```bash
# 更新并同步所有依赖项
go mod tidy

# 如果上述命令不能解决问题，可以尝试手动获取依赖
go get google.golang.org/grpc
go get google.golang.org/protobuf
```

Windows 环境如果遇到依赖问题时，可以尝试执行：

```bash
go mod download
go mod verify
```

### gRPC 版本兼容性问题

如果遇到以下错误：
```
undefined: grpc.SupportPackageIsVersion9
undefined: grpc.StaticMethod
```

这是由于生成的代码与当前安装的 gRPC 版本不兼容导致的。解决方法有两种：

1. 更新 gRPC 到最新版本（推荐）：
   ```bash
   go get google.golang.org/grpc@latest
   ```

2. 或者修改生成的代码以适应当前版本：
   - 在 `proto/hello_grpc.pb.go` 文件中将 `SupportPackageIsVersion9` 更改为 `SupportPackageIsVersion7`
   - 移除 `grpc.StaticMethod()` 调用

## 手动构建与运行

### 生成 gRPC 代码

```bash
make proto
```

Windows 环境下，如果无法执行 make 命令，可以直接运行以下命令：

```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/*.proto
```

这会从 proto 文件生成 Go 代码，包括:
- `proto/hello.pb.go`
- `proto/hello_grpc.pb.go`

### 构建服务端和客户端

```bash
make all
```

Windows 环境下的替代命令：

```bash
go build -o bin/server.exe ./server
go build -o bin/client.exe ./client
```

或单独构建:

```bash
make server
make client
```

Windows 环境下分别构建：

```bash
go build -o bin/server.exe ./server
go build -o bin/client.exe ./client
```

### 运行服务端

```bash
./bin/server
```

Windows 环境下：

```bash
bin\server.exe
```

服务端将在端口 50051 上监听连接。

### 运行客户端

基本问候调用:
```bash
./bin/client hello [名字]
```

JSON 处理调用:
```bash
./bin/client json
```

Windows 环境下：

```bash
bin\client.exe hello [名字]
bin\client.exe json
```

如果不提供名字参数，将使用默认值 "世界"。

示例:

```bash
./bin/client hello 张三
```

Windows 环境下：

```bash
bin\client.exe hello 张三
bin\client.exe json
```

## 清理

```bash
make clean
```

Windows 环境下：

```bash
del /Q bin\*.exe
del /Q proto\*.pb.go
```

## 项目扩展建议

1. 添加 TLS 安全连接
2. 实现更多 RPC 方法
3. 添加流式 RPC 示例
4. 添加错误处理和重试逻辑
5. 扩展 JSON 处理以支持更多数据类型和操作

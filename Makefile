.PHONY: proto clean server client all

all: proto server client

proto:
	@echo "生成 gRPC 代码..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/hello.proto

server:
	@echo "构建服务端..."
	go build -o bin/server server/main.go

client:
	@echo "构建客户端..."
	go build -o bin/client client/main.go

clean:
	@echo "清理生成文件..."
	rm -rf bin/
	rm -f proto/*.pb.go

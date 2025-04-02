package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/user/grpc-demo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// 服务实现
type server struct {
	pb.UnimplementedGreetServiceServer
}

// SayHello 实现了 gRPC 服务接口
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("收到请求: %v", in.GetName())
	return &pb.HelloResponse{Greeting: "你好, " + in.GetName()}, nil
}

// ProcessJson 处理 JSON 数据
func (s *server) ProcessJson(ctx context.Context, in *pb.JsonRequest) (*pb.JsonResponse, error) {
	log.Printf("收到 JSON 请求: %v", in.GetJsonData())
	// 验证metadata KEY
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &pb.JsonResponse{
			Success:  false,
			Message:  "无效的请求: 缺少认证信息",
			JsonData: "{}",
		}, nil
	}
	if len(md["key"]) == 0 || md["key"][0] != "sk-123456" {
		return &pb.JsonResponse{
			Success:  false,
			Message:  "无效的请求: 错误的认证信息",
			JsonData: "{}",
		}, nil
	}
	log.Printf("认证信息: %v", md["key"][0])
	// 设置响应头
	header := metadata.New(map[string]string{"x-requesr-id": "uuid-123456"})
	grpc.SendHeader(ctx, header)
	// 验证输入 JSON
	var requestData map[string]interface{}
	if err := json.Unmarshal([]byte(in.GetJsonData()), &requestData); err != nil {
		return &pb.JsonResponse{
			Success:  false,
			Message:  "无效的 JSON 格式: " + err.Error(),
			JsonData: "{}",
		}, nil
	}

	// 处理 JSON 数据 (这里是简单的示例，添加一个处理时间戳)
	requestData["processed"] = true
	requestData["server_message"] = "服务器已成功处理请求"

	// 如果请求中包含 name 字段，添加问候语
	if name, ok := requestData["name"].(string); ok {
		requestData["greeting"] = fmt.Sprintf("你好, %s!", name)
	}

	// 将处理后的数据转换回 JSON
	responseJSON, err := json.Marshal(requestData)
	if err != nil {
		return &pb.JsonResponse{
			Success:  false,
			Message:  "JSON 编码错误: " + err.Error(),
			JsonData: "{}",
		}, nil
	}

	return &pb.JsonResponse{
		Success:  true,
		Message:  "JSON 处理成功",
		JsonData: string(responseJSON),
	}, nil
}

// SayHelloStream 实现服务端流式 RPC
func (s *server) SayHelloStream(req *pb.NumberRequest, stream pb.GreetService_SayHelloStreamServer) error {
	log.Printf("收到来自客户端的数字: %v", req.GetNumber())

	// 发送从0到客户端请求数字的所有数字
	for i := 0; i <= int(req.GetNumber()); i++ {
		// 创建响应
		resp := &pb.NumberResponse{
			Number: int32(i),
		}
		// 通过流发送响应
		if err := stream.Send(resp); err != nil {
			return err
		}
		// 短暂停顿，便于观察
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("无法监听端口: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreetServiceServer(s, &server{})

	fmt.Println("gRPC 服务器已启动，正在监听端口 50051...")

	// 启动服务
	if err := s.Serve(lis); err != nil {
		log.Fatalf("无法服务: %v", err)
	}
}

package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	pb "github.com/user/grpc-demo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	// 建立连接到gRPC服务器
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("无法连接: %v", err)
	}
	defer conn.Close()

	// 创建服务的客户端
	client := pb.NewGreetServiceClient(conn)

	for {
		fmt.Println("\n请选择服务类型:")
		fmt.Println("1. 一元RPC - SayHello")
		fmt.Println("2. json RPC - ProcessJson")
		fmt.Println("3. 服务端流式RPC - SayHelloStream")
		fmt.Println("0. 退出")

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("请输入选项: ")
		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)

		switch option {
		case "1":
			callSayHello(client)
		case "2":
			callProcessJson(client)
		case "3":
			callSayHelloStream(client)
		case "0":
			fmt.Println("再见!")
			return
		default:
			fmt.Println("无效选项，请重试")
		}
	}
}

func callProcessJson(client pb.GreetServiceClient) {
	//构建JSON数据
	jsonData := map[string]interface{}{
		"name": "John Doe",
		"age":  30,
		"city": "New York",
	}
	jsonStr, err := json.Marshal(jsonData)
	if err != nil {
		log.Fatalf("无法编码JSON: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	md := metadata.Pairs("key", "sk-123456")
	ctx = metadata.NewOutgoingContext(ctx, md)
	var header metadata.MD
	r, err := client.ProcessJson(ctx, &pb.JsonRequest{JsonData: string(jsonStr)}, grpc.Header(&header))
	if err != nil {
		log.Fatalf("无法执行ProcessJson: %v", err)
	}
	for k, v := range header {
		log.Printf("响应头: %s: %s", k, v)
	}
	var responseData map[string]interface{}

	// 解析JSON字符串
	if err := json.Unmarshal([]byte(r.GetJsonData()), &responseData); err != nil {
		log.Fatalf("无法解析响应JSON: %v", err)
	}

	// 打印完整响应数据
	log.Printf("服务器响应数据: %v", responseData)

	// 根据响应状态处理
	if r.GetSuccess() {
		log.Printf("处理状态: 成功")

		// 如果responseData中有message字段，则打印
		if message, ok := responseData["message"]; ok {
			log.Printf("响应消息: %v", message)
		}
	} else {
		log.Printf("处理状态: 失败")
		log.Printf("错误信息: %s", r.GetMessage())
	}
}

func callSayHello(client pb.GreetServiceClient) {
	// 获取用户输入
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("请输入您的名字: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	// 联系服务器并打印出它的响应
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("无法执行SayHello: %v", err)
	}

	log.Printf("服务器响应: %s", r.GetGreeting())
}

func callSayHelloStream(client pb.GreetServiceClient) {
	// 获取用户输入
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("请输入一个数字: ")
	numStr, _ := reader.ReadString('\n')
	numStr = strings.TrimSpace(numStr)

	num, err := strconv.Atoi(numStr)
	if err != nil {
		log.Printf("无效的数字: %v", err)
		return
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 调用流式RPC
	stream, err := client.SayHelloStream(ctx, &pb.NumberRequest{Number: int32(num)})
	if err != nil {
		log.Fatalf("调用SayHelloStream时出错: %v", err)
	}

	fmt.Println("服务器流式响应:")
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			// 读取完成
			break
		}
		if err != nil {
			log.Fatalf("接收流时出错: %v", err)
		}

		fmt.Printf("%d\n", resp.GetNumber())
	}
}

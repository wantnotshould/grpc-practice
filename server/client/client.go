// Copyright ©2026 xiayoudi. All rights reserved.
// Author: xiayoudi
// Email: ur@xiaud.com

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	hello_v1 "grpc-practice/grpc_proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	addr := "localhost:8080" // 指定服务端的地址

	// 建立连接
	// insecure.NewCredentials() 表示不启用 SSL/TLS 加密，仅用于本地测试
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("无法创建客户端连接 [%s]: %v", addr, err)
	}

	// 确保在程序退出前关闭连接，释放网络资源
	defer conn.Close()

	// 初始化业务客户端
	// 这个 client 是由 protoc 插件帮我们生成的，它包装了底层的网络调用
	client := hello_v1.NewHelloServiceClient(conn)

	// 发起调用
	// 设置一个超时时间（防止服务端死掉导致客户端无限等待）
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// 构造请求参数，调用远程的 SayHello 函数
	result, err := client.SayHello(ctx, &hello_v1.HelloReq{
		Name:    "Perry",
		Message: "Hello Server!",
	})

	// 处理结果
	if err != nil {
		// 如果调用失败，打印具体的错误原因（如：连接超时、服务端拒绝等）
		log.Fatalf("调用 SayHello 失败: %v", err)
	}

	// 打印服务端返回的响应内容
	fmt.Printf("收到响应: 姓名=%s, 消息=%s\n", result.GetName(), result.GetMessage())
}

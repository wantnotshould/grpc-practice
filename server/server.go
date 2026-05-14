// Copyright ©2026 xiayoudi. All rights reserved.
// Author: xiayoudi
// Email: ur@xiaud.com

package main

import (
	"context"
	"fmt"
	"log"
	"net"

	hello_v1 "grpc-practice/grpc_proto"

	"google.golang.org/grpc"
)

// HelloServer 实现了 hello_grpc.HelloServiceServer 接口
type HelloServer struct {
	// 必须嵌入这个结构体，否则会报错：missing method mustEmbedUnimplemented...
	hello_v1.UnimplementedHelloServiceServer
}

// SayHello 是我们在 .proto 中定义的 RPC 函数的具体实现
func (s *HelloServer) SayHello(ctx context.Context, param *hello_v1.HelloReq) (*hello_v1.HelloResp, error) {
	// 打印接收到的请求
	log.Printf("Receive request from: %s, message: %s", param.Name, param.Message)

	// 业务逻辑：返回响应
	resp := &hello_v1.HelloResp{
		Name:    "Server",
		Message: fmt.Sprintf("Hello %s, I received your message!", param.Name),
	}

	return resp, nil
}

func main() {
	// 监听本地端口
	addr := ":8080"
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", addr, err)
	}

	// 创建 gRPC 服务器实例
	s := grpc.NewServer()

	// 将我们的服务实现注册到 gRPC 服务器中
	// 注意：这里传的是 &HelloServer{} 的指针
	hello_v1.RegisterHelloServiceServer(s, &HelloServer{})

	fmt.Printf("gRPC server is running at %s...\n", addr)

	// 启动服务（此调用是阻塞的）
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

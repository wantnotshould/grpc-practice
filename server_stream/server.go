// Copyright ©2026 xiayoudi. All rights reserved.
// Author: xiayoudi
// Email: ur@xiaud.com

package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	stream_v1 "grpc-practice/stream/stream/v1"

	"google.golang.org/grpc"
)

type StreamServer struct {
	stream_v1.UnimplementedStreamServiceServer
}

type ClientStreamServer struct {
	stream_v1.UnimplementedClientStreamServiceServer
}

// Fun 实现服务端流式接口
// 参数：param 是客户端发的请求，stream 是用来发送数据的“管子”
func (s *StreamServer) Fun(param *stream_v1.Request, stream stream_v1.StreamService_FunServer) error {
	fmt.Printf("收到请求: %s，准备开始发送流数据...\n", param.Name)

	// 模拟发送 5 条数据
	for i := range 5 {
		res := &stream_v1.Response{
			Name: fmt.Sprintf("这是第 %d 条数据，发给 %s", i+1, param.Name),
		}

		// 通过 stream.Send 发送单条响应
		if err := stream.Send(res); err != nil {
			return err
		}

		// 模拟耗时操作，每秒发一条
		time.Sleep(time.Second)
	}

	// 返回 nil 表示流发送完毕，通知客户端结束
	return nil
}

func (s *StreamServer) FileDownload(param *stream_v1.Request, stream stream_v1.StreamService_FileDownloadServer) error {
	// 模拟文件数据（实际开发中这里是 os.Open 打开真实文件）
	fileName := param.Name
	content := []byte("这是文件内容，假设它非常大，我们需要分片传输。" +
		"gRPC 会把这些数据切成小块发送给客户端...")

	// 定义每次发送的分片大小（例如 1KB）
	chunkSize := 1024

	log.Printf("客户端请求下载文件: %s", fileName)

	// 循环切片并发送
	for i := 0; i < len(content); i += chunkSize {
		end := i + chunkSize
		if end > len(content) {
			end = len(content)
		}

		// 构造响应消息
		resp := &stream_v1.FileDownloadResponse{
			FileName: fileName,
			Content:  content[i:end],
		}

		// 发送当前切片
		if err := stream.Send(resp); err != nil {
			return err
		}

		// 模拟网络延迟，方便观察流式效果
		time.Sleep(time.Millisecond * 500)
		log.Printf("已发送数据块，偏移量: %d/%d", end, len(content))
	}

	// 函数返回 nil，gRPC 会自动发送一个特殊的信号给客户端，标识流结束
	log.Println("文件发送完毕")
	return nil
}

// Chat 实现双向流
func (s *StreamServer) Chat(stream stream_v1.StreamService_ChatServer) error {
	for {
		// 接收客户端消息
		req, err := stream.Recv()
		if err == io.EOF {
			return nil // 客户端关闭了发送
		}
		if err != nil {
			return err
		}

		fmt.Printf("[%s]: %s\n", req.User, req.Text)

		// 立即回显/响应客户端
		reply := &stream_v1.ChatMessage{
			User: "系统管理员",
			Text: fmt.Sprintf("已收到 %s 的消息: %s", req.User, req.Text),
		}

		if err := stream.Send(reply); err != nil {
			return err
		}
	}
}

func (s *ClientStreamServer) Upload(stream stream_v1.ClientStreamService_UploadServer) error {
	var fileName string
	var totalSize int

	for {
		// 客户端流模式下，服务端必须调用 Recv() 接收每一块数据
		req, err := stream.Recv()

		// 判断客户端是否发送完毕 (io.EOF)
		if err == io.EOF {
			log.Printf("接收完成，文件名: %s, 总大小: %d", fileName, totalSize)

			// 【关键】客户端流模式，服务端最后必须调用 SendAndClose
			// 返回唯一的一个 Response 并关闭流
			return stream.SendAndClose(&stream_v1.Response{
				Name: fmt.Sprintf("上传成功，共接收 %d 字节", totalSize),
			})
		}

		if err != nil {
			return err
		}

		// 业务逻辑处理
		if fileName == "" {
			fileName = req.FileName
		}
		totalSize += len(req.Content)
	}
}

func main() {
	// 建立监听
	addr := ":8080"
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("端口监听失败: %v", err)
	}

	// 实例化 gRPC Server
	s := grpc.NewServer()

	// 注册服务
	// 注意：这里注册的是 StreamService，不是 SimpleService
	stream_v1.RegisterStreamServiceServer(s, &StreamServer{})
	stream_v1.RegisterClientStreamServiceServer(s, &ClientStreamServer{})

	fmt.Printf("gRPC 流式服务端启动成功，监听端口 %s...\n", addr)

	// 启动服务（阻塞运行）
	if err := s.Serve(listen); err != nil {
		log.Fatalf("运行失败: %v", err)
	}
}

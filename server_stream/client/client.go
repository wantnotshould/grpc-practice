// Copyright ©2026 xiayoudi. All rights reserved.
// Author: xiayoudi
// Email: ur@xiaud.com

package main

import (
	"context"
	"fmt"
	stream_v1 "grpc-practice/stream"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 建立连接 (gRPC 1.x 最新推荐用法 NewClient)
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	// 创建 Client 对象
	streamCli := stream_v1.NewStreamServiceClient(conn)
	uploadCli := stream_v1.NewClientStreamServiceClient(conn)

	// 调用普通的流式推送
	fmt.Println(">>> 普通消息")
	callFun(streamCli)

	fmt.Println("--------------------------")

	// 调用文件下载流
	fmt.Println(">>> 下载（服务端流式）")
	downloadFile(streamCli)

	fmt.Println(">>> 上传（客户端流式）")
	runUpload(uploadCli)

	fmt.Println(">>> 执行场景: 双向流对话")
	runChat(streamCli)
}

// 对应 rpc Fun(Request) returns (stream Response)
func callFun(client stream_v1.StreamServiceClient) {
	stream, err := client.Fun(context.Background(), &stream_v1.Request{Name: "Perry"})
	if err != nil {
		log.Printf("调用 Fun 失败: %v", err)
		return
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("【结束】Fun 消息接收完毕。")
			break
		}
		if err != nil {
			log.Printf("读取 Fun 流出错: %v", err)
			break
		}
		fmt.Printf("收到推送消息: %s\n", res.GetName())
	}
}

// 对应 rpc FileDownload(Request) returns (stream FileDownloadResponse)
func downloadFile(client stream_v1.StreamServiceClient) {
	req := &stream_v1.Request{Name: "movie.mp4"}
	stream, err := client.FileDownload(context.Background(), req)
	if err != nil {
		log.Printf("发起下载失败: %v", err)
		return
	}

	var fullContent []byte
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("\n【结束】文件下载成功，服务端已切断流。")
			break
		}
		if err != nil {
			log.Printf("下载过程中断: %v", err)
			break
		}

		// 将收到的字节分片（Content）拼接到一起
		fullContent = append(fullContent, res.Content...)

		// \r 可以让光标回到行首，实现动态进度条的效果
		fmt.Printf("\r正在下载 [%s] ... 已接收: %d 字节", res.FileName, len(fullContent))
	}
}

func runUpload(client stream_v1.ClientStreamServiceClient) {
	// 获取流对象
	stream, err := client.Upload(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 发送多次数据
	for i := range 5 {
		stream.Send(&stream_v1.FileUploadRequest{
			FileName: "test.txt",
			// Content:  []byte(fmt.Sprintf("第 %d 块数据", i)),
			Content: fmt.Appendf(nil, "第 %d 块数据", i),
		})
	}

	// 【关键】发送完毕并接收服务端那唯一的响应
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("服务端结果:", res.Name)
}

func runChat(client stream_v1.StreamServiceClient) {
	// 开启双向流
	stream, err := client.Chat(context.Background())
	if err != nil {
		log.Fatalf("开启对话失败: %v", err)
	}

	// 用于协调主线程等待
	waitc := make(chan struct{})

	// 开启协程：专门负责“听”（接收服务端的消息）
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("### 服务端关闭了对话 ###")
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("接收失败: %v", err)
			}
			fmt.Printf("收到回复 -> [%s]: %s\n", res.User, res.Text)
		}
	}()

	// 主线程：专门负责“说”（发送消息给服务端）
	msgs := []string{"你好", "在吗？", "gRPC 真好用", "再见"}
	for _, m := range msgs {
		err := stream.Send(&stream_v1.ChatMessage{
			User: "Perry",
			Text: m,
		})
		if err != nil {
			log.Fatalf("发送失败: %v", err)
		}
		time.Sleep(time.Second) // 模拟聊天间隔
	}

	// 发完后，关闭发送端流
	stream.CloseSend()

	// 等待接收协程处理完剩下的数据
	<-waitc
}

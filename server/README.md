```go
func main() {
    listen, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    s := grpc.NewServer()
    server := HelloServer{}
    hello_grpc.RegisterHelloServiceServer(s, &server)

    // 在独立的协程中启动服务
    go func() {
        fmt.Println("gRPC server running on :8080")
        if err := s.Serve(listen); err != nil {
            log.Fatalf("Failed to serve: %v", err)
        }
    }()

    // 监听系统退出信号 (Ctrl+C 或 kill)
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
    <-quit // 阻塞在这里，直到收到信号

    fmt.Println("Shutting down gRPC server...")

    // 优雅关闭：停止接收新请求，等待已有的请求处理完
    s.GracefulStop()

    fmt.Println("Server exited")
}
```
# gRPC-practice

## 准备

[protoc - GitHub](https://github.com/protocolbuffers/protobuf)

```bash
# 只负责数据
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# 只负责通信
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## 生成

```bash
# protoc -I . --go-grpc_out=. ./hello.proto

protoc -I . \
  --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  ./hello.proto

# 同时编译三个文件，这样插件才能处理它们之间的 import 依赖
protoc -I . \
  --go_out=. --go-grpc_out=. \
  --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative \
  common.proto order.proto video.proto

# 去掉 paths=source_relative
# protoc -I . \
#   --go_out=. \
#   --go-grpc_out=. \
#   common.proto order.proto video.proto

# 编译订单业务
# protoc --proto_path=. \
#     --go_out=. --go_opt=paths=source_relative \
#     --go-grpc_out=. --go-grpc_opt=paths=source_relative \
#     api/order/v1/order.proto

# 编译用户业务
# protoc --proto_path=. \
#     --go_out=. --go_opt=paths=source_relative \
#     --go-grpc_out=. --go-grpc_opt=paths=source_relative \
#     api/user/v1/user.proto
```

**决定结构**

`multiple/v1/*.pb.go`

```bash
protoc -I . --go_out=. --go-grpc_out=. *.proto
```

如果你希望 `.proto` 本身就在 `v1` 文件夹里（这是最正规的），请先创建文件夹

```bash
mkdir -p api/multiple/v1
mv *.proto api/multiple/v1/
# 然后在根目录执行 (带上 source_relative)
protoc -I . --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/multiple/v1/*.proto
```
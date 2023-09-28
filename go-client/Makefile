
.PHONY: build
build: user/user.pb.go user/user_grpc.pb.go client.go
	go build -o bin/client client.go

user/user.pb.go user/user_grpc.pb.go: proto/user/user.proto
	protoc -Iproto --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/user/user.proto

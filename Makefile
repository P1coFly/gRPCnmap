BINARY_NAME=grpcnmap

generate:
	protoc -I api api/NetVuln_v1/service.proto --go_out=pkg  --go_opt=paths=source_relative --go-grpc_out=pkg --go-grpc_opt=paths=source_relative
lint:
	golangci-lint run
build:
	go build -o ${BINARY_NAME} cmd/grpc-nmap/main.go
test:
	go test

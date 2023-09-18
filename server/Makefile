build:
	@go build -o bin/go-product

run: build
	@./bin/go-product

dev:
	@go run main.go

proto:
	@protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \ 
	proto/service.proto

.PHONY: proto

test:
	@go test -v ./...
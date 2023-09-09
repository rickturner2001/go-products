build:
	@go build -o bin/go-product

run:
	@./bin/go-product

dev:
	@go run main.go

test:
	@go test -v ./...
build: 
	@go build -o api .

run: build
	@./api

test: 
	@go test ./...

cover: 
	@go test -race -coverprofile=coverage.txt -covermode=atomic ./...
	@go tool cover -func coverage.txt | grep total

cover-html: cover
	@go tool cover -html coverage.txt -o coverage.html

lint:
	@LOG_LEVEL=error golangci-lint run ./...

proto:
	@protoc --proto_path=../protos ../protos/*.proto --go_out=pkg/services --go-grpc_out=pkg/services --experimental_allow_proto3_optional
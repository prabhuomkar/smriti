.PHONY: build
build: 
	@go build -o api --ldflags="-X 'api/internal/models.DefaultVersion=$(VERSION)' -X 'api/internal/models.DefaultGitSHA=$(shell git rev-parse HEAD)'"

.PHONY: run
run: build
	@./api

.PHONY: test
test: 
	@go test ./...

.PHONY: cover
cover: 
	@go test $$(go list ./... | grep -v api/pkg/services) -race -coverprofile=coverage.txt -covermode=atomic ./...
	@go tool cover -func coverage.txt | grep total

.PHONY: cover-html
cover-html: cover
	@go tool cover -html coverage.txt -o coverage.html

.PHONY: lint
lint:
	@LOG_LEVEL=error golangci-lint run ./...

.PHONY: proto
proto:
	@protoc --proto_path=../protos ../protos/*.proto --go_out=pkg/services --go-grpc_out=pkg/services --experimental_allow_proto3_optional
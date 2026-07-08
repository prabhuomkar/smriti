.PHONY: build
build: 
	@go build -o api --ldflags="-X 'api/internal/models.DefaultVersion=$(VERSION)' -X 'api/internal/models.DefaultGitSHA=$(shell git rev-parse HEAD)'"

.PHONY: migrate-up
migrate-up:
	@migrate -path=pkg/database/migrations -database "${DATABASE_URL}" -verbose up

.PHONY: migrate-down
migrate-down:
	@migrate -path=pkg/database/migrations -database "${DATABASE_URL}" -verbose down

.PHONY: run
run: build
	@./api

.PHONY: test
test: 
	@go test ./...

.PHONY: cover
cover: 
	@go test -race -coverprofile=coverage.txt -covermode=atomic -coverpkg=$($$(go list ./... | grep -v api/pkg/services) | tr ' ' ',') $$(go list ./... | grep -v api/pkg/services)
	@go tool cover -func coverage.txt | grep total

.PHONY: cover-html
cover-html: cover
	@go tool cover -html coverage.txt -o coverage.html

.PHONY: lint
lint:
	@LOG_LEVEL=error golangci-lint run ./...

.PHONY: proto
proto:
	@protoc --proto_path=./protos ./protos/*.proto --go_out=pkg/services --go-grpc_out=pkg/services
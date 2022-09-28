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
	@golangci-lint run ./...
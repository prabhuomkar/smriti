build: 
	@go build -o api .

run: build
	@./api

test: 
	@go test ./...

cover: 
	@go test -v -coverprofile cover.out ./...
	@go tool cover -func cover.out | grep total

cover-html: cover
	@go tool cover -html cover.out -o cover.html
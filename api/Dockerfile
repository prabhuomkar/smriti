FROM golang:1.18-alpine as builder
RUN mkdir -p /app
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o api .

FROM alpine
WORKDIR /app
COPY --from=builder /app/api .
EXPOSE 5001
EXPOSE 15001
CMD ["/app/api"]
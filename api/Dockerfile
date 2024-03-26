FROM golang:1.21-alpine as builder
RUN apk update && apk add git
RUN mkdir -p /app
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
ARG VERSION=dev
ARG GITSHA=-
RUN go build -o api --ldflags="-X 'api/internal/models.DefaultVersion=$VERSION' -X 'api/internal/models.DefaultGitSHA=$GITSHA'"

FROM alpine
WORKDIR /app
COPY --from=builder /app/api .
EXPOSE 5001
EXPOSE 15001
CMD ["/app/api"]
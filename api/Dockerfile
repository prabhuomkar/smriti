FROM golang:1.24-alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=dev
ARG GITSHA=-
RUN CGO_ENABLED=0 go build -o api --ldflags="-X 'api/internal/models.DefaultVersion=$VERSION' -X 'api/internal/models.DefaultGitSHA=$GITSHA'"

FROM gcr.io/distroless/static
WORKDIR /app
COPY --from=builder /app/api .
EXPOSE 5001
EXPOSE 15001
CMD ["/app/api"]
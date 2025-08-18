# 使用官方 Go 镜像编译
FROM golang:1.21 AS builder

WORKDIR /app
COPY . .
# 静态编译，避免依赖 glibc
RUN go mod init hello-k8s || true
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

# 使用更小的运行时镜像
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["./app"]

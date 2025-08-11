# 用官方 Go 运行环境作为基础镜像
FROM golang:1.22-alpine

# 安装 git（go mod 可能需要）
RUN apk add --no-cache git

# 设置容器内的工作目录
WORKDIR /app

# 先复制 go.mod 和 go.sum 并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 再复制所有代码
COPY . .

# 构建 Go 程序
RUN go build -o employee-crud main.go

# 容器启动时运行的命令
CMD ["./employee-crud"]

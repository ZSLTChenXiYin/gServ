# 使用 Go 1.24.11 作为构建环境
FROM golang:1.24.11-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的构建工具
RUN apk add --no-cache git ca-certificates tzdata

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o gServ main.go

# 使用精简的 Alpine 镜像作为运行环境
FROM alpine:3.20

# 设置工作目录
WORKDIR /app

# 安装必要的运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 从构建阶段复制可执行文件
COPY --from=builder /app/gServ /app/gServ

# 复制配置文件模板
COPY docker.gserv.conf.yaml /app/gserv.conf.yaml

# 复制资源文件夹
COPY res /app/res

# 设置环境变量
ENV TZ=Asia/Shanghai

# 暴露端口
EXPOSE 8080 9090

# 启动命令
CMD ["/app/gServ"]
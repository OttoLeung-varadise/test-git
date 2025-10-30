# 第一阶段：构建（编译 Go 程序）
FROM golang:1.25-alpine AS builder

# 启用 Go Modules
ENV GO111MODULE=on
# 设置 GOPROXY 加速依赖下载（可选，国内推荐）
ENV GOPROXY=https://goproxy.cn,direct
# 设置工作目录（容器内）
WORKDIR /app

# 复制依赖文件，优先下载依赖（利用 Docker 缓存）
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码到容器内
COPY . .

# 编译程序（静态链接，适配 alpine 环境）
# 说明：
# - CGO_ENABLED=0：禁用 CGO，避免依赖系统动态库
# - GOOS=linux：目标系统为 Linux
# - -a：强制重新编译所有包
# - -installsuffix cgo：避免 CGO 相关文件
# - -ldflags：优化二进制文件（去除调试信息、减小体积）
# - -o main：输出二进制文件名为 main
# - ./main.go：入口文件路径（根据你的项目调整）
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags "-s -w" \
    -o main ./main.go


# 第二阶段：运行（轻量镜像，仅包含二进制文件）
FROM alpine:3.18

# 安全加固：创建非 root 用户运行程序（避免用 root 权限）
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# 安装必要工具（可选，如时区、证书）
# - tzdata：用于设置时区
# - ca-certificates：用于 HTTPS 请求（如 GORM 连接 HTTPS 数据库）
RUN apk --no-cache add tzdata ca-certificates

# 设置时区（可选，根据需要调整）
ENV TZ=Asia/Shanghai

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件到当前镜像
COPY --from=builder /app/main .
# 复制配置文件（如果有，根据实际路径调整）
# COPY --from=builder /app/config ./config

# 赋予程序执行权限
RUN chmod +x ./main

# 切换到非 root 用户
USER appuser

# 暴露 GIN 默认端口（GIN 默认监听 8080，根据你的项目调整）
EXPOSE 8080

# 启动程序
CMD ["./main"]
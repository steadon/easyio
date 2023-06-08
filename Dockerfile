# 设置基础镜像
FROM golang:1.20.4-alpine as builder

# 设置工作目录
WORKDIR /app

# 将代码复制到容器中
COPY . /app

# 设置代理环境变量
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

# 构建项目
RUN go build -o main .

# 第二个阶段，用于构建最终镜像
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 设置容器启动命令
CMD ["./main"]

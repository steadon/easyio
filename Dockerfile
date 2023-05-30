# 设置基础镜像为 golang 的 alpine
FROM golang:1.20.4-alpine as builder

# 设置工作目录
WORKDIR /app

# 将代码复制到容器中
COPY . /app

# 设置交叉编译环境变量
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

# 构建项目
RUN go build -o main .

# 设置最终镜像的基础镜像为 alpine
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 设置代理环境变量
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

# 从上一阶段中复制二进制文件
COPY --from=builder /app/main /app/main

# 设置容器启动命令
CMD ["/app/main"]
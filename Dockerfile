# 设置基础镜像
FROM golang:1.20.4-alpine

# 设置工作目录
WORKDIR /app

# 将代码复制到容器中
COPY . /app

# 设置代理环境变量
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

# 构建项目
RUN go build -o main .

# 设置容器启动命令
CMD ["./main"]

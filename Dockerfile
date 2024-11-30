# 使用官方 Golang 1.20 版本镜像作为基础镜像
FROM golang:1.20 AS build

# 设置 Go 代理（使用国内代理 goproxy.cn）
ENV GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制当前目录中的所有内容到工作目录
COPY . .

# 编译 Go 程序
RUN go build -o aliyun_ddns

# 定义启动容器时运行的命令
CMD ["./aliyun_ddns"]

# 暴露容器的 3000 端口
EXPOSE 3000

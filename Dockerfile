# 使用官方 Golang 1.18 版本镜像作为基础镜像
FROM golang:1.19.3 AS build


# 设置代理
ENV ALL_PROXY=socks5://192.168.1.33:10803
ENV HTTP_PROXY=http://192.168.1.33:10804
ENV HTTPS_PROXY=http://192.168.1.33:10804

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

# 暴露容器的 8800 端口
EXPOSE 8800

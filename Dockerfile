# 使用官方 Golang 1.20 版本镜像作为构建阶段的基础镜像
FROM golang:1.22.2 AS builder

# 设置 Go 代理（使用国内代理 goproxy.cn）
ENV GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载 Go 依赖
RUN go mod download

# 复制当前目录中的所有内容到工作目录
COPY . .

# 编译 Go 程序，生成静态链接的二进制文件
RUN CGO_ENABLED=0 go build -o aliyun_ddns .

# 使用 scratch 作为基础镜像，最小化镜像大小
FROM scratch

# 将编译后的二进制文件从构建阶段复制到最小镜像中
COPY --from=builder /app/aliyun_ddns /aliyun_ddns

# 配置容器的默认命令
CMD ["/aliyun_ddns"]

# 暴露容器的 3000 端口
EXPOSE 3000



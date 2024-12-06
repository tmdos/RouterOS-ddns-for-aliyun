# RouterOS 阿里云ddns动态域名解析，支持ipv4和ipv6 ！

### 常见RouterOS阿里云ddns脚本多数使用别人提供的api，自己部署更放心！

## 一、部署方式

支持RouterOS Container部署 和 服务器Docker容器部署两种方式。

#### 1. RouterOS Containe 部署
- 在线拉取或下载镜像上传再部署
- 参考文档（待更新）

#### 2. 服务器 Docker 容器部署
- Docker 镜像：[Docker Hub](https://hub.docker.com/r/tmdos/aliyun_ddns)
- 创建Docker容器并运行
```
docker run -d --name aliyun_ddns -p 3000:3000 tmdos/aliyun_ddns
```
## 二、RouterOS 6-7.x 脚本代码（IPv4/IPv6） 

请根据自己的实际情况替换 URL 中的参数：
- AccessKeyID：你的阿里云 AccessKey ID。
- AccessKeySecret：你的阿里云 AccessKey Secret。
- RR：子域名（如：home）。
- DomainName：你的主域名（如：baidu.com）。
- local pppoe "pppoe-out1" 接口名称，(如pppoe-out1/ether1)。
- [IPV4 脚本](https://github.com/tmdos/RouterOS-ddns-for-aliyun/blob/master/IPv4-Script)
- [IPV6 脚本](https://github.com/tmdos/RouterOS-ddns-for-aliyun/blob/master/IPv6-Script)

- **注意！！脚本中的(http://192.168.x.x) 需要根据你部署的ip来填写，
          例如RouterOS Container部署，那么EVTH的ip就是需要填写ip..**

## 三、API请求
- URL:
- http://:3000/aliyun_ddns?AccessKeyID=XXXXXX&AccessKeySecret=XXXXXX&RR=XX&DomainName=XXX&IpAddr=XXX

# RouterOS 阿里云ddns动态域名解析，支持ipv4t和ipv6 ！
#### 基于 lsprain 大佬项目的修改版，感谢 lsprain 的分享，原项目可参考 [lsprain 的 GitHub](https://github.com/lsprain)！

### 常见ros阿里云ddns脚本多数使用别人提供的api，自己部署更放心！

## 一、部署方式

支持服务器部署及 Docker 容器部署两种方式。

### 1. 服务器部署
- 自行编译并运行
- 使用已编译好的 Release 版本

### 2. Docker 容器部署
- Docker 镜像：[Docker Hub](https://hub.docker.com/r/tmdos/aliyun_ddns)

##### 拉取最新版镜像
```
docker pull tmdos/aliyun_ddns:latest
```
##### 创建Docker容器并运行
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

- **注意！！脚本中的阿里云DDNS更新请求的URL(192.168.x.x),这个要改你成部署API服务器的ip**

## 三、阿里云API请求方式
- Method: POST
- URL: http://192.168.x.x:3000/aliyun_ddns?AccessKeyID=XXXXXX&AccessKeySecret=XXXXXX&RR=XX&DomainName=XXX&IpAddr=XXX



##### 再次感谢 lsprain 大佬的原始项目，感谢其提供的宝贵分享！

# RouterOS 阿里云 DDNS 动态域名解析 (支持 IPv4 / IPv6)

### 常见 RouterOS 阿里云 DDNS 脚本多数使用别人提供的 API，自己部署更放心！

## 一、部署方式

支持 RouterOS Container 部署 和 服务器 Docker 容器部署两种方式。

#### 1. RouterOS Container 部署 (离线包方式)
- 下载 Releases 中的 `aliyun_ddns.tar` 并上传到 RouterOS 的 File List 根目录。
- 开启 RouterOS Container 功能，并正常创建网络 (veth、bridge 等，此处略过)。
- **请根据你的 RouterOS 系统版本，选择下方对应的部署代码：**
> **🟢 【如果你是 RouterOS v7.21 及以上版本】（🌟 强烈推荐）**
> v7.21 支持直接把缓存写进路由器的运行内存 (tmpfs)，零磁盘读写，极大保护你的 U 盘/硬盘寿命！
```routeros
# 1. 创建内存盘挂载点 (无需指定本地路径)
/container/mounts/add list=ddns_tmp dst="/tmp"
/container/mounts/add list=ddns_run dst="/run"
# 2. 创建容器 (root-dir 记得指向你的外接盘，如 disk1)
/container/add file=aliyun_ddns.tar interface=veth1 root-dir=disk1/ddns mountlists=ddns_tmp,ddns_run shm-size=128M start-on-boot=yes logging=yes
```
🔵 【如果你是 RouterOS v7.20 及以下版本】
> 老版本不支持内存盘挂载，为了防止频繁读写损坏硬盘，请直接使用不带挂载的“纯净版”命令：
# 直接创建容器 (不带挂载参数，root-dir 记得指向你的外接盘，如 disk1)
```
/container/add file=aliyun_ddns.tar interface=veth1 root-dir=disk1/ddns start-on-boot=yes logging=yes
```
#### 2. 服务器 Docker 容器部署
- Docker 镜像：[Docker Hub](https://hub.docker.com/r/tmdos/aliyun_ddns)
- ----------
- 拉取镜像
```
docker pull tmdos/aliyun_ddns
```
- 创建Docker容器并运行
```
docker run -d --name aliyun_ddns -p 3000:3000 tmdos/aliyun_ddns
```
## 二outerOS 6-7.x 脚本配置与部署 
⚠️极其重要：代码中的 http://192.168.x.x:3000 必须根据你的实际部署情况修改！
- Linux 服务器 Docker 部署：填写 Linux 服务器的内网 IP。
- RouterOS Container 部署：填写你为容器（VETH 虚拟网卡）分配的内网 IP（例如 172.17.0.2:3000）
- 
### 请根据自己的实际情况替换 URL 中的参数：
- AccessKeyID：你的阿里云 AccessKey ID。
- AccessKeySecret：你的阿里云 AccessKey Secret。
- RR：子域名（如想解析 home.baidu.com，此处填 home）。
- DomainName：你的主域名（如：baidu.com）。
- local pppoe "pppoe-out1" 接口名称，(IPv4 一般为 pppoe-out1，IPv6 一般为 bridge1 或 ether1)。
- 
### 1. [IPv4 脚本](./IPv4-Script) 部署方式 (推荐：PPPoE 触发)
1. 在 WinBox 进入 System -> Scripts 新建脚本，命名为 ipv4-ddns-script，贴入修改好参数的代码并保存。
2. 在你当前拨号的 PPP Profile 的 Scripts -> On Up 框中填入以下触发代码：
```
:delay 35;
/system script run ipv4-ddns-script;
:log info "PPPoE 拨号成功，已运行 DDNS 更新脚本";
```
- 💡 注：延迟 35 秒是为了防止路由器刚开机时 Docker 容器尚未启动完毕，导致请求发送失败。)
- 
### 2. [IPv6 脚本](./IPv6-Script) 部署方式 (推荐：定时任务触发)
1. 在 WinBox 进入 `System -> Scripts` 新建脚本，命名为 `ipv6-ddns-script`，贴入修改好参数的完整 IPv6 代码并保存。
2. 进入 `System -> Scheduler` 新建计划任务。
3. Name 随意（如 `Aliyun-DDNS-v6`），**Interval 建议设为 `00:01:00`（1分钟执行一次）**。
4. 在 `On Event` 框中填入以下调用代码并保存：
- 💡 注：请放心设置为 1 分钟。当 IP 没变时，脚本仅在本地内存极速比对，会在几毫秒内瞬间退出，不会向外发送网络请求，极度节省路由器性能。当然，你也可以根据喜好改为 3 或 5 分钟

## 三、API请求
- URL:
- http://ip:3000/aliyun_ddns?AccessKeyID=XXXXXX&AccessKeySecret=XXXXXX&RR=XX&DomainName=XXX&IpAddr=XXX

## 四、致谢
- #### 本脚本基于 lsprain 大佬项目的修改版，原项目可参考 lsprain 的[GitHub](https://github.com/lsprain/Aliddns-Ros)！
- #### 感谢其提供的宝贵分享！

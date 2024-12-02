package main

import (
    "Aliddns-Ros/log-handler" // 引入自定义的日志中间件包
    "github.com/denverdino/aliyungo/dns" // 引入阿里云的 DNS SDK
    "github.com/gin-gonic/gin" // 引入 Gin 框架
    log "github.com/sirupsen/logrus" // 引入日志库
    "net/http"
)

// ConfigInfo 定义域名相关配置信息
type ConfigInfo struct {
    AccessKeyID     string
    AccessKeySecret string
    DomainName      string
    RR              string
    IpAddr          string
    IpType          string // 新增字段，用于区分 IPv4 或 IPv6
}

func main() {
    r := gin.Default()               // 初始化 Gin 路由
    r.Use(middlewares.Logger())       // 使用自定义的日志中间件
    r.GET("/aliyun_ddns", AddUpdateAliddns)  // 设置路由，调用 AddUpdateAliddns 函数
    r.Run(":3000")                    // 监听 3000 端口
}

// AddUpdateAliddns 处理 DDNS 更新请求
func AddUpdateAliddns(c *gin.Context) {
    conf := new(ConfigInfo)

    // 获取请求中的参数
    conf.AccessKeyID = c.Query("AccessKeyID")
    conf.AccessKeySecret = c.Query("AccessKeySecret")
    conf.DomainName = c.Query("DomainName")
    conf.RR = c.Query("RR")
    conf.IpAddr = c.Query("IpAddr")
    conf.IpType = c.Query("IpType") // 读取 IpType 参数

    log.Println("当前路由公网IP：" + conf.IpAddr)
    log.Println("进行阿里云登录……")

    // 创建阿里云 DNS 客户端
    client := dns.NewClient(conf.AccessKeyID, conf.AccessKeySecret)
    client.SetDebug(false) // 禁用调试模式

    // 获取已有的域名记录
    domainInfo := new(dns.DescribeDomainRecordsArgs)
    domainInfo.DomainName = conf.DomainName
    oldRecord, err := client.DescribeDomainRecords(domainInfo)
    if err != nil {
        log.Println("阿里云登录失败！请查看错误日志！", err)
        c.String(http.StatusOK, "loginerr") // 登录失败返回 "loginerr"
        return
    }
    log.Println("阿里云登录成功！")
    log.Println("进行域名及IP比对……")

    var exsitRecordID string
    // 查找匹配的记录
    for _, record := range oldRecord.DomainRecords.Record {
        if record.DomainName == conf.DomainName && record.RR == conf.RR {
            // 如果记录中的值和传入的 IP 地址相同，返回 "same" 表示无需更新
            if record.Value == conf.IpAddr {
                log.Println("当前配置解析地址与公网IP相同，不需要修改。")
                c.String(http.StatusOK, "same")
                return
            }
            exsitRecordID = record.RecordId
        }
    }

    // 根据 IpType 设置记录类型
    recordType := dns.ARecord // 默认使用 A 记录（IPv4）
    if conf.IpType == "IPv6" { // 如果是 IPv6，则使用 AAAA 记录
        recordType = dns.AAAARecord
    }

    // 更新现有记录
    if len(exsitRecordID) > 0 {
        updateRecord := new(dns.UpdateDomainRecordArgs)
        updateRecord.RecordId = exsitRecordID
        updateRecord.RR = conf.RR
        updateRecord.Value = conf.IpAddr
        updateRecord.Type = recordType
        rsp := new(dns.UpdateDomainRecordResponse)
        rsp, err := client.UpdateDomainRecord(updateRecord)
        if err != nil {
            log.Println("修改解析地址信息失败!", err)
            c.String(http.StatusOK, "iperr") // 更新失败返回 "iperr"
        } else {
           log.Println("修改解析地址信息成功!", rsp)
            c.String(http.StatusOK, "ip") // 更新成功返回 "ip"
        }
    } else {
        // 如果没有找到记录，则添加新记录
        newRecord := new(dns.AddDomainRecordArgs)
        newRecord.DomainName = conf.DomainName
        newRecord.RR = conf.RR
        newRecord.Value = conf.IpAddr
        newRecord.Type = recordType
        rsp := new(dns.AddDomainRecordResponse)
        rsp, err := client.AddDomainRecord(newRecord)
        if err != nil {
            log.Println("添加新域名解析失败！", err)
            c.String(http.StatusOK, "domainerr") // 添加失败返回 "domainerr"
        } else {
            log.Println("添加新域名解析成功！", rsp)
            c.String(http.StatusOK, "domain") // 添加成功返回 "domain"
        }
    }
}

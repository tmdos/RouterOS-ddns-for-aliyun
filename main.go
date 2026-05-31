package main

import (
    "aliyun_ddns/middlewares"
    "log"
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

// ConfigInfo 保存请求中的参数
type ConfigInfo struct {
    AccessKeyID     string
    AccessKeySecret string
    DomainName      string
    RR              string
    IpAddr          string
    IpType          string
}

// 设置记录参数的通用函数
func setRecordParams(record *alidns.AddDomainRecordRequest, conf *ConfigInfo, recordType string) {
    record.DomainName = conf.DomainName
    record.RR = conf.RR
    record.Value = conf.IpAddr
    record.Type = recordType
}

// AddUpdateAliddns 处理阿里云DNS更新请求
func AddUpdateAliddns(c *gin.Context) {
    conf := new(ConfigInfo)

    // 获取请求中的参数
    conf.AccessKeyID = c.Query("AccessKeyID")
    conf.AccessKeySecret = c.Query("AccessKeySecret")
    conf.DomainName = c.Query("DomainName")
    conf.RR = c.Query("RR")
    conf.IpAddr = c.Query("IpAddr")
    conf.IpType = c.Query("IpType")

    // 记录请求的日志
    log.Printf("Requesting Aliyun API with params: AccessKeyID=%s, DomainName=%s, RR=%s, IpAddr=%s, IpType=%s", 
        conf.AccessKeyID, conf.DomainName, conf.RR, conf.IpAddr, conf.IpType)

    log.Println("进行阿里云登录……")

    // 创建阿里云 DNS 客户端
    client, err := alidns.NewClientWithAccessKey("cn-hangzhou", conf.AccessKeyID, conf.AccessKeySecret)
    if err != nil {
        log.Printf("Failed to create Aliyun client: %v", err)
        c.String(http.StatusInternalServerError, "loginerr") // 登录失败返回 "loginerr"
        return
    }

    log.Println("阿里云登录成功！")
    log.Println("获取已有的域名记录……")

    // 获取已有的域名记录
    domainInfo := alidns.CreateDescribeDomainRecordsRequest()
    domainInfo.DomainName = conf.DomainName
    oldRecord, err := client.DescribeDomainRecords(domainInfo)
    if err != nil {
        log.Printf("Aliyun API Error: %v", err)
        c.String(http.StatusInternalServerError, "loginerr") // 登录失败返回 "loginerr"
        return
    }

    log.Println("获取域名记录成功！")
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
    recordType := "A" // 默认使用 A 记录（IPv4）
    if conf.IpType == "IPv6" || strings.Contains(conf.IpAddr, ":") { // 如果是 IPv6，则使用 AAAA 记录
        recordType = "AAAA"
    }

    if len(exsitRecordID) > 0 {
        // 更新现有记录
        updateRecord := alidns.CreateUpdateDomainRecordRequest()
        updateRecord.RecordId = exsitRecordID
        updateRecord.RR = conf.RR
        updateRecord.Value = conf.IpAddr
        updateRecord.Type = recordType
        log.Printf("Aliyun Update Request: %+v", updateRecord)
        rsp, err := client.UpdateDomainRecord(updateRecord)
        if err != nil {
            log.Printf("Aliyun API Error: %v", err)
            c.String(http.StatusInternalServerError, "iperr") // 更新失败返回 "iperr"
        } else {
            log.Printf("Aliyun API Response: %+v", rsp)
            c.String(http.StatusOK, "ip") // 更新成功返回 "ip"
        }
    } else {
        // 如果没有找到记录，则添加新记录
        newRecord := alidns.CreateAddDomainRecordRequest()
        setRecordParams(newRecord, conf, recordType)
        log.Printf("Aliyun Add Request: %+v", newRecord)
        rsp, err := client.AddDomainRecord(newRecord)
        if err != nil {
            log.Printf("Aliyun API Error: %v", err)
            c.String(http.StatusInternalServerError, "domainerr") // 添加失败返回 "domainerr"
        } else {
            log.Printf("Aliyun API Response: %+v", rsp)
            c.String(http.StatusOK, "domain") // 添加成功返回 "domain"
        }
    }
}

func main() {
    gin.SetMode(gin.ReleaseMode) // 设置 Gin 为发布模式
    r := gin.New()
    r.Use(gin.Recovery())
    r.Use(middlewares.Logger())
    r.GET("/aliyun_ddns", AddUpdateAliddns)
    r.Run(":3000") // 确保服务监听 3000 端口
}

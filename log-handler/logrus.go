package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func Logger() gin.HandlerFunc {
	log := logrus.New()

	// 设置输出文件路径
	filePath := "logs/"
	fileName := "Aliddns-API"
	file := filePath + fileName

	// 确保日志目录存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.MkdirAll(filePath, os.ModePerm)
	}

	// 设置日志级别
	log.SetLevel(logrus.DebugLevel)

	// 设置日志切割 rotatelogs
	writer, err := rotatelogs.New(
		file+"%Y%m%d.log",
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		log.Fatal(err)
	}

	// 使用 lfshook 将日志写入指定的文件
	writeMap := lfshook.WriterMap{
		logrus.PanicLevel: writer,
		logrus.FatalLevel: writer,
		logrus.ErrorLevel: writer,
		logrus.WarnLevel:  writer,
		logrus.InfoLevel:  writer,
		logrus.DebugLevel: writer,
	}

	// 配置 lfshook
	hook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05", // 设置日期格式
	})

	// 为 logrus 实例添加自定义 hook
	log.AddHook(hook)

	return func(c *gin.Context) {
		// 配置请求信息
		startTime := time.Now()
		c.Next()
		spendTime := time.Since(startTime).Milliseconds()
		ST := fmt.Sprintf("%d ms", spendTime) // API 调用耗时

		statusCode := c.Writer.Status()      // 状态码
		clientIP := c.ClientIP()              // 请求客户端的 IP
		dataSize := c.Writer.Size()           // 响应报文字节长度
		if dataSize < 0 {
			dataSize = 0
		}
		method := c.Request.Method           // 请求方法
		path := c.Request.RequestURI         // 请求 URL

		// 创建日志条目并添加字段
		entry := log.WithFields(logrus.Fields{
			"Status":    statusCode,
			"SpendTime": ST,
			"IP":        clientIP,
			"Method":    method,
			"Path":      path,
		})

		// 记录错误信息
		if len(c.Errors) > 0 {
			log.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		}

		// 根据状态码决定日志级别
		if statusCode >= 500 {
			entry.Error()
		} else if statusCode >= 400 {
			entry.Warn()
		} else {
			entry.Info()
		}
	}
}

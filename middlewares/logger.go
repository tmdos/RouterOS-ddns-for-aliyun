package middlewares

import (
    "fmt"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    rotatelogs "github.com/lestrrat-go/file-rotatelogs"
    "github.com/rifflock/lfshook"
    "github.com/sirupsen/logrus"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
    log := logrus.New()

    // 设置输出文件路径，改为可配置
    filePath := getEnv("LOG_PATH", "logs/")
    fileName := getEnv("LOG_FILE", "Aliddns-API")
    file := filePath + fileName

    // 确保日志目录存在
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
            log.Fatalf("Failed to create log directory: %v", err)
        }
    }

    // 设置日志级别为 Info，避免 Debug 输出，改为可配置
    level, err := logrus.ParseLevel(getEnv("LOG_LEVEL", "info"))
    if err != nil {
        log.Fatalf("Invalid log level: %v", err)
    }
    log.SetLevel(level)

    // 设置日志切割 rotatelogs
    writer, err := rotatelogs.New(
        file+"%Y%m%d.log",
        rotatelogs.WithMaxAge(7*24*time.Hour),
        rotatelogs.WithRotationTime(24*time.Hour),
    )
    if err != nil {
        log.Fatalf("Failed to initialize log rotation: %v", err)
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
        startTime := time.Now()
        c.Next()
        spendTime := time.Since(startTime).Milliseconds()
        ST := fmt.Sprintf("%d ms", spendTime)

        statusCode := c.Writer.Status()
        clientIP := c.ClientIP()
        dataSize := c.Writer.Size()
        if dataSize < 0 {
            dataSize = 0
        }
        method := c.Request.Method
        path := c.Request.RequestURI

        entry := log.WithFields(logrus.Fields{
            "Status":    statusCode,
            "SpendTime": ST,
            "IP":        clientIP,
            "Method":    method,
            "Path":      path,
        })

        if len(c.Errors) > 0 {
            entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
        }

        if statusCode >= 500 {
            entry.Error()
        } else if statusCode >= 400 {
            entry.Warn()
        } else {
            entry.Info()
        }
    }
}

// getEnv 获取环境变量值，如果未设置则使用默认值
func getEnv(key, defaultValue string) string {
    value, exists := os.LookupEnv(key)
    if !exists {
        return defaultValue
    }
    return value
}


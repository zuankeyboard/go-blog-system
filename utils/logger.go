package utils

import (
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Log 全局日志实例
var Log *logrus.Logger

// InitLogger 初始化日志
func InitLogger() {
	// 创建日志目录
	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		_ = os.Mkdir(logDir, 0755)
	}

	// 日志文件路径
	logFile := filepath.Join(logDir, "blog_"+time.Now().Format("20060102")+".log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic("日志文件创建失败: " + err.Error())
	}

	// 初始化logrus
	Log = logrus.New()
	// 同时输出到控制台和文件
	Log.SetOutput(os.Stdout)
	Log.SetOutput(file) // 仅输出到文件

	// 日志格式
	Log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceColors:     true,
	})

	// 日志级别
	if gin.Mode() == gin.DebugMode {
		Log.SetLevel(logrus.DebugLevel)
	} else {
		Log.SetLevel(logrus.InfoLevel)
	}
}

// GinLogger Gin请求日志中间件
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		fields := logrus.Fields{
			"method":  c.Request.Method,
			"path":    c.FullPath(),
			"status":  c.Writer.Status(),
			"latency": latency,
			"ip":      c.ClientIP(),
		}

		switch {
		case c.Writer.Status() >= 500:
			Log.WithFields(fields).Error("服务端错误")
		case c.Writer.Status() >= 400:
			Log.WithFields(fields).Warn("客户端错误")
		default:
			Log.WithFields(fields).Info("请求成功")
		}
	}
}

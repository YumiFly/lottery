// utils/logger.go
package utils

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func InitLogger() {
	// 设置 JSON 格式
	Logger.SetFormatter(&logrus.JSONFormatter{})

	// 创建日志目录
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		Logger.Fatalf("Failed to create log directory: %v", err)
	}

	// 打开日志文件
	logFile := filepath.Join(logDir, "app.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Logger.Fatalf("Failed to open log file: %v", err)
	}

	// 设置输出到文件和标准输出
	Logger.SetOutput(file) // 只写入文件
	// 如果需要同时输出到控制台，可以使用：
	// Logger.SetOutput(io.MultiWriter(file, os.Stdout))

	// 设置日志级别
	Logger.SetLevel(logrus.InfoLevel)
}

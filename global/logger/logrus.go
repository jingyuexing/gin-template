package logger

import (
	"io"
	"os"
	"path/filepath"

	"github.com/jingyuexing/go-utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// var Loggerus *logrus.Logger

func LogrusInit() *logrus.Logger {
	// 创建一个新的 Logger 实例
	Loggerus := logrus.New()

	Loggerus.SetOutput(os.Stdout)

	Loggerus.SetLevel(logrus.InfoLevel)

	Loggerus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,                  // 显示完整的时间戳
		TimestampFormat: "2006-01-02 15:04:05", // 自定义时间格式
	})

	// 确保日志目录存在
	logDir := "log"
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		Loggerus.Fatalf("无法创建日志目录: %v", err)
	}

	datetime := utils.NewDateTime()
	datetime.Format("YYYY-MM-DD")
	loggerInfoName := utils.Template("GOA-{date}-info.log", map[string]any{
		"date": datetime.Format("YYYY-MM-DD"),
	}, "{}")
	loggerErrorName := utils.Template("GOA-{date}-error.log", map[string]any{
		"date": datetime.Format("YYYY-MM-DD"),
	}, "{}")

	infoLogPath := filepath.Join(logDir,loggerInfoName)
	errorLogPath := filepath.Join(logDir, loggerErrorName)

	debugFile := &lumberjack.Logger{
        Filename:  infoLogPath, // 将 debug 和 info 日志写入同一个文件
        MaxSize:    10,                // 文件最大 10MB
        MaxBackups: 5,                 // 最多保留 5 个备份
        MaxAge:     30,                // 保留 30 天内的日志
        Compress:   true,              // 启用压缩
    }
    errorFile := &lumberjack.Logger{
        Filename:  errorLogPath, // 将 debug 和 info 日志写入同一个文件
        MaxSize:    10,                // 文件最大 10MB
        MaxBackups: 5,                 // 最多保留 5 个备份
        MaxAge:     30,                // 保留 30 天内的日志
        Compress:   true,              // 启用压缩
    }
    Loggerus.SetOutput(io.MultiWriter(debugFile,os.Stdout))
    Loggerus.SetLevel(logrus.DebugLevel)

    Loggerus.SetOutput(errorFile)
	return Loggerus
}

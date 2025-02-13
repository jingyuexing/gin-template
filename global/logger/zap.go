package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"template/global/config"

	"github.com/jingyuexing/go-utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(config config.System) *zap.Logger {
	// 确保日志目录存在
	if err := os.MkdirAll(config.LoggerPath, os.ModePerm); err != nil {
		fmt.Printf("无法创建日志目录: %v", err)
		return nil
	}

	appName := config.Name
	if appName == "" {
		appName = "app"
	}

	datetime := utils.NewDateTime()
	loggerInfoName := utils.Template(config.LoggerName, map[string]any{
		"date": datetime.Format("YYYY-MM-DD"),
		"app":  appName,
		"level": "info",
	}, "{}")
	loggerDebugName := utils.Template(config.LoggerName, map[string]any{
		"date": datetime.Format("YYYY-MM-DD"),
		"app":  appName,
		"level": "debug",
	}, "{}")
	loggerErrorName := utils.Template(config.LoggerName, map[string]any{
		"date": datetime.Format("YYYY-MM-DD"),
		"app":  appName,
		"level": "error",
	}, "{}")

	// 使用filepath.Join构建完整的日志文件路径
	debugLogPath := filepath.Join(config.LoggerPath, loggerDebugName)
	infoLogPath := filepath.Join(config.LoggerPath, loggerInfoName)
	errorLogPath := filepath.Join(config.LoggerPath, loggerErrorName)

	// 创建Encoder配置
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:   "message",                     // 日志内容的字段名
		LevelKey:     "level",                       // 日志级别的字段名
		TimeKey:      "time",                        // 时间的字段名
		CallerKey:    "caller",                      // 调用者信息的字段名
		EncodeCaller: zapcore.ShortCallerEncoder,    // 使用短的调用栈格式
		EncodeLevel:  zapcore.LowercaseLevelEncoder, // 小写日志级别
		EncodeTime:   zapcore.ISO8601TimeEncoder,    // 时间格式
	}

	// 创建不同级别的日志文件,使用完整路径
	debugFile := &lumberjack.Logger{
		Filename:   debugLogPath,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}

	infoFile := &lumberjack.Logger{
		Filename:   infoLogPath,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}

	errorFile := &lumberjack.Logger{
		Filename:   errorLogPath,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}

	// 创建不同日志级别的Core
	debugCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(debugFile),
		zapcore.DebugLevel, // Debug级别
	)

	infoCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(infoFile),
		zapcore.InfoLevel, // Info级别
	)

	errorCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(errorFile),
		zapcore.ErrorLevel, // Error级别
	)

	// 使用Tee将多个Core合并
	core := zapcore.NewTee(debugCore, infoCore, errorCore)

	// 创建Logger
	logger := zap.New(core)

	return logger
}

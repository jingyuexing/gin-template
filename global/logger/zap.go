package logger

import (
	"fmt"

	"github.com/jingyuexing/go-utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger() *zap.Logger {
	datetime := utils.NewDateTime()
	loggerInfoName := utils.Template("app-{date}-info.log", map[string]any{
		"date": datetime.Format("YYYY-MM-DD"),
	}, "{}")
	loggerDebugName := utils.Template("app-{date}-debug.log", map[string]any{
		"date": datetime.Format("YYYY-MM-DD"),
	}, "{}")
	loggerErrorName := utils.Template("app-{date}-error.log", map[string]any{
		"date": datetime.Format("YYYY-MM-DD"),
	}, "{}")

	// 构建文件名
	// filename := "app_" + utils.NewDateTime().Format("2006-01-02") + ".log"

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

	// 创建不同级别的日志文件
	debugFile, err := createLogFile(loggerDebugName)
	if err != nil {
		fmt.Printf("Error creating debug log file: %s", err.Error())
	}
	infoFile, err := createLogFile(loggerInfoName)
	if err != nil {
		fmt.Printf("Error creating info log file: %s", err.Error())
	}
	errorFile, err := createLogFile(loggerErrorName)
	if err != nil {
		fmt.Printf("Error creating error log file: %v", err.Error())
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

func createLogFile(filename string) (*lumberjack.Logger, error) {
	return &lumberjack.Logger{
		Filename:   filename, // 设置日志文件名
		MaxSize:    10,       // 最大日志文件大小（单位MB）
		MaxBackups: 5,        // 最大备份数量
		MaxAge:     30,       // 保留日志的最大天数
		Compress:   true,     // 启用日志压缩
	}, nil
}

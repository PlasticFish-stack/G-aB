package logger

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *zap.Logger
)

// InitLogger 初始化日志记录器
func InitLogger(isProduction bool) {
	// var err error
	logDir := "logs"
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		panic(err)
	}
	currentDate := time.Now().Format("20060102")
	logFilePath := filepath.Join(logDir, "app-"+currentDate+".log") // 根据日期生成文件名
	// 设置日志级别
	var level zapcore.Level
	if isProduction {
		level = zapcore.InfoLevel
	} else {
		level = zapcore.DebugLevel
	}

	// 配置日志轮转
	logFile := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    10,   // 每个文件最大10MB
		MaxBackups: 7,    // 保留7个备份
		MaxAge:     7,    // 备份保留7天
		Compress:   true, // 是否压缩
	}

	// 设置日志编码器
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	if !isProduction {
		encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	}

	// 创建核心日志
	core := zapcore.NewCore(encoder, zapcore.AddSync(logFile), level)

	// 创建 logger
	logger = zap.New(core)

	// 允许使用 zap的全局 logger
	zap.RedirectStdLog(logger)
}

// Info 记录 Info 级别的日志
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Warn 记录 Warn 级别的日志
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Error 记录 Error 级别的日志
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

// Sync 确保日志被写入
func Sync() {
	if logger != nil {
		_ = logger.Sync()
	}
}

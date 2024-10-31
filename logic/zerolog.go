package logic

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger zerolog.Logger

func InitLogger() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	currentDate := time.Now().Format("2006-01-02")
	logFilename := fmt.Sprintf("./log/%s.log", currentDate)
	// 创建一个日志文件输出，使用 lumberjack 实现轮换
	logFile := &lumberjack.Logger{
		Filename:   logFilename, // 日志文件路径
		MaxSize:    100,         // MB
		MaxBackups: 3,           // 保留的旧日志文件的最大数量
		MaxAge:     14,          // 保留旧日志的天数
		Compress:   true,        // 是否压缩旧日志
	}
	Logger = zerolog.New(logFile).With().Timestamp().Logger()
}
func LogInfo(msg string) {
	Logger.Info().Msg(msg)
}

// LogError 记录错误级别日志
func LogError(msg string) {
	Logger.Error().Msg(msg)
}

package tools

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"time"
)

type LoggerConfig struct {
	Mode       string // dev/test/prod
	LogDir     string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

// getEncoder 返回编码器
func GetEncoder(mode string) zapcore.Encoder {
	if mode == "dev" {
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "🕒",
			LevelKey:       "🔔",
			NameKey:        "模块",
			CallerKey:      "📍",
			MessageKey:     "💬",
			StacktraceKey:  "🧵",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    ColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
		return zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 生产环境用 JSON
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "module",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	return zapcore.NewJSONEncoder(encoderConfig)
}

// 彩色日志等级输出
func ColorLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch level {
	case zapcore.DebugLevel:
		enc.AppendString("\033[36m🐛 DEBUG\033[0m") // 蓝色
	case zapcore.InfoLevel:
		enc.AppendString("\033[32m✅ INFO\033[0m") // 绿色
	case zapcore.WarnLevel:
		enc.AppendString("\033[33m⚠️ WARN\033[0m") // 黄色
	case zapcore.ErrorLevel:
		enc.AppendString("\033[31m❌ ERROR\033[0m") // 红色
	case zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		enc.AppendString("\033[35m🔥 PANIC\033[0m") // 紫色
	default:
		enc.AppendString(level.String())
	}
}

// getLogWriter 返回日志写入器
func GetLogWriter(cfg *LoggerConfig) zapcore.WriteSyncer {
	date := time.Now().Format("2006-01-02")
	logPath := filepath.Join(cfg.LogDir, date)
	_ = os.MkdirAll(logPath, os.ModePerm)
	fmt.Println("日志输出路径:", logPath)
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filepath.Join(logPath, "server.log"),
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}
	return zapcore.AddSync(lumberJackLogger)
}

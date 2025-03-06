package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

/* zap 日志对象 */

var (
	Logger *zap.SugaredLogger
)

func InitLogger(debug bool) {
	level := zapcore.DebugLevel
	var syncWriter zapcore.WriteSyncer
	// 防止生产输出debug日志
	if !debug {
		level = zapcore.InfoLevel
	}
	// 开发日志输出到控制台
	if debug {
		syncWriter = os.Stdout
	} else {
		// 定时整理日志
		syncWriter = zapcore.AddSync(&lumberjack.Logger{
			Filename:  "./logs/error.log",
			MaxSize:   100,
			LocalTime: true,
			Compress:  true,
		})
	}

	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeTime = zapcore.EpochMillisTimeEncoder // 时间格式
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoder), syncWriter, zap.NewAtomicLevelAt(zapcore.Level(level)))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0))
	Logger = logger.Sugar()
}

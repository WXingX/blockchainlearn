package logger

import (
	"blog-management/config"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var Logger *zap.Logger
var atomicLevel zap.AtomicLevel // 动态调整日志级别

func InitLogger(cfg config.LogConfig) {
	//var writer zapcore.WriteSyncer
	logFileName := fmt.Sprintf("%s/%s_%s.log", cfg.FilePath, cfg.FileName, time.Now().Format("2006-01-02"))
	//创建日志目录
	_ = os.MkdirAll(cfg.FilePath, os.ModePerm)
	// zap 编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.CallerKey = ""

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// 日志级别
	//var level zapcore.Level
	//switch cfg.Level {
	//case "debug":
	//	level = zap.DebugLevel
	//case "warn":
	//	level = zap.WarnLevel
	//case "error":
	//	level = zap.ErrorLevel
	//default:
	//	level = zap.InfoLevel
	//}

	// 动态日志级别
	atomicLevel = zap.NewAtomicLevel()
	// 下面这里自动转换cfg.Level 为 zapcore.Level
	if err := atomicLevel.UnmarshalText([]byte(cfg.Level)); err != nil {
		atomicLevel.SetLevel(zap.InfoLevel)
	}

	// 文件输出（lumberjack）
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	})

	// 控制台输出
	consoleWriter := zapcore.AddSync(os.Stdout)

	var cores []zapcore.Core
	if cfg.Mode == "console" || cfg.Mode == "both" {
		cores = append(cores, zapcore.NewCore(consoleEncoder, consoleWriter, atomicLevel))
	}
	if cfg.Mode == "file" || cfg.Mode == "both" {
		cores = append(cores, zapcore.NewCore(fileEncoder, fileWriter, atomicLevel))
	}

	Logger = zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(1))
}

func SetLogLevel(level string) error {
	return atomicLevel.UnmarshalText([]byte(level))
}

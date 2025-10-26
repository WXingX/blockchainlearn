package xzap

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

const (
	loggerKey = "logger"
	customMsg = "custom"
)

type CtxLogger struct {
	logger *zap.Logger
	ctx    context.Context
}

func ToContext(ctx context.Context, logger *zap.Logger) context.Context {
	l := &CtxLogger{
		logger: logger,
	}
	return context.WithValue(ctx, loggerKey, l)
}

func WithContext(ctx context.Context) *CtxLogger {
	l, ok := ctx.Value(loggerKey).(*CtxLogger)
	if !ok || l == nil {
		return NewContextLogger(ctx)
	}

	l.ctx = ctx
	return l
}

func NewContextLogger(ctx context.Context) *CtxLogger {
	return &CtxLogger{
		logger: GetZapLogger(),
		ctx:    ctx,
	}
}

// Debug 调用 zap.Logger Debug
func (l *CtxLogger) Debug(msg string, fields ...zap.Field) {
	l.logger.WithOptions(zap.AddCallerSkip(1)).Debug(msg, fields...)
}

// Info 调用 zap.Logger Info
func (l *CtxLogger) Info(msg string, fields ...zap.Field) {
	l.logger.WithOptions(zap.AddCallerSkip(1)).Info(msg, fields...)
}

// Warn 调用 zap.Logger Warn
func (l *CtxLogger) Warn(msg string, fields ...zap.Field) {
	l.logger.WithOptions(zap.AddCallerSkip(1)).Warn(msg, fields...)
}

// Error 调用 zap.Logger Error
func (l *CtxLogger) Error(msg string, fields ...zap.Field) {
	l.logger.WithOptions(zap.AddCallerSkip(1)).Error(msg, fields...)
}

// Panic 调用 zap.Logger Panic
func (l *CtxLogger) Panic(msg string, fields ...zap.Field) {
	l.logger.WithOptions(zap.AddCallerSkip(1)).Panic(msg, fields...)
}

// Debugf 调用 zap.Logger Debug
func (l *CtxLogger) Debugf(format string, data ...interface{}) {
	l.Debug(customMsg, zap.String("content", fmt.Sprintf(format, data...)))
}

// Infof 调用 zap.Logger Info
func (l *CtxLogger) Infof(format string, data ...interface{}) {
	l.Info(customMsg, zap.String("content", fmt.Sprintf(format, data...)))
}

// Warnf 调用 zap.Logger Warn
func (l *CtxLogger) Warnf(format string, data ...interface{}) {
	l.Warn(customMsg, zap.String("content", fmt.Sprintf(format, data...)))
}

// Errorf 调用 zap.Logger Error
func (l *CtxLogger) Errorf(format string, data ...interface{}) {
	l.Error(customMsg, zap.String("content", fmt.Sprintf(format, data...)))
}

// Panicf 调用 zap.Logger Panic
func (l *CtxLogger) Panicf(format string, data ...interface{}) {
	l.Error(customMsg, zap.String("content", fmt.Sprintf(format, data...)))
}

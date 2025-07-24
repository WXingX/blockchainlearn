package database

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"time"
)

type ZapGormLogger struct {
	log *zap.Logger
}

func (l *ZapGormLogger) LogMode(level logger.LogLevel) logger.Interface { return l }
func (l *ZapGormLogger) Info(ctx context.Context, s string, i ...interface{}) {
	l.log.Sugar().Infof(s, i...)
}
func (l *ZapGormLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	l.log.Sugar().Warnf(s, i...)
}
func (l *ZapGormLogger) Error(ctx context.Context, s string, i ...interface{}) {
	l.log.Sugar().Errorf(s, i...)
}
func (l *ZapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	l.log.Info("SQL Trace", zap.String("sql", sql), zap.Duration("duration", time.Since(begin)),
		zap.Int64("rows", rows), zap.Error(err))
}

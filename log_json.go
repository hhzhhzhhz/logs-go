package logs_go

import (
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"io"
)

type LogJson interface {
	Sync() error
	Close() error
	Info(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	DPanic(msg string, fields ...zap.Field)
	Panic(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

func NewLogJson(l *zap.Logger, cs []io.Closer) LogJson {
	return &logJson{
		logger: l,
		closes: cs,
	}
}

type logJson struct {
	logger *zap.Logger
	closes []io.Closer
}

func (l *logJson) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *logJson) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *logJson) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *logJson) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *logJson) DPanic(msg string, fields ...zap.Field) {
	l.logger.DPanic(msg, fields...)
}

func (l *logJson) Panic(msg string, fields ...zap.Field) {
	l.logger.Panic(msg, fields...)
}

func (l *logJson) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

func (l *logJson) Sync() error {
	return l.logger.Core().Sync()
}

func (l *logJson) Close() error {
	var errs error
	for _, c := range l.closes {
		fn, ok := c.(interface{ Close() error })
		if ok {
			if err := fn.Close(); err != nil {
				errs = multierr.Append(errs, err)
			}
		}
	}
	return errs
}

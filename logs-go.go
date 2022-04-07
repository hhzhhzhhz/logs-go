package logs_go

import (
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"io"
	"logs-go/logger"
)

func NewLogger(l *zap.Logger, cs []io.Closer) logger.Logger {
	return &log{
		logger: l,
		closes: cs,
	}
}

type log struct {
	logger *zap.Logger
	closes []io.Closer
}

func (l *log) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *log) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *log) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *log) DPanic(msg string, fields ...zap.Field) {
	l.logger.DPanic(msg, fields...)
}

func (l *log) Panic(msg string, fields ...zap.Field) {
	l.logger.Panic(msg, fields...)
}

func (l *log) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

func (l *log) Close() error {
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

func (l *log) Sync() error {
	return l.logger.Core().Sync()
}

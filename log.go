package logs_go

import (
	"fmt"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"log"
	"os"
)

type Logf interface {
	io.Closer
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
	DPanic(format string, v ...interface{})
	Panic(format string, v ...interface{})
	Fatal(format string, v ...interface{})
}

type LogJ interface {
	io.Closer
	Sync() error
	Info(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	DPanic(msg string, fields ...zap.Field)
	Panic(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

func NewLogJ(l *zap.Logger, cs []io.Closer) LogJ {
	return &logj{
		logger: l,
		closes: cs,
	}
}

type logj struct {
	logger *zap.Logger
	closes []io.Closer
}

func (l *logj) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *logj) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *logj) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *logj) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *logj) DPanic(msg string, fields ...zap.Field) {
	l.logger.DPanic(msg, fields...)
}

func (l *logj) Panic(msg string, fields ...zap.Field) {
	l.logger.Panic(msg, fields...)
}

func (l *logj) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

func (l *logj) Sync() error {
	return l.logger.Core().Sync()
}

func (l *logj) Close() error {
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

func NewLogf(w io.Writer, cs []io.Closer, level zapcore.Level, errout io.Writer, calldepth int) Logf {
	l := log.New(w, "", log.Ldate|log.Ltime|log.Lshortfile)
	if calldepth < 2 {
		calldepth = 2
	}
	return &logf{
		l:         l,
		closes:    cs,
		level:     level,
		calldepth: calldepth,
	}
}

type logf struct {
	l         *log.Logger
	closes    []io.Closer
	level     zapcore.Level
	errorOut  io.Writer
	calldepth int
}

func (a *logf) Debug(format string, v ...interface{}) {
	if a.level <= zapcore.DebugLevel {
		if err := a.l.Output(a.calldepth, fmt.Sprintf("[DEBUG] "+format+"\n", v...)); err != nil {
			if a.errorOut != nil {
				a.errorOut.Write([]byte(err.Error()))
			}
		}
	}
}

func (a *logf) Info(format string, v ...interface{}) {
	if a.level <= zapcore.InfoLevel {
		if err := a.l.Output(a.calldepth, fmt.Sprintf("[INFO] "+format+"\n", v...)); err != nil {
			if a.errorOut != nil {
				a.errorOut.Write([]byte(err.Error()))
			}
		}
	}
}

func (a *logf) Warn(format string, v ...interface{}) {
	if a.level <= zapcore.WarnLevel {
		if err := a.l.Output(a.calldepth, fmt.Sprintf("[WARN] "+format+"\n", v)); err != nil {
			if a.errorOut != nil {
				a.errorOut.Write([]byte(err.Error()))
			}
		}
	}
}

func (a *logf) Error(format string, v ...interface{}) {
	if a.level <= zapcore.ErrorLevel {
		if err := a.l.Output(a.calldepth, fmt.Sprintf("[ERROR] "+format+"\n", v...)); err != nil {
			if a.errorOut != nil {
				a.errorOut.Write([]byte(err.Error()))
			}
		}
	}
}

func (a *logf) DPanic(format string, v ...interface{}) {
	if a.level <= zapcore.DPanicLevel {
		if err := a.l.Output(a.calldepth, fmt.Sprintf("[DPANIC] "+format+"\n", v...)); err != nil {
			if a.errorOut != nil {
				a.errorOut.Write([]byte(err.Error()))
			}
		}
	}
}

func (a *logf) Panic(format string, v ...interface{}) {
	if a.level <= zapcore.PanicLevel {
		if err := a.l.Output(a.calldepth, fmt.Sprintf("[FATAL] "+format+"\n", v...)); err != nil {
			if a.errorOut != nil {
				a.errorOut.Write([]byte(err.Error()))
			}
		}
	}
}

func (a *logf) Fatal(format string, v ...interface{}) {
	if a.level <= zapcore.FatalLevel {
		if err := a.l.Output(a.calldepth, fmt.Sprintf("[FATAL] "+format+"\n", v...)); err != nil {
			if a.errorOut != nil {
				a.errorOut.Write([]byte(err.Error()))
			}
		}
		os.Exit(1)
	}
}

func (a *logf) Close() error {
	var errs error
	for _, c := range a.closes {
		fn, ok := c.(interface{ Close() error })
		if ok {
			if err := fn.Close(); err != nil {
				errs = multierr.Append(errs, err)
			}
		}
	}
	return errs
}

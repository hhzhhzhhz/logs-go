package logs_go

import (
	"fmt"
	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"
	"io"
	"log"
	"os"
)

type LogSimple interface {
	io.Closer
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
	DPanic(format string, v ...interface{})
	Panic(format string, v ...interface{})
	Fatal(format string, v ...interface{})
}

func NewLogSimple(w io.Writer, cs []io.Closer, level zapcore.Level, errout io.Writer) LogSimple {
	l := log.New(w, "", log.Ldate|log.Ltime|log.Lshortfile)
	return &logSimple{
		l:      l,
		closes: cs,
		level:  level,
	}
}

type logSimple struct {
	l        *log.Logger
	closes   []io.Closer
	level    zapcore.Level
	errorOut io.Writer
}

func (a *logSimple) Info(format string, v ...interface{}) {
	if a.level >= zapcore.InfoLevel {
		if err := a.l.Output(2, fmt.Sprintf("[INFO] "+format+"\n", v...)); err != nil {
			if a.errorOut != nil {
				a.errorOut.Write([]byte(err.Error()))
			}
		}
	}
}

func (a *logSimple) Debug(format string, v ...interface{}) {
	if a.level >= zapcore.DebugLevel {
		if err := a.l.Output(2, fmt.Sprintf("[DEBUG] "+format+"\n", v...)); err != nil {
			if a.errorOut != nil {
				a.errorOut.Write([]byte(err.Error()))
			}
		}
	}
}

func (a *logSimple) Warn(format string, v ...interface{}) {
	if a.level >= zapcore.WarnLevel {
		if err := a.l.Output(2, fmt.Sprintf("[WARN] "+format+"\n", v)); err != nil {
			if a.errorOut != nil {
				a.errorOut.Write([]byte(err.Error()))
			}
		}
	}
}

func (a *logSimple) Error(format string, v ...interface{}) {
	if a.level >= zapcore.ErrorLevel {
		if err := a.l.Output(2, fmt.Sprintf("[ERROR] "+format+"\n", v...)); err != nil {
			if a.errorOut != nil {
				a.errorOut.Write([]byte(err.Error()))
			}
		}
	}
}

func (a *logSimple) DPanic(format string, v ...interface{}) {
	if a.level >= zapcore.DPanicLevel {
		if err := a.l.Output(2, fmt.Sprintf("[DPANIC] "+format+"\n", v...)); err != nil {
			if a.errorOut != nil {
				a.errorOut.Write([]byte(err.Error()))
			}
		}
	}
}

func (a *logSimple) Panic(format string, v ...interface{}) {
	if a.level >= zapcore.PanicLevel {
		if err := a.l.Output(2, fmt.Sprintf("[FATAL] "+format+"\n", v...)); err != nil {
			if a.errorOut != nil {
				a.errorOut.Write([]byte(err.Error()))
			}
		}
	}
}

func (a *logSimple) Fatal(format string, v ...interface{}) {
	if a.level >= zapcore.FatalLevel {
		if err := a.l.Output(2, fmt.Sprintf("[FATAL] "+format+"\n", v...)); err != nil {
			if a.errorOut != nil {
				a.errorOut.Write([]byte(err.Error()))
			}
		}
		os.Exit(1)
	}
}

func (a *logSimple) Close() error {
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

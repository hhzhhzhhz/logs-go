package logs_go

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"logs-go/formatime"
	fileout "logs-go/writer/files"
	"logs-go/writer/network"
	"logs-go/writer/stdout"
	"sort"
	"time"
)

const (
	simpleformat = "simple"
	jsonformat   = "json"
)

type writer func(c Config) (zapcore.WriteSyncer, error)

var writersFn = []writer{
	writerFileout,
	writerRsyslog,
	writerStdout,
}

func writerFileout(cfg Config) (zapcore.WriteSyncer, error) {
	var opts []fileout.Option
	if cfg.WriteFileout.GenerateRule == "" {
		return nil, nil
	}
	opts = append(opts, fileout.WithGenerationRule(cfg.WriteFileout.GenerateRule))
	opts = append(opts, fileout.WithMaxSize(cfg.WriteFileout.MaxSizeMb))
	opts = append(opts, fileout.WithMaxAge(cfg.WriteFileout.MaxAge))
	opts = append(opts, fileout.WithBufSize(cfg.WriteFileout.BufsizeMb))
	opts = append(opts, fileout.WithRotationTime(time.Duration(cfg.WriteFileout.RotationTime)))
	fw, err := fileout.NewFileout(cfg.WriteFileout.GenerateRule, opts...)
	if err != nil {
		return nil, err
	}
	return fw, nil
}

func writerRsyslog(cfg Config) (zapcore.WriteSyncer, error) {
	var opts []network.Option
	if cfg.WriteRsyslog.Addr == "" {
		return nil, nil
	}
	opts = append(opts, network.WithAddr(cfg.WriteRsyslog.Addr))
	if cfg.WriteRsyslog.NetworkTimeout > 0 {
		opts = append(opts, network.WithNetWorkTimeout(time.Duration(cfg.WriteRsyslog.NetworkTimeout)*time.Second))
	}
	opts = append(opts, network.WithLevle(network.Priority(cfg.Level.Level())))
	// rsyslog specification
	prefix := fmt.Sprintf("<%d>", network.LOG_LOCAL0+network.Priority(cfg.Level.Level()))
	opts = append(opts, network.WithCoder(network.NewRsyslogCoder(prefix)))
	sw := network.NewNetout(opts...)
	return sw, nil
}

func writerStdout(cfg Config) (zapcore.WriteSyncer, error) {
	if cfg.Stdout {
		return stdout.NewStdout(0), nil
	}
	return nil, nil

}

type Config struct {
	// json simple
	Format string `json:"format"`

	WriteFileout WriteFileout `json:"write_fileout"`

	WriteRsyslog WriteRsyslog `json:"write_rsyslog"`

	Stdout bool `json:"stdout"`
	// level
	Level zap.AtomicLevel `json:"level" yaml:"level"`
	// Development puts the jlog in development mode, which changes the
	// behavior of DPanicLevel and takes stacktraces more liberally.
	Development bool `json:"development" yaml:"development"`
	// DisableCaller stops annotating logs with the calling function's file
	// name and line number. By default, all logs are annotated.
	DisableCaller bool `json:"disableCaller" yaml:"disableCaller"`
	// DisableStacktrace completely disables automatic stacktrace capturing. By
	// default, stacktraces are captured for WarnLevel and above logs in
	// development and ErrorLevel and above in production.
	DisableStacktrace bool `json:"disableStacktrace" yaml:"disableStacktrace"`
	// EncoderConfig sets options for the chosen encoder. See
	// zapcore.EncoderConfig for details.
	EncoderConfig zapcore.EncoderConfig `json:"encoderConfig" yaml:"encoderConfig"`
	// OutputPaths is a list of URLs or file paths to write logging output to.
	// See Open for details.
	errorOut zapcore.WriteSyncer
	// InitialFields is a collection of fields to add to the root jlog.
	InitialFields map[string]interface{} `json:"initialFields" yaml:"initialFields"`
}

func NewJsonConfig() Config {
	return Config{
		Level:         zap.NewAtomicLevel(),
		Development:   false,
		EncoderConfig: NewProductionEncoderConfig(),
		errorOut:      stdout.NewStdout(0),
		Format:        jsonformat,
	}
}

func (c Config) BuildJsonLog() (LogJson, error) {
	if c.Format != jsonformat {
		return nil, fmt.Errorf("require %s output format bug %s", jsonformat, c.Format)
	}
	var writers []zapcore.WriteSyncer
	var core zapcore.Core
	var closes []io.Closer
	enc := zapcore.NewJSONEncoder(c.EncoderConfig)
	if c.Level == (zap.AtomicLevel{}) {
		return nil, fmt.Errorf("missing Level")
	}

	for _, wfn := range writersFn {
		wr, err := wfn(c)
		if err != nil {
			return nil, err
		}
		if wr != nil {
			writers = append(writers, wr)
			if close, ok := wr.(interface{ Close() error }); ok {
				closes = append(closes, close)
			}
		}
	}
	core = zapcore.NewCore(enc, zapcore.NewMultiWriteSyncer(writers...), c.Level)
	return NewLogJson(zap.New(core, c.buildOptions()...), closes), nil
}

func NewSimpleConfig() Config {
	return Config{
		Level:    zap.NewAtomicLevel(),
		errorOut: stdout.NewStdout(0),
		Format:   simpleformat,
	}
}

func (c Config) BuildSimpleLog() (LogSimple, error) {
	if c.Format != simpleformat {
		return nil, fmt.Errorf("require %s output format bug %s", simpleformat, c.Format)
	}
	var writers []zapcore.WriteSyncer
	var closes []io.Closer
	if c.Level == (zap.AtomicLevel{}) {
		return nil, fmt.Errorf("missing Level")
	}

	for _, wfn := range writersFn {
		wr, err := wfn(c)
		if err != nil {
			return nil, err
		}
		if wr != nil {
			writers = append(writers, wr)
			if close, ok := wr.(interface{ Close() error }); ok {
				closes = append(closes, close)
			}
		}
	}
	return NewLogSimple(zapcore.NewMultiWriteSyncer(writers...), closes, c.Level.Level(), c.errorOut), nil
}

func (c Config) buildOptions() []zap.Option {
	var opts []zap.Option

	if c.errorOut == nil {
		opts = append(opts, zap.ErrorOutput(stdout.NewStdout(0)))
	}

	if c.Development {
		opts = append(opts, zap.Development())
	}

	if !c.DisableCaller {
		opts = append(opts, zap.AddCaller())
	}

	stackLevel := zap.ErrorLevel
	if c.Development {
		stackLevel = zap.WarnLevel
	}
	if !c.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(stackLevel))
	}
	if len(c.InitialFields) > 0 {
		fs := make([]zap.Field, 0, len(c.InitialFields))
		keys := make([]string, 0, len(c.InitialFields))
		for k := range c.InitialFields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fs = append(fs, zap.Any(k, c.InitialFields[k]))
		}
		opts = append(opts, zap.Fields(fs...))
	}
	return opts
}

type WriteFileout struct {
	// file rule
	GenerateRule string `json:"generate_rule"`
	BufsizeMb    int    `json:"bufsize_mb"`
	MaxSizeMb    int    `json:"max_size_mb"`
	// day
	MaxAge int `json:"max_age"`
	// second
	RotationTime int `json:"rotation_time"`
}

type WriteRsyslog struct {
	// second
	NetworkTimeout int `json:"network_timeout"`
	// address
	Addr string `json:"addr"`
}

// NewProductionEncoderConfig returns an opinionated EncoderConfig for
// production environments.
func NewProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		CallerKey:      "caller",
		MessageKey:     "tag",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     formatime.RFC3339TimeEncoderKibana,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// NewTestingEncoderConfig returns an opinionated EncoderConfig for
// rsyslog
func NewTestEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "generated_time",
		LevelKey:       "level",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

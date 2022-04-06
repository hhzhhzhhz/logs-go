package example

import (
	"go.uber.org/zap"
	logs_go "logs-go"
	"testing"
)

func Test_log_file(t *testing.T) {
	cfg := logs_go.NewDefaultConfig()
	cfg.WriteFileout.GenerateRule = "./%Y-%d-%m/%H-log"
	cfg.Stdout = true
	l, err := cfg.Build()
	if err != nil {
		t.Error(err)
	}
	l.Info("Test_log_file", zap.String("out", "file"))
	l.Close()
}

func Test_log_rsyslog(t *testing.T) {
	fileds := map[string]interface{}{}
	fileds["@rsyslog_tag"] = "rsyslog_tag"
	cfg := logs_go.NewDefaultConfig()
	cfg.InitialFields = fileds
	cfg.WriteRsyslog.Addr = "127.0.0.1:65532"
	cfg.Stdout = true
	l, err := cfg.Build()
	if err != nil {
		t.Error(err)
	}
	l.Info("Test_log_rsyslog", zap.String("out", "rsyslog"))
	l.Close()
}

func Test_logs_stdout(t *testing.T) {
	cfg := logs_go.NewDefaultConfig()
	cfg.Stdout = true
	l, err := cfg.Build()
	if err != nil {
		t.Error(err)
	}
	l.Info("Test_logs_stdout", zap.String("out", "stdout"))
	l.Close()
}

package example

import (
	"go.uber.org/zap"
	logs_go "logs-go"
	"testing"
)

func Test_json(t *testing.T) {
	t.Run("file", func(t *testing.T) {
		cfg := logs_go.NewJsonConfig()
		cfg.WriteFileout.GenerateRule = "./%Y-%d-%m/%H-log"
		cfg.Stdout = true
		l, err := cfg.BuildJsonLog()
		if err != nil {
			t.Error(err)
		}
		l.Info("Test_log_file", zap.String("out", "file"))
		l.Close()
	})
	t.Run("rsyslog", func(t *testing.T) {
		fileds := map[string]interface{}{}
		fileds["@rsyslog_tag"] = "rsyslog_tag"
		cfg := logs_go.NewJsonConfig()
		cfg.InitialFields = fileds
		cfg.WriteRsyslog.Addr = "127.0.0.1:65532"
		cfg.Stdout = true
		l, err := cfg.BuildJsonLog()
		if err != nil {
			t.Error(err)
		}
		l.Info("Test_log_rsyslog", zap.String("out", "rsyslog"))
		l.Close()
	})
	t.Run("stdout", func(t *testing.T) {
		cfg := logs_go.NewJsonConfig()
		cfg.Stdout = true
		l, err := cfg.BuildJsonLog()
		if err != nil {
			t.Error(err)
		}
		l.Info("Test_logs_stdout", zap.String("out", "stdout"))
		l.Close()
	})
}

func Test_simple_stdout(t *testing.T) {
	t.Run("file", func(t *testing.T) {
		cfg := logs_go.NewSimpleConfig()
		cfg.WriteFileout.GenerateRule = "./%Y-%d-%m/%H-log"
		cfg.Stdout = true
		l, err := cfg.BuildSimpleLog()
		if err != nil {
			t.Error(err)
		}
		l.Info("Test_log_file %s", "file")
		l.Close()
	})
	t.Run("rsyslog", func(t *testing.T) {
		cfg := logs_go.NewSimpleConfig()
		cfg.WriteRsyslog.Addr = "127.0.0.1:65532"
		cfg.Stdout = true
		l, err := cfg.BuildSimpleLog()
		if err != nil {
			t.Error(err)
		}
		l.Info("Test_log_rsyslog", "rsyslog")
		l.Close()
	})
	t.Run("stdout", func(t *testing.T) {
		cfg := logs_go.NewSimpleConfig()
		cfg.Stdout = true
		l, err := cfg.BuildSimpleLog()
		if err != nil {
			t.Error(err)
		}
		l.Info("Test_logs_stdout %s", "stdout")
		l.Close()
	})
}

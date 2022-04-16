package logs_go

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
)

func Test_LogJ(t *testing.T) {
	t.Run("ryslog", func(t *testing.T) {
		t.Skip()
		fileds := map[string]interface{}{}
		fileds["@rsyslog_tag"] = "SVC_ID"
		cfg := NewLogJconfig()
		cfg.InitialFields = fileds
		cfg.WriteRsyslog.Addr = "127.0.0.1:65532"
		cfg.Stdout = true
		l, err := cfg.BuildLogJ()
		if err != nil {
			t.Error(err)
		}
		for i := 0; i < 100; i++ {
			index := i
			l.Info("Event.message", zap.String("msg_id", "message_id"), zap.Int("index", index))
		}
		l.Close()
	})

}

func Test_Logf(t *testing.T) {
	t.Run("stdout", func(t *testing.T) {
		cfg := NewLogfConfig()
		cfg.Stdout = true
		l, err := cfg.BuildLogf()
		if err != nil {
			t.Log(err.Error())
		}
		l.Info("The quick brown fox jumps over the lazy dog %s", "stdout")
		l.Close()
	})
	t.Run("disk", func(t *testing.T) {
		cfg := NewLogfConfig()
		cfg.WriteDisk.GenerateRule = "./%Y/logf"
		l, err := cfg.BuildLogf()
		if err != nil {
			t.Log(err.Error())
		}
		l.Info("The quick brown fox jumps over the lazy dog %s", "file")
		l.Close()
	})
	t.Run("syslog", func(t *testing.T) {
		t.Skip()
		cfg := NewLogfConfig()
		cfg.Stdout = true
		cfg.WriteRsyslog.Addr = "127.0.0.1:65532"
		l, err := cfg.BuildLogf()
		if err != nil {
			t.Error(err.Error())
		}
		l.Info("The quick brown fox jumps over the lazy dog %s", "ryslog")
		l.Close()
	})
	t.Run("debug", func(t *testing.T) {
		cfg := NewLogJconfig()
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		cfg.Stdout = true
		log, err := cfg.BuildLogJ()

		if err != nil {
			t.Error(err)
		}
		log.Debug("debug")
		log.Error("error")
		log.Info("info")
		log.Close()
	})

	t.Run("info", func(t *testing.T) {
		cfg := NewLogJconfig()
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		cfg.Stdout = true
		log, err := cfg.BuildLogJ()

		if err != nil {
			t.Error(err)
		}
		log.Debug("debug")
		log.Error("error")
		log.Info("info")
		log.Close()
	})

	t.Run("error", func(t *testing.T) {
		cfg := NewLogJconfig()
		cfg.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
		cfg.Stdout = true
		log, err := cfg.BuildLogJ()

		if err != nil {
			t.Error(err)
		}
		log.Debug("debug")
		log.Error("error")
		log.Info("info")
		log.Close()
	})
}

func Test_Readme(t *testing.T) {
	t.Run("rsyslog", func(t *testing.T) {
		cfg := NewLogfConfig()
		cfg.WriteRsyslog.Addr = "127.0.0.1:65532"
		cfg.Stdout = true
		l, err := cfg.BuildLogf()
		if err != nil {
			t.Error(err)
		}
		l.Info("rsyslog %s", "rsyslog")
		l.Close()
	})
	// output: 2022/04/16 14:31:37 log_test.go:185: [INFO] rsyslog rsyslog

	t.Run("rsyslog", func(t *testing.T) {
		fileds := map[string]interface{}{}
		fileds["@rsyslog_tag"] = "rsyslog_tag"
		cfg := NewLogJconfig()
		cfg.InitialFields = fileds
		cfg.WriteRsyslog.Addr = "127.0.0.1:65532"
		cfg.Stdout = true
		l, err := cfg.BuildLogJ()
		if err != nil {
			t.Error(err)
		}
		l.Info("rsyslog", zap.String("out", "rsyslog"))
		l.Close()
	})
	// output: {"level":"info","timestamp":"2022-04-16T14:31:57.338+08:00","caller":"logs-go/log.go:49","tag":"rsyslog","@rsyslog_tag":"rsyslog_tag","out":"rsyslog"}

	t.Run("disk", func(t *testing.T) {
		cfg := NewLogfConfig()
		cfg.WriteDisk.GenerateRule = "./%Y-%d-%m/%H-log"
		cfg.Stdout = true
		l, err := cfg.BuildLogf()
		if err != nil {
			t.Error(err)
		}
		l.Info("disk %s", "file")
		l.Close()
	})
	// 2022/04/16 14:28:23 log_test.go:215: [INFO] disk file

	t.Run("disk", func(t *testing.T) {
		cfg := NewLogJconfig()
		cfg.WriteDisk.GenerateRule = "./%Y-%d-%m/%H-log"
		cfg.Stdout = true
		l, err := cfg.BuildLogJ()
		if err != nil {
			t.Error(err)
		}
		l.Info("disk", zap.String("out", "file"))
		l.Close()
	})
	// output: {"level":"info","timestamp":"2022-04-16T14:28:34.688+08:00","caller":"logs-go/log.go:49","tag":"disk","out":"file"}

	t.Run("stdout", func(t *testing.T) {
		cfg := NewLogfConfig()
		cfg.Stdout = true
		l, err := cfg.BuildLogf()
		if err != nil {
			t.Error(err)
		}
		l.Info("stdout %s", "stdout")
		l.Close()
	})
	// output: 2022/04/16 14:26:33 log_test.go:240: [INFO] stdout stdout

	t.Run("stdout", func(t *testing.T) {
		cfg := NewLogJconfig()
		cfg.Stdout = true
		l, err := cfg.BuildLogJ()
		if err != nil {
			t.Error(err)
		}
		l.Info("stdout", zap.String("out", "stdout"))
		l.Close()
	})
	// output: {"level":"info","timestamp":"2022-04-16T14:26:53.911+08:00","caller":"logs-go/log.go:49","tag":"stdout","out":"stdout"}
}

func Test_Onec(t *testing.T) {
	t.Run("logj", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			DefaultLogJ().Info("test")
		}
		DefaultLogJ().Close()
	})

	t.Run("logf", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			DefaultLogJ().Info("test")
		}
		DefaultLogJ().Close()
	})
}

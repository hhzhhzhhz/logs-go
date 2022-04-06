package logs_go

import (
	"go.uber.org/zap"
	"testing"
)

func Test_Rsyslog(t *testing.T) {
	fileds := map[string]interface{}{}
	fileds["@rsyslog_tag"] = "service_id"
	cfg := NewDefaultConfig()
	//cfg.WriteFileout.GenerateRule = "./%Y/log"
	cfg.InitialFields = fileds
	cfg.WriteRsyslog.Addr = ""
	cfg.Stdout = true
	l, err := cfg.Build()
	if err != nil {
		t.Error(err)
	}
	for i := 0; i< 100; i++ {
		index := i
		//time.Sleep(1 * time.Second)
		//log.Error("message", zap.String("test_xx", "xx"), zap.Int("index", index))
		l.Info("event message", zap.String("msg_id", "message_id"), zap.Int("index", index))
	}
	l.Close()
}

func BenchmarkForFile(b *testing.B) {
	fileds := map[string]interface{}{}
	fileds["@rsyslog_tag"] = "service_id"
	cfg := NewDefaultConfig()
	cfg.WriteFileout.GenerateRule = "./%Y/log"
	cfg.InitialFields = fileds
	log, err := cfg.Build()
	if err != nil {
		b.Error(err)
	}

	b.ReportAllocs()
	b.SetParallelism(5)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next(){
			log.Info("The quick brown fox jumps over the lazy dog")
		}
	})
	log.Close()
}

func BenchmarkOneForFile(b *testing.B){
	fileds := map[string]interface{}{}
	fileds["@rsyslog_tag"] = "service_id"
	cfg := NewDefaultConfig()
	cfg.WriteFileout.GenerateRule = "./%Y/log"
	cfg.InitialFields = fileds
	log, err := cfg.Build()
	if err != nil {
		b.Error(err)
	}
	b.ReportAllocs()
	b.ResetTimer()

	for i:=0;i<b.N;i++ {
		log.Info("The quick brown fox jumps over the lazy dog")
	}
	log.Close()
}
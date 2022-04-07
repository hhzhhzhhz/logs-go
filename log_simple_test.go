package logs_go

import (
	"testing"
)

func Test_simple_stdout(t *testing.T) {
	cfg := NewSimpleConfig()
	cfg.Stdout = true
	l, err := cfg.BuildSimpleLog()
	if err != nil {
		t.Log(err.Error())
	}
	l.Info("The quick brown fox jumps over the lazy dog %s", "fileout")
	l.Close()

}

func Test_simple_fileout(t *testing.T) {
	cfg := NewSimpleConfig()
	cfg.WriteFileout.GenerateRule = "./%Y/simplelog"
	l, err := cfg.BuildSimpleLog()
	if err != nil {
		t.Log(err.Error())
	}
	l.Info("The quick brown fox jumps over the lazy dog %s", "file")
	l.Close()

}

//
func Test_simple_rsyslog(t *testing.T) {
	cfg := NewSimpleConfig()
	cfg.Stdout = true
	cfg.WriteRsyslog.Addr = "127.0.0.1:65532"
	l, err := cfg.BuildSimpleLog()
	if err != nil {
		t.Error(err.Error())
	}
	l.Info("The quick brown fox jumps over the lazy dog %s", "ryslog")
	l.Close()
}

func BenchmarkForSimpleFile(b *testing.B) {
	cfg := NewSimpleConfig()
	cfg.WriteFileout.GenerateRule = "./%Y/simple_log"
	log, err := cfg.BuildSimpleLog()
	if err != nil {
		b.Error(err)
	}

	b.ReportAllocs()
	b.SetParallelism(5)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Info("The quick brown fox jumps over the lazy dog")
		}
	})
	log.Close()
}

func BenchmarkOneForSimpleFile(b *testing.B) {
	cfg := NewSimpleConfig()
	cfg.WriteFileout.GenerateRule = "./%Y/simple_log"
	log, err := cfg.BuildSimpleLog()
	if err != nil {
		b.Error(err)
	}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		log.Info("The quick brown fox jumps over the lazy dog")
	}
	log.Close()
}

package logs_go

import "testing"

func BenchmarkLogJ(b *testing.B) {
	b.Run("disk", func(b *testing.B) {
		fileds := map[string]interface{}{}
		fileds["@rsyslog_tag"] = "service_id"
		cfg := NewLogJconfig()
		cfg.WriteDisk.Compress = false
		cfg.WriteDisk.GenerateRule = "./%Y/logj"
		cfg.InitialFields = fileds
		log, err := cfg.BuildLogJ()
		if err != nil {
			b.Error(err)
		}
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			log.Info("The quick brown fox jumps over the lazy dog")
		}
		log.Close()
	})

	b.Run("disk-Compress", func(b *testing.B) {
		fileds := map[string]interface{}{}
		fileds["@rsyslog_tag"] = "service_id"
		cfg := NewLogJconfig()
		cfg.WriteDisk.GenerateRule = "./%Y/logj"
		cfg.InitialFields = fileds
		cfg.WriteDisk.Compress = true
		log, err := cfg.BuildLogJ()
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
	})

}

func BenchmarkLogf(b *testing.B) {
	b.Run("disk-compress", func(b *testing.B) {
		cfg := NewLogJconfig()
		cfg.WriteDisk.GenerateRule = "./%Y/logf"
		cfg.WriteDisk.Compress = true
		log, err := cfg.BuildLogJ()
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
	})
}

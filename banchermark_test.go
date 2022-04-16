package logs_go

import "testing"

func BenchmarkLogJ(b *testing.B) {
	b.Run("disk", func(b *testing.B) {
		fileds := map[string]interface{}{}
		fileds["@rsyslog_tag"] = "service_id"
		cfg := NewLogJconfig()
		cfg.WriteFileout.Compress = false
		cfg.WriteFileout.GenerateRule = "./%Y/logj"
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
		cfg.WriteFileout.GenerateRule = "./%Y/logj"
		cfg.InitialFields = fileds
		cfg.WriteFileout.Compress = true
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
		cfg.WriteFileout.GenerateRule = "./%Y/logf"
		cfg.WriteFileout.Compress = true
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

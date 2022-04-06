package utils

import (
	"logs-go/strftime"
	"testing"
	"time"
)

func Benchmark_max(b *testing.B) {
	s, err := strftime.New("./log/%Y%m%d")
	t := time.Now()
	r := 10 * time.Second
	if err != nil {
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		GenRolaFileName(s, t, r, 10, false, "")
	}
}

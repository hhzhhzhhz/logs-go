package writer

import (
	"github.com/hhzhhzhhz/logs-go/strftime"
	"testing"
	"time"
)

func Benchmark_max(b *testing.B) {
	l, err := NewFileout("./log/%Y.log")
	if err != nil {
		b.Error(err)
	}
	b.Log(l.RuningInfo())
	for i := 0; i < b.N; i++ {
		if _, err := l.Write([]byte("hello world !!!\n")); err != nil {
			b.Error(err)
		}
	}
	l.Close()
}

func Benchmark_GenName(b *testing.B) {
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

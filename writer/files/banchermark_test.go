package fileout

import (
	"testing"
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
}

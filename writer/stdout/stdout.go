package stdout

import (
	"bufio"
	"os"
)

const (
	megabyte = 1024 * 1024
)

// NewStdout
func NewStdout(bufsizemb int) *Stdout {
	buf := 0
	if bufsizemb > 0 {
		buf = bufsizemb * megabyte
	}
	fr := os.Stdout
	w := bufio.NewWriterSize(fr, buf)
	return &Stdout{fr: fr, w: w}
}

type Stdout struct {
	w  *bufio.Writer
	fr *os.File
}

func (s *Stdout) Write(b []byte) (int, error) {
	return s.w.Write(b)
}

func (s *Stdout) Sync() error {
	if s.w != nil {
		s.w.Flush()
	}
	return s.fr.Sync()
}

func (s *Stdout) Close() error {
	if s.w != nil {
		s.w.Flush()
	}
	return s.fr.Close()
}

package writer

import (
	"bufio"
	"os"
)

// NewStdout
func NewStdout(bufsizemb int) *Stdout {
	buf := 512
	var w *bufio.Writer
	fr := os.Stdout
	if bufsizemb > 0 {
		buf = bufsizemb * megabyte
		w = bufio.NewWriterSize(fr, buf)
	}
	return &Stdout{fr: fr, w: w}
}

type Stdout struct {
	w  *bufio.Writer
	fr *os.File
}

func (s *Stdout) Write(b []byte) (int, error) {
	if s.w != nil {
		return s.w.Write(b)
	}
	return s.fr.Write(b)
}

func (s *Stdout) Sync() error {
	if s.w != nil {
		s.w.Flush()
	}
	return s.fr.Sync()
}

func (s *Stdout) Close() error {
	if s.w != nil {
		return s.w.Flush()
	}
	return nil
}

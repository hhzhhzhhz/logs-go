package network

import (
	"io"
)

type Network interface {
	io.WriteCloser
	//Writer() (io.Writer, error)
	ClearBuffer() error
}

type Coder interface {
	Encoder([]byte) []byte
}
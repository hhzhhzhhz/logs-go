package network

import (
	"io"
)

type Network interface {
	io.WriteCloser
	ClearBuffer() error
}

type Coder interface {
	Encoder([]byte) []byte
}

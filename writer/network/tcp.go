package network

import (
	"bufio"
	"net"
)

var megebyte = 1024 * 1024

func NewTcplog(opt *options) Network {
	r := &tcplog{
		addr: opt.addr,
		bufSize: opt.bufsize,
		coder: opt.coder,
	}
	r.connect()
	return r
}

type tcplog struct {
	addr string
	bufSize int
	w *bufio.Writer
	c net.Conn
	coder Coder
	err error
}

func (r *tcplog) Write(b []byte) (int, error) {
	if r.w == nil || r.err != nil {
		if err := r.connect(); err != nil {
			return 0, err
		}
	}

	if r.coder != nil {
		b = r.coder.Encoder(b)
	}
	n, err := r.w.Write(b)
	// todo check error
	if err != nil {
		r.err = err
	}
	return n, err
}
// bufsize
func (r *tcplog) bufsize() int {
	if r.bufSize <= 0 {
		return 1024
	}
	return r.bufSize * megebyte
}

// connect
func (r *tcplog) connect() error {
	conn, err := net.Dial("tcp", r.addr)
	if err != nil {
		r.err = err
		return err
	}
	if r.w != nil {
		r.w.Reset(conn)
	} else {
		r.w = bufio.NewWriterSize(conn, r.bufsize())
	}
	r.c = conn
	return nil
}

func (n *tcplog) Close() error {
	if n.w != nil {
		n.w.Flush()
	}
	if n.c != nil {
		return n.c.Close()
	}
	return nil
}

func (n *tcplog) ClearBuffer() error {
	if n.w != nil {
		return n.w.Flush()
	}
	return nil
}


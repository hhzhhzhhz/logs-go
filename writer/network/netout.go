package network

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"
)

type Option func(*options)

// WithwriteTimeout
//func WithwriteTimeout(timeout time.Duration) Option {
//	return func(o *options) {
//		o.writeTimeoutmill = timeout
//	}
//}

// WithNetWorkTimeout
func WithNetWorkTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.netTimeout = timeout
	}
}

// WithAddr
func WithAddr(addr string) Option {
	return func(o *options) {
		o.addr = addr
	}
}

// WithLevle
func WithLevle(level Priority) Option {
	return func(o *options) {
		o.level = level
	}
}

// WithBufsizeMb
func WithBufsizeMb(bufsize int) Option {
	return func(o *options) {
		o.bufsize = bufsize
	}
}

// WithCoder
func WithCoder(coder Coder) Option {
	return func(o *options) {
		o.coder = coder
	}
}

type options struct {
	// write timeout
	//writeTimeoutmill time.Duration
	// network timeout
	netTimeout time.Duration
	addr       string
	bufsize    int
	level      Priority
	coder      Coder
}

type netout struct {
	ctx    context.Context
	cancle context.CancelFunc
	opt    *options
	ch     chan []byte
	out    Network
	// Prevent async tasks from not exiting
	stop chan struct{}
	// Prevent asynchronous Ctrip data loss
	sw sync.WaitGroup
}

// NewNetout
func NewNetout(opts ...Option) *netout {
	opt := &options{}
	for _, f := range opts {
		f(opt)
	}
	ctx, cancle := context.WithCancel(context.Background())
	net := &netout{
		ctx:    ctx,
		cancle: cancle,
		ch:     make(chan []byte, 1024),
		opt:    opt,
		out:    NewTcplog(opt),
		stop:   make(chan struct{}, 1),
	}
	go net.running()
	return net
}

// wtimeout
//func (r *netout) wtimeout() time.Duration {
//	if r.opt.writeTimeoutmill > 0 {
//		return r.opt.writeTimeoutmill
//	}
//	return 10 * time.Millisecond
//}

// Write
func (r *netout) Write(b []byte) (int, error) {
	var buf bytes.Buffer
	buf.Write(b)
	select {
	case <-r.ctx.Done():
		return 0, fmt.Errorf("netout is closed")
	case r.ch <- buf.Bytes():
		return 0, nil
	}
	return 0, fmt.Errorf("write chan failed")
}

// running
func (r *netout) running() {
	for {
		r.sw.Add(1)
		select {
		case <-r.ctx.Done():
			r.sw.Done()
			r.stop <- struct{}{}
			return
		case b, ok := <-r.ch:
			if !ok {
				r.sw.Done()
				r.stop <- struct{}{}
				return
			}
			r.write(b)
			r.sw.Done()
		}
	}
}

// write write to network
func (r *netout) write(b []byte) (int, error) {
	n, err := r.out.Write(b)
	if err != nil {
		fmt.Println("failed to write cause: ", err.Error())
	}
	return n, nil
}

// Sync
func (r *netout) Sync() error {
	return r.out.ClearBuffer()
}

// Close
func (r *netout) Close() error {
	r.cancle()
	close(r.ch)
	<-r.stop
	r.sw.Wait()
	for b := range r.ch {
		r.write(b)
	}
	return r.out.Close()
}

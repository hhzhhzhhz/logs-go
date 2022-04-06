package network

import "testing"

func Test_rsyslog(t *testing.T) {
	c := NewTcplog(&options{
		addr: tcp_port,
	})
	if _, err := c.Write([]byte("hello world !!!")); err != nil {
		t.Error(err.Error())
	}
	c.Close()
}

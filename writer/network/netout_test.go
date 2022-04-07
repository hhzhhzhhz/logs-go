package network

import (
	"fmt"
	"net"
	"testing"
)

var (
	tcp_port = ":9876"
)

func init() {
	listen, err := net.Listen("tcp", tcp_port)
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				fmt.Println("listen accept error cause:", err.Error())
				continue
			}
			go handler(conn)

		}
		listen.Close()
	}()

}

func handler(c net.Conn) {
	buf := make([]byte, 1024*1024*4)
	for {
		c.Read(buf)
		//if err != nil {
		//	panic(err)
		//}
	}
}

func Test_Netout(t *testing.T) {
	c := NewNetout(WithAddr(tcp_port))
	if _, err := c.Write([]byte("hello world !!!")); err != nil {
		t.Error(err.Error())
	}
	c.Close()
}

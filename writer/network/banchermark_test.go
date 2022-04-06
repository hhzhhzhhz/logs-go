package network

import (
	"fmt"
	"net"
	"testing"
)


func init_() {
	listen, err := net.Listen( "tcp", tcp_port)
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				fmt.Println("listen accept error cause:", err.Error())
				break
			}
			go handler(conn)

		}
		listen.Close()
	}()

}

func Benchmark_max(b *testing.B) {
	c := NewNetout(WithAddr(tcp_port))
	for i := 0; i < b.N; i++ {
		if _, err := c.Write([]byte("hello world !!!")); err != nil {
			b.Error(err.Error())
		}
	}
	c.Close()
}

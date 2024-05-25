package proxy

import (
	"fmt"
	"net"
)

type MagicHost struct {
	Host string
	Port int
}

func Start(c chan MagicHost) {
	table := make(map[string]int)

	go func() {
		for {
			select {
			case p := <-c:
				table[p.Host] = p.Port
			}
		}
	}()

	go func() {
		addr, err := net.ResolveTCPAddr("tcp", ":8080")
		if err != nil {
			panic(err)
		}

		l, err := net.ListenTCP("tcp", addr)
		if err != nil {
			panic(err)
		}

		for {
			conn, err := l.Accept()
			if err != nil {
				fmt.Println(err)
				continue
			}

			p := table[conn.RemoteAddr().String()]
			go handle(conn, p)
		}
	}()
}

func handle(conn net.Conn, p int) {
	conn2, err := net.Dial("tcp", fmt.Sprintf(":%d", p))
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		defer conn.Close()
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println(err)
				return
			}
			conn2.Write(buf[:n])
		}
	}()

	go func() {
		defer conn2.Close()
		buf := make([]byte, 1024)
		for {
			n, err := conn2.Read(buf)
			if err != nil {
				fmt.Println(err)
				return
			}
			conn.Write(buf[:n])
		}
	}()
}

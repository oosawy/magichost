package daemon

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

type Daemon struct{}

type Args struct {
	Service string
}

type Reply struct {
	Port int
}

func (d *Daemon) Claim(args *Args, reply *Reply) error {
	return nil
}

func Do() {
	d := Daemon{}
	if err := rpc.Register(&d); err != nil {
		panic(err)
	}

	f := SocketFile()
	os.Remove(f)
	l, err := net.Listen("unix", f)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	rpc.HandleHTTP()
	if err := http.Serve(l, nil); err != nil {
		panic(err)
	}
}

func SocketFile() string {
	f := fmt.Sprintf("/run/user/%d/magichost.sock", os.Getuid())
	return f
}

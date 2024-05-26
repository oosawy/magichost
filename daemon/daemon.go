package daemon

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"

	"github.com/oosawy/magichost/proxy"
)

type Daemon struct {
	table map[string]int
}

type Args struct{}

type Reply struct{}

type ListReply struct {
	List []MagicHost
}

type MagicHost struct {
	Host string
	Port int
}

func (d *Daemon) List(args *Args, reply *ListReply) error {
	list := []MagicHost{}
	for k, v := range d.table {
		list = append(list, MagicHost{k, v})
	}
	reply.List = list

	return nil
}

type ClaimArgs struct {
	Host string
}

type ClaimReply struct {
	Port int
}

func (d *Daemon) Claim(args *ClaimArgs, reply *ClaimReply) error {
	l, err := net.ListenTCP("tcp", nil)
	if err != nil {
		return err
	}
	l.Close()
	port := l.Addr().(*net.TCPAddr).Port

	reply.Port = port

	d.table[args.Host] = port

	return nil
}

type RegisterArgs struct {
	Host string
	Port int
}

func Do() {
	println("magichost daemon starting")

	table := make(map[string]int)

	go proxy.Start(table)

	d := Daemon{table}
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

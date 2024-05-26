package daemon

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"

	"github.com/oosawy/magichost/dns"
	"github.com/oosawy/magichost/proxy"
)

type Daemon struct {
	table map[string]proxy.MagicEntry
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
		list = append(list, MagicHost{k, v.Port})
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

	stop := make(chan struct{})
	d.table[args.Host] = proxy.MagicEntry{Host: args.Host, Port: port, Stop: stop}

	go dns.Multicast(args.Host, stop)

	return nil
}

type ResolveArgs struct {
	Hostname string
}

type ResolveReply struct {
	Host string
}

func (d *Daemon) Resolve(args *ResolveArgs, reply *ResolveReply) error {
	e, ok := d.table[args.Hostname]
	if !ok {
		return fmt.Errorf("not found")
	}

	reply.Host = fmt.Sprintf("localhost:%d", e.Port)

	return nil
}

func Do() {
	println("magichost daemon starting")

	table := make(map[string]proxy.MagicEntry)

	go proxy.Start(table)

	d := Daemon{table}
	if err := rpc.Register(&d); err != nil {
		panic(err)
	}

	d.Claim(&ClaimArgs{Host: "foo.local"}, &ClaimReply{})

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

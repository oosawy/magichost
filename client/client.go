package client

import (
	"fmt"
	"net/rpc"

	"github.com/oosawy/magichost/daemon"
)

func Do(args []string) {
	switch args[0] {
	case "list":
		list()
	case "claim":
		claim(args[1])
	}
}

func list() {
	client, err := rpc.DialHTTP("unix", daemon.SocketFile())
	if err != nil {
		panic(err)
	}

	args := &daemon.Args{}
	reply := &daemon.ListReply{}

	if err = client.Call("Daemon.List", args, reply); err != nil {
		panic(err)
	}

	fmt.Println("List:", reply.List)
}

func claim(host string) {
	client, err := rpc.DialHTTP("unix", daemon.SocketFile())
	if err != nil {
		panic(err)
	}

	args := &daemon.ClaimArgs{
		Host: host,
	}
	reply := &daemon.ClaimReply{}

	if err = client.Call("Daemon.Claim", args, reply); err != nil {
		panic(err)
	}

	fmt.Println("Port:", reply.Port)
}

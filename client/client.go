package client

import (
	"fmt"
	"net/rpc"

	"github.com/oosawy/magichost/daemon"
)

func Do() {
	client, err := rpc.DialHTTP("unix", daemon.SocketFile())
	if err != nil {
		panic(err)
	}

	args := &daemon.Args{
		Service: "mh",
	}
	reply := &daemon.Reply{}

	if err = client.Call("Daemon.Claim", args, reply); err != nil {
		panic(err)
	}

	fmt.Println("Port:", reply.Port)
}

package main

import (
	"fmt"
	"os"

	"github.com/oosawy/magichost/client"
	"github.com/oosawy/magichost/daemon"
)

func main() {
	if len(os.Args) < 1 {
		fmt.Println("Usage: magichost [daemon]")
		os.Exit(1)
	}

	if len(os.Args) == 2 && (os.Args[1] == "daemon") {
		daemon.Do()
	} else {
		client.Do()
	}
}

func SocketFile() string {
	dir := fmt.Sprintf("/run/user/%d", os.Getuid())

	return dir + "/magichost.sock"
}

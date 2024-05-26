package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/oosawy/magichost/client"
	"github.com/oosawy/magichost/daemon"
	"github.com/oosawy/magichost/proxy"
)

func main() {
	if len(os.Args) < 1 {
		fmt.Println("Usage: magichost [daemon]")
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		client.Do()
		return
	}

	switch os.Args[1] {
	case "daemon":
		daemon.Do()
	case "proxy":
		c := make(chan proxy.MagicHost)
		proxy.Start(c)
		p, err := strconv.Atoi(os.Args[2])
		if err != nil {
			panic(err)
		}
		c <- proxy.MagicHost{
			Host: os.Args[3],
			Port: p,
		}

		quitChannel := make(chan os.Signal, 1)
		signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
		<-quitChannel
	}
}

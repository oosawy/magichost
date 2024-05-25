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

	switch os.Args[1] {
	case "daemon":
		daemon.Do()
	case "proxy":
		c := make(chan int)
		proxy.Start(c)
		p, err := strconv.Atoi(os.Args[2])
		if err != nil {
			panic(err)
		}
		c <- p

		quitChannel := make(chan os.Signal, 1)
		signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
		<-quitChannel
	default:
		client.Do()
	}
}

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
	if len(os.Args) < 2 {
		fmt.Println("usage: magichost [daemon]")
		fmt.Println("   or: magichost [client] (list|claim)")
		os.Exit(1)
		return
	}

	switch os.Args[1] {
	case "daemon":
		daemon.Do()
	case "proxy":
		table := map[string]int{}
		proxy.Start(table)
		p, err := strconv.Atoi(os.Args[2])
		if err != nil {
			panic(err)
		}
		table[os.Args[3]] = p

		quitChannel := make(chan os.Signal, 1)
		signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
		<-quitChannel
	case "client":
		client.Do(os.Args[2:])
	default:
		client.Do(os.Args[1:])
	}
}

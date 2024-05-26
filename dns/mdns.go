package dns

import (
	"fmt"
	"net"
	"os"

	"github.com/hashicorp/mdns"
)

func Multicast(hostname string, c chan struct{}) {
	hostName, err := os.Hostname()
	info := []string{"magichost: " + hostname}

	ips, err := net.LookupIP(hostName)
	if err != nil {
		panic(err)
	}
	service, err := mdns.NewMDNSService("foo", "_http._tcp", "local.", "foo.", 8080, ips, info)
	if err != nil {
		panic(err)
	}

	iface, err := loIf()
	if err != nil {
		panic(err)
	}

	server, err := mdns.NewServer(&mdns.Config{Zone: service, ForceUnicastResponses: true, LogEmptyResponses: false, Iface: iface})
	if err != nil {
		panic(err)
	}

	fmt.Println("mdns start")

	<-c
	defer server.Shutdown()
	defer func() {
		fmt.Println("mdns shutdown")
	}()
}

// func allIPs() ([]net.IP, error) {
// 	var ips []net.IP

// 	addrs, err := net.InterfaceAddrs()
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, ip := range addrs {
// 		ips = append(ips, ip.(*net.IPNet).IP)
// 	}

// 	return ips, nil
// }

func loIf() (*net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback != 0 {
			return &iface, nil
		}
	}
	return nil, fmt.Errorf("no loopback interface found")
}

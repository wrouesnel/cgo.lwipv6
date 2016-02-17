/*
	Connects to the VDE switch specified on the command line
 */
package main

import (
	"github.com/wrouesnel/cgo.lwipv6"
	"flag"
	"github.com/wrouesnel/go.log"
	"net"
	"fmt"
)

var (
	socket = flag.String("sock", "/tmp/vde.ctl", "Unix socket path to connect to VDE interface")
)

func main() {
	flag.Parse()

	lwipv6.Initialize(0)
	defer lwipv6.Finish()

	intf := lwipv6.CreateInterface(lwipv6.IF_VDE, *socket, 0)
	if intf == nil {
		log.Fatalln("Interface setup failed to:", *socket)
	}

	ip := net.IPv4(192,168,123,2)
	mask := net.IPMask(net.IPv4(255,255,255,0))
	if intf.AddAddress(net.IPNet{ip, mask}) != 0 {
		log.Fatalln("Adding address failed.")
	}

	if intf.IfUp(0) != 0 {
		log.Fatalln("Failed bringing interface up")
	}

	fmt.Println("Interface up successfully! Press any key to exit.")
	var input string
	fmt.Scanln(&input)
	fmt.Println("Exiting.")
}
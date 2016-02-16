// Go bindings for the LWIPv6 TCP/IP hybrid stack.
package lwipv6

/*
#cgo LDFLAGS: -llwipv6
#include <stddef.h>
#include <lwipv6.h>
*/
import "C"

import (
	//"encoding/binary"
)

//func IP4_ADDRX(ip4ax C.struct_ip_addr, a, b, c, d uint8)  {
//
//	((uint32)((a) & 0xff) << 24) | ((uint32)((b) & 0xff) << 16) | ((uint32)((c) & 0xff) << 8) | (uint32)((d) & 0xff)
//}
//
//func IP64_PREFIX () {
//	uint32(0xffff)
//
//	binary.BigEndian.
//}
//
//func IP6_ADDR(ipaddr, a,b,c,d,e,f,g,h) {
//
//}
//
//func IP64_MASKADDR(ipaddr, a,b,c,d) {
//
//}
//
//func IP_ADDR_IS_V4(ipaddr) {
//
//}

const (
	LWIP_STACK_FLAG_FORWARDING 	= 1
	LWIP_STACK_FLAG_USERFILTER 	= 0x2
	LWIP_STACK_FLAG_UF_NAT 		= 0x10000
)

const (
	/* Allows binding to TCP/UDP sockets below 1024 */
	LWIP_CAP_NET_BIND_SERVICE = 1<<10
	/* Allow broadcasting, listen to multicast */
	LWIP_CAP_NET_BROADCAST    = 1<<11
	/* Allow interface configuration */
	LWIP_CAP_NET_ADMIN        = 1<<12
	/* Allow use of RAW sockets */
	/* Allow use of PACKET sockets */
	LWIP_CAP_NET_RAW          = 1<<13
)

func lwip_init() {
	C.lwip_init()
}

func lwip_fini() {
	C.lwip_fini()
}


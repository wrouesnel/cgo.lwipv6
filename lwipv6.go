/*
 	Go bindings for the LWIPv6 TCP/IP hybrid stack.
 	The LWIPv6 stack is unclear on the ramifications of using multiple stacks, and Go's
 	concurrency story makes it even more unclear. Until we figure out some good test
 	scenarios, this library assumes you will be using exactly 1 stack and it's the default.

 	For simplicity, some redundant runtime reconfigure options are omitted for the time being.
  */
package lwipv6

/*
#cgo LDFLAGS: -llwipv6
#include <stdlib.h>
#include <stddef.h>
#include <lwipv6.h>
*/
import "C"

import (
	//"encoding/binary"
	"net"
	"encoding/binary"
	"unsafe"
)

type LWIPStackFlag uint32
const (
	LWIP_STACK_FLAG_FORWARDING LWIPStackFlag	= 1
	LWIP_STACK_FLAG_USERFILTER LWIPStackFlag	= 0x2
	LWIP_STACK_FLAG_UF_NAT 	   LWIPStackFlag 	= 0x10000
)

type LWIPStackCapability uint32
const (
	/* Allows binding to TCP/UDP sockets below 1024 */
	LWIP_CAP_NET_BIND_SERVICE LWIPStackCapability = 1<<10
	/* Allow broadcasting, listen to multicast */
	LWIP_CAP_NET_BROADCAST    LWIPStackCapability = 1<<11
	/* Allow interface configuration */
	LWIP_CAP_NET_ADMIN        LWIPStackCapability = 1<<12
	/* Allow use of RAW sockets */
	/* Allow use of PACKET sockets */
	LWIP_CAP_NET_RAW          LWIPStackCapability = 1<<13
)

type SlirpFlags uint32
const (
	SLIRP_LISTEN_UDP 		SlirpFlags = 0x1000
	SLIRP_LISTEN_TCP 		SlirpFlags = 0x2000
	SLIRP_LISTEN_UNIXSTREAM SlirpFlags = 0x3000
	SLIRP_LISTEN_TYPEMASK 	SlirpFlags = 0x7000
	SLIRP_LISTEN_ONCE 		SlirpFlags = 0x8000
)

var (
	stack *C.struct_stack
)

func IsInitialized() bool {
	return (stack != nil)
}

// Initialize the network stack. Multiple calls lead to a stack panic.
func Initialize(flags LWIPStackFlag) {
	if stack != nil {
		panic("LWIPv6 already initialized.")
	}

	// Initialize LWIPv6
	C.lwip_init()
	// Setup the default network stack
	stack = C.lwip_add_stack(C.ulong(flags))
	C.lwip_stack_set(stack)
}

// Finalizes and closes down the stack. This should only be done after all
// network operatons are finished.
func Finish() {
	C.lwip_stack_set(nil)
	C.lwip_del_stack(stack)
	stack = nil
	C.lwip_fini()
}

// Wraps an LWIP network interface. Returned by interface add methods.
type networkInterface struct {
	netif *C.struct_netif
}

// Convert a Go-lang IP address to an LWIP compatible one
func netIPToLWIP(netAddr net.IP) *C.struct_ip_addr {
	longIP := netAddr.To16()
	lwipAddr := new(C.struct_ip_addr)

	for idx := 0; idx < 4; idx++ {
		lwipAddr.addr[idx] = C.uint32_t(binary.BigEndian.Uint32(longIP[idx:idx+4]))
	}

	return lwipAddr
}

// Convert an LWIP compatible address to a net.IP
func lwipAddrToNetIP(lwipAddr C.struct_ip_addr) {

}

type LWIPInterfaceType int
const (
	IF_VDE		LWIPInterfaceType = iota
	IF_TAP		LWIPInterfaceType = iota
	IF_TUN		LWIPInterfaceType = iota
	IF_SLIRP	LWIPInterfaceType = iota
)

// Creates a virtual network interface
func CreateInterface(ifType LWIPInterfaceType, arg string, flags int ) *networkInterface {
	var netif *C.struct_netif
	lwipArg := C.CString(arg)	// FIXME: possible memory leak

	switch ifType {
	case IF_VDE:
		netif = C.lwip_add_vdeif(stack, unsafe.Pointer(lwipArg), C.int(flags))
	case IF_TAP:
		netif = C.lwip_add_tapif(stack, unsafe.Pointer(lwipArg), C.int(flags))
	case IF_TUN:
		netif = C.lwip_add_tunif(stack, unsafe.Pointer(lwipArg), C.int(flags))
	case IF_SLIRP:
		netif = C.lwip_add_slirpif(stack, unsafe.Pointer(lwipArg), C.int(flags))
	default:
		return nil
	}

	return &networkInterface{
		netif : netif,
	}
}

// Add an address to an interface
func (this *networkInterface) AddAddress(addr net.IPNet) int {
	lwipAddr := netIPToLWIP(addr.IP)
	lwipNetmask := netIPToLWIP((net.IP)(addr.Mask))

	return int(C.lwip_add_addr(this.netif, lwipAddr, lwipNetmask))
}

// Remove an address from an interface
func (this *networkInterface) DelAddress(addr net.IPNet) int {
	lwipAddr := netIPToLWIP(addr.IP)
	lwipNetmask := netIPToLWIP((net.IP)(addr.Mask))

	return int(C.lwip_del_addr(this.netif, lwipAddr, lwipNetmask))
}

func (this *networkInterface) IfUp(flags uint32) int {
	return int(C.lwip_ifup_flags(this.netif, C.int(flags)))
}

func (this *networkInterface) IfDown() int {
	return int(C.lwip_ifdown(this.netif))
}

/*
The decaying list of functions to implement

int lwip_add_route(struct stack *stack, struct ip_addr *addr, struct ip_addr *netmask, struct ip_addr *nexthop, struct netif *netif, int flags);
int lwip_del_route(struct stack *stack, struct ip_addr *addr, struct ip_addr *netmask, struct ip_addr *nexthop, struct netif *netif, int flags);

int lwip_accept(int s, struct sockaddr *addr, socklen_t *addrlen);
int lwip_bind(int s, const struct sockaddr *name, socklen_t namelen);
int lwip_shutdown(int s, int how);
int lwip_getpeername (int s, struct sockaddr *name, socklen_t *namelen);
int lwip_getsockname (int s, struct sockaddr *name, socklen_t *namelen);
int lwip_getsockopt (int s, int level, int optname, void *optval, socklen_t *optlen);
int lwip_setsockopt (int s, int level, int optname, const void *optval, socklen_t optlen);
int lwip_close(int s);
int lwip_connect(int s, const struct sockaddr *name, socklen_t namelen);
int lwip_listen(int s, int backlog);
ssize_t lwip_recv(int s, void *mem, int len, unsigned int flags);
ssize_t lwip_read(int s, void *mem, int len);
ssize_t lwip_recvfrom(int s, void *mem, int len, unsigned int flags,
struct sockaddr *from, socklen_t *fromlen);
ssize_t lwip_send(int s, const void *dataptr, int size, unsigned int flags);
ssize_t lwip_sendto(int s, const void *dataptr, int size, unsigned int flags,
const struct sockaddr *to, socklen_t tolen);
ssize_t lwip_recvmsg(int fd, struct msghdr *msg, int flags);
ssize_t lwip_sendmsg(int fd, const struct msghdr *msg, int flags);

int lwip_msocket(struct stack *stack, int domain, int type, int protocol);
int lwip_socket(int domain, int type, int protocol);
ssize_t lwip_write(int s, void *dataptr, int size);
int lwip_select(int maxfdp1, fd_set *readset, fd_set *writeset, fd_set *exceptset,
struct timeval *timeout);
int lwip_pselect(int maxfdp1, fd_set *readset, fd_set *writeset, fd_set *exceptset,
const struct timespec *timeout, const sigset_t *sigmask);
int lwip_poll(struct pollfd *fds, nfds_t nfds, int timeout);
int lwip_ppoll(struct pollfd *fds, nfds_t nfds,
const struct timespec *timeout, const sigset_t *sigmask);

int lwip_ioctl(int s, long cmd, void *argp);
int lwip_fcntl64(int s, int cmd, long arg);
int lwip_fcntl(int s, int cmd, long arg);

struct iovec;
ssize_t lwip_writev(int s, struct iovec *vector, int count);
ssize_t lwip_readv(int s, struct iovec *vector, int count);

void lwip_radv_load_config(struct stack *stack,FILE *filein);
int lwip_radv_load_configfile(struct stack *stack,void *arg);

int lwip_event_subscribe(lwipvoidfun cb, void *arg, int fd, int how);
*/
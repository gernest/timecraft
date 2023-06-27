//go:build wasip1

package timecraft

import (
	"context"
	"net"
	"syscall"
	"unsafe"

	"github.com/stealthrocket/net/wasip1"
)

//go:wasmimport wasi_snapshot_preview1 sock_setsockopt
//go:noescape
func setsockopt(fd int32, level uint32, name uint32, value unsafe.Pointer, valueLen uint32) syscall.Errno

// DialTLS connects to the given network address and performs a TLS handshake
// with the host. The resulting connection transparently performs encryption.
func DialTLS(ctx context.Context, network, addr string) (net.Conn, error) {
	switch network {
	case "tcp", "tcp4", "tcp6", "unix":
	default:
		return nil, net.UnknownNetworkError(network)
	}
	hostname, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	host := []byte(hostname)
	conn, err := wasip1.DialContext(ctx, network, addr)
	if err != nil {
		return nil, err
	}

	syscallConn := conn.(syscall.Conn)
	rawConn, err := syscallConn.SyscallConn()
	if err != nil {
		conn.Close()
		return nil, err
	}

	rawConn.Control(func(fd uintptr) {
		setsockopt(int32(fd), 0x74696d65, 1, unsafe.Pointer(unsafe.SliceData(host)), uint32(len(hostname)))
	})

	return conn, nil
}

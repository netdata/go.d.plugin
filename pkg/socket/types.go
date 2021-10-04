package socket

import (
	"crypto/tls"
	"time"
)

// Processor function passed to the Socket.Command function.
// It is passed by the caller to process a command's response
// line by line.
type Processor func([]byte) bool

// Config holds the network ip v4 or v6 address, port,
// Socket type(ip, tcp, udp, unix), timeout and TLS configuration
// for a Socket
type Config struct {
	Network Network
	Address string
	Timeout time.Duration
	TLSConf *tls.Config
}

// Network is a string alias for the available Socket types.
type Network string

const (
	// NetworkIP is used for IP sockets
	NetworkIP Network = "ip"
	// NetworkTCP is used for TCP sockets
	NetworkTCP Network = "tcp"
	// NetworkUDP is used for UDP sockets
	NetworkUDP Network = "udp"
	// NetworkUnix is used for UNIX sockets
	NetworkUnix Network = "unix"
)

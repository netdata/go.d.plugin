package socket

import "strings"

func IsUnixSocket(address string) bool {
	return strings.HasPrefix(address, "/") || strings.HasPrefix(address, "unix://")
}

func IsUdpSocket(address string) bool {
	return strings.HasPrefix(address, "udp://")
}

func networkType(address string) string {
	switch {
	case IsUnixSocket(address):
		return "unix"
	case IsUdpSocket(address):
		return "udp"
	default:
		return "tcp"
	}
}

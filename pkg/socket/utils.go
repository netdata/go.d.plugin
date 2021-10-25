package socket

import "strings"

func IsUnixSocket(address string) bool {
	return strings.HasPrefix(address, "/") || strings.HasPrefix(address, "unix://")
}

func IsUdpSocket(address string) bool {
	return strings.HasPrefix(address, "udp://")
}

func networkType(address string) Network {
	switch {
	case IsUnixSocket(address):
		return NetworkUnix
	case IsUdpSocket(address):
		return NetworkUDP
	default:
		return NetworkTCP
	}
}

package socket

import "strings"

func IsUnixSocket(address string) bool {
	return strings.HasPrefix(address, "/")
}

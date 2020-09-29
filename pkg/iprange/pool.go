package iprange

import (
	"math/big"
	"net"
	"strings"
)

type Pool []Range

func (p Pool) String() string {
	var b strings.Builder
	for _, r := range p {
		b.WriteString(r.String() + ",")
	}
	return strings.TrimSuffix(b.String(), ",")
}

func (p Pool) Size() *big.Int {
	size := big.NewInt(0)
	for _, r := range p {
		size.Add(size, r.Size())
	}
	return size
}

func (p Pool) Contains(ip net.IP) bool {
	for _, r := range p {
		if r.Contains(ip) {
			return true
		}
	}
	return false
}

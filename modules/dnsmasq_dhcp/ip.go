package dnsmasq_dhcp

import (
	"bytes"
	"fmt"
	"math/big"
	"net"
	"strings"
)

type ipNet struct {
	net *net.IPNet
}

func (n ipNet) String() string {
	return n.net.String()
}

func (n ipNet) contains(ip net.IP) bool {
	return n.net.Contains(ip)
}

func (n ipNet) size() *big.Int {
	ones, bits := n.net.Mask.Size()
	b := big.NewInt(1)
	for i := 0; i < bits-ones; i++ {
		b.Mul(b, big.NewInt(2))
	}
	return b
}

func parseIPNetwork(s string) (*ipNet, error) {
	parts := strings.Split(s, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("bad ip network format : %s", s)
	}

	_, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		return nil, fmt.Errorf("error on parsing %s : %v", s, err)
	}

	return &ipNet{ipnet}, nil
}

type ipRange struct {
	start net.IP
	end   net.IP
}

func (r ipRange) String() string {
	return fmt.Sprintf("%s-%s", r.start, r.end)
}

func (r ipRange) contains(ip net.IP) bool {
	inLower := bytes.Compare(ip, r.start) >= 0
	if !inLower {
		return false
	}
	inUpper := bytes.Compare(ip, r.end) <= 0
	return inLower && inUpper
}

func parseIPRange(s string) (*ipRange, error) {
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("bad ip range format : %s", s)
	}

	iprange := &ipRange{start: net.ParseIP(parts[0]), end: net.ParseIP(parts[1])}

	if !iprange.isValid() {
		return nil, fmt.Errorf("bad ip range format : %s", s)
	}

	return iprange, nil
}

func (r ipRange) isValid() bool {
	sameFam := (r.start.To4() == nil) == (r.end.To4() == nil)
	return r.start != nil && r.end != nil && bytes.Compare(r.end, r.start) >= 0 && sameFam
}

func (r ipRange) size() *big.Int {
	if !r.isValid() {
		big.NewInt(0)
	}
	start, end := r.start.To4(), r.end.To4()
	if start != nil && end != nil {
		return big.NewInt(int64(ipv4ToInt(end)) - int64(ipv4ToInt(start)) + 1)
	}

	return big.NewInt(0).Add(
		big.NewInt(0).Sub(
			big.NewInt(0).SetBytes(r.end),
			big.NewInt(0).SetBytes(r.start),
		),
		big.NewInt(1),
	)
}

func ipv4ToInt(ip net.IP) int32 {
	return int32(ip[0])<<24 |
		int32(ip[1])<<16 |
		int32(ip[2])<<8 |
		int32(ip[3])
}

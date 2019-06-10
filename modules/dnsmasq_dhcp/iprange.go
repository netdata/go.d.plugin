package dnsmasq_dhcp

import (
	"bytes"
	"fmt"
	"net"
	"strings"
)

type ipRange struct {
	start net.IP
	end   net.IP
}

func (r ipRange) Contains(ip net.IP) bool {
	inLower := bytes.Compare(ip, r.start) >= 0
	if !inLower {
		return false
	}
	inUpper := bytes.Compare(ip, r.end) <= 0
	return inLower && inUpper
}

func (r ipRange) isValid() bool {
	return r.start != nil && r.end != nil && bytes.Compare(r.end, r.start) >= 0
}

func (r ipRange) numOfIPs() int {
	if !r.isValid() {
		return -1
	}
	var num = 1
	ip := net.ParseIP(r.start.String())
	for !ip.Equal(r.end) {
		incIP(ip)
		num++
	}
	return num
}

func parseIPRange(s string) (*ipRange, error) {
	parts := strings.Split(s, ",")
	if len(parts) == 1 {
		parts = strings.Split(s, "-")
	}
	if len(parts) != 2 {
		return nil, fmt.Errorf("bad ip range format : %s", s)
	}

	start := net.ParseIP(strings.TrimSpace(parts[0]))
	end := net.ParseIP(strings.TrimSpace(parts[1]))

	r := &ipRange{start: start, end: end}

	if !r.isValid() {
		return nil, fmt.Errorf("bad ip range format : %s", s)
	}

	return r, nil
}

func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

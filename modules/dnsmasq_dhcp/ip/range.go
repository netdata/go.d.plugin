package ip

import (
	"bytes"
	"fmt"
	"math/big"
	"net"
	"strings"
)

// Type is type of Range.
type Type int

const (
	InvalidType Type = iota
	V4Type
	V6Type
)

func NewRange(s string) *Range {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return nil
	}

	r := Range{
		start: net.ParseIP(parts[0]),
		end:   net.ParseIP(parts[1]),
	}
	if !isRangeValid(r) {
		return nil
	}

	return &r
}

// Range represents IP Range.
type Range struct {
	start net.IP
	end   net.IP
}

// String returns Range string representation.
func (r Range) String() string {
	return fmt.Sprintf("%s-%s", r.start, r.end)
}

// Type returns Range IP type.
func (r Range) Type() Type {
	if r.start.To4() != nil && r.end.To4() != nil {
		return V4Type
	}
	if r.start.To16() != nil && r.end.To16() != nil {
		return V6Type
	}
	return InvalidType
}

// Contains reports whether net.IP is within Range.
func (r Range) Contains(ip net.IP) bool {
	inLower := bytes.Compare(ip, r.start) >= 0
	if !inLower {
		return false
	}
	inUpper := bytes.Compare(ip, r.end) <= 0
	return inLower && inUpper
}

// Hosts returns number of hosts addresses in the Range.
func (r Range) Hosts() *big.Int {
	switch r.Type() {
	default:
		return big.NewInt(0)
	case V4Type:
		return v4RangeSize(r)
	case V6Type:
		return v6RangeSize(r)
	}
}

// isRangeValid reports if the Range is valid.
func isRangeValid(r Range) bool {
	return r.Type() != InvalidType && bytes.Compare(r.end, r.start) >= 0
}

// v4RangeSize return ipv4 Range size.
func v4RangeSize(r Range) *big.Int {
	return big.NewInt(int64(v4ToInt(r.end)) - int64(v4ToInt(r.start)) + 1)
}

// v6RangeSize return ipv6 Range size.
func v6RangeSize(r Range) *big.Int {
	return big.NewInt(0).Add(
		big.NewInt(0).Sub(big.NewInt(0).SetBytes(r.end), big.NewInt(0).SetBytes(r.start)),
		big.NewInt(1),
	)
}

// v4ToInt converts net.IP to int32.
func v4ToInt(ip net.IP) int32 {
	ip = ip.To4()
	return int32(ip[0])<<24 | int32(ip[1])<<16 | int32(ip[2])<<8 | int32(ip[3])
}

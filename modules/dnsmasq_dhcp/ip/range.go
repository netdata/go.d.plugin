package ip

import (
	"bytes"
	"fmt"
	"math/big"
	"net"
	"strings"
)

// Range represents IP Range.
type Range struct {
	Start net.IP
	End   net.IP
}

// String returns Range string representation.
func (r Range) String() string { return fmt.Sprintf("%s-%s", r.Start, r.End) }

// Type returns Range IP type.
func (r Range) Type() Type {
	if r.Start.To4() != nil && r.End.To4() != nil {
		return V4Type
	}
	if r.Start.To16() != nil && r.End.To16() != nil {
		return V6Type
	}
	return UnknownType
}

// Contains reports whether net.IP is within Range.
func (r Range) Contains(ip net.IP) bool {
	inLower := bytes.Compare(ip, r.Start) >= 0
	if !inLower {
		return false
	}
	inUpper := bytes.Compare(ip, r.End) <= 0
	return inLower && inUpper
}

// Hosts returns number of hosts addresses in the Range.
func (r Range) Hosts() *big.Int {
	switch r.Type() {
	default:
		return big.NewInt(0)
	case V4Type:
		return V4RangeSize(r)
	case V6Type:
		return V6RangeSize(r)
	}
}

// IsRangeValid reports if Range is valid.
func IsRangeValid(r Range) bool {
	return r.Type() != UnknownType && bytes.Compare(r.End, r.Start) >= 0
}

// V4RangeSize return ipv4 Range size.
func V4RangeSize(r Range) *big.Int {
	return big.NewInt(int64(V4ToInt(r.End)) - int64(V4ToInt(r.Start)) + 1)
}

// V6RangeSize return ipv6 Range size.
func V6RangeSize(r Range) *big.Int {
	return big.NewInt(0).Add(
		big.NewInt(0).Sub(big.NewInt(0).SetBytes(r.End), big.NewInt(0).SetBytes(r.Start)),
		big.NewInt(1),
	)
}

// V4ToInt converts net.IP to int32.
func V4ToInt(ip net.IP) int32 {
	return int32(ip[0])<<24 | int32(ip[1])<<16 | int32(ip[2])<<8 | int32(ip[3])
}

func ParseRange(s string) *Range {
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		return nil
	}

	r := Range{
		Start: net.ParseIP(parts[0]),
		End:   net.ParseIP(parts[1]),
	}
	if !IsRangeValid(r) {
		return nil
	}

	return &r
}

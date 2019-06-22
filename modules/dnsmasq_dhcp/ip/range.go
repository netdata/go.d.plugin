package ip

import (
	"bytes"
	"fmt"
	"math/big"
	"net"
	"strings"
)

// Family represents IP Range family.
type Family uint8

const (
	InvalidFamily Family = iota
	V4Family
	V6Family
)

type IRange interface {
	Family() Family
	Contains(ip net.IP) bool
	Hosts() *big.Int
	fmt.Stringer
}

const separator = " "

var replacer = strings.NewReplacer(
	separator, "",
	",", separator,
	"-", separator,
)

// ParseRange parses s as an IP Range, returning the result.
// If s is not a valid textual representation of an IP Range,
// ParseRange returns nil.
func ParseRange(s string) IRange {
	if s == "" {
		return nil
	}

	var start, end string

	s = replacer.Replace(s)
	parts := strings.Split(s, separator)

	switch len(parts) {
	default:
		return nil
	case 1:
		start, end = parts[0], parts[0]
	case 2:
		start, end = parts[0], parts[1]
	}

	r := Range{
		Start: net.ParseIP(start),
		End:   net.ParseIP(end),
	}
	if !isRangeValid(r) {
		return nil
	}

	return &r
}

// Range represents IP Range.
type Range struct {
	Start net.IP
	End   net.IP
}

// String returns Range string representation.
func (r Range) String() string {
	return fmt.Sprintf("%s-%s", r.Start, r.End)
}

// Family returns IP Range family.
func (r Range) Family() Family {
	start := ipAddressFamily(r.Start)
	end := ipAddressFamily(r.End)
	if start != end || start == InvalidFamily {
		return InvalidFamily
	}
	return start
}

// Contains reports whether net.IP is within Range.
func (r Range) Contains(ip net.IP) bool {
	// in [lower:upper]
	return bytes.Compare(ip, r.Start) >= 0 && bytes.Compare(ip, r.End) <= 0
}

// Hosts returns number of hosts addresses in the Range.
// Hosts returns nil if Range is not valid.
func (r Range) Hosts() *big.Int {
	switch r.Family() {
	default:
		return nil
	case V4Family:
		return v4RangeSize(r)
	case V6Family:
		return v6RangeSize(r)
	}
}

// ipAddressFamily returns IP address family.
func ipAddressFamily(ip net.IP) Family {
	if ip.To16() == nil {
		return InvalidFamily
	}
	if ip.To4() == nil {
		return V6Family
	}
	return V4Family
}

// isRangeValid reports if the Range is valid.
func isRangeValid(r Range) bool {
	return r.Family() != InvalidFamily && bytes.Compare(r.End, r.Start) >= 0
}

// v4RangeSize returns ipv4 Range size.
func v4RangeSize(r Range) *big.Int {
	return big.NewInt(v4ToInt(r.End) - v4ToInt(r.Start) + 1)
}

// v6RangeSize returns ipv6 Range size.
func v6RangeSize(r Range) *big.Int {
	return big.NewInt(0).Add(
		big.NewInt(0).Sub(big.NewInt(0).SetBytes(r.End), big.NewInt(0).SetBytes(r.Start)),
		big.NewInt(1),
	)
}

// v4ToInt converts net.IP to int64.
func v4ToInt(ip net.IP) int64 {
	ip = ip.To4()
	return int64(ip[0])<<24 | int64(ip[1])<<16 | int64(ip[2])<<8 | int64(ip[3])
}

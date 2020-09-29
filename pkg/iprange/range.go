package iprange

import (
	"bytes"
	"fmt"
	"math/big"
	"net"
)

// Family represents IP Range family.
type Family uint8

const (
	V4Family Family = iota
	V6Family
)

type Range interface {
	Family() Family
	Contains(ip net.IP) bool
	Hosts() *big.Int
	fmt.Stringer
}

type v4Range struct {
	start net.IP
	end   net.IP
}

func (r v4Range) String() string          { return fmt.Sprintf("%s-%s", r.start, r.end) }
func (v4Range) Family() Family            { return V4Family }
func (r v4Range) Contains(ip net.IP) bool { return contains(r.start, r.end, ip) }
func (r v4Range) Hosts() *big.Int         { return v4RangeSize(r.start, r.end) }

type v6Range struct {
	start net.IP
	end   net.IP
}

func (r v6Range) String() string          { return fmt.Sprintf("%s-%s", r.start, r.end) }
func (v6Range) Family() Family            { return V6Family }
func (r v6Range) Contains(ip net.IP) bool { return contains(r.start, r.end, ip) }
func (r v6Range) Hosts() *big.Int         { return v6RangeSize(r.start, r.end) }

func contains(start, end, ip net.IP) bool {
	return bytes.Compare(ip, start) >= 0 && bytes.Compare(ip, end) <= 0
}

func v4RangeSize(start, end net.IP) *big.Int {
	return big.NewInt(v4ToInt(end) - v4ToInt(start) + 1)
}

func v6RangeSize(start, end net.IP) *big.Int {
	return big.NewInt(0).Add(
		big.NewInt(0).Sub(big.NewInt(0).SetBytes(end), big.NewInt(0).SetBytes(start)),
		big.NewInt(1),
	)
}

func v4ToInt(ip net.IP) int64 {
	ip = ip.To4()
	return int64(ip[0])<<24 | int64(ip[1])<<16 | int64(ip[2])<<8 | int64(ip[3])
}

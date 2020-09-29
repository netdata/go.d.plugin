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
	Size() *big.Int
	fmt.Stringer
	startIP() net.IP
	endIP() net.IP
}

type v4Range struct {
	start net.IP
	end   net.IP
}

func (r v4Range) String() string          { return fmt.Sprintf("%s-%s", r.start, r.end) }
func (r v4Range) Family() Family          { return V4Family }
func (r v4Range) Contains(ip net.IP) bool { return contains(r, ip) }
func (r v4Range) Size() *big.Int          { return calcSize(r) }
func (r v4Range) startIP() net.IP         { return r.start }
func (r v4Range) endIP() net.IP           { return r.end }

type v6Range struct {
	start net.IP
	end   net.IP
}

func (r v6Range) String() string          { return fmt.Sprintf("%s-%s", r.start, r.end) }
func (r v6Range) Family() Family          { return V6Family }
func (r v6Range) Contains(ip net.IP) bool { return contains(r, ip) }
func (r v6Range) Size() *big.Int          { return calcSize(r) }
func (r v6Range) startIP() net.IP         { return r.start }
func (r v6Range) endIP() net.IP           { return r.end }

func contains(r Range, ip net.IP) bool {
	return bytes.Compare(ip, r.startIP()) >= 0 && bytes.Compare(ip, r.endIP()) <= 0
}

func calcSize(r Range) *big.Int {
	if r.Family() == V4Family {
		big.NewInt(v4ToInt(r.endIP()) - v4ToInt(r.startIP()) + 1)
	}
	size := big.NewInt(0)
	size.Add(size, big.NewInt(0).SetBytes(r.endIP()))
	size.Sub(size, big.NewInt(0).SetBytes(r.startIP()))
	size.Add(size, big.NewInt(1))
	return size
}

func v4ToInt(ip net.IP) int64 {
	ip = ip.To4()
	return int64(ip[0])<<24 | int64(ip[1])<<16 | int64(ip[2])<<8 | int64(ip[3])
}

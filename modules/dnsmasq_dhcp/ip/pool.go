package ip

import (
	"fmt"
	"math/big"
	"net"
	"strings"
)

// Type represents type of DHCP pool.
type Type int

const (
	UnknownType Type = iota
	V4Type
	V6Type
)

// IRange is implemented by IP Ranges.
type IRange interface {
	Type() Type
	Hosts() *big.Int
	Contains(ip net.IP) bool
	fmt.Stringer
}

// NewPool creates new Pool.
func NewPool(iprange string) *Pool {
	r := ParseRange(iprange)
	if r == nil {
		return nil
	}

	return &Pool{
		IRange: r,
		leases: make(map[string]bool),
	}
}

type Pool struct {
	IRange
	leases map[string]bool
}

// NumOfLeases returns number of active dhcp leases.
func (p Pool) NumOfLeases() int64 { return int64(len(p.leases)) }

// Lease adds IP to leases database.
func (p *Pool) Lease(ip net.IP) { p.leases[ip.String()] = true }

// ResetLeases resets all leases.
func (p *Pool) ResetLeases() { p.leases = make(map[string]bool) }

// Utilization returns pool utilization in percent.
func (p Pool) Utilization() float64 {
	total := p.Hosts()
	if !total.IsInt64() {
		return 0
	}
	return float64(p.NumOfLeases()) * 100 / float64(total.Int64())
}

func ParseRange(iprange string) IRange {
	if strings.Contains(iprange, "-") {
		return NewRange(iprange)
	}
	return NewNet(iprange)
}

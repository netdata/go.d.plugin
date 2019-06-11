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

// IRange implements IP Range.
type IRange interface {
	Type() Type
	Hosts() *big.Int
	Contains(ip net.IP) bool
	fmt.Stringer
}

// RawDHCPPool RawDHCPPool.
type RawDHCPPool struct {
	Name  string `yaml:"name"`
	Range string `yaml:"range"`
}

// NewDHCPPool creates new DHCPPool.
func NewDHCPPool(raw RawDHCPPool) *DHCPPool {
	pool := ParseRange(raw.Range)
	if pool == nil {
		return nil
	}

	name := raw.Name
	if name == "" {
		name = raw.Range
	}

	return &DHCPPool{
		IRange: pool,
		name:   name,
		leases: make(map[string]bool),
	}
}

type DHCPPool struct {
	IRange
	name   string
	leases map[string]bool
}

// Name returns name.
func (d DHCPPool) Name() string { return d.name }

// NumOfLeases returns number of active dhcp leases.
func (d DHCPPool) NumOfLeases() int64 { return int64(len(d.leases)) }

// Lease adds IP to leases database.
func (d *DHCPPool) Lease(ip net.IP) { d.leases[ip.String()] = true }

// ResetLeases resets all leases.
func (d *DHCPPool) ResetLeases() { d.leases = make(map[string]bool) }

// Utilization returns pool utilization in percent.
func (d DHCPPool) Utilization() float64 {
	total := d.Hosts()
	if !total.IsInt64() {
		return 0
	}
	return float64(d.NumOfLeases()) * 100 / float64(total.Int64())
}

func ParseRange(iprange string) IRange {
	if strings.Contains(iprange, "-") {
		return NewRange(iprange)
	}
	return NewNet(iprange)
}

package ip

import (
	"fmt"
	"math/big"
	"net"
	"strings"
)

const (
	defaultV4Mask = "32"
	defaultV6Mask = "64"
)

// Net represents IP Network.
type Net struct {
	*net.IPNet
}

// Type returns Net IP type.
func (n Net) Type() Type {
	if n.IP.To4() != nil {
		return V4Type
	}
	if n.IP.To16() != nil {
		return V6Type
	}
	return UnknownType
}

// Size returns number of available IPs in the Net.
func (n Net) Size() *big.Int {
	var (
		ones, bits = n.IPNet.Mask.Size()
		zero       = big.NewInt(0)
		size       = big.NewInt(1)
		two        = big.NewInt(2)
	)
	if ones == 0 && bits == 0 {
		return zero
	}
	for i := 0; i < bits-ones; i++ {
		size.Mul(size, two)
	}
	if size.Sub(size, two).Cmp(zero) <= 0 {
		return zero
	}
	return size
}

func ParseNet(s string) *Net {
	_, ipnet, err := net.ParseCIDR(addMaskToIP(s))
	if err != nil || ipnet == nil {
		return nil
	}

	return &Net{IPNet: ipnet}
}

func addMaskToIP(s string) string {
	if strings.Contains(s, "/") {
		return s
	}
	var (
		ip   net.IP
		mask = defaultV4Mask
	)
	if ip = net.ParseIP(s); ip == nil {
		return ""
	}
	if ip.To4() == nil {
		mask = defaultV6Mask
	}

	return fmt.Sprintf("%s/%s", ip, mask)
}

package ip

import (
	"fmt"
	"math/big"
	"net"
	"strings"
)

type Net struct {
	*net.IPNet
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
	s = strings.TrimSpace(s)

	if !strings.Contains(s, "/") {
		ip := net.ParseIP(s)
		if ip == nil {
			return nil
		}
		s = fmt.Sprintf("%s/%s", ip, ip.DefaultMask())
	}

	parts := strings.Split(s, "/")
	if len(parts) != 2 {
		return nil
	}

	_, ipnet, err := net.ParseCIDR(s)
	if err != nil || ipnet == nil {
		return nil
	}

	return &Net{IPNet: ipnet}
}

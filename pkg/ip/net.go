package ip

//import (
//	"fmt"
//	"math/big"
//	"net"
//	"strings"
//)
//
//const (
//	defaultV4Mask = "32"
//	defaultV6Mask = "64"
//)
//
//func NewNet(s string) *Net {
//	if s == "" {
//		return nil
//	}
//	_, ipnet, err := net.ParseCIDR(addMaskToIP(s))
//	if err != nil || ipnet == nil {
//		return nil
//	}
//	return &Net{IPNet: ipnet}
//}
//
//// Net represents IP Network.
//type Net struct {
//	*net.IPNet
//}
//
//// Type returns Net IP type.
//func (n Net) Type() Type {
//	if n.IP.To4() != nil {
//		return V4Type
//	}
//	if n.IP.To16() != nil {
//		return V6Type
//	}
//	return InvalidType
//}
//
//// Hosts returns number of hosts addresses in the Net.
//func (n Net) Hosts() *big.Int {
//	var (
//		ones, bits = n.Mask.Size()
//		zero       = big.NewInt(0)
//		hosts      = big.NewInt(1)
//		two        = big.NewInt(2)
//	)
//	if ones == 0 && bits == 0 {
//		return zero
//	}
//	for i := 0; i < bits-ones; i++ {
//		hosts.Mul(hosts, two)
//	}
//
//	switch n.Type() {
//	default:
//		return zero
//	case V4Type:
//		if hosts.Sub(hosts, two).Cmp(zero) < 0 {
//			return zero
//		}
//		return hosts
//	case V6Type:
//		return hosts
//	}
//}
//
//func addMaskToIP(s string) string {
//	if strings.Contains(s, "/") {
//		return s
//	}
//	var (
//		ip   net.IP
//		mask = defaultV4Mask
//	)
//	if ip = net.ParseIP(s); ip == nil {
//		return ""
//	}
//	if ip.To4() == nil {
//		mask = defaultV6Mask
//	}
//
//	return fmt.Sprintf("%s/%s", ip, mask)
//}

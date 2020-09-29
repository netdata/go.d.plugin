package iprange

import (
	"bytes"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/apparentlymart/go-cidr/cidr"
)

func ParseRanges(value string) ([]Range, error) {
	parts := strings.Fields(value)
	if len(parts) == 0 {
		return nil, nil
	}

	var ranges []Range
	for _, v := range parts {
		r, err := ParseRange(v)
		if err != nil {
			return nil, err
		}

		if r != nil {
			ranges = append(ranges, r)
		}
	}
	return ranges, nil
}

var (
	reRange = regexp.MustCompile("^[0-9a-fA-F.:-]+$")            // addr | addr-addr
	reCIDR  = regexp.MustCompile("^[0-9a-fA-F.:-]+/[0-9]{1,3}$") // addr/prefix_length
	reV4Net = regexp.MustCompile("^[0-9.-]+/[0-9.-]{8,}$")       // v4_addr/mask
)

func ParseRange(s string) (Range, error) {
	if s == "" {
		return nil, nil
	}

	var r Range
	switch {
	case reRange.MatchString(s):
		r = parseRange(s)
	case reCIDR.MatchString(s):
		r = parseCIDR(s)
	case reV4Net.MatchString(s):
		r = parseV4Network(s)
	}
	if r == nil {
		return nil, fmt.Errorf("ip range (%s) invalid syntax", s)
	}
	return r, nil
}

func parseRange(s string) Range {
	var start, end net.IP
	switch parts := strings.Split(s, "-"); len(parts) {
	case 1:
		start, end = net.ParseIP(parts[0]), net.ParseIP(parts[0])
	case 2:
		start, end = net.ParseIP(parts[0]), net.ParseIP(parts[1])
	default:
		return nil
	}

	switch {
	case start.To4() != nil && end.To4() != nil && bytes.Compare(end, start) >= 0:
		return v4Range{start: start, end: end}
	case start.To16() != nil && end.To16() != nil && bytes.Compare(end, start) >= 0:
		return v6Range{start: start, end: end}
	default:
		return nil
	}
}

func parseCIDR(s string) Range {
	ip, network, err := net.ParseCIDR(s)
	if err != nil {
		return nil
	}

	start, end := cidr.AddressRange(network)
	ones, _ := network.Mask.Size()

	if isV4 := ip.To4() != nil; isV4 && ones < 31 || ones < 127 {
		start = cidr.Inc(start)
		end = cidr.Dec(end)
	}

	return parseRange(fmt.Sprintf("%s-%s", start, end))
}

func parseV4Network(s string) Range {
	idx := strings.LastIndexByte(s, '/')
	if idx == -1 {
		return nil
	}

	address, mask := s[:idx], s[idx+1:]

	ip := net.ParseIP(mask).To4()
	if ip == nil {
		return nil
	}

	ones, _ := net.IPv4Mask(ip[0], ip[1], ip[2], ip[3]).Size()
	if ones == 0 {
		return nil
	}

	return parseCIDR(fmt.Sprintf("%s/%s", address, strconv.Itoa(ones)))
}

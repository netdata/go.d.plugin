package dnsmasq_dhcp

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"sort"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/iprange"
)

func (d *DnsmasqDHCP) autodetection() error {
	configs := findConfigurationFiles(d.ConfPath, d.ConfDir)

	err := d.autodetectDHCPRanges(configs)
	if err != nil {
		return err
	}

	d.autodetectStaticIPs(configs)
	return nil
}

func (d *DnsmasqDHCP) autodetectDHCPRanges(configs []*configFile) error {
	var ipranges []iprange.Range
	var parsed string
	seen := make(map[string]bool)

	for _, conf := range configs {
		d.Debugf("looking in '%s'", conf.path)

		for _, value := range conf.get("dhcp-range") {
			d.Debugf("found dhcp-range '%s'", value)
			if parsed = parseDHCPRangeValue(value); parsed == "" || seen[parsed] {
				continue
			}
			seen[parsed] = true

			r, err := iprange.ParseRange(parsed)
			if r == nil || err != nil {
				d.Warningf("error on parsing dhcp-range '%s', skipping it", parsed)
				continue
			}

			d.Debugf("adding dhcp-range '%s'", parsed)
			ipranges = append(ipranges, r)
		}
	}

	if len(ipranges) == 0 {
		return errors.New("haven't found any dhcp-ranges")
	}

	// order: ipv4, ipv6
	sort.Slice(
		ipranges,
		func(i, j int) bool { return ipranges[i].Family() < ipranges[j].Family() },
	)

	d.ranges = ipranges
	return nil
}

func (d *DnsmasqDHCP) autodetectStaticIPs(configs []*configFile) {
	seen := make(map[string]bool)
	var parsed string

	for _, conf := range configs {
		d.Debugf("looking in '%s'", conf.path)

		for _, value := range conf.get("dhcp-host") {
			d.Debugf("found dhcp-host '%s'", value)
			if parsed = parseDHCPHostValue(value); parsed == "" || seen[parsed] {
				continue
			}
			seen[parsed] = true

			v := net.ParseIP(parsed)
			if v == nil {
				d.Warningf("error on parsing dhcp-host '%s', skipping it", parsed)
				continue
			}

			d.Debugf("adding dhcp-host '%s'", parsed)
			d.staticIPs = append(d.staticIPs, v)
		}
	}
}

/*
Examples:
  - 192.168.0.50,192.168.0.150,12h
  - 192.168.0.50,192.168.0.150,255.255.255.0,12h
  - set:red,1.1.1.50,1.1.2.150, 255.255.252.0
  - 192.168.0.0,static
  - 1234::2,1234::500, 64, 12h
  - 1234::2,1234::500
  - 1234::2,1234::500, slaac
  - 1234::,ra-only
  - 1234::,ra-names
  - 1234::,ra-stateless
*/
var reDHCPRange = regexp.MustCompile(`([0-9a-f.:]+),([0-9a-f.:]+)`)

func parseDHCPRangeValue(s string) (r string) {
	if strings.Contains(s, "ra-stateless") {
		return
	}

	match := reDHCPRange.FindStringSubmatch(s)
	if match == nil {
		return
	}

	start, end := net.ParseIP(match[1]), net.ParseIP(match[2])
	if start == nil || end == nil {
		return
	}

	return fmt.Sprintf("%s-%s", start, end)
}

/*
Examples:
  - 11:22:33:44:55:66,192.168.0.60
  - 11:22:33:44:55:66,fred,192.168.0.60,45m
  - 11:22:33:44:55:66,12:34:56:78:90:12,192.168.0.60
  - bert,192.168.0.70,infinite
  - id:01:02:02:04,192.168.0.60
  - id:ff:00:00:00:00:00:02:00:00:02:c9:00:f4:52:14:03:00:28:05:81,192.168.0.61
  - id:marjorie,192.168.0.60
  - id:00:01:00:01:16:d2:83:fc:92:d4:19:e2:d8:b2, fred, [1234::5]
*/
var (
	reDHCPHostV4 = regexp.MustCompile(`(?:[0-9]{1,3}\.){3}[0-9]{1,3}`)
	reDHCPHostV6 = regexp.MustCompile(`\[([0-9a-f.:]+)\]`)
)

func parseDHCPHostValue(s string) (r string) {
	if strings.Contains(s, "[") {
		return strings.Trim(reDHCPHostV6.FindString(s), "[]")
	}
	return reDHCPHostV4.FindString(s)
}

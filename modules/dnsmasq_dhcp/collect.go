package dnsmasq_dhcp

import (
	"bufio"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"strings"
)

func (d *DnsmasqDHCP) collect() (map[string]int64, error) {
	f, err := os.Open(d.LeasesPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if d.modTime.Equal(fi.ModTime()) {
		for k, v := range d.mx {
			fmt.Println(k, v)
		}
		return d.mx, nil
	}
	d.modTime = fi.ModTime()

	mx := make(map[string]int64)

	for _, ip := range findIPs(f) {
		for _, r := range d.ranges {
			if !r.Contains(ip) {
				continue
			}
			mx[r.String()]++
			break
		}
	}

	for _, r := range d.ranges {
		name := r.String()
		numOfIps, ok := mx[name]
		if !ok {
			mx[name] = 0
		}

		hosts := r.Hosts()
		if !hosts.IsInt64() {
			continue
		}

		mx[name+"_utilization"] = int64(calcPercent(numOfIps, hosts) * 1000)
	}

	d.mx = mx

	return mx, nil
}

func findIPs(r io.Reader) []net.IP {
	var ips []net.IP
	s := bufio.NewScanner(r)

	for s.Scan() {
		parts := strings.Fields(s.Text())
		if len(parts) != 5 {
			continue
		}

		ip := net.ParseIP(parts[2])
		if ip == nil {
			continue
		}
		ips = append(ips, ip)
	}

	return ips
}

func calcPercent(ips int64, hosts *big.Int) float64 {
	return float64(ips) * 100 / float64(hosts.Int64())
}

package dnsmasq_dhcp

import (
	"bufio"
	"io"
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
		return d.mx, nil
	}
	d.modTime = fi.ModTime()

	mx := make(map[string]int64)

	for _, lease := range findLeases(f) {
		for _, r := range d.ranges {
			if !r.Contains(lease) {
				continue
			}
			mx[r.String()]++
			break
		}
	}

	for _, r := range d.ranges {
		name := r.String()
		v, ok := mx[name]
		if !ok {
			mx[name] = 0
		}

		mx[name+"_utilization"] = 0
		h := r.Hosts()
		if !h.IsInt64() {
			continue
		}

		mx[name+"_utilization"] = int64(float64(v) * 100 / float64(h.Int64()) * 1000)
	}

	d.mx = mx

	return mx, nil
}

func findLeases(r io.Reader) []net.IP {
	var leases []net.IP
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
		leases = append(leases, ip)
	}

	return leases
}

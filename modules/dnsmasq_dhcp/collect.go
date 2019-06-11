package dnsmasq_dhcp

import (
	"bufio"
	"fmt"
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

	_ = fi

	//if d.modTime.Equal(fi.ModTime()) {
	//	return nil, nil
	//}

	s := bufio.NewScanner(f)
	for s.Scan() {
		d.collectLine(s.Text())
	}

	for _, pool := range d.pools {
		fmt.Println(pool, pool.Hosts(), pool.NumOfLeases(), pool.Utilization())
	}

	return nil, nil
}

func (d *DnsmasqDHCP) collectLine(line string) {
	// 1560248031 08:00:27:61:3c:ee 1.1.2.27 debian8 *
	// 1560252212 660684014 1234::20b * 00:01:00:01:24:90:cf:a3:08:00:27:61:3c:ee

	parts := strings.Fields(line)
	if len(parts) != 5 {
		return
	}

	ip := net.ParseIP(parts[2])
	if ip == nil {
		return
	}

	for _, pool := range d.pools {
		if pool.Contains(ip) {
			pool.Lease(ip)
		}
	}
}

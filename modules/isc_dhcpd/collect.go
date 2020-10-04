package isc_dhcpd

import (
	"math/big"
	"os"
)

/*
dhcpd.leases db (file), see details: https://kb.isc.org/docs/en/isc-dhcp-44-manual-pages-dhcpdleases#dhcpdleases

Every time a prepare is acquired, renewed or released, its new value is recorded at the end of the prepare file.
So if more than one declaration appears for a given prepare, the last one in the file is the current one.

In order to prevent the prepare database from growing without bound, the file is rewritten from time to time.
First, a temporary lease database is created and all known leases are dumped to it.
Then, the old lease database is renamed DBDIR/dhcpd.leases~.
Finally, the newly written lease database is moved into place.

In order to process both DHCPv4 and DHCPv6 messages you will need to run two separate instances of the dhcpd process.
Each of these instances will need it’s own prepare file.
*/

func (d *DHCPd) collect() (map[string]int64, error) {
	fi, err := os.Stat(d.LeasesPath)
	if err != nil {
		return nil, err
	}

	if d.leasesModTime.Equal(fi.ModTime()) {
		d.Debugf("leases file is not modified, returning cached metrics ('%s')", d.LeasesPath)
		return d.collected, nil
	}

	d.leasesModTime = fi.ModTime()

	leases, err := parseDHCPdLeasesFile(d.LeasesPath)
	if err != nil {
		return nil, err
	}

	activeLeases := removeInactiveLeases(leases)
	d.Debugf("found total/active %d/%d leases ('%s')", len(leases), len(activeLeases), d.LeasesPath)

	for _, pool := range d.pools {
		collectPool(d.collected, pool, activeLeases)
	}
	d.collected["active_leases_total"] = int64(len(activeLeases))

	return d.collected, nil
}

func collectPool(collected map[string]int64, pool ipPool, leases []leaseEntry) {
	var n int64
	for _, l := range leases {
		if pool.addresses.Contains(l.ip) {
			n++
		}
	}
	collected["pool_"+pool.name+"_active_leases"] = n
	collected["pool_"+pool.name+"_utilization"] = calcPoolUtilizationPercentage(pool.addresses.Size(), n)
}

const precision = 100

func calcPoolUtilizationPercentage(size *big.Int, leases int64) int64 {
	if leases == 0 || !size.IsInt64() {
		return 0
	}
	if size.Int64() == 0 {
		return 100 * precision
	}
	return int64(float64(leases) / float64(size.Int64()) * 100 * precision)
}

func removeInactiveLeases(leases []leaseEntry) (active []leaseEntry) {
	active = leases[:0]
	for _, l := range leases {
		if l.bindingState == "active" {
			active = append(active, l)
		}
	}
	return active
}

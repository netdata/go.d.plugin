package isc_dhcpd

import "github.com/netdata/go-orchestrator/module"

var (
	dhcpdCharts = module.Charts {
		dhcpPollsUtilization.Copy(),
		dhcpPollsActiveLeases.Copy(),
	}

	dhcpPollsUtilization = module.Chart {
		ID : "pools_utilization",
		Title : "Pools Utilization",
		Units : "percentage",
		Fam : "utilization",
		Ctx : "isc_dhcpd.utilization",
	}

	dhcpPollsActiveLeases = module.Chart {
		ID : "pools_active_leases",
		Title : "Active Leases Per Pool",
		Units : "leases",
		Fam : "active leases",
		Ctx : "isc_dhcpd.active_leases",
	}

	dhcpPollsTotalLeases = module.Chart {
		ID : "pools_total_leases",
		Title : "All Active Leases",
		Units : "leases",
		Fam : "active leases",
		Ctx : "isc_dhcpd.leases_total",
	}
)
package vsphere

import (
	"fmt"
	"strings"

	rs "github.com/netdata/go.d.plugin/modules/vsphere/resources"

	"github.com/netdata/go-orchestrator"
	"github.com/netdata/go-orchestrator/module"
)

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Chart is an alias for module.Chart
	Chart = module.Chart
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

const (
	hostPrio = orchestrator.DefaultJobPriority
	vmPrio   = hostPrio + 200
)

var (
	vmCharts = func() Charts {
		cs := Charts{}
		panicIf(cs.Add(vmCPUCharts...))
		panicIf(cs.Add(vmMemCharts...))
		panicIf(cs.Add(vmNetCharts...))
		panicIf(cs.Add(vmDiskCharts...))
		panicIf(cs.Add(vmSystemCharts...))
		return cs
	}()

	vmCPUCharts = Charts{
		{
			ID:    "vm_%s_%s_cpu_usage_total",
			Title: "Cpu Usage Total",
			Units: "percentage",
			Fam:   "vm %s",
			Ctx:   "vsphere.cpu_usage_total",
			Dims: Dims{
				{ID: "%s_%s_cpu.usage.average", Name: "used", Div: 100},
			},
		},
	}
	// Ref: https://www.vmware.com/support/developer/converter-sdk/conv51_apireference/memory_counters.html
	vmMemCharts = Charts{
		{
			ID:    "vm_%s_%s_mem_usage_percentage",
			Title: "Memory Usage Percentage",
			Units: "percentage",
			Fam:   "vm %s",
			Ctx:   "vsphere.mem_usage_percentage",
			Dims: Dims{
				{ID: "%s_%s_mem.usage.average", Name: "used", Div: 100},
			},
		},
		{
			ID:    "vm_%s_%s_mem_usage",
			Title: "Memory Usage",
			Units: "KB",
			Fam:   "vm %s",
			Ctx:   "vsphere.mem_usage",
			Dims: Dims{
				{ID: "%s_%s_mem.granted.average", Name: "granted"},
				{ID: "%s_%s_mem.consumed.average", Name: "consumed"},
				{ID: "%s_%s_mem.active.average", Name: "active"},
				{ID: "%s_%s_mem.shared.average", Name: "shared"},
			},
		},
		{
			ID:    "vm_%s_%s_mem_swap_rate",
			Title: "Memory Swap Rate",
			Units: "KB/s",
			Fam:   "vm %s",
			Ctx:   "vsphere.mem_swap_rate",
			Dims: Dims{
				{ID: "%s_%s_mem.swapinRate.average", Name: "in"},
				{ID: "%s_%s_mem.swapoutRate.average", Name: "out"},
			},
		},
		{
			ID:    "vm_%s_%s_mem_swap",
			Title: "Memory Swap",
			Units: "KB",
			Fam:   "vm %s",
			Ctx:   "vsphere.mem_swap",
			Dims: Dims{
				{ID: "%s_%s_mem.swapped.average", Name: "swapped"},
			},
		},
	}
	vmNetCharts = Charts{
		{
			ID:    "vm_%s_%s_net_bandwidth_total",
			Title: "Network Bandwidth Total",
			Units: "KB/s",
			Fam:   "vm %s",
			Ctx:   "vsphere.net_bandwidth_total",
			Dims: Dims{
				{ID: "%s_%s_net.bytesRx.average", Name: "rx"},
				{ID: "%s_%s_net.bytesTx.average", Name: "tx", Mul: -1},
			},
		},
		{
			ID:    "vm_%s_%s_net_packets_total",
			Title: "Network Packets Total",
			Units: "packets",
			Fam:   "vm %s",
			Ctx:   "vsphere.net_packets_total",
			Dims: Dims{
				{ID: "%s_%s_net.packetsRx.summation", Name: "rx"},
				{ID: "%s_%s_net.packetsTx.summation", Name: "tx", Mul: -1},
			},
		},
		{
			ID:    "vm_%s_%s_net_drops_total",
			Title: "Network Drops Total",
			Units: "dropped packets",
			Fam:   "vm %s",
			Ctx:   "vsphere.net_drops_total",
			Dims: Dims{
				{ID: "%s_%s_net.droppedRx.summation", Name: "rx"},
				{ID: "%s_%s_net.droppedTx.summation", Name: "tx", Mul: -1},
			},
		},
	}
	vmDiskCharts = Charts{
		{
			ID:    "vm_%s_%s_disk_usage_total",
			Title: "Disk Usage Total",
			Units: "KB/s",
			Fam:   "vm %s",
			Ctx:   "vsphere.disk_usage_total",
			Dims: Dims{
				{ID: "%s_%s_disk.read.average", Name: "read"},
				{ID: "%s_%s_disk.write.average", Name: "write", Mul: -1},
			},
		},
		{
			ID:    "vm_%s_%s_disk_max_latency",
			Title: "Disk Max Latency",
			Units: "ms",
			Fam:   "vm %s",
			Ctx:   "vsphere.disk_max_latency",
			Dims: Dims{
				{ID: "%s_%s_disk.maxTotalLatency.latest", Name: "latency"},
			},
		},
	}
	vmSystemCharts = Charts{
		{
			ID:    "vm_%s_%s_system_uptime",
			Title: "System Uptime",
			Units: "seconds",
			Fam:   "vm %s",
			Ctx:   "vsphere.system_uptime",
			Dims: Dims{
				{ID: "%s_%s_sys.uptime.latest", Name: "time"},
			},
		},
	}
)

var (
	hostCharts = func() Charts {
		cs := Charts{}
		panicIf(cs.Add(hostCPUCharts...))
		panicIf(cs.Add(hostMemCharts...))
		panicIf(cs.Add(hostNetCharts...))
		panicIf(cs.Add(hostDiskCharts...))
		panicIf(cs.Add(hostSystemCharts...))
		return cs
	}()
	hostCPUCharts = Charts{
		{
			ID:    "host_%s_%s_cpu_usage_total",
			Title: "Cpu Usage Total",
			Units: "percentage",
			Fam:   "host %s",
			Ctx:   "vsphere.cpu_usage_total",
			Dims: Dims{
				{ID: "%s_%s_cpu.usage.average", Name: "used", Div: 100},
			},
		},
	}
	// Ref: https://www.vmware.com/support/developer/converter-sdk/conv51_apireference/memory_counters.html
	hostMemCharts = Charts{
		{
			ID:    "host_%s_%s_mem_usage_percentage",
			Title: "Memory Usage Percentage",
			Units: "percentage",
			Fam:   "host %s",
			Ctx:   "vsphere.mem_usage_percentage",
			Dims: Dims{
				{ID: "%s_%s_mem.usage.average", Name: "used", Div: 100},
			},
		},
		{
			ID:    "host_%s_%s_mem_usage",
			Title: "Memory Usage",
			Units: "KB",
			Fam:   "host %s",
			Ctx:   "vsphere.mem_usage",
			Dims: Dims{
				{ID: "%s_%s_mem.granted.average", Name: "granted"},
				{ID: "%s_%s_mem.consumed.average", Name: "consumed"},
				{ID: "%s_%s_mem.active.average", Name: "active"},
				{ID: "%s_%s_mem.shared.average", Name: "shared"},
				{ID: "%s_%s_mem.sharedcommon.average", Name: "sharedcommon"},
			},
		},
		{
			ID:    "host_%s_%s_mem_swap_rate",
			Title: "Memory Swap Rate",
			Units: "KB/s",
			Fam:   "host %s",
			Ctx:   "vsphere.mem_swap_rate",
			Dims: Dims{
				{ID: "%s_%s_mem.swapinRate.average", Name: "in"},
				{ID: "%s_%s_mem.swapoutRate.average", Name: "out"},
			},
		},
	}
	hostNetCharts = Charts{
		{
			ID:    "host_%s_%s_net_bandwidth_total",
			Title: "Network Bandwidth Total",
			Units: "KB/s",
			Fam:   "host %s",
			Ctx:   "vsphere.net_bandwidth_total",
			Dims: Dims{
				{ID: "%s_%s_net.bytesRx.average", Name: "rx"},
				{ID: "%s_%s_net.bytesTx.average", Name: "tx", Mul: -1},
			},
		},
		{
			ID:    "host_%s_%s_net_packets_total",
			Title: "Network Packets Total",
			Units: "packets",
			Fam:   "host %s",
			Ctx:   "vsphere.net_packets_total",
			Dims: Dims{
				{ID: "%s_%s_net.packetsRx.summation", Name: "rx"},
				{ID: "%s_%s_net.packetsTx.summation", Name: "tx", Mul: -1},
			},
		},
		{
			ID:    "host_%s_%s_net_drops_total",
			Title: "Network Drops Total",
			Units: "dropped packets",
			Fam:   "host %s",
			Ctx:   "vsphere.net_drops_total",
			Dims: Dims{
				{ID: "%s_%s_net.droppedRx.summation", Name: "rx"},
				{ID: "%s_%s_net.droppedTx.summation", Name: "tx", Mul: -1},
			},
		},
		{
			ID:    "host_%s_%s_net_errors_total",
			Title: "Network Errors Total",
			Units: "errors",
			Fam:   "host %s",
			Ctx:   "vsphere.net_errors_total",
			Dims: Dims{
				{ID: "%s_%s_net.errorsRx.summation", Name: "rx"},
				{ID: "%s_%s_net.errorsTx.summation", Name: "tx", Mul: -1},
			},
		},
	}
	hostDiskCharts = Charts{
		{
			ID:    "host_%s_%s_disk_usage_total",
			Title: "Disk Usage Total",
			Units: "KB/s",
			Fam:   "host %s",
			Ctx:   "vsphere.disk_usage_total",
			Dims: Dims{
				{ID: "%s_%s_disk.read.average", Name: "read"},
				{ID: "%s_%s_disk.write.average", Name: "write", Mul: -1},
			},
		},
		{
			ID:    "host_%s_%s_disk_max_latency",
			Title: "Disk Max Latency",
			Units: "ms",
			Fam:   "host %s",
			Ctx:   "vsphere.disk_max_latency",
			Dims: Dims{
				{ID: "%s_%s_disk.maxTotalLatency.latest", Name: "latency"},
			},
		},
	}
	hostSystemCharts = Charts{
		{
			ID:    "host_%s_%s_system_uptime",
			Title: "System Uptime",
			Units: "seconds",
			Fam:   "host %s",
			Ctx:   "vsphere.system_uptime",
			Dims: Dims{
				{ID: "%s_%s_sys.uptime.latest", Name: "time"},
			},
		},
	}
)

func (vs *VSphere) updateCharts() {
	vs.updateHostsCharts()
	vs.updateVMsCharts()
}

func (vs *VSphere) updateHostsCharts() {
	for _, h := range vs.resources.Hosts {
		if vs.charted[h.ID] {
			continue
		}

		vs.charted[h.ID] = true
		cs := newHostCharts(h)
		panicIf(vs.charts.Add(*cs...))
	}
}

func newHostCharts(host *rs.Host) *Charts {
	cs := hostCharts.Copy()
	for i, c := range *cs {
		setHostChart(c, host, hostPrio+i)
	}
	return cs
}

func setHostChart(chart *Chart, host *rs.Host, prio int) {
	chart.Priority = prio
	chart.ID = fmt.Sprintf(chart.ID, host.ID, host.Name)
	chart.Fam = fmt.Sprintf(chart.Fam, host.Name)
	for _, d := range chart.Dims {
		d.ID = fmt.Sprintf(d.ID, host.ID, host.Name)
	}
}

func (vs *VSphere) updateVMsCharts() {
	for _, v := range vs.resources.VMs {
		if vs.charted[v.ID] {
			continue
		}

		vs.charted[v.ID] = true
		cs := newVMCHarts(v)
		panicIf(vs.charts.Add(*cs...))
	}
}

func newVMCHarts(vm *rs.VM) *Charts {
	cs := vmCharts.Copy()
	for i, c := range *cs {
		setVMChart(c, vm, vmPrio+i)
	}
	return cs
}

func setVMChart(chart *Chart, vm *rs.VM, prio int) {
	chart.Priority = prio
	chart.ID = fmt.Sprintf(chart.ID, vm.ID, vm.Name)
	chart.Fam = fmt.Sprintf(chart.Fam, vm.Name)
	for _, d := range chart.Dims {
		d.ID = fmt.Sprintf(d.ID, vm.ID, vm.Name)
	}
}

func (vs *VSphere) removeFromCharts(prefix string) {
	for _, c := range *vs.charts {
		if strings.HasPrefix(c.ID, prefix) {
			c.MarkRemove()
			c.MarkNotCreated()
		}
	}
}

//func findMetricSeriesByPrefix(ms []performance.MetricSeries, prefix string) []performance.MetricSeries {
//	from := sort.Search(len(ms), func(i int) bool { return ms[i].Name >= prefix })
//
//	if from == len(ms) || !strings.HasPrefix(ms[from].Name, prefix) {
//		return nil
//	}
//
//	until := from + 1
//	for until < len(ms) && strings.HasPrefix(ms[until].Name, prefix) {
//		until++
//	}
//	return ms[from:until]
//}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

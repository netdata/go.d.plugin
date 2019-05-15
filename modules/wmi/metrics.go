package wmi

import (
	mtx "github.com/netdata/go.d.plugin/pkg/metrics"
	"sort"
	"strconv"
)

func newMetrics() *metrics { return &metrics{} }

type metrics struct {
	CPU *cpu     `stm:"cpu"`
	Net *network `stm:"net"`
}

type (
	cpu struct {
		// Mode represents time that processor spent in different modes.
		PercentIdleTime       mtx.Gauge `stm:"idle,1000,1"`
		PercentInterruptTime  mtx.Gauge `stm:"interrupt,1000,1"`
		PercentDPCTime        mtx.Gauge `stm:"dpc,1000,1"`
		PercentPrivilegedTime mtx.Gauge `stm:"privileged,1000,1"`
		PercentUserTime       mtx.Gauge `stm:"user,1000,1"`
		// Cores is the per core statistics.
		Cores cpuCores `stm:"core"`
	}

	cpuCores []*cpuCore

	cpuCore struct {
		STMKey string
		ID     string
		id     int
		// Total number of received and serviced deferred procedure calls (DPCs).
		DPCsQueuedPerSec mtx.Gauge `stm:"dpc,1000,1"`
		// Total number of received and serviced hardware interrupts.
		InterruptsPerSec mtx.Gauge `stm:"interrupts,1000,1"`
		// Mode represents time that processor spent in different modes.
		PercentIdleTime       mtx.Gauge `stm:"idle,1000,1"`
		PercentInterruptTime  mtx.Gauge `stm:"interrupt,1000,1"`
		PercentDPCTime        mtx.Gauge `stm:"dpc,1000,1"`
		PercentPrivilegedTime mtx.Gauge `stm:"privileged,1000,1"`
		PercentUserTime       mtx.Gauge `stm:"user,1000,1"`
		// CState represents time spent in low-power idle state.
		PercentC1Time mtx.Gauge `stm:"c1,1000,1"`
		PercentC2Time mtx.Gauge `stm:"c2,1000,1"`
		PercentC3Time mtx.Gauge `stm:"c3,1000,1"`
	}
)

type (
	network struct {
		NICs netNICs `stm:""`
	}

	netNICs []*netNIC

	netNIC struct {
		STMKey                   string
		ID                       string
		BytesReceivedTotal       mtx.Gauge `stm:"bytes_received,1000,1"`
		BytesSentTotal           mtx.Gauge `stm:"bytes_sent,1000,1"`
		BytesTotal               mtx.Gauge `stm:"bytes_total,1000,1"`
		PacketsOutboundDiscarded mtx.Gauge `stm:"packets_outbound_discarded,1000,1"`
		PacketsOutboundErrors    mtx.Gauge `stm:"packets_outbound_errors,1000,1"`
		PacketsTotal             mtx.Gauge `stm:"packets_total,1000,1"`
		PacketsReceivedDiscarded mtx.Gauge `stm:"packets_received_discarded,1000,1"`
		PacketsReceivedErrors    mtx.Gauge `stm:"packets_received_errors,1000,1"`
		PacketsReceivedTotal     mtx.Gauge `stm:"packets_received_total,1000,1"`
		PacketsReceivedUnknown   mtx.Gauge `stm:"packets_received_unknown,1000,1"`
		PacketsSentTotal         mtx.Gauge `stm:"packets_sent_total,1000,1"`
		CurrentBandwidth         mtx.Gauge `stm:"current_bandwidth"`
	}
)

func newCPUCore(id string) *cpuCore { return &cpuCore{STMKey: id, ID: id, id: getCPUIntID(id)} }

func (cc *cpuCores) get(id string, createIfNotExist bool) (core *cpuCore) {
	for _, c := range *cc {
		if c.ID == id {
			return c
		}
	}
	if createIfNotExist {
		core = newCPUCore(id)
		*cc = append(*cc, core)
	}
	return core
}

func (cc *cpuCores) sort() { sort.Slice(*cc, func(i, j int) bool { return (*cc)[i].id < (*cc)[j].id }) }

func newNIC(id string) *netNIC { return &netNIC{STMKey: id, ID: id} }

func (ns *netNICs) get(id string, createIfNotExist bool) (nic *netNIC) {
	for _, n := range *ns {
		if n.ID == id {
			return n
		}
	}
	if createIfNotExist {
		nic = newNIC(id)
		*ns = append(*ns, nic)
	}
	return nic
}

func getCPUIntID(id string) int {
	if id == "" {
		return -1
	}
	v, err := strconv.Atoi(string(id[len(id)-1]))
	if err != nil {
		return -1
	}
	return v
}

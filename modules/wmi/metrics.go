package wmi

import (
	mtx "github.com/netdata/go.d.plugin/pkg/metrics"
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
		PercentC1Time mtx.Gauge `stm:"cstate_c1,1000,1"`
		PercentC2Time mtx.Gauge `stm:"cstate_c2,1000,1"`
		PercentC3Time mtx.Gauge `stm:"cstate_c3,1000,1"`
	}
)

type (
	network struct {
		NICs nics
	}

	nics []*nic

	nic struct {
		STMKey                   string
		ID                       string
		BytesReceivedTotal       mtx.Gauge `stm:"bytes_received"`
		BytesSentTotal           mtx.Gauge `stm:"bytes_sent"`
		BytesTotal               mtx.Gauge `stm:"bytes_total"`
		PacketsOutboundDiscarded mtx.Gauge `stm:"packets_outbound_discarded"`
		PacketsOutboundErrors    mtx.Gauge `stm:"packets_outbound_errors"`
		PacketsTotal             mtx.Gauge `stm:"packets_total"`
		PacketsReceivedDiscarded mtx.Gauge `stm:"packets_received_discarded"`
		PacketsReceivedErrors    mtx.Gauge `stm:"packets_received_errors"`
		PacketsReceivedTotal     mtx.Gauge `stm:"packets_received_total"`
		PacketsReceivedUnknown   mtx.Gauge `stm:"packets_received_unknown"`
		PacketsSentTotal         mtx.Gauge `stm:"packets_sent_total"`
		CurrentBandwidth         mtx.Gauge `stm:"current_bandwidth"`
	}
)

func newCPUCore(id string) *cpuCore { return &cpuCore{STMKey: id, ID: id} }

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

func newNIC(id string) *nic { return &nic{STMKey: id, ID: id} }

func (ns *nics) get(id string, createIfNotExist bool) (n *nic) {
	for _, n := range *ns {
		if n.ID == id {
			return n
		}
	}
	if createIfNotExist {
		n = newNIC(id)
		*ns = append(*ns, n)
	}
	return n
}

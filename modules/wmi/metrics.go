package wmi

import mtx "github.com/netdata/go.d.plugin/pkg/metrics"

func newMetrics() *metrics { return &metrics{} }

type metrics struct {
	CPU struct {
		// Time that processor spent in different modes.
		Time  cpuModes `stm:"time"`
		Cores cpuCores `stm:"core"`
	} `stm:"cpu"`
}

type (
	cpuModes struct {
		DPC        mtx.Gauge `stm:"dpc,1000,1"`
		Idle       mtx.Gauge `stm:"idle,1000,1"`
		Interrupt  mtx.Gauge `stm:"interrupt,1000,1"`
		Privileged mtx.Gauge `stm:"privileged,1000,1"`
		User       mtx.Gauge `stm:"user,1000,1"`
	}

	cpuCore struct {
		STMKey string
		// Core id
		ID string
		// Total number of received and serviced deferred procedure calls (DPCs).
		DPCs mtx.Gauge `stm:"dpc,1000,1"`
		// Total number of received and serviced hardware interrupts.
		Interrupts mtx.Gauge `stm:"interrupts,1000,1"`
		// Time that processor spent in different modes.
		Time cpuModes `stm:"time"`
		// Time spent in low-power idle state.
		CState struct {
			C1 mtx.Gauge `stm:"c1,1000,1"`
			C2 mtx.Gauge `stm:"c2,1000,1"`
			C3 mtx.Gauge `stm:"c3,1000,1"`
		} `stm:"cstate"`
	}

	cpuCores []*cpuCore
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

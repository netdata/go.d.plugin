package nvidia_smi

import (
	"fmt"

	"github.com/netdata/go-orchestrator/module"
	"github.com/netdata/go.d.plugin/pkg/nvml"
)

func init() {
	creator := module.Creator{
		Defaults: module.Defaults{
			Disabled: true,
		},
		Create: func() module.Module { return New() },
	}

	module.Register("nvidia_smi", creator)
}

// New creates Nvsmi with default values
func New() *Nvsmi {
	return &Nvsmi{
		metrics: make(map[string]int64),
	}
}

// Nvsmi module struct
type Nvsmi struct {
	module.Base // should be embedded by every module
	devices     []Device
	metrics     map[string]int64
}

// Device embedded nvml device
type Device struct {
	*nvml.Device
	ID    string
	Model string
}

// Cleanup makes cleanup
func (n *Nvsmi) Cleanup() {
	nvml.Shutdown()
}

// Init makes initialization
func (n *Nvsmi) Init() bool {
	nvml.Init()
	return true
}

// Check makes check
func (n *Nvsmi) Check() bool {
	version, err := nvml.GetDriverVersion()
	if err != nil {
		n.Errorf("failed to get driver version, %v", err)
		return false
	}
	n.Debugf("GPU driver version: %s", version)

	count, err := nvml.GetDeviceCount()
	if err != nil {
		n.Errorf("get device count failed, %v", err)
		return false
	}
	n.Debugf("GPU device count: %d", count)

	if count < 1 {
		n.Info("no GPU device present")
		return false
	}

	for i := uint(0); i < count; i++ {
		device, err := nvml.NewDevice(i)
		if err != nil {
			n.Errorf("Error getting device %d: %v\n", i, err)
		}
		d := Device{Device: device, Model: *device.Model, ID: fmt.Sprintf("gpu%d", i)}
		n.devices = append(n.devices, d)
	}
	return true
}

// Charts creates Charts dynamically (each gpu device will be a family)
func (n *Nvsmi) Charts() *Charts {
	allCharts := Charts{}
	for _, d := range n.devices {
		cs := *charts.Copy()
		family := fmt.Sprintf("%s %s", d.ID, d.Model)
		for _, c := range cs {
			c.ID = getRealID(c.ID, d.ID)
			c.Fam = family
			for _, dim := range c.Dims {
				dim.ID = getRealID(dim.ID, d.ID)
			}
			allCharts = append(allCharts, c)
		}
	}

	return &allCharts
}

// Collect collects metrics
func (n *Nvsmi) Collect() map[string]int64 {
	for i, d := range n.devices {
		st, err := d.Status()
		if err != nil || st == nil {
			n.Errorf("failed to get device %d status, %v", i, err)
		}

		// n.Debugf("GPU device %d status: %5d %5d %5d %5d %5d %5d %5d %5d %5d",
		// 	i, *st.Memory.ECCErrors.Device, *st.Utilization.GPU, *st.Utilization.Memory,
		// 	*st.Utilization.Encoder, *st.Utilization.Decoder, *st.Clocks.Memory, *st.Clocks.Cores, *st.PCI.Throughput.TX, *st.PCI.Throughput.RX)
		if st.PCI.Throughput.RX != nil {
			n.metrics[getRealID("pci_rx_throughput", d.ID)] = (int64)(*st.PCI.Throughput.RX)
		}
		if st.PCI.Throughput.TX != nil {
			n.metrics[getRealID("pci_tx_throughput", d.ID)] = (int64)(*st.PCI.Throughput.TX)
		}
		if st.Utilization.GPU != nil {
			n.metrics[getRealID("gpu_core_util", d.ID)] = (int64)(*st.Utilization.GPU)
		}
		if st.Utilization.Memory != nil {
			n.metrics[getRealID("memory_util", d.ID)] = (int64)(*st.Utilization.Memory)
		}
		if st.Utilization.Decoder != nil {
			n.metrics[getRealID("decoder_util", d.ID)] = (int64)(*st.Utilization.Decoder)
		}
		if st.Utilization.Encoder != nil {
			n.metrics[getRealID("encoder_util", d.ID)] = (int64)(*st.Utilization.Encoder)
		}
		if st.Memory.ECCErrors.Device != nil {
			n.metrics[getRealID("ecc_errors_device", d.ID)] = (int64)(*st.Memory.ECCErrors.Device)
		}
		if st.Memory.ECCErrors.L1Cache != nil {
			n.metrics[getRealID("ecc_errors_l1cache", d.ID)] = (int64)(*st.Memory.ECCErrors.L1Cache)
		}
		if st.Memory.ECCErrors.L2Cache != nil {
			n.metrics[getRealID("ecc_errors_l2cache", d.ID)] = (int64)(*st.Memory.ECCErrors.L2Cache)
		}
		if st.Memory.Global.Free != nil {
			n.metrics[getRealID("memory_usage_free", d.ID)] = (int64)(*st.Memory.Global.Free)
		}
		if st.Memory.Global.Used != nil {
			n.metrics[getRealID("memory_usage_used", d.ID)] = (int64)(*st.Memory.Global.Used)
		}
		if st.Clocks.Cores != nil {
			n.metrics[getRealID("core_clock", d.ID)] = (int64)(*st.Clocks.Cores)
		}
		if st.Clocks.Memory != nil {
			n.metrics[getRealID("mem_clock", d.ID)] = (int64)(*st.Clocks.Memory)
		}
		if st.Power != nil {
			n.metrics[getRealID("gpu_power", d.ID)] = (int64)(*st.Power)
		}
		if st.Temperature != nil {
			n.metrics[getRealID("gpu_temp", d.ID)] = (int64)(*st.Temperature)
		}
		computeCount := 0
		graphicsCount := 0
		computeAndGraphicsCount := 0
		for _, p := range st.Processes {
			switch p.Type {
			case nvml.Compute:
				computeCount++
			case nvml.Graphics:
				graphicsCount++
			case nvml.ComputeAndGraphics:
				computeAndGraphicsCount++
			}
		}
		n.metrics[getRealID("gpu_compute_process", d.ID)] = (int64)(computeCount)
		n.metrics[getRealID("gpu_graphics_process", d.ID)] = (int64)(graphicsCount)
		n.metrics[getRealID("gpu_compute_and_graphics_process", d.ID)] = (int64)(computeAndGraphicsCount)
	}
	return n.metrics
}

func getRealID(key string, model string) string {
	return fmt.Sprintf("%s_%s", key, model)
}

package nvidia_smi

import "github.com/netdata/go-orchestrator/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

var charts = Charts{
	{
		ID:    "pci_bandwidth",
		Title: "PCI Express Bandwidth Utilization",
		Units: "KB/s",
		Type:  module.Area,
		Ctx:   "nvidia_smi.pci_bandwidth",
		Dims: Dims{
			{ID: "pci_rx_throughput", Name: "rx", Algo: module.Absolute, Div: 1, Mul: 1},
			{ID: "pci_tx_throughput", Name: "tx", Algo: module.Absolute, Div: 1, Mul: -1},
		},
	},
	{
		ID:    "gpu_core",
		Title: "GPU Utilization",
		Units: "%",
		Type:  module.Line,
		Ctx:   "nvidia_smi.gpu_utilization",
		Dims: Dims{
			{ID: "gpu_core_util", Name: "utilization", Algo: module.Absolute, Div: 1, Mul: 1},
		},
	},
	{
		ID:    "mem_util",
		Title: "Memory Bandwidth Utilization",
		Units: "%",
		Type:  module.Line,
		Ctx:   "nvidia_smi.mem_utilization",
		Dims: Dims{
			{ID: "memory_util", Name: "utilization", Algo: module.Absolute, Div: 1, Mul: 1},
		},
	},
	{
		ID:    "codec_util",
		Title: "Codec Utilization",
		Units: "%",
		Type:  module.Line,
		Ctx:   "nvidia_smi.codec_utilization",
		Dims: Dims{
			{ID: "encoder_util", Name: "encoder", Algo: module.Absolute, Div: 1, Mul: 1},
			{ID: "decoder_util", Name: "decoder", Algo: module.Absolute, Div: 1, Mul: 1},
		},
	},
	{
		ID:    "mem_allocated",
		Title: "Memory Allocated",
		Units: "MB",
		Type:  module.Line,
		Ctx:   "nvidia_smi.memory_allocated",
		Dims: Dims{
			{ID: "memory_usage_used", Name: "used", Algo: module.Absolute, Div: 1, Mul: 1},
			{ID: "memory_usage_free", Name: "free", Algo: module.Absolute, Div: 1, Mul: 1},
		},
	},
	{
		ID:    "temperature",
		Title: "Temperature",
		Units: "celsius",
		Type:  module.Line,
		Ctx:   "nvidia_smi.temperature",
		Dims: Dims{
			{ID: "gpu_temp", Name: "temp", Algo: module.Absolute, Div: 1, Mul: 1},
		},
	},
	{
		ID:    "clock",
		Title: "Clock Frequencies",
		Units: "MHz",
		Type:  module.Line,
		Ctx:   "nvidia_smi.clocks",
		Dims: Dims{
			{ID: "core_clock", Name: "core", Algo: module.Absolute, Div: 1, Mul: 1},
			{ID: "mem_clock", Name: "mem", Algo: module.Absolute, Div: 1, Mul: 1},
		},
	},
	{
		ID:    "power",
		Title: "Power Utilization",
		Units: "Watts",
		Type:  module.Line,
		Ctx:   "nvidia_smi.power",
		Dims: Dims{
			{ID: "gpu_power", Name: "power", Algo: module.Absolute, Div: 1, Mul: 1},
		},
	},
	{
		ID:    "ecc_errors",
		Title: "Memory ECC Erros",
		Units: "counts/s",
		Type:  module.Line,
		Ctx:   "nvidia_smi.ecc_errors",
		Dims: Dims{
			{ID: "ecc_errors_device", Name: "device", Algo: module.Incremental, Div: 1, Mul: 1},
			{ID: "ecc_errors_l1cache", Name: "l1cache", Algo: module.Incremental, Div: 1, Mul: 1},
			{ID: "ecc_errors_l2cache", Name: "l2cache", Algo: module.Incremental, Div: 1, Mul: 1},
		},
	},
	{
		ID:    "process",
		Title: "Running GPU Processes",
		Units: "counts",
		Type:  module.Line,
		Ctx:   "nvidia_smi.process",
		Dims: Dims{
			{ID: "gpu_compute_process", Name: "compute", Algo: module.Absolute, Div: 1, Mul: 1},
			{ID: "gpu_graphics_process", Name: "graphics", Algo: module.Absolute, Div: 1, Mul: 1},
			{ID: "gpu_compute_and_graphics_process", Name: "combined", Algo: module.Absolute, Div: 1, Mul: 1},
		},
	},
}

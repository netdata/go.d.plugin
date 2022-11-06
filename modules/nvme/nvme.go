package nvme

import (
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	module.Register("nvme", module.Creator{
		Defaults: module.Defaults{
			//Disabled:    true,
			UpdateEvery: 5,
		},
		Create: func() module.Module { return New() },
	})
}

func New() *NVMe {
	return &NVMe{
		Config: Config{
			BinaryPath: "nvme",
			Timeout:    web.Duration{Duration: time.Second * 5},
		},
		charts:           &module.Charts{},
		newNVMeCLI:       newNVMeCLIExec,
		devicePaths:      make(map[string]bool),
		listDevicesEvery: time.Minute * 10,
	}

}

type Config struct {
	Timeout    web.Duration
	BinaryPath string `yaml:"binary_path"`
}

type (
	NVMe struct {
		module.Base
		Config `yaml:",inline"`

		charts *module.Charts

		newNVMeCLI func(cfg Config) (*nvmeCLIExec, error)
		exec       nvmeCLI

		devicePaths      map[string]bool
		listDevicesTime  time.Time
		listDevicesEvery time.Duration
	}
	nvmeCLI interface {
		list() (*nvmeDeviceList, error)
		smartLog(devicePath string) (*nvmeDeviceSmartLog, error)
	}
)

func (n *NVMe) Init() bool {
	return true
}

func (n *NVMe) Check() bool {
	return len(n.Collect()) > 0
}

func (n *NVMe) Charts() *module.Charts {
	return n.charts
}

func (n *NVMe) Collect() map[string]int64 {
	mx, err := n.collect()
	if err != nil {
		n.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (n *NVMe) Cleanup() {}

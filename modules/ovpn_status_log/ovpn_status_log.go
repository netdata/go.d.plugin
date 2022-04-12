package ovpn_status_log

import (
	//"fmt"
	"os"

	"github.com/netdata/go.d.plugin/agent/module"
)

const (
	defaultFilePath = "/var/log/openvpn/status.log"
)

func init() {
	creator := module.Creator{
		Defaults: module.Defaults{
			Disabled: true,
		},
		Create: func() module.Module { return New() },
	}

	module.Register("ovpn_status_log", creator)
}

// New creates VPNStatus with default values.
func New() *VPNStatus {
	config := Config{
		StatusPath: defaultFilePath,
	}
	return &VPNStatus{
		Config: config,
		charts: charts.Copy(),
		//collectedUsers: make(map[string]bool),
	}
}

// Config is the VPNStatus module configuration.
type Config struct {
	StatusPath string `yaml:"log_path"`
}

// VPNStatus VPNStatus module.
type VPNStatus struct {
	module.Base
	Config `yaml:",inline"`
	charts *module.Charts
	//fileHandle     *os.File
	//collectedUsers map[string]bool
}

// Cleanup makes cleanup.
func (o *VPNStatus) Cleanup() {
}

// Init makes initialization.
func (o *VPNStatus) Init() bool {
	if _, err := os.Stat(o.StatusPath); err != nil {
		o.Errorf("status log file doesn't exist: %v", err)
		return false
	}

	f, err := os.Open(o.StatusPath)
	if err != nil {
		o.Errorf("Error opening file: %v", err)
		return false
	}
	f.Close()
	o.Errorf("----------Error opening file: %v", err)
	return true
}

// Check makes check.
func (o *VPNStatus) Check() bool {
	return true
}

// Charts creates Charts.
func (o VPNStatus) Charts() *module.Charts {
	return o.charts
}

// Collect collects metrics.
func (o *VPNStatus) Collect() map[string]int64 {
	mx, err := o.collect()
	if err != nil {
		o.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

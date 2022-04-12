package openvpn_status

import (
	"os"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/matcher"
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

	module.Register("openvpn_status", creator)
}

// New creates OpenVPNStatus with default values.
func New() *OpenVPNStatus {
	config := Config{
		StatusPath: defaultFilePath,
	}
	return &OpenVPNStatus{
		Config:         config,
		charts:         charts.Copy(),
		collectedUsers: make(map[string]bool),
	}
}

// Config is the OpenVPNStatus module configuration.
type Config struct {
	StatusPath   string             `yaml:"log_path"`
	PerUserStats matcher.SimpleExpr `yaml:"per_user_stats"`
}

// OpenVPNStatus OpenVPNStatus module.
type OpenVPNStatus struct {
	module.Base
	Config         `yaml:",inline"`
	charts         *module.Charts
	collectedUsers map[string]bool
	perUserMatcher matcher.Matcher
}

// Init makes initialization.
func (o *OpenVPNStatus) Init() bool {
	if _, err := os.Stat(o.StatusPath); err != nil {
		o.Errorf("status file read error: %v", err)
		return false
	}

	if !o.PerUserStats.Empty() {
		m, err := o.PerUserStats.Parse()
		if err != nil {
			o.Errorf("error on creating per user stats matcher : %v", err)
			return false
		}
		o.perUserMatcher = matcher.WithCache(m)
	}

	return true
}

// Check makes check.
func (o *OpenVPNStatus) Check() bool {
	return true
}

// Charts creates Charts.
func (o OpenVPNStatus) Charts() *module.Charts {
	return o.charts
}

// Collect collects metrics.
func (o *OpenVPNStatus) Collect() map[string]int64 {
	mx, err := o.collect()
	if err != nil {
		o.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

// Cleanup makes cleanup.
func (o *OpenVPNStatus) Cleanup() {
}

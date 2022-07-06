// SPDX-License-Identifier: GPL-3.0-or-later

package openvpn_status_log

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

	module.Register("openvpn_status_log", creator)
}

// New creates OpenVPNStatusLog with default values.
func New() *OpenVPNStatusLog {
	config := Config{
		StatusPath: defaultFilePath,
	}
	return &OpenVPNStatusLog{
		Config:         config,
		charts:         charts.Copy(),
		collectedUsers: make(map[string]bool),
	}
}

// Config is the OpenVPNStatusLog module configuration.
type Config struct {
	StatusPath   string             `yaml:"log_path"`
	PerUserStats matcher.SimpleExpr `yaml:"per_user_stats"`
}

// OpenVPNStatusLog OpenVPNStatusLog module.
type OpenVPNStatusLog struct {
	module.Base
	Config         `yaml:",inline"`
	charts         *module.Charts
	collectedUsers map[string]bool
	perUserMatcher matcher.Matcher
}

// Init makes initialization.
func (o *OpenVPNStatusLog) Init() bool {
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
func (o *OpenVPNStatusLog) Check() bool {
	if _, err := os.Stat(o.StatusPath); err != nil {
		o.Errorf("file read error: %v", err)
		return false
	}

	return true
}

// Charts creates Charts.
func (o OpenVPNStatusLog) Charts() *module.Charts {
	return o.charts
}

// Collect collects metrics.
func (o *OpenVPNStatusLog) Collect() map[string]int64 {
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
func (o *OpenVPNStatusLog) Cleanup() {
}

// SPDX-License-Identifier: GPL-3.0-or-later

package chrony

import (
	"net"
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
)

type (
	Config struct {
		Protocol string `yaml:"protocol"`
		Address  string `yaml:"address"`
		Timeout  int    `yaml:"timeout"` // Millisecond
	}

	// Chrony is the main collector for chrony
	Chrony struct {
		module.Base   // should be embedded by every module
		Config        `yaml:",inline"`
		chronyVersion uint8
		latestSource  net.IP
		conn          net.Conn
		charts        *module.Charts
	}
)

var (
	// chronyCmdAddr is the chrony local port
	chronyDefaultProtocol = "udp"
	chronyDefaultCmdAddr  = "127.0.0.1:323"
	chronyDefaultTimeout  = 1
)

func init() {
	creator := module.Creator{
		Defaults: module.Defaults{
			Disabled: true,
		},
		Create: func() module.Module { return New() },
	}

	module.Register("chrony", creator)
}

// New creates Chrony exposing local status of a chrony daemon
func New() *Chrony {
	return &Chrony{
		Config: Config{
			Protocol: chronyDefaultProtocol,
			Address:  chronyDefaultCmdAddr,
			Timeout:  1,
		},
		charts:       &charts,
		latestSource: net.IPv4zero,
	}
}

// Cleanup makes cleanup
func (c *Chrony) Cleanup() {
}

// Init makes initialization
func (c *Chrony) Init() bool {
	if c.Timeout <= 0 {
		c.Timeout = chronyDefaultTimeout
	}

	conn, err := net.DialTimeout(c.Protocol, c.Address, time.Duration(c.Timeout)*time.Millisecond)
	if err != nil {
		c.Errorf(
			"unable connect to chrony addr %s:%s err: %s, is chrony up and running?",
			c.Protocol, c.Address, err)
		return false
	}

	c.conn = conn
	return true
}

// Check makes check
func (c *Chrony) Check() bool {
	err := c.applyChronyVersion()
	if err != nil {
		c.Errorf("get chrony version failed with err: %s", err)
		return false
	}

	return true
}

// Charts creates Charts dynamically
func (c *Chrony) Charts() *module.Charts {
	return c.charts
}

// Collect collects metrics
func (c *Chrony) Collect() map[string]int64 {
	// collect all we need and sent Exception to sentry
	res := map[string]int64{"running": 0}

	if !c.running() {
		return res
	}
	res["running"] = 1

	tra := c.collectTracking()
	for k, v := range tra {
		res[k] = v
	}

	act := c.collectActivity()
	for k, v := range act {
		res[k] = v
	}

	return res
}

func (c *Chrony) running() bool {
	err := c.submitEmptyRequest()
	if err != nil {
		c.Errorf("contract chrony failed with err: %s", err)
		return false
	}
	return true
}

func (c *Chrony) collectTracking() (res map[string]int64) {
	res = make(map[string]int64)
	tracking, err := c.fetchTracking()
	if err != nil {
		c.Errorf("fetch tracking status failed: %s", err)
		res["running"] = 0
		return
	}
	c.Debugf(tracking.String())

	res["running"] = 1
	res["stratum"] = (int64)(tracking.Stratum)
	res["leap_status"] = (int64)(tracking.LeapStatus)
	res["root_delay"] = (int64)(tracking.RootDelay.Int64())
	res["root_dispersion"] = (int64)(tracking.RootDispersion.Int64())
	res["skew"] = (int64)(tracking.SkewPpm.Int64())
	res["frequency"] = (int64)(tracking.FreqPpm.Int64())
	res["last_offset"] = (int64)(tracking.LastOffset.Int64())
	res["rms_offset"] = (int64)(tracking.RmsOffset.Int64())
	res["update_interval"] = (int64)(tracking.LastUpdateInterval.Int64())
	res["current_correction"] = (int64)(tracking.LastUpdateInterval.Int64())
	res["ref_timestamp"] = tracking.RefTime.Time().Unix()

	sourceIp := tracking.Ip.Ip()

	if !sourceIp.Equal(c.latestSource) {
		chart := c.charts.Get("source")
		_ = chart.AddDim(&module.Dim{
			ID: sourceIp.String(), Name: sourceIp.String(), Algo: module.Absolute, Div: 1, Mul: 1,
		})
		_ = chart.RemoveDim(c.latestSource.String())

		// you should let go.d.plugin know that something has been changed, and print dimension again.
		chart.MarkNotCreated()

		c.Debugf("source change from %s to %s", c.latestSource, sourceIp)
		c.latestSource = sourceIp
	}
	res[c.latestSource.String()] = 1

	if sourceIp.Equal(net.IPv4zero) || sourceIp.Equal(net.IPv6zero) {
		c.Warningf("chrony not select valid upstream")
	}

	return
}

func (c *Chrony) collectActivity() (res map[string]int64) {
	res = make(map[string]int64)
	activity, err := c.fetchActivity()
	if err != nil {
		c.Errorf("fetch activity status failed: %s", err)
		return
	}
	c.Debug(activity.String())

	res["online_sources"] = int64(activity.Online)
	res["offline_sources"] = int64(activity.Offline)
	res["burst_online_sources"] = int64(activity.BurstOnline)
	res["burst_offline_sources"] = int64(activity.BurstOffline)
	res["unresolved_sources"] = int64(activity.Unresolved)

	if activity.Online == 0 {
		c.Warningf("chrony have no available upstream")
	}
	return res
}

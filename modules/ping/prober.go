// SPDX-License-Identifier: GPL-3.0-or-later

package ping

import (
	"fmt"
	"time"

	"github.com/netdata/go.d.plugin/logger"

	probing "github.com/prometheus-community/pro-bing"
)

func newPingProber(conf pingProberConfig, log *logger.Logger) prober {
	return &pingProber{
		privileged: conf.privileged,
		packets:    conf.packets,
		interval:   conf.interval,
		deadline:   conf.deadline,
		Logger:     log,
	}
}

type pingProberConfig struct {
	privileged bool
	packets    int
	interval   time.Duration
	deadline   time.Duration
}

type pingProber struct {
	privileged bool
	packets    int
	interval   time.Duration
	deadline   time.Duration
	*logger.Logger
}

func (p *pingProber) ping(host string) (*probing.Statistics, error) {
	pr := probing.New(host)

	if err := pr.Resolve(); err != nil {
		return nil, fmt.Errorf("DNS lookup '%s' : %v", host, err)
	}

	pr.RecordRtts = false
	pr.Interval = p.interval
	pr.Count = p.packets
	pr.Timeout = p.deadline
	pr.SetPrivileged(p.privileged)
	pr.SetLogger(nil)

	if err := pr.Run(); err != nil {
		return nil, fmt.Errorf("pinging host '%s' (ip %s): %v", pr.Addr(), pr.IPAddr(), err)
	}

	stats := pr.Statistics()

	p.Debugf("ping stats for host '%s' (ip '%s'): %+v", pr.Addr(), pr.IPAddr(), stats)

	return stats, nil
}

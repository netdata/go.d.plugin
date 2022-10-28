// SPDX-License-Identifier: GPL-3.0-or-later

package ping

import (
	"errors"
	"time"
)

func (p *Ping) validateConfig() error {
	if len(p.Hosts) == 0 {
		return errors.New("'hosts' can't be empty")
	}
	if p.SendPackets <= 0 {
		return errors.New("'send_packets' can't be <= 0")
	}
	return nil
}

func (p *Ping) initProber() (prober, error) {
	deadline := time.Millisecond * time.Duration(float64(p.UpdateEvery)*1.5*1000)
	if deadline.Milliseconds() == 0 {
		return nil, errors.New("zero ping deadline")
	}

	conf := pingProberConfig{
		privileged: p.Privileged,
		packets:    p.SendPackets,
		interval:   p.Interval.Duration,
		deadline:   deadline,
	}

	return p.newProber(conf, p.Logger), nil
}

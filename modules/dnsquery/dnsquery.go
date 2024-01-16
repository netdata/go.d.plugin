// SPDX-License-Identifier: GPL-3.0-or-later

package dnsquery

import (
	_ "embed"
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/miekg/dns"
)

//go:embed "config_schema.json"
var configSchema string

func init() {
	module.Register("dns_query", module.Creator{
		JobConfigSchema: configSchema,
		Defaults: module.Defaults{
			UpdateEvery: 5,
		},
		Create: func() module.Module { return New() },
	})
}

func New() *DNSQuery {
	return &DNSQuery{
		Config: Config{
			Timeout:     web.Duration(time.Second * 2),
			Network:     "udp",
			RecordTypes: []string{"A"},
			Port:        53,
		},
		newDNSClient: func(network string, timeout time.Duration) dnsClient {
			return &dns.Client{
				Net:         network,
				ReadTimeout: timeout,
			}
		},
	}
}

type Config struct {
	UpdateEvery int          `yaml:"update_every" json:"update_every"`
	Timeout     web.Duration `yaml:"timeout" json:"timeout"`
	Domains     []string     `yaml:"domains" json:"domains"`
	Servers     []string     `yaml:"servers" json:"servers"`
	Network     string       `yaml:"network" json:"network"`
	RecordType  string       `yaml:"record_type" json:"record_type"`
	RecordTypes []string     `yaml:"record_types" json:"record_types"`
	Port        int          `yaml:"port" json:"port"`
}

type (
	DNSQuery struct {
		module.Base
		Config `yaml:",inline" json:""`

		charts *module.Charts

		dnsClient    dnsClient
		newDNSClient func(network string, duration time.Duration) dnsClient

		recordTypes map[string]uint16
	}
	dnsClient interface {
		Exchange(msg *dns.Msg, address string) (response *dns.Msg, rtt time.Duration, err error)
	}
)

func (d *DNSQuery) Configuration() any {
	return d.Config
}

func (d *DNSQuery) Init() error {
	if err := d.verifyConfig(); err != nil {
		d.Errorf("config validation: %v", err)
		return err
	}

	rt, err := d.initRecordTypes()
	if err != nil {
		d.Errorf("init record type: %v", err)
		return err
	}
	d.recordTypes = rt

	charts, err := d.initCharts()
	if err != nil {
		d.Errorf("init charts: %v", err)
		return err
	}
	d.charts = charts

	return nil
}

func (d *DNSQuery) Check() error {
	return nil
}

func (d *DNSQuery) Charts() *module.Charts {
	return d.charts
}

func (d *DNSQuery) Collect() map[string]int64 {
	mx, err := d.collect()
	if err != nil {
		d.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (d *DNSQuery) Cleanup() {}

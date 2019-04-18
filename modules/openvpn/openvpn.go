package openvpn

import (
	"strings"
	"time"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

const (
	defaultAddress        = "127.0.0.1:7505"
	defaultConnectTimeout = time.Second * 2
	defaultReadTimeout    = time.Second * 2
	defaultWriteTimeout   = time.Second * 2
)

func init() {
	creator := module.Creator{
		Defaults: module.Defaults{
			Disabled: true,
		},
		Create: func() module.Module { return New() },
	}

	module.Register("openvpn", creator)
}

// New creates OpenVPN with default values.
func New() *OpenVPN {
	config := Config{
		Address: defaultAddress,
		Timeouts: timeouts{
			Connect: web.Duration{Duration: defaultConnectTimeout},
			Read:    web.Duration{Duration: defaultReadTimeout},
			Write:   web.Duration{Duration: defaultWriteTimeout},
		},
	}
	return &OpenVPN{Config: config}
}

// Config is the OpenVPN module configuration.
type Config struct {
	Address  string
	Timeouts timeouts `yaml:",inline"`
}

type timeouts struct {
	Connect web.Duration `yaml:"connect_timeout"`
	Read    web.Duration `yaml:"read_timeout"`
	Write   web.Duration `yaml:"write_timeout"`
}

// OpenVPN OpenVPN module.
type OpenVPN struct {
	module.Base
	Config    `yaml:",inline"`
	apiClient apiClient
}

// Cleanup makes cleanup.
func (o *OpenVPN) Cleanup() {
	if o.apiClient == nil {
		return
	}
	_ = o.apiClient.disconnect()
}

// Init makes initialization.
func (o *OpenVPN) Init() bool {
	var network = "tcp"
	if strings.HasPrefix(o.Address, "/") {
		network = "unix"
	}
	config := clientConfig{
		network: network,
		address: o.Address,
		timeouts: clientTimeouts{
			connect: o.Timeouts.Connect.Duration,
			read:    o.Timeouts.Read.Duration,
			write:   o.Timeouts.Write.Duration,
		},
	}
	o.Infof("using network: %s, address: %s, connect timeout: %s, read timeout: %s, write timeout: %s",
		network, o.Address, o.Timeouts.Connect.Duration, o.Timeouts.Read.Duration, o.Timeouts.Write.Duration)
	o.apiClient = newClient(config)

	return true
}

// Check makes check.
func (o *OpenVPN) Check() bool {
	if !o.apiClient.isConnected() {
		err := o.apiClient.connect()
		if err != nil {
			o.Error(err)
			return false
		}
	}
	ver, err := o.collectVersion()
	if err != nil {
		o.Error(err)
		o.Cleanup()
		return false
	}
	o.Infof("connected to OpenVPN v%d.%d.%d, Management v%d",
		ver.major, ver.minor, ver.patch, ver.management)
	return true
}

// Charts creates Charts.
func (OpenVPN) Charts() *Charts { return charts.Copy() }

// Collect collects metrics.
func (o *OpenVPN) Collect() map[string]int64 {
	mx, err := o.collect()

	if err != nil {
		o.Error(err)
		return nil
	}

	return mx
}

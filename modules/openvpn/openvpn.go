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
)

func init() {
	creator := module.Creator{
		// TODO: enable
		// DisabledByDefault:true,
		Create: func() module.Module { return New() },
	}

	module.Register("openvpn", creator)
}

// New creates OpenVPN with default values.
func New() *OpenVPN {
	config := Config{
		Address:        defaultAddress,
		ConnectTimeout: web.Duration{Duration: defaultConnectTimeout},
		ReadTimeout:    web.Duration{Duration: defaultReadTimeout},
	}
	return &OpenVPN{Config: config}
}

// Config is the OpenVPN module configuration.
type Config struct {
	Address        string
	ConnectTimeout web.Duration `yaml:"connect_timeout"`
	ReadTimeout    web.Duration `yaml:"read_timeout"`
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
	if len(o.Address) == 0 {
		o.Error("mandatory 'address' parameter is not set")
		return false
	}
	var network = "tcp"
	if strings.HasPrefix(o.Address, "/") {
		network = "unix"
	}
	o.apiClient = newAPIClient(clientConfig{
		network:        network,
		address:        o.Address,
		connectTimeout: o.ConnectTimeout.Duration,
		readTimeout:    o.ReadTimeout.Duration,
	})
	return true
}

// Check makes check.
func (o *OpenVPN) Check() bool { return len(o.Collect()) > 0 }

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

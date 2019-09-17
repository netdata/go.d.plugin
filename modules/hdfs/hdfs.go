package hdfs

import (
	"time"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("hdfs", creator)
}

// New creates HDFS with default values.
func New() *HDFS {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{
				UserURL: "http://127.0.0.1:9870/jmx",
			},
			Client: web.Client{
				Timeout: web.Duration{Duration: time.Second}},
		},
	}

	return &HDFS{
		Config: config,
		charts: charts.Copy(),
	}
}

// Config is the HDFS module configuration.
type Config struct {
	web.HTTP `yaml:",inline"`
}

// HDFS HDFS module.
type HDFS struct {
	module.Base
	Config `yaml:",inline"`

	client *client
	charts *Charts
}

// Cleanup makes cleanup.
func (HDFS) Cleanup() {}

// Init makes initialization.
func (h *HDFS) Init() bool {
	httpClient, err := web.NewHTTPClient(h.Client)
	if err != nil {
		h.Error(err)
		return false
	}

	h.client = newClient(httpClient, h.Request)
	return true
}

// Check makes check.
func (h HDFS) Check() bool {
	return len(h.Collect()) > 0
}

// Charts returns Charts.
func (h HDFS) Charts() *module.Charts {
	return h.charts
}

// Collect collects metrics.
func (h *HDFS) Collect() map[string]int64 {
	mx, err := h.collect()

	if err != nil {
		h.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}

	return mx
}

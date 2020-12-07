package energid

import (
	"net/http"
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/web"
)

type Config struct {
	web.HTTP `yaml:",inline"`
}

type Energid struct {
	module.Base
	Config `yaml:",inline"`

	httpClient *http.Client
	charts     *module.Charts
}

func New() *Energid {
	return &Energid{
		Config: Config{
			HTTP: web.HTTP{
				Request: web.Request{
					URL: "http://127.0.0.1:9796",
				},
				Client: web.Client{
					Timeout: web.Duration{Duration: time.Second},
				},
			},
		},
	}
}

func (d *Energid) Init() bool {
	err := d.validateConfig()
	if err != nil {
		d.Errorf("config validation: %v", err)
		return false
	}

	client, err := d.initHTTPClient()
	if err != nil {
		d.Errorf("init HTTP client: %v", err)
		return false
	}
	d.httpClient = client

	cs, err := d.initCharts()
	if err != nil {
		d.Errorf("init charts: %v", err)
		return false
	}
	d.charts = cs

	return true
}

func (d *Energid) Check() bool {
	return len(d.Collect()) > 0
}

func (d *Energid) Charts() *module.Charts {
	return d.charts
}

func (d *Energid) Collect() map[string]int64 {
	ms, err := d.collect()
	if err != nil {
			d.Error(err)
	}

	if len(ms) == 0 {
			return nil
	}

	return ms
}

func (d *Energid) Cleanup() {
	if d.httpClient == nil {
			return
	}

	d.httpClient.CloseIdleConnections()
}

package dnsdist

import (
	"net/http"

	"github.com/netdata/go.d.plugin/agent/module"
)

func init() {
	module.Register("dnsdist", module.Creator{
		Defaults: module.Defaults{
			UpdateEvery: 1,
		},
		Create: func() module.Module { return New() },
	})
}

type (
	Config struct {
		Url string   `yaml:"url"`
		User string  `yaml:"user"`
		Pass string  `yaml:"pass"`
		Headers []HeaderValues `yaml:"headers"`
	}

	HeaderValues struct {
		Name string   `yaml:"name"`
		Value string  `yaml:"value"`
	}
)

type DNSdist struct {
	module.Base
	Config `yaml:",inline"`

	httpClient    *http.Client
	charts        *module.Charts
	collected     map[string]int64
}

func New() *DNSdist {
	return &DNSdist{
		Config: Config {
			Url: "http://127.0.0.1:5053/jsonstat?command=stats",
			User: "netdata",
			Pass: "netdata",
			Headers: nil,
			/*
			HTTP: web.HTTP{
				Request: web.Request{
					URL: "http://127.0.0.1:5053/jsonstat?command=stats",
				},
				Client: web.Client{
					Timeout: web.Duration{Duration: time.Second},
				},
			},
			*/
		},
		collected: make(map[string]int64),
	}
}

func (d *DNSdist) Init() bool {
	err := d.validateConfig()
	if err != nil {
		d.Errorf("Config validation: %v", err)
		return false
	}

	/*
	client, err := d.initHTTPClient()
	if err != nil {
		d.Errorf("init HTTP client: %v", err)
		return false
	}
	d.httpClient = client
	*/

	c, err := d.initCharts()
	if err != nil {
		d.Errorf("init charts: %v", err)
		return false
	}
	d.charts = c

	return true
}

func (d *DNSdist) Check() bool {
	return len(d.Collect()) > 0
}

func (d *DNSdist) Charts() *module.Charts {
	return d.charts
}

func (d *DNSdist) Collect() map[string]int64 {
	return nil
}

func (d *DNSdist) Cleanup() {
	if d.httpClient == nil {
		return 
	}

	d.httpClient.CloseIdleConnections()
}
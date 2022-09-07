package phpfpm

import (
	"time"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go.d.plugin/agent/module"
)

func init() {
	module.Register("phpfpm", module.Creator{
		Create: func() module.Module { return New() },
	})
}

func New() *Phpfpm {
	return &Phpfpm{
		Config: Config{
			HTTP: web.HTTP{
				Request: web.Request{
					URL: "http://127.0.0.1/status?full&json",
				},
				Client: web.Client{
					Timeout: web.Duration{Duration: time.Second},
				},
			},
			FcgiPath: "/status",
		},
	}
}

type (
	Config struct {
		web.HTTP `yaml:",inline"`
		Socket   string `yaml:"socket"`
		Address  string `yaml:"address"`
		FcgiPath string `yaml:"fcgi_path"`
	}
	Phpfpm struct {
		module.Base
		Config `yaml:",inline"`

		client client
	}
)

func (p *Phpfpm) Init() bool {
	c, err := p.initClient()
	if err != nil {
		p.Errorf("init client: %v", err)
		return false
	}
	p.client = c
	return true
}

func (p *Phpfpm) Check() bool {
	return len(p.Collect()) > 0
}

func (Phpfpm) Charts() *Charts {
	return charts.Copy()
}

func (p *Phpfpm) Collect() map[string]int64 {
	mx, err := p.collect()
	if err != nil {
		p.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (Phpfpm) Cleanup() {}

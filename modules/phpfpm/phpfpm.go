package phpfpm

import (
	"errors"
	"os"
	"time"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Defaults: module.Defaults{
			Disabled: true,
		},
		Create: func() module.Module { return New() },
	}

	module.Register("phpfpm", creator)
}

func New() *Phpfpm {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{UserURL: "http://127.0.0.1/status?full&json"},
			Client:  web.Client{Timeout: web.Duration{Duration: time.Second}},
		},
	}

	return &Phpfpm{
		Config: config,
	}
}

type (
	Config struct {
		web.HTTP `yaml:",inline"`
		Socket     string `yaml:"socket"`
		env 		map[string]string

	}
	Phpfpm struct {
		module.Base
		Config `yaml:",inline"`
		client client
	}
)

func (p *Phpfpm) validateHttpConfig() error {
	if err := p.ParseUserURL(); err != nil {
		return err
	}
	if p.URL.Host == "" {
		return errors.New("URL is not set")
	}
	return nil
}

func (p *Phpfpm) initClient() error {
	cl, err := web.NewHTTPClient(p.Client)
	if err != nil {
		return err
	}

	p.client = newClient(cl, p.Request)
	return nil
}

func (p *Phpfpm) initSocket() error {

	env := make(map[string]string)

	env["SCRIPT_NAME"] = "/status"
	env["SCRIPT_FILENAME"] = "/status"
	env["SERVER_SOFTWARE"] = "go / fcgiclient "
	env["REMOTE_ADDR"] = "127.0.0.1"
	env["QUERY_STRING"] = "json&full"
	env["REQUEST_METHOD"] = "GET"
	env["CONTENT_TYPE"] = "application/json"
	p.Config.Socket = p.Socket

	p.client = &socketClient{
		socket:  p.Socket,
		env:     env,
		timeout: p.Timeout.Duration,
	}
	return nil

}
func (p *Phpfpm) validateSocketConfig() bool {

	if _, err := os.Stat(p.Socket); err != nil {
		p.Errorf("the socket does not exist: %v", err)
		return false
	}
	return true
}

func (p *Phpfpm) Init() bool {
	if len(p.Socket) > 0 {
		if p.validateSocketConfig() {
			err := p.initSocket()
			p.Debugf("using Socket %s", p.Socket)
			p.Debugf("using timeout: %s", p.Timeout.Duration)
			return err == nil
		}
		return false
	}


	if err := p.validateHttpConfig(); err != nil {
		p.Errorf("error on validating config: %v", err)
		return false
	}
	if err := p.initClient(); err != nil {
		p.Errorf("error on initializing client: %v", err)
		return false
	}

	p.Debugf("using URL %s", p.URL)
	p.Debugf("using timeout: %s", p.Timeout.Duration)
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

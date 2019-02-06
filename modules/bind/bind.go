package bind

import (
	"fmt"
	"github.com/netdata/go-orchestrator/module"
	"github.com/netdata/go.d.plugin/pkg/web"
	"net/url"
	"time"
)

func init() {
	creator := module.Creator{
		// DisabledByDefault: true,
		Create: func() module.Module { return New() },
	}

	module.Register("bind", creator)
}

const (
	defaultURL         = "http://100.127.0.91:8080/json/v1"
	defaultHTTPTimeout = time.Second
)

// New creates Bind with default values.
func New() *Bind {
	return &Bind{
		HTTP: web.HTTP{
			Request: web.Request{URL: defaultURL},
			Client:  web.Client{Timeout: web.Duration{Duration: defaultHTTPTimeout}},
		},
	}
}

type bindAPIClient interface {
	serverStats() (*serverStats, error)
}

// Bind bind module.
type Bind struct {
	module.Base

	web.HTTP `yaml:",inline"`

	bindAPIClient
}

// Cleanup makes cleanup.
func (Bind) Cleanup() {}

// Init makes initialization.
func (b *Bind) Init() bool {
	if b.URL == "" {
		b.Error("URL not set")
		return false
	}

	client, err := web.NewHTTPClient(b.Client)

	if err != nil {
		b.Error("error on creating http client : %v", err)
		return false
	}

	addr, err := url.Parse(b.URL)

	if err != nil {
		b.Errorf("error on parsing URL %s : %v", b.URL, err)
		return false
	}

	switch addr.Path {
	default:
		b.Errorf("URL %s is wrong", b.URL)
		return false
	case "":
		b.Error("WIP")
		return false
	case "/xml/v2":
		b.Error("WIP")
		return false
	case "/xml/v3":
		b.Error("WIP")
		return false
	case "/json/v1":
		b.bindAPIClient = &jsonClient{request: b.Request, httpClient: client}
	}

	return true
}

// Check makes check.
func (Bind) Check() bool {
	return true
}

// Charts creates Charts.
func (Bind) Charts() *Charts {
	return &module.Charts{}
}

// Collect collects metrics.
func (b *Bind) Collect() map[string]int64 {
	s, err := b.serverStats()

	if err != nil {
		b.Error(err)
	}

	fmt.Println(s)

	return nil
}

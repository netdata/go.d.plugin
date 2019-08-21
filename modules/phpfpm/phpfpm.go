package phpfpm

import (
	"math"
	"time"

	"github.com/netdata/go.d.plugin/pkg/stm"
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

const (
	defaultURL         = "http://127.0.0.1/status?full&json"
	defaultHTTPTimeout = time.Second
)

// Config is the php-fpm module configuration.
type Config struct {
	web.HTTP `yaml:",inline"`
}

// Phpfpm collets php-fpm metrics.
type Phpfpm struct {
	module.Base

	Config `yaml:",inline"`

	client *client
}

// New returns a php-fpm module with default values.
func New() *Phpfpm {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{UserURL: defaultURL},
			Client:  web.Client{Timeout: web.Duration{Duration: defaultHTTPTimeout}},
		},
	}

	return &Phpfpm{
		Config: config,
	}
}

// Init makes initialization.
func (n *Phpfpm) Init() bool {
	if err := n.ParseUserURL(); err != nil {
		n.Errorf("error on parsing url '%s' : %v", n.UserURL, err)
		return false
	}

	if n.URL.Host == "" {
		n.Error("URL is not set")
		return false
	}

	client, err := web.NewHTTPClient(n.Client)
	if err != nil {
		n.Error(err)
		return false
	}

	n.client = newClient(client, n.Request, )

	n.Debugf("using URL %s", n.URL)
	n.Debugf("using timeout: %s", n.Timeout.Duration)

	return true
}

// Check checks the module can collect metrics.
func (n *Phpfpm) Check() bool {
	return len(n.Collect()) > 0
}

// Charts creates Charts.
func (*Phpfpm) Charts() *Charts {
	return charts.Copy()
}

// Collect returns collected metrics.
func (n *Phpfpm) Collect() map[string]int64 {
	status, err := n.client.Status()
	if err != nil {
		n.Error(err)
		return nil
	}

	data := stm.ToMap(status)
	if len(status.Processes) == 0 {
		return data
	}

	statProcesses(data, status.Processes, "ReqDur", func(p proc) int64 { return p.Duration })
	statProcesses(data, status.Processes, "ReqCpu", func(p proc) int64 { return int64(p.CPU) })
	statProcesses(data, status.Processes, "ReqMem", func(p proc) int64 { return p.Memory })

	return data
}

// Cleanup makes cleanup.
func (*Phpfpm) Cleanup() {}

type accessor func(p proc) int64

func statProcesses(m map[string]int64, procs []proc, met string, acc accessor) {
	var sum, count, min, max int64
	for _, proc := range procs {
		if proc.State != "Idle" {
			continue
		}

		val := acc(proc)
		sum += val
		count += 1
		if count == 1 {
			min, max = val, val
			continue
		}
		min = int64(math.Min(float64(min), float64(val)))
		max = int64(math.Max(float64(max), float64(val)))
	}

	m["min"+met] = min
	m["max"+met] = max
	m["avg"+met] = sum / count
}

package portcheck

import (
	"fmt"
	"sort"
	"time"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Defaults: module.Defaults{
			Disabled:    true,
			UpdateEvery: 5,
		},
		Create: func() module.Module { return New() },
	}

	module.Register("portcheck", creator)
}

const (
	defaultHTTPTimeout = time.Second
)

// New creates PortCheck with default values
func New() *PortCheck {
	return &PortCheck{
		Timeout: web.Duration{Duration: defaultHTTPTimeout},

		task:     make(chan *port),
		taskDone: make(chan struct{}),
		ports:    make([]*port, 0),
		metrics:  make(map[string]int64),
	}

}

type state string

var (
	success state = "success"
	timeout state = "timeout"
	failed  state = "failed"
)

type port struct {
	number      int
	state       state
	inState     int
	updateEvery int
	latency     time.Duration
}

func (p *port) setState(s state) {
	changed := p.state != s

	if changed {
		p.inState = p.updateEvery
		p.state = s
	} else {
		p.inState += p.updateEvery
	}
}

func (p port) stateText() string {
	switch p.state {
	case success:
		return fmt.Sprintf("success_%d", p.number)
	case timeout:
		return fmt.Sprintf("timeout_%d", p.number)
	case failed:
		return fmt.Sprintf("failed_%d", p.number)
	}
	panic("unknown state")
}

func newPort(number, updateEvery int) *port {
	return &port{
		number:      number,
		updateEvery: updateEvery,
	}
}

// PortCheck portcheck module
type PortCheck struct {
	module.Base

	Host        string       `yaml:"host" validate:"required"`
	Ports       []int        `yaml:"ports" validate:"required,gte=1"`
	Timeout     web.Duration `yaml:"timeout"`
	UpdateEvery int          `yaml:"update_every"`

	task     chan *port
	taskDone chan struct{}

	ports   []*port
	workers []*worker

	metrics map[string]int64
}

// Cleanup makes cleanup
func (tc *PortCheck) Cleanup() {
	if len(tc.workers) == 0 {
		return
	}
	close(tc.task)
	tc.workers = make([]*worker, 0)
}

// Init makes initialization
func (tc *PortCheck) Init() bool {
	sort.Ints(tc.Ports)

	for _, p := range tc.Ports {
		tc.ports = append(tc.ports, newPort(p, tc.UpdateEvery))
		tc.workers = append(tc.workers, newWorker(tc.Host, tc.Timeout.Duration, tc.task, tc.taskDone))
	}

	tc.Debugf("using host %s", tc.Host)
	tc.Debugf("using ports %v", tc.Ports)
	tc.Debugf("using HTTP timeout: %s", tc.Timeout.Duration)

	return true
}

// Check makes check
func (PortCheck) Check() bool {
	return true
}

// Charts creates    charts
func (tc PortCheck) Charts() *Charts {
	var charts module.Charts

	for _, p := range tc.Ports {
		_ = charts.Add(chartsTemplate(p)...)
	}

	return &charts
}

// Collect collects metrics
func (tc *PortCheck) Collect() map[string]int64 {
	for _, p := range tc.ports {
		tc.task <- p
	}

	for i := 0; i < len(tc.ports); i++ {
		<-tc.taskDone
	}

	for _, p := range tc.ports {
		tc.metrics[fmt.Sprintf("success_%d", p.number)] = 0
		tc.metrics[fmt.Sprintf("failed_%d", p.number)] = 0
		tc.metrics[fmt.Sprintf("timeout_%d", p.number)] = 0

		tc.metrics[p.stateText()] = 1
		tc.metrics[fmt.Sprintf("instate_%d", p.number)] = int64(p.inState)
		tc.metrics[fmt.Sprintf("latency_%d", p.number)] = int64(p.latency)
	}

	return tc.metrics
}

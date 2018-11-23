package portcheck

import (
	"net"
	"sort"
	"time"

	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/utils"
)

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
	if p.state != s {
		p.inState = p.updateEvery
		p.state = s
	} else {
		p.inState += p.updateEvery
	}
}

func (p port) stateText() string {
	switch p.state {
	case success:
		return sprintf("success_%d", p.number)
	case timeout:
		return sprintf("timeout_%d", p.number)
	case failed:
		return sprintf("failed_%d", p.number)
	}
	panic("unknown state")
}

func newPort(number, updateEvery int) *port {
	return &port{
		number:      number,
		updateEvery: updateEvery,
	}
}

type PortCheck struct {
	modules.Base

	Host        string         `yaml:"host" validate:"required"`
	Ports       []int          `yaml:"ports" validate:"required,gte=1"`
	Timeout     utils.Duration `yaml:"timeout"`
	UpdateEvery int            `yaml:"update_every"`

	doCh   chan *port
	doneCh chan struct{}

	ports   []*port
	workers []*worker

	data map[string]int64
}

func New() *PortCheck {
	return &PortCheck{
		doCh:   make(chan *port),
		doneCh: make(chan struct{}),
		ports:  make([]*port, 0),
		data:   make(map[string]int64),
	}

}

func (pc *PortCheck) Init() bool {
	if pc.Timeout.Duration == 0 {
		pc.Timeout.Duration = time.Second
	}
	pc.Debugf("using timeout: %s", pc.Timeout.Duration)

	ips, err := net.LookupIP(pc.Host)
	if err != nil {
		return false
	}

	pc.Host = ips[len(ips)-1].String()
	pc.Debugf("using %s:%v", pc.Host, pc.Ports)

	sort.Ints(pc.Ports)

	for _, p := range pc.Ports {
		pc.ports = append(pc.ports, newPort(p, pc.UpdateEvery))

		worker := newWorker(pc.Host, pc.Timeout.Duration, pc.doCh, pc.doneCh)
		pc.workers = append(pc.workers, worker)
		go worker.start()
	}

	return true
}

func (pc PortCheck) Check() bool {
	return true
}

func (pc *PortCheck) Cleanup() {
	for _, worker := range pc.workers {
		worker.stop()
	}
	close(pc.doCh)
	close(pc.doneCh)
}

func (pc PortCheck) GetCharts() *modules.Charts {
	charts := modules.Charts{}
	for _, p := range pc.Ports {
		charts.Add(chartsTemplate(p)...)
	}

	return &charts
}

func (pc *PortCheck) GetData() map[string]int64 {
	for _, p := range pc.ports {
		pc.doCh <- p
	}

	for i := 0; i < len(pc.ports); i++ {
		<-pc.doneCh
	}

	for _, p := range pc.ports {
		pc.data[sprintf("success_%d", p.number)] = 0
		pc.data[sprintf("failed_%d", p.number)] = 0
		pc.data[sprintf("timeout_%d", p.number)] = 0

		pc.data[p.stateText()] = 1
		pc.data[sprintf("instate_%d", p.number)] = int64(p.inState)
		pc.data[sprintf("latency_%d", p.number)] = int64(p.latency)

	}
	return pc.data
}

func init() {
	modules.Register(
		"portcheck",
		modules.Creator{
			UpdateEvery: 5,
			Create:      func() modules.Module { return New() },
		},
	)
}

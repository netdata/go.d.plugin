package tcpcheck

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

type TcpCheck struct {
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

func New() *TcpCheck {
	return &TcpCheck{
		doCh:   make(chan *port),
		doneCh: make(chan struct{}),
		ports:  make([]*port, 0),
		data:   make(map[string]int64),
	}

}

func (tc *TcpCheck) Init() bool {
	if tc.Timeout.Duration == 0 {
		tc.Timeout.Duration = time.Second
	}
	tc.Debugf("using timeout: %s", tc.Timeout.Duration)

	ips, err := net.LookupIP(tc.Host)
	if err != nil {
		return false
	}

	tc.Host = ips[len(ips)-1].String()
	tc.Debugf("using %s:%v", tc.Host, tc.Ports)

	sort.Ints(tc.Ports)

	for _, p := range tc.Ports {
		tc.ports = append(
			tc.ports,
			newPort(p, tc.UpdateEvery),
		)

		tc.workers = append(
			tc.workers,
			newWorker(tc.Host, tc.Timeout.Duration, tc.doCh, tc.doneCh),
		)
	}

	return true
}

func (tc TcpCheck) Check() bool {
	return true
}

func (tc *TcpCheck) Cleanup() {
	for _, worker := range tc.workers {
		worker.stop()
	}
}

func (tc TcpCheck) GetCharts() *modules.Charts {
	charts := modules.Charts{}
	for _, p := range tc.Ports {
		charts.Add(chartsTemplate(p)...)
	}

	return &charts
}

func (tc *TcpCheck) GetData() map[string]int64 {
	for _, p := range tc.ports {
		tc.doCh <- p
	}

	for i := 0; i < len(tc.ports); i++ {
		<-tc.doneCh
	}

	for _, p := range tc.ports {
		tc.data[sprintf("success_%d", p.number)] = 0
		tc.data[sprintf("failed_%d", p.number)] = 0
		tc.data[sprintf("timeout_%d", p.number)] = 0

		tc.data[p.stateText()] = 1
		tc.data[sprintf("instate_%d", p.number)] = int64(p.inState)
		tc.data[sprintf("latency_%d", p.number)] = int64(p.latency)

	}
	return tc.data
}

func init() {
	creator := modules.Creator{
		UpdateEvery: 5,
		Create:      func() modules.Module { return New() },
	}

	modules.Register("tcpcheck", creator)
}

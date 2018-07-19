package portcheck

import (
	"net"
	"sort"
	"time"

	"github.com/l2isbad/go.d.plugin/internal/modules"
	"github.com/l2isbad/go.d.plugin/internal/pkg/utils"
)

const (
	_ = iota
	success
	timeout
	failed
)

type state struct {
	curr     int
	prev     int
	duration int
	u        int
}

func (s *state) set(v int) {
	if v != s.curr {
		s.duration = s.u
		s.prev = s.curr
		s.curr = v
	} else {
		s.duration += s.u
	}
}

type port struct {
	number  int
	state   state
	latency time.Duration
}

func (p port) stateText() string {
	switch p.state.curr {
	case success:
		return sprintf("success_%d", p.number)
	case timeout:
		return sprintf("timeout_%d", p.number)
	case failed:
		return sprintf("failed_%d", p.number)
	}
	return ""
}

func newPort(p, u int) *port {
	return &port{
		number: p,
		state:  state{u: u},
	}
}

type PortCheck struct {
	modules.Charts
	modules.BaseConfHook
	modules.Logger

	Host    string         `yaml:"host,required"`
	Ports   []int          `yaml:"ports,required"`
	Timeout utils.Duration `yaml:"timeout"`

	do    chan *port
	done  chan struct{}
	ports []*port
	data  map[string]int64
}

func (pc *PortCheck) Check() bool {
	ips, err := net.LookupIP(pc.Host)
	if err != nil {
		return false
	}

	pc.Host = ips[len(ips)-1].String()
	pc.Debugf("Using %s:%v", pc.Host, pc.Ports)

	if pc.Timeout.Duration == 0 {
		pc.Timeout.Duration = time.Second
	}
	pc.Debugf("Using timeout: %s", pc.Timeout.Duration)

	sort.Ints(pc.Ports)
	for _, p := range pc.Ports {
		pc.ports = append(pc.ports, newPort(p, pc.GetUpdateEvery()))
		pc.AddMany(charts(p))

		go worker(pc.Host, pc.Timeout.Duration, pc.do, pc.done)
	}

	return true
}

func (pc *PortCheck) GetData() map[string]int64 {
	for _, p := range pc.ports {
		pc.do <- p
	}

	for i := 0; i < len(pc.ports); i++ {
		<-pc.done
	}

	for _, p := range pc.ports {
		pc.data[sprintf("success_%d", p.number)] = 0
		pc.data[sprintf("failed_%d", p.number)] = 0
		pc.data[sprintf("timeout_%d", p.number)] = 0

		pc.data[p.stateText()] = int64(p.state.curr)
		pc.data[sprintf("instate_%d", p.number)] = int64(p.state.duration)
		pc.data[sprintf("latency_%d", p.number)] = int64(p.latency)

	}
	return pc.data
}

func worker(host string, dialTimeout time.Duration, doCh chan *port, doneCh chan struct{}) {
	for p := range doCh {
		t := time.Now()
		c, err := net.DialTimeout("tcp", sprintf("%s:%d", host, p.number), dialTimeout)
		p.latency = time.Since(t)

		if err == nil {
			p.state.set(success)
			c.Close()
		} else {
			v, ok := err.(interface{ Timeout() bool })

			if ok && v.Timeout() {
				p.state.set(timeout)
			} else {
				p.state.set(failed)
			}
		}

		doneCh <- struct{}{}
	}
}

func init() {
	modules.SetDefault().SetUpdateEvery(5)

	f := func() modules.Module {
		return &PortCheck{
			do:    make(chan *port),
			done:  make(chan struct{}),
			ports: make([]*port, 0),
			data:  make(map[string]int64)}
	}
	modules.Add(f)
}

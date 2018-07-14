package portcheck

import (
	"fmt"
	"net"
	"sort"
	"time"

	"github.com/l2isbad/go.d.plugin/internal/modules"
	"github.com/l2isbad/go.d.plugin/internal/pkg/charts/raw"
	"github.com/l2isbad/go.d.plugin/internal/pkg/utils"
)

const (
	conSuccess = iota
	conTimeout
	conFailed
)

type (
	Charts      = raw.Charts
	Order       = raw.Order
	Definitions = raw.Definitions
	Chart       = raw.Chart
	Options     = raw.Options
	Dimensions  = raw.Dimensions
	Dimension   = raw.Dimension
)

func chartsTemplate(port int) *Charts {
	family := fmt.Sprintf("port %d", port)
	return &raw.Charts{
		Order: Order{
			fmt.Sprintf("status_%d", port),
			fmt.Sprintf("instate_%d", port),
			fmt.Sprintf("latency_%d", port)},
		Definitions: Definitions{
			Chart{
				ID:      fmt.Sprintf("status_%d", port),
				Options: Options{"Port Check Status", "boolean", family, "portcheck.status"},
				Dimensions: Dimensions{
					Dimension{fmt.Sprintf("success_%d", port), "success"},
					Dimension{fmt.Sprintf("failed_%d", port), "failed"},
					Dimension{fmt.Sprintf("timeout_%d", port), "timeout"},
				},
			},
			Chart{
				ID:      fmt.Sprintf("instate_%d", port),
				Options: Options{"Current State Duration", "seconds", family, "portcheck.instate"},
				Dimensions: Dimensions{
					Dimension{fmt.Sprintf("instate_%d", port), "time"},
				},
			},
			Chart{
				ID:      fmt.Sprintf("latency_%d", port),
				Options: Options{"TCP Connect Latency", "ms", family, "portcheck.latency"},
				Dimensions: Dimensions{
					Dimension{fmt.Sprintf("latency_%d", port), "time", "", 1, 1e6},
				},
			},
		},
	}
}

type port struct {
	num        int
	curState   int
	prevState  int
	inState    int
	conLatency time.Duration
}

func (p *port) setState(s int) {
	p.prevState = p.curState
	p.curState = s
}

func (p *port) addStateTo(m map[string]int64) {
	switch p.curState {
	case conSuccess:
		m[fmt.Sprintf("success_%d", p.num)] = 1
	case conTimeout:
		m[fmt.Sprintf("timeout_%d", p.num)] = 1
	default:
		m[fmt.Sprintf("failed_%d", p.num)] = 1
	}
}

type PortCheck struct {
	modules.Charts
	modules.BaseConfHook
	modules.Logger

	Host    string         `yaml:"host"`
	Ports   []int          `yaml:"ports"`
	Timeout utils.Duration `yaml:"timeout"`

	doCh   chan *port
	doneCh chan *port
	ports  []*port
	data   map[string]int64
}

func (pc *PortCheck) Check() bool {
	ips, err := net.LookupIP(pc.Host)
	if err != nil {
		return false
	}
	pc.Host = ips[len(ips)-1].String()
	pc.Logger.Debugf("Using %s:%v", pc.Host, pc.Ports)

	if pc.Timeout.Duration == 0 {
		pc.Timeout.Duration = time.Duration(pc.GetUpdateEvery()) * time.Second
		pc.Logger.Warningf("timeout not specified. Setting to %s", pc.Timeout.Duration)
	}
	sort.Ints(pc.Ports)
	for _, p := range pc.Ports {
		pc.ports = append(pc.ports, &port{num: p})
		pc.AddMany(chartsTemplate(p))
		go connWorker(pc.Host, pc.Timeout.Duration, pc.doCh, pc.doneCh)
	}
	return true
}

func (pc *PortCheck) GetData() map[string]int64 {
	for _, p := range pc.ports {
		pc.data[fmt.Sprintf("success_%d", p.num)] = 0
		pc.data[fmt.Sprintf("failed_%d", p.num)] = 0
		pc.data[fmt.Sprintf("timeout_%d", p.num)] = 0
		pc.doCh <- p
	}

	for i := 0; i < len(pc.ports); i++ {
		p := <-pc.doneCh
		p.addStateTo(pc.data)
		if p.curState != p.prevState {
			p.inState = pc.GetUpdateEvery()
		} else {
			p.inState += pc.GetUpdateEvery()
		}
		pc.data[fmt.Sprintf("instate_%d", p.num)] = int64(p.inState)
		pc.data[fmt.Sprintf("latency_%d", p.num)] = int64(p.conLatency)
	}
	return pc.data
}

func connWorker(host string, timeout time.Duration, doCh chan *port, doneCh chan *port) {
	for p := range doCh {
		t := time.Now()
		c, err := net.DialTimeout(
			"tcp",
			fmt.Sprintf("%s:%d", host, p.num),
			timeout)
		p.conLatency = time.Since(t)
		switch err {
		case nil:
			p.setState(conSuccess)
			c.Close()
			doneCh <- p
		default:
			if v, ok := err.(interface{ Timeout() bool }); ok && v.Timeout() {
				p.setState(conTimeout)
				doneCh <- p
			} else {
				p.setState(conFailed)
				doneCh <- p
			}
		}
	}
}

func init() {
	modules.SetDefault(modules.UpdateEvery).Set(5)

	f := func() modules.Module {
		return &PortCheck{
			doCh:   make(chan *port),
			doneCh: make(chan *port),
			ports:  make([]*port, 0),
			data:   make(map[string]int64)}
	}
	modules.Add(f)
}

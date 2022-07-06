// SPDX-License-Identifier: GPL-3.0-or-later

package dnsquery

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/miekg/dns"
	"github.com/netdata/go.d.plugin/agent/module"
)

func init() {
	creator := module.Creator{
		Defaults: module.Defaults{
			UpdateEvery: 5,
		},
		Create: func() module.Module { return New() },
	}

	module.Register("dns_query", creator)
}

const (
	defaultTimeout    = time.Second * 2
	defaultNetwork    = "udp"
	defaultRecordType = "A"
	defaultPort       = 53
)

// New creates DNSQuery with default values
func New() *DNSQuery {
	return &DNSQuery{
		Timeout:    web.Duration{Duration: defaultTimeout},
		Network:    defaultNetwork,
		RecordType: defaultRecordType,
		Port:       defaultPort,

		task:             make(chan task),
		taskDone:         make(chan struct{}),
		exchangerFactory: newExchanger,
		servers:          make([]*server, 0),
		workers:          make([]*worker, 0),
	}
}

type exchanger interface {
	Exchange(msg *dns.Msg, address string) (response *dns.Msg, rtt time.Duration, err error)
}

func newExchanger(network string, timeout time.Duration) exchanger {
	return &dns.Client{
		Net:         network,
		ReadTimeout: timeout,
	}
}

type server struct {
	id   string
	name string
	port int

	resp *dns.Msg
	rtt  time.Duration
	err  error
}

// DNSQuery dnsquery module
type DNSQuery struct {
	module.Base

	Domains    []string
	Servers    []string
	Network    string
	RecordType string `yaml:"record_type"`
	Port       int
	Timeout    web.Duration

	task     chan task
	taskDone chan struct{}

	exchangerFactory func(network string, duration time.Duration) exchanger

	rtype   uint16
	servers []*server
	workers []*worker
}

// Cleanup makes cleanup
func (d *DNSQuery) Cleanup() {
	if len(d.workers) == 0 {
		return
	}
	close(d.task)
	d.workers = make([]*worker, 0)
}

func (d *DNSQuery) setup() error {
	if len(d.Domains) == 0 {
		return errors.New("no domains specified")
	}

	if len(d.Servers) == 0 {
		return errors.New("no servers specified")
	}

	if !(d.Network == "" || d.Network == "udp" || d.Network == "tcp" || d.Network == "tcp-tls") {
		return fmt.Errorf("wrong network transport : %s", d.Network)
	}

	rtype, err := parseRecordType(d.RecordType)

	if err != nil {
		return fmt.Errorf("error on parsing record type : %s", err)
	}
	d.rtype = rtype

	return nil
}

// Init makes initialization
func (d *DNSQuery) Init() bool {
	if err := d.setup(); err != nil {
		d.Error(err)
		return false
	}

	exch := d.exchangerFactory(d.Network, d.Timeout.Duration)

	for _, srv := range d.Servers {
		d.servers = append(d.servers, &server{id: serverNameReplacer.Replace(srv), name: srv, port: d.Port})
		// newWorker spawns worker goroutine
		d.workers = append(d.workers, newWorker(exch, d.task, d.taskDone))
	}

	return true
}

// Check makes check
func (DNSQuery) Check() bool {
	return true
}

// Charts creates Charts
func (d DNSQuery) Charts() *Charts {
	charts := charts.Copy()

	for _, srv := range d.servers {
		chart := charts.Get("query_time")
		dim := &Dim{ID: srv.id, Name: srv.name, Div: 1000000}

		if err := chart.AddDim(dim); err != nil {
			d.Errorf("error on creating charts : %s", err)
			return nil
		}
	}

	return charts
}

// Collect collects metrics
func (d *DNSQuery) Collect() map[string]int64 {
	domain := randomDomain(d.Domains)
	d.Debugf("current domain : %s", domain)

	for _, srv := range d.servers {
		d.task <- task{server: srv, domain: domain, rtype: d.rtype}
	}

	for range d.servers {
		<-d.taskDone
	}

	metrics := make(map[string]int64)

	for _, srv := range d.servers {
		if srv.resp != nil && srv.resp.Rcode != dns.RcodeSuccess {
			d.Errorf("invalid answer from %s after %s query for %s", srv.name, d.RecordType, domain)
			continue
		}

		if srv.err != nil {
			d.Debugf("error on querying %s after %s query for %s : %s", srv.name, d.RecordType, domain, srv.err)
			continue
		}

		metrics[srv.id] = srv.rtt.Nanoseconds()
	}

	if len(metrics) == 0 {
		return nil
	}

	return metrics
}

func parseRecordType(recordType string) (uint16, error) {
	var rtype uint16

	switch recordType {
	case "A":
		rtype = dns.TypeA
	case "AAAA":
		rtype = dns.TypeAAAA
	case "ANY":
		rtype = dns.TypeANY
	case "CNAME":
		rtype = dns.TypeCNAME
	case "MX":
		rtype = dns.TypeMX
	case "NS":
		rtype = dns.TypeNS
	case "PTR":
		rtype = dns.TypePTR
	case "SOA":
		rtype = dns.TypeSOA
	case "SPF":
		rtype = dns.TypeSPF
	case "SRV":
		rtype = dns.TypeSRV
	case "TXT":
		rtype = dns.TypeTXT
	default:
		return 0, fmt.Errorf("unknown record type : %s", recordType)
	}

	return rtype, nil
}

func randomDomain(domains []string) string {
	rand.Seed(time.Now().UnixNano())
	return domains[rand.Intn(len(domains))]
}

var serverNameReplacer = strings.NewReplacer(".", "_")

package bind

import (
	"net/url"
	"strings"
	"time"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
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
		charts: &Charts{},
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
	charts *Charts
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
		b.bindAPIClient = newJSONClient(client, b.Request)
	}

	return true
}

// Check makes check.
func (Bind) Check() bool {
	return true
}

// Charts creates Charts.
func (b Bind) Charts() *Charts {
	return b.charts
}

// Collect collects metrics.
func (b *Bind) Collect() map[string]int64 {
	metrics := make(map[string]int64)

	s, err := b.serverStats()
	if err != nil {
		b.Error(err)
	}
	b.collectServerStats(metrics, s)

	return metrics
}

func (b *Bind) collectServerStats(metrics map[string]int64, stats *serverStats) {
	var chart *Chart

	if len(stats.NSStats) > 0 {
		var chartID, dimName string

		for k, v := range stats.NSStats {
			switch {
			default:
				continue
			case k == "Requestv4":
				dimName = "IPv4"
				chartID = keyReceivedRequests
			case k == "Requestv6":
				dimName = "IPv6"
				chartID = keyReceivedRequests
			case k == "QryFailure":
				dimName = "failures"
				chartID = keyQueryFailures
			case k == "QryUDP":
				dimName = "UDP"
				chartID = keyProtocolsQueries
			case k == "QryTCP":
				dimName = "TCP"
				chartID = keyProtocolsQueries
			case k == "QrySuccess":
				dimName = "queries"
				chartID = keyQueriesSuccess
			case strings.HasSuffix(k, "QryRej"):
				dimName = k
				chartID = keyQueryFailuresDetail
			case strings.HasPrefix(k, "Qry"):
				dimName = k
				chartID = keyQueriesAnalysis
			case strings.HasPrefix(k, "Update"):
				dimName = k
				chartID = keyReceivedUpdates
			}

			if !b.charts.Has(chartID) {
				_ = b.charts.Add(charts[chartID].Copy())
			}

			chart = b.charts.Get(chartID)

			if !chart.HasDim(k) {
				_ = chart.AddDim(&Dim{ID: k, Name: dimName, Algo: module.Incremental})
				chart.MarkNotCreated()
			}

			delete(stats.NSStats, k)
			metrics[k] = v
		}
	}

	if len(stats.NSStats) > 0 {
		if !b.charts.Has(keyNSStats) {
			_ = b.charts.Add(charts[keyNSStats].Copy())
		}

		chart = b.charts.Get(keyNSStats)

		for k, v := range stats.NSStats {
			if !chart.HasDim(k) {
				_ = chart.AddDim(&Dim{ID: k, Algo: module.Incremental})
				chart.MarkNotCreated()
			}

			metrics[k] = v
		}
	}

	if len(stats.OpCodes) > 0 {
		if !b.charts.Has(keyInOpCodes) {
			_ = b.charts.Add(charts[keyInOpCodes].Copy())
		}

		chart = b.charts.Get(keyInOpCodes)

		for k, v := range stats.OpCodes {
			if !chart.HasDim(k) {
				_ = chart.AddDim(&Dim{ID: k, Algo: module.Incremental})
				chart.MarkNotCreated()
			}

			metrics[k] = v
		}
	}

	if len(stats.QTypes) > 0 {
		if !b.charts.Has(keyInQTypes) {
			_ = b.charts.Add(charts[keyInQTypes].Copy())
		}

		chart = b.charts.Get(keyInQTypes)

		for k, v := range stats.QTypes {
			if !chart.HasDim(k) {
				_ = chart.AddDim(&Dim{ID: k, Algo: module.Incremental})
				chart.MarkNotCreated()
			}

			metrics[k] = v
		}
	}

	if len(stats.SockStats) > 0 {
		if !b.charts.Has(keyInSockStats) {
			_ = b.charts.Add(charts[keyInSockStats].Copy())
		}

		chart = b.charts.Get(keyInSockStats)

		for k, v := range stats.SockStats {
			if !chart.HasDim(k) {
				_ = chart.AddDim(&Dim{ID: k, Algo: module.Incremental})
				chart.MarkNotCreated()
			}

			metrics[k] = v
		}
	}
}

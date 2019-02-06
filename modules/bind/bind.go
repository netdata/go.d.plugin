package bind

import (
	"net/url"
	"strings"
	"time"

	"github.com/netdata/go.d.plugin/pkg/matcher"
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

	web.HTTP    `yaml:",inline"`
	ViewsFilter string `yaml:"views_filter"`

	views matcher.Matcher
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

	if b.ViewsFilter != "" {
		if b.views, err = matcher.Parse(b.ViewsFilter); err != nil {
			b.Errorf("error on creating views matcher : %v", err)
			return false
		}
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
		for k, v := range stats.NSStats {
			var (
				algo    = module.Incremental
				dimName = k
				chartID = ""
			)

			switch {
			default:
				continue
			case k == "RecursClients":
				dimName = "clients"
				chartID = keyRecursiveClients
				algo = module.Absolute
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
				chartID = keyQueryFailuresDetail
			case strings.HasPrefix(k, "Qry"):
				chartID = keyQueriesAnalysis
			case strings.HasPrefix(k, "Update"):
				chartID = keyReceivedUpdates
			}

			if !b.charts.Has(chartID) {
				_ = b.charts.Add(charts[chartID].Copy())
			}

			chart = b.charts.Get(chartID)

			if !chart.HasDim(k) {
				_ = chart.AddDim(&Dim{ID: k, Name: dimName, Algo: algo})
				chart.MarkNotCreated()
			}

			delete(stats.NSStats, k)
			metrics[k] = v
		}
	}

	for _, v := range []struct {
		item map[string]int64
		key  string
	}{
		{item: stats.NSStats, key: keyNSStats},
		{item: stats.OpCodes, key: keyInOpCodes},
		{item: stats.QTypes, key: keyInQTypes},
		{item: stats.SockStats, key: keyNSStats},
	} {
		if !b.charts.Has(v.key) {
			_ = b.charts.Add(charts[v.key].Copy())
		}

		chart = b.charts.Get(v.key)

		for key, val := range v.item {
			if !chart.HasDim(key) {
				_ = chart.AddDim(&Dim{ID: key, Algo: module.Incremental})
				chart.MarkNotCreated()
			}

			metrics[key] = val
		}
	}
}

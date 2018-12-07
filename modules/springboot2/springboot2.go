package springboot2

import (
	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/utils"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	modules.Register("springboot2", modules.Creator{
		Create: func() modules.Module { return New() },
	})
}

// New returns SpringBoot2 instance with default values
func New() *SpringBoot2 {
	return &SpringBoot2{}
}

// SpringBoot2 Spring boot 2 module
type SpringBoot2 struct {
	modules.Base

	web.HTTP `yaml:",inline"`

	prom prometheus.Prometheus
}

type metrics struct {
	ThreadsDaemon int64 `stm:"threads_daemon"`
	Threads       int64 `stm:"threads"`

	Resp1xx int64 `stm:"resp_1xx"`
	Resp2xx int64 `stm:"resp_2xx"`
	Resp3xx int64 `stm:"resp_3xx"`
	Resp4xx int64 `stm:"resp_4xx"`
	Resp5xx int64 `stm:"resp_5xx"`
}

// Cleanup Cleanup
func (SpringBoot2) Cleanup() {}

// Init makes initialization
func (s *SpringBoot2) Init() bool {
	s.prom = prometheus.New(s.CreateHTTPClient(), s.RawRequest)
	return true
}

// Check makes check
func (s *SpringBoot2) Check() bool {
	rawMetrics, err := s.prom.Scrape()
	if err != nil {
		s.Error(err)
		return false
	}
	jvmMemory := rawMetrics.FindByName("jvm_memory_used_bytes")

	return len(jvmMemory) > 0
}

// Charts creates Charts
func (SpringBoot2) Charts() *Charts {
	return charts.Copy()
}

// GatherMetrics gathers metrics
func (s *SpringBoot2) GatherMetrics() map[string]int64 {
	rawMetrics, err := s.prom.Scrape()
	if err != nil {
		return nil
	}

	var m metrics
	m.ThreadsDaemon = int64(rawMetrics.FindByName("jvm_threads_daemon").Max())
	m.Threads = int64(rawMetrics.FindByName("jvm_threads_live").Max())

	for _, metric := range rawMetrics.FindByName("http_server_requests_seconds_count") {
		status := metric.Labels.Get("status")
		if status == "" {
			continue
		}
		value := int64(metric.Value)
		switch status[0] {
		case '1':
			m.Resp1xx += value
		case '2':
			m.Resp2xx += value
		case '3':
			m.Resp3xx += value
		case '4':
			m.Resp4xx += value
		case '5':
			m.Resp5xx += value
		}
	}
	return utils.ToMap(m)
}

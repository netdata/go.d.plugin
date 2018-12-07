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

type data struct {
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

// Init Init
func (s *SpringBoot2) Init() bool {
	s.prom = prometheus.New(s.CreateHTTPClient(), s.RawRequest)
	return true
}

// Check Check
func (s *SpringBoot2) Check() bool {
	metrics, err := s.prom.Scrape()
	if err != nil {
		s.Error(err)
		return false
	}
	jvmMemory := metrics.FindByName("jvm_memory_used_bytes")

	return len(jvmMemory) > 0
}

func (SpringBoot2) Charts() *Charts {
	return charts.Copy()
}

// GatherMetrics GatherMetrics
func (s *SpringBoot2) GatherMetrics() map[string]int64 {
	metrics, err := s.prom.Scrape()
	if err != nil {
		return nil
	}

	var d data
	d.ThreadsDaemon = int64(metrics.FindByName("jvm_threads_daemon").Max())
	d.Threads = int64(metrics.FindByName("jvm_threads_live").Max())

	for _, metric := range metrics.FindByName("http_server_requests_seconds_count") {
		status := metric.Labels.Get("status")
		if status == "" {
			continue
		}
		switch status[0] {
		case '1':
			d.Resp1xx++
		case '2':
			d.Resp2xx++
		case '3':
			d.Resp3xx++
		case '4':
			d.Resp4xx++
		case '5':
			d.Resp5xx++
		}
	}

	return utils.ToMap(d)
}

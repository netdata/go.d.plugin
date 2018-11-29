package springboot2

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"

	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/utils"
	"github.com/netdata/go.d.plugin/pkg/web"
)

// New returns Springboot2 instance with default values
func New() *Springboot2 {
	return &Springboot2{}
}

// Springboot2 Spring boot 2 module
type Springboot2 struct {
	modules.Base

	web.HTTP `yaml:",inline"`

	prom prometheus.Prometheus
}

type data struct {
	ThreadsDaemon int64 `stm:"threads_daemon"`
	Threads       int64 `stm:"threads"`
}

// Init Init
func (s *Springboot2) Init() bool {
	s.prom = prometheus.New(s.CreateHTTPClient(), s.RawRequest)
	return true
}

// Check Check
func (s *Springboot2) Check() bool {
	metrics, err := s.prom.Scrape()
	if err != nil {
		s.Error(err)
		return false
	}
	jvmMemory := metrics.FindByName("jvm_memory_used_bytes")

	return len(jvmMemory) > 0
}

func (Springboot2) Charts() *Charts {
	return charts.Copy()
}

// GatherMetrics GatherMetrics
func (s *Springboot2) GatherMetrics() map[string]int64 {
	metrics, err := s.prom.Scrape()
	if err != nil {
		return nil
	}

	var d data
	d.ThreadsDaemon = int64(metrics.FindByName("jvm_threads_daemon").Max())
	d.Threads = int64(metrics.FindByName("jvm_threads_live").Max())

	return utils.ToMap(d)
}

func init() {
	creator := modules.Creator{
		Create: func() modules.Module { return New() },
	}

	modules.Register("springboot2", creator)
}

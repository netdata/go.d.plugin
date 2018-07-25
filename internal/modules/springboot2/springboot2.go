package springboot2

import (
	"github.com/l2isbad/go.d.plugin/internal/modules"
	"github.com/l2isbad/go.d.plugin/internal/pkg/helpers/prometheus"
	"github.com/l2isbad/go.d.plugin/internal/pkg/helpers/web"
	"github.com/l2isbad/go.d.plugin/internal/pkg/utils"
)

// Springboot2 Spring boot 2 plugin
type Springboot2 struct {
	modules.Charts
	modules.Logger

	web.Request `yaml:",inline"`
	web.Client  `yaml:",inline"`

	prom prometheus.Prometheus
}

type data struct {
	ThreadsDaemon int64 `stm:"threads_daemon"`
	Threads       int64 `stm:"threads"`
}

// Check Check
func (s *Springboot2) Check() bool {
	s.AddMany(charts)
	s.prom = prometheus.New(s.Client.CreateHttpClient(), s.Request)
	metrics, err := s.prom.GetMetrics()
	if err != nil {
		s.Error(err)
		return false
	}
	jvmMemory := metrics.FindByName("jvm_memory_used_bytes")
	return len(jvmMemory) > 0
}

// GetData GetData
func (s *Springboot2) GetData() map[string]int64 {
	metrics, err := s.prom.GetMetrics()
	if err != nil {
		return nil
	}

	var d data
	d.ThreadsDaemon = int64(metrics.FindByName("jvm_threads_daemon").Max())
	d.Threads = int64(metrics.FindByName("jvm_threads_live").Max())
	return utils.StrToMap(d)
}

func init() {
	f := func() modules.Module {
		return &Springboot2{}
	}
	modules.Add(f)
}

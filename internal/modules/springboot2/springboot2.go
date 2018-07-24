package springboot2

import (
	"github.com/l2isbad/go.d.plugin/internal/modules"
	"github.com/l2isbad/go.d.plugin/internal/pkg/helpers/prometheus"
	"github.com/l2isbad/go.d.plugin/internal/pkg/helpers/web"
)

type Springboot2 struct {
	modules.Charts
	modules.Logger

	web.Request `yaml:",inline"`
	web.Client  `yaml:",inline"`

	prom prometheus.Prometheus
	data map[string]int64
}

func (s *Springboot2) Check() bool {
	return false
}

func (s *Springboot2) GetData() map[string]int64 {
	return nil
}

func init() {
	f := func() modules.Module {
		return &Springboot2{
			data: make(map[string]int64),
		}
	}
	modules.Add(f)
}

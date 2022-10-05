package pardot_zfs

import (
	"github.com/netdata/go.d.plugin/agent/module"
)

type Module interface {
	Init() bool
	Check() bool
	Charts() *module.Charts
	Collect() map[string]int64
	Cleanup()
}

type zfsMetric struct {
	module.Base
	pools []string
}

// Init enables metric collection - it must return true for the metric collection to happen
func (z *zfsMetric) Init() bool {
	return z.init()
}

func (z *zfsMetric) Check() bool {
	return len(z.Collect()) > 0
}

func (z *zfsMetric) Charts() *module.Charts {
	return z.createcharts()
}

func (z *zfsMetric) Collect() map[string]int64 {
	return z.collect()
}

func (z *zfsMetric) Cleanup() {}

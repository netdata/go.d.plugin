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

type PardotZFS struct {
	module.Base
	pools []string
}

// Init enables metric collection - it must return true for the metric collection to happen
func (p *PardotZFS) Init() bool {
	return p.init()
}

func (p *PardotZFS) Check() bool {
	return len(p.Collect()) > 0
}

func (p *PardotZFS) Charts() *module.Charts {
	return p.createcharts()
}

func (p *PardotZFS) Collect() map[string]int64 {
	return p.collect()
}

func (z *PardotZFS) Cleanup() {}

package zfs

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

type ZFS struct {
	module.Base
	pools []string
}

// Init enables metric collection - it must return true for the metric collection to happen
func (z *ZFS) Init() bool {
	return z.init()
}

func (z *ZFS) Check() bool {
	return len(z.Collect()) > 0
}

func (z *ZFS) Charts() *module.Charts {
	return z.createcharts()
}

func (z *ZFS) Collect() map[string]int64 {
	return z.collect()
}

func (z *ZFS) Cleanup() {}

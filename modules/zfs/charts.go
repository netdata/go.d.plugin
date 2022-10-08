package zfs

import (
	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/logger"
)

func (z *ZFS) createcharts() *module.Charts {
	chart := &module.Chart{
		ID:    "zfs",
		Title: "ZFS Free Space Fragmentation",
		Units: "percent",
		Fam:   "zfs",
		Ctx:   "zfs.fragmentation",
	}

	for _, v := range z.pools {
		d := module.Dim{
			ID: v,
		}
		err := chart.AddDim(&d)
		if err != nil {
			logger.Panicf("failed to add dim to chart: %v\n", err)
		}
	}

	charts := &module.Charts{}
	err := charts.Add(chart)
	if err != nil {
		logger.Panicf("failed to create module.Charts: %v\n", err)
	}

	return charts
}

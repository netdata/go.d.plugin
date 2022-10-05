package pardot_zfs

import "github.com/netdata/go.d.plugin/agent/module"

func (z *zfsMetric) createcharts() *module.Charts {
	chart := &module.Chart{
		ID:    "pardot_zfs",
		Title: "ZFS Free Space Fragmentation",
		Units: "percent",
		Fam:   "pardot_zfs",
	}

	for _, v := range z.pools {
		d := module.Dim{
			ID:   v,
			Name: v,
		}
		err := chart.AddDim(&d)
		if err != nil {
			panic("failed to add dim to chart")
		}
	}

	charts := &module.Charts{}
	err := charts.Add(chart)
	if err != nil {
		panic("failed to create module.Charts")
	}

	return charts
}

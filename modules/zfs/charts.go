package pardot_zfs

import "github.com/netdata/go.d.plugin/agent/module"

func (z *ZFS) createcharts() *module.Charts {
	chart := &module.Chart{
		ID:    "zfs",
		Title: "ZFS Free Space Fragmentation",
		Units: "percent",
		Fam:   "zfs",
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

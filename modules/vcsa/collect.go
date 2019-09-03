package vcsa

import (
	"sync"

	"github.com/netdata/go.d.plugin/pkg/stm"
)

type healthMetrics struct {
	ApplMgmt         *int `stm:"appl_mgmt"`
	DatabaseStorage  *int `stm:"database_storage"`
	Load             *int `stm:"load"`
	Mem              *int `stm:"mem"`
	SoftwarePackages *int `stm:"software_packages"`
	Storage          *int `stm:"storage"`
	Swap             *int `stm:"swap"`
	System           *int `stm:"system"`
}

func (vc *VCenter) collect() (map[string]int64, error) {
	err := vc.client.Ping()
	if err != nil {
		return nil, err
	}

	mx := &healthMetrics{}
	vc.scrapeHealth(mx, true)

	return stm.ToMap(mx), nil
}

func (vc *VCenter) scrapeHealth(mx *healthMetrics, doConcurrently bool) {
	type task func(*healthMetrics)

	var tasks = []task{
		vc.scrapeApplMgmt,
		vc.scrapeDatabaseStorage,
		vc.scrapeLoad,
		vc.scrapeMem,
		vc.scrapeSoftwarePackages,
		vc.scrapeStorage,
		vc.scrapeSwap,
		vc.scrapeSystem,
	}

	wg := &sync.WaitGroup{}
	wrap := func(call task) task {
		return func(metrics *healthMetrics) {
			call(metrics)
			wg.Done()
		}
	}
	for _, task := range tasks {
		if doConcurrently {
			wg.Add(1)
			task = wrap(task)
			go task(mx)
		} else {
			task(mx)
		}
	}
	wg.Wait()
}

//	The vCenter Server Appliance API offers health status indicators for several key components of the appliance:
// - green  The component is healthy.
// - yellow The component is healthy, but may have some problems.
// - orange The component is degraded, and may have serious problems.
// - red The component is unavailable, or will stop functioning soon.
// - gray No health data is available.
func decodeHealth(v string) int {
	switch v {
	default:
		return -1
	case "green":
		return 0
	case "yellow":
		return 1
	case "orange":
		return 2
	case "red":
		return 3
	case "gray":
		return 4
	}
}

func (vc *VCenter) scrapeApplMgmt(mx *healthMetrics) {
	v, err := vc.client.ApplMgmt()
	if err != nil {
		vc.Error(err)
		return
	}
	i := decodeHealth(v)
	mx.ApplMgmt = &i
}

func (vc *VCenter) scrapeDatabaseStorage(mx *healthMetrics) {
	v, err := vc.client.DatabaseStorage()
	if err != nil {
		vc.Error(err)
		return
	}
	i := decodeHealth(v)
	mx.DatabaseStorage = &i
}

func (vc *VCenter) scrapeLoad(mx *healthMetrics) {
	v, err := vc.client.Load()
	if err != nil {
		vc.Error(err)
		return
	}
	i := decodeHealth(v)
	mx.Load = &i
}

func (vc *VCenter) scrapeMem(mx *healthMetrics) {
	v, err := vc.client.Mem()
	if err != nil {
		vc.Error(err)
		return
	}
	i := decodeHealth(v)
	mx.Mem = &i
}

func (vc *VCenter) scrapeSoftwarePackages(mx *healthMetrics) {
	v, err := vc.client.SoftwarePackages()
	if err != nil {
		vc.Error(err)
		return
	}
	i := decodeHealth(v)
	mx.SoftwarePackages = &i
}

func (vc *VCenter) scrapeStorage(mx *healthMetrics) {
	v, err := vc.client.Storage()
	if err != nil {
		vc.Error(err)
		return
	}
	i := decodeHealth(v)
	mx.Storage = &i
}

func (vc *VCenter) scrapeSwap(mx *healthMetrics) {
	v, err := vc.client.Swap()
	if err != nil {
		vc.Error(err)
		return
	}
	i := decodeHealth(v)
	mx.Swap = &i
}

func (vc *VCenter) scrapeSystem(mx *healthMetrics) {
	v, err := vc.client.System()
	if err != nil {
		vc.Error(err)
		return
	}
	i := decodeHealth(v)
	mx.System = &i
}

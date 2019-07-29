package vsphere

import (
	"time"

	rs "github.com/netdata/go.d.plugin/modules/vsphere/resources"
)

func (vs *VSphere) goDiscovery(runEvery time.Duration) *task {
	discovery := func() {
		res, err := vs.Discover()
		if err != nil {
			vs.Errorf("error on discovering : %v", err)
			return
		}
		vs.consumeDiscovered(res)
	}
	return newTask(discovery, runEvery)
}

func (vs *VSphere) consumeDiscovered(res *rs.Resources) {
	vs.collectionLock.Lock()
	defer vs.collectionLock.Unlock()

	for _, h := range res.Hosts {
		if _, ok := vs.discoveredHosts[h.ID]; !ok {
			vs.discoveredHosts[h.ID] = 0
		}
	}
	for _, v := range res.VMs {
		if _, ok := vs.discoveredVMs[v.ID]; !ok {
			vs.discoveredVMs[v.ID] = 0
		}
	}

	vs.resources = res
}

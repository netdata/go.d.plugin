package vsphere

import (
	rs "github.com/netdata/go.d.plugin/modules/vsphere/resources"
)

func (vs *VSphere) goDiscovery() {
	if vs.discoveryTask != nil {
		vs.discoveryTask.stop()
	}
	vs.Infof("starting discovery process, will do discovery every %s", vs.DiscoveryInterval)

	job := func() {
		err := vs.discoverOnce()
		if err != nil {
			vs.Errorf("error on discovering : %v", err)
		}
	}
	vs.discoveryTask = newTask(job, vs.DiscoveryInterval.Duration)
}

func (vs *VSphere) discoverOnce() error {
	res, err := vs.Discover()
	if err != nil {
		return err
	}
	vs.consumeDiscovered(res)
	return nil
}

func (vs *VSphere) consumeDiscovered(res *rs.Resources) {
	vs.collectionLock.Lock()
	defer vs.collectionLock.Unlock()

	vs.updateDiscoveredItems(res)
	vs.resources = res
}

func (vs *VSphere) updateDiscoveredItems(res *rs.Resources) {
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
}

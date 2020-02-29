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
	vs.resources = res
}

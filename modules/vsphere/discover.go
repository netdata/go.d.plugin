package vsphere

import (
	"time"
)

func (vs *VSphere) goDiscovery(runEvery time.Duration) *task {
	discovery := func() {
		res, err := vs.Discover()
		if err != nil {
			vs.Errorf("error on discovering : %v", err)
			return
		}
		vs.resLock.Lock()
		defer vs.resLock.Unlock()
		vs.resources = res
	}
	return newTask(discovery, runEvery)
}

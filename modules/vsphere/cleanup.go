package vsphere

const (
	failedMax = 10
)

func (vs *VSphere) cleanupResources() {
	vs.cleanupHosts()
	vs.cleanupVMs()
}

func (vs *VSphere) cleanupHosts() {
	for k, v := range vs.failedUpdatesHosts {
		if v < failedMax {
			continue
		}
		delete(vs.chartedHosts, k)
		delete(vs.failedUpdatesHosts, k)
		host := vs.resources.Hosts.Get(k)
		if host == nil {
			continue
		}
		vs.removeHostFromCharts(host)
	}
}

func (vs *VSphere) cleanupVMs() {
	for k, v := range vs.failedUpdatesVms {
		if v < failedMax {
			continue
		}
		delete(vs.chartedHosts, k)
		delete(vs.failedUpdatesHosts, k)
		vm := vs.resources.VMs.Get(k)
		if vm == nil {
			continue
		}
		vs.removeVMFromCharts(vm)
	}
}

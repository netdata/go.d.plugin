package vsphere

func (vs *VSphere) collect() (map[string]int64, error) {
	mx := make(map[string]int64)

	vs.resLock.Lock()
	defer vs.resLock.Unlock()

	defer vs.updateCharts()
	defer vs.cleanupResources()

	err := vs.collectHosts(mx)
	if err != nil {
		return mx, err
	}

	err = vs.collectVMs(mx)
	if err != nil {
		return mx, err
	}

	return mx, nil
}

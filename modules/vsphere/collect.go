package vsphere

func (vs *VSphere) collect() (map[string]int64, error) {
	vs.collectionLock.Lock()
	defer vs.collectionLock.Unlock()

	defer vs.removeStale()

	mx := make(map[string]int64)
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

const (
	failedUpdatesLimit = 10
)

func (vs *VSphere) removeStale() {
	for k, v := range vs.discoveredHosts {
		if v < failedUpdatesLimit {
			continue
		}
		vs.removeFromCharts(k)
		delete(vs.charted, k)
		delete(vs.discoveredHosts, k)
	}
	for k, v := range vs.discoveredVMs {
		if v < failedUpdatesLimit {
			continue
		}
		vs.removeFromCharts(k)
		delete(vs.charted, k)
		delete(vs.discoveredVMs, k)
	}
}

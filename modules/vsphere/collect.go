package vsphere

func (vs *VSphere) collect() (map[string]int64, error) {
	mx := make(map[string]int64)

	vs.resLock.Lock()
	defer vs.resLock.Unlock()

	defer vs.updateCharts()
	defer vs.removeStale()

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
	failedMax = 10
)

func (vs *VSphere) removeStale() {
	for k, v := range vs.collectedHosts {
		if v < failedMax {
			continue
		}
		delete(vs.charted, k)
		delete(vs.collectedHosts, k)
		vs.removeFromCharts(k)
	}
	for k, v := range vs.collectedVMs {
		if v < failedMax {
			continue
		}
		delete(vs.charted, k)
		delete(vs.collectedVMs, k)
		vs.removeFromCharts(k)
	}
}

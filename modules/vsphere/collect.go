package vsphere

func (vs *VSphere) collect() (map[string]int64, error) {
	mx := make(map[string]int64)

	//err := vs.collectHosts(mx)
	//if err != nil {
	//	return mx, err
	//}
	//vs.updateHostsCharts()

	err := vs.collectVMs(mx)
	if err != nil {
		return mx, err
	}
	vs.updateVMsCharts()

	return mx, nil
}

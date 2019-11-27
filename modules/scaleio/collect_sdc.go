package scaleio

func (s *ScaleIO) collectSdcStats(mx *metrics, stats selectedStatistics) {
	mx.Sdc = make(map[string]sdcStatistics, len(stats.Sdc))
	for k, v := range stats.Sdc {
		var m sdcStatistics
		m.BW.set(
			calcBW(v.UserDataReadBwc),
			calcBW(v.UserDataWriteBwc),
		)
		m.IOPS.set(
			calcIOPS(v.UserDataReadBwc),
			calcIOPS(v.UserDataWriteBwc),
		)
		m.IOSize.set(
			calcIOSize(v.UserDataReadBwc),
			calcIOSize(v.UserDataWriteBwc),
		)
		m.MappedVolumes.Set(v.NumOfMappedVolumes)
		mx.Sdc[k] = m
	}
}

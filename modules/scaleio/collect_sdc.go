package scaleio

import "github.com/netdata/go.d.plugin/modules/scaleio/client"

func (s *ScaleIO) collectSdc(mx *metrics, ss client.SelectedStatistics) {
	mx.Sdc = make(map[string]sdcMetrics, len(ss.Sdc))

	for id, stats := range ss.Sdc {
		sdc, ok := s.discovered.sdc[id]
		if !ok {
			continue
		}
		var m sdcMetrics
		m.BW.set(
			calcBW(stats.UserDataReadBwc),
			calcBW(stats.UserDataWriteBwc),
		)
		m.IOPS.set(
			calcIOPS(stats.UserDataReadBwc),
			calcIOPS(stats.UserDataWriteBwc),
		)
		m.IOSize.set(
			calcIOSize(stats.UserDataReadBwc),
			calcIOSize(stats.UserDataWriteBwc),
		)
		m.MappedVolumes = stats.NumOfMappedVolumes
		m.MDMConnectionState = isSdcConnected(sdc.MdmConnectionState)

		mx.Sdc[id] = m
	}
}

func isSdcConnected(state string) bool {
	return state == "Connected"
}

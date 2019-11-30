package scaleio

import "github.com/netdata/go.d.plugin/modules/scaleio/client"

func (s *ScaleIO) collectSdc(mx *metrics, stats client.SelectedStatistics) {
	mx.Sdc = make(map[string]sdcStatistics, len(stats.Sdc))

	for k, v := range stats.Sdc {
		sdc, ok := s.discovered.sdc[k]
		if !ok {
			continue
		}
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
		m.MappedVolumes = v.NumOfMappedVolumes
		m.MDMConnectionState = isSdcConnected(sdc.MdmConnectionState)

		mx.Sdc[k] = m
	}
}

func isSdcConnected(state string) bool {
	return state == "Connected"
}

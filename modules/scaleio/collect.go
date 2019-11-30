package scaleio

import (
	"github.com/netdata/go.d.plugin/modules/scaleio/client"
	"github.com/netdata/go.d.plugin/pkg/stm"
)

const discoveryEvery = 5

func (s *ScaleIO) collect() (map[string]int64, error) {
	s.runs += 1
	if !s.lastDiscoveryOK || s.runs%discoveryEvery == 0 {
		if err := s.discovery(); err != nil {
			return nil, err
		}
	}

	stats, err := s.client.SelectedStatistics(query)
	if err != nil {
		return nil, err
	}

	var mx metrics
	s.collectSystemOverview(&mx, stats)
	s.collectSdc(&mx, stats)
	s.collectStoragePool(&mx, stats)
	s.updateCharts()
	return stm.ToMap(mx), nil
}

func (s *ScaleIO) discovery() error {
	ins, err := s.client.Instances()
	if err != nil {
		s.lastDiscoveryOK = false
		return err
	}

	s.discovered.pool = make(map[string]client.StoragePool, len(ins.StoragePoolList))
	for _, pool := range ins.StoragePoolList {
		s.discovered.pool[pool.ID] = pool
	}
	s.discovered.sdc = make(map[string]client.Sdc, len(ins.SdcList))
	for _, sdc := range ins.SdcList {
		s.discovered.sdc[sdc.ID] = sdc
	}
	s.lastDiscoveryOK = true
	return nil
}

package scaleio

import (
	"github.com/netdata/go.d.plugin/pkg/stm"
)

func (s *ScaleIO) collect() (map[string]int64, error) {
	var mx metrics
	err := s.collectSystemOverview(&mx)
	if err != nil {
		return nil, err
	}

	return stm.ToMap(mx), nil
}

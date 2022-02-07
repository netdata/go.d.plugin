package snmp

import (
	//"fmt"

	"github.com/gosnmp/gosnmp"
	"github.com/netdata/go.d.plugin/agent/module"
)

func (s *SNMP) collect() (map[string]int64, error) {
	collected := make(map[string]int64)

	for _, chart := range *s.Charts() {
		s.collectChart(collected, chart)
	}

	return collected, nil
}

func (s *SNMP) collectChart(collected map[string]int64, chart *module.Chart) {
	params := s.Config.SNMPClient

	for _, d := range chart.Dims {
		oid := []string{d.ID}
		result, err := params.Get(oid)

		if err != nil {
			s.Warningf("Cannot get SNMP data: %v", err)
			continue
		}
		a := gosnmp.ToBigInt(result.Variables[0].Value).Int64()
		collected[d.ID] = a
	}
}

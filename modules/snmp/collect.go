package snmp

import (
	"github.com/gosnmp/gosnmp"
)

func (s *SNMP) collect() (map[string]int64, error) {
	collected := make(map[string]int64)
	var all_oid []string

	//build oid chart
	for _, chart := range *s.Charts() {
		for _, d := range chart.Dims {
			all_oid = append(all_oid, d.ID)
		}
	}

	if err := s.collectChart(collected, all_oid); err != nil {
		return nil, err
	}

	return collected, nil
}

func (s *SNMP) collectChart(collected map[string]int64, oid_s []string) error {
	params := s.Config.SNMPClient
	if len(oid_s) > s.Config.MaxOIDs {
		if err := s.collectChart(collected, oid_s[s.Config.MaxOIDs:]); err != nil {
			return err
		}
		oid_s = oid_s[:s.Config.MaxOIDs]
	}

	result, err := params.Get(oid_s)

	if err != nil {
		s.Errorf("Cannot get SNMP data: %v", err)
		return err
	}
	for i, oid := range oid_s {
		collected[oid] = gosnmp.ToBigInt(result.Variables[i].Value).Int64()
	}
	return nil
}

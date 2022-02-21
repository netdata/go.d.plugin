package snmp

import (
	"github.com/gosnmp/gosnmp"
)

func (s *SNMP) collect() (map[string]int64, error) {
	collected := make(map[string]int64)
	var all_oid []string

	err := s.SNMPClient.Connect()
	if err != nil {
		s.Errorf("SNMP Connect fail: %v", err)
		return nil, err
	}
	defer s.SNMPClient.Close()

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

func (s *SNMP) collectChart(collected map[string]int64, OIDs []string) error {
	if len(OIDs) > s.Options.MaxOIDs {
		if err := s.collectChart(collected, OIDs[s.Options.MaxOIDs:]); err != nil {
			return err
		}
		OIDs = OIDs[:s.Options.MaxOIDs]
	}

	result, err := s.SNMPClient.Get(OIDs)

	if err != nil {
		s.Errorf("Cannot get SNMP data: %v", err)
		return err
	}

	for i, oid := range OIDs {
		//TODO: map oid to result.Variables[i].Name necessary??
		collected[oid] = gosnmp.ToBigInt(result.Variables[i].Value).Int64()
	}
	return nil
}

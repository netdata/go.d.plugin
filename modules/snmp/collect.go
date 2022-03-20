package snmp

import (
	"github.com/gosnmp/gosnmp"
)

func (s *SNMP) collect() (map[string]int64, error) {
	collected := make(map[string]int64)
	var oids []string

	for _, chart := range *s.Charts() {
		for _, d := range chart.Dims {
			oids = append(oids, d.ID)
		}
	}

	if err := s.collectChart(collected, oids); err != nil {
		return nil, err
	}

	return collected, nil
}

func (s *SNMP) collectChart(collected map[string]int64, oids []string) error {
	if len(oids) > s.Options.MaxOIDs {
		if err := s.collectChart(collected, oids[s.Options.MaxOIDs:]); err != nil {
			return err
		}
		oids = oids[:s.Options.MaxOIDs]
	}

	result, err := s.snmpHandler.Get(oids)

	if err != nil {
		s.Errorf("Cannot get SNMP data: %v", err)
		return err
	}

	for i, oid := range oids {
		switch result.Variables[i].Type {
		case gosnmp.NoSuchInstance, gosnmp.NoSuchObject:
			s.Debugf("Skipping OID %s, no such object", oid)
			continue
		default:
			collected[oid] = gosnmp.ToBigInt(result.Variables[i].Value).Int64()
		}
	}
	return nil
}

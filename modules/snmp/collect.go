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

	for i, end := 0, 0; i < len(oids); i += s.Options.MaxOIDs {
		if end = i + s.Options.MaxOIDs; end > len(oids) {
			end = len(oids)
		}
		if err := s.collectOIDs(collected, oids[i:end]); err != nil {
			return nil, err
		}
	}

	return collected, nil
}

func (s *SNMP) collectOIDs(collected map[string]int64, oids []string) error {
	result, err := s.snmpClient.Get(oids)
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

package snmp

import (
	"github.com/gosnmp/gosnmp"
)

func (s *SNMP) collect() (map[string]int64, error) {
	collected := make(map[string]int64)

	if err := s.collectOIDs(collected); err != nil {
		return nil, err
	}

	return collected, nil
}

func (s *SNMP) collectOIDs(collected map[string]int64) error {
	for i, end := 0, 0; i < len(s.oids); i += s.Options.MaxOIDs {
		if end = i + s.Options.MaxOIDs; end > len(s.oids) {
			end = len(s.oids)
		}

		oids := s.oids[i:end]
		resp, err := s.snmpClient.Get(oids)
		if err != nil {
			s.Errorf("cannot get SNMP data: %v", err)
			return err
		}

		for i, oid := range oids {
			switch resp.Variables[i].Type {
			case gosnmp.NoSuchInstance, gosnmp.NoSuchObject:
				s.Debugf("skipping OID %s, no such object", oid)
			default:
				collected[oid] = gosnmp.ToBigInt(resp.Variables[i].Value).Int64()
			}
		}
	}

	return nil
}

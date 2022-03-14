package snmp

import (
	"time"

	gosnmp "github.com/gosnmp/gosnmp"
)

var snmpHandler = gosnmp.NewHandler

func (s *SNMP) initSNMPClient() bool {
	snmpClient := snmpHandler()

	//Default SNMP connection params
	snmpClient.SetTarget(s.Name)
	snmpClient.SetPort(uint16(s.Options.Port))
	snmpClient.SetMaxOids(s.Options.MaxOIDs)
	snmpClient.SetLogger(gosnmp.NewLogger(s.Logger))
	snmpClient.SetTimeout(time.Duration(s.Options.Timeout) * time.Second)

	switch s.Options.Version {
	case 1:
		snmpClient.SetCommunity(*s.Community)
		snmpClient.SetVersion(gosnmp.Version1)

	case 2:
		snmpClient.SetCommunity(*s.Community)
		snmpClient.SetVersion(gosnmp.Version2c)

	case 3:
		snmpClient.SetVersion(gosnmp.Version3)
		snmpClient.SetSecurityModel(gosnmp.UserSecurityModel)
		snmpClient.SetMsgFlags(gosnmp.SnmpV3MsgFlags(s.User.Level))
		snmpClient.SetSecurityParameters(&gosnmp.UsmSecurityParameters{
			UserName:                 s.User.Name,
			AuthenticationProtocol:   gosnmp.SnmpV3AuthProtocol(s.User.AuthProto),
			AuthenticationPassphrase: s.User.AuthKey,
			PrivacyProtocol:          gosnmp.SnmpV3PrivProtocol(s.User.PrivProto),
			PrivacyPassphrase:        s.User.PrivKey,
		})

	default:
		s.Errorf("invalid SNMP version: %d", s.Options.Version)
		return false
	}

	err := snmpClient.Connect()
	if err != nil {
		s.Errorf("SNMP Connect fail: %v", err)
		return false
	}
	s.SNMPClient = snmpClient

	return true
}

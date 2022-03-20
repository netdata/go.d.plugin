package snmp

import (
	"fmt"
	"time"

	"github.com/gosnmp/gosnmp"
)

var newSNMPClient = gosnmp.NewHandler

func (s SNMP) initSNMPClient() (gosnmp.Handler, error) {
	snmpClient := newSNMPClient()

	snmpClient.SetTarget(s.Hostname)
	snmpClient.SetPort(uint16(s.Options.Port))
	snmpClient.SetMaxOids(s.Options.MaxOIDs)
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
		snmpClient.SetMsgFlags(generateSNMPv3MsgFlags(s.User.Level))
		snmpClient.SetSecurityParameters(&gosnmp.UsmSecurityParameters{
			UserName:                 s.User.Name,
			AuthenticationProtocol:   gosnmp.SnmpV3AuthProtocol(s.User.AuthProto),
			AuthenticationPassphrase: s.User.AuthKey,
			PrivacyProtocol:          gosnmp.SnmpV3PrivProtocol(s.User.PrivProto),
			PrivacyPassphrase:        s.User.PrivKey,
		})
	default:
		return nil, fmt.Errorf("invalid SNMP version: %d", s.Options.Version)

	}

	return snmpClient, nil
}

func generateSNMPv3MsgFlags(level int) gosnmp.SnmpV3MsgFlags {
	flag := gosnmp.NoAuthNoPriv
	switch level {
	case 1:
		flag = gosnmp.NoAuthNoPriv
	case 2:
		flag = gosnmp.AuthNoPriv
	case 3:
		flag = gosnmp.AuthPriv
	}
	return flag
}

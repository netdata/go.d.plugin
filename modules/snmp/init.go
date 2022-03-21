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
		snmpClient.SetMsgFlags(safeConvertSNMPv3UserLevel(s.User.Level))
		snmpClient.SetSecurityParameters(&gosnmp.UsmSecurityParameters{
			UserName:                 s.User.Name,
			AuthenticationProtocol:   safeConvertSNMPv3UserAuthProtocol(s.User.AuthProto),
			AuthenticationPassphrase: s.User.AuthKey,
			PrivacyProtocol:          safeConvertSNMPV3UserPrivProtocol(s.User.PrivProto),
			PrivacyPassphrase:        s.User.PrivKey,
		})
	default:
		return nil, fmt.Errorf("invalid SNMP version: %d", s.Options.Version)

	}

	return snmpClient, nil
}

func safeConvertSNMPv3UserLevel(level string) gosnmp.SnmpV3MsgFlags {
	v, _ := convertSNMPv3UserLevel(level)
	return v
}

func convertSNMPv3UserLevel(level string) (gosnmp.SnmpV3MsgFlags, error) {
	switch level {
	case "1", "noAuthNoPriv":
		return gosnmp.NoAuthNoPriv, nil
	case "2", "authNoPriv":
		return gosnmp.AuthNoPriv, nil
	case "3", "authPriv":
		return gosnmp.AuthPriv, nil
	default:
		return gosnmp.NoAuthNoPriv, fmt.Errorf("invalid snmpv3 user level value (%s)", level)
	}
}

func safeConvertSNMPv3UserAuthProtocol(protocol string) gosnmp.SnmpV3AuthProtocol {
	v, _ := convertSNMPv3UserAuthProtocol(protocol)
	return v
}

func convertSNMPv3UserAuthProtocol(protocol string) (gosnmp.SnmpV3AuthProtocol, error) {
	switch protocol {
	case "1", "none":
		return gosnmp.NoAuth, nil
	case "2", "md5":
		return gosnmp.MD5, nil
	case "3", "sha":
		return gosnmp.SHA, nil
	case "4", "sha224":
		return gosnmp.SHA224, nil
	case "5", "sha256":
		return gosnmp.SHA256, nil
	case "6", "sha384":
		return gosnmp.SHA384, nil
	case "7", "sha512":
		return gosnmp.SHA512, nil
	default:
		return gosnmp.NoAuth, fmt.Errorf("invalid snmpv3 user auth protocol value (%s)", protocol)
	}
}

func safeConvertSNMPV3UserPrivProtocol(protocol string) gosnmp.SnmpV3PrivProtocol {
	v, _ := convertSNMPV3UserPrivProtocol(protocol)
	return v
}

func convertSNMPV3UserPrivProtocol(protocol string) (gosnmp.SnmpV3PrivProtocol, error) {
	switch protocol {
	case "1", "noPriv":
		return gosnmp.NoPriv, nil
	case "2", "des":
		return gosnmp.DES, nil
	case "3", "aes":
		return gosnmp.AES, nil
	case "4", "aes192":
		return gosnmp.AES192, nil
	case "5", "aes256":
		return gosnmp.AES256, nil
	case "6", "aes192c":
		return gosnmp.AES192C, nil
	case "7", "aes256c":
		return gosnmp.AES256C, nil
	default:
		return gosnmp.NoPriv, fmt.Errorf("invalid snmpv3 user priv protocol value (%s)", protocol)
	}
}

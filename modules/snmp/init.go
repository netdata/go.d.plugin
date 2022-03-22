package snmp

import (
	"errors"
	"fmt"
	"time"

	"github.com/gosnmp/gosnmp"
)

var newSNMPClient = gosnmp.NewHandler

func (s SNMP) validateConfig() error {
	if len(s.ChartsInput) == 0 {
		return errors.New("'user.name' is required when using SNMPv3 but not set")
	}

	if s.Options.Version == int(gosnmp.Version3) {
		if s.User.Name == "" {
			return errors.New("'user.name' is required when using SNMPv3 but not set")
		}
		if _, err := parseSNMPv3SecurityLevel(s.User.SecurityLevel); err != nil {
			return err
		}
		if _, err := parseSNMPv3AuthProtocol(s.User.AuthProto); err != nil {
			return err
		}
		if _, err := parseSNMPv3PrivProtocol(s.User.PrivProto); err != nil {
			return err
		}
	}

	return nil
}

func (s SNMP) initSNMPClient() (gosnmp.Handler, error) {
	snmpClient := newSNMPClient()

	if s.Hostname == "" {
		s.Warningf("'hostname' not set, using the default value: '%s'", defaultHostname)
		snmpClient.SetTarget(defaultHostname)
	} else {
		snmpClient.SetTarget(s.Hostname)
	}
	if s.Options.Port <= 0 || s.Options.Port > 65535 {
		s.Warningf("'options.port' is invalid, changing to the default value: '%d' => '%d'", s.Options.Port, defaultPort)
		snmpClient.SetPort(defaultPort)
	} else {
		snmpClient.SetPort(uint16(s.Options.Port))
	}
	if s.Options.Retries < 1 || s.Options.Retries > 10 {
		s.Warningf("'options.retries' is invalid, changing to the default value: '%d' => '%d'", s.Options.Retries, defaultRetries)
		snmpClient.SetRetries(defaultRetries)
	} else {
		snmpClient.SetRetries(s.Options.Retries)
	}
	if s.Options.Timeout < 1 {
		s.Warningf("'options.timeout' is invalid, changing to the default value: '%d' => '%d'", s.Options.Timeout, defaultTimeout)
		snmpClient.SetTimeout(defaultTimeout * time.Second)
	} else {
		snmpClient.SetTimeout(time.Duration(s.Options.Timeout) * time.Second)
	}
	if s.Options.MaxOIDs < 1 {
		s.Warningf("'options.max_request_size' is invalid, changing to the default value: '%d' => '%d'", s.Options.MaxOIDs, defaultMaxOIDs)
		snmpClient.SetMaxOids(defaultMaxOIDs)
	} else {
		snmpClient.SetMaxOids(s.Options.MaxOIDs)
	}

	snmpVersion := s.Options.Version
	if snmpVersion < int(gosnmp.Version1) || snmpVersion > int(gosnmp.Version3) {
		s.Warningf("'options.version' is invalid, changing to the default value: '%d' => '%d'",
			s.Options.Version, defaultVersion)
		snmpVersion = defaultVersion
	}
	community := s.Community
	if community == "" && (snmpVersion == int(gosnmp.Version1) || snmpVersion == int(gosnmp.Version2c)) {
		s.Warningf("'community' not set, using the default value: '%s'", defaultCommunity)
		community = defaultCommunity
	}

	switch snmpVersion {
	case 1:
		snmpClient.SetCommunity(community)
		snmpClient.SetVersion(gosnmp.Version1)
	case 2:
		snmpClient.SetCommunity(community)
		snmpClient.SetVersion(gosnmp.Version2c)
	case 3:
		snmpClient.SetVersion(gosnmp.Version3)
		snmpClient.SetSecurityModel(gosnmp.UserSecurityModel)
		snmpClient.SetMsgFlags(safeParseSNMPv3SecurityLevel(s.User.SecurityLevel))
		snmpClient.SetSecurityParameters(&gosnmp.UsmSecurityParameters{
			UserName:                 s.User.Name,
			AuthenticationProtocol:   safeParseSNMPv3AuthProtocol(s.User.AuthProto),
			AuthenticationPassphrase: s.User.AuthKey,
			PrivacyProtocol:          safeParseSNMPv3PrivProtocol(s.User.PrivProto),
			PrivacyPassphrase:        s.User.PrivKey,
		})
	default:
		return nil, fmt.Errorf("invalid SNMP version: %d", s.Options.Version)
	}

	return snmpClient, nil
}

func (s SNMP) initOIDs() (oids []string) {
	for _, c := range *s.charts {
		for _, d := range c.Dims {
			oids = append(oids, d.ID)
		}
	}
	return oids
}

func safeParseSNMPv3SecurityLevel(level string) gosnmp.SnmpV3MsgFlags {
	v, _ := parseSNMPv3SecurityLevel(level)
	return v
}

func parseSNMPv3SecurityLevel(level string) (gosnmp.SnmpV3MsgFlags, error) {
	switch level {
	case "1", "none", "noAuthNoPriv", "":
		return gosnmp.NoAuthNoPriv, nil
	case "2", "authNoPriv":
		return gosnmp.AuthNoPriv, nil
	case "3", "authPriv":
		return gosnmp.AuthPriv, nil
	default:
		return gosnmp.NoAuthNoPriv, fmt.Errorf("invalid snmpv3 user security level value (%s)", level)
	}
}

func safeParseSNMPv3AuthProtocol(protocol string) gosnmp.SnmpV3AuthProtocol {
	v, _ := parseSNMPv3AuthProtocol(protocol)
	return v
}

func parseSNMPv3AuthProtocol(protocol string) (gosnmp.SnmpV3AuthProtocol, error) {
	switch protocol {
	case "1", "none", "noAuth", "":
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

func safeParseSNMPv3PrivProtocol(protocol string) gosnmp.SnmpV3PrivProtocol {
	v, _ := parseSNMPv3PrivProtocol(protocol)
	return v
}

func parseSNMPv3PrivProtocol(protocol string) (gosnmp.SnmpV3PrivProtocol, error) {
	switch protocol {
	case "1", "none", "noPriv", "":
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

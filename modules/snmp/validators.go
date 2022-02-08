package snmp

import (
	"errors"
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
)

func appendError(err error, msg string) error {
	var e error
	if err == nil {
		if msg == "" {
			return nil
		}
		e = errors.New(msg)
	} else {
		if msg == "" {
			return err
		}
		e = fmt.Errorf("%w %s", err, msg)
	}
	return e
}

func (d Dimension) validateConfig() error {
	var err error
	err = nil
	if d.Name == "" {
		err = appendError(err, "invalid or missing value: dimension.name;")
	}
	if d.OID == "" {
		err = appendError(err, "missing value: dimension.oid;")
	}
	if d.Algorithm == "" ||
		(d.Algorithm != string(module.Incremental) &&
			d.Algorithm != string(module.PercentOfIncremental) &&
			d.Algorithm != string(module.PercentOfAbsolute)) {
		err = appendError(err, "invalid or missing value: dimension.algorithm;")
	}
	if d.Multiplier == 0 {
		err = appendError(err, "integer set to 0: dimension.multiplier;")
	}
	if d.Divisor == 0 {
		err = appendError(err, "integer set to 0: dimension.divisor;")
	}
	return err
}

func (u User) validateConfig() error {
	var err error
	err = nil
	if u.Name == "" {
		err = appendError(err, "missing value: user.name;")
	}
	if u.Level < 1 || u.Level > 3 {
		err = appendError(err, "invalid range of value: user.level;")
	}
	if u.PrivProto < 1 || u.PrivProto > 2 {
		err = appendError(err, "invalid range of value: user.priv_proto;")
	}
	if u.AuthProto < 1 || u.AuthProto > 3 {
		err = appendError(err, "invalid range of value: user.auth_proto;")
	}
	return err
}

func (c ChartsConfig) validateConfig() error {
	var err error
	err = nil
	if c.Title == "" {
		err = appendError(err, "missing value: charts.title;")
	}
	if c.Priority < 1 {
		err = appendError(err, "invalid value: charts.priority;")
	}
	if c.Dimensions != nil {
		for _, d := range c.Dimensions {
			if e := d.validateConfig(); e != nil {
				err = appendError(err, e.Error())
			}
		}
	}

	return err
}

func (o Options) validateConfig() error {
	var err error
	err = nil
	if o.Port <= 0 && o.Port > 65535 {
		err = appendError(err, "invalid range of value: options.port;")
	}
	if o.Version < 1 || o.Version > 3 {
		err = appendError(err, "invalid range of value: options.versions;")
	}
	if o.Retries < 0 || o.Retries > 100 {
		err = appendError(err, "invalid range of value: options.retries;")
	}
	if o.Timeout < 0 {
		err = appendError(err, "invalid value: options.timeout;")
	}

	return err
}

func (s SNMP) validateConfig() error {
	var err error
	err = nil
	if s.Options != nil {
		if e := s.Options.validateConfig(); e != nil {
			err = appendError(err, e.Error())
		}
		if s.Options.Version == 3 {
			if s.User == nil {
				err = appendError(err, "SNMP v3 missing user credentials;")
			}
		} else {
			if s.Community == nil {
				err = appendError(err, "SNMP v1/2 missing community value;")
			}
		}
	}
	if s.User != nil {
		if e := s.User.validateConfig(); e != nil {
			err = appendError(err, e.Error())
		}
	}
	if s.Settings != nil {
		if e := s.Settings[0].validateConfig(); e != nil {
			err = appendError(err, e.Error())
		}
	}

	return err
}

package x509check

import (
	"crypto/x509"
	"time"
)

func (x *X509Check) collect() (map[string]int64, error) {
	certs, err := x.Gather()

	if err != nil {
		x.Error(err)
		return nil, nil
	}

	if len(certs) == 0 {
		x.Error("no certificate was provided by '%s'", x.Config.Source)
		return nil, nil
	}

	mx := map[string]int64{
		"expiry":                         calcExpiry(certs),
		"days_until_expiration_warning":  x.DaysUntilWarn,
		"days_until_expiration_critical": x.DaysUntilCrit,
	}

	return mx, nil
}

func calcExpiry(certs []*x509.Certificate) int64 {
	now := time.Now()
	return int64(certs[0].NotAfter.Sub(now).Seconds())
}

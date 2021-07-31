package mongo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth_connectionString(t *testing.T) {
	auth := Auth{
		SSL: SSL{
			Ssl:              true,
			SslCertReqs:      true,
			SslCaCerts:       "certs.file",
			SslCrlFile:       "certs.file",
			SslCertfile:      "certs.file",
			SslKeyfile:       "certs.file",
			SslPemPassphrase: "pass",
		},
		Host:   "localhost",
		Port:   27017,
		Authdb: "admin",
		User:   "user",
		Pass:   "pass",
	}
	connectionString := auth.connectionString()
	assert.Equal(
		t,
		"mongodb://user:pass@localhost:27017/authSource=admin&ssl=true&tlsCAFile=certs.file&tlsCertificateKeyFile=certs.file&tlsCertificateKeyFilePassword=pass",
		connectionString,
	)
}

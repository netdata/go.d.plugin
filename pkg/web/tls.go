package web

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
)

// ClientTLSConfig represents the standard client TLS config.
type ClientTLSConfig struct {
	TLSCA              string `yaml:"tls_ca"`
	TLSCert            string `yaml:"tls_cert"`
	TLSKey             string `yaml:"tls_key"`
	InsecureSkipVerify bool   `yaml:"tls_skip_verify"`
}

// NewTLSConfig returns a tls.Config, may be nil without error if TLS is not
// configured.
func NewTLSConfig(config ClientTLSConfig) (*tls.Config, error) {
	if config.TLSCA == "" && config.TLSKey == "" && config.TLSCert == "" && !config.InsecureSkipVerify {
		return nil, nil
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: config.InsecureSkipVerify,
		Renegotiation:      tls.RenegotiateNever,
	}

	if config.TLSCA != "" {
		pool, err := makeCertPool([]string{config.TLSCA})
		if err != nil {
			return nil, err
		}
		tlsConfig.RootCAs = pool
	}

	if config.TLSCert != "" && config.TLSKey != "" {
		if err := loadCertificate(tlsConfig, config.TLSCert, config.TLSKey); err != nil {
			return nil, err
		}
	}

	return tlsConfig, nil
}

func makeCertPool(certFiles []string) (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	for _, certFile := range certFiles {
		pem, err := ioutil.ReadFile(certFile)
		if err != nil {
			return nil, fmt.Errorf("could not read certificate %q: %v", certFile, err)
		}
		if !pool.AppendCertsFromPEM(pem) {
			return nil, fmt.Errorf("could not parse any PEM certificates %q: %v", certFile, err)
		}
	}
	return pool, nil
}

func loadCertificate(config *tls.Config, certFile, keyFile string) error {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("could not load keypair %s:%s: %v", certFile, keyFile, err)
	}

	config.Certificates = []tls.Certificate{cert}
	return nil
}

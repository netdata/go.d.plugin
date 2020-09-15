package tlscfg

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
)

// TLSConfig represents the standard client TLS configuration.
type TLSConfig struct {
	TLSCA              string `yaml:"tls_ca"`
	TLSCert            string `yaml:"tls_cert"`
	TLSKey             string `yaml:"tls_key"`
	InsecureSkipVerify bool   `yaml:"tls_skip_verify"`
}

// NewTLSConfig creates a tls.Config, may be nil without an error if TLS is not configured.
func NewTLSConfig(cfg TLSConfig) (*tls.Config, error) {
	if cfg.TLSCA == "" && cfg.TLSKey == "" && cfg.TLSCert == "" && !cfg.InsecureSkipVerify {
		return nil, nil
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: cfg.InsecureSkipVerify,
		Renegotiation:      tls.RenegotiateNever,
	}

	if cfg.TLSCA != "" {
		pool, err := makeCertPool([]string{cfg.TLSCA})
		if err != nil {
			return nil, err
		}
		tlsConfig.RootCAs = pool
	}

	if cfg.TLSCert != "" && cfg.TLSKey != "" {
		if err := loadCertificate(tlsConfig, cfg.TLSCert, cfg.TLSKey); err != nil {
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

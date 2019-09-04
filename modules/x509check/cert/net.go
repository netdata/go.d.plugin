package cert

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/url"
	"time"
)

func NewNet(url *url.URL, tlsConfig *tls.Config, timeout time.Duration) *Net {
	return &Net{
		URL:       url,
		TLSConfig: tlsConfig,
		Timeout:   timeout,
	}
}

type Net struct {
	URL       *url.URL
	TLSConfig *tls.Config
	Timeout   time.Duration
}

func (n Net) Gather() ([]*x509.Certificate, error) {
	ipConn, err := net.DialTimeout(n.URL.Scheme, n.URL.Host, n.Timeout)
	if err != nil {
		return nil, fmt.Errorf("error on dial to '%s' : %v", n.URL, err)
	}
	defer ipConn.Close()

	conn := tls.Client(ipConn, n.TLSConfig.Clone())
	defer conn.Close()

	if err := conn.Handshake(); err != nil {
		return nil, fmt.Errorf("error on ssl handshake with '%s' : %v", n.URL, err)
	}

	certs := conn.ConnectionState().PeerCertificates
	return certs, nil
}

package cert

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/smtp"
	"net/url"
	"time"
)

func NewSMTP(url *url.URL, tlsConfig *tls.Config, timeout time.Duration) *SMTP {
	return &SMTP{
		URL:       url,
		TLSConfig: tlsConfig,
		Timeout:   timeout,
	}
}

type SMTP struct {
	URL       *url.URL
	TLSConfig *tls.Config
	Timeout   time.Duration
}

func (s SMTP) Gather() ([]*x509.Certificate, error) {
	ipConn, err := net.DialTimeout(s.URL.Scheme, s.URL.Host, s.Timeout)
	if err != nil {
		return nil, fmt.Errorf("error on dial to '%s' : %v", s.URL, err)
	}
	defer ipConn.Close()

	host, _, _ := net.SplitHostPort(s.URL.Host)
	smtpClient, err := smtp.NewClient(ipConn, host)
	if err != nil {
		return nil, fmt.Errorf("error on creating smtp client : %v", err)
	}
	defer smtpClient.Quit()

	err = smtpClient.StartTLS(s.TLSConfig.Clone())
	if err != nil {
		return nil, fmt.Errorf("error on startTLS with '%s' : %v", s.URL, err)
	}

	conn, ok := smtpClient.TLSConnectionState()
	if !ok {
		return nil, fmt.Errorf("startTLS didn't succeed")
	}
	return conn.PeerCertificates, nil
}

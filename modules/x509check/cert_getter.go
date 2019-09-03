package x509check

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/smtp"
	"net/url"
	"time"

	"github.com/netdata/go.d.plugin/pkg/web"
)

var supportedSchemes = []string{
	"file",
	"https",
	"tcp",
	"tcp4",
	"tcp6",
	"udp",
	"udp4",
	"udp6",
	"smtp",
}

func createTLSConfig(config web.ClientTLSConfig, hostname string) (*tls.Config, error) {
	tlsCfg, err := web.NewTLSConfig(config)
	if err != nil {
		return nil, err
	}

	if tlsCfg == nil {
		tlsCfg = &tls.Config{}
	}
	tlsCfg.ServerName = hostname
	return tlsCfg, nil
}

func newCertGetter(config Config) (certGetter, error) {
	if config.Source == "" {
		return nil, errors.New("'source' parameter is mandatory, but it's not set")
	}

	u, err := url.Parse(config.Source)

	if err != nil {
		return nil, fmt.Errorf("error on parsing source : %v", err)
	}

	switch u.Scheme {
	case "file":
		return newFileCertGetter(u.Path), nil
	case "https":
		u.Scheme = "tcp"
		fallthrough
	case "udp", "udp4", "udp6", "tcp", "tcp4", "tcp6":
		tlsCfg, err := createTLSConfig(config.ClientTLSConfig, u.Hostname())
		if err != nil {
			return nil, fmt.Errorf("error on creating tls config : %v", err)
		}
		return newURLCertGetter(u, tlsCfg, config.Timeout.Duration), nil
	case "smtp":
		u.Scheme = "tcp"
		tlsCfg, err := createTLSConfig(config.ClientTLSConfig, u.Hostname())
		if err != nil {
			return nil, fmt.Errorf("error on creating tls config : %v", err)
		}
		return newSMTPCertGetter(u, tlsCfg, config.Timeout.Duration), nil

	}

	return nil, fmt.Errorf("unsupported scheme in '%s', supported schemes : %v", u, supportedSchemes)
}

type certGetter interface {
	getCert() ([]*x509.Certificate, error)
}

func newFileCertGetter(path string) *fileCertGetter {
	return &fileCertGetter{path: path}
}

type fileCertGetter struct {
	path string
}

func (fg fileCertGetter) getCert() ([]*x509.Certificate, error) {
	content, err := ioutil.ReadFile(fg.path)
	if err != nil {
		return nil, fmt.Errorf("error on reading '%s' : %v", fg.path, err)
	}

	block, _ := pem.Decode(content)
	if block == nil {
		return nil, fmt.Errorf("error on decoding '%s' : %v", fg.path, err)
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error on parsing certigicate '%s' : %v", fg.path, err)
	}

	return []*x509.Certificate{cert}, nil
}

func newURLCertGetter(url *url.URL, tlsCfg *tls.Config, timeout time.Duration) *urlCertGetter {
	return &urlCertGetter{
		url:     url,
		tlsCfg:  tlsCfg,
		timeout: timeout,
	}
}

type urlCertGetter struct {
	url     *url.URL
	tlsCfg  *tls.Config
	timeout time.Duration
}

func (ug urlCertGetter) getCert() ([]*x509.Certificate, error) {
	ipConn, err := net.DialTimeout(ug.url.Scheme, ug.url.Host, ug.timeout)
	if err != nil {
		return nil, fmt.Errorf("error on dial to '%s' : %v", ug.url, err)
	}

	defer ipConn.Close()

	conn := tls.Client(ipConn, ug.tlsCfg.Clone())

	defer conn.Close()

	if err := conn.Handshake(); err != nil {
		return nil, fmt.Errorf("error on ssl handshake with '%s' : %v", ug.url, err)
	}

	certs := conn.ConnectionState().PeerCertificates

	return certs, nil
}

func newSMTPCertGetter(url *url.URL, tlsCfg *tls.Config, timeout time.Duration) *smtpCertGetter {
	return &smtpCertGetter{
		url:     url,
		tlsCfg:  tlsCfg,
		timeout: timeout,
	}
}

type smtpCertGetter struct {
	url     *url.URL
	tlsCfg  *tls.Config
	timeout time.Duration
}

func (sg smtpCertGetter) getCert() ([]*x509.Certificate, error) {
	ipConn, err := net.DialTimeout(sg.url.Scheme, sg.url.Host, sg.timeout)
	if err != nil {
		return nil, fmt.Errorf("error on dial to '%s' : %v", sg.url, err)
	}
	defer ipConn.Close()

	host, _, _ := net.SplitHostPort(sg.url.Host)

	smtpClient, err := smtp.NewClient(ipConn, host)
	if err != nil {
		return nil, fmt.Errorf("error on creating smtp client : %v", err)
	}
	defer smtpClient.Quit()

	err = smtpClient.StartTLS(sg.tlsCfg.Clone())
	if err != nil {
		return nil, fmt.Errorf("error on startTLS with '%s' : %v", sg.url, err)
	}

	conn, ok := smtpClient.TLSConnectionState()
	if !ok {
		return nil, fmt.Errorf("startTLS didn't succedd")
	}

	return conn.PeerCertificates, nil
}

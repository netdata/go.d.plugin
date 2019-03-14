package x509check

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"time"

	"github.com/netdata/go.d.plugin/pkg/web"
)

var supportedSchemes = []string{
	//"file",
	"https",
	"tcp",
	"tcp4",
	"tcp6",
	"udp",
	"udp4",
	"udp6",
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
	// TODO: not tested
	//case "file":
	//	return newFileCertGetter(u.Path), nil
	case "https":
		u.Scheme = "tcp"
		fallthrough
	case "udp", "udp4", "udp6", "tcp", "tcp4", "tcp6":
		tlsCfg, err := web.NewTLSConfig(config.ClientTLSConfig)

		if err != nil {
			return nil, fmt.Errorf("error on creating tls config : %v", err)
		}

		if tlsCfg == nil {
			tlsCfg = &tls.Config{}
		}

		tlsCfg.ServerName = u.Hostname()

		return newURLCertGetter(u, tlsCfg, config.Timeout.Duration), nil
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

	conn := tls.Client(ipConn, ug.tlsCfg)

	defer conn.Close()

	if err := conn.Handshake(); err != nil {
		return nil, fmt.Errorf("error on ssl handshake with '%s' : %v", ug.url, err)
	}

	certs := conn.ConnectionState().PeerCertificates

	return certs, nil
}

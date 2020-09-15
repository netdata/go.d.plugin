package web

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/netdata/go.d.plugin/pkg/tlscfg"
)

// ErrRedirectAttempted indicates that a redirect occurred.
var ErrRedirectAttempted = errors.New("redirect")

// Client is a struct that contains the fields that are needed fore creating HTTPClient.
type Client struct {
	Timeout           Duration `yaml:"timeout"`
	NotFollowRedirect bool     `yaml:"not_follow_redirects"`
	ProxyURL          string   `yaml:"proxy_url"`
	tlscfg.TLSConfig  `yaml:",inline"`
}

// NewHTTPClient creates new HTTPClient.
func NewHTTPClient(client Client) (*http.Client, error) {
	tlsConfig, err := tlscfg.NewTLSConfig(client.TLSConfig)
	if err != nil {
		return nil, fmt.Errorf("error on creating TLS config : %v", err)
	}

	transport := &http.Transport{
		Proxy:               proxyFunc(client.ProxyURL),
		TLSClientConfig:     tlsConfig,
		DialContext:         (&net.Dialer{Timeout: client.Timeout.Duration}).DialContext,
		TLSHandshakeTimeout: client.Timeout.Duration,
	}

	return &http.Client{
		Timeout:       client.Timeout.Duration,
		Transport:     transport,
		CheckRedirect: redirectFunc(client.NotFollowRedirect),
	}, nil
}

func redirectFunc(notFollowRedirect bool) func(req *http.Request, via []*http.Request) error {
	if notFollowRedirect {
		return func(_ *http.Request, _ []*http.Request) error { return ErrRedirectAttempted }
	}
	return nil
}

func proxyFunc(rawProxyURL string) func(r *http.Request) (*url.URL, error) {
	if rawProxyURL == "" {
		return http.ProxyFromEnvironment
	}

	proxyURL, err := url.Parse(rawProxyURL)
	if err != nil {
		return func(_ *http.Request) (*url.URL, error) { return nil, fmt.Errorf("invalid proxy: %v", err) }
	}
	return func(r *http.Request) (*url.URL, error) { return proxyURL, nil }
}

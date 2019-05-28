package web

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
)

// ErrRedirectAttempted indicates that a redirect occurred.
var ErrRedirectAttempted = errors.New("redirect")

// Client is a struct that contains the fields that are needed fore creating HTTPClient.
type Client struct {
	Timeout           Duration `yaml:"timeout"`              // default is zero (no timeout) must be tuned by modules
	NotFollowRedirect bool     `yaml:"not_follow_redirects"` // default is follow
	ProxyURL          string   `yaml:"proxy_url"`
	ClientTLSConfig   `yaml:",inline"`
}

// NewHTTPClient creates new HTTPClient.
func NewHTTPClient(client Client) (*http.Client, error) {
	tlsConfig, err := NewTLSConfig(client.ClientTLSConfig)

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

func proxyFunc(proxyurl string) func(r *http.Request) (*url.URL, error) {
	if proxyurl == "" {
		return http.ProxyFromEnvironment
	}

	proxyURL, err := url.Parse(proxyurl)
	if err != nil {
		return func(_ *http.Request) (*url.URL, error) { return nil, fmt.Errorf("invalid proxy: %v", err) }
	}

	return func(r *http.Request) (*url.URL, error) { return proxyURL, nil }
}

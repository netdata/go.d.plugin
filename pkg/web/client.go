package web

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// HTTPClient is the interface that wraps the Do method.
type HTTPClient interface {
	Do(r *http.Request) (*http.Response, error)
}

type clientFunc func(r *http.Request) (*http.Response, error)

// Do calls f(r).
func (f clientFunc) Do(r *http.Request) (*http.Response, error) {
	return f(r)
}

// Client is a struct that contains the fields that are needed fore creating HTTPClient.
type Client struct {
	Timeout           Duration `yaml:"timeout"`              // default is zero (no timeout) must be tuned by modules
	NotFollowRedirect bool     `yaml:"not_follow_redirects"` // default is follow
	SkipVerify        bool     `yaml:"skip_tls_verify"`      // default is verify
	ProxyURL          string   `yaml:"proxy_url"`
	ProxyUsername     string   `yaml:"proxy_username"`
	ProxyPassword     string   `yaml:"proxy_password"`
}

// NewHTTPClient creates new HTTPClient.
func NewHTTPClient(client Client) HTTPClient {
	if client.ProxyUsername == "" || client.ProxyPassword == "" {
		return newHTTPClient(client)
	}
	return newHTTPClientProxyAuth(client)
}

// TODO: TLSClientConfig
func newHTTPClient(client Client) HTTPClient {
	httpClient := &http.Client{
		Timeout: client.Timeout.Duration,
		Transport: &http.Transport{
			Proxy:           proxyFunc(client.ProxyURL),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: client.SkipVerify},
		},
	}

	if client.NotFollowRedirect {
		httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error { return errors.New("redirect") }
	}

	return httpClient
}

func newHTTPClientProxyAuth(client Client) HTTPClient {
	auth := base64.StdEncoding.EncodeToString([]byte(client.ProxyUsername + ":" + client.ProxyPassword))
	httpClient := newHTTPClient(client)

	f := func(req *http.Request) (*http.Response, error) {
		req.Header.Set("Proxy-Authorization", "Basic "+auth)

		return httpClient.Do(req)
	}

	return clientFunc(f)
}

func proxyFunc(rawProxyURL string) func(r *http.Request) (*url.URL, error) {
	if rawProxyURL == "" {
		return http.ProxyFromEnvironment
	}

	proxyURL, err := url.Parse(rawProxyURL)
	if err != nil {
		return func(r *http.Request) (*url.URL, error) { return nil, fmt.Errorf("invalid proxy: %s", err) }
	}

	return func(r *http.Request) (*url.URL, error) { return proxyURL, nil }
}

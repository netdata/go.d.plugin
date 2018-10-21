package web

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/l2isbad/go.d.plugin/modules/helpers/utils"
)

// Client is the interface that wraps the Do method.
type Client interface {
	Do(r *http.Request) (*http.Response, error)
}

type clientFunc func(r *http.Request) (*http.Response, error)

// Do calls f(r).
func (f clientFunc) Do(r *http.Request) (*http.Response, error) {
	return f(r)
}

// RawClient is a struct that contains the fields that are needed to create Client.
type RawClient struct {
	Timeout           utils.Duration `yaml:"timeout"`              // default is zero (no timeout) must be tuned by modules
	NotFollowRedirect bool           `yaml:"not_follow_redirects"` // default is follow
	SkipVerify        bool           `yaml:"skip_tls_verify"`      // default is verify
	ProxyURL          string         `yaml:"proxy_url"`
	ProxyUsername     string         `yaml:"proxy_username"`
	ProxyPassword     string         `yaml:"proxy_password"`
}

// CreateHTTPClient returns new Client.
func (r RawClient) CreateHTTPClient() Client {
	if r.ProxyUsername == "" || r.ProxyPassword == "" {
		return r.create()
	}
	return r.createWithProxyAuth()
}

// TODO: TLSClientConfig
func (r RawClient) create() Client {
	client := &http.Client{
		Timeout: r.Timeout.Duration,
		Transport: &http.Transport{
			Proxy:           getProxyFunc(r.ProxyURL),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: r.SkipVerify},
		}}

	if r.NotFollowRedirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error { return errors.New("redirect") }
	}
	return client
}

func (r RawClient) createWithProxyAuth() Client {
	auth := base64.StdEncoding.EncodeToString([]byte(r.ProxyUsername + ":" + r.ProxyPassword))
	client := r.create()

	return clientFunc(
		func(req *http.Request) (*http.Response, error) {
			req.Header.Set("Proxy-Authorization", "Basic "+auth)
			return client.Do(req)
		})
}

func getProxyFunc(u string) func(r *http.Request) (*url.URL, error) {
	if u == "" {
		return http.ProxyFromEnvironment
	}
	proxyURL, err := url.Parse(u)
	if err != nil {
		return func(r *http.Request) (*url.URL, error) { return nil, fmt.Errorf("invalid proxy: %s", err) }
	}
	return func(r *http.Request) (*url.URL, error) { return proxyURL, nil }
}

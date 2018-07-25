package web

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/l2isbad/go.d.plugin/internal/pkg/utils"
)

type Client struct {
	FollowRedirect bool           `yaml:"follow_redirects"`
	Timeout        utils.Duration `yaml:"timeout"`
	ProxyURL       string         `yaml:"proxy_url"`
	TLSVerify      bool           `yaml:"tls_verify"`
}

func (c Client) CreateHttpClient() *http.Client {
	client := &http.Client{
		Timeout: c.Timeout.Duration,
		Transport: &http.Transport{
			Proxy:           getProxyFunc(c.ProxyURL),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !c.TLSVerify},
		}}

	if !c.FollowRedirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error { return errors.New("redirect") }
	}
	return client
}

func getProxyFunc(u string) func(r *http.Request) (*url.URL, error) {
	if u == "" {
		return http.ProxyFromEnvironment
	}
	proxyURL, err := url.Parse(u)
	if err != nil || proxyURL.Scheme != "http" && proxyURL.Scheme != "https" {
		return func(r *http.Request) (*url.URL, error) { return nil, fmt.Errorf("invalid proxy: %s", err) }
	}
	return func(r *http.Request) (*url.URL, error) { return proxyURL, nil }
}

package url_helper

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/l2isbad/go.d.plugin/shared"
)

type Client struct {
	FollowRedirect bool            `toml:"follow_redirects"`
	Timeout        shared.Duration `toml:"timeout"`
	ProxyUrl       string          `toml:"proxy_url"`
	TLSVerify      bool            `toml:"tls_verify"`
}

func CreateHttpClient(c *Client) *http.Client {
	client := &http.Client{
		Timeout: c.Timeout.Duration,
		Transport: &http.Transport{
			Proxy:           getProxyFunc(c.ProxyUrl),
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
	proxyUrl, err := url.Parse(u)
	if err != nil || proxyUrl.Scheme != "http" && proxyUrl.Scheme != "https" {
		return func(r *http.Request) (*url.URL, error) { return nil, fmt.Errorf("invalid proxy: %s", err) }
	}
	return func(r *http.Request) (*url.URL, error) { return proxyUrl, nil }
}

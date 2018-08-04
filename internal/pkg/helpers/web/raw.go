package web

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/l2isbad/go.d.plugin/internal/pkg/utils"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// RawWeb RawWeb
type RawWeb struct {
	RawRequest `yaml:",inline"`
	RawClient  `yaml:",inline"`
}

// RawRequest is a structure that contains the fields that are needed to create *http.Request
type RawRequest struct {
	URL    string `yaml:"url" validate:"required,url"`
	Body   string `yaml:"body"`
	Method string `yaml:"method" validate:"isdefault|oneof=GET POST"`
}

// RawClient is a structure that contains the fields that are needed to create Client
type RawClient struct {
	Header        map[string]string `yaml:"headers"`
	Username      string            `yaml:"username"`
	Password      string            `yaml:"password"`
	ProxyUsername string            `yaml:"proxy_username"`
	ProxyPassword string            `yaml:"proxy_password"`

	FollowRedirect bool           `yaml:"follow_redirects"`
	Timeout        utils.Duration `yaml:"timeout"`
	ProxyURL       string         `yaml:"proxy_url"`
	TLSVerify      bool           `yaml:"tls_verify"`
}

// CreateHTTPRequest returns new *http.Requests and error if any encountered
func (r RawRequest) CreateHTTPRequest() (*http.Request, error) {
	var body io.Reader
	if r.Body != "" {
		body = strings.NewReader(r.Body)
	}
	req, err := http.NewRequest(r.Method, r.URL, body)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// CreateHTTPClient returns new Client
func (r RawClient) CreateHTTPClient() Client {
	var ds []decorator

	if r.Username != "" && r.Password != "" {
		ds = append(ds, authorization(r.Username, r.Password))
	}

	if r.ProxyUsername != "" && r.ProxyPassword != "" {
		ds = append(ds, proxyAuthorization(r.ProxyUsername, r.ProxyPassword))
	}

	if len(r.Header) != 0 {
		ds = append(ds, header(r.Header))
	}

	return decorate(r.create(), ds...)
}

func (r RawClient) create() *http.Client {
	client := &http.Client{
		Timeout: r.Timeout.Duration,
		Transport: &http.Transport{
			Proxy:           getProxyFunc(r.ProxyURL),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !r.TLSVerify},
		}}

	if !r.FollowRedirect {
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

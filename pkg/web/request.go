package web

import (
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Request is a struct that contains the fields that are needed to newHTTPClient *http.Request.
type Request struct {
	URL           *url.URL          `yaml:"-"`
	UserURL       string            `yaml:"url"`
	Body          string            `yaml:"body"`
	Method        string            `yaml:"method"`
	Headers       map[string]string `yaml:"headers"`
	Username      string            `yaml:"username"`
	Password      string            `yaml:"password"`
	ProxyUsername string            `yaml:"proxy_username"`
	ProxyPassword string            `yaml:"proxy_password"`
}

// Copy makes full copy of Request.
func (r Request) Copy() Request {
	h := make(map[string]string)
	for k, v := range r.Headers {
		h[k] = v
	}

	var u *url.URL
	if r.URL != nil {
		u, _ = url.Parse(r.URL.String())
	}

	r.URL = u
	r.Headers = h

	return r
}

// ParseUserURL parses UserURL into *url.URL and sets URL.
func (r *Request) ParseUserURL() (err error) {
	r.URL, err = url.Parse(r.UserURL)
	return err
}

// NewHTTPRequest creates a new *http.Requests based on Request fields
// and returns *http.Requests and error if any encountered.
func NewHTTPRequest(req Request) (*http.Request, error) {
	var body io.Reader
	if req.Body != "" {
		body = strings.NewReader(req.Body)
	}

	var u = req.UserURL
	if req.URL != nil {
		u = req.URL.String()
	}

	httpReq, err := http.NewRequest(req.Method, u, body)
	if err != nil {
		return nil, err
	}

	if req.Username != "" || req.Password != "" {
		httpReq.SetBasicAuth(req.Username, req.Password)
	}

	if req.ProxyUsername != "" && req.ProxyPassword != "" {
		basicAuth := base64.StdEncoding.EncodeToString([]byte(req.ProxyUsername + ":" + req.ProxyPassword))
		httpReq.Header.Set("Proxy-Authorization", "Basic "+basicAuth)
	}

	for k, v := range req.Headers {
		if k == "host" || k == "Host" {
			httpReq.Host = v
			continue
		}
		httpReq.Header.Set(k, v)
	}

	return httpReq, nil
}

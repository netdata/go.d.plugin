package web

import (
	"encoding/base64"
	"io"
	"net/http"
	"strings"
)

// Request is a struct that contains the fields that are needed to newHTTPClient *http.Request.
type Request struct {
	URL           string            `yaml:"url"`
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
	if r.Headers != nil {
		headers := make(map[string]string, len(r.Headers))
		for k, v := range r.Headers {
			headers[k] = v
		}
		r.Headers = headers
	}
	return r
}

// NewHTTPRequest creates a new *http.Requests based on Request fields
// and returns *http.Requests and error if any encountered.
func NewHTTPRequest(req Request) (*http.Request, error) {
	var body io.Reader
	if req.Body != "" {
		body = strings.NewReader(req.Body)
	}

	httpReq, err := http.NewRequest(req.Method, req.URL, body)
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
		switch k {
		case "host", "Host":
			httpReq.Host = v
		default:
			httpReq.Header.Set(k, v)
		}
	}
	return httpReq, nil
}

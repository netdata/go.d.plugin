package web

import (
	"io"
	"net/http"
	"strings"
)

// RawRequest is a struct that contains the fields that are needed to create *http.Request.
type RawRequest struct {
	URL      string            `yaml:"url" validate:"required,url"`
	Body     string            `yaml:"body"`
	Method   string            `yaml:"method" validate:"isdefault|oneof=GET POST HEAD PUT BATCH"`
	Headers  map[string]string `yaml:"headers"`
	Username string            `yaml:"username"`
	Password string            `yaml:"password"`
}

// CreateHTTPRequest creates a new *http.Requests based RawRequest fields
// and returns *http.Requests and error if any encountered.
func (r RawRequest) CreateHTTPRequest() (*http.Request, error) {
	var body io.Reader
	if r.Body != "" {
		body = strings.NewReader(r.Body)
	}
	req, err := http.NewRequest(r.Method, r.URL, body)
	if err != nil {
		return nil, err
	}
	if r.Username != "" && r.Password != "" {
		req.SetBasicAuth(r.Username, r.Password)
	}
	if len(r.Headers) != 0 {
		for k, v := range r.Headers {
			req.Header.Set(k, v)
		}
	}

	return req, nil
}

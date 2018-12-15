package web

import (
	"io"
	"net/http"
	"strings"
)

// Request is a struct that contains the fields that are needed to newHTTPClient *http.Request.
type Request struct {
	URI      string            `yaml:"-"`
	URL      string            `yaml:"url" validate:"required,url"`
	Body     string            `yaml:"body"`
	Method   string            `yaml:"method" validate:"isdefault|oneof=GET POST HEAD PUT BATCH"`
	Headers  map[string]string `yaml:"headers"`
	Username string            `yaml:"username"`
	Password string            `yaml:"password"`
}

// NewHTTPRequest creates a new *http.Requests based Request fields
// and returns *http.Requests and error if any encountered.
func NewHTTPRequest(req Request) (*http.Request, error) {
	var body io.Reader

	if req.Body != "" {
		body = strings.NewReader(req.Body)
	}

	httpReq, err := http.NewRequest(req.Method, joinURL(req.URL, req.URI), body)
	if err != nil {
		return nil, err
	}

	if req.Username != "" && req.Password != "" {
		httpReq.SetBasicAuth(req.Username, req.Password)
	}

	for k, v := range req.Headers {
		if k == "host" {
			httpReq.Host = v
			continue
		}
		httpReq.Header.Set(k, v)
	}

	return httpReq, nil
}

func joinURL(url, uri string) string {
	if uri == "" || url == "" {
		return url
	}

	if strings.HasSuffix(url, "/") {
		url = url[0 : len(url)-1]
	}

	if strings.HasPrefix(uri, "/") {
		uri = uri[1:]
	}

	return url + "/" + uri
}

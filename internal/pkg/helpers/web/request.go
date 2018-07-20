package web

import (
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	Url           string            `yaml:"url" validate:"required"`
	Body          string            `yaml:"body"`
	Header        map[string]string `yaml:"headers"`
	Method        string            `yaml:"method" validate:"isdefault|oneof=GET POST"`
	Username      string            `yaml:"username"`
	Password      string            `yaml:"password"`
	ProxyUsername string            `yaml:"proxy_username"`
	ProxyPassword string            `yaml:"proxy_password"`
}

func CreateRequest(r *Request) (*http.Request, error) {
	// URL is the only thing needed to create the Request
	if len(r.Url) == 0 {
		return nil, errors.New("empty URL")
	}

	u, err := url.Parse(r.Url)
	if err != nil {
		return nil, err
	}
	// TODO better URL validation needed ?
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, errors.New("unsupported scheme")
	}

	var body io.Reader
	if r.Body != "" {
		body = strings.NewReader(r.Body)
	}

	req, err := http.NewRequest(r.Method, r.Url, body)
	if err != nil {
		return nil, err
	}
	if r.Username != "" && r.Password != "" {
		req.SetBasicAuth(r.Username, r.Password)
	}

	if r.ProxyUsername != "" && r.ProxyPassword != "" {
		setProxyBasicAuth(req, r.ProxyUsername, r.ProxyPassword)
	}

	for k, v := range r.Header {
		req.Header.Add(k, v)
		if k == "Host" || k == "host" {
			req.Host = v
		}
	}
	return req, nil
}

func setProxyBasicAuth(r *http.Request, u, p string) {
	r.Header.Set("Proxy-Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(u+":"+p)))
}

package web

import (
	"encoding/base64"
	"net/http"
)

// Client Client
type Client interface {
	Do(r *http.Request) (*http.Response, error)
}

type clientFunc func(r *http.Request) (*http.Response, error)

func (f clientFunc) Do(r *http.Request) (*http.Response, error) {
	return f(r)
}

type decorator func(Client) Client

func authorization(user, password string) decorator {
	return func(c Client) Client {
		f := func(r *http.Request) (*http.Response, error) {
			r.SetBasicAuth(user, password)
			return c.Do(r)
		}
		return clientFunc(f)
	}
}

func proxyAuthorization(user, password string) decorator {
	return func(c Client) Client {
		f := func(r *http.Request) (*http.Response, error) {
			r.Header.Set("Proxy-Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(user+":"+password)))
			return c.Do(r)
		}
		return clientFunc(f)
	}
}

func header(h map[string]string) decorator {
	return func(c Client) Client {
		f := func(r *http.Request) (*http.Response, error) {
			for k, v := range h {
				r.Header.Set(k, v)
			}
			return c.Do(r)
		}
		return clientFunc(f)
	}
}

func decorate(c Client, ds ...decorator) Client {
	decorated := c
	for _, d := range ds {
		decorated = d(decorated)
	}
	return decorated
}

package web

import (
	"encoding/base64"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	username      = "user"
	password      = "password"
	headerKey     = "X-Api-Key"
	headerValue   = "secret"
	proxyUsername = "proxyUser"
	proxyPassword = "proxyPassword"
)

func TestRawRequest_CreateRequest(t *testing.T) {
	req, err := NewHTTPRequest(Request{
		Username: username,
		Password: password,
		Headers: map[string]string{
			headerKey: headerValue,
		},
		ProxyUsername: proxyUsername,
		ProxyPassword: proxyPassword,
	})

	assert.Nil(t, err)
	assert.IsType(t, (*http.Request)(nil), req)

	user, pass, ok := req.BasicAuth()
	assert.True(t, ok)
	assert.True(t, user == username && pass == password)

	user, pass, ok = parseBasicAuth(req.Header.Get("Proxy-Authorization"))
	assert.True(t, ok)
	assert.True(t, user == proxyUsername && pass == proxyPassword)

	assert.True(t, req.Header.Get(headerKey) == headerValue)
}

func parseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	// Case insensitive prefix match. See Issue 22736.
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}

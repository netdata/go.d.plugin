package web

import (
	"encoding/base64"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testURI           = "testURI"
	testURL           = "testURL"
	testMethod        = "testMethod"
	testUsername      = "user"
	testPassword      = "testPassword"
	testHeaderKey     = "X-Api-Key"
	testHeaderValue   = "secret"
	testProxyUsername = "proxyUser"
	testProxyPassword = "testProxyPassword"
)

func TestRequest_Copy(t *testing.T) {
	var r Request
	r.URL = testURI
	r.URL = testURL
	r.Method = testMethod
	r.Headers = map[string]string{
		testHeaderKey: testHeaderValue,
	}
	r.Username = testUsername
	r.Password = testPassword
	r.ProxyUsername = testProxyUsername
	r.ProxyPassword = testProxyPassword

	rr := r.Copy()
	assert.Equal(t, r, *rr)
	rr.Headers[""] = ""
	assert.NotEqual(t, r, *rr)
}

func TestNewHTTPRequest(t *testing.T) {
	req, err := NewHTTPRequest(Request{
		Username: testUsername,
		Password: testPassword,
		Headers: map[string]string{
			testHeaderKey: testHeaderValue,
		},
		ProxyUsername: testProxyUsername,
		ProxyPassword: testProxyPassword,
	})

	assert.Nil(t, err)
	assert.IsType(t, (*http.Request)(nil), req)

	user, pass, ok := req.BasicAuth()
	assert.True(t, ok)
	assert.Equal(t, testUsername, user)
	assert.Equal(t, testPassword, pass)

	user, pass, ok = parseBasicAuth(req.Header.Get("Proxy-Authorization"))
	assert.True(t, ok)

	assert.Equal(t, testProxyUsername, user)
	assert.Equal(t, testProxyPassword, pass)
	assert.Equal(t, testHeaderValue, req.Header.Get(testHeaderKey))
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

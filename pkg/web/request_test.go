package web

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testPath          = "/api/v1/info"
	testURL           = "http://127.0.0.1"
	testMethod        = "POST"
	testUsername      = "user"
	testPassword      = "password"
	testHeaderKey     = "X-Api-Key"
	testHeaderValue   = "secret"
	testProxyUsername = "proxy_user"
	testProxyPassword = "proxy_password"
)

func TestRequest_Copy(t *testing.T) {
	var r1 Request
	r1.URL = &url.URL{Path: testPath}
	r1.UserURL = testURL
	r1.Method = testMethod
	r1.Headers = map[string]string{
		testHeaderKey: testHeaderValue,
	}
	r1.Username = testUsername
	r1.Password = testPassword
	r1.ProxyUsername = testProxyUsername
	r1.ProxyPassword = testProxyPassword

	r2 := r1.Copy()
	r3 := r1.Copy()
	assert.Equal(t, r1, r2)
	r2.Headers[""] = ""
	assert.NotEqual(t, r1, r2)
	r3.URL.Path = ""
	assert.NotEqual(t, r1, r3)
}

func TestNewHTTPRequest(t *testing.T) {
	r := Request{
		UserURL:  testURL,
		Method:   testMethod,
		Username: testUsername,
		Password: testPassword,
		Headers: map[string]string{
			testHeaderKey: testHeaderValue,
		},
		ProxyUsername: testProxyUsername,
		ProxyPassword: testProxyPassword,
	}

	assert.NoError(t, r.ParseUserURL())
	assert.NotNil(t, r.URL)
	r.URL.Path = testPath

	req, err := NewHTTPRequest(r)
	assert.NoError(t, err)
	assert.IsType(t, (*http.Request)(nil), req)

	assert.Equal(t, testMethod, req.Method)

	user, pass, ok := req.BasicAuth()
	assert.True(t, ok)
	assert.Equal(t, testUsername, user)
	assert.Equal(t, testPassword, pass)

	user, pass, ok = parseBasicAuth(req.Header.Get("Proxy-Authorization"))
	assert.True(t, ok)
	assert.Equal(t, testProxyUsername, user)
	assert.Equal(t, testProxyPassword, pass)
	assert.Equal(t, testHeaderValue, req.Header.Get(testHeaderKey))

	assert.Equal(t, testURL+testPath, req.URL.String())
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

package web

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	username      = "user"
	password      = "password"
	proxyUsername = "proxyUser"
	proxyPassword = "proxyPassword"
	headerKey     = "X-Api-Key"
	headerValue   = "secret"
)

func TestRawRequest_CreateRequest(t *testing.T) {
	rawRequest := RawRequest{}
	r, err := rawRequest.CreateHTTPRequest()
	assert.Nil(t, err)
	assert.IsType(t, (*http.Request)(nil), r)
}

func TestRawClient_CreateClient(t *testing.T) {
	rawClient := RawClient{
		Header: map[string]string{
			headerKey: headerValue,
		},
		Username:      username,
		Password:      password,
		ProxyUsername: proxyUsername,
		ProxyPassword: proxyPassword,
	}

	client := rawClient.CreateHTTPClient()

	req, _ := http.NewRequest("GET", "", nil)

	_, _ = client.Do(req)

	user, pass, ok := req.BasicAuth()

	assert.True(t, ok)
	assert.Equal(t, username, user)
	assert.Equal(t, password, pass)
	assert.Equal(t, req.Header.Get(headerKey), headerValue)
	assert.NotEmpty(t, req.Header.Get("Proxy-Authorization"))
}

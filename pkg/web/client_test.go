package web

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	proxyUsername = "proxyUser"
	proxyPassword = "proxyPassword"
)

func TestNewHTTPClient(t *testing.T) {
	req, _ := http.NewRequest("GET", "", nil)

	// W/o Proxy Authorization
	var client Client

	httpClient := NewHTTPClient(client)

	assert.Implements(t, (*HTTPClient)(nil), httpClient)

	_, _ = httpClient.Do(req)

	assert.Empty(t, req.Header.Get("Proxy-Authorization"))

	// W/ Proxy Authorization

	client.ProxyUsername = proxyUsername
	client.ProxyPassword = proxyPassword

	httpClient = NewHTTPClient(client)

	assert.Implements(t, (*HTTPClient)(nil), httpClient)

	_, _ = httpClient.Do(req)

	assert.NotEmpty(t, req.Header.Get("Proxy-Authorization"))

}

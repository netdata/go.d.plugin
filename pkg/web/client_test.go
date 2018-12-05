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

func TestRawClient_CreateHTTPClient(t *testing.T) {
	req, _ := http.NewRequest("GET", "", nil)
	// Without Proxy Authorization
	rawClient := new(RawClient)

	client := rawClient.CreateHTTPClient()
	assert.Implements(t, (*Client)(nil), client)

	_, _ = client.Do(req)
	assert.Empty(t, req.Header.Get("Proxy-Authorization"))

	rawClient.ProxyUsername = proxyUsername
	rawClient.ProxyPassword = proxyPassword

	// With Proxy Authorization
	client = rawClient.CreateHTTPClient()
	assert.Implements(t, (*Client)(nil), client)

	_, _ = client.Do(req)
	assert.NotEmpty(t, req.Header.Get("Proxy-Authorization"))
}

package web

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//var (
//	proxyUsername = "proxyUser"
//	proxyPassword = "proxyPassword"
//)

func TestNewHTTPClient(t *testing.T) {
	client := NewHTTPClient(Client{
		Timeout:           Duration{Duration: time.Second * 5},
		NotFollowRedirect: true,
		TLSSkipVerify:     true,
		ProxyURL:          "http://127.0.0.1:3128",
	})

	assert.IsType(t, (*http.Client)(nil), client)
	assert.Equal(t, time.Second*5, client.Timeout)
	assert.NotNil(t, client.CheckRedirect)
}

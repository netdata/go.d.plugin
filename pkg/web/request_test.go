package web

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	username    = "user"
	password    = "password"
	headerKey   = "X-Api-Key"
	headerValue = "secret"
)

func TestRawRequest_CreateRequest(t *testing.T) {
	req := Request{
		Username: username,
		Password: password,
		Headers: map[string]string{
			headerKey: headerValue,
		},
	}

	httpReq, err := NewHTTPRequest(req)

	assert.Nil(t, err)
	assert.IsType(t, (*http.Request)(nil), httpReq)

	user, pass, ok := httpReq.BasicAuth()

	assert.True(t, ok)
	assert.True(t, user == username && pass == password)
	assert.True(t, httpReq.Header.Get(headerKey) == headerValue)
}

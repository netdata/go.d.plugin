package httpcheck

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.IsType(t, (*HTTPCheck)(nil), New())
}

func TestHTTPCheck_Init(t *testing.T) {
	hc := New()

	assert.True(t, hc.Init())

	hc.ResponseMatch = "(?:qwe))"

	assert.False(t, hc.Init())
}

func TestHTTPCheck_Check(t *testing.T) {
	assert.True(t, New().Check())
}

func TestHTTPCheck_GetCharts(t *testing.T) {
	assert.NotNil(t, New().GetCharts())
}

func TestHTTPCheck_ResponseSuccess(t *testing.T) {
	hc := New()
	msg := "hello"

	var resp http.Response

	resp.Body = nopCloser{bytes.NewBufferString(msg)}
	resp.StatusCode = http.StatusOK

	hc.processOKResponse(&resp)

	assert.Equal(
		t,
		data{
			Success:        1,
			Failed:         0,
			Timeout:        0,
			BadContent:     0,
			BadStatus:      0,
			ResponseTime:   0,
			ResponseLength: len(msg),
		},
		hc.data,
	)
}

func TestHTTPCheck_ResponseSuccessBadContent(t *testing.T) {
	hc := New()
	hc.ResponseMatch = "no match"
	require.True(t, hc.Init())

	msg := "hello"

	var resp http.Response

	resp.Body = nopCloser{bytes.NewBufferString(msg)}
	resp.StatusCode = http.StatusOK

	hc.processOKResponse(&resp)

	assert.Equal(
		t,
		data{
			Success:        1,
			Failed:         0,
			Timeout:        0,
			BadContent:     1,
			BadStatus:      0,
			ResponseTime:   0,
			ResponseLength: len(msg),
		},
		hc.data,
	)
}

func TestHTTPCheck_ResponseTimeout(t *testing.T) {
	hc := New()

	var err net.Error = timeoutError{}

	hc.processErrResponse(err)

	assert.Equal(
		t,
		data{
			Success:        0,
			Failed:         0,
			Timeout:        1,
			BadContent:     0,
			BadStatus:      0,
			ResponseTime:   0,
			ResponseLength: 0,
		},
		hc.data,
	)
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

type timeoutError struct {
}

func (r timeoutError) Timeout() bool {
	return true
}

func (r timeoutError) Error() string {
	return ""
}

func (r timeoutError) Temporary() bool {
	return true
}

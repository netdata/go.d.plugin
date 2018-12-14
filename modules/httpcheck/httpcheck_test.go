package httpcheck

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	assert.IsType(t, (*HTTPCheck)(nil), New())
}

func TestHTTPCheck_Init(t *testing.T) {
	mod := New()

	assert.True(t, mod.Init())

	mod.ResponseMatch = "(?:qwe))"

	assert.False(t, mod.Init())
}

func TestHTTPCheck_Check(t *testing.T) {
	assert.True(t, New().Check())
}

func TestHTTPCheck_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestHTTPCheck_Cleanup(t *testing.T) {
	New().Cleanup()
}

func TestHTTPCheck_Collect(t *testing.T) {
	mod := New()

	srv := httptest.NewServer(myHandler{})
	defer srv.Close()

	mod.URL = srv.URL
	mod.Init()
	assert.NotNil(t, mod.Collect())
}

func TestHTTPCheck_ResponseSuccess(t *testing.T) {
	mod := New()
	msg := "hello"

	var resp http.Response

	resp.Body = nopCloser{bytes.NewBufferString(msg)}
	resp.StatusCode = http.StatusOK

	mod.processOKResponse(&resp)

	assert.Equal(
		t,
		metrics{
			Success:        1,
			Failed:         0,
			Timeout:        0,
			BadContent:     0,
			BadStatus:      0,
			ResponseTime:   0,
			ResponseLength: len(msg),
		},
		mod.metrics,
	)
}

func TestHTTPCheck_ResponseSuccessBadContent(t *testing.T) {
	mod := New()
	mod.ResponseMatch = "no match"
	require.True(t, mod.Init())

	msg := "hello"

	var resp http.Response

	resp.Body = nopCloser{bytes.NewBufferString(msg)}
	resp.StatusCode = http.StatusOK

	mod.processOKResponse(&resp)

	assert.Equal(
		t,
		metrics{
			Success:        1,
			Failed:         0,
			Timeout:        0,
			BadContent:     1,
			BadStatus:      0,
			ResponseTime:   0,
			ResponseLength: len(msg),
		},
		mod.metrics,
	)
}

func TestHTTPCheck_ResponseTimeout(t *testing.T) {
	mod := New()

	var err net.Error = timeoutError{}

	mod.processErrResponse(err)

	assert.Equal(
		t,
		metrics{
			Success:        0,
			Failed:         0,
			Timeout:        1,
			BadContent:     0,
			BadStatus:      0,
			ResponseTime:   0,
			ResponseLength: 0,
		},
		mod.metrics,
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

type myHandler struct{}

func (myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

package httpcheck

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/netdata/go-orchestrator/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	assert.Implements(t, (*module.Module)(nil), New())
}

func TestHTTPCheck_Cleanup(t *testing.T) {
	New().Cleanup()
}

func TestHTTPCheck_Init(t *testing.T) {
	job := New()
	job.UserURL = "http://127.0.0.1:38001"

	// OK case
	assert.True(t, job.Init())
	assert.NotNil(t, job.client)

	// NG case
	job.ResponseMatch = "(?:qwe))"
	assert.False(t, job.Init())
}

func TestHTTPCheck_Check(t *testing.T) {
	job := New()
	job.UserURL = "http://127.0.0.1:38001"

	require.True(t, job.Init())
	assert.True(t, job.Check())
}

func TestHTTPCheck_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
	assert.NoError(t, module.CheckCharts(*New().Charts()...))
}

func TestHTTPCheck_Collect(t *testing.T) {
	job := New()

	ts := httptest.NewServer(myHandler{})
	defer ts.Close()

	job.UserURL = ts.URL
	require.True(t, job.Init())
	assert.NotNil(t, job.Collect())
}

func TestHTTPCheck_ResponseSuccess(t *testing.T) {
	job := New()
	msg := "hello"

	var resp http.Response

	resp.Body = nopCloser{bytes.NewBufferString(msg)}
	resp.StatusCode = http.StatusOK

	job.processOKResponse(&resp)

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
		job.metrics,
	)
}

func TestHTTPCheck_ResponseSuccessInvalidContent(t *testing.T) {
	job := New()
	job.ResponseMatch = "no match"
	job.UserURL = "http://127.0.0.1:38001"
	require.True(t, job.Init())

	msg := "hello"

	var resp http.Response

	resp.Body = nopCloser{bytes.NewBufferString(msg)}
	resp.StatusCode = http.StatusOK

	job.processOKResponse(&resp)

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
		job.metrics,
	)
}

func TestHTTPCheck_ResponseTimeout(t *testing.T) {
	job := New()

	var err net.Error = timeoutError{}

	job.processErrResponse(err)

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
		job.metrics,
	)
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

type timeoutError struct{}

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

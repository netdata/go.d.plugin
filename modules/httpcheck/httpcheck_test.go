package httpcheck

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"net/url"
	"testing"

	"github.com/netdata/go.d.plugin/pkg/stm"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	job := New()
	assert.Implements(t, (*module.Module)(nil), job)
	assert.Equal(t, defaultHTTPTimeout, job.Timeout.Duration)
	assert.Equal(t, defaultAcceptedStatuses, job.AcceptedStatuses)
}

func TestHTTPCheck_Cleanup(t *testing.T) { New().Cleanup() }

func TestHTTPCheck_Init(t *testing.T) {
	job := New()

	job.UserURL = "http://127.0.0.1:38001"
	assert.True(t, job.Init())
	assert.NotNil(t, job.client)
}

func TestHTTPCheck_InitNG(t *testing.T) {
	job := New()

	assert.False(t, job.Init())
	job.UserURL = "http://127.0.0.1:38001"
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
	assert.False(t, New().Charts().Has(bodyLengthChart.ID))

	job := New()
	job.UserURL = "http://127.0.0.1"
	job.ResponseMatch = "1"
	require.True(t, job.Init())
	assert.True(t, job.Charts().Has(bodyLengthChart.ID))
}

func TestHTTPCheck_Collect(t *testing.T) {
	job := New()

	job.UserURL = "http://127.0.0.1"
	job.ResponseMatch = "hello"
	require.True(t, job.Init())

	job.client = clientFunc(func(r *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       nopCloser{bytes.NewBufferString("hello")},
		}
		return resp, nil
	})
	assert.Equal(
		t,
		stm.ToMap(metrics{Status: status{Success: true}}),
		job.Collect(),
	)
}

func TestHTTPCheck_Collect_TimeoutError(t *testing.T) {
	job := New()

	job.UserURL = "http://127.0.0.1"
	job.ResponseMatch = "hello"
	require.True(t, job.Init())

	job.client = clientFunc(func(r *http.Request) (*http.Response, error) { return nil, timeoutError{} })
	assert.Equal(
		t,
		stm.ToMap(metrics{Status: status{Timeout: true}}),
		job.Collect(),
	)

}

func TestHTTPCheck_Collect_DNSLookupError(t *testing.T) {
	job := New()

	job.UserURL = "http://127.0.0.1"
	job.ResponseMatch = "hello"
	require.True(t, job.Init())

	job.client = clientFunc(func(r *http.Request) (*http.Response, error) {
		return nil, net.Error(&url.Error{Err: &net.OpError{Err: &net.DNSError{}}})
	})
	assert.Equal(
		t,
		stm.ToMap(metrics{Status: status{DNSLookupError: true}}),
		job.Collect(),
	)
}

func TestHTTPCheck_Collect_AddressParseError(t *testing.T) {
	job := New()

	job.UserURL = "http://127.0.0.1"
	job.ResponseMatch = "hello"
	require.True(t, job.Init())

	job.client = clientFunc(func(r *http.Request) (*http.Response, error) {
		return nil, net.Error(&url.Error{Err: &net.OpError{Err: &net.ParseError{}}})
	})
	assert.Equal(
		t,
		stm.ToMap(metrics{Status: status{ParseAddressError: true}}),
		job.Collect(),
	)

}

func TestHTTPCheck_Collect_RedirectError(t *testing.T) {
	job := New()

	job.UserURL = "http://127.0.0.1"
	job.ResponseMatch = "hello"
	require.True(t, job.Init())

	job.client = clientFunc(func(r *http.Request) (*http.Response, error) {
		return nil, net.Error(&url.Error{Err: web.ErrRedirectAttempted})
	})
	assert.Equal(
		t,
		stm.ToMap(metrics{Status: status{RedirectError: true}}),
		job.Collect(),
	)
}

func TestHTTPCheck_Collect_BadContentError(t *testing.T) {
	job := New()

	job.UserURL = "http://127.0.0.1"
	job.ResponseMatch = "hello"
	require.True(t, job.Init())

	job.client = clientFunc(func(r *http.Request) (*http.Response, error) {
		resp := &http.Response{StatusCode: http.StatusOK, Body: nopCloser{bytes.NewBufferString("good bye")}}
		return resp, nil
	})
	assert.Equal(
		t,
		stm.ToMap(metrics{Status: status{BadContent: true}}),
		job.Collect(),
	)
}

func TestHTTPCheck_Collect_BadStatusError(t *testing.T) {
	job := New()

	job.UserURL = "http://127.0.0.1"
	job.ResponseMatch = "hello"
	require.True(t, job.Init())

	job.client = clientFunc(func(r *http.Request) (*http.Response, error) {
		resp := &http.Response{StatusCode: http.StatusBadGateway}
		return resp, nil
	})
	assert.Equal(
		t,
		stm.ToMap(metrics{Status: status{BadStatusCode: true}}),
		job.Collect(),
	)
}

type clientFunc func(r *http.Request) (*http.Response, error)

func (f clientFunc) Do(r *http.Request) (*http.Response, error) { return f(r) }

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

type timeoutError struct{}

func (r timeoutError) Timeout() bool { return true }

func (r timeoutError) Error() string { return "" }

func (r timeoutError) Temporary() bool { return true }

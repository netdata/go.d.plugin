// SPDX-License-Identifier: GPL-3.0-or-later

package httpcheck

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/netdata/go.d.plugin/pkg/stm"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testURL = "http://127.0.0.1:38001"
)

func TestNew(t *testing.T) {
	job := New()
	assert.Implements(t, (*module.Module)(nil), job)
}

func TestHTTPCheck_Cleanup(t *testing.T) { New().Cleanup() }

func TestHTTPCheck_Init(t *testing.T) {
	job := New()

	job.URL = testURL
	assert.True(t, job.Init())
	assert.NotNil(t, job.client)
}

func TestHTTPCheck_InitNG(t *testing.T) {
	job := New()

	assert.False(t, job.Init())
	job.URL = testURL
	job.ResponseMatch = "(?:qwe))"
	assert.False(t, job.Init())
}

func TestHTTPCheck_Check(t *testing.T) {
	job := New()
	job.URL = testURL

	require.True(t, job.Init())
	assert.True(t, job.Check())
}

func TestHTTPCheck_Charts(t *testing.T) {
	job := New()
	job.URL = testURL
	require.True(t, job.Init())
	assert.NotNil(t, job.Charts())
}

func TestHTTPCheck_Collect(t *testing.T) {
	job := New()
	body := "hello"

	job.URL = testURL
	job.ResponseMatch = body
	job.UpdateEvery = 15
	require.True(t, job.Init())

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       nopCloser{bytes.NewBufferString(body)},
	}
	job.client = newClientFunc(resp, nil)

	assert.Equal(
		t,
		stm.ToMap(metrics{
			Status:         status{Success: true},
			InState:        job.UpdateEvery,
			ResponseLength: len(body),
			ResponseTime:   0,
		}),
		job.Collect(),
	)
}

func TestHTTPCheck_Collect_TimeoutError(t *testing.T) {
	job := New()

	job.URL = testURL
	require.True(t, job.Init())

	job.client = newClientFunc(nil, timeoutError{})
	assert.Equal(
		t,
		stm.ToMap(metrics{Status: status{Timeout: true}}),
		job.Collect(),
	)

}

//func TestHTTPCheck_Collect_DNSLookupError(t *testing.T) {
//	job := New()
//
//	job.UserURL = testURL
//	require.True(t, job.Init())
//
//	err := net.Error(&url.Error{Err: &net.OpError{Err: &net.DNSError{}}})
//	job.client = newClientFunc(nil, err)
//	assert.Equal(
//		t,
//		stm.ToMap(metrics{Status: status{DNSLookupError: true}}),
//		job.Collect(),
//	)
//}

//func TestHTTPCheck_Collect_AddressParseError(t *testing.T) {
//	job := New()
//
//	job.UserURL = testURL
//	require.True(t, job.Init())
//
//	err := net.Error(&url.Error{Err: &net.OpError{Err: &net.ParseError{}}})
//	job.client = newClientFunc(nil, err)
//	assert.Equal(
//		t,
//		stm.ToMap(metrics{Status: status{ParseAddressError: true}}),
//		job.Collect(),
//	)
//
//}

//func TestHTTPCheck_Collect_RedirectError(t *testing.T) {
//	job := New()
//
//	job.UserURL = testURL
//	require.True(t, job.Init())
//
//	err := net.Error(&url.Error{Err: web.ErrRedirectAttempted})
//	job.client = newClientFunc(nil, err)
//	assert.Equal(
//		t,
//		stm.ToMap(metrics{Status: status{RedirectError: true}}),
//		job.Collect(),
//	)
//}

func TestHTTPCheck_Collect_BadContentError(t *testing.T) {
	job := New()
	body := "hello"

	job.URL = testURL
	job.ResponseMatch = "not match"
	require.True(t, job.Init())

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       nopCloser{bytes.NewBufferString(body)},
	}
	job.client = newClientFunc(resp, nil)
	assert.Equal(
		t,
		stm.ToMap(metrics{
			Status:         status{BadContent: true},
			ResponseLength: len(body),
		}),
		job.Collect(),
	)
}

func TestHTTPCheck_Collect_BadStatusError(t *testing.T) {
	job := New()

	job.URL = testURL
	require.True(t, job.Init())

	resp := &http.Response{StatusCode: http.StatusBadGateway}
	job.client = newClientFunc(resp, nil)
	assert.Equal(
		t,
		stm.ToMap(metrics{Status: status{BadStatusCode: true}}),
		job.Collect(),
	)
}

func TestHTTPCheck_Collect_InState(t *testing.T) {
	job := New()
	goodBody := "hello"
	badBody := "goodbye"

	job.URL = testURL
	job.ResponseMatch = goodBody
	job.UpdateEvery = 15
	require.True(t, job.Init())

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       nopCloser{bytes.NewBufferString(goodBody)},
	}
	job.client = newClientFunc(resp, nil)

	assert.Equal(
		t,
		stm.ToMap(metrics{
			Status:         status{Success: true},
			InState:        job.UpdateEvery,
			ResponseLength: len(goodBody),
			ResponseTime:   0,
		}),
		job.Collect(),
	)

	resp = &http.Response{
		StatusCode: http.StatusOK,
		Body:       nopCloser{bytes.NewBufferString(goodBody)},
	}
	job.client = newClientFunc(resp, nil)

	assert.Equal(
		t,
		stm.ToMap(metrics{
			Status:         status{Success: true},
			InState:        job.UpdateEvery * 2,
			ResponseLength: len(goodBody),
			ResponseTime:   0,
		}),
		job.Collect(),
	)

	resp = &http.Response{
		StatusCode: http.StatusOK,
		Body:       nopCloser{bytes.NewBufferString(badBody)},
	}
	job.client = newClientFunc(resp, nil)

	assert.Equal(
		t,
		stm.ToMap(metrics{
			Status:         status{BadContent: true},
			InState:        job.UpdateEvery,
			ResponseLength: len(badBody),
			ResponseTime:   0,
		}),
		job.Collect(),
	)
}

type clientFunc func(r *http.Request) (*http.Response, error)

func (f clientFunc) Do(r *http.Request) (*http.Response, error) { return f(r) }

func newClientFunc(resp *http.Response, err error) clientFunc {
	return func(r *http.Request) (*http.Response, error) { return resp, err }
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

type timeoutError struct{}

func (r timeoutError) Timeout() bool { return true }

func (r timeoutError) Error() string { return "" }

func (r timeoutError) Temporary() bool { return true }

package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testToken   = "token"
	testVersion = "2.5"
)

func TestNewClient(t *testing.T) {
	client := NewClient(&http.Client{}, web.Request{})
	assert.IsType(t, (*Client)(nil), client)
	assert.NotNil(t, client.httpClient)
	assert.NotNil(t, client.token)
}

func TestClient_IsLoggedIn(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	req := web.Request{UserURL: ts.URL}
	_ = req.ParseUserURL()

	client := NewClient(&http.Client{}, req)

	require.NoError(t, client.Login())
	assert.True(t, client.IsLoggedIn())
}

func TestClient_Login(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	req := web.Request{UserURL: ts.URL}
	_ = req.ParseUserURL()

	client := NewClient(&http.Client{}, req)

	require.NoError(t, client.Login())
	assert.Equal(t, testToken, client.token.get())
}

func TestClient_Logout(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	req := web.Request{UserURL: ts.URL}
	_ = req.ParseUserURL()

	client := NewClient(&http.Client{}, req)

	require.NoError(t, client.Login())
	require.True(t, client.IsLoggedIn())
	assert.NoError(t, client.Logout())
	assert.False(t, client.IsLoggedIn())
}

func TestClient_GetAPIVersion(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	req := web.Request{UserURL: ts.URL}
	_ = req.ParseUserURL()

	client := NewClient(&http.Client{}, req)

	ver, err := client.GetAPIVersion()
	require.NoError(t, err)
	assert.Equal(t, &Version{Major: 2, Minor: 5}, ver)
}

func TestClient_GetSelectedStatistics(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	req := web.Request{UserURL: ts.URL}
	_ = req.ParseUserURL()

	client := NewClient(&http.Client{}, req)
	require.NoError(t, client.Login())
	dst := &struct {
		A, B int
	}{}
	require.NoError(t, client.GetSelectedStatistics(dst, ""))
	assert.Equal(t, 1, dst.A)
	assert.Equal(t, 2, dst.B)
}

func newTestServer() *httptest.Server {
	handle := func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		default:
			w.WriteHeader(http.StatusBadRequest)
		case PATHLogin:
			_, _ = w.Write([]byte(testToken))
		case PATHLogout:
		case PATHVersion:
			_, _ = w.Write([]byte(testVersion))
		case PATHSelectedStatistics:
			_, _ = w.Write([]byte(`{"A": 1, "B": 2}`))
		}
	}

	return httptest.NewServer(http.HandlerFunc(handle))
}

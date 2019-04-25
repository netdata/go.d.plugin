package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	client := NewClient(&http.Client{}, web.Request{})
	assert.IsType(t, (*Client)(nil), client)
	assert.NotNil(t, client.httpClient)
	assert.NotNil(t, client.token)
}

func TestClient_IsLoggedIn(t *testing.T) {
	secret := "secret token"

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				default:
					w.WriteHeader(http.StatusBadRequest)
				case PATHLogin:
					_, _ = w.Write([]byte(secret))
				}
			}))
	defer ts.Close()

	req := web.Request{UserURL: ts.URL}
	_ = req.ParseUserURL()

	client := NewClient(&http.Client{}, req)

	require.NoError(t, client.Login())
	assert.True(t, client.IsLoggedIn())

}

func TestClient_Login(t *testing.T) {
	secret := "secret token"

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				default:
					w.WriteHeader(http.StatusBadRequest)
				case PATHLogin:
					_, _ = w.Write([]byte(secret))
				}
			}))
	defer ts.Close()

	req := web.Request{UserURL: ts.URL}
	_ = req.ParseUserURL()

	client := NewClient(&http.Client{}, req)

	require.NoError(t, client.Login())
	assert.Equal(t, secret, client.token.get())
}

func TestClient_Logout(t *testing.T) {
	secret := "secret token"

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				default:
					w.WriteHeader(http.StatusBadRequest)
				case PATHLogin:
					_, _ = w.Write([]byte(secret))
				case PATHLogout:
				}
			}))
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
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				default:
					w.WriteHeader(http.StatusBadRequest)
				case PATHVersion:
					_, _ = w.Write([]byte("2.5"))
				}
			}))
	defer ts.Close()

	req := web.Request{UserURL: ts.URL}
	_ = req.ParseUserURL()

	client := NewClient(&http.Client{}, req)

	ver, err := client.GetAPIVersion()
	require.NoError(t, err)
	assert.Equal(t, &Version{Major: 2, Minor: 5}, ver)
}

func TestClient_GetSelectedStatistics(t *testing.T) {
	secret := "secret token"
	query := "{query: query}"

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				default:
					w.WriteHeader(http.StatusBadRequest)
				case PATHLogin:
					_, _ = w.Write([]byte(secret))
				case PATHSelectedStatistics:
					_, pass, _ := r.BasicAuth()
					bs, _ := ioutil.ReadAll(r.Body)
					ct := r.Header.Get("Content-Type")

					switch {
					default:
						_, _ = w.Write([]byte(`{"A": 1, "B": 2}`))
					case string(bs) != query:
						err := apiError{
							Message: fmt.Sprintf("wrong query, expect %s, got %s", query, string(bs)),
						}
						b, _ := json.Marshal(err)
						w.WriteHeader(http.StatusBadRequest)
						_, _ = w.Write(b)
					case r.Method != http.MethodPost:
						err := apiError{
							Message: fmt.Sprintf("wrong req method, expect %s, got %s", http.MethodPost, r.Method),
						}
						b, _ := json.Marshal(err)
						w.WriteHeader(http.StatusBadRequest)
						_, _ = w.Write(b)
					case pass != secret:
						err := apiError{
							Message: fmt.Sprintf("wrong password, expect %s, got %s", secret, pass),
						}
						b, _ := json.Marshal(err)
						w.WriteHeader(http.StatusBadRequest)
						_, _ = w.Write(b)
					case ct != "application/json":
						err := apiError{
							Message: fmt.Sprintf("wrong content type, expect %s, got %s", "application/json", ct),
						}
						b, _ := json.Marshal(err)
						w.WriteHeader(http.StatusBadRequest)
						_, _ = w.Write(b)
					}
				}
			}))
	defer ts.Close()

	req := web.Request{UserURL: ts.URL}
	_ = req.ParseUserURL()

	client := NewClient(&http.Client{}, req)
	require.NoError(t, client.Login())
	dst := &struct {
		A, B int
	}{}
	require.NoError(t, client.GetSelectedStatistics(dst, query))
	assert.Equal(t, 1, dst.A)
	assert.Equal(t, 2, dst.B)
}

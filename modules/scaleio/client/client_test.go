package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	New(nil, web.Request{})
}

func TestClient_Login(t *testing.T) {
	srv := newTestScaleIOServer()
	defer srv.Close()

	client := New(nil, web.Request{
		UserURL:  srv.URL,
		Username: testUser,
		Password: testPassword,
	})

	assert.NoError(t, client.Login())
	assert.Equal(t, testToken, client.token.get())
}

func TestClient_Logout(t *testing.T) {
	srv := newTestScaleIOServer()
	defer srv.Close()

	client := New(nil, web.Request{
		UserURL:  srv.URL,
		Username: testUser,
		Password: testPassword,
	})

	require.NoError(t, client.Login())

	assert.NoError(t, client.Logout())
	assert.False(t, client.token.isSet())

}

func TestClient_LoggedIn(t *testing.T) {
	srv := newTestScaleIOServer()
	defer srv.Close()

	client := New(nil, web.Request{
		UserURL:  srv.URL,
		Username: testUser,
		Password: testPassword,
	})

	assert.False(t, client.LoggedIn())
	assert.NoError(t, client.Login())
	assert.True(t, client.LoggedIn())
}

func TestClient_APIVersion(t *testing.T) {
	srv := newTestScaleIOServer()
	defer srv.Close()

	client := New(nil, web.Request{
		UserURL:  srv.URL,
		Username: testUser,
		Password: testPassword,
	})

	err := client.Login()
	require.NoError(t, err)

	version, err := client.APIVersion()
	assert.NoError(t, err)
	assert.Equal(t, Version{Major: 2, Minor: 5}, version)
}

func TestClient_Instances(t *testing.T) {
	srv := newTestScaleIOServer()
	defer srv.Close()

	client := New(nil, web.Request{
		UserURL:  srv.URL,
		Username: testUser,
		Password: testPassword,
	})

	err := client.Login()
	require.NoError(t, err)

	instances, err := client.Instances()
	assert.NoError(t, err)
	assert.Equal(t, testInstances, instances)
}

func TestClient_Instances_RetryOnExpiredToken(t *testing.T) {
	srv := newTestScaleIOServer()
	defer srv.Close()

	client := New(nil, web.Request{
		UserURL:  srv.URL,
		Username: testUser,
		Password: testPassword,
	})

	instances, err := client.Instances()
	assert.NoError(t, err)
	assert.Equal(t, testInstances, instances)
}

func TestClient_SelectedStatistics(t *testing.T) {
	srv := newTestScaleIOServer()
	defer srv.Close()

	client := New(nil, web.Request{
		UserURL:  srv.URL,
		Username: testUser,
		Password: testPassword,
	})

	err := client.Login()
	require.NoError(t, err)

	stats, err := client.SelectedStatistics(SelectedStatisticsQuery{})
	assert.NoError(t, err)
	assert.Equal(t, testStatistics, stats)
}

func TestClient_SelectedStatistics_RetryOnExpiredToken(t *testing.T) {
	srv := newTestScaleIOServer()
	defer srv.Close()

	client := New(nil, web.Request{
		UserURL:  srv.URL,
		Username: testUser,
		Password: testPassword,
	})

	stats, err := client.SelectedStatistics(SelectedStatisticsQuery{})
	assert.Equal(t, testStatistics, stats)
	assert.NoError(t, err)
	assert.Equal(t, testStatistics, stats)
}

var (
	testUser      = "user"
	testPassword  = "password"
	testVersion   = "2.5"
	testToken     = "token"
	testInstances = Instances{
		StoragePoolList: []StoragePool{
			{ID: "id1", Name: "Marketing", SparePercentage: 10},
			{ID: "id2", Name: "Finance", SparePercentage: 10},
		},
		SdcList: []Sdc{
			{ID: "id1", SdcIp: "10.0.0.1", MdmConnectionState: "Connected"},
			{ID: "id2", SdcIp: "10.0.0.2", MdmConnectionState: "Connected"},
		},
	}
	testStatistics = SelectedStatistics{
		System:      SystemStatistics{NumOfDevices: 1},
		Sdc:         map[string]SdcStatistics{"id1": {}, "id2": {}},
		StoragePool: map[string]StoragePoolStatistics{"id1": {}, "id2": {}},
	}
)

func newTestScaleIOServer() *httptest.Server {
	return httptest.NewServer(mockScaleIOServer{})
}

type mockScaleIOServer struct{}

func (s mockScaleIOServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/api/") {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch r.URL.Path {
	default:
		w.WriteHeader(http.StatusBadRequest)
	case "/api/login":
		s.handleLogin(w, r)
	case "/api/logout":
		s.handleLogout(w, r)
	case "/api/version":
		s.handleVersion(w, r)
	case "/api/instances":
		s.handleInstances(w, r)
	case "/api/instances/querySelectedStatistics":
		s.handleQuerySelectedStatistics(w, r)
	}
}

func (mockScaleIOServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	if user, pass, ok := r.BasicAuth(); !ok || user != testUser || pass != testPassword {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, _ = w.Write([]byte(testToken))
}

func (mockScaleIOServer) handleLogout(w http.ResponseWriter, r *http.Request) {
	if _, pass, ok := r.BasicAuth(); !ok || pass != testToken {
		w.WriteHeader(http.StatusUnauthorized)
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (mockScaleIOServer) handleVersion(w http.ResponseWriter, r *http.Request) {
	if _, pass, ok := r.BasicAuth(); !ok || pass != testToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, _ = w.Write([]byte(testVersion))
}

func (mockScaleIOServer) handleInstances(w http.ResponseWriter, r *http.Request) {
	if _, pass, ok := r.BasicAuth(); !ok || pass != testToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	b, _ := json.Marshal(testInstances)
	_, _ = w.Write(b)
}

func (mockScaleIOServer) handleQuerySelectedStatistics(w http.ResponseWriter, r *http.Request) {
	if _, pass, ok := r.BasicAuth(); !ok || pass != testToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&SelectedStatisticsQuery{}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	b, _ := json.Marshal(testStatistics)
	_, _ = w.Write(b)
}

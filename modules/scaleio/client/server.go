package client

import (
	"encoding/json"
	"net/http"
	"strings"
)

// MockScaleIOAPIServer represents VxFlex OS Gateway.
type MockScaleIOAPIServer struct {
	User       string
	Password   string
	Token      string
	Version    string
	Instances  Instances
	Statistics SelectedStatistics
}

func (s MockScaleIOAPIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (s MockScaleIOAPIServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	if user, pass, ok := r.BasicAuth(); !ok || user != s.User || pass != s.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, _ = w.Write([]byte(s.Token))
}

func (s MockScaleIOAPIServer) handleLogout(w http.ResponseWriter, r *http.Request) {
	if _, pass, ok := r.BasicAuth(); !ok || pass != s.Token {
		w.WriteHeader(http.StatusUnauthorized)
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (s MockScaleIOAPIServer) handleVersion(w http.ResponseWriter, r *http.Request) {
	if _, pass, ok := r.BasicAuth(); !ok || pass != s.Token {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, _ = w.Write([]byte(s.Version))
}

func (s MockScaleIOAPIServer) handleInstances(w http.ResponseWriter, r *http.Request) {
	if _, pass, ok := r.BasicAuth(); !ok || pass != s.Token {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	b, _ := json.Marshal(s.Instances)
	_, _ = w.Write(b)
}

func (s MockScaleIOAPIServer) handleQuerySelectedStatistics(w http.ResponseWriter, r *http.Request) {
	if _, pass, ok := r.BasicAuth(); !ok || pass != s.Token {
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
	b, _ := json.Marshal(s.Statistics)
	_, _ = w.Write(b)
}

package parser

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogLine_Verify(t *testing.T) {
	parser, err := newLTSVParser(Config{
		TimeMultiplier: time.Second.Seconds(),
		LTSV: LTSVConfig{
			FieldDelimiter: ' ',
			ValueDelimiter: '=',
		},
	}, nil)
	require.NoError(t, err)
	tests := []struct {
		name string
		line string
	}{
		{"", `client=127.0.0.1 method=GET uri=/ status=200 resp_size=10`},
		{"missing some mandatory fields", ``},
		{"invalid vhost", `client=127.0.0.1 method=GET uri=/ status=200 resp_size=10 vhost=my_server.com`},
		{"invalid client", `client=my_server method=GET uri=/ status=200 resp_size=10`},
		{"invalid method", `client=127.0.0.1 method=foo uri=/ status=200 resp_size=10`},
		{"invalid URI", `client=127.0.0.1 method=GET uri=foo status=200 resp_size=10`},
		{"invalid protocol", `client=127.0.0.1 method=GET uri=/ status=200 resp_size=10 version=a`},
		{"invalid status", `client=127.0.0.1 method=GET uri=/ status=50 resp_size=10`},
		{"invalid response size", `client=127.0.0.1 method=GET uri=/ status=200 resp_size=-10`},
		{"invalid request size", `client=127.0.0.1 method=GET uri=/ status=200 resp_size=10 req_size=-10`},
		{"invalid response time", `client=127.0.0.1 method=GET uri=/ status=200 resp_size=10 resp_time=-10`},
		{"invalid upstream response time", `client=127.0.0.1 method=GET uri=/ status=200 resp_size=10 upstream_resp_time=-10`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log, err := parser.Parse([]byte(tt.line))
			require.NoError(t, err)
			err = log.Verify()
			if tt.name == "" {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.name)
			}
		})
	}
}

func Test_LogLine_assignRequest(t *testing.T) {
	tests := []struct {
		name        string
		wantMethod  string
		wantUri     string
		wantVersion string
		wantErr     bool
	}{
		{"GET / HTTP/1.1", "GET", "/", "1.1", false},
		{"GET / HTTP/1.0", "GET", "/", "1.0", false},
		{"GET / HTTP/2", "GET", "/", "2", false},
		{"GET /ihs.gif HTTP/1.1", "GET", "/ihs.gif", "1.1", false},
		{"GET no_version", "", "", "", true},
		{"GET /invalid_version http/1.1", "GET", "/invalid_version", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := emptyLogLine
			err := log.assignRequest(tt.name)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantMethod, log.Method)
				assert.Equal(t, tt.wantUri, log.URI)
				assert.Equal(t, tt.wantVersion, log.Version)
			}
		})
	}
}

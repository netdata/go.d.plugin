package weblog

import (
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
const (
	fieldVhost         = "vhost"
	fieldPort          = "port"
	fieldVhostWithPort = "vhost_port"
	fieldReqScheme     = "req_scheme"
	fieldReqClient     = "req_client"
	fieldRequest       = "request"
	fieldReqMethod     = "req_method"
	fieldReqURI        = "req_uri"
	fieldReqProto      = "req_proto"
	fieldReqSize       = "req_size"
	fieldRespStatus    = "resp_status"
	fieldRespSize      = "resp_size"
	fieldRespTime      = "resp_time"
	fieldUpsRespTime   = "ups_resp_time"
	fieldCustom        = "custom"
)
*/

func TestLogLine_Assign(t *testing.T) {
	type testCase struct {
		v    string
		line string
		err  error
	}
	tests := []struct {
		name  string
		vars  []string
		cases []testCase
	}{
		{
			"vhost",
			[]string{"host", "http_host", "v"},
			[]testCase{
				{v: "1.1.1.1", line: "vhost=1.1.1.1"},
				{},
			},
		},
		{
			"vhost with port",
			[]string{"host:$server_port", "v:%p"},
			[]testCase{
				{v: "1.1.1.1:80", line: "vhost=1.1.1.1 port=80"},
				{v: "1.1.1.1:", line: "vhost=1.1.1.1 port="},
				{v: "debian10.debian:81", line: "vhost=debian10.debian port=81"},
				{},
				{v: "1.1.1.1", err: errBadVhostPort},
				{v: "1.1.1.1 80", err: errBadVhostPort},
				{v: "invalid", err: errBadVhostPort},
			},
		},
		{
			"request",
			[]string{"request", "r"},
			[]testCase{
				{v: "GET / HTTP/1.1", line: "req_method=GET req_uri=/ req_proto=1.1"},
				{v: "OPTIONS / HTTP/1.0", line: "req_method=OPTIONS req_uri=/ req_proto=1.0"},
				{v: "HEAD / HTTP/2", line: "req_method=HEAD req_uri=/ req_proto=2"},
				{v: "POST /ihs.gif HTTP/1.1", line: "req_method=POST req_uri=/ihs.gif req_proto=1.1"},
				{},
				{v: "GET no_version", err: errBadRequest},
				{v: "GOT / HTTP/2", err: errBadRequest},
				{v: "GET /invalid_version http/1.1", err: errBadRequest},
			},
		},
		{
			"request protocol",
			[]string{"server_protocol", "H"},
			[]testCase{
				{v: "HTTP/1.0", line: "req_proto=1.0"},
				{v: "HTTP/1.1", line: "req_proto=1.1"},
				{v: "HTTP/2", line: "req_proto=2"},
				{},
				{v: "1.1", err: errBadReqProto},
				{v: "http/1.1", err: errBadReqProto},
			},
		},
		{
			"response status",
			[]string{"status", "s", ">s"},
			[]testCase{
				{v: "-"},
				{v: "200", line: "resp_status=200"},
				{},
				{v: "200 ", err: errBadRespStatus},
				{v: "0.222", err: errBadRespStatus},
				{v: "localhost", err: errBadRespStatus},
			},
		},
		{
			"request size",
			[]string{"request_length", "I"},
			[]testCase{
				{v: "-", line: "req_size=0"},
				{v: "1000", line: "req_size=1000"},
				{v: "100.222", err: errBadReqSize},
				{v: "number", err: errBadReqSize},
			},
		},
		{
			"response size",
			[]string{"bytes_sent", "body_bytes_sent", "O", "B"},
			[]testCase{
				{v: "-", line: "resp_size=0"},
				{v: "1000", line: "resp_size=1000"},
				{v: "100.222", err: errBadRespSize},
				{v: "number", err: errBadRespSize},
			},
		},
		{
			"response time",
			[]string{"request_time", "D"},
			[]testCase{
				{v: "-"},
				{v: "100222", line: "resp_time=100222"},
				{v: "100.222", line: "resp_time=100222000"},
				{v: "0.333,0.444,0.555", err: errBadRespTime},
				{v: "number", err: errBadRespTime},
			},
		},
		{
			"upstream response time",
			[]string{"upstream_response_time"},
			[]testCase{
				{v: "-"},
				{v: "100222", line: "ups_resp_time=100222"},
				{v: "100.222", line: "ups_resp_time=100222000"},
				{v: "0.333,0.444,0.555", line: "ups_resp_time=333000"},
				{v: "number", err: errBadUpstreamRespTime},
			},
		},
	}

	for _, tt := range tests {
		for _, v := range tt.vars {
			for i, c := range tt.cases {
				name := fmt.Sprintf("name=%s|var='%s'|v='%s'(%d)", tt.name, v, c.v, i+1)

				t.Run(name, func(t *testing.T) {
					line := newEmptyLogLine()
					err := line.Assign(v, c.v)

					if c.err != nil {
						require.Error(t, err)
						assert.True(t, errors.Is(err, c.err))
					} else {
						require.NoError(t, err)
					}

					expected := prepareLogLine(t, c.line)
					assert.Equal(t, expected, *line)
				})
			}
		}
	}
}

func TestLogLine_verify(t *testing.T) {
	type testCase struct {
		line string
		err  error
	}
	tests := []struct {
		name  string
		cases []testCase
	}{
		{
			"vhost",
			[]testCase{
				{line: "vhost=192.168.0.1"},
				{line: "vhost=debian10.debian"},
				{line: "vhost=1ce:1ce::babe"},
				{line: "vhost=localhost"},
				{line: "vhost="},
				{},
				{line: "\"vhost=localhost \"", err: errBadVhost},
				{line: "vhost=invalid_vhost", err: errBadVhost},
				{line: "vhost=http://192.168.0.1/", err: errBadVhost},
			},
		},
		{
			"port",
			[]testCase{
				{line: "port=80"},
				{line: "port=8081"},
				{line: "port="},
				{},
				{line: "\"port=80 \"", err: errBadPort},
				{line: "port=79", err: errBadPort},
				{line: "port=0.0.0.0", err: errBadPort},
			},
		},
		{
			"request scheme",
			[]testCase{
				{line: "req_scheme=http"},
				{line: "req_scheme=https"},
				{line: "req_scheme="},
				{},
				{line: "req_scheme=not_https", err: errBadReqScheme},
				{line: "req_scheme=HTTP", err: errBadReqScheme},
				{line: "req_scheme=HTTPS", err: errBadReqScheme},
				{line: "req_scheme=10", err: errBadReqScheme},
			},
		},
		{
			"request method",
			[]testCase{
				{line: "req_method=GET"},
				{line: "req_method=POST"},
				{line: "req_method=TRACE"},
				{line: "req_method=OPTIONS"},
				{line: "req_method=CONNECT"},
				{line: "req_method=DELETE"},
				{line: "req_method=PUT"},
				{line: "req_method=PATCH"},
				{line: "req_method=HEAD"},
				{line: "req_method="},
				{},
				{line: "req_method=Get", err: errBadReqMethod},
				{line: "req_method=get", err: errBadReqMethod},
				{line: "req_method=-", err: errBadReqMethod},
			},
		},
	}

	for _, tt := range tests {
		for i, c := range tt.cases {
			name := fmt.Sprintf("name=%s|line='%s'(%d)", tt.name, c.line, i+1)

			t.Run(name, func(t *testing.T) {
				line := prepareLogLine(t, c.line)
				if name != "resp_status" {
					line.respStatus = 200
				}
				err := line.verify()

				if c.err != nil {
					require.Error(t, err)
					assert.True(t, errors.Is(err, c.err))
				} else {
					require.NoError(t, err)
				}
			})
		}
	}
}

func Test_logLine_verify_reqURI(t *testing.T) {
	tests := []struct {
		uri     string
		wantErr bool
	}{
		{"/", false},
		{"/respStatus?full&json", false},
		{"/icons/openlogo-75.png", false},
		{"", false},
		{emptyString, false},

		{"respStatus?full&json", true},
		{"localhost ", true},
		{"http://192.168.0.1/", true},
	}

	line := newEmptyLogLine()
	line.respStatus = 200

	for _, tt := range tests {
		t.Run(tt.uri, func(t *testing.T) {
			line.reqURI = tt.uri
			err := line.verify()

			if tt.wantErr {
				require.True(t, errors.Is(err, errBadReqURI))
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_logLine_verify_reqProto(t *testing.T) {
	tests := []struct {
		proto   string
		wantErr bool
	}{
		{"1", false},
		{"1.0", false},
		{"1.1", false},
		{"2.0", false},
		{"2", false},
		{"", false},
		{emptyString, false},

		{"0.9", true},
		{"1 ", true},
		{"1.1.1", true},
		{"2.2", true},
		{"localhost", true},
	}

	line := newEmptyLogLine()
	line.respStatus = 200

	for _, tt := range tests {
		t.Run(tt.proto, func(t *testing.T) {
			line.reqProto = tt.proto
			err := line.verify()

			if tt.wantErr {
				require.True(t, errors.Is(err, errBadReqProto))
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_logLine_verify_reqSize(t *testing.T) {
	tests := []struct {
		size    int
		wantErr bool
	}{
		{0, false},
		{100, false},
		{100000, false},
		{emptyNumber, false},

		{-1, true},
	}

	line := newEmptyLogLine()
	line.respStatus = 200

	for _, tt := range tests {
		name := strconv.Itoa(tt.size)
		t.Run(name, func(t *testing.T) {
			line.reqSize = tt.size
			err := line.verify()

			if tt.wantErr {
				require.True(t, errors.Is(err, errBadReqSize))
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_logLine_verify_respStatus(t *testing.T) {
	tests := []struct {
		status  int
		wantErr error
	}{
		{100, nil},
		{200, nil},
		{300, nil},
		{400, nil},
		{500, nil},
		{600, nil},

		{emptyNumber, errMandatoryField},
		{-1, errBadRespStatus},
		{99, errBadRespStatus},
		{601, errBadRespStatus},
	}

	line := newEmptyLogLine()

	for _, tt := range tests {
		name := strconv.Itoa(tt.status)
		t.Run(name, func(t *testing.T) {
			line.respStatus = tt.status
			err := line.verify()

			if tt.wantErr != nil {
				require.True(t, errors.Is(err, tt.wantErr))
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_logLine_verify_respSize(t *testing.T) {
	tests := []struct {
		size    int
		wantErr bool
	}{
		{0, false},
		{100, false},
		{1e9, false},
		{emptyNumber, false},

		{-1, true},
	}

	line := newEmptyLogLine()
	line.respStatus = 200

	for _, tt := range tests {
		name := strconv.Itoa(tt.size)
		t.Run(name, func(t *testing.T) {
			line.respSize = tt.size
			err := line.verify()

			if tt.wantErr {
				require.True(t, errors.Is(err, errBadRespSize))
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_logLine_verify_respTime(t *testing.T) {
	tests := []struct {
		time    float64
		wantErr bool
	}{
		{0, false},
		{100, false},
		{1e9, false},
		{0.000, false},
		{100.123, false},
		{emptyNumber, false},

		{-1, true},
	}

	line := newEmptyLogLine()
	line.respStatus = 200

	for _, tt := range tests {
		name := fmt.Sprintf("%f", tt.time)
		t.Run(name, func(t *testing.T) {
			line.respTime = tt.time
			err := line.verify()

			if tt.wantErr {
				require.True(t, errors.Is(err, errBadRespTime))
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_logLine_verify_upstreamRespTime(t *testing.T) {
	tests := []struct {
		time    float64
		wantErr bool
	}{
		{0, false},
		{100, false},
		{1e9, false},
		{0.000, false},
		{100.123, false},
		{emptyNumber, false},

		{-1, true},
	}

	line := newEmptyLogLine()
	line.respStatus = 200

	for _, tt := range tests {
		name := fmt.Sprintf("%f", tt.time)
		t.Run(name, func(t *testing.T) {
			line.upsRespTime = tt.time
			err := line.verify()

			if tt.wantErr {
				require.True(t, errors.Is(err, errBadUpstreamRespTime))
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func prepareLogLine(t *testing.T, from string) logLine {
	if from == "" {
		return *newEmptyLogLine()
	}

	t.Helper()
	r := csv.NewReader(strings.NewReader(from))
	r.Comma = ' '
	line := newEmptyLogLine()

	rs, err := r.Read()
	require.NoError(t, err)

	for _, v := range rs {
		parts := strings.Split(v, "=")
		require.Len(t, parts, 2)
		field, val := parts[0], parts[1]
		switch field {
		case fieldVhost:
			line.vhost = val
		case fieldPort:
			line.port = val
		case fieldReqScheme:
			line.reqScheme = val
		case fieldReqClient:
			line.reqClient = val
		case fieldReqMethod:
			line.reqMethod = val
		case fieldReqURI:
			line.reqURI = val
		case fieldReqProto:
			line.reqProto = val
		case fieldReqSize:
			i, err := strconv.Atoi(val)
			require.NoError(t, err)
			line.reqSize = i
		case fieldRespStatus:
			i, err := strconv.Atoi(val)
			require.NoError(t, err)
			line.respStatus = i
		case fieldRespSize:
			i, err := strconv.Atoi(val)
			require.NoError(t, err)
			line.respSize = i
		case fieldRespTime:
			i, err := strconv.ParseFloat(val, 64)
			require.NoError(t, err)
			line.respTime = i
		case fieldUpsRespTime:
			i, err := strconv.ParseFloat(val, 64)
			require.NoError(t, err)
			line.upsRespTime = i
		case fieldCustom:
			line.custom = val
		default:
			t.Fatalf("cant prepare logLine, unknown field: %s", field)
		}
	}
	return *line
}

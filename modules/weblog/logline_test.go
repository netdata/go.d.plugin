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

const (
	emptyLogLine = "emptyLogLine"
	emptyStr     = ""
	hyphen       = "-"
)

func TestLogLine_Assign(t *testing.T) {
	type testCase struct {
		v        string
		wantLine string
		wantErr  error
	}
	tests := []struct {
		name  string
		vars  []string
		cases []testCase
	}{
		{
			"Vhost",
			[]string{
				"host",
				"http_host",
				"v",
			},
			[]testCase{
				{v: "1.1.1.1", wantLine: "vhost=1.1.1.1"},
				{v: "::1", wantLine: "vhost=::1"},
				{v: "[::1]", wantLine: "vhost=::1"},
				{v: "1ce:1ce::babe", wantLine: "vhost=1ce:1ce::babe"},
				{v: "[1ce:1ce::babe]", wantLine: "vhost=1ce:1ce::babe"},
				{v: "localhost", wantLine: "vhost=localhost"},
				{v: "debian10.debian", wantLine: "vhost=debian10.debian"},
				{v: "my_vhost", wantLine: "vhost=my_vhost"},
				{v: emptyStr, wantLine: emptyLogLine},
				{v: hyphen, wantLine: emptyLogLine},
			},
		},
		{
			"Port",
			[]string{
				"server_port",
				"p",
			},
			[]testCase{
				{v: "80", wantLine: "port=80"},
				{v: "8081", wantLine: "port=8081"},
				{v: "30000", wantLine: "port=30000"},
				{v: emptyStr, wantLine: emptyLogLine},
				{v: hyphen, wantLine: emptyLogLine},
				{v: "-1", wantLine: emptyLogLine, wantErr: errBadPort},
				{v: "0", wantLine: emptyLogLine, wantErr: errBadPort},
				{v: "50000", wantLine: emptyLogLine, wantErr: errBadPort},
			},
		},
		{
			"Vhost with port",
			[]string{
				"host:$server_port",
				"v:%p",
			},
			[]testCase{
				{v: "1.1.1.1:80", wantLine: "vhost=1.1.1.1 port=80"},
				{v: "::1:80", wantLine: "vhost=::1 port=80"},
				{v: "[::1]:80", wantLine: "vhost=::1 port=80"},
				{v: "1ce:1ce::babe:80", wantLine: "vhost=1ce:1ce::babe port=80"},
				{v: "debian10.debian:81", wantLine: "vhost=debian10.debian port=81"},
				{v: emptyStr, wantLine: emptyLogLine},
				{v: hyphen, wantLine: emptyLogLine},
				{v: "1.1.1.1", wantLine: emptyLogLine, wantErr: errBadVhostPort},
				{v: "1.1.1.1:", wantLine: emptyLogLine, wantErr: errBadVhostPort},
				{v: "1.1.1.1 80", wantLine: emptyLogLine, wantErr: errBadVhostPort},
				{v: "1.1.1.1:20", wantLine: emptyLogLine, wantErr: errBadVhostPort},
				{v: "1.1.1.1:50000", wantLine: emptyLogLine, wantErr: errBadVhostPort},
			},
		},
		{
			"Scheme",
			[]string{
				"scheme",
			},
			[]testCase{
				{v: "http", wantLine: "req_scheme=http"},
				{v: "https", wantLine: "req_scheme=https"},
				{v: emptyStr, wantLine: emptyLogLine},
				{v: hyphen, wantLine: emptyLogLine},
				{v: "HTTP", wantLine: emptyLogLine, wantErr: errBadReqScheme},
				{v: "HTTPS", wantLine: emptyLogLine, wantErr: errBadReqScheme},
			},
		},
		{
			"Client",
			[]string{
				"remote_addr",
				"a",
				"h",
			},
			[]testCase{
				{v: "1.1.1.1", wantLine: "req_client=1.1.1.1"},
				{v: "debian10", wantLine: "req_client=debian10"},
				{v: emptyStr, wantLine: emptyLogLine},
				{v: hyphen, wantLine: emptyLogLine},
			},
		},
		{
			"Request",
			[]string{
				"request",
				"r",
			},
			[]testCase{
				{v: "GET / HTTP/1.0", wantLine: "req_method=GET req_url=/ req_proto=1.0"},
				{v: "HEAD /ihs.gif HTTP/1.0", wantLine: "req_method=HEAD req_url=/ihs.gif req_proto=1.0"},
				{v: "POST /ihs.gif HTTP/1.0", wantLine: "req_method=POST req_url=/ihs.gif req_proto=1.0"},
				{v: "PUT /ihs.gif HTTP/1.0", wantLine: "req_method=PUT req_url=/ihs.gif req_proto=1.0"},
				{v: "PATCH /ihs.gif HTTP/1.0", wantLine: "req_method=PATCH req_url=/ihs.gif req_proto=1.0"},
				{v: "DELETE /ihs.gif HTTP/1.0", wantLine: "req_method=DELETE req_url=/ihs.gif req_proto=1.0"},
				{v: "OPTIONS /ihs.gif HTTP/1.0", wantLine: "req_method=OPTIONS req_url=/ihs.gif req_proto=1.0"},
				{v: "TRACE /ihs.gif HTTP/1.0", wantLine: "req_method=TRACE req_url=/ihs.gif req_proto=1.0"},
				{v: "CONNECT ip.cn:443 HTTP/1.1", wantLine: "req_method=CONNECT req_url=ip.cn:443 req_proto=1.1"},
				{v: "GET / HTTP/1.1", wantLine: "req_method=GET req_url=/ req_proto=1.1"},
				{v: "GET / HTTP/2", wantLine: "req_method=GET req_url=/ req_proto=2"},
				{v: "GET / HTTP/2.0", wantLine: "req_method=GET req_url=/ req_proto=2.0"},
				{v: "GET /invalid_version http/1.1", wantLine: "req_method=GET req_url=/invalid_version", wantErr: errBadRequest},
				{v: emptyStr, wantLine: emptyLogLine},
				{v: hyphen, wantLine: emptyLogLine},
				{v: "GET no_version", wantLine: emptyLogLine, wantErr: errBadRequest},
				{v: "GOT / HTTP/2", wantLine: emptyLogLine, wantErr: errBadRequest},
				{v: "get / HTTP/2", wantLine: emptyLogLine, wantErr: errBadRequest},
				{v: "x04\x01\x00P$3\xFE\xEA\x00", wantLine: emptyLogLine, wantErr: errBadRequest},
			},
		},
		{
			"Method",
			[]string{
				"request_method",
				"m",
			},
			[]testCase{
				{v: "GET", wantLine: "req_method=GET"},
				{v: "HEAD", wantLine: "req_method=HEAD"},
				{v: "POST", wantLine: "req_method=POST"},
				{v: "PUT", wantLine: "req_method=PUT"},
				{v: "PATCH", wantLine: "req_method=PATCH"},
				{v: "DELETE", wantLine: "req_method=DELETE"},
				{v: "OPTIONS", wantLine: "req_method=OPTIONS"},
				{v: "TRACE", wantLine: "req_method=TRACE"},
				{v: "CONNECT", wantLine: "req_method=CONNECT"},
				{v: emptyStr, wantLine: emptyLogLine},
				{v: hyphen, wantLine: emptyLogLine},
				{v: "GET no_version", wantLine: emptyLogLine, wantErr: errBadReqMethod},
				{v: "GOT / HTTP/2", wantLine: emptyLogLine, wantErr: errBadReqMethod},
				{v: "get / HTTP/2", wantLine: emptyLogLine, wantErr: errBadReqMethod},
			},
		},
		{
			"URL",
			[]string{
				"request_uri",
				"U",
			},
			[]testCase{
				{v: "/server-status?auto", wantLine: "req_url=/server-status?auto"},
				{v: "/default.html", wantLine: "req_url=/default.html"},
				{v: "10.0.0.1:3128", wantLine: "req_url=10.0.0.1:3128"},
				{v: emptyStr, wantLine: emptyLogLine},
				{v: hyphen, wantLine: emptyLogLine},
			},
		},
		{
			"Protocol",
			[]string{
				"server_protocol",
				"H",
			},
			[]testCase{
				{v: "HTTP/1.0", wantLine: "req_proto=1.0"},
				{v: "HTTP/1.1", wantLine: "req_proto=1.1"},
				{v: "HTTP/2", wantLine: "req_proto=2"},
				{v: "HTTP/2.0", wantLine: "req_proto=2.0"},
				{v: emptyStr, wantLine: emptyLogLine},
				{v: hyphen, wantLine: emptyLogLine},
				{v: "1.1", wantLine: emptyLogLine, wantErr: errBadReqProto},
				{v: "http/1.1", wantLine: emptyLogLine, wantErr: errBadReqProto},
			},
		},
		{
			"Status",
			[]string{
				"status",
				"s",
				">s",
			},
			[]testCase{
				{v: "100", wantLine: "resp_status=100"},
				{v: "200", wantLine: "resp_status=200"},
				{v: "300", wantLine: "resp_status=300"},
				{v: "400", wantLine: "resp_status=400"},
				{v: "500", wantLine: "resp_status=500"},
				{v: "600", wantLine: "resp_status=600"},
				{v: emptyStr, wantLine: emptyLogLine},
				{v: hyphen, wantLine: emptyLogLine},
				{v: "99", wantLine: emptyLogLine, wantErr: errBadRespStatus},
				{v: "601", wantLine: emptyLogLine, wantErr: errBadRespStatus},
				{v: "200 ", wantLine: emptyLogLine, wantErr: errBadRespStatus},
				{v: "0.222", wantLine: emptyLogLine, wantErr: errBadRespStatus},
				{v: "localhost", wantLine: emptyLogLine, wantErr: errBadRespStatus},
			},
		},
		{
			"Request size",
			[]string{
				"request_length",
				"I",
			},
			[]testCase{
				{v: "15", wantLine: "req_size=15"},
				{v: "1000000", wantLine: "req_size=1000000"},
				{v: emptyStr, wantLine: emptyLogLine},
				{v: hyphen, wantLine: emptyLogLine},
				{v: "-1", wantLine: emptyLogLine, wantErr: errBadReqSize},
				{v: "100.222", wantLine: emptyLogLine, wantErr: errBadReqSize},
				{v: "invalid", wantLine: emptyLogLine, wantErr: errBadReqSize},
			},
		},
		{
			"Response size",
			[]string{
				"bytes_sent",
				"body_bytes_sent",
				"O",
				"B",
			},
			[]testCase{
				{v: "15", wantLine: "resp_size=15"},
				{v: "1000000", wantLine: "resp_size=1000000"},
				{v: emptyStr, wantLine: emptyLogLine},
				{v: hyphen, wantLine: "resp_size=0"},
				{v: "-1", wantLine: emptyLogLine, wantErr: errBadRespSize},
				{v: "100.222", wantLine: emptyLogLine, wantErr: errBadRespSize},
				{v: "invalid", wantLine: emptyLogLine, wantErr: errBadRespSize},
			},
		},
		{
			"Response time",
			[]string{
				"request_time",
				"D",
			},
			[]testCase{
				{v: "100222", wantLine: "resp_time=100222"},
				{v: "100.222", wantLine: "resp_time=100222000"},
				{v: emptyStr, wantLine: emptyLogLine},
				{v: hyphen, wantLine: emptyLogLine},
				{v: "-1", wantLine: emptyLogLine, wantErr: errBadRespTime},
				{v: "0.333,0.444,0.555", wantLine: emptyLogLine, wantErr: errBadRespTime},
				{v: "number", wantLine: emptyLogLine, wantErr: errBadRespTime},
			},
		},
		{
			"Upstream response time",
			[]string{
				"upstream_response_time",
			},
			[]testCase{
				{v: "100222", wantLine: "ups_resp_time=100222"},
				{v: "100.222", wantLine: "ups_resp_time=100222000"},
				{v: "0.333,0.444,0.555", wantLine: "ups_resp_time=333000"},
				{v: emptyStr, wantLine: emptyLogLine},
				{v: hyphen, wantLine: emptyLogLine},
				{v: "-1", wantLine: emptyLogLine, wantErr: errBadUpstreamRespTime},
				{v: "number", wantLine: emptyLogLine, wantErr: errBadUpstreamRespTime},
			},
		},
	}

	for _, tt := range tests {
		for _, varName := range tt.vars {
			for i, tc := range tt.cases {
				name := fmt.Sprintf("[%s:%d]var='%s'|val='%s'", tt.name, i+1, varName, tc.v)
				t.Run(name, func(t *testing.T) {

					line := newEmptyLogLine()
					err := line.Assign(varName, tc.v)

					if tc.wantErr != nil {
						require.Error(t, err)
						assert.True(t, errors.Is(err, tc.wantErr))
					} else {
						require.NoError(t, err)
					}

					expected := prepareLogLine(t, tc.wantLine)
					assert.Equal(t, expected, *line)
				})
			}
		}
	}
}

func TestLogLine_verify(t *testing.T) {
	type testCase struct {
		line    string
		wantErr error
	}
	tests := []struct {
		name  string
		cases []testCase
	}{
		{
			"Vhost",
			[]testCase{
				{line: "vhost=192.168.0.1"},
				{line: "vhost=debian10.debian"},
				{line: "vhost=1ce:1ce::babe"},
				{line: "vhost=localhost"},
				{line: emptyLogLine},
				{line: "vhost=invalid_vhost", wantErr: errBadVhost},
				{line: "vhost=http://192.168.0.1/", wantErr: errBadVhost},
			},
		},
		{
			"Port",
			[]testCase{
				{line: "port=80"},
				{line: "port=8081"},
				{line: emptyLogLine},
				{line: "port=79", wantErr: errBadPort},
				{line: "port=50000", wantErr: errBadPort},
				{line: "port=0.0.0.0", wantErr: errBadPort},
			},
		},
		{
			"Vhost with port",
			[]testCase{
				{line: "vhost=1.1.1.1 port=80"},
				{line: "vhost=::1 port=8081"},
				{line: "vhost=1ce:1ce::babe port=81"},
				{line: "vhost=debian10.debian port=81"},
				{line: emptyLogLine},
				{line: "vhost=::1 port=79", wantErr: errBadPort},
				{line: "vhost=::1 port=50000", wantErr: errBadPort},
				{line: "vhost=::1 port=0.0.0.0", wantErr: errBadPort},
			},
		},
		{
			"Scheme",
			[]testCase{
				{line: "req_scheme=http"},
				{line: "req_scheme=https"},
				{line: emptyLogLine},
				{line: "req_scheme=not_https", wantErr: errBadReqScheme},
				{line: "req_scheme=HTTP", wantErr: errBadReqScheme},
				{line: "req_scheme=HTTPS", wantErr: errBadReqScheme},
				{line: "req_scheme=10", wantErr: errBadReqScheme},
			},
		},
		{
			"Client",
			[]testCase{
				{line: "req_client=1.1.1.1"},
				{line: "req_client=::1"},
				{line: "req_client=1ce:1ce::babe"},
				{line: "req_client=localhost"},
				{line: emptyLogLine},
				{line: "req_client=debian10.debian", wantErr: errBadReqClient},
				{line: "req_client=invalid", wantErr: errBadReqClient},
			},
		},
		{
			"Method",
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
				{line: emptyLogLine},
				{line: "req_method=Get", wantErr: errBadReqMethod},
				{line: "req_method=get", wantErr: errBadReqMethod},
				{line: "req_method=-", wantErr: errBadReqMethod},
			},
		},
		{
			"URL",
			[]testCase{
				{line: "req_url=/"},
				{line: "req_url=/status?full&json"},
				{line: "req_url=/icons/openlogo-75.png"},
				{line: "req_method=CONNECT req_url=http://192.168.0.1/"},
				{line: emptyLogLine},
				{line: "req_url=status?full&json", wantErr: errBadReqURI},
				{line: "\"req_url=/ \"", wantErr: errBadReqURI},
				{line: "req_url=http://192.168.0.1/", wantErr: errBadReqURI},
			},
		},
		{
			"Protocol",
			[]testCase{
				{line: "req_proto=1"},
				{line: "req_proto=1.0"},
				{line: "req_proto=1.1"},
				{line: "req_proto=2.0"},
				{line: "req_proto=2"},
				{line: emptyLogLine},
				{line: "req_proto=0.9", wantErr: errBadReqProto},
				{line: "req_proto=1.1.1", wantErr: errBadReqProto},
				{line: "req_proto=2.2", wantErr: errBadReqProto},
				{line: "req_proto=localhost", wantErr: errBadReqProto},
			},
		},
		{
			"Status",
			[]testCase{
				{line: "resp_status=100"},
				{line: "resp_status=200"},
				{line: "resp_status=300"},
				{line: "resp_status=400"},
				{line: "resp_status=500"},
				{line: "resp_status=600"},
				{line: emptyLogLine, wantErr: errMandatoryField},
				{line: "resp_status=-1", wantErr: errBadRespStatus},
				{line: "resp_status=99", wantErr: errBadRespStatus},
				{line: "resp_status=601", wantErr: errBadRespStatus},
			},
		},
		{
			"Request size",
			[]testCase{
				{line: "req_size=0"},
				{line: "req_size=100"},
				{line: "req_size=1000000"},
				{line: emptyLogLine},
				{line: "req_size=-1", wantErr: errBadReqSize},
			},
		},
		{
			"Response size",
			[]testCase{
				{line: "resp_size=0"},
				{line: "resp_size=100"},
				{line: "resp_size=1000000"},
				{line: emptyLogLine},
				{line: "resp_size=-1", wantErr: errBadRespSize},
			},
		},
		{
			"Response time",
			[]testCase{
				{line: "resp_time=0"},
				{line: "resp_time=0.000"},
				{line: "resp_time=100"},
				{line: "resp_time=1000000.123"},
				{line: emptyLogLine},
				{line: "resp_time=-1", wantErr: errBadRespTime},
			},
		},
		{
			"Upstream response time",
			[]testCase{
				{line: "ups_resp_time=0"},
				{line: "ups_resp_time=0.000"},
				{line: "ups_resp_time=100"},
				{line: "ups_resp_time=1000000.123"},
				{line: emptyLogLine},
				{line: "ups_resp_time=-1", wantErr: errBadUpstreamRespTime},
			},
		},
	}

	for _, tt := range tests {
		for i, c := range tt.cases {
			name := fmt.Sprintf("name=%s|wantLine='%s'(%d)", tt.name, c.line, i+1)

			t.Run(name, func(t *testing.T) {
				line := prepareLogLine(t, c.line)
				if tt.name != "Status" {
					line.respStatus = 200
				}

				err := line.verify()

				if c.wantErr != nil {
					require.Error(t, err)
					assert.True(t, errors.Is(err, c.wantErr))
				} else {
					require.NoError(t, err)
				}
			})
		}
	}
}

func prepareLogLine(t *testing.T, from string) logLine {
	require.NotZero(t, from)

	if from == emptyLogLine {
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
		case fieldReqURL:
			line.reqURL = val
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

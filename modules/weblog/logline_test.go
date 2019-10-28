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
	zeroLogLine = "zeroLogLine"
	emptyStr    = ""
	hyphen      = "-"
)

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
			"Vhost",
			[]string{
				"host",
				"http_host",
				"v",
			},
			[]testCase{
				{v: "1.1.1.1", line: "vhost=1.1.1.1"},
				{v: "::1", line: "vhost=::1"},
				{v: "[::1]", line: "vhost=::1"},
				{v: "1ce:1ce::babe", line: "vhost=1ce:1ce::babe"},
				{v: "[1ce:1ce::babe]", line: "vhost=1ce:1ce::babe"},
				{v: "localhost", line: "vhost=localhost"},
				{v: "debian10.debian", line: "vhost=debian10.debian"},
				{v: "my_vhost", line: "vhost=my_vhost"},
				{v: emptyStr, line: zeroLogLine},
				{v: hyphen, line: zeroLogLine},
			},
		},
		{
			"Port",
			[]string{
				"server_port",
				"p",
			},
			[]testCase{
				{v: "80", line: "port=80"},
				{v: "8081", line: "port=8081"},
				{v: "30000", line: "port=30000"},
				{v: emptyStr, line: zeroLogLine},
				{v: hyphen, line: zeroLogLine},
				{v: "-1", err: errBadPort},
				{v: "0", err: errBadPort},
				{v: "50000", err: errBadPort},
			},
		},
		{
			"Vhost with port",
			[]string{
				"host:$server_port",
				"v:%p",
			},
			[]testCase{
				{v: "1.1.1.1:80", line: "vhost=1.1.1.1 port=80"},
				{v: "::1:80", line: "vhost=::1 port=80"},
				{v: "[::1]:80", line: "vhost=::1 port=80"},
				{v: "1ce:1ce::babe:80", line: "vhost=1ce:1ce::babe port=80"},
				{v: "debian10.debian:81", line: "vhost=debian10.debian port=81"},
				{v: emptyStr, line: zeroLogLine},
				{v: hyphen, line: zeroLogLine},
				{v: "1.1.1.1", err: errBadVhostPort},
				{v: "1.1.1.1:", err: errBadVhostPort},
				{v: "1.1.1.1 80", err: errBadVhostPort},
				{v: "1.1.1.1:20", err: errBadVhostPort},
				{v: "1.1.1.1:50000", err: errBadVhostPort},
			},
		},
		{
			"Scheme",
			[]string{
				"scheme",
			},
			[]testCase{
				{v: "http", line: "req_scheme=http"},
				{v: "https", line: "req_scheme=https"},
				{v: emptyStr, line: zeroLogLine},
				{v: hyphen, line: zeroLogLine},
				{v: "HTTP", err: errBadReqScheme},
				{v: "HTTPS", err: errBadReqScheme},
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
				{v: "1.1.1.1", line: "req_client=1.1.1.1"},
				{v: "debian10", line: "req_client=debian10"},
				{v: emptyStr, line: zeroLogLine},
				{v: hyphen, line: zeroLogLine},
			},
		},
		{
			"Request",
			[]string{
				"request",
				"r",
			},
			[]testCase{
				{v: "GET / HTTP/1.0", line: "req_method=GET req_url=/ req_proto=1.0"},
				{v: "HEAD /ihs.gif HTTP/1.0", line: "req_method=HEAD req_url=/ihs.gif req_proto=1.0"},
				{v: "POST /ihs.gif HTTP/1.0", line: "req_method=POST req_url=/ihs.gif req_proto=1.0"},
				{v: "PUT /ihs.gif HTTP/1.0", line: "req_method=PUT req_url=/ihs.gif req_proto=1.0"},
				{v: "PATCH /ihs.gif HTTP/1.0", line: "req_method=PATCH req_url=/ihs.gif req_proto=1.0"},
				{v: "DELETE /ihs.gif HTTP/1.0", line: "req_method=DELETE req_url=/ihs.gif req_proto=1.0"},
				{v: "OPTIONS /ihs.gif HTTP/1.0", line: "req_method=OPTIONS req_url=/ihs.gif req_proto=1.0"},
				{v: "TRACE /ihs.gif HTTP/1.0", line: "req_method=TRACE req_url=/ihs.gif req_proto=1.0"},
				{v: "CONNECT ip.cn:443 HTTP/1.1", line: "req_method=CONNECT req_url=ip.cn:443 req_proto=1.1"},
				{v: "GET / HTTP/1.1", line: "req_method=GET req_url=/ req_proto=1.1"},
				{v: "GET / HTTP/2", line: "req_method=GET req_url=/ req_proto=2"},
				{v: "GET / HTTP/2.0", line: "req_method=GET req_url=/ req_proto=2.0"},
				{v: "GET /invalid_version http/1.1", line: "req_method=GET req_url=/invalid_version", err: errBadRequest},
				{v: emptyStr, line: zeroLogLine},
				{v: hyphen, line: zeroLogLine},
				{v: "GET no_version", err: errBadRequest},
				{v: "GOT / HTTP/2", err: errBadRequest},
				{v: "get / HTTP/2", err: errBadRequest},
				{v: "x04\x01\x00P$3\xFE\xEA\x00", err: errBadRequest},
			},
		},
		{
			"Request method",
			[]string{
				"request_method",
				"m",
			},
			[]testCase{
				{v: "GET", line: "req_method=GET"},
				{v: "HEAD", line: "req_method=HEAD"},
				{v: "POST", line: "req_method=POST"},
				{v: "PUT", line: "req_method=PUT"},
				{v: "PATCH", line: "req_method=PATCH"},
				{v: "DELETE", line: "req_method=DELETE"},
				{v: "OPTIONS", line: "req_method=OPTIONS"},
				{v: "TRACE", line: "req_method=TRACE"},
				{v: "CONNECT", line: "req_method=CONNECT"},
				{v: emptyStr, line: zeroLogLine},
				{v: hyphen, line: zeroLogLine},
				{v: "GET no_version", err: errBadReqMethod},
				{v: "GOT / HTTP/2", err: errBadReqMethod},
				{v: "get / HTTP/2", err: errBadReqMethod},
			},
		},
		{
			"Request url",
			[]string{
				"request_uri",
				"U",
			},
			[]testCase{
				{v: "/server-status?auto", line: "req_url=/server-status?auto"},
				{v: "/default.html", line: "req_url=/default.html"},
				{v: "10.0.0.1:3128", line: "req_url=10.0.0.1:3128"},
				{v: emptyStr, line: zeroLogLine},
				{v: hyphen, line: zeroLogLine},
			},
		},
		{
			"Request protocol",
			[]string{
				"server_protocol",
				"H",
			},
			[]testCase{
				{v: "HTTP/1.0", line: "req_proto=1.0"},
				{v: "HTTP/1.1", line: "req_proto=1.1"},
				{v: "HTTP/2", line: "req_proto=2"},
				{v: "HTTP/2.0", line: "req_proto=2.0"},
				{v: emptyStr, line: zeroLogLine},
				{v: hyphen, line: zeroLogLine},
				{v: "1.1", err: errBadReqProto},
				{v: "http/1.1", err: errBadReqProto},
			},
		},
		{
			"Response status",
			[]string{
				"status",
				"s",
				">s",
			},
			[]testCase{
				{v: "100", line: "resp_status=100"},
				{v: "200", line: "resp_status=200"},
				{v: "300", line: "resp_status=300"},
				{v: "400", line: "resp_status=400"},
				{v: "500", line: "resp_status=500"},
				{v: "600", line: "resp_status=600"},
				{v: "99", err: errBadRespStatus},
				{v: "601", err: errBadRespStatus},
				{v: emptyStr, line: zeroLogLine},
				{v: hyphen, line: zeroLogLine},
				{v: "200 ", err: errBadRespStatus},
				{v: "0.222", err: errBadRespStatus},
				{v: "localhost", err: errBadRespStatus},
			},
		},
		{
			"Request size",
			[]string{
				"request_length",
				"I",
			},
			[]testCase{
				{v: "15", line: "req_size=15"},
				{v: "1000000", line: "req_size=1000000"},
				{v: emptyStr, line: zeroLogLine},
				{v: hyphen, line: zeroLogLine},
				{v: "-1", err: errBadReqSize},
				{v: "100.222", err: errBadReqSize},
				{v: "invalid", err: errBadReqSize},
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
				{v: "15", line: "resp_size=15"},
				{v: "1000000", line: "resp_size=1000000"},
				{v: emptyStr, line: zeroLogLine},
				{v: hyphen, line: "resp_size=0"},
				{v: "-1", err: errBadRespSize},
				{v: "100.222", err: errBadRespSize},
				{v: "invalid", err: errBadRespSize},
			},
		},
		{
			"Response time",
			[]string{
				"request_time",
				"D",
			},
			[]testCase{
				{v: "100222", line: "resp_time=100222"},
				{v: "100.222", line: "resp_time=100222000"},
				{v: emptyStr, line: zeroLogLine},
				{v: hyphen, line: zeroLogLine},
				{v: "-1", err: errBadRespTime},
				{v: "0.333,0.444,0.555", err: errBadRespTime},
				{v: "number", err: errBadRespTime},
			},
		},
		{
			"Upstream response time",
			[]string{
				"upstream_response_time",
			},
			[]testCase{
				{v: "100222", line: "ups_resp_time=100222"},
				{v: "100.222", line: "ups_resp_time=100222000"},
				{v: "0.333,0.444,0.555", line: "ups_resp_time=333000"},
				{v: emptyStr, line: zeroLogLine},
				{v: hyphen, line: zeroLogLine},
				{v: "-1", err: errBadUpstreamRespTime},
				{v: "number", err: errBadUpstreamRespTime},
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

					if tc.err != nil {
						require.Error(t, err)
						assert.True(t, errors.Is(err, tc.err))
					} else {
						require.NoError(t, err)
					}

					expected := prepareLogLine(t, tc.line)
					assert.Equal(t, expected, *line)
				})
			}
		}
	}
}

//func TestLogLine_verify(t *testing.T) {
//	type testCase struct {
//		line string
//		err  error
//	}
//	tests := []struct {
//		name  string
//		cases []testCase
//	}{
//		{
//			"vhost",
//			[]testCase{
//				{line: "vhost=192.168.0.1"},
//				{line: "vhost=debian10.debian"},
//				{line: "vhost=1ce:1ce::babe"},
//				{line: "vhost=localhost"},
//				{line: emptyLine},
//				{line: "\"vhost=localhost \"", err: errBadVhost},
//				{line: "vhost=invalid_vhost", err: errBadVhost},
//				{line: "vhost=http://192.168.0.1/", err: errBadVhost},
//			},
//		},
//		{
//			"port",
//			[]testCase{
//				{line: "port=80"},
//				{line: "port=8081"},
//				{line: emptyLine},
//				{line: "\"port=80 \"", err: errBadPort},
//				{line: "port=79", err: errBadPort},
//				{line: "port=0.0.0.0", err: errBadPort},
//			},
//		},
//		{
//			"request scheme",
//			[]testCase{
//				{line: "req_scheme=http"},
//				{line: "req_scheme=https"},
//				{line: emptyLine},
//				{line: "req_scheme=not_https", err: errBadReqScheme},
//				{line: "req_scheme=HTTP", err: errBadReqScheme},
//				{line: "req_scheme=HTTPS", err: errBadReqScheme},
//				{line: "req_scheme=10", err: errBadReqScheme},
//			},
//		},
//		{
//			"request method",
//			[]testCase{
//				{line: "req_method=GET"},
//				{line: "req_method=POST"},
//				{line: "req_method=TRACE"},
//				{line: "req_method=OPTIONS"},
//				{line: "req_method=CONNECT"},
//				{line: "req_method=DELETE"},
//				{line: "req_method=PUT"},
//				{line: "req_method=PATCH"},
//				{line: "req_method=HEAD"},
//				{line: emptyLine},
//				{line: "req_method=Get", err: errBadReqMethod},
//				{line: "req_method=get", err: errBadReqMethod},
//				{line: "req_method=-", err: errBadReqMethod},
//			},
//		},
//		{
//			"request uri",
//			[]testCase{
//				{line: "req_uri=/"},
//				{line: "req_uri=/status?full&json"},
//				{line: "req_uri=/icons/openlogo-75.png"},
//				{line: emptyLine},
//				{line: "req_uri=status?full&json", err: errBadReqURI},
//				{line: "\"req_uri=/ \"", err: errBadReqURI},
//				{line: "req_uri=http://192.168.0.1/", err: errBadReqURI},
//			},
//		},
//		{
//			"request protocol",
//			[]testCase{
//				{line: "req_proto=1"},
//				{line: "req_proto=1.0"},
//				{line: "req_proto=1.1"},
//				{line: "req_proto=2.0"},
//				{line: "req_proto=2"},
//				{line: emptyLine},
//				{line: "req_proto=0.9", err: errBadReqProto},
//				{line: "req_proto=1.1.1", err: errBadReqProto},
//				{line: "req_proto=2.2", err: errBadReqProto},
//				{line: "req_proto=localhost", err: errBadReqProto},
//			},
//		},
//		{
//			"request size",
//			[]testCase{
//				{line: "req_size=0"},
//				{line: "req_size=100"},
//				{line: "req_size=1000000"},
//				{line: emptyLine},
//				{line: "req_size=-1", err: errBadReqSize},
//			},
//		},
//		{
//			"response size",
//			[]testCase{
//				{line: "resp_size=0"},
//				{line: "resp_size=100"},
//				{line: "resp_size=1000000"},
//				{line: emptyLine},
//				{line: "resp_size=-1", err: errBadRespSize},
//			},
//		},
//		{
//			"response time",
//			[]testCase{
//				{line: "resp_time=0"},
//				{line: "resp_time=0.000"},
//				{line: "resp_time=100"},
//				{line: "resp_time=1000000.123"},
//				{line: emptyLine},
//				{line: "resp_time=-1", err: errBadRespTime},
//			},
//		},
//		{
//			"upstream response time",
//			[]testCase{
//				{line: "ups_resp_time=0"},
//				{line: "ups_resp_time=0.000"},
//				{line: "ups_resp_time=100"},
//				{line: "ups_resp_time=1000000.123"},
//				{line: emptyLine},
//				{line: "ups_resp_time=-1", err: errBadUpstreamRespTime},
//			},
//		},
//		{
//			"response status",
//			[]testCase{
//				{line: "resp_status=100"},
//				{line: "resp_status=200"},
//				{line: "resp_status=300"},
//				{line: "resp_status=400"},
//				{line: "resp_status=500"},
//				{line: "resp_status=600"},
//				{line: emptyLine, err: errMandatoryField},
//				{line: "resp_status=-1", err: errBadRespStatus},
//				{line: "resp_status=99", err: errBadRespStatus},
//				{line: "resp_status=601", err: errBadRespStatus},
//			},
//		},
//	}
//
//	for _, tt := range tests {
//		for i, c := range tt.cases {
//			name := fmt.Sprintf("name=%s|line='%s'(%d)", tt.name, c.line, i+1)
//
//			t.Run(name, func(t *testing.T) {
//				line := prepareLogLine(t, c.line)
//				if tt.name != "response status" {
//					line.respStatus = 200
//				}
//				err := line.verify()
//
//				if c.err != nil {
//					require.Error(t, err)
//					assert.True(t, errors.Is(err, c.err))
//				} else {
//					require.NoError(t, err)
//				}
//			})
//		}
//	}
//}

func prepareLogLine(t *testing.T, from string) logLine {
	if from == zeroLogLine || from == emptyStr {
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

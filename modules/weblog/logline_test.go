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
)

func TestLogLine_Assign(t *testing.T) {
	type test struct {
		input    string
		wantLine string
		wantErr  error
	}
	type groupTest struct {
		name   string
		fields []string
		tests  []test
	}
	tests := []groupTest{
		{
			name: "Vhost",
			fields: []string{
				"host",
				"http_host",
				"v",
			},
			tests: []test{
				{input: "1.1.1.1", wantLine: "vhost=1.1.1.1"},
				{input: "::1", wantLine: "vhost=::1"},
				{input: "[::1]", wantLine: "vhost=::1"},
				{input: "1ce:1ce::babe", wantLine: "vhost=1ce:1ce::babe"},
				{input: "[1ce:1ce::babe]", wantLine: "vhost=1ce:1ce::babe"},
				{input: "localhost", wantLine: "vhost=localhost"},
				{input: "debian10.debian", wantLine: "vhost=debian10.debian"},
				{input: "my_vhost", wantLine: "vhost=my_vhost"},
				{input: emptyStr, wantLine: emptyLogLine},
				{input: hyphen, wantLine: emptyLogLine},
			},
		},
		{
			name: "Server Port",
			fields: []string{
				"server_port",
				"p",
			},
			tests: []test{
				{input: "80", wantLine: "port=80"},
				{input: "8081", wantLine: "port=8081"},
				{input: "30000", wantLine: "port=30000"},
				{input: emptyStr, wantLine: emptyLogLine},
				{input: hyphen, wantLine: emptyLogLine},
				{input: "-1", wantLine: emptyLogLine, wantErr: errBadPort},
				{input: "0", wantLine: emptyLogLine, wantErr: errBadPort},
				{input: "50000", wantLine: emptyLogLine, wantErr: errBadPort},
			},
		},
		{
			name: "Vhost With Port",
			fields: []string{
				"host:$server_port",
				"v:%p",
			},
			tests: []test{
				{input: "1.1.1.1:80", wantLine: "vhost=1.1.1.1 port=80"},
				{input: "::1:80", wantLine: "vhost=::1 port=80"},
				{input: "[::1]:80", wantLine: "vhost=::1 port=80"},
				{input: "1ce:1ce::babe:80", wantLine: "vhost=1ce:1ce::babe port=80"},
				{input: "debian10.debian:81", wantLine: "vhost=debian10.debian port=81"},
				{input: emptyStr, wantLine: emptyLogLine},
				{input: hyphen, wantLine: emptyLogLine},
				{input: "1.1.1.1", wantLine: emptyLogLine, wantErr: errBadVhostPort},
				{input: "1.1.1.1:", wantLine: emptyLogLine, wantErr: errBadVhostPort},
				{input: "1.1.1.1 80", wantLine: emptyLogLine, wantErr: errBadVhostPort},
				{input: "1.1.1.1:20", wantLine: emptyLogLine, wantErr: errBadVhostPort},
				{input: "1.1.1.1:50000", wantLine: emptyLogLine, wantErr: errBadVhostPort},
			},
		},
		{
			name: "Scheme",
			fields: []string{
				"scheme",
			},
			tests: []test{
				{input: "http", wantLine: "req_scheme=http"},
				{input: "https", wantLine: "req_scheme=https"},
				{input: emptyStr, wantLine: emptyLogLine},
				{input: hyphen, wantLine: emptyLogLine},
				{input: "HTTP", wantLine: emptyLogLine, wantErr: errBadReqScheme},
				{input: "HTTPS", wantLine: emptyLogLine, wantErr: errBadReqScheme},
			},
		},
		{
			name: "Client",
			fields: []string{
				"remote_addr",
				"a",
				"h",
			},
			tests: []test{
				{input: "1.1.1.1", wantLine: "req_client=1.1.1.1"},
				{input: "debian10", wantLine: "req_client=debian10"},
				{input: emptyStr, wantLine: emptyLogLine},
				{input: hyphen, wantLine: emptyLogLine},
			},
		},
		{
			name: "Request",
			fields: []string{
				"request",
				"r",
			},
			tests: []test{
				{input: "GET / HTTP/1.0", wantLine: "req_method=GET req_url=/ req_proto=1.0"},
				{input: "HEAD /ihs.gif HTTP/1.0", wantLine: "req_method=HEAD req_url=/ihs.gif req_proto=1.0"},
				{input: "POST /ihs.gif HTTP/1.0", wantLine: "req_method=POST req_url=/ihs.gif req_proto=1.0"},
				{input: "PUT /ihs.gif HTTP/1.0", wantLine: "req_method=PUT req_url=/ihs.gif req_proto=1.0"},
				{input: "PATCH /ihs.gif HTTP/1.0", wantLine: "req_method=PATCH req_url=/ihs.gif req_proto=1.0"},
				{input: "DELETE /ihs.gif HTTP/1.0", wantLine: "req_method=DELETE req_url=/ihs.gif req_proto=1.0"},
				{input: "OPTIONS /ihs.gif HTTP/1.0", wantLine: "req_method=OPTIONS req_url=/ihs.gif req_proto=1.0"},
				{input: "TRACE /ihs.gif HTTP/1.0", wantLine: "req_method=TRACE req_url=/ihs.gif req_proto=1.0"},
				{input: "CONNECT ip.cn:443 HTTP/1.1", wantLine: "req_method=CONNECT req_url=ip.cn:443 req_proto=1.1"},
				{input: "GET / HTTP/1.1", wantLine: "req_method=GET req_url=/ req_proto=1.1"},
				{input: "GET / HTTP/2", wantLine: "req_method=GET req_url=/ req_proto=2"},
				{input: "GET / HTTP/2.0", wantLine: "req_method=GET req_url=/ req_proto=2.0"},
				{input: "GET /invalid_version http/1.1", wantLine: "req_method=GET req_url=/invalid_version", wantErr: errBadReqProto},
				{input: emptyStr, wantLine: emptyLogLine},
				{input: hyphen, wantLine: emptyLogLine},
				{input: "GET no_version", wantLine: emptyLogLine, wantErr: errBadRequest},
				{input: "GOT / HTTP/2", wantLine: emptyLogLine, wantErr: errBadReqMethod},
				{input: "get / HTTP/2", wantLine: emptyLogLine, wantErr: errBadReqMethod},
				{input: "x04\x01\x00P$3\xFE\xEA\x00", wantLine: emptyLogLine, wantErr: errBadRequest},
			},
		},
		{
			name: "Request HTTP Method",
			fields: []string{
				"request_method",
				"m",
			},
			tests: []test{
				{input: "GET", wantLine: "req_method=GET"},
				{input: "HEAD", wantLine: "req_method=HEAD"},
				{input: "POST", wantLine: "req_method=POST"},
				{input: "PUT", wantLine: "req_method=PUT"},
				{input: "PATCH", wantLine: "req_method=PATCH"},
				{input: "DELETE", wantLine: "req_method=DELETE"},
				{input: "OPTIONS", wantLine: "req_method=OPTIONS"},
				{input: "TRACE", wantLine: "req_method=TRACE"},
				{input: "CONNECT", wantLine: "req_method=CONNECT"},
				{input: emptyStr, wantLine: emptyLogLine},
				{input: hyphen, wantLine: emptyLogLine},
				{input: "GET no_version", wantLine: emptyLogLine, wantErr: errBadReqMethod},
				{input: "GOT / HTTP/2", wantLine: emptyLogLine, wantErr: errBadReqMethod},
				{input: "get / HTTP/2", wantLine: emptyLogLine, wantErr: errBadReqMethod},
			},
		},
		{
			name: "Request URL",
			fields: []string{
				"request_uri",
				"U",
			},
			tests: []test{
				{input: "/server-status?auto", wantLine: "req_url=/server-status?auto"},
				{input: "/default.html", wantLine: "req_url=/default.html"},
				{input: "10.0.0.1:3128", wantLine: "req_url=10.0.0.1:3128"},
				{input: emptyStr, wantLine: emptyLogLine},
				{input: hyphen, wantLine: emptyLogLine},
			},
		},
		{
			name: "Request HTTP Protocol",
			fields: []string{
				"server_protocol",
				"H",
			},
			tests: []test{
				{input: "HTTP/1.0", wantLine: "req_proto=1.0"},
				{input: "HTTP/1.1", wantLine: "req_proto=1.1"},
				{input: "HTTP/2", wantLine: "req_proto=2"},
				{input: "HTTP/2.0", wantLine: "req_proto=2.0"},
				{input: emptyStr, wantLine: emptyLogLine},
				{input: hyphen, wantLine: emptyLogLine},
				{input: "1.1", wantLine: emptyLogLine, wantErr: errBadReqProto},
				{input: "http/1.1", wantLine: emptyLogLine, wantErr: errBadReqProto},
			},
		},
		{
			name: "Response Status Code",
			fields: []string{
				"status",
				"s",
				">s",
			},
			tests: []test{
				{input: "100", wantLine: "resp_code=100"},
				{input: "200", wantLine: "resp_code=200"},
				{input: "300", wantLine: "resp_code=300"},
				{input: "400", wantLine: "resp_code=400"},
				{input: "500", wantLine: "resp_code=500"},
				{input: "600", wantLine: "resp_code=600"},
				{input: emptyStr, wantLine: emptyLogLine},
				{input: hyphen, wantLine: emptyLogLine},
				{input: "99", wantLine: emptyLogLine, wantErr: errBadRespCode},
				{input: "601", wantLine: emptyLogLine, wantErr: errBadRespCode},
				{input: "200 ", wantLine: emptyLogLine, wantErr: errBadRespCode},
				{input: "0.222", wantLine: emptyLogLine, wantErr: errBadRespCode},
				{input: "localhost", wantLine: emptyLogLine, wantErr: errBadRespCode},
			},
		},
		{
			name: "Request Size",
			fields: []string{
				"request_length",
				"I",
			},
			tests: []test{
				{input: "15", wantLine: "req_size=15"},
				{input: "1000000", wantLine: "req_size=1000000"},
				{input: emptyStr, wantLine: emptyLogLine},
				{input: hyphen, wantLine: "req_size=0"},
				{input: "-1", wantLine: emptyLogLine, wantErr: errBadReqSize},
				{input: "100.222", wantLine: emptyLogLine, wantErr: errBadReqSize},
				{input: "invalid", wantLine: emptyLogLine, wantErr: errBadReqSize},
			},
		},
		{
			name: "Response Size",
			fields: []string{
				"bytes_sent",
				"body_bytes_sent",
				"O",
				"B",
				"b",
			},
			tests: []test{
				{input: "15", wantLine: "resp_size=15"},
				{input: "1000000", wantLine: "resp_size=1000000"},
				{input: emptyStr, wantLine: emptyLogLine},
				{input: hyphen, wantLine: "resp_size=0"},
				{input: "-1", wantLine: emptyLogLine, wantErr: errBadRespSize},
				{input: "100.222", wantLine: emptyLogLine, wantErr: errBadRespSize},
				{input: "invalid", wantLine: emptyLogLine, wantErr: errBadRespSize},
			},
		},
		{
			name: "Request Processing Time",
			fields: []string{
				"request_time",
				"D",
			},
			tests: []test{
				{input: "100222", wantLine: "req_proc_time=100222"},
				{input: "100.222", wantLine: "req_proc_time=100222000"},
				{input: emptyStr, wantLine: emptyLogLine},
				{input: hyphen, wantLine: emptyLogLine},
				{input: "-1", wantLine: emptyLogLine, wantErr: errBadReqProcTime},
				{input: "0.333,0.444,0.555", wantLine: emptyLogLine, wantErr: errBadReqProcTime},
				{input: "number", wantLine: emptyLogLine, wantErr: errBadReqProcTime},
			},
		},
		{
			name: "Upstream Response Time",
			fields: []string{
				"upstream_response_time",
			},
			tests: []test{
				{input: "100222", wantLine: "ups_resp_time=100222"},
				{input: "100.222", wantLine: "ups_resp_time=100222000"},
				{input: "0.333,0.444,0.555", wantLine: "ups_resp_time=333000"},
				{input: emptyStr, wantLine: emptyLogLine},
				{input: hyphen, wantLine: emptyLogLine},
				{input: "-1", wantLine: emptyLogLine, wantErr: errBadUpsRespTime},
				{input: "number", wantLine: emptyLogLine, wantErr: errBadUpsRespTime},
			},
		},
		{
			name: "SSL Protocol",
			fields: []string{
				"ssl_protocol",
			},
			tests: []test{
				{input: "SSLv2", wantLine: "ssl_proto=SSLv2"},
				{input: "TLSv1", wantLine: "ssl_proto=TLSv1"},
				{input: "TLSv1.1", wantLine: "ssl_proto=TLSv1.1"},
				{input: "TLSv1.2", wantLine: "ssl_proto=TLSv1.2"},
				{input: "TLSv1.3", wantLine: "ssl_proto=TLSv1.3"},
				{input: emptyStr, wantLine: emptyLogLine},
				{input: hyphen, wantLine: emptyLogLine},
				{input: "-1", wantLine: emptyLogLine, wantErr: errBadSSLProto},
				{input: "invalid", wantLine: emptyLogLine, wantErr: errBadSSLProto},
			},
		},
		{
			name: "SSL Cipher Suite",
			fields: []string{
				"ssl_cipher",
			},
			tests: []test{
				{input: "ECDHE-RSA-AES256-SHA", wantLine: "ssl_cipher_suite=ECDHE-RSA-AES256-SHA"},
				{input: "DHE-RSA-AES256-SHA", wantLine: "ssl_cipher_suite=DHE-RSA-AES256-SHA"},
				{input: "AES256-SHA", wantLine: "ssl_cipher_suite=AES256-SHA"},
				{input: "PSK-RC4-SHA", wantLine: "ssl_cipher_suite=PSK-RC4-SHA"},
				{input: emptyStr, wantLine: emptyLogLine},
				{input: hyphen, wantLine: emptyLogLine},
				{input: "-1", wantLine: emptyLogLine, wantErr: errBadSSLCipherSuite},
				{input: "invalid", wantLine: emptyLogLine, wantErr: errBadSSLCipherSuite},
			},
		},
		{
			name: "Custom Fields",
			fields: []string{
				"custom_field",
			},
			tests: []test{
				{input: "POST", wantLine: "custom_field=POST"},
				{input: "/example.com", wantLine: "custom_field=/example.com"},
				{input: "HTTP/1.1", wantLine: "custom_field=HTTP/1.1"},
				{input: "0.333,0.444,0.555", wantLine: "custom_field=0.333,0.444,0.555"},
				{input: "-1", wantLine: "custom_field=-1"},
				{input: "invalid", wantLine: "custom_field=invalid"},
				{input: emptyStr, wantLine: emptyLogLine},
				{input: hyphen, wantLine: emptyLogLine},
			},
		},
		{
			name: "Custom Fields Not Exist",
			fields: []string{
				"custom_field_not_exist",
			},
			tests: []test{
				{input: "POST", wantLine: emptyLogLine},
				{input: "/example.com", wantLine: emptyLogLine},
				{input: "HTTP/1.1", wantLine: emptyLogLine},
				{input: "0.333,0.444,0.555", wantLine: emptyLogLine},
				{input: "-1", wantLine: emptyLogLine},
				{input: "invalid", wantLine: emptyLogLine},
				{input: emptyStr, wantLine: emptyLogLine},
				{input: hyphen, wantLine: emptyLogLine},
			},
		},
	}

	for _, tt := range tests {
		for _, varName := range tt.fields {
			for i, tc := range tt.tests {
				name := fmt.Sprintf("[%s:%d]field='%s'|input='%s'", tt.name, i+1, varName, tc.input)
				t.Run(name, func(t *testing.T) {

					line := newEmptyLogLineWithFields()
					err := line.Assign(varName, tc.input)

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
	type test struct {
		line    string
		wantErr error
	}
	tests := []struct {
		name  string
		tests []test
	}{
		{
			name: "Vhost",
			tests: []test{
				{line: "vhost=192.168.0.1"},
				{line: "vhost=debian10.debian"},
				{line: "vhost=1ce:1ce::babe"},
				{line: "vhost=localhost"},
				{line: "vhost=invalid_vhost", wantErr: errBadVhost},
				{line: "vhost=http://192.168.0.1/", wantErr: errBadVhost},
			},
		},
		{
			name: "Server Port",
			tests: []test{
				{line: "port=80"},
				{line: "port=8081"},
				{line: "port=79", wantErr: errBadPort},
				{line: "port=50000", wantErr: errBadPort},
				{line: "port=0.0.0.0", wantErr: errBadPort},
			},
		},
		{
			name: "Vhost With Port",
			tests: []test{
				{line: "vhost=1.1.1.1 port=80"},
				{line: "vhost=::1 port=8081"},
				{line: "vhost=1ce:1ce::babe port=81"},
				{line: "vhost=debian10.debian port=81"},
				{line: "vhost=::1 port=79", wantErr: errBadPort},
				{line: "vhost=::1 port=50000", wantErr: errBadPort},
				{line: "vhost=::1 port=0.0.0.0", wantErr: errBadPort},
			},
		},
		{
			name: "Scheme",
			tests: []test{
				{line: "req_scheme=http"},
				{line: "req_scheme=https"},
				{line: "req_scheme=not_https", wantErr: errBadReqScheme},
				{line: "req_scheme=HTTP", wantErr: errBadReqScheme},
				{line: "req_scheme=HTTPS", wantErr: errBadReqScheme},
				{line: "req_scheme=10", wantErr: errBadReqScheme},
			},
		},
		{
			name: "Client",
			tests: []test{
				{line: "req_client=1.1.1.1"},
				{line: "req_client=::1"},
				{line: "req_client=1ce:1ce::babe"},
				{line: "req_client=localhost"},
				{line: "req_client=debian10.debian", wantErr: errBadReqClient},
				{line: "req_client=invalid", wantErr: errBadReqClient},
			},
		},
		{
			name: "Request HTTP Method",
			tests: []test{
				{line: "req_method=GET"},
				{line: "req_method=POST"},
				{line: "req_method=TRACE"},
				{line: "req_method=OPTIONS"},
				{line: "req_method=CONNECT"},
				{line: "req_method=DELETE"},
				{line: "req_method=PUT"},
				{line: "req_method=PATCH"},
				{line: "req_method=HEAD"},
				{line: "req_method=Get", wantErr: errBadReqMethod},
				{line: "req_method=get", wantErr: errBadReqMethod},
				{line: "req_method=-", wantErr: errBadReqMethod},
			},
		},
		{
			name: "Request URL",
			tests: []test{
				{line: "req_url=/"},
				{line: "req_url=/status?full&json"},
				{line: "req_url=/icons/openlogo-75.png"},
				{line: "req_method=CONNECT req_url=http://192.168.0.1/"},
				{line: "req_url=status?full&json", wantErr: errBadReqURL},
				{line: "\"req_url=/ \"", wantErr: errBadReqURL},
				{line: "req_url=http://192.168.0.1/", wantErr: errBadReqURL},
			},
		},
		{
			name: "Request HTTP Protocol",
			tests: []test{
				{line: "req_proto=1"},
				{line: "req_proto=1.0"},
				{line: "req_proto=1.1"},
				{line: "req_proto=2.0"},
				{line: "req_proto=2"},
				{line: "req_proto=0.9", wantErr: errBadReqProto},
				{line: "req_proto=1.1.1", wantErr: errBadReqProto},
				{line: "req_proto=2.2", wantErr: errBadReqProto},
				{line: "req_proto=localhost", wantErr: errBadReqProto},
			},
		},
		{
			name: "Response Status Code",
			tests: []test{
				{line: "resp_code=100"},
				{line: "resp_code=200"},
				{line: "resp_code=300"},
				{line: "resp_code=400"},
				{line: "resp_code=500"},
				{line: "resp_code=600"},
				{line: "resp_code=-1", wantErr: errBadRespCode},
				{line: "resp_code=99", wantErr: errBadRespCode},
				{line: "resp_code=601", wantErr: errBadRespCode},
			},
		},
		{
			name: "Request size",
			tests: []test{
				{line: "req_size=0"},
				{line: "req_size=100"},
				{line: "req_size=1000000"},
				{line: "req_size=-1", wantErr: errBadReqSize},
			},
		},
		{
			name: "Response size",
			tests: []test{
				{line: "resp_size=0"},
				{line: "resp_size=100"},
				{line: "resp_size=1000000"},
				{line: "resp_size=-1", wantErr: errBadRespSize},
			},
		},
		{
			name: "Request Processing Time",
			tests: []test{
				{line: "req_proc_time=0"},
				{line: "req_proc_time=0.000"},
				{line: "req_proc_time=100"},
				{line: "req_proc_time=1000000.123"},
				{line: "req_proc_time=-1", wantErr: errBadReqProcTime},
			},
		},
		{
			name: "Upstream Response Time",
			tests: []test{
				{line: "ups_resp_time=0"},
				{line: "ups_resp_time=0.000"},
				{line: "ups_resp_time=100"},
				{line: "ups_resp_time=1000000.123"},
				{line: "ups_resp_time=-1", wantErr: errBadUpsRespTime},
			},
		},
		{
			name: "Upstream Response Time",
			tests: []test{
				{line: "ups_resp_time=0"},
				{line: "ups_resp_time=0.000"},
				{line: "ups_resp_time=100"},
				{line: "ups_resp_time=1000000.123"},
				{line: "ups_resp_time=-1", wantErr: errBadUpsRespTime},
			},
		},
		{
			name: "SSL Protocol",
			tests: []test{
				{line: "ssl_proto=SSLv2"},
				{line: "ssl_proto=TLSv1"},
				{line: "ssl_proto=TLSv1.1"},
				{line: "ssl_proto=TLSv1.2"},
				{line: "ssl_proto=TLSv1.3"},
				{line: "ssl_proto=invalid", wantErr: errBadSSLProto},
			},
		},
		{
			name: "SSL Cipher Suite",
			tests: []test{
				{line: "ssl_cipher_suite=ECDHE-RSA-AES256-SHA"},
				{line: "ssl_cipher_suite=DHE-RSA-AES256-SHA"},
				{line: "ssl_cipher_suite=AES256-SHA"},
				{line: "ssl_cipher_suite=PSK-RC4-SHA"},
				{line: "ssl_cipher_suite=invalid", wantErr: errBadSSLCipherSuite},
			},
		},
		{
			name: "Custom Fields",
			tests: []test{
				{line: "custom_field=POST"},
				{line: "custom_field=/example.com"},
				{line: "custom_field=HTTP/1.1"},
				{line: "custom_field=0.333,0.444,0.555"},
				{line: "custom_field=-1"},
				{line: "custom_field=invalid"},
			},
		},
		{
			name: "Empty Line",
			tests: []test{
				{line: "custom_field_not_exit=POST", wantErr: errEmptyLine},
				{line: emptyLogLine, wantErr: errEmptyLine},
			},
		},
	}

	for _, tt := range tests {
		for i, c := range tt.tests {
			name := fmt.Sprintf("name=%s|line='%s'(%d)", tt.name, c.line, i+1)

			t.Run(name, func(t *testing.T) {
				line := prepareLogLine(t, c.line)

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

func newEmptyLogLineWithFields() *logLine {
	l := newEmptyLogLine()
	l.custom.fields = map[string]struct{}{"custom_field": {}}
	return l
}

func prepareLogLine(t *testing.T, from string) logLine {
	require.NotZero(t, from)

	if from == emptyLogLine {
		return *newEmptyLogLineWithFields()
	}

	t.Helper()
	r := csv.NewReader(strings.NewReader(from))
	r.Comma = ' '
	line := newEmptyLogLineWithFields()

	rs, err := r.Read()
	require.NoError(t, err)

	for _, v := range rs {
		parts := strings.Split(v, "=")
		require.Len(t, parts, 2)
		field, val := parts[0], parts[1]
		switch field {
		case "vhost":
			line.vhost = val
		case "port":
			line.port = val
		case "req_scheme":
			line.reqScheme = val
		case "req_client":
			line.reqClient = val
		case "req_method":
			line.reqMethod = val
		case "req_url":
			line.reqURL = val
		case "req_proto":
			line.reqProto = val
		case "req_size":
			i, err := strconv.Atoi(val)
			require.NoError(t, err)
			line.reqSize = i
		case "resp_code":
			i, err := strconv.Atoi(val)
			require.NoError(t, err)
			line.respCode = i
		case "resp_size":
			i, err := strconv.Atoi(val)
			require.NoError(t, err)
			line.respSize = i
		case "req_proc_time":
			i, err := strconv.ParseFloat(val, 64)
			require.NoError(t, err)
			line.reqProcTime = i
		case "ups_resp_time":
			i, err := strconv.ParseFloat(val, 64)
			require.NoError(t, err)
			line.upsRespTime = i
		case "ssl_proto":
			line.sslProto = val
		case "ssl_cipher_suite":
			line.sslCipherSuite = val
		case "custom_field":
			line.custom.values = append(line.custom.values, customValue{name: field, value: val})
		}
	}
	return *line
}

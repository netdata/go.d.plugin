package squidlog

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogLine_Assign(t *testing.T) {
	type subTest struct {
		input    string
		wantLine logLine
		wantErr  error
	}
	type test struct {
		name  string
		field string
		cases []subTest
	}
	tests := []test{
		{
			name:  "Response Time",
			field: fieldRespTime,
			cases: []subTest{
				{input: "0", wantLine: logLine{respTime: 0}},
				{input: "1000", wantLine: logLine{respTime: 1000}},
				{input: "", wantLine: emptyLogLine},
				{input: "-1", wantLine: emptyLogLine, wantErr: errBadRespTime},
				{input: "0.000", wantLine: emptyLogLine, wantErr: errBadRespTime},
				{input: hyphen, wantLine: emptyLogLine, wantErr: errBadRespTime},
			},
		},
		{
			name:  "Client Address",
			field: fieldClientAddr,
			cases: []subTest{
				{input: "127.0.0.1", wantLine: logLine{clientAddr: "127.0.0.1"}},
				{input: "::1", wantLine: logLine{clientAddr: "::1"}},
				{input: "kadr20.m1.netdata.lan", wantLine: logLine{clientAddr: "kadr20.m1.netdata.lan"}},
				{input: "±!@#$%^&*()", wantLine: logLine{clientAddr: "±!@#$%^&*()"}},
				{input: "", wantLine: emptyLogLine},
				{input: hyphen, wantLine: emptyLogLine, wantErr: errBadClientAddr},
			},
		},
		{
			name:  "Cache Code",
			field: fieldCacheCode,
			cases: []subTest{
				{input: "TCP_MISS", wantLine: logLine{cacheCode: "TCP_MISS"}},
				{input: "TCP_DENIED", wantLine: logLine{cacheCode: "TCP_DENIED"}},
				{input: "TCP_CLIENT_REFRESH_MISS", wantLine: logLine{cacheCode: "TCP_CLIENT_REFRESH_MISS"}},
				{input: "UDP_MISS_NOFETCH", wantLine: logLine{cacheCode: "UDP_MISS_NOFETCH"}},
				{input: "UDP_INVALID", wantLine: logLine{cacheCode: "UDP_INVALID"}},
				{input: "NONE", wantLine: logLine{cacheCode: "NONE"}},
				{input: "", wantLine: emptyLogLine},
				{input: hyphen, wantLine: emptyLogLine, wantErr: errBadCacheCode},
				{input: "NONE_MISS", wantLine: emptyLogLine, wantErr: errBadCacheCode},
			},
		},
	}

	for _, tt := range tests {
		for i, tc := range tt.cases {
			name := fmt.Sprintf("[%s:%d]field='%s'|input='%s'", tt.name, i+1, tt.field, tc.input)
			t.Run(name, func(t *testing.T) {

				line := newEmptyLogLine()
				err := line.Assign(tt.field, tc.input)

				if tc.wantErr != nil {
					require.Error(t, err)
					assert.Truef(t, errors.Is(err, tc.wantErr), "expected '%v' error, got '%v'", tc.wantErr, err)
				} else {
					require.NoError(t, err)
				}

				expected := prepareLogLine(tt.field, tc.wantLine)
				assert.Equal(t, expected, *line)
			})
		}
	}
}

func TestLogLine_verify(t *testing.T) {

}

func prepareLogLine(field string, template logLine) logLine {
	if template.empty() {
		return template
	}

	var line logLine
	line.reset()

	switch field {
	case fieldRespTime:
		line.respTime = template.respTime
	case fieldClientAddr:
		line.clientAddr = template.clientAddr
	case fieldCacheCode:
		line.cacheCode = template.cacheCode
	case fieldHTTPCode:
		line.httpCode = template.httpCode
	case fieldRespSize:
		line.respSize = template.respSize
	case fieldReqMethod:
		line.reqMethod = template.reqMethod
	case fieldHierCode:
		line.hierCode = template.hierCode
	case fieldMimeType:
		line.mimeType = template.mimeType
	case fieldServerAddr:
		line.serverAddr = template.serverAddr
	case fieldResultCode:
		line.cacheCode = template.cacheCode
		line.httpCode = template.httpCode
	case fieldHierarchy:
		line.hierCode = template.hierCode
		line.serverAddr = template.serverAddr
	}
	return line
}

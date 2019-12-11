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
				{input: hyphen, wantLine: emptyLogLine, wantErr: errBadRespTime},
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
					assert.True(t, errors.Is(err, tc.wantErr))
				} else {
					require.NoError(t, err)
				}

				expected := prepareLogLine(t, tt.field, tc.wantLine)
				assert.Equal(t, expected, *line)
			})
		}
	}
}

func TestLogLine_verify(t *testing.T) {

}

func prepareLogLine(t *testing.T, field string, template logLine) logLine {
	if template == emptyLogLine {
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

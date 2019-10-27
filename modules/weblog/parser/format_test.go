package parser

//
//import (
//	"testing"
//	"time"
//
//	"github.com/stretchr/testify/require"
//
//	"github.com/stretchr/testify/assert"
//)
//
//func TestParse(t *testing.T) {
//	type args struct {
//		timeUnit  time.Duration
//		logFormat string
//	}
//	tests := []struct {
//		name string
//		args args
//		want Format
//	}{
//		{"common",
//			args{time.Microsecond, `$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent`},
//			Format{
//				TimeScale:        1,
//				maxIndex:         8,
//				client:       0,
//				Request:          5,
//				respCode:           6,
//				BytesSent:        7,
//				Host:             -1,
//				ReqTime:          -1,
//				upstreamRespTime: -1,
//				ReqLength:        -1,
//				custom:           -1,
//			}},
//		{"combined",
//			args{time.Microsecond, `$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"`},
//			Format{
//				TimeScale:        1,
//				maxIndex:         10,
//				client:       0,
//				Request:          5,
//				respCode:           6,
//				BytesSent:        7,
//				Host:             -1,
//				ReqTime:          -1,
//				upstreamRespTime: -1,
//				ReqLength:        -1,
//				custom:           -1,
//			}},
//		{"custom",
//			args{time.Microsecond, `$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" $uid_got $request_time $http_host`},
//			Format{
//				TimeScale:        1,
//				maxIndex:         13,
//				client:       0,
//				Request:          5,
//				respCode:           6,
//				BytesSent:        7,
//				Host:             12,
//				ReqTime:          11,
//				upstreamRespTime: -1,
//				ReqLength:        -1,
//				custom:           -1,
//			}},
//		{"custom_xff",
//			args{time.Microsecond, `$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" $uid_got $request_time $http_host "$http_x_forwarded_for"`},
//			Format{
//				TimeScale:        1,
//				maxIndex:         14,
//				client:       0,
//				Request:          5,
//				respCode:           6,
//				BytesSent:        7,
//				Host:             12,
//				ReqTime:          11,
//				upstreamRespTime: -1,
//				ReqLength:        -1,
//				custom:           -1,
//			}},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			tt.want.Raw = tt.args.logFormat
//			assert.Equal(t, tt.want, *NewFormat(1, tt.args.logFormat))
//		})
//	}
//}
//
//func TestFormat_Parse(t *testing.T) {
//	var (
//		common    = NewFormat(time.Microsecond.Seconds(), `$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent`)
//		combined  = NewFormat(time.Microsecond.Seconds(), `$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"`)
//		custom    = NewFormat(time.Microsecond.Seconds(), `$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" $uid_got $request_time $http_host`)
//		customXFF = NewFormat(time.Microsecond.Seconds(), `$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" $uid_got $request_time $http_host "$http_x_forwarded_for"`)
//	)
//
//	type args struct {
//		record []string
//	}
//	tests := []struct {
//		name                                        string
//		line                                        []string
//		want                                        LogLine
//		commonOK, combinedOK, customOK, customXFFOK bool
//	}{
//		{"simple common", []string{`10.131.201.180`, `-`, `-`, `[07/Mar/2002:15:46:25`, `+0900]`, `GET / HTTP/1.1`, `200`, `1620`},
//			LogLine{
//				client: `10.131.201.180`,
//				Request:    `GET / HTTP/1.1`,
//				reqHTTPMethod:     `GET`,
//				reqURI:        "/",
//				reqHTTPVersion:    "1.1",
//				respCode:     200,
//				BytesSent:  1620,
//				ReqLength:  -1,
//				ReqTime:    -1,
//			},
//			true, false, false, false,
//		},
//		{"simple combined", []string{`10.131.201.180`, `-`, `-`, `[07/Mar/2002:15:46:25`, `+0900]`, `GET / HTTP/1.1`, `200`, `1620`, `-`, `-`},
//			LogLine{
//				client: `10.131.201.180`,
//				Request:    `GET / HTTP/1.1`,
//				reqHTTPMethod:     `GET`,
//				reqURI:        "/",
//				reqHTTPVersion:    "1.1",
//				respCode:     200,
//				BytesSent:  1620,
//				ReqLength:  -1,
//				ReqTime:    -1,
//			},
//			true, true, false, false,
//		},
//		{"simple custom", []string{`10.131.201.180`, `-`, `-`, `[07/Mar/2002:15:46:25`, `+0900]`, `GET / HTTP/1.1`, `200`, `1620`, `-`, `-`, `-`, `128`, `www.example.com`},
//			LogLine{
//				client: `10.131.201.180`,
//				Request:    `GET / HTTP/1.1`,
//				reqHTTPMethod:     `GET`,
//				reqURI:        "/",
//				reqHTTPVersion:    "1.1",
//				respCode:     200,
//				BytesSent:  1620,
//				ReqLength:  -1,
//				ReqTime:    0.000128,
//				Host:       "www.example.com",
//			},
//			true, true, true, false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name+"_common", func(t *testing.T) {
//			want := tt.want
//			want.ReqTime = -1
//			want.Host = ""
//			got, err := common.Parse(tt.line)
//			if tt.commonOK {
//				assert.Equal(t, want, got)
//			} else {
//				assert.Error(t, err)
//			}
//		})
//		t.Run(tt.name+"_combined", func(t *testing.T) {
//			want := tt.want
//			want.ReqTime = -1
//			want.Host = ""
//			got, err := combined.Parse(tt.line)
//			if tt.combinedOK {
//				assert.Equal(t, want, got)
//			} else {
//				assert.Error(t, err)
//			}
//		})
//		t.Run(tt.name+"_custom", func(t *testing.T) {
//			got, err := custom.Parse(tt.line)
//			if tt.customOK {
//				assert.Equal(t, tt.want, got)
//			} else {
//				assert.Error(t, err)
//			}
//		})
//		t.Run(tt.name+"_customXFF", func(t *testing.T) {
//			got, err := customXFF.Parse(tt.line)
//			if tt.customXFFOK {
//				assert.Equal(t, tt.want, got)
//			} else {
//				assert.Error(t, err)
//			}
//		})
//	}
//}
//
//func Test_parseRequest(t *testing.T) {
//	tests := []struct {
//		name        string
//		wantMethod  string
//		wantUri     string
//		wantVersion string
//		wantErr     bool
//	}{
//		{"GET / HTTP/1.1", "GET", "/", "1.1", false},
//		{"GET / HTTP/1.0", "GET", "/", "1.0", false},
//		{"GET / HTTP/2", "GET", "/", "2", false},
//		{"GET /ihs.gif HTTP/1.1", "GET", "/ihs.gif", "1.1", false},
//		{"GET no_version", "", "", "", true},
//		{"GET /invalid_version http/1.1", "GET", "/invalid_version", "", true},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			gotMethod, gotUri, gotVersion, err := parseRequest(tt.name)
//			if tt.wantErr {
//				require.Error(t, err)
//			} else {
//				require.NoError(t, err)
//			}
//			assert.Equal(t, tt.wantMethod, gotMethod)
//			assert.Equal(t, tt.wantUri, gotUri)
//			assert.Equal(t, tt.wantVersion, gotVersion)
//		})
//	}
//}

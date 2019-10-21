package weblog

import (
	"errors"
	"regexp"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/logs"
)

// log_format vcombined '$host:$server_port '
//        '$remote_addr - $remote_user [$time_local] '
//        '"$request" $status $body_bytes_sent '
//        '"$http_referer" "$http_user_agent"';

// LogFormat "%v:%p %h %l %u %t \"%r\" %>s %O \"%{Referer}i\" \"%{User-Agent}i\"" vhost_combined
// LogFormat "%h %l %u %t \"%r\" %>s %O \"%{Referer}i\" \"%{User-Agent}i\"" combined
// LogFormat "%h %l %u %t \"%r\" %>s %O" common

/*
| name               | nginx                   | apache    |
|--------------------|-------------------------|-----------|
| vhost              | $http                   | %v        | name of the server which accepted a request
| port               | $server_port            | %p        | port of the server which accepted a request
| client             | $remote_addr            | %a (%h)   | apache %h: logs the IP address if HostnameLookups is Off
| request            | $request                | %r        | req_method + req_uri + req_protocol
| req_method         | $request_method         | %m        |
| req_uri            | $request_uri            | %U        | nginx: w/ queries, apache: w/o
| req_proto          | $server_protocol        | %H        | request protocol, usually “HTTP/1.0”, “HTTP/1.1”, or “HTTP/2.0”
| resp_status        | $status                 | %s (%>s)  | response status
| req_size           | $request_length         | $I        | request length (including request line, header, and request body), apache: need mod_logio
| resp_size          | $bytes_sent             | %O        | number of bytes sent to a client, including headers
| resp_size          | $body_bytes_sent        | %B        | number of bytes sent to a client, not including headers
| req_time           | $request_time           | %D        | the time taken to serve the request. Apache: in microseconds, nginx: in seconds with a milliseconds resolution
| ups_resp_time      | $upstream_response_time | -         | keeps time spent on receiving the response from the upstream server; the time is kept in seconds with millisecond resolution. Times of several responses are separated by commas and colons
| custom             | -                       | -         |
*/

var (
	reLTSV = regexp.MustCompile(`^[a-zA-Z0-9]+:[^\t]*(\t[a-zA-Z0-9]+:[^\t]*)*$`)

	csvVhostCombined = `$host:$server_port $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"`

	csvCommon       = `      $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent`
	csvCombined     = `      $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"`
	csvCustom1      = `      $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" $uid_got                  $request_time`
	csvCustom2      = `      $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent $request_length $request_time'`
	csvCustom3      = `      $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent $request_length $request_time      "$upstream_response_time"`
	csvVhostCommon  = `$host $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent`
	csvVhostCustom1 = `$host $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" $uid_got                  $request_time`
	csvVhostCustom2 = `$host $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent $request_length $request_time`
	csvVhostCustom3 = `$host $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent $request_length $request_time      "$upstream_response_time"`

	guessOrder = []string{
		//csvVhostCustom1,
		//csvVhostCustom3,
		//csvVhostCustom2,
		//csvVhostCombined,
		//csvVhostCommon,
		//csvCustom1,
		//csvCustom3,
		//csvCustom2,
		//csvCombined,
		csvCommon,
	}
)

const (
	typeAuto = "auto"
)

func (w *WebLog) newParser(record []byte) (logs.Parser, error) {
	w.Parser.CSV.CheckField = checkCSVFormatField

	if w.Parser.LogType == typeAuto {
		return w.guessParser(record)
	}
	return logs.NewParser(w.Parser, w.file)
}

func (w *WebLog) guessParser(record []byte) (logs.Parser, error) {
	if reLTSV.Match(record) {
		return logs.NewLTSVParser(w.Parser.LTSV, w.file)
	}

	for _, format := range guessOrder {
		cfg := w.Parser.CSV
		cfg.Format = format
		cfg.TrimLeadingSpace = true

		parser, err := logs.NewCSVParser(cfg, w.file)
		if err != nil {
			return nil, err
		}

		line := newEmptyLogLine()
		if err := parser.Parse(record, line); err != nil {
			continue
		}

		if err = line.Verify(); err != nil {
			continue
		}
		return parser, nil
	}
	return nil, errors.New("cannot determine log format")
}

func checkCSVFormatField(name string) (newName string, valid bool, offset int) {
	if name == "[$time_local]" {
		offset = 1
		return
	}

	if !isValidVar(name) {
		return
	}

	newName = name[1:]
	valid = true
	return
}

func isValidVar(v string) bool { return len(v) > 1 && (isNginxVar(v) || isApacheVar(v)) }

func isNginxVar(v string) bool { return strings.HasPrefix(v, "$") }

func isApacheVar(v string) bool { return strings.HasPrefix(v, "%") }

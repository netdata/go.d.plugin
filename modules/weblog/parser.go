package weblog

import (
	"errors"
	"regexp"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/logs"
)

/*
Default apache log format:
 - "%v:%p %h %l %u %t \"%r\" %>s %O \"%{Referer}i\" \"%{User-Agent}i\"" vhost_combined
 - "%h %l %u %t \"%r\" %>s %O \"%{Referer}i\" \"%{User-Agent}i\"" combined
 - "%h %l %u %t \"%r\" %>s %O" common

Default nginx log format:
 - '$remote_addr - $remote_user [$time_local] '
   '"$request" $status $body_bytes_sent '
   '"$http_referer" "$http_user_agent"' combined

Netdata recommends:
 Nginx:
  - '$remote_addr - $remote_user [$time_local] '
    '"$request" $status $body_bytes_sent '
    '$request_length $request_time $upstream_response_time '
    '"$http_referer" "$http_user_agent"'

 Apache:
  - "%h %l %u %t \"%r\" %>s %B %I %D \"%{Referer}i\" \"%{User-Agent}i\""
*/

var (
	csvCommon       = `                   $remote_addr - - [$time_local] "$request" $status $body_bytes_sent`
	csvCustom1      = `                   $remote_addr - - [$time_local] "$request" $status $body_bytes_sent     $request_length $request_time`
	csvCustom2      = `                   $remote_addr - - [$time_local] "$request" $status $body_bytes_sent     $request_length $request_time $upstream_response_time`
	csvCustom3      = `                   $remote_addr - - [$time_local] "$request" $status $body_bytes_sent - - $request_length $request_time`
	csvCustom4      = `                   $remote_addr - - [$time_local] "$request" $status $body_bytes_sent - - $request_length $request_time $upstream_response_time`
	csvVhostCommon  = `$host:$server_port $remote_addr - - [$time_local] "$request" $status $body_bytes_sent`
	csvVhostCustom1 = `$host:$server_port $remote_addr - - [$time_local] "$request" $status $body_bytes_sent     $request_length $request_time`
	csvVhostCustom2 = `$host:$server_port $remote_addr - - [$time_local] "$request" $status $body_bytes_sent     $request_length $request_time $upstream_response_time`
	csvVhostCustom3 = `$host:$server_port $remote_addr - - [$time_local] "$request" $status $body_bytes_sent - - $request_length $request_time`
	csvVhostCustom4 = `$host:$server_port $remote_addr - - [$time_local] "$request" $status $body_bytes_sent - - $request_length $request_time $upstream_response_time`

	guessOrder = []string{
		csvVhostCustom4,
		csvVhostCustom3,
		csvVhostCustom2,
		csvVhostCustom1,
		csvVhostCommon,
		csvCustom4,
		csvCustom3,
		csvCustom2,
		csvCustom1,
		csvCommon,
	}
)

const (
	typeAuto = "auto"
)

var (
	reLTSV = regexp.MustCompile(`^[a-zA-Z0-9]+:[^\t]*(\t[a-zA-Z0-9]+:[^\t]*)*$`)
)

func (w *WebLog) newParser(record []byte) (logs.Parser, error) {
	if w.Parser.LogType == typeAuto {
		return w.guessParser(record)
	}
	return logs.NewParser(w.Parser, w.file)
}

func (w *WebLog) guessParser(record []byte) (logs.Parser, error) {
	w.Debug("starting log format auto detection")
	if reLTSV.Match(record) {
		return logs.NewLTSVParser(w.Parser.LTSV, w.file)
	}
	return w.guessCSVParser(record)
}

func (w *WebLog) guessCSVParser(record []byte) (logs.Parser, error) {
	for _, format := range guessOrder {
		format = cleanCSVFormat(format)
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

		if err = line.verify(); err != nil {
			continue
		}
		return parser, nil
	}
	return nil, errors.New("cannot determine log format")
}

func checkCSVFormatField(name string) (newName string, offset int, valid bool) {
	if name == "[$time_local]" || name == "$time_local" {
		return "", 1, false
	}
	if !isValidVar(name) {
		return "", 0, false
	}
	return name[1:], 0, true
}

func isValidVar(v string) bool {
	return len(v) > 1 && (isNginxVar(v) || isApacheVar(v))
}

func isNginxVar(v string) bool {
	return strings.HasPrefix(v, "$")
}

func isApacheVar(v string) bool {
	return strings.HasPrefix(v, "%")
}

func cleanCSVFormat(format string) string {
	return strings.Join(strings.Fields(format), " ")
}

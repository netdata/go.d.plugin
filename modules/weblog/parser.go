package weblog

import (
	"errors"
	"io"
	"regexp"

	"github.com/netdata/go.d.plugin/pkg/logs"
)

var (
	reLTSV = regexp.MustCompile(`^[a-zA-Z0-9]+:[^\t]*(\t[a-zA-Z0-9]+:[^\t]*)*$`)

	csvCommon = `           $remote_addr - $remote_user [time local] "$request" $status $body_bytes_sent`
	//csvCombined      = `           $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"`
	//csvCustom1       = `           $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" $uid_got                  $request_time`
	//csvCustom2       = `           $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent $request_length $request_time'`
	//csvCustom3       = `           $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent $request_length $request_time      "$upstream_response_time"`
	//csvVhostCommon   = `$http_host $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent`
	//csvVhostCombined = `$http_host $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"`
	//csvVhostCustom1  = `$http_host $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" $uid_got                  $request_time`
	//csvVhostCustom2  = `$http_host $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent $request_length $request_time`
	//csvVhostCustom3  = `$http_host $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent $request_length $request_time      "$upstream_response_time"`

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

func (w *WebLog) guessParser(record []byte) logs.Guesser {
	f := func(config logs.ParserConfig, in io.Reader) (parser logs.Parser, e error) {
		if reLTSV.Match(record) {
			return logs.NewLTSVParser(config.LTSV, in)
		}

		for _, format := range guessOrder {
			cfg := config.CSV
			cfg.Format = format
			cfg.TrimLeadingSpace = true

			parser, err := logs.NewCSVParser(cfg, in)
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
	return logs.GuessFunc(f)
}

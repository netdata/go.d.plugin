package weblog

import (
	"errors"
	"fmt"
	"io"
	"regexp"

	"github.com/netdata/go.d.plugin/pkg/logs/parse"
)

var (
	//reSpace = regexp.MustCompile(`\s+`)
	reLTSV = regexp.MustCompile(`^[a-zA-Z0-9]+:[^\t]*(\t[a-zA-Z0-9]+:[^\t]*)*$`)

	csvCommon = `           $remote_addr - $remote_user [time local] "$request" $resp_status $body_bytes_sent`
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

//func removeDuplicateSpaces(s string) string {
//	return reSpace.ReplaceAllString(s, " ")
//}

var defaultParserConfig = parse.Config{
	LogType: parse.TypeAuto,
	CSV: parse.CSVConfig{
		Delimiter: ' ',
	},
	LTSV: parse.LTSVConfig{
		FieldDelimiter: '\t',
		ValueDelimiter: ':',
	},
	RegExp: parse.RegExpConfig{},
}

func guessParser(config parse.Config, in io.Reader, record []byte) (parse.Parser, error) {
	if reLTSV.Match(record) {
		return parse.NewLTSVParser(config.LTSV, in)
	}

	for _, format := range guessOrder {
		cfg := config.CSV
		cfg.Format = format
		cfg.TrimLeadingSpace = true

		parser, err := parse.NewCSVParser(cfg, in)
		if err != nil {
			return nil, err
		}

		line := newEmptyLogLine()
		err = parser.Parse(record, line)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(line)
		if err = line.Verify(); err != nil {
			fmt.Println(err)
			continue
		}
		return parser, nil
	}
	return nil, errors.New("cannot determine log format")
}

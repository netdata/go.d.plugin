package parser

//
//import (
//	"io"
//	"regexp"
//
//	"golang.org/x/xerrors"
//)
//
//type Error struct{ msg string }
//
//func (e Error) Error() string { return e.msg }
//
//type (
//	Parser interface {
//		ReadLine() (LogLine, error)
//		Parse(line []byte) (LogLine, error)
//	}
//)
//
//var (
//	reLTSV = regexp.MustCompile(`^[a-zA-Z0-9]+:[^\t]*(\t[a-zA-Z0-9]+:[^\t]*)*$`)
//
//	csvCommon        = `           $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent`
//	csvCombined      = `           $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"`
//	csvCustom1       = `           $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" $uid_got                  $request_time'`
//	csvCustom2       = `           $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent $request_length $request_time'`
//	csvCustom3       = `           $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent $request_length $request_time      "$upstream_response_time"'`
//	csvVhostCommon   = `$http_host $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent`
//	csvVhostCombined = `$http_host $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"`
//	csvVhostCustom1  = `$http_host $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" $uid_got                  $request_time`
//	csvVhostCustom2  = `$http_host $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent $request_length $request_time`
//	csvVhostCustom3  = `$http_host $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent $request_length $request_time      "$upstream_response_time"`
//
//	guessOrder = []string{
//		csvVhostCustom1,
//		csvVhostCustom3,
//		csvVhostCustom2,
//		csvVhostCombined,
//		csvVhostCommon,
//		csvCustom1,
//		csvCustom3,
//		csvCustom2,
//		csvCombined,
//		csvCommon,
//	}
//)
//
//func NewParser(config Config, in io.Reader, line []byte) (Parser, error) {
//	switch config.LogType {
//	case TypeAuto:
//		return guessParser(config, in, line)
//	case TypeCSV:
//		return newCSVParser(config, in)
//	case TypeLTSV:
//		return newLTSVParser(config, in)
//	case TypeRegExp:
//		return newRegExpParser(config, in)
//	default:
//		return nil, xerrors.Errorf("invalid type: %q", config.LogType)
//	}
//}
//
//func guessParser(config Config, in io.Reader, line []byte) (Parser, error) {
//	if reLTSV.Match(line) {
//		return newLTSVParser(config, in)
//	}
//	for _, format := range guessOrder {
//		cfg := config
//		cfg.CSV.Format = format
//		parser, _ := newCSVParser(cfg, in)
//		log, err := parser.Parse(line)
//		if err == nil && log.verify() == nil {
//			return parser, nil
//		}
//	}
//	return nil, xerrors.New("cannot determine log format")
//}

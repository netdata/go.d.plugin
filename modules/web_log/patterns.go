package web_log

import (
	"regexp"
	"strings"
)

var reRequest = regexp.MustCompile(`(?P<method>[A-Z]+) (?P<url>[^ ]+) [A-Z]+/(?P<http_version>\d(?:.\d)?)`)

var (
	lastHop = strings.Join([]string{
		`(?P<address>[\da-f.:]+|localhost) -.*?"`,
		`(?P<request>[^"]*)" `,
		`(?P<code>[1-9]\d{2}) `,
		`(?P<bytes_sent>\d+|-)`,
	}, "")
	apacheV1 = strings.Join([]string{
		`(?P<address>[\da-f.:]+|localhost) -.*?"`,
		`(?P<request>[^"]*)" `,
		`(?P<code>[1-9]\d{2}) `,
		`(?P<bytes_sent>\d+|-) `,
		`(?P<resp_length>\d+|-) `,
		`(?P<resp_time>\d+) `,
	}, "")
	apacheV2 = strings.Join([]string{
		`(?P<address>[\da-f.:]+|localhost) -.*?"`,
		`(?P<request>[^"]*)" `,
		`(?P<code>[1-9]\d{2}) `,
		`(?P<bytes_sent>\d+|-) .*? `,
		`(?P<resp_length>\d+|-) `,
		`(?P<resp_time>\d+)(?: |$)`,
	}, "")
	nginxV1 = strings.Join([]string{
		`(?P<address>[\da-f.:]+) -.*?"`,
		`(?P<request>[^"]*)" `,
		`(?P<code>[1-9]\d{2}) `,
		`(?P<bytes_sent>\d+) `,
		`(?P<resp_length>\d+) `,
		`(?P<resp_time>\d+\.\d+) `,
		`(?P<resp_time_upstream>[\d.-]+) `,
	}, "")
	nginxV2 = strings.Join([]string{
		`(?P<address>[\da-f.:]+) -.*?"`,
		`(?P<request>[^"]*)" `,
		`(?P<code>[1-9]\d{2}) `,
		`(?P<bytes_sent>\d+) `,
		`(?P<resp_length>\d+) `,
		`(?P<resp_time>\d+\.\d+) `,
	}, "")
	nginxV3 = strings.Join([]string{
		`(?P<address>[\da-f.:]+) -.*?"`,
		`(?P<request>[^"]*)" `,
		`(?P<code>[1-9]\d{2}) `,
		`(?P<bytes_sent>\d+) .*? `,
		`(?P<resp_length>\d+) `,
		`(?P<resp_time>\d+\.\d+)`,
	}, "")
)

var patterns = []*regexp.Regexp{
	regexp.MustCompile(apacheV1),
	regexp.MustCompile(apacheV2),
	regexp.MustCompile(nginxV1),
	regexp.MustCompile(nginxV2),
	regexp.MustCompile(nginxV3),
	regexp.MustCompile(lastHop),
}

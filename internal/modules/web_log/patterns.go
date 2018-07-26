package web_log

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/l2isbad/go.d.plugin/internal/pkg/utils"
)

var (
	mandatoryKey = keyCode
	reRequest    = regexp.MustCompile(`(?P<method>[A-Z]+) (?P<url>[^ ]+) [A-Z]+/(?P<http_version>\d(?:.\d)?)`)
)

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

func getPattern(custom string, line []byte) (*regexp.Regexp, error) {
	if custom == "" {
		for _, p := range patterns {
			if p.Match(line) {
				return p, nil
			}
		}
		return nil, errors.New("can not find appropriate regex, consider using 'custom_log_format' feature")
	}
	r, err := regexp.Compile(custom)
	if err != nil {
		return nil, err
	}
	if len(r.SubexpNames()) == 1 {
		return nil, errors.New("custom regex contains no named groups (?P<subgroup_name>)")
	}

	if !utils.StringSlice(r.SubexpNames()).Include(mandatoryKey) {
		return nil, fmt.Errorf("custom regex missing mandatory key '%s'", mandatoryKey)
	}

	if !r.Match(line) {
		return nil, errors.New("custom regex match fails")
	}

	return r, nil
}

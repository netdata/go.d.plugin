package nginx

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/stm"
)

const (
	connActive  = "connActive"
	connAccepts = "connAccepts"
	connHandled = "connHandled"
	requests    = "requests"
	requestTime = "requestTime"
	connReading = "connReading"
	connWriting = "connWriting"
	connWaiting = "connWaiting"
)

var (
	nginxSeq = []string{
		connActive,
		connAccepts,
		connHandled,
		requests,
		connReading,
		connWriting,
		connWaiting,
	}
	tengineSeq = []string{
		connActive,
		connAccepts,
		connHandled,
		requests,
		requestTime,
		connReading,
		connWriting,
		connWaiting,
	}

	reStatus = regexp.MustCompile(`^Active connections: ([0-9]+)\n[^\d]+([0-9]+) ([0-9]+) ([0-9]+) ?([0-9]+)?\nReading: ([0-9]+) Writing: ([0-9]+) Waiting: ([0-9]+)`)
)

func (n *Nginx) collect() (map[string]int64, error) {
	raw, err := n.apiClient.getStubStatus()

	if err != nil {
		return nil, err
	}

	var status stubStatus

	err = parseLines(&status, raw)

	if err != nil {
		return nil, fmt.Errorf("error on parsing response : %v", err)
	}

	return stm.ToMap(status), nil
}

func parseLines(status *stubStatus, lines []string) error {
	s := strings.Join(lines, "\n")
	parsed := reStatus.FindStringSubmatch(s)

	if len(parsed) == 0 {
		return fmt.Errorf("can't parse '%v'", lines)
	}

	parsed = parsed[1:]
	var seq []string

	switch len(parsed) {
	default:
		return fmt.Errorf("invalid number of fields, got %d, expect %d or %d", len(parsed), len(nginxSeq), len(tengineSeq))
	case len(nginxSeq):
		seq = nginxSeq
	case len(tengineSeq):
		seq = tengineSeq
	}

	for i, key := range seq {
		value := parsed[i]
		if value == "" {
			continue
		}
		switch key {
		default:
			return fmt.Errorf("unknown key in seq : %s", key)
		case connActive:
			status.Connections.Active = mustParseInt(value)
		case connAccepts:
			status.Connections.Accepts = mustParseInt(value)
		case connHandled:
			status.Connections.Handled = mustParseInt(value)
		case requests:
			status.Requests.Total = mustParseInt(value)
		case connReading:
			status.Connections.Reading = mustParseInt(value)
		case connWriting:
			status.Connections.Writing = mustParseInt(value)
		case connWaiting:
			status.Connections.Waiting = mustParseInt(value)
		case requestTime:
			v := mustParseInt(value)
			status.Requests.Time = &v
		}
	}

	return nil
}

func mustParseInt(value string) int64 {
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(err)
	}
	return v
}

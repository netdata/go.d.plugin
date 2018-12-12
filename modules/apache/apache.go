package apache

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/web"
)

// -- Extended On --
// Total Accesses: 7
// Total kBytes: 5
// Uptime: 6
// ReqPerSec: 1.16667
// BytesPerSec: 853.333
// BytesPerReq: 731.429
// BusyWorkers: 1
// IdleWorkers: 49
// ConnsTotal: 1
// ConnsAsyncWriting: 0
// ConnsAsyncKeepAlive: 1
// ConnsAsyncClosing: 0

// -- Extended Off --
// BusyWorkers: 1
// IdleWorkers: 49
// ConnsTotal: 1
// ConnsAsyncWriting: 0
// ConnsAsyncKeepAlive: 1
// ConnsAsyncClosing: 0

func init() {
	creator := modules.Creator{
		Create: func() modules.Module { return New() },
	}

	modules.Register("apache", creator)
}

var assignment = map[string]string{
	"BytesPerReq":         "size_req",
	"IdleWorkers":         "idle",
	"IdleServers":         "dle_servers",
	"BusyWorkers":         "busy",
	"BusyServers":         "busy_servers",
	"ReqPerSec":           "requests_sec",
	"BytesPerSec":         "size_sec",
	"Total Accesses":      "requests",
	"Total kBytes":        "sent",
	"ConnsTotal":          "connections",
	"ConnsAsyncKeepAlive": "keepalive",
	"ConnsAsyncClosing":   "closing",
	"ConnsAsyncWriting":   "writing",
}

// New creates Apache with default values
func New() *Apache {
	return &Apache{
		metrics: make(map[string]int64),
	}
}

// Apache apache module
type Apache struct {
	modules.Base // should be embedded by every module

	web.HTTP `yaml:",inline"`

	request *http.Request
	client  web.Client

	metrics map[string]int64
}

func (Apache) Cleanup() {

}

func (a *Apache) Init() bool {
	req, err := a.CreateHTTPRequest()

	if err != nil {
		a.Error(err)
		return false
	}

	if a.Timeout.Duration == 0 {
		a.Timeout.Duration = time.Second
	}

	a.request = req
	a.client = a.CreateHTTPClient()

	return true
}

func (Apache) Check() bool {
	return false
}

func (Apache) Charts() *modules.Charts {
	return nil
}

func (a *Apache) GatherMetrics() map[string]int64 {
	resp, err := a.doRequest()

	if err != nil {
		a.Error(err)
		return nil
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if err := a.parseResponse(resp); err != nil {
		a.Error(err)
		return nil
	}

	return a.metrics
}

func (a *Apache) doRequest() (*http.Response, error) {
	return a.client.Do(a.request)
}

func (a *Apache) parseResponse(resp *http.Response) error {
	s := bufio.NewScanner(resp.Body)
	var parsed int

	for s.Scan() {
		if err := parseLine(s.Text(), a.metrics); err != nil {
			continue
		}
		parsed++
	}

	if parsed == 0 {
		return errors.New("unparsable data")
	}

	return nil
}

func parseLine(line string, metrics map[string]int64) error {
	if !strings.Contains(line, ":") {
		return errors.New("bad line format")
	}

	parts := strings.SplitN(line, ":", 2)

	if len(parts) != 2 {
		return errors.New("bad line format")
	}

	key := strings.Replace(parts[0], " ", "", -1)
	value := strings.TrimSpace(parts[1])

	if newKey, ok := assignment[key]; ok {
		key = newKey
	}

	switch key {
	default:
		v, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		metrics[key] = int64(v)
	case "size_req", "requests_sec", "size_sec":
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err

		}
		metrics[key] = int64(v * 100000)
	case "Scoreboard":
		parseScoreboard(value, metrics)
	}

	return nil
}

func parseScoreboard(scoreboard string, metrics map[string]int64) {
	//  “_” Waiting for Connection
	// “S” Starting up
	// “R” Reading Request
	// “W” Sending Reply
	// “K” Keepalive (read)
	// “D” DNS Lookup
	// “C” Closing connection
	// “L” Logging
	// “G” Gracefully finishing
	// “I” Idle cleanup of worker
	//“.” Open slot with no current process

	var waiting, open, S, R, W, K, D, C, L, G, I int64

	for _, s := range strings.Split(scoreboard, "") {

		switch s {
		case "_":
			waiting++
		case "S":
			S++
		case "R":
			R++
		case "W":
			W++
		case "K":
			K++
		case "D":
			D++
		case "C":
			C++
		case "L":
			L++
		case "G":
			G++
		case "I":
			I++
		case ".":
			open++
		}
	}

	metrics["scoreboard_waiting"] = waiting
	metrics["scoreboard_starting"] = S
	metrics["scoreboard_reading"] = R
	metrics["scoreboard_sending"] = W
	metrics["scoreboard_keepalive"] = K
	metrics["scoreboard_dns_lookup"] = D
	metrics["scoreboard_closing"] = C
	metrics["scoreboard_logging"] = L
	metrics["scoreboard_finishing"] = G
	metrics["scoreboard_idle_cleanup"] = I
	metrics["scoreboard_open"] = open
}

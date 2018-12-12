package apache

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	creator := modules.Creator{
		Create: func() modules.Module { return New() },
	}

	modules.Register("apache", creator)
}

const (
	totalAccesses       = "Total Accesses" // extended stats
	totalkBytes         = "Total kBytes"   // extended stats
	cpuLoad             = "CPULoad"        // extended stats
	uptime              = "Uptime"         // extended stats
	reqPerSec           = "ReqPerSec"      // extended stats
	bytesPerSec         = "BytesPerSec"    // extended stats
	bytesPerReq         = "BytesPerReq"    // extended stats
	busyWorkers         = "BusyWorkers"
	idleWorkers         = "IdleWorkers"
	connsTotal          = "ConnsTotal"
	connsAsyncKeepAlive = "ConnsAsyncKeepAlive"
	connsAsyncClosing   = "ConnsAsyncClosing"
	connsAsyncWriting   = "ConnsAsyncWriting"
	scoreBoard          = "Scoreboard"
)

var assignment = map[string]string{
	totalAccesses:       "requests",
	totalkBytes:         "sent",
	reqPerSec:           "requests_sec",
	bytesPerSec:         "size_sec",
	bytesPerReq:         "size_req",
	busyWorkers:         "busy",
	idleWorkers:         "idle",
	connsTotal:          "connections",
	connsAsyncKeepAlive: "keepalive",
	connsAsyncClosing:   "closing",
	connsAsyncWriting:   "writing",
}

// New creates Apache with default values
func New() *Apache {
	return &Apache{}
}

// Apache apache module
type Apache struct {
	modules.Base // should be embedded by every module

	web.HTTP `yaml:",inline"`

	request *http.Request
	client  web.Client

	extendedStats bool

	metrics map[string]int64
}

func (Apache) Cleanup() {}

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

func (a *Apache) Check() bool {
	if len(a.GatherMetrics()) == 0 {
		return false
	}

	_, a.extendedStats = a.metrics[assign(totalAccesses)]

	if !a.extendedStats {
		a.Info("extended status is disabled, not all metrics are available")
	}

	return true
}

func (a Apache) Charts() *modules.Charts {
	charts := charts.Copy()

	if !a.extendedStats {
		charts.Remove("requests")
		charts.Remove("net")
		charts.Remove("reqpersec")
		charts.Remove("bytespersec")
		charts.Remove("bytesperreq")
	}

	return charts
}

func (a *Apache) GatherMetrics() map[string]int64 {
	resp, err := a.doRequest()

	if err != nil {
		a.Errorf("error on request to %s : %s", a.request.URL, err)
		return nil
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		a.Errorf("%s returned HTTP status %d", a.request.URL, resp.StatusCode)
		return nil
	}

	// Extended On -> Extended Off
	a.metrics = make(map[string]int64)

	if err := a.parseResponse(resp); err != nil {
		a.Errorf("error on parse response : %s", err)
		return nil
	}

	return a.metrics
}

func (a *Apache) doRequest() (*http.Response, error) {
	return a.client.Do(a.request)
}

func (a *Apache) parseResponse(resp *http.Response) error {
	s := bufio.NewScanner(resp.Body)

	for s.Scan() {
		if err := parseLine(s.Text(), a.metrics); err != nil {
			return err
		}
	}

	return nil
}

func parseLine(line string, metrics map[string]int64) error {
	parts := strings.SplitN(line, ":", 2)

	if len(parts) != 2 {
		return fmt.Errorf("bad format : %s", line)
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	switch key {
	case cpuLoad, uptime:
	case totalAccesses, totalkBytes, busyWorkers, idleWorkers, connsTotal:
		fallthrough
	case connsAsyncWriting, connsAsyncKeepAlive, connsAsyncClosing:
		v, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		metrics[assign(key)] = int64(v)
	case reqPerSec, bytesPerSec, bytesPerReq:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		metrics[assign(key)] = int64(v * 100000)
	case scoreBoard:
		parseScoreboard(value, metrics)
	default:
		return fmt.Errorf("unknown key: %s", key)
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

func assign(key string) string {
	if v, ok := assignment[key]; ok {
		return v
	}
	return key
}

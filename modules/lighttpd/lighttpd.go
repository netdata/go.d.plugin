package lighttpd

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
	"github.com/netdata/go.d.plugin/pkg/utils"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	creator := modules.Creator{
		Create: func() modules.Module { return New() },
	}

	modules.Register("lighttpd", creator)
}

const (
	totalAccesses = "Total Accesses"
	totalkBytes   = "Total kBytes"
	uptime        = "Uptime"
	busyServers   = "BusyServers"
	idleServers   = "IdleServers"
	scoreBoard    = "Scoreboard"
)

var assignment = map[string]string{
	totalAccesses: "requests",
	totalkBytes:   "sent",
	busyServers:   "busy",
	idleServers:   "idle",
	uptime:        "uptime",
}

// New creates Lighttpd with default values
func New() *Lighttpd {
	return &Lighttpd{
		HTTP: web.HTTP{
			RawRequest: web.RawRequest{
				URL: "http://localhost/server-status?auto",
			},
			RawClient: web.RawClient{
				Timeout: utils.Duration{Duration: time.Second},
			},
		},
		metrics: make(map[string]int64),
	}
}

// Lighttpd lighttpd module
type Lighttpd struct {
	modules.Base // should be embedded by every module

	web.HTTP `yaml:",inline"`

	request *http.Request
	client  web.Client

	metrics map[string]int64
}

// Cleanup makes cleanup
func (Lighttpd) Cleanup() {}

// Init makes initialization
func (a *Lighttpd) Init() bool {
	req, err := a.CreateHTTPRequest()

	if err != nil {
		a.Error(err)
		return false
	}

	a.request = req
	a.client = a.CreateHTTPClient()

	return true
}

// Check makes check
func (a *Lighttpd) Check() bool {
	return len(a.Collect()) > 0
}

// Charts creates Charts
func (a Lighttpd) Charts() *modules.Charts {
	return charts.Copy()
}

// Collect collects metrics
func (a *Lighttpd) Collect() map[string]int64 {
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

	if err := a.parseResponse(resp); err != nil {
		a.Errorf("error on parse response : %s", err)
		return nil
	}

	return a.metrics
}

func (a *Lighttpd) doRequest() (*http.Response, error) {
	return a.client.Do(a.request)
}

func (a *Lighttpd) parseResponse(resp *http.Response) error {
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
	case uptime, totalAccesses, totalkBytes, busyServers, idleServers:
		v, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		metrics[assign(key)] = int64(v)
	case scoreBoard:
		parseScoreboard(value, metrics)
	default:
		return fmt.Errorf("unknown key: %s", key)
	}
	return nil
}

func parseScoreboard(scoreboard string, metrics map[string]int64) {
	// Descriptions from https://blog.serverdensity.com/monitor-lighttpd/
	//
	// “.” = Opening the TCP connection (connect)
	// “C” = Closing the TCP connection if no other HTTP request will use it (close)
	// “E” = hard error
	// “k” = Keeping the TCP connection open for more HTTP requests from the same client to avoid the TCP handling overhead (keep-alive)
	// “r” = Read the content of the HTTP request (read)
	// “R” = Read the content of the HTTP request (read-POST)
	// “W” = Write the HTTP response to the socket (write)
	// “h” = Decide action to take with the request (handle-request)
	// “q” = Start of HTTP request (request-start)
	// “Q” = End of HTTP request (request-end)
	// “s” = Start of the HTTP request response (response-start)
	// “S” = End of the HTTP request response (response-end)
	// “_” Waiting for Connection (NOTE: not sure, copied the description from apache score board)

	var waiting, open, C, E, k, r, R, W, h, q, Q, s, S int64

	for _, v := range strings.Split(scoreboard, "") {

		switch v {
		case "_":
			waiting++
		case ".":
			open++
		case "C":
			C++
		case "E":
			E++
		case "k":
			k++
		case "r":
			r++
		case "R":
			R++
		case "W":
			W++
		case "h":
			h++
		case "q":
			q++
		case "Q":
			Q++
		case "s":
			s++
		case "S":
			S++
		}
	}

	metrics["scoreboard_waiting"] = waiting
	metrics["scoreboard_open"] = open
	metrics["scoreboard_close"] = C
	metrics["scoreboard_hard_error"] = E
	metrics["scoreboard_keepalive"] = k
	metrics["scoreboard_read"] = r
	metrics["scoreboard_read_post"] = R
	metrics["scoreboard_write"] = W
	metrics["scoreboard_handle_request"] = h
	metrics["scoreboard_request_start"] = q
	metrics["scoreboard_request_end"] = Q
	metrics["scoreboard_response_start"] = s
	metrics["scoreboard_response_end"] = S
}

func assign(key string) string {
	if v, ok := assignment[key]; ok {
		return v
	}
	return key
}

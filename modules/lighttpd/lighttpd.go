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
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	creator := modules.Creator{
		Create: func() modules.Module { return New() },
	}

	modules.Register("lighttpd", creator)
}

// New creates Lighttpd with default values
func New() *Lighttpd {
	var (
		defURL         = "http://localhost/server-status?auto"
		defHTTPTimeout = time.Second
	)

	return &Lighttpd{
		HTTP: web.HTTP{
			Request: web.Request{URL: defURL},
			Client:  web.Client{Timeout: web.Duration{Duration: defHTTPTimeout}},
		},
		metrics: make(map[string]int64),
	}
}

const (
	totalAccesses = "Total Accesses"
	totalkBytes   = "Total kBytes"
	uptime        = "Uptime"
	busyServers   = "BusyServers"
	idleServers   = "IdleServers"
	scoreBoard    = "Scoreboard"
)

// Lighttpd lighttpd module
type Lighttpd struct {
	modules.Base // should be embedded by every module

	web.HTTP `yaml:",inline"`

	request *http.Request
	client  web.HTTPClient

	metrics map[string]int64
}

// Cleanup makes cleanup
func (Lighttpd) Cleanup() {}

// Init makes initialization
func (l *Lighttpd) Init() bool {
	if !strings.HasSuffix(l.URL, "?auto") {
		l.Errorf("invalid status page URL %s, must end with '?auto'", l.URL)
		return false
	}

	var err error

	// create HTTP request
	if l.request, err = web.NewHTTPRequest(l.Request); err != nil {
		l.Errorf("error on creating request to %s : %s", l.URL, err)
		return false
	}

	// create HTTP client
	l.client = web.NewHTTPClient(l.Client)

	// post Init debug info
	l.Debugf("using URL %s", l.request.URL)
	l.Debugf("using timeout: %s", l.Timeout.Duration)

	return true
}

// Check makes check
func (l *Lighttpd) Check() bool {
	return len(l.Collect()) > 0
}

// Charts creates Charts
func (l Lighttpd) Charts() *modules.Charts {
	return charts.Copy()
}

// Collect collects metrics
func (l *Lighttpd) Collect() map[string]int64 {
	resp, err := l.doRequest()

	if err != nil {
		l.Errorf("error on request to %s : %s", l.request.URL, err)
		return nil
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		l.Errorf("%s returned HTTP status %d", l.request.URL, resp.StatusCode)
		return nil
	}

	if err := l.parseResponse(resp); err != nil {
		l.Errorf("error on parse response : %s", err)
		return nil
	}

	return l.metrics
}

func (l *Lighttpd) doRequest() (*http.Response, error) {
	return l.client.Do(l.request)
}

func (l *Lighttpd) parseResponse(resp *http.Response) error {
	s := bufio.NewScanner(resp.Body)

	for s.Scan() {
		if err := parseLine(s.Text(), l.metrics); err != nil {
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

var assignment = map[string]string{
	totalAccesses: "requests",
	totalkBytes:   "sent",
	busyServers:   "busy",
	idleServers:   "idle",
	uptime:        "uptime",
}

func assign(key string) string {
	if v, ok := assignment[key]; ok {
		return v
	}
	return key
}

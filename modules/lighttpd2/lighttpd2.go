package lighttpd2

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

	modules.Register("lighttpd2", creator)
}

// New creates Lighttpd2 with default values
func New() *Lighttpd2 {
	var (
		defURL         = "http://localhost/server-status?format=plain"
		defHTTPTimeout = time.Second
	)

	return &Lighttpd2{
		HTTP: web.HTTP{
			Request: web.Request{URL: defURL},
			Client:  web.Client{Timeout: web.Duration{Duration: defHTTPTimeout}},
		},
		metrics: make(map[string]int64),
	}
}

// Lighttpd2 lighttpd module
type Lighttpd2 struct {
	modules.Base // should be embedded by every module

	web.HTTP `yaml:",inline"`

	request *http.Request
	client  *http.Client

	metrics map[string]int64
}

// Cleanup makes cleanup
func (Lighttpd2) Cleanup() {}

// Init makes initialization
func (l *Lighttpd2) Init() bool {
	if !strings.HasSuffix(l.URL, "?format=plain") {
		l.Errorf("invalid status page URL %s, must end with '?format=plain'", l.URL)
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
func (l *Lighttpd2) Check() bool {
	return len(l.Collect()) > 0
}

// Charts creates Charts
func (l Lighttpd2) Charts() *modules.Charts {
	return charts.Copy()
}

// Collect collects metrics
func (l *Lighttpd2) Collect() map[string]int64 {
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

func (l *Lighttpd2) doRequest() (*http.Response, error) {
	return l.client.Do(l.request)
}

func (l *Lighttpd2) parseResponse(resp *http.Response) error {
	var parsed int

	s := bufio.NewScanner(resp.Body)

	for s.Scan() {
		line := s.Text()

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if err := parseLine(line, l.metrics); err != nil {
			return err
		}
		parsed++
	}

	if parsed == 0 {
		return fmt.Errorf("nothing has been parsed")
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

	if !validKeys[key] {
		return fmt.Errorf("unknown key: %s", key)
	}

	v, err := strconv.Atoi(value)
	if err != nil {
		return err
	}

	metrics[key] = int64(v)

	return nil
}

var validKeys = map[string]bool{
	"uptime":                          true,
	"memory_usage":                    true,
	"requests_abs":                    true,
	"traffic_out_abs":                 true,
	"traffic_in_abs":                  true,
	"connections_abs":                 true,
	"requests_avg":                    true,
	"traffic_out_avg":                 true,
	"traffic_in_avg":                  true,
	"connections_avg":                 true,
	"requests_avg_5sec":               true,
	"traffic_out_avg_5sec":            true,
	"traffic_in_avg_5sec":             true,
	"connections_avg_5sec":            true,
	"connection_state_start":          true,
	"connection_state_read_header":    true,
	"connection_state_handle_request": true,
	"connection_state_write_response": true,
	"connection_state_keep_alive":     true,
	"connection_state_upgraded":       true,
	"status_1xx":                      true,
	"status_2xx":                      true,
	"status_3xx":                      true,
	"status_4xx":                      true,
	"status_5xx":                      true,
}

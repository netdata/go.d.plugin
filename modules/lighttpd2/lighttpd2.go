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
	"github.com/netdata/go.d.plugin/pkg/utils"
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
	return &Lighttpd2{
		HTTP: web.HTTP{
			RawRequest: web.RawRequest{
				URL: "http://localhost/server-status?format=plain",
			},
			RawClient: web.RawClient{
				Timeout: utils.Duration{Duration: time.Second},
			},
		},
		metrics: make(map[string]int64),
	}
}

// Lighttpd2 lighttpd module
type Lighttpd2 struct {
	modules.Base // should be embedded by every module

	web.HTTP `yaml:",inline"`

	request *http.Request
	client  web.Client

	metrics map[string]int64
}

// Cleanup makes cleanup
func (Lighttpd2) Cleanup() {}

// Init makes initialization
func (a *Lighttpd2) Init() bool {
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
func (a *Lighttpd2) Check() bool {
	return len(a.Collect()) > 0
}

// Charts creates Charts
func (a Lighttpd2) Charts() *modules.Charts {
	return charts.Copy()
}

// Collect collects metrics
func (a *Lighttpd2) Collect() map[string]int64 {
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

func (a *Lighttpd2) doRequest() (*http.Response, error) {
	return a.client.Do(a.request)
}

func (a *Lighttpd2) parseResponse(resp *http.Response) error {
	var parsed int

	s := bufio.NewScanner(resp.Body)

	for s.Scan() {
		line := s.Text()

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if err := parseLine(line, a.metrics); err != nil {
			return err
		}
		parsed++
	}

	if parsed == 0 {
		return fmt.Errorf("nothing parsed")
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

package nginx

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
	"github.com/netdata/go.d.plugin/pkg/utils"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	creator := modules.Creator{
		Create: func() modules.Module { return New() },
	}

	modules.Register("nginx", creator)
}

type metrics struct {
	Active   int `stm:"active"`
	Requests int `stm:"requests"`
	Reading  int `stm:"reading"`
	Writing  int `stm:"writing"`
	Waiting  int `stm:"waiting"`
	Accepts  int `stm:"accepts"`
	Handled  int `stm:"handled"`
}

// New creates Nginx with default values
func New() *Nginx {
	return &Nginx{}
}

// Nginx nginx module
type Nginx struct {
	modules.Base // should be embedded by every module

	web.HTTP `yaml:",inline"`

	request *http.Request
	client  web.Client

	metrics metrics
}

// Cleanup makes cleanup
func (Nginx) Cleanup() {}

// Init makes initialization
func (n *Nginx) Init() bool {
	req, err := n.CreateHTTPRequest()

	if err != nil {
		n.Error(err)
		return false
	}

	if n.Timeout.Duration == 0 {
		n.Timeout.Duration = time.Second
	}

	n.request = req
	n.client = n.CreateHTTPClient()

	return true
}

// Check makes check
func (n *Nginx) Check() bool {
	return n.GatherMetrics() != nil
}

// Charts creates Charts
func (Nginx) Charts() *Charts {
	return charts.Copy()
}

// GatherMetrics gathers metrics
func (n *Nginx) GatherMetrics() map[string]int64 {
	resp, err := n.doRequest()

	if err != nil {
		n.Error(err)
		return nil
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		n.Errorf("%s returned HTTP status %d", n.request.URL.String(), resp.StatusCode)
		return nil
	}

	if err := n.processResponse(resp); err != nil {
		n.Error(err)
		return nil
	}

	return utils.ToMap(n.metrics)
}

func (n *Nginx) doRequest() (*http.Response, error) {
	return n.client.Do(n.request)
}

func (n *Nginx) processResponse(resp *http.Response) error {
	// Active connections: 2
	//server accepts handled requests
	// 2 2 3
	//Reading: 0 Writing: 1 Waiting: 1

	reader := bufio.NewReader(resp.Body)

	// Active connections: 2
	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	// [Active connections: 2]
	slice := strings.Fields(line)
	if len(slice) != 3 {
		return errors.New("not enough fields")
	}

	// active connections
	n.metrics.Active, err = strconv.Atoi(slice[2])
	if err != nil {
		return err
	}

	// server accepts handled requests
	line, err = reader.ReadString('\n')
	if err != nil {
		return err
	}

	// 2 2 3
	line, err = reader.ReadString('\n')
	if err != nil {
		return err
	}

	// [2 2 3]
	slice = strings.Fields(line)

	if len(slice) != 3 {
		return err
	}

	n.metrics.Accepts, err = strconv.Atoi(slice[0])
	if err != nil {
		return err
	}

	n.metrics.Handled, err = strconv.Atoi(slice[1])
	if err != nil {
		return err
	}

	n.metrics.Requests, err = strconv.Atoi(slice[2])
	if err != nil {
		return err
	}

	//Reading: 0 Writing: 1 Waiting: 1
	line, err = reader.ReadString('\n')

	if err != nil {
		return err
	}

	slice = strings.Fields(line)

	if len(slice) != 6 {
		return errors.New("not enough fields")
	}

	n.metrics.Reading, err = strconv.Atoi(slice[1])
	if err != nil {
		return err
	}

	n.metrics.Writing, err = strconv.Atoi(slice[3])
	if err != nil {
		return err
	}

	n.metrics.Waiting, err = strconv.Atoi(slice[5])
	if err != nil {
		return err
	}

	return nil
}

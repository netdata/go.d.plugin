package activemq

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	creator := modules.Creator{
		Create: func() modules.Module { return New() },
	}

	modules.Register("activemq", creator)
}

var (
	uriStats    = "/%s/xml/%s.jsp"
	queuesStats = "queues"
	topicsStats = "topics"
)

var (
	defURL         = "http://127.0.0.1:8161"
	defHTTPTimeout = time.Second
)

// New creates Example with default values
func New() *Activemq {
	return &Activemq{
		HTTP: web.HTTP{
			Request: web.Request{URL: defURL},
			Client:  web.Client{Timeout: web.Duration{Duration: defHTTPTimeout}},
		},

		charts:       &Charts{},
		activeQueues: make(map[string]bool),
		activeTopics: make(map[string]bool),
		metrics:      make(map[string]int64),
	}
}

type topics struct {
	XMLName xml.Name `xml:"topics"`
	Items   []topic  `xml:"topic"`
}

type topic struct {
	XMLName xml.Name `xml:"topic"`
	Name    string   `xml:"name,attr"`
	Stats   stats    `xml:"stats"`
}

type queues struct {
	XMLName xml.Name `xml:"queues"`
	Items   []queue  `xml:"queue"`
}

type queue struct {
	XMLName xml.Name `xml:"queue"`
	Name    string   `xml:"name,attr"`
	Stats   stats    `xml:"stats"`
}

type stats struct {
	XMLName       xml.Name `xml:"stats"`
	Size          int      `xml:"size,attr"`
	ConsumerCount int      `xml:"consumerCount,attr"`
	EnqueueCount  int      `xml:"enqueueCount,attr"`
	DequeueCount  int      `xml:"dequeueCount,attr"`
}

// Activemq activemq module
type Activemq struct {
	modules.Base

	web.HTTP `yaml:",inline"`

	Webadmin string `yaml:"webadmin"`

	reqQueues *http.Request
	reqTopics *http.Request
	client    *http.Client

	activeQueues map[string]bool
	activeTopics map[string]bool

	charts  *Charts
	metrics map[string]int64
}

// Cleanup makes cleanup
func (Activemq) Cleanup() {}

// Init makes initialization
func (a *Activemq) Init() bool {
	if a.Webadmin == "" {
		a.Error("webadmin root path not specified")
		return false
	}

	if err := a.createRequests(); err != nil {
		a.Error(err)
		return false
	}

	a.client = web.NewHTTPClient(a.Client)

	return true
}

// Check makes check
func (a *Activemq) Check() bool {
	return len(a.Collect()) > 0
}

// Charts creates Charts
func (a Activemq) Charts() *Charts {
	return a.charts
}

// Collect collects metrics
func (a *Activemq) Collect() map[string]int64 {
	a.metrics = make(map[string]int64)

	var (
		q   *queues
		t   *topics
		err error
	)

	if q, err = a.collectQueues(); err != nil {
		a.Error(err)
		return nil
	}

	if t, err = a.collectTopics(); err != nil {
		a.Error(err)
		return nil
	}

	if err = a.processQueues(q); err != nil {
		a.Error(err)
		return nil
	}

	if err = a.processTopics(t); err != nil {
		a.Error(err)
		return nil
	}

	return a.metrics
}

func (a *Activemq) createRequests() (err error) {
	a.URI = fmt.Sprintf(uriStats, a.Webadmin, queuesStats)
	a.reqQueues, err = web.NewHTTPRequest(a.Request)

	if err != nil {
		return fmt.Errorf("error on creating HTTP request : %s", err)
	}

	a.URI = fmt.Sprintf(uriStats, a.Webadmin, topicsStats)
	a.reqTopics, err = web.NewHTTPRequest(a.Request)

	if err != nil {
		return fmt.Errorf("error on creating HTTP request : %s", err)
	}

	return nil
}

func (a *Activemq) doRequest(req *http.Request) (*http.Response, error) {
	return a.client.Do(req)
}

func (a *Activemq) getData(req *http.Request) ([]byte, error) {
	resp, err := a.doRequest(req)

	if err != nil {
		return nil, fmt.Errorf("error on request to %s : %s", req.URL, err)
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s returned HTTP status %d", req.URL, resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

func (a *Activemq) collectQueues() (*queues, error) {
	b, err := a.getData(a.reqQueues)

	if err != nil {
		return nil, err
	}

	var q queues

	if err := xml.Unmarshal(b, &q); err != nil {
		return nil, fmt.Errorf("error on decoding resp from %s : %s", a.reqQueues.URL, err)
	}

	return &q, nil

}

func (a *Activemq) collectTopics() (*topics, error) {
	b, err := a.getData(a.reqTopics)

	if err != nil {
		return nil, err
	}

	var t topics

	if err := xml.Unmarshal(b, &t); err != nil {
		return nil, fmt.Errorf("error on decoding resp from %s : %s", a.reqTopics.URL, err)
	}

	return &t, nil
}

func (a *Activemq) processQueues(q *queues) error {

	return nil
}

func (a *Activemq) processTopics(t *topics) error {

	return nil
}

func (a *Activemq) addQueueTopicCharts(name string, typ string) {
	charts := charts.Copy()

	for _, chart := range *charts {
		chart.ID = fmt.Sprintf("%s_%s", name, chart.ID)
		chart.Fam = typ

		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf("%s_%s", name, dim.ID)
		}
	}

	_ = a.charts.Add(*charts...)
}

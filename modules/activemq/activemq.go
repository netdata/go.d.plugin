package activemq

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
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
	keyQueues   = "queues"
	keyTopics   = "topics"
	keyAdvisory = "Advisory"

	nameReplacer = strings.NewReplacer(".", "_", " ", "")
)

var (
	defMaxQueues   = 999
	defMaxTopics   = 999
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

		MaxQueues: defMaxQueues,
		MaxTopics: defMaxTopics,

		charts:       &Charts{},
		activeQueues: make(map[string]bool),
		activeTopics: make(map[string]bool),
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
	Size          int64    `xml:"size,attr"`
	ConsumerCount int64    `xml:"consumerCount,attr"`
	EnqueueCount  int64    `xml:"enqueueCount,attr"`
	DequeueCount  int64    `xml:"dequeueCount,attr"`
}

// Activemq activemq module
type Activemq struct {
	modules.Base

	web.HTTP `yaml:",inline"`

	Webadmin  string `yaml:"webadmin"`
	MaxQueues int
	MaxTopics int

	reqQueues *http.Request
	reqTopics *http.Request
	client    *http.Client

	activeQueues map[string]bool
	activeTopics map[string]bool

	charts *Charts
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
	metrics := make(map[string]int64)

	var (
		queues queues
		topics topics
		err    error
	)

	if err = a.collect(a.reqQueues, &queues); err != nil {
		a.Error(err)
		return nil
	}

	if err = a.collect(a.reqTopics, &topics); err != nil {
		a.Error(err)
		return nil
	}

	a.processQueues(queues, metrics)
	a.processTopics(topics, metrics)

	return metrics
}

func (a *Activemq) createRequests() (err error) {
	a.URI = fmt.Sprintf(uriStats, a.Webadmin, keyQueues)
	a.reqQueues, err = web.NewHTTPRequest(a.Request)

	if err != nil {
		return fmt.Errorf("error on creating HTTP request : %s", err)
	}

	a.URI = fmt.Sprintf(uriStats, a.Webadmin, keyTopics)
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

func (a *Activemq) collect(req *http.Request, elem interface{}) error {
	b, err := a.getData(req)

	if err != nil {
		return err
	}

	if err := xml.Unmarshal(b, elem); err != nil {
		return fmt.Errorf("error on decoding resp from %s : %s", req.URL, err)
	}

	return nil
}

func (a *Activemq) processQueues(queues queues, metrics map[string]int64) {
	var (
		count   = len(a.activeQueues)
		updated = make(map[string]bool)
	)

	for _, q := range queues.Items {
		if strings.Contains(q.Name, keyAdvisory) {
			continue
		}

		if !a.activeQueues[q.Name] {
			if a.MaxQueues != 0 && count > a.MaxQueues {
				continue
			}
			a.activeQueues[q.Name] = true
			a.addQueueTopicCharts(q.Name, keyQueues)
		}

		rname := nameReplacer.Replace(q.Name)

		metrics["queue_"+rname+"_consumers"] = q.Stats.ConsumerCount
		metrics["queue_"+rname+"_enqueued"] = q.Stats.EnqueueCount
		metrics["queue_"+rname+"_dequeued"] = q.Stats.DequeueCount
		metrics["queue_"+rname+"_unprocessed"] = q.Stats.EnqueueCount - q.Stats.DequeueCount

		updated[q.Name] = true
	}

	for name := range a.activeQueues {
		if !updated[name] {
			delete(a.activeQueues, name)
			a.removeQueueTopicCharts(name, keyQueues)
		}
	}
}

func (a *Activemq) processTopics(topics topics, metrics map[string]int64) {
	var (
		count   = len(a.activeTopics)
		updated = make(map[string]bool)
	)

	for _, t := range topics.Items {
		if strings.Contains(t.Name, keyAdvisory) {
			continue
		}

		if !a.activeTopics[t.Name] {
			if a.MaxTopics != 0 && count > a.MaxTopics {
				continue
			}
			a.activeTopics[t.Name] = true
			a.addQueueTopicCharts(t.Name, keyTopics)
		}

		name := nameReplacer.Replace(t.Name)

		metrics["topic_"+name+"_consumers"] = t.Stats.ConsumerCount
		metrics["topic_"+name+"_enqueued"] = t.Stats.EnqueueCount
		metrics["topic_"+name+"_dequeued"] = t.Stats.DequeueCount
		metrics["topic_"+name+"_unprocessed"] = t.Stats.EnqueueCount - t.Stats.DequeueCount

		updated[t.Name] = true
	}

	for name := range a.activeTopics {
		if !updated[name] {
			delete(a.activeTopics, name)
			a.removeQueueTopicCharts(name, keyTopics)
		}
	}
}

func (a *Activemq) addQueueTopicCharts(name, typ string) {
	rname := nameReplacer.Replace(name)

	charts := charts.Copy()

	for _, chart := range *charts {
		chart.ID = fmt.Sprintf(chart.ID, typ, rname)
		chart.Title = fmt.Sprintf(chart.Title, name)
		chart.Fam = typ

		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, typ, rname)
		}
	}

	_ = a.charts.Add(*charts...)

}

func (a *Activemq) removeQueueTopicCharts(name, typ string) {
	rname := nameReplacer.Replace(name)

	chart := a.charts.Get(fmt.Sprintf("%s_%s_messages", typ, rname))
	chart.Obsolete = true
	chart.MarkNotCreated()
	chart.MarkRemove()

	chart = a.charts.Get(fmt.Sprintf("%s_%s_consumers", typ, rname))
	chart.Obsolete = true
	chart.MarkNotCreated()
	chart.MarkRemove()
}

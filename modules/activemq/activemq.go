package activemq

import (
	"encoding/xml"
	"fmt"
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
func (Activemq) Check() bool {
	return false
}

// Charts creates Charts
func (Activemq) Charts() *Charts {
	return nil
}

// Collect collects metrics
func (Activemq) Collect() map[string]int64 {
	return nil
}

func (a *Activemq) createRequests() (err error) {
	a.URI = fmt.Sprintf("/%s/xml/%s.jsp", a.Webadmin, "queues")
	if a.reqQueues, err = web.NewHTTPRequest(a.Request); err != nil {
		return fmt.Errorf("error on creating HTTP request : %s", err)
	}

	a.URI = fmt.Sprintf("/%s/xml/%s.jsp", a.Webadmin, "topics")
	if a.reqTopics, err = web.NewHTTPRequest(a.Request); err != nil {
		return fmt.Errorf("error on creating HTTP request : %s", err)
	}

	return nil
}

package activemq

import (
	"encoding/xml"

	"github.com/netdata/go.d.plugin/modules"
)

func init() {
	creator := modules.Creator{
		Create: func() modules.Module { return New() },
	}

	modules.Register("activemq", creator)
}

// New creates Example with default values
func New() *Activemq {
	return &Activemq{}
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
}

// Cleanup makes cleanup
func (Activemq) Cleanup() {}

// Init makes initialization
func (Activemq) Init() bool {
	return false
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

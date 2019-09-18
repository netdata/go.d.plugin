package hdfs

import (
	"bytes"
	"errors"
	"time"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("hdfs", creator)
}

// New creates HDFS with default values.
func New() *HDFS {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{
				UserURL: "http://127.0.0.1:8081/jmx",
			},
			Client: web.Client{
				Timeout: web.Duration{Duration: time.Second}},
		},
	}

	return &HDFS{
		Config: config,
	}
}

type nodeType int

const (
	nameNodeType nodeType = iota
	dataNodeType
	unknownNodeType
)

// Config is the HDFS module configuration.
type Config struct {
	web.HTTP `yaml:",inline"`
}

// HDFS HDFS module.
type HDFS struct {
	module.Base
	Config `yaml:",inline"`

	nodeType
	client *client
}

// Cleanup makes cleanup.
func (HDFS) Cleanup() {}

func (h HDFS) createClient() (*client, error) {
	httpClient, err := web.NewHTTPClient(h.Client)
	if err != nil {
		return nil, err
	}

	return newClient(httpClient, h.Request), nil
}

func (h *HDFS) determineNodeType() (nodeType, error) {
	var raw rawJMX
	err := h.client.doOKWithDecodeJSON(&raw)
	if err != nil {
		return -1, err
	}

	if raw.isEmpty() {
		return -1, errors.New("empty response")
	}

	jvm := raw.findJvm()
	if jvm == nil {
		return -1, errors.New("couldn't find jvm in response")
	}

	v, ok := jvm["tag.ProcessName"]
	if !ok {
		return -1, errors.New("couldn't find process name in response")
	}

	switch {
	default:
		return unknownNodeType, nil
	case bytes.Equal(v, []byte("\"NameNode\"")):
		return nameNodeType, nil
	case bytes.Equal(v, []byte("\"DataNode\"")):
		return dataNodeType, nil
	}
}

// Init makes initialization.
func (h *HDFS) Init() bool {
	cl, err := h.createClient()
	if err != nil {
		h.Errorf("error on creating client : %v", err)
		return false
	}
	h.client = cl

	t, err := h.determineNodeType()
	if err != nil {
		h.Errorf("error on node type determination : %v", err)
		return false
	}
	h.nodeType = t

	return true
}

// Check makes check.
func (h HDFS) Check() bool {
	return len(h.Collect()) > 0
}

// Charts returns Charts.
func (h HDFS) Charts() *Charts {
	switch h.nodeType {
	default:
		return unknownNodeCharts.Copy()
	case nameNodeType:
		return nameNodeCharts.Copy()
	case dataNodeType:
		return dataNodeCharts.Copy()
	}
}

// Collect collects metrics.
func (h *HDFS) Collect() map[string]int64 {
	mx, err := h.collect()

	if err != nil {
		h.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}

	return mx
}

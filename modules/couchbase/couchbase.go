package couchbase

import (
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	module.Register("couchbase", module.Creator{
		Defaults: module.Defaults{
			UpdateEvery: 10,
		},
		Create: func() module.Module { return New() },
	})
}

func New() *Couchbase {
	return &Couchbase{
		Config: Config{
			HTTP: web.HTTP{
				Request: web.Request{
					URL: "http://127.0.0.1:8091",
				},
				Client: web.Client{
					Timeout: web.Duration{Duration: time.Second * 5},
				},
			},
		},
		client: newCouchbaseClient(),
	}
}

type (
	Config struct {
		web.HTTP `yaml:",inline"`
	}

	Couchbase struct {
		module.Base
		Config `yaml:",inline"`

		client      cbClient
		conn        cbConnection
		charts      *module.Charts
		bucketNames []string
	}
)

func (cb *Couchbase) Cleanup() {
	if cb.conn != nil {
		cb.conn = nil
	}
}

func (cb *Couchbase) Init() bool {
	bucketNames, err := cb.collectBucketNames()
	if err != nil {
		cb.Errorf("init bucketNames: %v", err)
		return false
	}
	cb.bucketNames = bucketNames

	charts, err := cb.initCharts()
	if err != nil {
		cb.Errorf("init charts: %v", err)
		return false
	}

	cb.charts = charts
	return true
}

func (cb *Couchbase) Check() bool {
	return len(cb.Collect()) > 0
}

func (cb *Couchbase) Charts() *Charts {
	return cb.charts
}

func (cb *Couchbase) Collect() map[string]int64 {
	mx, err := cb.collect()
	if err != nil {
		cb.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

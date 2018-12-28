package solr

import (
	"github.com/netdata/go.d.plugin/modules"
)

func init() {
	creator := modules.Creator{
		Create: func() modules.Module { return New() },
	}

	modules.Register("solr", creator)
}

// New creates Solr with default values
func New() *Solr {
	return &Solr{}
}

// Solr solr module
type Solr struct {
	modules.Base // should be embedded by every module
}

// Cleanup makes cleanup
func (Solr) Cleanup() {}

// Init makes initialization
func (Solr) Init() bool {
	return false
}

// Check makes check
func (Solr) Check() bool {
	return false
}

// Charts creates Charts
func (Solr) Charts() *Charts {
	return nil
}

// Collect collects metrics
func (Solr) Collect() map[string]int64 {
	return nil
}

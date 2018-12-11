package apache

import "github.com/netdata/go.d.plugin/modules"

func init() {
	creator := modules.Creator{
		Create: func() modules.Module { return New() },
	}

	modules.Register("apache", creator)
}

// New creates Apache with default values
func New() *Apache {
	return &Apache{}
}

// Apache apache module
type Apache struct {
	modules.Base // should be embedded by every module
}

func (Apache) Cleanup() {

}

func (Apache) Init() bool {
	return false
}

func (Apache) Check() bool {
	return false
}

func (Apache) Charts() *modules.Charts {
	return nil
}

func (Apache) GatherMetrics() map[string]int64 {
	return nil
}

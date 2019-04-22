package apache

import (
	"fmt"

	"github.com/netdata/go.d.plugin/pkg/stm"
)

func (a *Apache) collect() (map[string]int64, error) {
	status, err := a.apiClient.getServerStatus()

	if err != nil {
		return nil, err
	}

	mx := stm.ToMap(status)

	if len(mx) == 0 {
		return nil, fmt.Errorf("nothing was collected from %s", a.URL)
	}

	return mx, nil
}

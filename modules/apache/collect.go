package apache

import "github.com/netdata/go.d.plugin/pkg/stm"

func (a *Apache) collect() (map[string]int64, error) {
	status, err := a.apiClient.getServerStatus()

	if err != nil {
		return nil, err
	}

	return stm.ToMap(status), nil
}

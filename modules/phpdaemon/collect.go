package phpdaemon

import "github.com/netdata/go.d.plugin/pkg/stm"

func (p *PHPDaemon) collect() (map[string]int64, error) {
	s, err := p.client.queryFullStatus()

	if err != nil {
		return nil, err
	}

	// https://github.com/kakserpom/phpdaemon/blob/master/PHPDaemon/Core/Daemon.php
	// see getStateOfWorkers()
	s.Initialized = s.Idle - (s.Init + s.Preinit)
	s.Total = s.Alive + s.Shutdown

	return stm.ToMap(s), nil
}

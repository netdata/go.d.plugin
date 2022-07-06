// SPDX-License-Identifier: GPL-3.0-or-later

package rabbitmq

import (
	"fmt"
	"sync"

	"github.com/netdata/go.d.plugin/pkg/stm"
)

type metrics struct {
	overview *overviewStats
	node     *nodeStats
	vhosts   vhostsStats
}

func (m metrics) hasOverviewStats() bool {
	return m.overview != nil
}

func (m metrics) hasNodeStats() bool {
	return m.node != nil
}

func (m metrics) hasVhostsStats() bool {
	return len(m.vhosts) > 0
}

func (m metrics) hasStats() bool {
	return m.hasOverviewStats() || m.hasNodeStats() || m.hasVhostsStats()
}

func (r *RabbitMQ) collect() (map[string]int64, error) {
	mx := r.scrape(true)

	if !mx.hasStats() {
		return nil, nil
	}

	r.updateCharts(mx)

	return r.parse(mx), nil
}

func (r *RabbitMQ) scrape(doConcurrently bool) *metrics {
	type task func(*metrics)

	var tasks = []task{
		r.scrapeOverviewStats,
		r.scrapeNodeStats,
		r.scrapeVhostsStats,
	}

	wg := &sync.WaitGroup{}
	wrap := func(call task) task {
		return func(metrics *metrics) {
			call(metrics)
			wg.Done()
		}
	}

	var mx metrics
	for _, task := range tasks {
		if doConcurrently {
			wg.Add(1)
			task = wrap(task)
			go task(&mx)
		} else {
			task(&mx)
		}
	}
	wg.Wait()

	return &mx
}

func (r *RabbitMQ) scrapeOverviewStats(mx *metrics) {
	v, err := r.client.scrapeOverview()
	if err != nil {
		r.Error(err)
		return
	}
	mx.overview = v
}

func (r *RabbitMQ) scrapeNodeStats(mx *metrics) {
	v, err := r.client.scrapeNodeStats()
	if err != nil {
		r.Error(err)
		return
	}
	mx.node = v
}

func (r *RabbitMQ) scrapeVhostsStats(mx *metrics) {
	v, err := r.client.scrapeVhostsStats()
	if err != nil {
		r.Error(err)
		return
	}
	mx.vhosts = v
}

func (r *RabbitMQ) parse(mx *metrics) map[string]int64 {
	ms := make(map[string]int64)

	if mx.hasOverviewStats() {
		for k, v := range stm.ToMap(mx.overview) {
			ms[k] = v
		}
	}
	if mx.hasNodeStats() {
		for k, v := range stm.ToMap(mx.node) {
			ms[k] = v
		}
	}
	if mx.hasVhostsStats() {
		for _, vhost := range mx.vhosts {
			for k, v := range stm.ToMap(vhost) {
				ms[fmt.Sprintf("vhost_%s_%s", vhost.Name, k)] = v
			}
		}
	}

	return ms
}

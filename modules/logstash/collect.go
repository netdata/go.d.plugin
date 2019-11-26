package logstash

import "github.com/netdata/go.d.plugin/pkg/stm"

func (l *Logstash) collect() (map[string]int64, error) {
	jvmStats, err := l.apiClient.jvmStats()
	if err != nil {
		return nil, err
	}

	for id := range jvmStats.Pipelines {
		chartID := "pipeline_" + id + "_event"
		if !l.Charts().Has(chartID) {
			if err := l.Charts().Add(createPipelineChart(id)...); err != nil {
				l.Warningf("create charts for '%s' pipeline: %v", id, err)
			}
		}
	}

	return stm.ToMap(jvmStats), nil
}

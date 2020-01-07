package logstash

import "github.com/netdata/go.d.plugin/pkg/stm"

func (l *Logstash) collect() (map[string]int64, error) {
	stats, err := l.client.jvmStats()
	if err != nil {
		return nil, err
	}

	l.collectPipelines(stats.Pipelines)
	return stm.ToMap(stats), nil
}

func (l *Logstash) collectPipelines(pipelines map[string]pipeline) {
	if len(pipelines) == 0 {
		return
	}

	set := make(map[string]bool)
	for id := range pipelines {
		set[id] = true
		if !l.collectedPipelines[id] {
			l.collectedPipelines[id] = true
			l.addPipelineCharts(id)
		}
	}

	for id := range l.collectedPipelines {
		if !set[id] {
			delete(l.collectedPipelines, id)
			l.removePipelineCharts(id)
		}
	}
}

func (l *Logstash) addPipelineCharts(id string) {
	err := l.Charts().Add(*pipelineCharts(id)...)
	if err != nil {
		l.Warningf("can't add pipeline '%s' charts: %v", id, err)
	}
}

func (l *Logstash) removePipelineCharts(id string) {
	for _, chart := range *pipelineCharts(id) {
		chart = l.Charts().Get(chart.ID)
		if chart == nil {
			l.Warningf("can't remove pipeline '%s' charts: chart is not found", id)
			continue
		}
		chart.MarkRemove()
		chart.MarkNotCreated()
	}
}

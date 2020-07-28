package prometheus

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

type split interface {
	chartID(pm prometheus.Metric) string
	dimID(pm prometheus.Metric) string
	dimName(pm prometheus.Metric) string
	doReSplit(pms prometheus.Metrics) bool
}

type anySplit struct {
	chartIDFunc   func(pm prometheus.Metric) string
	dimIDFunc     func(pm prometheus.Metric) string
	dimNameFunc   func(pm prometheus.Metric) string
	doReSplitFunc func(pms prometheus.Metrics) bool
}

func (s anySplit) chartID(pm prometheus.Metric) string   { return s.chartIDFunc(pm) }
func (s anySplit) dimID(pm prometheus.Metric) string     { return s.dimIDFunc(pm) }
func (s anySplit) dimName(pm prometheus.Metric) string   { return s.dimNameFunc(pm) }
func (s anySplit) doReSplit(pms prometheus.Metrics) bool { return s.doReSplitFunc(pms) }

type histogramSplit struct{}

func (s histogramSplit) chartID(pm prometheus.Metric) string { return joinLabelsExcept(pm, "le") }
func (s histogramSplit) dimID(pm prometheus.Metric) string   { return joinLabels(pm) }
func (s histogramSplit) dimName(pm prometheus.Metric) string { return pm.Labels.Get("le") }
func (s histogramSplit) doReSplit(_ prometheus.Metrics) bool { return false }

type summarySplit struct{}

func (s summarySplit) chartID(pm prometheus.Metric) string { return joinLabelsExcept(pm, "quantile") }
func (s summarySplit) dimID(pm prometheus.Metric) string   { return joinLabels(pm) }
func (s summarySplit) dimName(pm prometheus.Metric) string { return pm.Labels.Get("quantile") }
func (s summarySplit) doReSplit(_ prometheus.Metrics) bool { return false }

func newAnySplit(pms prometheus.Metrics) (*anySplit, error) {
	if s := newAnySplitSpecialCase(pms); s != nil {
		return s, nil
	}
	return newAnySplitGrouped(pms)
}

func newAnySplitSpecialCase(pms prometheus.Metrics) *anySplit {
	pm := pms[0]
	// name + special + something
	if pm.Labels.Len() < 3 {
		return nil
	}
	if !strings.Contains(pm.Name(), "_cpu_") {
		return nil
	}

	if pm.Labels.Has("cpu") {
		return &anySplit{
			chartIDFunc:   func(pm prometheus.Metric) string { return joinLabelsOnly(pm, "__name__", "cpu") },
			dimIDFunc:     func(pm prometheus.Metric) string { return joinLabels(pm) },
			dimNameFunc:   func(pm prometheus.Metric) string { return joinLabelsExcept(pm, "__name__", "cpu") },
			doReSplitFunc: func(pms prometheus.Metrics) bool { return false },
		}
	}
	if pm.Labels.Has("core") {
		return &anySplit{
			chartIDFunc:   func(pm prometheus.Metric) string { return joinLabelsOnly(pm, "__name__", "core") },
			dimIDFunc:     func(pm prometheus.Metric) string { return joinLabels(pm) },
			dimNameFunc:   func(pm prometheus.Metric) string { return joinLabelsExcept(pm, "__name__", "core") },
			doReSplitFunc: func(pms prometheus.Metrics) bool { return false },
		}
	}
	return nil
}

const (
	maxChartsPerMetric = 20
	desiredDim         = 20
	maxDim             = desiredDim + 10
)

func newAnySplitGrouped(pms prometheus.Metrics) (*anySplit, error) {
	numOfCharts := desiredNumOfCharts(len(pms))
	if numOfCharts > maxChartsPerMetric {
		return nil, fmt.Errorf("to many charts, got %d charts, max %d", numOfCharts, maxChartsPerMetric)
	}

	s := &anySplit{
		dimIDFunc:     func(pm prometheus.Metric) string { return joinLabels(pm) },
		dimNameFunc:   func(pm prometheus.Metric) string { return joinLabelsExcept(pm, "__name__") },
		doReSplitFunc: func(pms prometheus.Metrics) bool { return maxNumOfCharts(len(pms)) > numOfCharts },
	}

	if numOfCharts == 1 {
		s.chartIDFunc = func(pm prometheus.Metric) string { return pm.Name() }
	} else {
		var current uint64
		cache := make(map[string]uint64)
		s.chartIDFunc = func(pm prometheus.Metric) string {
			id := joinLabels(pm)
			if group, ok := cache[id]; ok {
				return pm.Name() + "_group" + strconv.FormatUint(group, 10)
			}
			if current >= numOfCharts {
				current = 0
			}
			current++
			cache[id] = current - 1
			return pm.Name() + "_group" + strconv.FormatUint(current-1, 10)
		}
	}

	return s, nil
}

func newHistogramSplit() *histogramSplit {
	return &histogramSplit{}
}

func newSummarySplit() *summarySplit {
	return &summarySplit{}
}

func desiredNumOfCharts(numOfSeries int) (num uint64) {
	num = uint64(numOfSeries / desiredDim)
	if numOfSeries%desiredDim != 0 {
		num++
	}
	return num
}

func maxNumOfCharts(numOfSeries int) (num uint64) {
	num = uint64(numOfSeries / maxDim)
	if numOfSeries%maxDim != 0 {
		num++
	}
	return num
}

func joinLabels(pm prometheus.Metric) string {
	return joinLabelsIf(pm, false)
}

func joinLabelsOnly(pm prometheus.Metric, only string, other ...string) string {
	return joinLabelsIf(pm, true, append(other, only)...)
}

func joinLabelsExcept(pm prometheus.Metric, except string, other ...string) string {
	return joinLabelsIf(pm, false, append(other, except)...)
}

func joinLabelsIf(pm prometheus.Metric, shouldContain bool, labels ...string) string {
	// {__name__="name",value1="value1",value1="value2"} => name|value1=value1,value2=value2
	var id strings.Builder
	var comma bool
loop:
	for i, label := range pm.Labels {
		if len(labels) > 0 {
			if ok := contains(label.Name, labels); (!ok && shouldContain) || (ok && !shouldContain) {
				continue loop
			}
		}

		if i == 0 {
			id.WriteString(label.Value)
			continue
		}
		if id.Len() > 0 {
			if !comma {
				id.WriteString("|")
			} else {
				id.WriteString(",")
			}
		}
		comma = true
		id.WriteString(label.Name)
		id.WriteString("=")
		id.WriteString(label.Value)
	}
	return id.String()
}

func contains(value string, in []string) bool {
	switch len(in) {
	case 0:
		return false
	case 1:
		return value == in[0]
	default:
		return value == in[0] || contains(value, in[1:])
	}
}

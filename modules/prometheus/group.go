// SPDX-License-Identifier: GPL-3.0-or-later

package prometheus

import (
	"strconv"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/prometheus/prometheus/model/labels"
)

type grouper interface {
	chartID(pm prometheus.Metric) string
	dimID(pm prometheus.Metric) string
	dimName(pm prometheus.Metric) string
}

type anyGrouper struct {
	chartIDFunc func(pm prometheus.Metric) string
	dimIDFunc   func(pm prometheus.Metric) string
	dimNameFunc func(pm prometheus.Metric) string
}

func (g anyGrouper) chartID(pm prometheus.Metric) string { return g.chartIDFunc(pm) }
func (g anyGrouper) dimID(pm prometheus.Metric) string   { return g.dimIDFunc(pm) }
func (g anyGrouper) dimName(pm prometheus.Metric) string { return g.dimNameFunc(pm) }

var (
	defaultAnyGrouping = anyGrouper{
		chartIDFunc: func(pm prometheus.Metric) string { return pm.Name() },
		dimIDFunc:   func(pm prometheus.Metric) string { return joinLabels(pm) },
		dimNameFunc: func(pm prometheus.Metric) string { return joinLabelsExcept(pm, labels.MetricName) },
	}
	defaultHistogramGrouping = anyGrouper{
		chartIDFunc: func(pm prometheus.Metric) string { return joinLabelsExcept(pm, "le") },
		dimIDFunc:   func(pm prometheus.Metric) string { return joinLabels(pm) },
		dimNameFunc: func(pm prometheus.Metric) string { return pm.Labels.Get("le") },
	}
	defaultSummaryGrouping = anyGrouper{
		chartIDFunc: func(pm prometheus.Metric) string { return joinLabelsExcept(pm, "quantile") },
		dimIDFunc:   func(pm prometheus.Metric) string { return joinLabels(pm) },
		dimNameFunc: func(pm prometheus.Metric) string { return pm.Labels.Get("quantile") },
	}
)

func newGroupingSplitN(grp grouper, numOfGroups uint64) grouper {
	if numOfGroups <= 1 {
		return grp
	}

	var current uint64
	cache := make(map[uint64]uint64)
	return anyGrouper{
		chartIDFunc: func(pm prometheus.Metric) string {
			hash := pm.Labels.Hash()
			if id, ok := cache[hash]; ok {
				return grp.chartID(pm) + "_group" + strconv.FormatUint(id, 10)
			}
			if current >= numOfGroups {
				current = 0
			}
			id := current
			current++
			cache[hash] = id
			return grp.chartID(pm) + "_group" + strconv.FormatUint(id, 10)
		},
		dimIDFunc:   grp.dimID,
		dimNameFunc: grp.dimName,
	}
}

func newGroupingGroupedBy(lbs ...string) grouper {
	return anyGrouper{
		chartIDFunc: func(pm prometheus.Metric) string {
			return joinLabelsOnly(pm, labels.MetricName, lbs...)
		},
		dimIDFunc: func(pm prometheus.Metric) string {
			return joinLabels(pm)
		},
		dimNameFunc: func(pm prometheus.Metric) string {
			return joinLabelsExcept(pm, labels.MetricName, lbs...)
		},
	}
}

func joinLabels(pm prometheus.Metric) string {
	return joinLabelsIf(pm, false)
}

func joinLabelsOnly(pm prometheus.Metric, label string, otherLabels ...string) string {
	return joinLabelsIf(pm, true, append(otherLabels, label)...)
}

func joinLabelsExcept(pm prometheus.Metric, label string, otherLabels ...string) string {
	return joinLabelsIf(pm, false, append(otherLabels, label)...)
}

func joinLabelsIf(pm prometheus.Metric, shouldContain bool, lbs ...string) string {
	// {__name__="name",value1="value1",value1="value2"} => name|value1=value1,value2=value2
	var id strings.Builder
	var comma bool
loop:
	for i, label := range pm.Labels {
		if len(lbs) > 0 {
			if ok := contains(lbs, label.Name); (!ok && shouldContain) || (ok && !shouldContain) {
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

		lval := label.Value
		if strings.IndexByte(lval, ' ') != -1 {
			lval = spaceReplacer.Replace(lval)
		}
		if strings.IndexByte(lval, '\\') != -1 {
			if lval = decodeLabelValue(lval); strings.IndexByte(lval, '\\') != -1 {
				lval = backslashReplacer.Replace(lval)
			}
		}
		id.WriteString(lval)
	}
	return id.String()
}

func decodeLabelValue(value string) string {
	v, err := strconv.Unquote("\"" + value + "\"")
	if err != nil {
		return value
	}
	return v
}

func contains(in []string, value string) bool {
	switch len(in) {
	case 0:
		return false
	case 1:
		return value == in[0]
	default:
		return value == in[0] || contains(in[1:], value)
	}
}

// TODO: fix
var (
	spaceReplacer     = strings.NewReplacer(" ", "_")
	backslashReplacer = strings.NewReplacer(`\`, "_")
)

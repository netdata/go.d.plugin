// SPDX-License-Identifier: GPL-3.0-or-later

package prometheus

import (
	"fmt"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/prometheus"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/prometheus/prometheus/model/textparse"
)

type (
	Charts = module.Charts
	Dims   = module.Dims
)

var statsCharts = Charts{
	{
		ID:    "collect_statistics",
		Title: "Collect Statistics",
		Units: "num",
		Fam:   "collection",
		Ctx:   "prometheus.collect_statistics",
		Dims: Dims{
			{ID: "series"},
			{ID: "metrics"},
			{ID: "charts"},
		},
	},
}

func anyChart(id, app string, pm prometheus.Metric, meta prometheus.Metadata) *module.Chart {
	units := extractUnits(pm.Name())
	if isIncremental(pm, meta) && !isIncrementalUnitsException(units) {
		units += "/s"
	}
	cType := module.Line
	if strings.HasPrefix(units, "bytes") {
		cType = module.Area
	}
	return &module.Chart{
		ID:    id,
		Title: chartTitle(pm, meta),
		Units: units,
		Fam:   chartFamily(pm),
		Ctx:   chartContext(app, pm),
		Type:  cType,
	}
}

func isIncrementalUnitsException(units string) bool {
	switch units {
	case "seconds", "time":
		return true
	}
	return false
}

func summaryChart(id, app string, pm prometheus.Metric, meta prometheus.Metadata) *module.Chart {
	return &module.Chart{
		ID:    id,
		Title: chartTitle(pm, meta),
		Units: "observations",
		Fam:   chartFamily(pm),
		Ctx:   chartContext(app, pm),
		Type:  module.Stacked,
	}
}

func histogramChart(id, app string, pm prometheus.Metric, meta prometheus.Metadata) *module.Chart {
	return summaryChart(id, app, pm, meta)
}

func chartTitle(pm prometheus.Metric, meta prometheus.Metadata) string {
	if help := meta.Help(pm.Name()); help != "" {
		// ' used to wrap external plugins api messages, netdata parser cant handle ' inside ''
		return strings.Replace(help, "'", "\"", -1)
	}
	return fmt.Sprintf("Metric \"%s\"", pm.Name())
}

func chartFamily(pm prometheus.Metric) (fam string) {
	if strings.HasPrefix(pm.Name(), "go_") {
		return "go"
	}
	if strings.HasPrefix(pm.Name(), "process_") {
		return "process"
	}
	if parts := strings.SplitN(pm.Name(), "_", 3); len(parts) < 3 {
		fam = pm.Name()
	} else {
		fam = parts[0] + "_" + parts[1]
	}

	// remove number suffix if any
	// load1, load5, load15 => load
	i := len(fam) - 1
	for i >= 0 && fam[i] >= '0' && fam[i] <= '9' {
		i--
	}
	if i > 0 {
		return fam[:i+1]
	}
	return fam
}

func chartContext(app string, pm prometheus.Metric) string {
	if app == "" {
		return fmt.Sprintf("prometheus.%s", pm.Name())
	}
	return fmt.Sprintf("prometheus.%s.%s", app, pm.Name())
}

func anyChartDimension(id, name string, pm prometheus.Metric, meta prometheus.Metadata) *module.Dim {
	algorithm := module.Absolute
	if isIncremental(pm, meta) {
		algorithm = module.Incremental
	}
	return &module.Dim{
		ID:   id,
		Name: name,
		Algo: algorithm,
		Div:  precision,
	}
}

func summaryChartDimension(id, name string) *module.Dim {
	return &module.Dim{
		ID:   id,
		Name: name,
		Algo: module.Incremental,
		Div:  precision,
	}
}

func histogramChartDim(id, name string) *module.Dim {
	return &module.Dim{
		ID:   id,
		Name: name,
		Algo: module.Incremental,
		Div:  precision,
	}
}

func extractUnits(metric string) string {
	// https://prometheus.io/docs/practices/naming/#metric-names
	// ...must have a single unit (i.e. do not mix seconds with milliseconds, or seconds with bytes).
	// ...should have a suffix describing the unit, in plural form.
	// Note that an accumulating count has total as a suffix, in addition to the unit if applicable

	idx := strings.LastIndexByte(metric, '_')
	if idx == -1 {
		return "events"
	}
	switch suffix := metric[idx:]; suffix {
	case "_total", "_sum", "_count":
		return extractUnits(metric[:idx])
	}
	switch units := metric[idx+1:]; units {
	case "hertz":
		return "Hz"
	default:
		return units
	}
}

func isIncremental(pm prometheus.Metric, meta prometheus.Metadata) bool {
	switch meta.Type(pm.Name()) {
	case textparse.MetricTypeCounter,
		textparse.MetricTypeHistogram,
		textparse.MetricTypeSummary:
		return true
	}
	switch {
	case strings.HasSuffix(pm.Name(), "_total"),
		strings.HasSuffix(pm.Name(), "_sum"),
		strings.HasSuffix(pm.Name(), "_count"):
		return true
	}
	return false
}

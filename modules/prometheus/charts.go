package prometheus

import (
	"fmt"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/prometheus"

	"github.com/netdata/go-orchestrator/module"
)

type (
	Charts = module.Charts
	Dims   = module.Dims
)

var statsCharts = Charts{
	{
		ID:    "collections_statistics",
		Title: "Statistics",
		Units: "num",
		Fam:   "collection",
		Ctx:   "prometheus.collections_statistics",
		Dims: Dims{
			{ID: "series", Name: "series"},
			{ID: "metrics", Name: "metrics"},
			{ID: "charts", Name: "charts"},
		},
	},
}

func anyChart(id string, pm prometheus.Metric, meta prometheus.Metadata) *module.Chart {
	units := extractUnits(pm.Name())
	if isIncremental(pm) {
		units += "/s"
	}
	return &module.Chart{
		ID:    id,
		Title: chartTitle(pm, meta),
		Units: units,
		Fam:   pm.Name(),
		Ctx:   "prometheus." + pm.Name(),
		Type:  module.Line,
	}
}

func summaryChart(id string, pm prometheus.Metric, meta prometheus.Metadata) *module.Chart {
	return &module.Chart{
		ID:    id,
		Title: chartTitle(pm, meta),
		Units: "observations",
		Fam:   pm.Name(),
		Ctx:   "prometheus." + pm.Name(),
		Type:  module.Stacked,
	}
}

func summaryPercentChart(id string, pm prometheus.Metric, meta prometheus.Metadata) *module.Chart {
	return &module.Chart{
		ID:    id,
		Title: chartTitle(pm, meta),
		Units: "%",
		Fam:   pm.Name(),
		Ctx:   "prometheus." + pm.Name() + "_percentage",
		Type:  module.Stacked,
	}
}

func histogramChart(id string, pm prometheus.Metric, meta prometheus.Metadata) *module.Chart {
	return summaryChart(id, pm, meta)
}

func histogramPercentChart(id string, pm prometheus.Metric, meta prometheus.Metadata) *module.Chart {
	return summaryPercentChart(id, pm, meta)
}

func chartTitle(pm prometheus.Metric, meta prometheus.Metadata) string {
	if v := meta.Help(pm.Name()); v != "" {
		// ' used to wrap external plugins api messages, netdata parser cant handle ' inside ''
		return strings.Replace(v, "'", "\"", -1)
	}
	return fmt.Sprintf("Metric \"%s\"", pm.Name())
}

func anyDimension(id string, pm prometheus.Metric) *module.Dim {
	name := id
	// name => name
	// name|value1=value1|value2=value2 => value1=value1|value2=value2
	if name != pm.Name() && strings.HasPrefix(name, pm.Name()) {
		name = id[len(pm.Name())+1:]
	}
	algorithm := module.Absolute
	if isIncremental(pm) {
		algorithm = module.Incremental
	}
	return &module.Dim{
		ID:   id,
		Name: name,
		Algo: algorithm,
		Div:  precision,
	}
}

func summaryChartDimension(id string, pm prometheus.Metric) *module.Dim {
	return &module.Dim{
		ID:   id,
		Name: pm.Labels.Get("quantile"),
		Algo: module.Incremental,
		Div:  precision,
	}
}

func summaryPercentChartDim(id string, pm prometheus.Metric) *module.Dim {
	return &module.Dim{
		ID:   id,
		Name: pm.Labels.Get("quantile"),
		Algo: module.PercentOfIncremental,
		Div:  precision,
	}
}

func histogramChartDim(id string, pm prometheus.Metric) *module.Dim {
	return &module.Dim{
		ID:   id,
		Name: pm.Labels.Get("le"),
		Algo: module.Incremental,
		Div:  precision,
	}
}

func histogramPercentChartDim(id string, pm prometheus.Metric) *module.Dim {
	return &module.Dim{
		ID:   id,
		Name: pm.Labels.Get("le"),
		Algo: module.PercentOfIncremental,
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
	return metric[idx+1:]
}

func isIncremental(pm prometheus.Metric) bool {
	return strings.HasSuffix(pm.Name(), "_total")
}

package prometheus

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/prometheus"

	"github.com/prometheus/prometheus/model/labels"
)

const (
	prioDefault   = module.Priority
	prioGORuntime = prioDefault + 10
)

func (p *Prometheus) addGaugeChart(id, name, help string, labels labels.Labels) {
	units := getChartUnits(name)

	cType := module.Line
	if strings.HasSuffix(units, "bytes") {
		cType = module.Area
	}

	chart := &module.Chart{
		ID:       id,
		Title:    getChartTitle(name, help),
		Units:    units,
		Fam:      getChartFamily(name),
		Ctx:      "prometheus." + name,
		Type:     cType,
		Priority: getChartPriority(name),
		Dims: module.Dims{
			{ID: id, Name: name, Div: precision},
		},
	}

	for _, lbl := range labels {
		chart.Labels = append(chart.Labels,
			module.Label{Key: lbl.Name, Value: lbl.Value},
		)
	}

	if err := p.Charts().Add(chart); err != nil {
		p.Warning(err)
	}
}

func (p *Prometheus) addCounterChart(id, name, help string, labels labels.Labels) {
	units := getChartUnits(name)

	switch units {
	case "seconds", "time":
	default:
		units += "/s"
	}

	cType := module.Line
	if strings.HasSuffix(units, "bytes/s") {
		cType = module.Area
	}

	chart := &module.Chart{
		ID:       id,
		Title:    getChartTitle(name, help),
		Units:    units,
		Fam:      getChartFamily(name),
		Ctx:      "prometheus." + name,
		Type:     cType,
		Priority: getChartPriority(name),
		Dims: module.Dims{
			{ID: id, Name: name, Algo: module.Incremental, Div: precision},
		},
	}
	for _, lbl := range labels {
		chart.Labels = append(chart.Labels,
			module.Label{Key: lbl.Name, Value: lbl.Value},
		)
	}

	if err := p.Charts().Add(chart); err != nil {
		p.Warning(err)
	}
}

func (p *Prometheus) addSummaryChart(id, name, help string, labels labels.Labels, quantiles []prometheus.Quantile) {
	units := getChartUnits(name)

	switch units {
	case "seconds", "time":
	default:
		units += "/s"
	}

	chart := &module.Chart{
		ID:       id,
		Title:    getChartTitle(name, help),
		Units:    units,
		Fam:      getChartFamily(name),
		Ctx:      "prometheus." + name,
		Priority: getChartPriority(name),
	}
	for _, v := range quantiles {
		s := strconv.FormatFloat(v.Quantile(), 'f', -1, 64)
		chart.Dims = append(chart.Dims, &module.Dim{
			ID:   fmt.Sprintf("%s_quantile=%s", id, s),
			Name: fmt.Sprintf("quantile_%s", s),
			Algo: module.Incremental,
			Div:  precision,
		})
	}
	for _, lbl := range labels {
		chart.Labels = append(chart.Labels, module.Label{
			Key:   lbl.Name,
			Value: lbl.Value,
		})
	}

	if err := p.Charts().Add(chart); err != nil {
		p.Warning(err)
	}
}

func (p *Prometheus) addHistogramChart(id, name, help string, labels labels.Labels, buckets []prometheus.Bucket) {
	chart := &module.Chart{
		ID:       id,
		Title:    getChartTitle(name, help),
		Units:    "observations/s",
		Fam:      getChartFamily(name),
		Ctx:      "prometheus." + name,
		Priority: getChartPriority(name),
	}
	for _, v := range buckets {
		s := strconv.FormatFloat(v.UpperBound(), 'f', -1, 64)
		chart.Dims = append(chart.Dims, &module.Dim{
			ID:   fmt.Sprintf("%s_bucket=%s", id, s),
			Name: fmt.Sprintf("bucket_%s", s),
			Algo: module.Incremental,
		})
	}
	for _, lbl := range labels {
		chart.Labels = append(chart.Labels, module.Label{
			Key:   lbl.Name,
			Value: lbl.Value,
		})
	}

	if err := p.Charts().Add(chart); err != nil {
		p.Warning(err)
	}
}

func getChartTitle(name, help string) string {
	if help == "" {
		return fmt.Sprintf("SeriesSample \"%s\"", name)
	}

	help = strings.Replace(help, "'", "", -1)
	if strings.HasSuffix(help, ".") {
		help = help[:len(help)-1]
	}
	return help
}

func getChartFamily(metric string) (fam string) {
	if strings.HasPrefix(metric, "go_") {
		return "go"
	}
	if strings.HasPrefix(metric, "process_") {
		return "process"
	}
	if parts := strings.SplitN(metric, "_", 3); len(parts) < 3 {
		fam = metric
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

func getChartUnits(metric string) string {
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
		return getChartUnits(metric[:idx])
	}
	switch units := metric[idx+1:]; units {
	case "hertz":
		return "Hz"
	default:
		return units
	}
}

func getChartPriority(name string) int {
	if strings.HasPrefix(name, "go_") || strings.HasPrefix(name, "process_") {
		return prioGORuntime
	}
	return prioDefault
}

package prometheus

import (
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/prometheus"

	"github.com/prometheus/prometheus/pkg/textparse"
)

const (
	precision          = 1000
	maxChartsPerMetric = 20
	desiredDim         = 10
	maxDim             = desiredDim + 10
)

// TODO: proper cleanup stale charts
func (p *Prometheus) collect() (map[string]int64, error) {
	pms, err := p.prom.Scrape()
	if err != nil {
		return nil, err
	}
	if len(pms) == 0 {
		p.Warningf("endpoint '%s' returned 0 time series", p.UserURL)
		return nil, nil
	}

	names, metricSet := buildMetricSet(pms)
	meta := p.prom.Metadata()
	mx := make(map[string]int64)

	for _, name := range names {
		metrics := metricSet[name]
		if len(metrics) == 0 || p.skipMetrics[name] {
			continue
		}

		switch meta.Type(name) {
		case textparse.MetricTypeGauge, textparse.MetricTypeCounter:
			p.collectAny(mx, metrics, meta)
		case textparse.MetricTypeSummary:
			p.collectSummary(mx, metrics, meta)
		case textparse.MetricTypeHistogram:
			p.collectHistogram(mx, metrics, meta)
		case textparse.MetricTypeUnknown:
			pm := metrics[0]
			switch {
			case pm.Labels.Get("quantile") != "":
				p.collectSummary(mx, metrics, meta)
			case pm.Labels.Get("le") != "":
				p.collectHistogram(mx, metrics, meta)
			default:
				p.collectAny(mx, metrics, meta)
			}
		}
	}
	p.Debugf("time series: %d, metrics: %d, charts: %d", len(pms), len(names), len(*p.Charts()))
	mx["series"] = int64(len(pms))
	mx["metrics"] = int64(len(names))
	mx["charts"] = int64(len(*p.Charts()))
	return mx, nil
}

func (p *Prometheus) collectAny(mx map[string]int64, pms prometheus.Metrics, meta prometheus.Metadata) {
	name := pms[0].Name()
	if !p.cache.has(name) {
		num := desiredNumOfCharts(len(pms))
		if num > maxChartsPerMetric {
			p.skipMetrics[name] = true
			p.Infof("skip metric '%s', it would produce %d charts (max %d)", name, num, maxChartsPerMetric)
			return
		}

		p.cache.put(name, &cacheEntry{
			split:         anySplitFunc(pms),
			desiredCharts: num,
			charts:        make(chartsCache),
			dims:          make(dimsCache),
		})
	}

	cache := p.cache.get(name)

	if maxNumOfCharts(len(pms)) > cache.desiredCharts {
		p.cleanupMetric(name)
		p.collectAny(mx, pms, meta)
	}

	for _, pm := range pms {
		chartID := cache.split(pm)
		dimID := joinLabels(pm)

		mx[dimID] = int64(pm.Value * precision)

		if !cache.hasChart(chartID) {
			chart := anyChart(chartID, pm, meta)
			cache.putChart(chartID, chart)
			_ = p.Charts().Add(chart)
		}
		if !cache.hasDim(dimID, chartID) {
			cache.putDim(dimID, chartID)
			chart := cache.getChart(chartID)
			dim := anyDimension(dimID, pm)
			_ = chart.AddDim(dim)
			chart.MarkNotCreated()
		}
	}
}

func (p *Prometheus) collectSummary(mx map[string]int64, pms prometheus.Metrics, meta prometheus.Metadata) {
	sortSummary(pms)
	name := pms[0].Name()
	if !p.cache.has(name) {
		p.cache.put(name, &cacheEntry{
			split:  summarySplit,
			charts: make(chartsCache),
			dims:   make(dimsCache),
		})
	}

	cache := p.cache.get(name)

	for _, pm := range pms {
		if math.IsNaN(pm.Value) {
			continue
		}

		chartID := cache.split(pm)
		percChartID := chartID + "_percentage"
		dimID := joinLabels(pm)

		mx[dimID] = int64(pm.Value * precision)

		if !cache.hasChart(chartID) {
			chart := summaryChart(chartID, pm, meta)
			cache.putChart(chartID, chart)
			if err := p.Charts().Add(chart); err != nil {
				p.Warning(err)
			}
		}
		if !cache.hasDim(dimID, chartID) {
			cache.putDim(dimID, chartID)
			chart := cache.getChart(chartID)
			dim := summaryChartDimension(dimID, pm)
			if err := chart.AddDim(dim); err != nil {
				p.Warning(err)
			}
			chart.MarkNotCreated()
		}

		if !cache.hasChart(percChartID) {
			chart := summaryPercentChart(percChartID, pm, meta)
			cache.putChart(percChartID, chart)
			if err := p.Charts().Add(chart); err != nil {
				p.Warning(err)
			}
		}
		if !cache.hasDim(dimID, percChartID) {
			cache.putDim(dimID, percChartID)
			chart := cache.getChart(percChartID)
			dim := summaryPercentChartDim(dimID, pm)
			if err := chart.AddDim(dim); err != nil {
				p.Warning(err)
			}
			chart.MarkNotCreated()
		}
	}
}

func (p *Prometheus) collectHistogram(mx map[string]int64, pms prometheus.Metrics, meta prometheus.Metadata) {
	sortHistogram(pms)
	name := pms[0].Name()
	if !p.cache.has(name) {
		p.cache.put(name, &cacheEntry{
			split:  histogramSplit,
			charts: make(chartsCache),
			dims:   make(dimsCache),
		})
	}

	set := make(map[string]float64)
	cache := p.cache.get(name)

	for _, pm := range pms {
		chartID := cache.split(pm)
		percChartID := chartID + "_percentage"
		dimID := joinLabels(pm)

		// {handler="/",le="0.1"} 1
		// {handler="/",le="0.2"} 2
		// {handler="/",le="0.4"} 3
		// le="0.4" = 3 - 2 (le="0.4" - le="0.2")
		if v, ok := set[chartID]; !ok {
			mx[dimID] = int64(pm.Value * precision)
		} else {
			mx[dimID] = int64((pm.Value - v) * precision)
		}
		set[chartID] = pm.Value

		if !cache.hasChart(chartID) {
			chart := histogramChart(chartID, pm, meta)
			cache.putChart(chartID, chart)
			if err := p.Charts().Add(chart); err != nil {
				p.Warning(err)
			}
		}
		if !cache.hasDim(dimID, chartID) {
			cache.putDim(dimID, chartID)
			chart := cache.getChart(chartID)
			dim := histogramChartDim(dimID, pm)
			if err := chart.AddDim(dim); err != nil {
				p.Warning(err)
			}
			chart.MarkNotCreated()
		}

		if !cache.hasChart(percChartID) {
			chart := histogramPercentChart(percChartID, pm, meta)
			cache.putChart(percChartID, chart)
			if err := p.Charts().Add(chart); err != nil {
				p.Warning(err)
			}
		}
		if !cache.hasDim(dimID, percChartID) {
			cache.putDim(dimID, percChartID)
			chart := cache.getChart(percChartID)
			dim := histogramPercentChartDim(dimID, pm)
			if err := chart.AddDim(dim); err != nil {
				p.Warning(err)
			}
			chart.MarkNotCreated()
		}
	}
}

func (p *Prometheus) cleanupMetric(name string) {
	if !p.cache.has(name) {
		return
	}
	defer p.cache.remove(name)

	for _, chart := range p.cache.get(name).charts {
		for _, dim := range chart.Dims {
			_ = chart.MarkDimRemove(dim.ID, true)
		}
		chart.MarkRemove()
		chart.MarkNotCreated()
	}
}

// TODO: should be done by prom pkg
func buildMetricSet(pms prometheus.Metrics) (names []string, metrics map[string]prometheus.Metrics) {
	names = make([]string, 0, len(pms))
	metrics = make(map[string]prometheus.Metrics)

	for _, pm := range pms {
		if _, ok := metrics[pm.Name()]; !ok {
			names = append(names, pm.Name())
		}
		metrics[pm.Name()] = append(metrics[pm.Name()], pm)
	}
	return names, metrics
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

func anySplitFunc(pms prometheus.Metrics) func(metric prometheus.Metric) string {
	num := desiredNumOfCharts(len(pms))
	if num == 1 {
		return func(pm prometheus.Metric) string {
			return pm.Name()
		}
	}

	var current uint64
	cache := make(map[string]uint64)
	return func(pm prometheus.Metric) string {
		str := joinLabels(pm)
		if v, ok := cache[str]; ok {
			return pm.Name() + "_group" + strconv.FormatUint(v, 10)
		}
		if current >= num {
			current = 0
		}
		current++
		cache[str] = current - 1
		return pm.Name() + "_group" + strconv.FormatUint(current-1, 10)
	}
}

func summarySplit(pm prometheus.Metric) string {
	return joinLabels(pm, "quantile")
}

func histogramSplit(pm prometheus.Metric) string {
	return joinLabels(pm, "le")
}

func sortSummary(pms prometheus.Metrics) {
	sort.Slice(pms, func(i, j int) bool {
		iVal, _ := strconv.ParseFloat(pms[i].Labels.Get("quantile"), 64)
		jVal, _ := strconv.ParseFloat(pms[j].Labels.Get("quantile"), 64)
		return iVal < jVal
	})
}

func sortHistogram(pms prometheus.Metrics) {
	sort.Slice(pms, func(i, j int) bool {
		iStr := pms[i].Labels.Get("le")
		jStr := pms[j].Labels.Get("le")
		if iStr == "+Inf" {
			return false
		}
		if jStr == "+Inf" {
			return true
		}
		iVal, _ := strconv.ParseFloat(iStr, 64)
		jVal, _ := strconv.ParseFloat(jStr, 64)
		return iVal < jVal
	})
}

func joinLabels(pm prometheus.Metric, skip ...string) string {
	// {__name__="name",value1="value1",value1="value2"} => name|value1=value1,value2=value2
	var id strings.Builder
	var comma bool
loop:
	for i, label := range pm.Labels {
		if i == 0 {
			id.WriteString(label.Value)
			continue
		}
		for _, name := range skip {
			if label.Name == name {
				continue loop
			}
		}
		if !comma {
			id.WriteString("|")
		} else {
			id.WriteString(",")
		}
		comma = true
		id.WriteString(label.Name)
		id.WriteString("=")
		id.WriteString(label.Value)
	}
	return id.String()
}
